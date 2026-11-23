package oauth2

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
)

type authKey struct{}

// TokenInfo is the information of a validated access token.
// It is shared between requests when cached, treat it as read-only.
type TokenInfo struct {
	Active   bool
	Subject  string
	ClientID string
	Username string
	Issuer   string
	Audience []string
	Scopes   []string
	Expiry   time.Time
	// Extra holds the raw claims returned by the authorization server.
	Extra map[string]any
}

// HasScopes reports whether the token is granted all the given scopes.
func (t *TokenInfo) HasScopes(scopes ...string) bool {
	for _, s := range scopes {
		if !slices.Contains(t.Scopes, s) {
			return false
		}
	}
	return true
}

// Validator validates an access token.
//
// Returning a *errors.Error passes it to the caller as is, any other error
// is treated as an invalid token.
type Validator interface {
	Validate(ctx context.Context, token string) (*TokenInfo, error)
}

// ValidatorFunc is an adapter to use a function as a Validator.
type ValidatorFunc func(ctx context.Context, token string) (*TokenInfo, error)

// Validate calls f(ctx, token).
func (f ValidatorFunc) Validate(ctx context.Context, token string) (*TokenInfo, error) {
	return f(ctx, token)
}

// ServerOption is oauth2 server option.
type ServerOption func(*serverOptions)

type serverOptions struct {
	realm string
	now   func() time.Time
}

// WithRealm with the realm of the WWW-Authenticate challenge.
func WithRealm(realm string) ServerOption {
	return func(o *serverOptions) {
		o.realm = realm
	}
}

// Server is a server oauth2 middleware. It validates the bearer token and
// puts the token info into the context, which can be read by FromContext.
func Server(v Validator, opts ...ServerOption) middleware.Middleware {
	o := &serverOptions{now: time.Now}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, ErrWrongContext
			}
			if v == nil {
				return nil, ErrMissingValidator
			}
			token, ok := bearerToken(tr.RequestHeader().Get(authorizationKey))
			if !ok {
				setChallenge(tr, o.realm, "", "")
				return nil, ErrMissingToken
			}
			info, err := v.Validate(ctx, token)
			if err == nil {
				err = checkToken(info, o.now())
			}
			if err != nil {
				se := toError(err)
				// RFC 6750: invalid or expired tokens get an invalid_token challenge.
				if se.Code == ErrInvalidToken.Code {
					setChallenge(tr, o.realm, challengeInvalidToken, "")
				}
				return nil, se
			}
			return handler(NewContext(ctx, info), req)
		}
	}
}

func bearerToken(auth string) (string, bool) {
	scheme, token, ok := strings.Cut(auth, " ")
	if !ok || !strings.EqualFold(scheme, bearerWord) {
		return "", false
	}
	token = strings.TrimSpace(token)
	return token, token != ""
}

func checkToken(info *TokenInfo, now time.Time) error {
	if info == nil || !info.Active {
		return ErrInvalidToken
	}
	if !info.Expiry.IsZero() && now.After(info.Expiry) {
		return ErrTokenExpired
	}
	return nil
}

func toError(err error) *errors.Error {
	se := new(errors.Error)
	if errors.As(err, &se) {
		return se
	}
	return ErrInvalidToken.WithCause(err)
}

// NewContext put token info into context
func NewContext(ctx context.Context, info *TokenInfo) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext extract token info from context
func FromContext(ctx context.Context) (info *TokenInfo, ok bool) {
	info, ok = ctx.Value(authKey{}).(*TokenInfo)
	return
}
