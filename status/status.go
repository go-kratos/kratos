package status

import (
	"fmt"
)

// Status contains an error response from the server.
type Status struct {
	// Code is the HTTP response status code and will always be populated.
	Code int `json:"code"`
	// Message is the server response message and is only populated when
	// explicitly referenced by the JSON server response.
	Message string `json:"message"`
	// Details provide more context to an error.
	Details []interface{} `json:"details"`
}

// WithDetails provided details messages appended to the errors.
func (e *Status) WithDetails(details ...interface{}) {
	e.Details = append(e.Details, details...)
}

// Is each error in a chain for a match with a target value.
func (e *Status) Is(target error) bool {
	te, ok := target.(*ErrorInfo)
	if !ok {
		return false
	}
	for _, d := range e.Details {
		if err, ok := d.(*ErrorInfo); ok {
			return err.Reason == te.Reason && err.Domain == te.Domain
		}
	}
	return false
}

func (e *Status) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s details = %+v", e.Code, e.Message, e.Details)
}

// New returns a Status representing c and msg.
func New(code int, message string, details ...interface{}) error {
	return &Status{Code: code, Message: message, Details: details}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(code int, format string, a ...interface{}) error {
	return New(code, fmt.Sprintf(format, a...))
}
