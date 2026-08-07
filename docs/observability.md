# Observability: OpenTelemetry, ADOT, X-Ray & CloudWatch

custapi is instrumented with [OpenTelemetry](https://opentelemetry.io/) for
both traces and metrics. In production it ships telemetry to an
**AWS Distro for OpenTelemetry (ADOT) Collector** running as a sidecar
container in the ECS task, which forwards:

- **Traces → AWS X-Ray**
- **Metrics → CloudWatch** (via the collector's `awsemf` exporter, which
  writes CloudWatch embedded metric format log events that CloudWatch turns
  into metrics automatically)

## How it's wired up (code)

| Concern | Where |
|---|---|
| SDK setup (providers, resource, propagator, shutdown) | `internal/telemetry/telemetry.go` |
| Custom business metrics + manual tracer | `internal/telemetry/metrics.go` |
| Config / env vars | `internal/config/config.go` (`Config.Telemetry`) |
| HTTP tracing + RED metrics | `github.com/gofiber/contrib/otelfiber/v2` middleware, registered first in `internal/routes/routes.go` |
| DB tracing + metrics | `gorm.io/plugin/opentelemetry/tracing` plugin, registered in `internal/database/db.go` |
| Log ↔ trace correlation | `internal/middleware/middleware.go` (`Logger`) prints `trace_id`/`span_id` from `c.UserContext()` |

Key implementation details worth knowing before touching this code:

- **`c.UserContext()`, not `c.Context()`.** otelfiber stores the
  request's span on the Fiber context via `c.SetUserContext()`. Handlers
  pass `c.UserContext()` into services/repositories so that DB spans (via
  `db.WithContext(ctx)`, already used everywhere in `internal/repositories`)
  and any manual spans are correctly parented to the HTTP request span. If
  you add a new handler, use `c.UserContext()`, not `c.Context()`.
- **X-Ray-compatible trace IDs.** `sdktrace.WithIDGenerator(xray.NewIDGenerator())`
  makes every trace ID embed a timestamp the way X-Ray requires — the
  default random ID generator's IDs would be silently dropped by X-Ray.
- **Propagator** is a composite of the AWS X-Ray header format
  (`X-Amzn-Trace-Id`), W3C `traceparent`, and `baggage`, so custapi both
  reads/writes X-Ray-style headers (e.g. from an ALB) and interops with any
  upstream/downstream service using plain W3C trace context.
- **Resource attributes** come from `service.name`/`version`/
  `deployment.environment.name` plus host/OS/process detectors and the
  [ECS resource detector](https://pkg.go.dev/go.opentelemetry.io/contrib/detectors/aws/ecs)
  (`aws.ecs.*`, `container.*`, `cloud.*`). The ECS detector silently no-ops
  outside of ECS, so this is safe to run locally too.
- **Disabling telemetry** (e.g. running `go run cmd/main.go` locally with no
  collector listening on `4317`) is a single env var: `OTEL_ENABLED=false`.
  All instrumentation call sites (`telemetry.RecordUserCreated`, etc.) stay
  safe no-ops in that mode.

## Configuration

All of it is environment-driven (see `.env.example`):

| Env var | Default | Purpose |
|---|---|---|
| `OTEL_ENABLED` | `true` | Master on/off switch. |
| `OTEL_SERVICE_NAME` | `custapi` | `service.name` resource attribute. |
| `OTEL_SERVICE_VERSION` | `dev` | `service.version`; set to the release tag/SHA in CI. |
| `DEPLOYMENT_ENVIRONMENT` | `development` | `deployment.environment.name` (`production`, `qa`, ...). |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `localhost:4317` | ADOT collector's OTLP/gRPC endpoint. Always `localhost` in ECS — sidecars share the task's network namespace. |
| `OTEL_EXPORTER_OTLP_INSECURE` | `true` | Plaintext gRPC to the sidecar (no TLS needed over loopback). |
| `OTEL_TRACES_SAMPLER_RATIO` | `1.0` | Fraction of root spans sampled (parent-based, so any upstream sampling decision is respected). Lower this if request volume in production gets expensive to trace at 100%. |
| `OTEL_METRIC_EXPORT_INTERVAL_MS` | `60000` | How often metrics are pushed to the collector. |

## Deploying the ADOT sidecar on ECS

Reference files live in `deploy/otel/`:

- [`ecs-task-definition.example.json`](../deploy/otel/ecs-task-definition.example.json) —
  a two-container Fargate task definition: `custapi` + `aws-otel-collector`.
- [`otel-collector-config.yaml`](../deploy/otel/otel-collector-config.yaml) —
  a custom ADOT config (OTLP receiver → X-Ray + CloudWatch EMF exporters).
  ADOT's image already ships a working default at
  `/etc/ecs/ecs-default-config.yaml`; only supply a custom config if you
  need a fixed CloudWatch namespace/log group or extra processors.

### Steps

1. **Get the config to the sidecar.** The example task definition passes it
   via the `AOT_CONFIG_CONTENT` environment variable, sourced from an SSM
   Parameter Store `SecureString`/`String` parameter:

   ```bash
   aws ssm put-parameter \
     --name /custapi/otel-collector-config \
     --type String \
     --value file://deploy/otel/otel-collector-config.yaml
   ```

   Alternatively, drop `command`/`AOT_CONFIG_CONTENT` from the task
   definition entirely to use ADOT's built-in default config.

2. **Grant IAM permissions.** The **task role** (used by the collector to
   call AWS APIs — not the execution role) needs at minimum:
   - `xray:PutTraceSegments`, `xray:PutTelemetryRecords` (X-Ray)
   - `cloudwatch:PutMetricData` — only if using the `awscloudwatch` exporter
     instead of `awsemf`
   - `logs:CreateLogGroup`, `logs:CreateLogStream`, `logs:PutLogEvents` —
     `awsemf` writes metrics as CloudWatch Logs EMF events
   - `ssm:GetParameters` — if loading the collector config from SSM

   AWS publishes a managed policy that covers the X-Ray + CloudWatch pieces:
   `AWSXRayDaemonWriteAccess` plus `CloudWatchAgentServerPolicy` (or a
   narrower custom policy with just the actions above).

3. **Register and deploy** the task definition, wiring in real values for
   `<ACCOUNT_ID>`, `<REGION>`, `<owner>`, and the `DATA_SOURCE_NAME` SSM
   parameter:

   ```bash
   aws ecs register-task-definition --cli-input-json file://deploy/otel/ecs-task-definition.example.json
   ```

4. **Verify.** After a deploy, hit a few endpoints and check:
   - **X-Ray console** → Traces: a `custapi` service map node with spans for
     the HTTP route and any DB queries underneath it.
   - **CloudWatch** → Metrics → the `CustAPI` namespace (or whatever
     `namespace` you set in `otel-collector-config.yaml`):
     `http.server.request.duration`, `db.client.operation.duration`,
     `custapi.users.created`, etc.
   - **CloudWatch Logs** → `/ecs/custapi` log group: both the `app` and
     `adot` log streams should show steady output with no exporter
     connection errors.

## Local development

Telemetry works out of the box against nothing (no collector) as long as
`OTEL_ENABLED=false` is set — instrumentation becomes inert. To see real
traces/metrics locally instead, run the ADOT collector as a container
pointed at your own AWS credentials:

```bash
docker run --rm -p 4317:4317 -p 4318:4318 \
  -e AWS_REGION=<region> \
  -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY -e AWS_SESSION_TOKEN \
  -v "$(pwd)/deploy/otel/otel-collector-config.yaml:/etc/otel-config.yaml" \
  public.ecr.aws/aws-observability/aws-otel-collector:latest \
  --config=/etc/otel-config.yaml
```

Then run custapi with `OTEL_ENABLED=true` (the default) and its default
`OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317` — traces and metrics will land
in your AWS account's X-Ray/CloudWatch exactly as they would from ECS.
