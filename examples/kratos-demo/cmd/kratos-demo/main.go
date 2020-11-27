package main

import (
	"github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/codes"
	"github.com/go-kratos/kratos/v2/examples/kratos-demo/errors"
	"github.com/go-kratos/kratos/v2/status"
)

func main() {
}

func error01() error {
	return status.BadRequest(
		&status.ErrorInfo{Reason: codes.Kratos_MissingField.String()},
	)
}

func error02() error {
	return status.InternalServerError(
		&errors.ErrMissingField,
	)
}

func error03() error {
	err := errors.ErrMissingField
	err.WithMetadata("name", "empty")
	return status.InternalServerError(
		&err,
	)
}

func error04() error {
	return status.InternalServerError(
		errors.RequestBlocked("my_resources"),
	)
}
