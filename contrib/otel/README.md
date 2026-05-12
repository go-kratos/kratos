# OpenTelemetry contrib

This module keeps OpenTelemetry integrations out of the core Kratos module.

## Packages

- `github.com/go-kratos/kratos/contrib/otel/v2/log`: slog bridge, usually imported as `otel` for `otel.NewLogger`.
- `github.com/go-kratos/kratos/contrib/otel/v2/tracing`: tracing middleware and trace slog attributes.
- `github.com/go-kratos/kratos/contrib/otel/v2/metrics`: metrics middleware and OTel metric helpers.

## Logger

```go
import otel "github.com/go-kratos/kratos/contrib/otel/v2/log"

logger := otel.NewLogger("helloworld")
```

Use the core log builder when the logger also needs fixed attrs or filtering:

```go
import (
	"log/slog"

	otel "github.com/go-kratos/kratos/contrib/otel/v2/log"
	"github.com/go-kratos/kratos/v2/log"
)

logger := log.NewLogger(
	log.WithHandler(otel.NewHandler("helloworld")),
	log.WithExtractor(otel.TraceAttrs),
	log.WithAttrs(slog.String("service.name", "helloworld")),
)
```

Log, tracing, and metrics stay as shallow subpackages because they expose common
names such as `NewLogger`, `Server`, `Client`, and `Option`.
