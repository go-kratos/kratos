package dataplatform

import (
	"context"
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDataplatformsetParams(t *testing.T) {
	convey.Convey("setParams", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.setParams()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataplatformGetArchiveByMID(t *testing.T) {
	convey.Convey("GetArchiveByMID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ids, err := d.GetArchiveByMID(c, query)
			ctx.Convey("Then err should be nil.ids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataplatformSend(t *testing.T) {
	convey.Convey("Send", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.Send(c, query)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataplatformSendSpyRequest(t *testing.T) {
	convey.Convey("SendSpyRequest", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.SendSpyRequest(c, query)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataplatformSendBGMRequest(t *testing.T) {
	convey.Convey("SendBGMRequest", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			infos, err := d.SendBGMRequest(c, query)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldBeNil)
			})
		})
	})
}

func TestDataplatformNewRequest(t *testing.T) {
	convey.Convey("NewRequest", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			url    = ""
			realIP = ""
			res    = interface{}(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.NewRequest(c, url, realIP, d.setParams(), res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataplatformsign(t *testing.T) {
	convey.Convey("sign", t, func(ctx convey.C) {
		var (
			params url.Values
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			query, err := d.sign(params)
			ctx.Convey("Then err should be nil.query should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(query, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDataplatformencode(t *testing.T) {
	convey.Convey("encode", t, func(ctx convey.C) {
		var (
			v url.Values
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.encode(v)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
