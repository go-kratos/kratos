package main

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/codes"
	myerr "github.com/go-kratos/kratos/v2/examples/kratos-demo/errors"
)

func main() {
}

func error01() error {
	return errors.InvalidArgument(
		codes.Kratos_MissingField.String(),
		"name is missing",
	)
}

func error02() error {
	return myerr.ErrMissingField
}

func error03() error {
	return myerr.RequestBlocked("this is a message")
}
