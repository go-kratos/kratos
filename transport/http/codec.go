package http

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/gorilla/mux"
)

const baseContentType = "application"

func contentSubtype(contentType string) string {
	if contentType == baseContentType {
		return ""
	}
	if !strings.HasPrefix(contentType, baseContentType) {
		return ""
	}
	// guaranteed since != baseContentType and has baseContentType prefix
	switch contentType[len(baseContentType)] {
	case '/', ';':
		// this will return true for "application/grpc+" or "application/grpc;"
		// which the previous validContentType function tested to be valid, so we
		// just say that no content-subtype is specified in this case
		return contentType[len(baseContentType)+1:]
	default:
		return ""
	}
}

func codecForReq(req *http.Request) (Marshaler, error) {
	contentType := req.Header.Get("content-type")
	cc := encoding.GetCodec(contentSubtype(contentType))
	if cc == nil {
		return nil, errors.InvalidArgument("Errors_UnknownCodec", contentType)
	}
	return &codec{codec: cc, req: req}, nil
}

// Marshaler defines the interface HTTP uses to encode and decode messages.
type Marshaler interface {
	PathParams() map[string]string
	ReadHeader() http.Header
	ReadBody() ([]byte, error)
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(v interface{}) error
}

type codec struct {
	codec encoding.Codec
	req   *http.Request
}

func (c *codec) PathParams() map[string]string {
	return mux.Vars(c.req)
}

func (c *codec) ReadHeader() http.Header {
	return c.req.Header
}

func (c *codec) Marshal(v interface{}) ([]byte, error) {
	return c.codec.Marshal(v)
}

func (c *codec) ReadBody() ([]byte, error) {
	data, err := ioutil.ReadAll(c.req.Body)
	if err != nil {
		return nil, errors.DataLoss("Errors_DataLoss", err.Error())
	}
	defer c.req.Body.Close()
	return data, nil
}

func (c *codec) Unmarshal(v interface{}) error {
	data, err := c.ReadBody()
	if err != nil {
		return err
	}
	if err = c.codec.Unmarshal(data, v); err != nil {
		return errors.InvalidArgument("Errors_CodecUnmarshal", err.Error())
	}
	return nil
}
