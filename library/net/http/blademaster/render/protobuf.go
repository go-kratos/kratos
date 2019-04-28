package render

import (
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

var pbContentType = []string{"application/x-protobuf"}

// Render (PB) writes data with protobuf ContentType.
func (r PB) Render(w http.ResponseWriter) error {
	if r.TTL <= 0 {
		r.TTL = 1
	}
	return writePB(w, r)
}

// WriteContentType write protobuf ContentType.
func (r PB) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, pbContentType)
}

func writePB(w http.ResponseWriter, obj PB) (err error) {
	var pbBytes []byte
	writeContentType(w, pbContentType)

	if pbBytes, err = proto.Marshal(&obj); err != nil {
		err = errors.WithStack(err)
		return
	}

	if _, err = w.Write(pbBytes); err != nil {
		err = errors.WithStack(err)
	}
	return
}
