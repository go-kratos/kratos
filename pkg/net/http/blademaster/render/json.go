package render

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// JSON common json struct.
type JSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TTL     int         `json:"ttl"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, obj interface{}) (err error) {
	var jsonBytes []byte
	writeContentType(w, jsonContentType)
	if jsonBytes, err = json.Marshal(obj); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = w.Write(jsonBytes); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// Render (JSON) writes data with json ContentType.
func (r JSON) Render(w http.ResponseWriter) error {
	// FIXME(zhoujiahui): the TTL field will be configurable in the future
	if r.TTL <= 0 {
		r.TTL = 1
	}
	return writeJSON(w, r)
}

// WriteContentType write json ContentType.
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// MapJSON common map json struct.
type MapJSON map[string]interface{}

// Render (MapJSON) writes data with json ContentType.
func (m MapJSON) Render(w http.ResponseWriter) error {
	return writeJSON(w, m)
}

// WriteContentType write json ContentType.
func (m MapJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
