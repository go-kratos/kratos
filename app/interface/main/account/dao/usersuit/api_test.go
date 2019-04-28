package usersuit

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_Group(t *testing.T) {
	convey.Convey("Group", t, func() {
		groups, err := d.Group(context.TODO(), "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(groups, convey.ShouldNotBeEmpty)
	})
}
