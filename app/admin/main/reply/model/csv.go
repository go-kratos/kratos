package model

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// CSV Comma-Separated Values struct.
type CSV struct {
	Data  [][]string
	Title string
}

var csvContentType = []string{"text/csv; charset=utf-8"}

// Render (CSV) writes data with CSV ContentType.
func (c CSV) Render(w http.ResponseWriter) (err error) {
	c.WriteContentType(w)
	writer := csv.NewWriter(w)
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	writer.Write([]string{string(bomUtf8[:])})
	writer.WriteAll(c.Data)
	if err = writer.Error(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func writeContentType(w http.ResponseWriter, value []string, title string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
	if title != "" {
		header["Content-Disposition"] = append(header["Content-Disposition"], fmt.Sprintf("attachment; filename=%s.csv", title))
	}
}

// WriteContentType write CSV ContentType.
func (c CSV) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, csvContentType, c.Title)
}
