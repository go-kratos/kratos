package jwt

// TokenProvider provider jwt token
type TokenProvider interface {
	GetToken() string
}
