package main

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	apierr "github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/errors"
	myerr "github.com/go-kratos/kratos/v2/examples/kratos-demo/errors"
)

func main() {
	err := error01()
	if apierr.IsMissingField(err) {
		// TODO
	}
	if errors.Reason(err) == apierr.Kratos_MissingField.String() {
		// TODO
	}

	message := errors.ReasonMessage(err)

	fmt.Println(message)
}

func error01() error {
	return errors.InvalidArgument(
		apierr.Kratos_MissingField.String(),
		"name is missing",
	)
}

func error02() error {
	return myerr.MissingField()
}

func error03() error {
	return myerr.RequestBlocked("this is a message")
}
