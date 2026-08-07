package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/hoshina-dev/custapi"

// Tracer returns custapi's tracer for creating manual spans around
// operations that automatic instrumentation (otelfiber for HTTP, the GORM
// plugin for SQL) doesn't cover, such as CPU-bound business logic.
func Tracer() trace.Tracer {
	return otel.Tracer(instrumentationName)
}

// Business metrics layered on top of the automatic RED metrics
// (http.server.* from otelfiber, db.client.* from the GORM plugin): these
// track domain events regardless of which HTTP route or query produced
// them, so dashboards/alarms don't need to reconstruct "was a user created"
// from route + status code combinations.
var (
	usersCreated metric.Int64Counter
	usersDeleted metric.Int64Counter
	orgsCreated  metric.Int64Counter
	orgsDeleted  metric.Int64Counter
)

// initMetrics resolves custapi's meter from whatever MeterProvider is
// currently registered globally (the real OTLP-backed one, or the SDK's
// no-op default when telemetry is disabled) and registers the counters
// above against it. It must run after otel.SetMeterProvider so the
// instruments are bound to the intended provider.
func initMetrics() error {
	meter := otel.Meter(instrumentationName)

	var err error
	if usersCreated, err = meter.Int64Counter(
		"custapi.users.created",
		metric.WithDescription("Number of users successfully created"),
		metric.WithUnit("{user}"),
	); err != nil {
		return fmt.Errorf("custapi.users.created counter: %w", err)
	}

	if usersDeleted, err = meter.Int64Counter(
		"custapi.users.deleted",
		metric.WithDescription("Number of users successfully deleted"),
		metric.WithUnit("{user}"),
	); err != nil {
		return fmt.Errorf("custapi.users.deleted counter: %w", err)
	}

	if orgsCreated, err = meter.Int64Counter(
		"custapi.organizations.created",
		metric.WithDescription("Number of organizations successfully created"),
		metric.WithUnit("{organization}"),
	); err != nil {
		return fmt.Errorf("custapi.organizations.created counter: %w", err)
	}

	if orgsDeleted, err = meter.Int64Counter(
		"custapi.organizations.deleted",
		metric.WithDescription("Number of organizations successfully deleted"),
		metric.WithUnit("{organization}"),
	); err != nil {
		return fmt.Errorf("custapi.organizations.deleted counter: %w", err)
	}

	return nil
}

// RecordUserCreated increments the users-created counter. Safe to call
// even if telemetry is disabled (it becomes a no-op).
func RecordUserCreated(ctx context.Context) {
	usersCreated.Add(ctx, 1)
}

// RecordUserDeleted increments the users-deleted counter.
func RecordUserDeleted(ctx context.Context) {
	usersDeleted.Add(ctx, 1)
}

// RecordOrganizationCreated increments the organizations-created counter.
func RecordOrganizationCreated(ctx context.Context) {
	orgsCreated.Add(ctx, 1)
}

// RecordOrganizationDeleted increments the organizations-deleted counter.
func RecordOrganizationDeleted(ctx context.Context) {
	orgsDeleted.Add(ctx, 1)
}
