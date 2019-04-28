package pb

import (
	"strconv"

	"go-common/library/ecode"

	any "github.com/golang/protobuf/ptypes/any"
)

func (e *Error) Error() string {
	return strconv.FormatInt(int64(e.GetErrCode()), 10)
}

// Code is the code of error.
func (e *Error) Code() int {
	return int(e.GetErrCode())
}

// Message is error message.
func (e *Error) Message() string {
	return e.GetErrMessage()
}

// Equal compare whether two errors are equal.
func (e *Error) Equal(ec error) bool {
	return ecode.Cause(ec).Code() == e.Code()
}

// Details return error details.
func (e *Error) Details() []interface{} {
	return []interface{}{e.GetErrDetail()}
}

// From will convert ecode.Codes to pb.Error.
//
// Deprecated: please use ecode.Error
func From(ec ecode.Codes) *Error {
	var detail *any.Any
	if details := ec.Details(); len(details) > 0 {
		detail, _ = details[0].(*any.Any)
	}
	return &Error{
		ErrCode:    int32(ec.Code()),
		ErrMessage: ec.Message(),
		ErrDetail:  detail,
	}
}
