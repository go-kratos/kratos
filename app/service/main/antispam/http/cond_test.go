package http_test

import (
	"errors"
	"fmt"
	"testing"

	"go-common/app/service/main/antispam/http"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	cases := []struct {
		cond        *http.Condition
		expectedErr error
	}{
		{&http.Condition{Search: " ", Order: ""}, nil},
		{&http.Condition{Search: "foo", Order: ""}, nil},
		{&http.Condition{Search: "bar", Order: "xxx"}, errors.New("Order by should be 'ASC' or 'DESC' but got(XXX)")},
		{&http.Condition{Search: "bar", Order: "asc"}, nil},
		{&http.Condition{Search: "bar", Order: "ASC"}, nil},
		{&http.Condition{Search: "bar", Order: "DESC"}, nil},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("Search(%q) Order(%q)", c.cond.Search, c.cond.Order), func(t *testing.T) {
			assert := assert.New(t)
			err := c.cond.Valid()
			assert.Equal(c.expectedErr, err, fmt.Sprintf("cond.Valid() = %v, want %v", err, c.expectedErr))
		})
	}
}
