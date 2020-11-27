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
func (s *Status) WithDetails(details ...interface{}) {
	s.Details = append(s.Details, details...)
}

// Is each error in a chain for a match with a target value.
func (s *Status) Is(target error) bool {
	st, ok := target.(*Status)
	if !ok {
		return false
	}
	if s.Code != st.Code {
		return false
	}
	var (
		err1 *ErrorInfo
		err2 *ErrorInfo
	)
	for _, d := range s.Details {
		if e, ok := d.(*ErrorInfo); ok {
			err1 = e
			break
		}
	}
	for _, d := range st.Details {
		if e, ok := d.(*ErrorInfo); ok {
			err2 = e
			break
		}
	}
	if err1 != nil && err2 != nil &&
		err1.Reason == err2.Reason && err1.Domain == err2.Domain {
		return true
	}
	if err1 == nil && err2 == nil {
		return s.Code == st.Code
	}
	return false
}

func (s *Status) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s details = %+v", s.Code, s.Message, s.Details)
}

// ErrorInfo the cause of the error with structured details.
type ErrorInfo struct {
	Reason   string            `json:"reason"`
	Domain   string            `json:"domain"`
	Metadata map[string]string `json:"metadata"`
}

// WithMetadata .
func (e *ErrorInfo) WithMetadata(k, v string) {
	if e.Metadata == nil {
		e.Metadata = map[string]string{k: v}
	} else {
		e.Metadata[k] = v
	}
}

func (e *ErrorInfo) String() string {
	return fmt.Sprintf("error: reason = %s domain = %s metadata = %+v", e.Reason, e.Domain, e.Metadata)
}

// New returns a Status representing c and msg.
func New(code int, message string, details ...interface{}) error {
	return &Status{Code: code, Message: message, Details: details}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(code int, format string, a ...interface{}) error {
	return New(code, fmt.Sprintf(format, a...))
}
