package http

import (
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
)

const (
	baseContentType    = "application"
	defaultContentType = "application/json"
)

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

// RequestCodec returns request codec.
func RequestCodec(req *http.Request) (encoding.Codec, error) {
	contentType := req.Header.Get("content-type")
	codec := encoding.GetCodec(contentSubtype(contentType))
	if codec == nil {
		return nil, errors.InvalidArgument("Errors_UnknownCodec", contentType)
	}
	return codec, nil
}

// ResponseCodec returns response codec.
func ResponseCodec(req *http.Request) (string, encoding.Codec, error) {
	accepts := req.Header.Values("accept")
	for _, contentType := range accepts {
		if codec := encoding.GetCodec(contentSubtype(contentType)); codec != nil {
			return contentType, codec, nil
		}
	}
	if codec := encoding.GetCodec("json"); codec != nil {
		return defaultContentType, codec, nil
	}
	return "", nil, errors.InvalidArgument("Error_UnknownCodec", strings.Join(accepts, "; "))
}
