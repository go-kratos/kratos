package http

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	imageContentType        = []string{"image/jpeg"}
	_                Render = Image{}
)

// Render http reponse render.
type Render interface {
	// Render render it to http response writer.
	Render(http.ResponseWriter) error
	// WriteContentType write content-type to http response writer.
	WriteContentType(w http.ResponseWriter)
}

// Image Image.
type Image struct {
	Body []byte
}

// WriteContentType write json ContentType.
func (j Image) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, imageContentType)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Render (JSON) writes data with json ContentType.
func (j Image) Render(w http.ResponseWriter) (err error) {
	if _, err = w.Write(j.Body); err != nil {
		err = errors.WithStack(err)
	}
	return
}
