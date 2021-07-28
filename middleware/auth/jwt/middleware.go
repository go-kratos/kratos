package jwt

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type jwtKey string

const (

	// JWTContextKey holds the key used to store a JWT in the context.
	JWTContextKey jwtKey = "JWTToken"

	// JWTClaimsContextKey holds the key used to store the JWT Claims in the context.
	JWTClaimsContextKey jwtKey = "JWTClaims"
)

var (
	ErrMissingJwtToken        = errors.New("JWT is missing")
	ErrMissingAccessSecret    = errors.New("AccessSecret is missing")
	ErrMissionSigningMethod   = errors.New("SigningMethod is missing")
	ErrTokenInvalid           = errors.New("Token is invalid")
	ErrUnSupportSigningMethod = errors.New("Wrong signing method")
)

//Option is jwt option.
type Option func(*options)

type options struct {
	AccessSecret         string
	AccessExpireInSecond uint32
	SigningMethod        jwt.SigningMethod
}

//WithAccessSecret with access secret option.
func WithAccessSecret(accessSecret string) Option {
	return func(o *options) {
		o.AccessSecret = accessSecret
	}
}

//WithAccessExpire with access expire option.
func WithAccessExpire(second uint32) Option {
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
func Server(opts ...Option) middleware.Middleware {
	o := initOptions(opts...)
	parser := NewParser(o)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			jwtToken := fromHeader(ctx, req)
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

func initOptions(opts ...Option) options {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	if o.AccessSecret == "" {
		panic(ErrMissingAccessSecret)
	}
	if o.SigningMethod == nil {
		panic(ErrMissionSigningMethod)
	}
	return o
}

//fromHeader get token from header
func fromHeader(ctx context.Context, req interface{}) string {
	var jwtToken string
	if request, ok := req.(http.Request); ok {
		jwtToken = request.Header.Get("Authorization")
	} else if md, ok := metadata.FromIncomingContext(ctx); ok {
		jwtToken = md.Get("Authorization")[0]
	}
	return jwtToken
}
