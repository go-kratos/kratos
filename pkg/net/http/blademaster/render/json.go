package render

import (
	"encoding/json"
	"html/template"
	"net/http"
	"reflect"
	"unsafe"

	"github.com/pkg/errors"
)

var jsonContentType = []string{"application/json; charset=utf-8"}
var jsonpContentType = []string{"application/javascript; charset=utf-8"}

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

// JsonpJSON contains the given interface object its callback.
type JsonpJSON struct {
	Callback string
	Data     JSON
}

// Render (JsonpJSON) marshals the given interface object and writes it and its callback with custom ContentType.
func (r JsonpJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		_, err = w.Write(ret)
		return err
	}

	callback := template.JSEscapeString(r.Callback)
	_, err = w.Write(StringToBytes(callback))
	if err != nil {
		return err
	}
	_, err = w.Write(StringToBytes("("))
	if err != nil {
		return err
	}
	_, err = w.Write(ret)
	if err != nil {
		return err
	}
	_, err = w.Write(StringToBytes(");"))
	if err != nil {
		return err
	}

	return nil
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

// WriteContentType (JsonpJSON) writes Javascript ContentType.
func (r JsonpJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpContentType)
}
