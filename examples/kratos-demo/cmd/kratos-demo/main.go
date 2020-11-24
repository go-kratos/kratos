package main

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	ec "github.com/go-kratos/kratos/v2/examples/kratos-demo/api/errors"
	code "github.com/go-kratos/kratos/v2/examples/kratos-demo/errors"
)

func main() {
	err := errors.BadRequest(
		errors.ErrorInfo{Reason: ec.Errors_MissingField.String()},
	)
	err = code.WithErrorInfo(ec.Errors_RequestBlocked.String(), map[string]string{"field_name": "name"})

	fmt.Println(err)
}
