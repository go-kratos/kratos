# OpenTelemetry contrib

This module keeps OpenTelemetry integrations out of the core Kratos module.

## Packages

- `github.com/go-kratos/kratos/contrib/otel/v3/log`: slog bridge, usually imported as `otel` for `otel.NewLogger`.
- `github.com/go-kratos/kratos/contrib/otel/v3/tracing`: tracing middleware and trace slog attributes.
- `github.com/go-kratos/kratos/contrib/otel/v3/metrics`: metrics middleware and OTel metric helpers.

## Logger

```go
import otel "github.com/go-kratos/kratos/contrib/otel/v3/log"

logger := otel.NewLogger("helloworld")
```

Use `WithLogOptions` when the logger also needs core Kratos log builder
options such as fixed attrs or filtering:

```go
import (
	"log/slog"

	otel "github.com/go-kratos/kratos/contrib/otel/v3/log"
	"github.com/go-kratos/kratos/v3/log"
)

logger := otel.NewLogger(
	"helloworld",
	otel.WithLogOptions(
		log.WithAttrs(slog.String("service.name", "helloworld")),
		log.WithFilter(log.FilterKey("password")),
	),
)
```

Log, tracing, and metrics stay as shallow subpackages because they expose common
names such as `NewLogger`, `Server`, `Client`, and `Option`.
