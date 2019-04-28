package dao

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/library/ecode"
	"testing"
)

func TestDao_GetMetal(t *testing.T) {
	Convey("GetMetal", t, func() {
		once.Do(startService)
		var uid int64 = 1
		metal, err := d.GetMetal(ctx, uid)
		So(err, ShouldBeNil)
		So(metal, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func TestDao_ModifyMetal(t *testing.T) {
	Convey("ModifyMetal", t, func() {
		once.Do(startService)

		var uid int64 = 1
		metal, err := d.GetMetal(ctx, uid)
		So(err, ShouldBeNil)

		if metal < 10 {
			d.ModifyMetal(ctx, uid, 10, 0, nil)
		}

		var coins int64 = -5

		success, code, err := d.ModifyMetal(ctx, uid, coins, 500, "ut")
		So(code, ShouldEqual, 0)
		So(success, ShouldEqual, true)
		So(err, ShouldBeNil)

		coins = 0 - coins

		success, code, err = d.ModifyMetal(ctx, uid, coins, 0, nil)
		So(code, ShouldEqual, 0)
		So(success, ShouldEqual, true)
		So(err, ShouldBeNil)

		nmetal, err := d.GetMetal(ctx, uid)
		So(err, ShouldBeNil)
		So(metal, ShouldEqual, nmetal)
	})

	Convey("ModifyMetal not enough", t, func() {
		once.Do(startService)

		var uid int64 = 1
		metal, err := d.GetMetal(ctx, uid)
		So(err, ShouldBeNil)

		coins := 0 - int64(metal+1)
		_, _, err = d.ModifyMetal(ctx, uid, coins, 400, nil)
		So(err, ShouldEqual, ecode.CoinNotEnough)

	})
}
