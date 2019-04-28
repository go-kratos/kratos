package dao

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetOrder(t *testing.T) {
	Convey("Get Order", t, func() {
		_, err := d.GetOrder(context.TODO(), 10001)
		So(err, ShouldBeNil)
	})
}
