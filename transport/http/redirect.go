package http

type redirect struct {
	Code int
	URL  string
}

func (r *redirect) Redirect() (int, string) {
	return r.Code, r.URL
}

// NewRedirect new a redirect with url, which may be a path relative to the request path.
// The provided code should be in the 3xx range and is usually StatusMovedPermanently, StatusFound or StatusSeeOther.
// If the Content-Type header has not been set, Redirect sets it to "text/html; charset=utf-8" and writes a small HTML body.
// Setting the Content-Type header to any value, including nil, disables that behavior.
func NewRedirect(code int, url string) Redirector {
	return &redirect{Code: code, URL: url}
}
