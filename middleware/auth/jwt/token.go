package jwt

// TokenManager manager the jwt token.
type TokenManager interface {
	Token() string
}
