package lic

import (
	model "go-common/app/job/main/tv/model/pgc"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLicBuildLic(t *testing.T) {
	var (
		sign  = ""
		ps    = []*model.PS{}
		count = int(0)
	)
	convey.Convey("BuildLic", t, func(ctx convey.C) {
		p1 := BuildLic(sign, ps, count)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestLicRandStringBytesRmndr(t *testing.T) {
	var (
		n = int(0)
	)
	convey.Convey("RandStringBytesRmndr", t, func(ctx convey.C) {
		p1 := RandStringBytesRmndr(n)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestLicDelLic(t *testing.T) {
	var (
		sign   = ""
		prefix = ""
		sid    = int64(0)
	)
	convey.Convey("DelLic", t, func(ctx convey.C) {
		p1 := DelLic(sign, prefix, sid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestLicDelEpLic(t *testing.T) {
	var (
		prefix = ""
		sign   = ""
		delEps = []int{}
	)
	convey.Convey("DelEpLic", t, func(ctx convey.C) {
		p1 := DelEpLic(prefix, sign, delEps)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
