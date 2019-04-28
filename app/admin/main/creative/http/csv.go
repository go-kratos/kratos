package http

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// CSVContentType
var (
	CSVContentType        = []string{"application/csv"}
	_              Render = CSV{}
)

// Render http reponse render.
type Render interface {
	Render(http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
}

// CSV str.
type CSV struct {
	Content []byte
	Title   string
}

// WriteContentType fn
func (j CSV) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, CSVContentType, j.Title)
}

func writeContentType(w http.ResponseWriter, value []string, title string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
	header["Content-Disposition"] = append(header["Content-Disposition"], fmt.Sprintf("attachment; filename=%s.csv", title))
}

// Render (JSON) writes data with json ContentType.
func (j CSV) Render(w http.ResponseWriter) (err error) {
	if _, err = w.Write(j.Content); err != nil {
		err = errors.WithStack(err)
	}
	return
}
