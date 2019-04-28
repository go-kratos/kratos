package render

import (
	"encoding/xml"
	"net/http"

	"github.com/pkg/errors"
)

// XML common xml struct.
type XML struct {
	Code    int
	Message string
	Data    interface{}
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

// Render (XML) writes data with xml ContentType.
func (r XML) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	if err = xml.NewEncoder(w).Encode(r.Data); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// WriteContentType write xml ContentType.
func (r XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, xmlContentType)
}
