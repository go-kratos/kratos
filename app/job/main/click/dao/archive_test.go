package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_MaxAID(t *testing.T) {
	Convey("MaxAID", t, func() {
		d.MaxAID(context.TODO())
	})
}
