package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/protobuf/types/known/anypb"
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

// DefaultRequestDecoder default request decoder.
func DefaultRequestDecoder(ctx context.Context, in interface{}, req *http.Request) error {
	codec, err := RequestCodec(req)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.DataLoss("Errors_DataLoss", err.Error())
	}
	defer req.Body.Close()
	if err = codec.Unmarshal(data, in); err != nil {
		return errors.InvalidArgument("Errors_CodecUnmarshal", err.Error())
	}
	return nil
}

// DefaultResponseEncoder is default response encoder.
func DefaultResponseEncoder(ctx context.Context, out interface{}, res http.ResponseWriter, req *http.Request) error {
	contentType, codec, err := ResponseCodec(req)
	if err != nil {
		return err
	}
	data, err := codec.Marshal(out)
	if err != nil {
		return err
	}
	res.Header().Set("content-type", contentType)
	res.Write(data)
	return nil
}

// DefaultErrorEncoder is default errors encoder.
func DefaultErrorEncoder(ctx context.Context, err error, res http.ResponseWriter, req *http.Request) {
	status, se := StatusError(err)
	e := &Error{
		Error: &Error_Status{
			Code:    se.Code,
			Message: se.Message,
		},
	}
	for _, detail := range se.Details {
		if any, err := anypb.New(detail); err == nil {
			e.Error.Details = append(e.Error.Details, any)
		}
	}
	contentType, codec, err := ResponseCodec(req)
	if err != nil {
		data, _ := json.Marshal(se)
		res.Header().Set("content-type", contentType)
		res.WriteHeader(status)
		res.Write(data)
	} else {
		data, _ := codec.Marshal(e)
		res.Header().Set("content-type", contentType)
		res.WriteHeader(status)
		res.Write(data)
	}
}
