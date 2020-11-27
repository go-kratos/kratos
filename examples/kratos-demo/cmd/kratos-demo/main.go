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
	return errors.ErrMissingField
}

func error03() error {
	return errors.RequestBlocked("my_resources")
}
