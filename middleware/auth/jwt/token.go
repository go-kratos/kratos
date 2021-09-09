package jwt

// KeyProvider provider the key that sign jwt token
type KeyProvider interface {
	Key() []byte
}
