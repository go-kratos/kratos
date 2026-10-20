package oauth2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/health/grpc_health_v1"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/selector"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	khttp "github.com/go-kratos/kratos/v3/transport/http"
)

const testAccessToken = "e2e-token"

// newTokenServer mocks the token endpoint of an authorization server.
func newTokenServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id, secret, ok := r.BasicAuth(); !ok || id != "app" || secret != "secret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": testAccessToken,
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
}

func testValidator() Validator {
	return ValidatorFunc(func(_ context.Context, token string) (*TokenInfo, error) {
		if token != testAccessToken {
			return &TokenInfo{Active: false}, nil
		}
		return &TokenInfo{Active: true, Subject: "app", Scopes: []string{"ping.read"}}, nil
	})
}

func testTokenSource(tokenURL string) *clientcredentials.Config {
	return &clientcredentials.Config{ClientID: "app", ClientSecret: "secret", TokenURL: tokenURL}
}

func TestHTTPEndToEnd(t *testing.T) {
	tokenSrv := newTokenServer(t)
	defer tokenSrv.Close()

	srv := khttp.NewServer(
		khttp.Address("127.0.0.1:0"),
		khttp.Middleware(
			Server(testValidator()),
			selector.Server(RequireScopes("ping.admin")).Path("/admin").Build(),
		),
	)
	handle := func(op string) khttp.HandlerFunc {
		return func(ctx khttp.Context) error {
			khttp.SetOperation(ctx, op)
			h := ctx.Middleware(func(ctx context.Context, _ any) (any, error) {
				info, _ := FromContext(ctx)
				return map[string]string{"subject": info.Subject}, nil
			})
			out, err := h(ctx, nil)
			if err != nil {
				return err
			}
			return ctx.Result(http.StatusOK, out)
		}
	}
	srv.Route("/").GET("/ping", handle("/ping"))
	srv.Route("/").GET("/admin", handle("/admin"))
	endpoint, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = srv.Start(context.Background()) }()
	defer func() { _ = srv.Stop(context.Background()) }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newClient := func(ms ...middleware.Middleware) *khttp.Client {
		t.Helper()
		client, err := khttp.NewClient(ctx, khttp.WithEndpoint(endpoint.Host), khttp.WithMiddleware(ms...))
		if err != nil {
			t.Fatal(err)
		}
		return client
	}

	client := newClient(Client(testTokenSource(tokenSrv.URL).TokenSource(ctx)))
	defer client.Close()
	var out map[string]string
	if err := client.Invoke(ctx, http.MethodGet, "/ping", nil, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["subject"] != "app" {
		t.Fatalf("subject = %q, want app", out["subject"])
	}
	if err := client.Invoke(ctx, http.MethodGet, "/admin", nil, &out); kerrors.Code(err) != http.StatusForbidden {
		t.Fatalf("admin: err = %v, want 403", err)
	}

	anonymous := newClient()
	defer anonymous.Close()
	if err := anonymous.Invoke(ctx, http.MethodGet, "/ping", nil, &out); kerrors.Code(err) != http.StatusUnauthorized {
		t.Fatalf("anonymous: err = %v, want 401", err)
	}
}

func TestGRPCEndToEnd(t *testing.T) {
	tokenSrv := newTokenServer(t)
	defer tokenSrv.Close()

	// The built-in health service is enough to drive the unary interceptor.
	srv := grpc.NewServer(grpc.Address("127.0.0.1:0"), grpc.Middleware(Server(testValidator())))
	endpoint, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = srv.Start(context.Background()) }()
	defer func() { _ = srv.Stop(context.Background()) }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	check := func(ms ...middleware.Middleware) error {
		t.Helper()
		conn, err := grpc.NewClient(ctx,
			grpc.WithEndpoint(endpoint.Host),
			grpc.WithMiddleware(ms...),
			grpc.WithHealthCheck(false),
		)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		_, err = grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		return err
	}

	if err := check(Client(testTokenSource(tokenSrv.URL).TokenSource(ctx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The client receives a gRPC status, convert it back to a kratos error.
	if se := kerrors.FromError(check()); se == nil || se.Reason != ErrMissingToken.Reason {
		t.Fatalf("anonymous: err = %v, want %v", se, ErrMissingToken)
	}
}
