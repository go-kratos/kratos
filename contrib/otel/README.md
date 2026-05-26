# OpenTelemetry contrib

This module keeps OpenTelemetry integrations out of the core Kratos module.

## Packages

- `github.com/go-kratos/kratos/contrib/otel/v3/log`: slog handler bridge, usually imported as `otel` for `otel.NewHandler`.
- `github.com/go-kratos/kratos/contrib/otel/v3/tracing`: tracing middleware and trace slog attributes.
- `github.com/go-kratos/kratos/contrib/otel/v3/metrics`: metrics middleware and OTel metric helpers.

## Logger

```go
import (
	otel "github.com/go-kratos/kratos/contrib/otel/v3/log"
	"github.com/go-kratos/kratos/v3/log"
)

logger := log.NewLogger(otel.NewHandler("helloworld"))
```

Use the core Kratos log builder when the logger also needs fixed attrs or
filtering:

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

Log, tracing, and metrics stay as shallow subpackages because they expose common
names such as `NewHandler`, `Server`, `Client`, and `Option`.
