package blademaster

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	allowOriginHosts = []string{
		".bilibili.com",
		".biligame.com",
		".bilibili.co",
		".im9.com",
		".acg.tv",
		".hdslb.com",
	}
)

// CORSConfig represents all available options for the middleware.
type CORSConfig struct {
	AllowAllOrigins bool

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// Default value is []
	AllowOrigins []string

	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowMethods []string

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	AllowHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposeHeaders []string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge time.Duration
}

type cors struct {
	allowAllOrigins  bool
	allowCredentials bool
	allowOriginFunc  func(string) bool
	allowOrigins     []string
	normalHeaders    http.Header
	preflightHeaders http.Header
}

type converter func(string) string

// Validate is check configuration of user defined.
func (c *CORSConfig) Validate() error {
	if c.AllowAllOrigins && (c.AllowOriginFunc != nil || len(c.AllowOrigins) > 0) {
		return errors.New("conflict settings: all origins are allowed. AllowOriginFunc or AllowedOrigins is not needed")
	}
	if !c.AllowAllOrigins && c.AllowOriginFunc == nil && len(c.AllowOrigins) == 0 {
		return errors.New("conflict settings: all origins disabled")
	}
	for _, origin := range c.AllowOrigins {
		if origin != "*" && !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
			return errors.New("bad origin: origins must either be '*' or include http:// or https://")
		}
	}
	return nil
}

// CORS returns the location middleware with default configuration.
func CORS() HandlerFunc {
	config := &CORSConfig{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           time.Duration(0),
		AllowOriginFunc: func(origin string) bool {
			for _, host := range allowOriginHosts {
				if strings.HasSuffix(strings.ToLower(origin), host) {
					return true
				}
			}
			return false
		},
	}
	return newCORS(config)
}

// newCORS returns the location middleware with user-defined custom configuration.
func newCORS(config *CORSConfig) HandlerFunc {
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}
	cors := &cors{
		allowOriginFunc:  config.AllowOriginFunc,
		allowAllOrigins:  config.AllowAllOrigins,
		allowCredentials: config.AllowCredentials,
		allowOrigins:     normalize(config.AllowOrigins),
		normalHeaders:    generateNormalHeaders(config),
		preflightHeaders: generatePreflightHeaders(config),
	}

	return func(c *Context) {
		cors.applyCORS(c)
	}
}

func (cors *cors) applyCORS(c *Context) {
	origin := c.Request.Header.Get("Origin")
	if len(origin) == 0 {
		// request is not a CORS request
		return
	}
	if !cors.validateOrigin(origin) {
		log.V(5).Info("The request's Origin header `%s` does not match any of allowed origins.", origin)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if c.Request.Method == "OPTIONS" {
		cors.handlePreflight(c)
		defer c.AbortWithStatus(200)
	} else {
		cors.handleNormal(c)
	}

	if !cors.allowAllOrigins {
		header := c.Writer.Header()
		header.Set("Access-Control-Allow-Origin", origin)
	}
}

func (cors *cors) validateOrigin(origin string) bool {
	if cors.allowAllOrigins {
		return true
	}
	for _, value := range cors.allowOrigins {
		if value == origin {
			return true
		}
	}
	if cors.allowOriginFunc != nil {
		return cors.allowOriginFunc(origin)
	}
	return false
}

func (cors *cors) handlePreflight(c *Context) {
	header := c.Writer.Header()
	for key, value := range cors.preflightHeaders {
		header[key] = value
	}
}

func (cors *cors) handleNormal(c *Context) {
	header := c.Writer.Header()
	for key, value := range cors.normalHeaders {
		header[key] = value
	}
}

func generateNormalHeaders(c *CORSConfig) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}

	// backport support for early browsers
	if len(c.AllowMethods) > 0 {
		allowMethods := convert(normalize(c.AllowMethods), strings.ToUpper)
		value := strings.Join(allowMethods, ",")
		headers.Set("Access-Control-Allow-Methods", value)
	}

	if len(c.ExposeHeaders) > 0 {
		exposeHeaders := convert(normalize(c.ExposeHeaders), http.CanonicalHeaderKey)
		headers.Set("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ","))
	}
	if c.AllowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Vary", "Origin")
	}
	return headers
}

func generatePreflightHeaders(c *CORSConfig) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.AllowMethods) > 0 {
		allowMethods := convert(normalize(c.AllowMethods), strings.ToUpper)
		value := strings.Join(allowMethods, ",")
		headers.Set("Access-Control-Allow-Methods", value)
	}
	if len(c.AllowHeaders) > 0 {
		allowHeaders := convert(normalize(c.AllowHeaders), http.CanonicalHeaderKey)
		value := strings.Join(allowHeaders, ",")
		headers.Set("Access-Control-Allow-Headers", value)
	}
	if c.MaxAge > time.Duration(0) {
		value := strconv.FormatInt(int64(c.MaxAge/time.Second), 10)
		headers.Set("Access-Control-Max-Age", value)
	}
	if c.AllowAllOrigins {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		// Always set Vary headers
		// see https://github.com/rs/cors/issues/10,
		// https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001

		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")
	}
	return headers
}

func normalize(values []string) []string {
	if values == nil {
		return nil
	}
	distinctMap := make(map[string]bool, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.ToLower(value)
		if _, seen := distinctMap[value]; !seen {
			normalized = append(normalized, value)
			distinctMap[value] = true
		}
	}
	return normalized
}

func convert(s []string, c converter) []string {
	var out []string
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}
