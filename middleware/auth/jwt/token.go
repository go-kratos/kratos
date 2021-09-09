package jwt

// TokenProvider provide all info that would be used to sign jwt token
type TokenProvider interface {
	AccessSecretKey() []byte
}
