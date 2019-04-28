package account

import (
	"context"
	"testing"

	accmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestAccountCard3(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(27515256)
	)
	convey.Convey("Card3", t, func(ctx convey.C) {
		res, err := d.Card3(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		res, err = d.Card3(c, 777777777777)
		ctx.So(err, convey.ShouldEqual, ecode.AccessDenied)
		ctx.So(res, convey.ShouldBeNil)
	})
}

func TestAccountIsVip(t *testing.T) {
	var (
		cardFalse = &accmdl.Card{}
		// card.Vip.Type == 0 || card.Vip.Status == 0 || card.Vip.Status == 2 || card.Vip.Status == 3
		cardTrue = &accmdl.Card{
			Vip: accmdl.VipInfo{
				Type:   1,
				Status: 1,
			},
		}
	)
	convey.Convey("IsVip", t, func(ctx convey.C) {
		p1 := IsVip(cardFalse)
		ctx.Convey("Then p1 should be false.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldBeFalse)
		})
		p2 := IsVip(cardTrue)
		ctx.Convey("Then p2 should be true.", func(ctx convey.C) {
			ctx.So(p2, convey.ShouldBeTrue)
		})
	})
}
