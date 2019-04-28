package up

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpUpInfo(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(2089809)
		from = int(1)
		ip   = "127.0.0.1"
	)
	Convey("UpInfo", t, func(ctx C) {
		res, err := d.UpInfo(c, mid, from, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(res, ShouldNotBeNil)
		})
	})
}

func TestUpUpSwitch(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(2089809)
		from = int(1)
		ip   = "127.0.0.1"
	)
	Convey("UpSwitch", t, func(ctx C) {
		res, err := d.UpSwitch(c, mid, from, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(res, ShouldNotBeNil)
		})
	})
}

func TestUpSetUpSwitch(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(0)
		state = int(0)
		from  = int(0)
		ip    = ""
	)
	Convey("SetUpSwitch", t, func(ctx C) {
		res, err := d.SetUpSwitch(c, mid, state, from, ip)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(res, ShouldNotBeNil)
		})
	})
}

func TestDao_UpSpecialGroups(t *testing.T) {
	var c = context.Background()
	Convey("UpSpecialGroups", t, func(ctx C) {
		httpMock("GET", d.c.Host.API+_upSpecialGroupURI).Reply(200).JSON(`{"code":0,"data":[]}`)
		_, err := d.UpSpecialGroups(c, 2089809)
		So(err, ShouldBeNil)
	})
}

func TestUpSpecial(t *testing.T) {
	var (
		c   = context.Background()
		res map[int64]int64
		err error
	)
	Convey("UpSpecial", t, func(ctx C) {
		res, err = d.UpSpecial(c, 17)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
