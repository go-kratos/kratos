package render

import (
	"net/http"

	"github.com/pkg/errors"
)

// Redirect render for redirect to specified location.
type Redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

// Render (Redirect) redirect to specified location.
func (r Redirect) Render(w http.ResponseWriter) error {
	if (r.Code < 300 || r.Code > 308) && r.Code != 201 {
		return errors.Errorf("Cannot redirect with status code %d", r.Code)
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}

// WriteContentType noneContentType.
func (r Redirect) WriteContentType(http.ResponseWriter) {}
