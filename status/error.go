package status

import "fmt"

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

func (e *ErrorInfo) Error() string {
	return fmt.Sprintf("error: reason = %s domain = %s metadata = %+v", e.Reason, e.Domain, e.Metadata)
}
