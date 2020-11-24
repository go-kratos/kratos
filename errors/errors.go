package errors

import "fmt"

// Error contains an error response from the server.
type Error struct {
	// Code is the HTTP response status code and will always be populated.
	Code int `json:"code"`
	// Message is the server response message and is only populated when
	// explicitly referenced by the JSON server response.
	Message string `json:"message"`
	// Details provide more context to an error.
	Details []interface{} `json:"details"`
}

// WithDetails provided details messages appended to the errors.
func (e *Error) WithDetails(details ...interface{}) {
	e.Details = append(e.Details, details...)
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s details = %+v", e.Code, e.Message, e.Details)
}

// ErrorInfo the cause of the error with structured details.
type ErrorInfo struct {
	Reason   string            `json:"reason"`
	Domain   string            `json:"domain"`
	Metadata map[string]string `json:"metadata"`
}

func (e *ErrorInfo) Error() string {
	return fmt.Sprintf("error: reason = %s domain = %s metadata = %+v", e.Reason, e.Domain, e.Metadata)
}

// New returns a Status representing c and msg.
func New(code int, message string, details ...interface{}) error {
	return &Error{Code: code, Message: message, Details: details}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(code int, format string, a ...interface{}) error {
	return New(code, fmt.Sprintf(format, a...))
}
