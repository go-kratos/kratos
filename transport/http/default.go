package http

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	code, se := StatusError(err)
	contentType, codec, err := ResponseCodec(req)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	gs := status.Newf(codes.Code(se.Code), "%s: %s", se.Reason, se.Message)
	gs, err = gs.WithDetails(&errdetails.ErrorInfo{
		Reason:   se.Reason,
		Metadata: map[string]string{"message": se.Message},
	})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := codec.Marshal(gs.Proto())
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", contentType)
	res.WriteHeader(code)
	res.Write(data)
}
