package validate

import (
	"context"
	"testing"

	"google.golang.org/protobuf/proto"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/testdata/validate"
	"github.com/go-kratos/kratos/v2/middleware"
)

type testcase struct {
	name string
	req  proto.Message
	err  bool
}

func TestTable(t *testing.T) {
	var mock middleware.Handler = func(context.Context, interface{}) (interface{}, error) { return nil, nil }

	tests := []testcase{
		{
			name: "valid_legacy",
			req:  &validate.Legacy{Name: "testcase", Age: 19},
			err:  false,
		},
		{
			name: "invalid_legacy1",
			req:  &validate.Legacy{Name: "testcase", Age: 10},
			err:  true,
		},
		{
			name: "invalid_legacy2",
			req:  &validate.Legacy{Name: "test", Age: 100},
			err:  true,
		},
		{
			name: "valid_mixed",
			req:  &validate.Mixed{Name: "testcase", Age: 19},
			err:  false,
		},
		{
			name: "invalid_mixed1",
			req:  &validate.Mixed{Name: "testcase", Age: 10},
			err:  true,
		},
		{
			name: "invalid_mixed2",
			req:  &validate.Mixed{Name: "test", Age: 100},
			err:  true,
		},
		{
			name: "valid_modern",
			req:  &validate.Modern{Name: "testcase", Age: 19},
			err:  false,
		},
		{
			name: "invalid_modern1",
			req:  &validate.Modern{Name: "testcase", Age: 10},
			err:  true,
		},
		{
			name: "invalid_modern2",
			req:  &validate.Modern{Name: "test", Age: 100},
			err:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handle := Validator()(mock)
			_, err := handle(context.Background(), test.req)
			expect := test.err
			actual := kratoserrors.IsBadRequest(err)
			if expect != actual {
				t.Errorf("case %s expect %v, actual %v, err %v", test.name, expect, actual, err)
			}
		})
	}
}
