package dao

import (
	"context"
	"go-common/app/interface/main/growup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetActivity(t *testing.T) {
	convey.Convey("GetActivity", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO creative_activity(id, name) VALUES(1000, 'dao-test')")
			ac, err := d.GetActivity(c, id)
			ctx.Convey("Then err should be nil.ac should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ac, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListUpActivity(t *testing.T) {
	convey.Convey("ListUpActivity", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_activity(mid,activity_id,state) VALUES(1001, 1000, 2)")
			ups, err := d.ListUpActivity(c, id)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSignUpActivity(t *testing.T) {
	convey.Convey("SignUpActivity", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			up = &model.UpBonus{
				MID:        1002,
				Nickname:   "test",
				ActivityID: 1000,
				State:      2,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "DELETE FROM up_activity WHERE mid = 1002")
			rows, err := d.SignUpActivity(c, up)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
