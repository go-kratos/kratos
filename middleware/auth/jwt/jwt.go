package jwt

import (
	"github.com/golang-jwt/jwt"
)

//NewJwtSigner create jwt token
func NewJwtSigner(config options) func() {

	return nil
}

//NewParser parse jwt token
func NewParser(config options) func(token string) (*jwt.Token, error) {
	return func(token string) (*jwt.Token, error) {
		return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if token.Method != config.SigningMethod {
				return nil, ErrUnSupportSigningMethod
			}
			return config.AccessSecret, nil
		})
	}
}
