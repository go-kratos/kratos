package jwt

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

const (

	//bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	//bearerFormat authorization token format
	bearerFormat string = "Bearer %s"
)

var (
	ErrMissingJwtToken        = errors.Unauthorized("Missing info", "JWT token is missing")
	ErrMissingAccessSecret    = errors.Unauthorized("Missing info", "AccessSecret is missing")
	ErrTokenInvalid           = errors.Unauthorized("Token invalid", "Token is invalid")
	ErrTokenExpired           = errors.Unauthorized("Token invalid", "JWT token has expired")
	ErrTokenFormat            = errors.Unauthorized("Token invalid", "JWT token format error")
	ErrTokenParseFail         = errors.Unauthorized("Something wrong", "Fail to parse JWT token ")
	ErrUnSupportSigningMethod = errors.Unauthorized("Something wrong", "Wrong signing method")
)

//Option is jwt option.
type Option func(*JwtMethod)

type JwtMethod struct {
	AccessSecret         string
	AccessExpireInSecond time.Duration
	SigningMethod        jwt.SigningMethod
}

//WithAccessExpire with access expire option.
func WithAccessExpire(second time.Duration) Option {
	return func(o *JwtMethod) {
		o.AccessExpireInSecond = second
	}
}

//WithSigningMethod with signing method option.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *JwtMethod) {
		o.SigningMethod = method
	}
}

func (j JwtMethod) ParseToken(jwtToken string) (interface{}, error) {
	/*check the access secret*/
	if j.AccessSecret == "" {
		return nil, ErrMissingAccessSecret
	}
	auths := strings.Split(jwtToken, " ")
	if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
		return nil, ErrMissingJwtToken
	}
	jwtToken = auths[1]
	/*parse token*/
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.AccessSecret), nil
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
	} else if token.Method != j.SigningMethod {
		return nil, ErrUnSupportSigningMethod
	}
	return token, err
}

//NewJWTParser create a jwt token parser.
func NewJWTParser(accessSecret string, opts ...Option) *JwtMethod {
	method := &JwtMethod{
		AccessSecret: accessSecret,
	}
	for _, opt := range opts {
		opt(method)
	}
	if method.SigningMethod == nil {
		method.SigningMethod = jwt.SigningMethodHS256
	}
	return method
}
