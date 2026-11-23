# OAuth2 Middleware

OAuth2 bearer token middleware for kratos, built on [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2). Works with both HTTP and gRPC transports.

```shell
go get github.com/go-kratos/kratos/contrib/middleware/oauth2/v3
```

## Client

`Client` gets an access token from an `oauth2.TokenSource` and sets the `Authorization` header. Tokens are cached and refreshed on expiry.

```go
import (
	"golang.org/x/oauth2/clientcredentials"

	oauth2mw "github.com/go-kratos/kratos/contrib/middleware/oauth2/v3"
)

cfg := &clientcredentials.Config{
	ClientID:     "app",
	ClientSecret: "secret",
	TokenURL:     "https://auth.example.com/oauth2/token",
	Scopes:       []string{"order.read"},
}
src := cfg.TokenSource(ctx)

httpClient, _ := http.NewClient(ctx, http.WithEndpoint("127.0.0.1:8000"), http.WithMiddleware(oauth2mw.Client(src)))
conn, _ := grpc.NewClient(ctx, grpc.WithEndpoint("127.0.0.1:9000"), grpc.WithMiddleware(oauth2mw.Client(src)))
```

Middleware does not run for `http.Client.Do`. To cover it, use the transport layer instead:

```go
// HTTP: configure TLS on Base, WithTLSConfig only applies to *http.Transport.
http.WithTransport(&oauth2.Transport{Source: src, Base: tlsTransport})

// gRPC: requires TLS. WithOptions replaces previously set dial options.
grpc.WithOptions(grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: src}))
```

## Server

`Server` validates the bearer token with a `Validator` and stores the `*TokenInfo` in the context. Failures return 401 with an RFC 6750 `WWW-Authenticate` challenge.

`NewIntrospectionValidator` validates tokens against an RFC 7662 introspection endpoint:

```go
validator := oauth2mw.NewIntrospectionValidator(
	"https://auth.example.com/oauth2/introspect", nil,
	oauth2mw.WithBasicAuth("resource-server", "secret"),
	oauth2mw.WithCacheTTL(time.Minute),
)

srv := http.NewServer(
	http.Middleware(
		oauth2mw.Server(validator, oauth2mw.WithRealm("api")),
		// Per-operation scopes, using the generated Operation constants.
		selector.Server(oauth2mw.RequireScopes("order.write")).
			Path(v1.OperationOrderCreateOrder).Build(),
	),
)
```

Read the token info in a handler:

```go
info, ok := oauth2mw.FromContext(ctx)
```

To authenticate the introspection call with an access token instead of basic auth, pass `clientcredentials.Config.Client(ctx)` as the client. For JWT access tokens validated with JWKS, use [contrib/middleware/jwt](../jwt) instead.
