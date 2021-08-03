package jwt

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
	"strings"
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

	//bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	//bearerFormat authorization token format
	bearerFormat string = "Bearer %s"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized("Missing info", "JWT token is missing")
	ErrMissingAccessSecret    = errors.Unauthorized("Missing info", "AccessSecret is missing")
	ErrNeedTokenProvider      = errors.Unauthorized("Missing info", "Token provider is missing")
	ErrTokenInvalid           = errors.Unauthorized("Token invalid", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("Token invalid", "JWT token has expired")
	ErrTokenFormat            = errors.Unauthorized("Token invalid", "JWT token format error")
	ErrTokenParseFail         = errors.Unauthorized("Something wrong", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("Something wrong", "Wrong signing method")
	ErrWrongContext           = errors.Unauthorized("Something wrong", "Wrong context for middelware")
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
	parser := newParser(o)
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
		auths := strings.Split(jwtToken, " ")
		if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
			return ""
		}
		jwtToken = auths[1]
	}
	return jwtToken
}

//toHeader add token to header
func toHeader(ctx context.Context, token string) error {
	if clientContext, ok := transport.FromClientContext(ctx); ok {
		clientContext.RequestHeader().Set(JWTHeaderKey, fmt.Sprintf(bearerFormat, token))
	} else {
		return ErrWrongContext
	}
	return nil
}

//newParser create a jwt token parser.
func newParser(o *options) func(tokenStr string) (*jwt.Token, error) {
	return func(tokenStr string) (*jwt.Token, error) {
		/*check the access secret*/
		if o.AccessSecret == "" {
			return nil, ErrMissingAccessSecret
		}
		/*parse token*/
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(o.AccessSecret), nil
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
		} else if token.Method != o.SigningMethod {
			return nil, ErrUnSupportSigningMethod
		}
		return token, err
	}
}
