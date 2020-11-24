package main

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	ec "github.com/go-kratos/kratos/v2/examples/kratos-demo/api/errors"
)

func main() {
	err := errors.BadRequest(
		errors.ErrorInfo{Reason: ec.Errors_MissingField.String()},
	)

	fmt.Println(err)
}
