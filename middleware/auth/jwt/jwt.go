package jwt

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
)

type authkey string

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// HeaderKey holds the key used to store the JWT Token in the request header.
	HeaderKey string = "Authorization"

	// RefreshTokenHeaderKey holds the key used to store the JWT Refresh Token in the request header
	RefreshTokenHeaderKey string = "RefreshToken"

	// InfoKey holds the key used to store the auth info in the context
	InfoKey authkey = "AuthInfo"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized("Missing info", "JWT token is missing")
	ErrMissingAccessSecret    = errors.Unauthorized("Missing info", "AccessSecret is missing")
	ErrTokenInvalid           = errors.Unauthorized("Token invalid", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("Token invalid", "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("Something wrong", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("Something wrong", "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized("Something wrong", "Wrong context for middelware")
	ErrNeedTokenProvider      = errors.Unauthorized("Missing info", "Token provider is missing")
)

// Option is jwt option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	accessSecret          string
	signingMethod         jwt.SigningMethod
	authHeaderKey         string
	refreshTokenHeaderKey string
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithAuthHeaderKey set key that hold auth token in header
func WithAuthHeaderKey(headerKey string) Option {
	return func(options *options) {
		options.authHeaderKey = headerKey
	}
}

// WithRefreshTokenHeaderKey set key that hold refresh token in header
func WithRefreshTokenHeaderKey(headerKey string) Option {
	return func(o *options) {
		o.refreshTokenHeaderKey = headerKey
	}
}

// Server is a server auth middleware
func Server(accessSecret string, opts ...Option) middleware.Middleware {
	o := initOptions(accessSecret, opts...)
	parser := newParser(o)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				tokenInfo, err := parser(header.RequestHeader().Get(o.authHeaderKey))
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, InfoKey, tokenInfo)
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// Client is a client jwt middleware
func Client(provider TokenProvider, opts ...Option) middleware.Middleware {
	o := initOptions("", opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if provider == nil {
				return nil, ErrNeedTokenProvider
			}
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(o.authHeaderKey, provider.GetToken())
				clientContext.RequestHeader().Set(o.refreshTokenHeaderKey, provider.GetRefreshToken())
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// initOptions init the option
func initOptions(accessSecret string, opts ...Option) *options {
	o := &options{
		accessSecret:          accessSecret,
		authHeaderKey:         HeaderKey,
		signingMethod:         jwt.SigningMethodHS256,
		refreshTokenHeaderKey: RefreshTokenHeaderKey,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// newParser create a jwt token parser.
func newParser(o *options) func(jwtToken string) (*jwt.Token, error) {
	return func(jwtToken string) (*jwt.Token, error) {
		/*check the access secret*/
		if o.accessSecret == "" {
			return nil, ErrMissingAccessSecret
		}
		auths := strings.Split(jwtToken, " ")
		if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
			return nil, ErrMissingJwtToken
		}
		jwtToken = auths[1]
		/*parse token*/
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(o.accessSecret), nil
		})
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					/*token format error*/
					return nil, ErrTokenInvalid
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					/*Token is either expired or not active yet*/
					return nil, ErrTokenExpired
				} else {
					return nil, ErrTokenParseFail
				}
			}
		} else if !token.Valid {
			return nil, ErrTokenInvalid
		} else if token.Method != o.signingMethod {
			return nil, ErrUnSupportSigningMethod
		}
		return token, err
	}
}
