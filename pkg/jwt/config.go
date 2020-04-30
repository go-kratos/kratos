package jwt

import "github.com/dgrijalva/jwt-go"

const (
	//DefaultContextKey jwt
	DefaultContextKey = "jwt"
)

// Config is a struct for specifying configuration options for the jwt middleware.
type Config struct {
	// The function that will return the Key to validate the JWT.
	// It can be either a shared secret or a public key.
	// Default value: nil
	ValidationKeyGetter jwt.Keyfunc
	// The name of the property in the request where the user (&token) information
	// from the JWT will be stored.
	// Default value: "jwt"
	ContextKey string
	// The function that will be called when there's an error validating the token
	// Default value:
	ErrorHandler errorHandler
	// A boolean indicating if the credentials are required or not
	// Default value: false
	CredentialsOptional bool
	// A function that extracts the token from the request
	// Default: FromAuthHeader (i.e., from Authorization header as bearer token)
	Extractor TokenExtractor
	// Debug flag turns on debugging output
	// Default: false
	Debug bool
	// When set, all requests with the OPTIONS method will use authentication
	// if you enable this option you should register your route with iris.Options(...) also
	// Default: false
	EnableAuthOnOptions bool
	// When set, the middelware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	// Default: nil
	SigningMethod jwt.SigningMethod
	// When set, the expiration time of token will be check every time
	// if the token was expired, expiration error will be returned
	// Default: false
	Expiration bool
}
