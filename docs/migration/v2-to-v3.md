# Migrating from Kratos v2 to v3

Kratos v3 cleans up historical coupling in the core framework and makes several previously implicit dependencies explicit. This guide summarizes the main upgrade work tracked in [go-kratos/kratos#3820](https://github.com/go-kratos/kratos/issues/3820).

## 1. Update Module Paths

Update Kratos imports from `github.com/go-kratos/kratos/v2` to `github.com/go-kratos/kratos/v3`.

```go
// v2
import "github.com/go-kratos/kratos/v2"

// v3
import "github.com/go-kratos/kratos/v3"
```

Then refresh dependencies:

```shell
go get github.com/go-kratos/kratos/v3@latest
go mod tidy
```

Contrib modules also use `/v3` import paths, for example `github.com/go-kratos/kratos/contrib/middleware/jwt/v3`.

## 2. Choose the JSON Codec Explicitly

In v2, `encoding/json` handled both ordinary Go JSON values and `proto.Message` values with protobuf JSON semantics. In v3, the core JSON codecs are split:

- `github.com/go-kratos/kratos/v3/encoding/json` registers the standard-library JSON codec as `json`.
- `github.com/go-kratos/kratos/v3/encoding/protojson` registers the protobuf JSON codec as `protojson`.
- `github.com/go-kratos/kratos/contrib/encoding/json/v3` keeps the v2-compatible `json` codec behavior for migration.

New v3 code should prefer explicit imports:

```go
import (
	_ "github.com/go-kratos/kratos/v3/encoding/json"
	_ "github.com/go-kratos/kratos/v3/encoding/protojson"
)
```

If a service depends on v2 behavior where the `json` codec also handles protobuf messages, use the contrib compatibility codec while migrating:

```go
import _ "github.com/go-kratos/kratos/contrib/encoding/json/v3"
```

Do not register both JSON codecs for the same process unless you intentionally want the later import initialization to replace the earlier `json` codec registration.

## 3. Migrate Logging to slog

Kratos v3 uses the standard-library `log/slog` APIs. Code that depends on `log.Logger`, `log.Helper`, `log.Valuer`, `log.NewStdLogger`, or trace/service helper fields should migrate to `*slog.Logger`, `log.NewHandler`, and `log.NewLogger`.

```go
import (
	"log/slog"
	"os"

	"github.com/go-kratos/kratos/v3/log"
)

logger := log.NewLogger(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}),
	log.WithFilter(log.FilterKey("password")),
).With(
	slog.String("service.name", "helloworld"),
	slog.String("service.version", "v1.0.0"),
)

log.SetDefault(logger)
```

Kratos application and middleware options now accept `*slog.Logger`.

```go
app := kratos.New(
	kratos.Name("helloworld"),
	kratos.Logger(logger),
)
```

OpenTelemetry logging is provided by the contrib handler:

```go
import (
	"github.com/go-kratos/kratos/v3/log"
	otel "github.com/go-kratos/kratos/contrib/otel/v3/log"
)

logger := log.NewLogger(otel.NewHandler("helloworld"))
```

## 4. Update JWT Middleware Imports

JWT middleware moved out of core so projects only pull `github.com/golang-jwt/jwt/v5` when they use JWT.

```go
// v2
import "github.com/go-kratos/kratos/v2/middleware/auth/jwt"

// v3
import "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
```

After changing imports, run:

```shell
go get github.com/go-kratos/kratos/contrib/middleware/jwt/v3@latest
go mod tidy
```

## 5. Review Circuit Breaker Customization

The default circuit breaker no longer depends on `github.com/go-kratos/aegis`. Default usage does not require code changes:

```go
handler := circuitbreaker.Client()(next)
```

If your service injects a custom Aegis breaker, add Aegis as an explicit dependency and adapt the injection point to `WithBreakerFactory`:

```go
handler := circuitbreaker.Client(
	circuitbreaker.WithBreakerFactory(func() circuitbreaker.CircuitBreaker {
		return newBreaker()
	}),
)(next)
```

## 6. Replace Direct HTTP binding Imports

The exported `transport/http/binding` package was removed. Generated `_http.pb.go` files do not need manual changes after regeneration.

For hand-written code:

- Replace `binding.EncodeURL` with `http.BuildPath`.
- Use `transport/http.Context` methods such as `Bind`, `BindVars`, `BindQuery`, and `BindForm`.
- Use `encoding/form` directly only for low-level query/form encoding needs.

```go
path := http.BuildPath("/v1/users/{id}", req)
```

## 7. Regenerate Generated Code

After dependency and import updates, regenerate Kratos generated files:

```shell
go generate ./...
go mod tidy
```

Regeneration is especially important for services using generated HTTP clients or servers because v3 generated code no longer imports the removed HTTP binding package.

## 8. Validate the Upgrade

Run project tests and linters before shipping the migration:

```shell
go test ./...
go vet ./...
```

For this repository, use:

```shell
make test
make lint
```

## Migration Checklist

- [ ] Update imports from `/v2` to `/v3`.
- [ ] Choose `encoding/json`, `encoding/protojson`, or the contrib compatibility JSON codec.
- [ ] Replace old Kratos logging helpers with `log/slog`-based APIs.
- [ ] Move JWT imports to `contrib/middleware/jwt/v3`.
- [ ] Review custom circuit breaker dependencies.
- [ ] Replace direct `transport/http/binding` usage.
- [ ] Regenerate generated code.
- [ ] Run tests, lint, and `go mod tidy`.
