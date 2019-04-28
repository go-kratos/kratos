package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetQidCache(t *testing.T) {
	convey.Convey("SetQidCache", t, func() {
		err := d.SetQidCache(context.TODO(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}
