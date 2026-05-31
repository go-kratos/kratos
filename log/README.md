# Logger

## Usage

### slog

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

log.SetDefault(logger)

logger.InfoContext(ctx, "user created",
	"user_id", userID,
	"service.name", "helloworld",
)
```

### Global logger

Common global helpers are still available for gradual migration. The signatures
now mirror slog: the first argument is the message, followed by key/value pairs
or `slog.Attr` values.

```go
log.Info("started")
log.Info("listening", "addr", addr)
log.Info("service started", "service.name", "helloworld", "service.version", "v1.0.0")
log.InfoContext(ctx, "user created", "user_id", userID)
```

### Builder

`log.NewHandler` builds a default handler. `log.NewLogger` wraps an existing
handler with Kratos decorators. Attach fixed service attrs with `logger.With`.

```go
logger := log.NewLogger(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}),
	log.WithFilter(log.FilterKey("password")), // redact sensitive keys
).With(
	slog.String("service.id", id),
	slog.String("service.name", name),
	slog.String("service.version", version),
)
log.SetDefault(logger)
```

### Context attrs

Attach attributes to a `context.Context` and they will flow through any
ctx-aware log call automatically.

```go
ctx = log.ContextWithAttrs(ctx, slog.String("request_id", id))
log.InfoContext(ctx, "handling request")
```

### OpenTelemetry

```go
import (
	otel "github.com/go-kratos/kratos/contrib/otel/v3/log"
	"github.com/go-kratos/kratos/v3/log"
)

logger := log.NewLogger(otel.NewHandler("helloworld"))
log.SetDefault(logger)
```

The `github.com/go-kratos/kratos/contrib/otel/v3/log` handler bridges slog records to
OpenTelemetry Logs. Use the core log builder when you need Kratos logger options:

```go
import (
	"log/slog"

	otel "github.com/go-kratos/kratos/contrib/otel/v3/log"
	"github.com/go-kratos/kratos/v3/log"
)

logger := log.NewLogger(
	otel.NewHandler("helloworld"),
	log.WithFilter(log.FilterKey("password")),
).With(slog.String("service.name", "helloworld"))
```
