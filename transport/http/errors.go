package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/errors"
)

var errMapping = map[int]int{
	http.StatusOK:                  0,
	http.StatusBadRequest:          3,
	http.StatusRequestTimeout:      4,
	http.StatusNotFound:            5,
	http.StatusConflict:            6,
	http.StatusForbidden:           7,
	http.StatusPreconditionFailed:  9,
	http.StatusNotImplemented:      12,
	http.StatusInternalServerError: 13,
	http.StatusServiceUnavailable:  14,
	http.StatusUnauthorized:        16,
}

func httpError(err error) *errors.StatusError {
	se, ok := err.(*errors.StatusError)
	if !ok {
		se = &errors.StatusError{
			Code:    2,
			Message: "Unknown",
		}
	}
	if code, ok := errMapping[se.Code]; ok {
		se.Code = code
	} else {
		se.Code = http.StatusInternalServerError
	}
	return se
}

// ErrUnknownCodec returns a unknown codec error.
func ErrUnknownCodec(message string) error {
	return errors.InvalidArgument("Errors_UnknownCodec", message)
}

// ErrDataLoss returns a data loss error.
func ErrDataLoss(message string) error {
	return errors.InvalidArgument("Errors_DataLoss", message)
}

// ErrCodecUnmarshal returns a codec unmarshal error.
func ErrCodecUnmarshal(message string) error {
	return errors.InvalidArgument("Errors_CodecUnmarshal", message)
}

// ErrCodecMarshal returns a codec marshal error.
func ErrCodecMarshal(message string) error {
	return errors.InvalidArgument("Errors_CodecMarshal", message)
}

// DefaultErrorHandler is default errors handler.
func DefaultErrorHandler(ctx context.Context, err error, codec encoding.Codec, w http.ResponseWriter) {
	se := httpError(err)
	w.WriteHeader(se.Code)
	if codec != nil {
		b, _ := codec.Marshal(se)
		w.Write(b)
	}
}
