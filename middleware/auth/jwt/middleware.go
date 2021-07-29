package jwt

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
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
func Server(accessSecret string, signingMethod jwt.SigningMethod, opts ...Option) middleware.Middleware {
	checkParam(accessSecret, signingMethod)
	o := initOptions(opts...)
	o.AccessSecret = accessSecret
	o.SigningMethod = signingMethod
	parser := NewParser(o)
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

func checkParam(accessSecret string, signingMethod jwt.SigningMethod) {
	if accessSecret == "" {
		panic(ErrMissingAccessSecret)
	}
	if signingMethod == nil {
		panic(ErrMissionSigningMethod)
	}
}

func initOptions(opts ...Option) options {
	o := options{}
	for _, opt := range opts {
		opt(&o)
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
