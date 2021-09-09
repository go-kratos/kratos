package jwt

import (
	"context"
	"fmt"
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

	// InfoKey holds the key used to store the auth info in the context
	InfoKey authkey = "AuthInfo"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized("UNAUTHORIZED", "JWT token is missing")
	ErrMissingKeyFunc         = errors.Unauthorized("UNAUTHORIZED", "keyFunc is missing")
	ErrTokenInvalid           = errors.Unauthorized("UNAUTHORIZED", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("UNAUTHORIZED", "JWT token has expired")
	ErrTokenParseFail         = errors.Unauthorized("UNAUTHORIZED", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("UNAUTHORIZED", "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized("UNAUTHORIZED", "Wrong context for middelware")
	ErrNeedTokenManager       = errors.Unauthorized("UNAUTHORIZED", "Token manager is missing")
)

// Option is jwt option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	signingMethod jwt.SigningMethod
	claims        jwt.Claims
}

// WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithClaims with customer claim
func WithClaims(claims jwt.Claims) Option {
	return func(o *options) {
		o.claims = claims
	}
}

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(keyFunc jwt.Keyfunc, opts ...Option) middleware.Middleware {
	o := &options{
		signingMethod: jwt.SigningMethodHS256,
		claims:        jwt.StandardClaims{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if header, ok := transport.FromServerContext(ctx); ok {
				/*check the access secret*/
				if keyFunc == nil {
					return nil, ErrMissingKeyFunc
				}
				auths := strings.Split(header.RequestHeader().Get(HeaderKey), " ")
				if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
					return nil, ErrMissingJwtToken
				}
				jwtToken := auths[1]
				/*parse token*/
				tokenInfo, err := jwt.Parse(jwtToken, keyFunc)
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
				} else if !tokenInfo.Valid {
					return nil, ErrTokenInvalid
				} else if tokenInfo.Method != o.signingMethod {
					return nil, ErrUnSupportSigningMethod
				}
				ctx = context.WithValue(ctx, InfoKey, tokenInfo)
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}

// Client is a client jwt middleware.
func Client(tokenManager TokenManager) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tokenManager == nil {
				return nil, ErrNeedTokenManager
			}
			if clientContext, ok := transport.FromClientContext(ctx); ok {
				clientContext.RequestHeader().Set(HeaderKey, fmt.Sprintf(bearerFormat, tokenManager.Token()))
				return handler(ctx, req)
			}
			return nil, ErrWrongContext
		}
	}
}
