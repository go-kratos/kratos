package render

import (
	"net/http"

	"github.com/pkg/errors"
)

// Data common bytes struct.
type Data struct {
	ContentType string
	Data        [][]byte
}

// Render (Data) writes data with custom ContentType.
func (r Data) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	for _, d := range r.Data {
		if _, err = w.Write(d); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	return
}

// WriteContentType writes data with custom ContentType.
func (r Data) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, []string{r.ContentType})
}
