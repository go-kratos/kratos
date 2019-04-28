package dao

import (
	"go-common/app/job/main/relation/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpdateFollowingCache(t *testing.T) {
	var (
		r = &model.Relation{}
	)
	convey.Convey("UpdateFollowingCache", t, func(ctx convey.C) {
		err := d.UpdateFollowingCache(r)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
