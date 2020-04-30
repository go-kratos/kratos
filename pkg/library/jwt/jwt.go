package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kratos/kratos/pkg/ecode"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"strings"
	"time"
)

// iris provides some basic middleware, most for your learning courve.
// You can use any net/http compatible middleware with iris.FromStd wrapper.
//
// JWT net/http video tutorial for golang newcomers: https://www.youtube.com/watch?v=dgJFeqeXVKw
//
// Unlike the other middleware, this middleware was cloned from external source: https://github.com/auth0/go-jwt-middleware
// (because it used "context" to define the user but we don't need that so a simple iris.FromStd wouldn't work as expected.)
// jwt_test.go also didn't created by me:
// 28 Jul 2016
// @heralight heralight add jwt unit test.
//
// So if this doesn't works for you just try other net/http compatible middleware and bind it via `iris.FromStd(myHandlerWithNext)`,
// It's here for your learning curve.

// A function called whenever an error is encountered
type errorHandler func(*bm.Context, string)

// TokenExtractor is a function that takes a context as input and returns
// either a token or an error.  An error should only be returned if an attempt
// to specify a token was found, but the information was somehow incorrectly
// formed.  In the case where a token is simply not present, this should not
// be treated as an error.  An empty string should be returned in that case.
type TokenExtractor func(*bm.Context) (string, error)

// Middleware the middleware for JSON Web tokens authentication method
type Middleware struct {
	Config Config
}

// OnError default error handler
func OnError(ctx *bm.Context, err string) {
	_ = ecode.Errorf(501, err)
}

// New constructs a new Secure instance with supplied options.
func New(cfg ...Config) *Middleware {
	var c Config
	if len(cfg) == 0 {
		c = Config{}
	} else {
		c = cfg[0]
	}
	if c.ContextKey == "" {
		c.ContextKey = DefaultContextKey
	}
	if c.ErrorHandler == nil {
		c.ErrorHandler = OnError
	}
	if c.Extractor == nil {
		c.Extractor = FromAuthHeader
	}
	return &Middleware{Config: c}
}

func (m *Middleware) logf(format string, args ...interface{}) {
	if m.Config.Debug {
		log.Info(format, args...)
	}
}

// Get returns the user (&token) information for this client/request
func (m *Middleware) Get(ctx *bm.Context) *jwt.Token {
	return ctx.Value(m.Config.ContextKey).(*jwt.Token)
}

// Serve the middleware's action
func (m *Middleware) Serve(ctx *bm.Context) {
	if err := m.CheckJWT(ctx); err != nil {
		ctx.AbortWithStatus(401)
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

// FromAuthHeader is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(ctx *bm.Context) (string, error) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// FromParameter returns a function that extracts the token from the specified
// query string parameter
func FromParameter(param string) TokenExtractor {
	return func(ctx *bm.Context) (string, error) {
		return ctx.Params.ByName(param), nil
	}
}

// FromFirst returns a function that runs multiple token extractors and takes the
// first token it finds
func FromFirst(extractors ...TokenExtractor) TokenExtractor {
	return func(ctx *bm.Context) (string, error) {
		for _, ex := range extractors {
			token, err := ex(ctx)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", nil
	}
}

// CheckJWT the main functionality, checks for token
func (m *Middleware) CheckJWT(ctx *bm.Context) error {
	if !m.Config.EnableAuthOnOptions {
		if ctx.Request.Method == "OPTIONS" {
			return nil
		}
	}

	// Use the specified token extractor to extract a token from the request
	token, err := m.Config.Extractor(ctx)

	// If debugging is turned on, log the outcome
	if err != nil {
		m.logf("Error extracting JWT: %v", err)
	} else {
		m.logf("Token extracted: %s", token)
	}

	// If an error occurs, call the error handler and return an error
	if err != nil {
		m.Config.ErrorHandler(ctx, err.Error())
		return fmt.Errorf("error extracting token: %v", err)
	}

	// If the token is empty...
	if token == "" {
		// Check if it was required
		if m.Config.CredentialsOptional {
			m.logf("  no credentials found (CredentialsOptional=true)")
			// No error, just no token (and that is ok given that CredentialsOptional is true)
			return nil
		}

		// If we get here, the required token is missing
		errorMsg := "required authorization token not found"
		m.Config.ErrorHandler(ctx, errorMsg)
		m.logf("  error: No credentials found (CredentialsOptional=false)")
		return fmt.Errorf(errorMsg)
	}

	// Now parse the token

	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	// Check if there was an error in parsing...
	if err != nil {
		m.logf("Error parsing token: %v", err)
		m.Config.ErrorHandler(ctx, err.Error())
		return fmt.Errorf("error parsing token: %v", err)
	}

	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		m.logf("Error validating token algorithm: %s", message)
		m.Config.ErrorHandler(ctx, errors.New(message).Error())
		return fmt.Errorf("error validating token algorithm: %s", message)
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid {
		m.logf("Token is invalid")
		m.Config.ErrorHandler(ctx, "The token isn't valid")
		return fmt.Errorf("token is invalid")
	}

	if m.Config.Expiration {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			if expired := claims.VerifyExpiresAt(time.Now().Unix(), true); !expired {
				return fmt.Errorf("token is expired")
			}
		}
	}

	m.logf("JWT: %v", parsedToken)

	// If we get here, everything worked and we can set the
	// user property in context.
	ctx.Set(m.Config.ContextKey, parsedToken)

	return nil
}

/**
 * jwtHandler
 */
func MidJwt() *Middleware {
	return New(Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}
