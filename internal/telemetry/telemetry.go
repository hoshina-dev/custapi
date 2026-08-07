// Package telemetry wires up the OpenTelemetry SDK for custapi.
//
// It is designed around running behind the AWS Distro for OpenTelemetry
// (ADOT) Collector as a sidecar container in an ECS task:
//
//   - Traces use the AWS X-Ray ID generator and propagator so trace/span IDs
//     are valid X-Ray trace IDs, and are batched to the collector over
//     OTLP/gRPC. ADOT forwards them to AWS X-Ray.
//   - Metrics are exported periodically over OTLP/gRPC to the same
//     collector, which converts them to CloudWatch EMF and publishes them
//     to CloudWatch.
//   - The ECS resource detector adds aws.ecs.*, container.* and cloud.*
//     resource attributes automatically when running inside an ECS task; it
//     is a no-op everywhere else (e.g. local development).
//
// Sidecar containers in an ECS task share the same network namespace as the
// application container, so the collector is always reachable on localhost.
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// Config controls how the OpenTelemetry SDK is configured.
type Config struct {
	// Enabled toggles the real SDK on or off. When false, Setup wires up
	// no-op global providers so instrumentation call sites stay cheap and
	// side-effect free (useful for local development without a collector).
	Enabled bool

	ServiceName    string
	ServiceVersion string
	Environment    string

	// OTLPEndpoint is the ADOT collector's OTLP/gRPC endpoint, e.g.
	// "localhost:4317". In an ECS task this is always reachable on
	// localhost because sidecars share the task's network namespace.
	OTLPEndpoint string
	// OTLPInsecure disables TLS for the OTLP/gRPC connection. This is safe
	// (and expected) for a sidecar reachable only via localhost.
	OTLPInsecure bool

	// TracesSamplerRatio is the fraction (0.0-1.0) of root spans sampled.
	// Applied via a parent-based sampler, so any sampling decision made
	// upstream is always respected.
	TracesSamplerRatio float64

	// MetricExportInterval controls how often metrics are pushed to the
	// collector.
	MetricExportInterval time.Duration
}

// Shutdown flushes and stops all telemetry providers. It should be called
// once during application shutdown, with a bounded context so a stuck
// exporter can't hang process termination indefinitely.
type Shutdown func(context.Context) error

func noopShutdown(context.Context) error { return nil }

// Setup wires up the global OpenTelemetry tracer and meter providers and
// registers custapi's custom business metrics. Callers must invoke the
// returned Shutdown exactly once when the application is stopping.
func Setup(ctx context.Context, cfg Config) (Shutdown, error) {
	if !cfg.Enabled {
		log.Println("telemetry: disabled, using no-op tracer/meter providers")
		if err := initMetrics(); err != nil {
			return noopShutdown, fmt.Errorf("telemetry: registering custom metrics: %w", err)
		}
		return noopShutdown, nil
	}

	res, err := newResource(ctx, cfg)
	if err != nil {
		// A schema URL conflict between resource detectors is non-fatal:
		// the resulting resource still carries every attribute, just
		// without a single schema URL annotation. Log and continue.
		log.Printf("telemetry: partial resource detection: %v", err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		xray.Propagator{},
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracerProvider, err := newTracerProvider(ctx, cfg, res)
	if err != nil {
		return noopShutdown, fmt.Errorf("telemetry: setting up tracer provider: %w", err)
	}
	otel.SetTracerProvider(tracerProvider)

	meterProvider, err := newMeterProvider(ctx, cfg, res)
	if err != nil {
		return noopShutdown, fmt.Errorf("telemetry: setting up meter provider: %w", err)
	}
	otel.SetMeterProvider(meterProvider)

	if err := initMetrics(); err != nil {
		return noopShutdown, fmt.Errorf("telemetry: registering custom metrics: %w", err)
	}

	log.Printf(
		"telemetry: enabled (service=%s version=%s env=%s otlp_endpoint=%s)",
		cfg.ServiceName, cfg.ServiceVersion, cfg.Environment, cfg.OTLPEndpoint,
	)

	return func(shutdownCtx context.Context) error {
		return errors.Join(
			tracerProvider.Shutdown(shutdownCtx),
			meterProvider.Shutdown(shutdownCtx),
		)
	}, nil
}

// newResource describes this process: static service identity plus
// host/OS/process/ECS facts gathered from the environment it runs in.
func newResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentNameKey.String(cfg.Environment),
		),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithDetectors(ecs.NewResourceDetector()),
	)
}

func newTracerProvider(ctx context.Context, cfg Config, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint)}
	if cfg.OTLPInsecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
	}

	return sdktrace.NewTracerProvider(
		// X-Ray requires trace IDs that embed the segment's start time; the
		// standard random ID generator produces IDs X-Ray would reject.
		sdktrace.WithIDGenerator(xray.NewIDGenerator()),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.TracesSamplerRatio))),
		sdktrace.WithBatcher(exporter),
	), nil
}

func newMeterProvider(ctx context.Context, cfg Config, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint)}
	if cfg.OTLPInsecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP metric exporter: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(cfg.MetricExportInterval))

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	), nil
}
