package v2

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV2UpdateConfValue(t *testing.T) {
	var (
		ID    = int64(0)
		value = ""
	)
	convey.Convey("UpdateConfValue", t, func(ctx convey.C) {
		err := d.UpdateConfValue(ID, value)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestV2UpdateConfState(t *testing.T) {
	var (
		ID     = int64(855)
		state  = int8(0)
		state2 = int8(2)
	)
	convey.Convey("UpdateConfState", t, func(ctx convey.C) {
		err := d.UpdateConfState(ID, state)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	convey.Convey("UpdateConfState restoration", t, func(ctx convey.C) {
		err := d.UpdateConfState(ID, state2)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestV2ConfigsByIDs(t *testing.T) {
	var (
		ids = []int64{855, 788}
	)
	convey.Convey("ConfigsByIDs", t, func(ctx convey.C) {
		confs, err := d.ConfigsByIDs(ids)
		ctx.Convey("Then err should be nil.confs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(confs, convey.ShouldNotBeNil)
		})
	})
}
