package blademaster

import (
	"net/url"
	"regexp"
	"strings"

	"go-common/library/log"
)

var (
	_allowHosts = []string{
		".bilibili.com",
		".bilibili.co",
		".biligame.com",
		".im9.com",
		".acg.tv",
		".hdslb.com",
	}
	_allowPatterns = []string{
		// match by wechat appid
		`^http(?:s)?://([\w\d]+\.)?servicewechat.com/(wx7564fd5313d24844|wx618ca8c24bf06c33)`,
	}

	validations = []func(*url.URL) bool{}
)

func matchHostSuffix(suffix string) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return strings.HasSuffix(strings.ToLower(uri.Host), suffix)
	}
}

func matchPattern(pattern *regexp.Regexp) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return pattern.MatchString(strings.ToLower(uri.String()))
	}
}

// addHostSuffix add host suffix into validations
func addHostSuffix(suffix string) {
	validations = append(validations, matchHostSuffix(suffix))
}

// addPattern add referer pattern into validations
func addPattern(pattern string) {
	validations = append(validations, matchPattern(regexp.MustCompile(pattern)))
}

func init() {
	for _, r := range _allowHosts {
		addHostSuffix(r)
	}
	for _, p := range _allowPatterns {
		addPattern(p)
	}
}

// CSRF returns the csrf middleware to prevent invalid cross site request.
// Only referer is checked currently.
func CSRF() HandlerFunc {
	return func(c *Context) {
		referer := c.Request.Header.Get("Referer")
		params := c.Request.Form
		cross := (params.Get("callback") != "" && params.Get("jsonp") == "jsonp") || (params.Get("cross_domain") != "")
		if referer == "" {
			if !cross {
				return
			}
			log.V(5).Info("The request's Referer header is empty.")
			c.AbortWithStatus(403)
			return
		}
		illegal := true
		if uri, err := url.Parse(referer); err == nil && uri.Host != "" {
			for _, validate := range validations {
				if validate(uri) {
					illegal = false
					break
				}
			}
		}
		if illegal {
			log.V(5).Info("The request's Referer header `%s` does not match any of allowed referers.", referer)
			c.AbortWithStatus(403)
			return
		}
	}
}
