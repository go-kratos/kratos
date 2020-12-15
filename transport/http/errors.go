package http

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
)

var errMapping = map[int]int{
	0:  http.StatusOK,
	1:  http.StatusInternalServerError,
	2:  http.StatusInternalServerError,
	3:  http.StatusBadRequest,
	4:  http.StatusRequestTimeout,
	5:  http.StatusNotFound,
	6:  http.StatusConflict,
	7:  http.StatusForbidden,
	9:  http.StatusPreconditionFailed,
	10: http.StatusConflict,
	11: http.StatusBadRequest,
	12: http.StatusNotImplemented,
	13: http.StatusInternalServerError,
	14: http.StatusServiceUnavailable,
	15: http.StatusInternalServerError,
	16: http.StatusUnauthorized,
}

// StatusError converts error to http error.
func StatusError(err error) *errors.StatusError {
	se, ok := err.(*errors.StatusError)
	if !ok {
		se = &errors.StatusError{
			Code:    2,
			Message: "Unknown: " + err.Error(),
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
func DefaultErrorHandler(ctx context.Context, err error, m Marshaler, w http.ResponseWriter) {
	se := StatusError(err)
	w.WriteHeader(se.Code)
	if m != nil {
		b, _ := m.Marshal(se)
		w.Write(b)
	}
}
