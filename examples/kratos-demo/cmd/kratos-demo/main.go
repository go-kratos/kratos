package main

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/codes"
	myerr "github.com/go-kratos/kratos/v2/examples/kratos-demo/errors"
)

func main() {
	err := error01()
	if codes.IsMissingField(err) {
		// TODO
	}
	if errors.ReasonForError(err) == codes.Kratos_MissingField.String() {
		// TODO
	}
}

func error01() error {
	return errors.InvalidArgument(
		codes.Kratos_MissingField.String(),
		"name is missing",
	)
}

func error02() error {
	return myerr.MissingField()
}

func error03() error {
	return myerr.RequestBlocked("this is a message")
}
