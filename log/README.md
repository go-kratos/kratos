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
log.Infof("listening on %s", addr)
log.Info("service started", "service.name", "helloworld", "service.version", "v1.0.0")
log.InfoContext(ctx, "user created", "user_id", userID)
```

### Builder

`log.NewLogger` assembles a fully wired `*slog.Logger` with the kratos
defaults. Pass a handler directly when you already have one; attach fixed
service attrs with `log.WithAttrs`.

```go
logger := log.NewLogger(
	log.WithHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})),
	log.WithAttrs(
		slog.String("service.id", id),
		slog.String("service.name", name),
		slog.String("service.version", version),
	),
	log.WithFilter(log.FilterKey("password")), // redact sensitive keys
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
	"log/slog"

	otel "github.com/go-kratos/kratos/contrib/otel/v2/log"
	"github.com/go-kratos/kratos/v2/log"
)

logger := otel.NewLogger("helloworld")
log.SetDefault(logger)
```

The `github.com/go-kratos/kratos/contrib/otel/v2/log` handler bridges slog records to
OpenTelemetry Logs. Use `otel.WithLogOptions` when you need Kratos logger options:

```go
logger := otel.NewLogger(
	"helloworld",
	otel.WithLogOptions(
		log.WithAttrs(slog.String("service.name", "helloworld")),
		log.WithFilter(log.FilterKey("password")),
	),
)
```

## Third party log library

Adapters that wrap an existing logger accept core builder options directly on
`NewLogger`. Remote-service adapters keep their connection options and expose
`WithLogOptions` for core builder options.

### zap

```shell
go get -u github.com/go-kratos/kratos/contrib/log/zap/v2
```

```go
logger := kratoszap.NewLogger(
	zapLogger,
	log.WithAttrs(slog.String("service.name", "helloworld")),
)
```

### logrus

```shell
go get -u github.com/go-kratos/kratos/contrib/log/logrus/v2
```

```go
logger := kratoslogrus.NewLogger(logrusLogger)
```

### fluent

```shell
go get -u github.com/go-kratos/kratos/contrib/log/fluent/v2
```

```go
logger, err := kratosfluent.NewLogger(
	"tcp://127.0.0.1:24224",
	kratosfluent.WithLogOptions(log.WithAttrs(slog.String("service.name", "helloworld"))),
)
```

### aliyun

```shell
go get -u github.com/go-kratos/kratos/contrib/log/aliyun/v2
```

```go
logger, err := kratosaliyun.NewLogger(
	kratosaliyun.WithProject("project"),
	kratosaliyun.WithLogstore("app"),
	kratosaliyun.WithLogOptions(log.WithAttrs(slog.String("service.name", "helloworld"))),
)
```
