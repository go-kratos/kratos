package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDialogAll(t *testing.T) {
	convey.Convey("DialogAll", t, func() {
		_, err := d.DialogAll(context.TODO())
		convey.So(err, convey.ShouldBeNil)
	})
}
