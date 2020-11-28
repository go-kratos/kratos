package codes

import (
	"github.com/go-kratos/kratos/v2/errors"
)

// IsRequestBlocked determines if err is an error which indicates that the request failed due to an RequestBlocked.
func IsRequestBlocked(err error) bool {
	return errors.Reason(err) == Kratos_RequestBlocked.String()
}

// IsMissingField determines if err is an error which indicates that the request failed due to an MissingField.
func IsMissingField(err error) bool {
	return errors.Reason(err) == Kratos_MissingField.String()
}
