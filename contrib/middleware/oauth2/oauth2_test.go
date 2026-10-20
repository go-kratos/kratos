package oauth2

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"

	kerrors "github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/transport"
)

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string        { return http.Header(hc).Get(key) }
func (hc headerCarrier) Set(key string, value string) { http.Header(hc).Set(key, value) }
func (hc headerCarrier) Add(key string, value string) { http.Header(hc).Add(key, value) }
func (hc headerCarrier) Values(key string) []string   { return http.Header(hc).Values(key) }
func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

type mockTransport struct {
	reqHeader   headerCarrier
	replyHeader headerCarrier
}

func newMockTransport() *mockTransport {
	return &mockTransport{reqHeader: headerCarrier{}, replyHeader: headerCarrier{}}
}

func (tr *mockTransport) Kind() transport.Kind            { return transport.KindHTTP }
func (tr *mockTransport) Endpoint() string                { return "" }
func (tr *mockTransport) Operation() string               { return "/test.Service/Method" }
func (tr *mockTransport) RequestHeader() transport.Header { return tr.reqHeader }
func (tr *mockTransport) ReplyHeader() transport.Header   { return tr.replyHeader }

func okHandler(ctx context.Context, _ any) (any, error) {
	info, _ := FromContext(ctx)
	return info, nil
}

type countingSource struct {
	calls int
	token *oauth2.Token
	err   error
}

func (s *countingSource) Token() (*oauth2.Token, error) {
	s.calls++
	return s.token, s.err
}

func TestClient(t *testing.T) {
	src := &countingSource{token: &oauth2.Token{AccessToken: "abc", TokenType: "bearer", Expiry: time.Now().Add(time.Hour)}}
	h := Client(src)(func(context.Context, any) (any, error) { return "ok", nil })

	for range 3 {
		tr := newMockTransport()
		ctx := transport.NewClientContext(context.Background(), tr)
		if _, err := h(ctx, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := tr.reqHeader.Get(authorizationKey); got != "Bearer abc" {
			t.Fatalf("Authorization = %q, want %q", got, "Bearer abc")
		}
	}
	if src.calls != 1 {
		t.Fatalf("token source called %d times, want 1 (token should be reused)", src.calls)
	}
}

func TestClientErrors(t *testing.T) {
	ctx := transport.NewClientContext(context.Background(), newMockTransport())
	tests := []struct {
		name string
		mw   func() (any, error)
		want *kerrors.Error
	}{
		{"nil source", func() (any, error) { return Client(nil)(okHandler)(ctx, nil) }, ErrMissingTokenSource},
		{"wrong context", func() (any, error) {
			return Client(oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "a"}))(okHandler)(context.Background(), nil)
		}, ErrWrongContext},
		{"source error", func() (any, error) {
			return Client(&countingSource{err: errors.New("boom")})(okHandler)(ctx, nil)
		}, ErrTokenSource},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.mw(); !kerrors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestServer(t *testing.T) {
	now := time.Now()
	validator := ValidatorFunc(func(_ context.Context, token string) (*TokenInfo, error) {
		switch token {
		case "good":
			return &TokenInfo{Active: true, Subject: "alice", Scopes: []string{"read", "write"}, Expiry: now.Add(time.Hour)}, nil
		case "inactive":
			return &TokenInfo{Active: false}, nil
		case "expired":
			return &TokenInfo{Active: true, Expiry: now.Add(-time.Minute)}, nil
		case "unavailable":
			return nil, ErrIntrospection
		default:
			return nil, errors.New("unknown token")
		}
	})

	tests := []struct {
		name      string
		auth      string
		want      *kerrors.Error
		challenge string
	}{
		{"valid", "Bearer good", nil, ""},
		{"lowercase scheme", "bearer good", nil, ""},
		{"missing", "", ErrMissingToken, `Bearer realm="api"`},
		{"wrong scheme", "Basic good", ErrMissingToken, `Bearer realm="api"`},
		{"empty token", "Bearer ", ErrMissingToken, `Bearer realm="api"`},
		{"inactive", "Bearer inactive", ErrInvalidToken, `Bearer realm="api", error="invalid_token"`},
		{"unknown", "Bearer nope", ErrInvalidToken, `Bearer realm="api", error="invalid_token"`},
		{"expired", "Bearer expired", ErrTokenExpired, `Bearer realm="api", error="invalid_token"`},
		{"validator unavailable", "Bearer unavailable", ErrIntrospection, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := newMockTransport()
			if tt.auth != "" {
				tr.reqHeader.Set(authorizationKey, tt.auth)
			}
			ctx := transport.NewServerContext(context.Background(), tr)
			reply, err := Server(validator, WithRealm("api"))(okHandler)(ctx, nil)
			if tt.want == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if info := reply.(*TokenInfo); info.Subject != "alice" {
					t.Fatalf("subject = %q, want alice", info.Subject)
				}
				return
			}
			if !kerrors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
			if got := tr.replyHeader.Get(wwwAuthenticateKey); got != tt.challenge {
				t.Fatalf("WWW-Authenticate = %q, want %q", got, tt.challenge)
			}
		})
	}
}

func TestServerErrors(t *testing.T) {
	if _, err := Server(nil)(okHandler)(context.Background(), nil); !kerrors.Is(err, ErrWrongContext) {
		t.Fatalf("err = %v, want %v", err, ErrWrongContext)
	}
	ctx := transport.NewServerContext(context.Background(), newMockTransport())
	if _, err := Server(nil)(okHandler)(ctx, nil); !kerrors.Is(err, ErrMissingValidator) {
		t.Fatalf("err = %v, want %v", err, ErrMissingValidator)
	}
}

func TestRequireScopes(t *testing.T) {
	info := &TokenInfo{Active: true, Scopes: []string{"read", "write"}}
	tests := []struct {
		name      string
		ctx       context.Context
		scopes    []string
		want      *kerrors.Error
		challenge string
	}{
		{"granted", NewContext(context.Background(), info), []string{"read"}, nil, ""},
		{"all granted", NewContext(context.Background(), info), []string{"read", "write"}, nil, ""},
		{"no scopes required", NewContext(context.Background(), info), nil, nil, ""},
		{
			name:      "insufficient",
			ctx:       NewContext(context.Background(), info),
			scopes:    []string{"read", "admin"},
			want:      ErrInsufficientScope,
			challenge: `Bearer error="insufficient_scope", scope="read admin"`,
		},
		{"no token info", context.Background(), []string{"read"}, ErrMissingToken, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := newMockTransport()
			ctx := transport.NewServerContext(tt.ctx, tr)
			_, err := RequireScopes(tt.scopes...)(okHandler)(ctx, nil)
			if tt.want == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !kerrors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
			if got := tr.replyHeader.Get(wwwAuthenticateKey); got != tt.challenge {
				t.Fatalf("WWW-Authenticate = %q, want %q", got, tt.challenge)
			}
		})
	}
}
