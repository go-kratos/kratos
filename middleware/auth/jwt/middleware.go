package jwt

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
	"time"
)

type jwtKey string

const (

	// JWTContextKey holds the key used to store a JWT in the context.
	JWTContextKey jwtKey = "JWTToken"

	// JWTClaimsContextKey holds the key used to store the JWT Claims in the context.
	JWTClaimsContextKey jwtKey = "JWTClaims"

	//JWTHeaderKey holds the key used to store the JWT Token in the request header.
	JWTHeaderKey string = "Authorization"
)

var (
	ErrMissingJwtToken        = errors.New("JWT is missing")
	ErrMissingAccessSecret    = errors.New("AccessSecret is missing")
	ErrTokenInvalid           = errors.New("Token is invalid")
	ErrUnSupportSigningMethod = errors.New("Wrong signing method")
	ErrNeedTokenProvider      = errors.New("Miss token provider")
	ErrWrongContext           = errors.New("Wrong context for middelware")
)

//Option is jwt option.
type Option func(*options)

type options struct {
	AccessSecret         string
	AccessExpireInSecond time.Duration
	SigningMethod        jwt.SigningMethod
}

//WithAccessExpire with access expire option.
func WithAccessExpire(second time.Duration) Option {
	return func(o *options) {
		o.AccessExpireInSecond = second
	}
}

//WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.SigningMethod = method
	}
}

//Server is an server jwt middleware
func Server(accessSecret string, opts ...Option) middleware.Middleware {
	o := initOptions(accessSecret, opts...)
	parser := func(token string) (*jwt.Token, error) {
		return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if token.Method != o.SigningMethod {
				return nil, ErrUnSupportSigningMethod
			}
			if o.AccessSecret == "" {
				return nil, ErrMissingAccessSecret
			}
			return []byte(o.AccessSecret), nil
		})
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			jwtToken := fromHeader(ctx)
			if jwtToken == "" {
				return nil, ErrMissingJwtToken
			}
			token, err := parser(jwtToken)
			if err != nil {
				return nil, err
			}
			if !token.Valid {
				return nil, ErrTokenInvalid
			}
			ctx = context.WithValue(ctx, JWTClaimsContextKey, token.Claims)
			return handler(ctx, req)
		}
	}
}

//Client is an client jwt middleware
func Client(provider TokenProvider) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if provider == nil {
				return nil, ErrNeedTokenProvider
			}
			err := toHeader(ctx, provider.GetToken())
			if err != nil {
				return nil, err
			}
			return handler(ctx, req)
		}
	}
}

//TokenProvider provider jwt token
type TokenProvider interface {
	GetToken() string
}

//initOptions init jwt option. And set the default option
func initOptions(accessSecret string, opts ...Option) *options {
	o := &options{
		AccessSecret: accessSecret,
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.SigningMethod == nil {
		o.SigningMethod = jwt.SigningMethodHS256
	}
	return o
}

//fromHeader get token from header
func fromHeader(ctx context.Context) string {
	var jwtToken string
	if serverContext, ok := transport.FromServerContext(ctx); ok {
		jwtToken = serverContext.RequestHeader().Get(JWTHeaderKey)
	}
	return jwtToken
}

//toHeader add token to header
func toHeader(ctx context.Context, token string) error {
	if clientContext, ok := transport.FromClientContext(ctx); ok {
		clientContext.RequestHeader().Set(JWTHeaderKey, token)
	} else {
		return ErrWrongContext
	}
	return nil
}
