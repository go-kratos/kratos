package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
)

const baseContentType = "application"

func contentType(subtype string) string {
	return strings.Join([]string{baseContentType, subtype}, "/")
}

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
		if i := strings.Index(contentType, ";"); i != -1 {
			return contentType[len(baseContentType)+1 : i]
		}
		return contentType[len(baseContentType)+1:]
	default:
		return ""
	}
}

func defaultRequestDecoder(req *http.Request, v interface{}) error {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	subtype := contentSubtype(req.Header.Get("content-type"))
	codec := encoding.GetCodec(subtype)
	if codec == nil {
		return fmt.Errorf("decoding request failed unknown content-type: %s", subtype)
	}
	return codec.Unmarshal(data, v)
}

func defaultResponseEncoder(res http.ResponseWriter, req *http.Request, v interface{}) error {
	subtype := contentSubtype(req.Header.Get("accept"))
	codec := encoding.GetCodec(subtype)
	if codec == nil {
		codec = encoding.GetCodec("json")
	}
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	res.Header().Set("content-type", contentType(codec.Name()))
	res.Write(data)
	return nil
}

func defaultErrorEncoder(res http.ResponseWriter, req *http.Request, err error) {
	se, code := StatusError(err)
	subtype := contentSubtype(req.Header.Get("accept"))
	codec := encoding.GetCodec(subtype)
	if codec == nil {
		codec = encoding.GetCodec("json")
	}
	data, err := codec.Marshal(se)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", contentType(codec.Name()))
	res.WriteHeader(code)
	res.Write(data)
}
