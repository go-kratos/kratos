package oauth2

import (
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/transport"
)

const (
	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// authorizationKey holds the key used to store the access token in the request header.
	authorizationKey string = "Authorization"

	// wwwAuthenticateKey holds the key of the RFC 6750 challenge in the reply header.
	wwwAuthenticateKey string = "WWW-Authenticate"
)

// Each error has its own reason, so they can be told apart with errors.Is.
var (
	ErrMissingToken       = errors.Unauthorized("OAUTH2_TOKEN_MISSING", "OAuth2 bearer token is missing")
	ErrInvalidToken       = errors.Unauthorized("OAUTH2_TOKEN_INVALID", "OAuth2 token is invalid")
	ErrTokenExpired       = errors.Unauthorized("OAUTH2_TOKEN_EXPIRED", "OAuth2 token has expired")
	ErrMissingValidator   = errors.Unauthorized("OAUTH2_VALIDATOR_MISSING", "Token validator is missing")
	ErrMissingTokenSource = errors.Unauthorized("OAUTH2_TOKEN_SOURCE_MISSING", "Token source is missing")
	ErrTokenSource        = errors.Unauthorized("OAUTH2_TOKEN_SOURCE", "Can not get token from token source")
	ErrWrongContext       = errors.Unauthorized("OAUTH2_WRONG_CONTEXT", "Wrong context for middleware")
	ErrInsufficientScope  = errors.Forbidden("OAUTH2_INSUFFICIENT_SCOPE", "OAuth2 token has insufficient scope")
	ErrIntrospection      = errors.ServiceUnavailable("OAUTH2_INTROSPECTION", "Can not introspect OAuth2 token")
)

// RFC 6750 section 3.1 error codes.
const (
	challengeInvalidToken      = "invalid_token"
	challengeInsufficientScope = "insufficient_scope"
)

// setChallenge writes an RFC 6750 WWW-Authenticate challenge into the reply header.
func setChallenge(tr transport.Transporter, realm, code, scope string) {
	h := tr.ReplyHeader()
	if h == nil {
		return
	}
	params := make([]string, 0, 3)
	if realm != "" {
		params = append(params, fmt.Sprintf("realm=%q", realm))
	}
	if code != "" {
		params = append(params, fmt.Sprintf("error=%q", code))
	}
	if scope != "" {
		params = append(params, fmt.Sprintf("scope=%q", scope))
	}
	value := bearerWord
	if len(params) > 0 {
		value += " " + strings.Join(params, ", ")
	}
	h.Set(wwwAuthenticateKey, value)
}
