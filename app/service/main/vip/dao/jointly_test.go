package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoJointlys(t *testing.T) {
	var (
		c   = context.TODO()
		now = time.Now().Unix()
	)
	convey.Convey("DaoJointlys", t, func(ctx convey.C) {
		_, err := d.Jointlys(c, now)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
