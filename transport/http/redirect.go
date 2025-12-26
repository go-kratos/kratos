package http

var (
	_ error      = (*redirect)(nil)
	_ Redirector = (*redirect)(nil)
)

type redirect struct {
	URL  string
	Code int
}

func (r *redirect) Redirect() (string, int) {
	return r.URL, r.Code
}

func (r *redirect) Error() string {
	return "redirect to " + r.URL
}

// NewRedirect new a redirect with url, which may be a path relative to the request path.
// The provided code should be in the 3xx range and is usually StatusMovedPermanently, StatusFound or StatusSeeOther.
// If the Content-Type header has not been set, Redirect sets it to "text/html; charset=utf-8" and writes a small HTML body.
// Setting the Content-Type header to any value, including nil, disables that behavior.
func NewRedirect(url string, code int) *redirect {
	return &redirect{URL: url, Code: code}
}
