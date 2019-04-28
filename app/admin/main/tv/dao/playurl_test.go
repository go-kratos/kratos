package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPlayurl(t *testing.T) {
	var (
		ctx       = context.Background()
		cid       = int(0)
		normalStr = `{"Code":0,"Durl":[{"URL":"test"}]}`
		httpStr   = `{"Code":-400,"Durl":[{"URL":"test"}]}`
		emptyStr  = `{"Code":0,"Durl":[]}`
	)
	convey.Convey("Playurl", t, func(cx convey.C) {
		cx.Convey("Normal Situation. Then err should be nil.playurl should not be nil.", func(cx convey.C) {
			httpMock("GET", d.c.Cfg.PlayurlAPI).Reply(200).JSON(normalStr)
			playurl, err := d.Playurl(ctx, cid)
			cx.So(err, convey.ShouldBeNil)
			cx.So(playurl, convey.ShouldNotBeNil)
		})
		cx.Convey("Http Err Situation. Then err should be nil.playurl should not be nil.", func(cx convey.C) {
			httpMock("GET", d.c.Cfg.PlayurlAPI).Reply(200).JSON(httpStr)
			_, err := d.Playurl(ctx, cid)
			cx.So(err, convey.ShouldNotBeNil)
		})
		cx.Convey("Empty Durl Situation. Then err should be nil.playurl should not be nil.", func(cx convey.C) {
			httpMock("GET", d.c.Cfg.PlayurlAPI).Reply(200).JSON(emptyStr)
			_, err := d.Playurl(ctx, cid)
			cx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUPlayurl(t *testing.T) {
	var (
		ctx       = context.Background()
		cid       = int(0)
		normalStr = `{"Code":0,"Durl":[{"URL":"test"}]}`
		httpStr   = `{"Code":-400,"Durl":[{"URL":"test"}]}`
		emptyStr  = `{"Code":0,"Durl":[]}`
	)
	convey.Convey("UPlayurl", t, func(cx convey.C) {
		cx.Convey("Normal Situation. Then err should be nil.playurl should not be nil.", func(cx convey.C) {
			httpMock("GET", d.c.Cfg.UPlayurlAPI).Reply(200).JSON(normalStr)
			playurl, err := d.UPlayurl(ctx, cid)
			cx.So(err, convey.ShouldBeNil)
			cx.So(playurl, convey.ShouldNotBeNil)
		})
		cx.Convey("Http Err Situation. Then err should be nil.playurl should not be nil.", func(cx convey.C) {
			httpMock("GET", d.c.Cfg.UPlayurlAPI).Reply(200).JSON(httpStr)
			_, err := d.UPlayurl(ctx, cid)
			cx.So(err, convey.ShouldNotBeNil)
		})
		cx.Convey("Empty Durl Situation. Then err should be nil.playurl should not be nil.", func(cx convey.C) {
			httpMock("GET", d.c.Cfg.UPlayurlAPI).Reply(200).JSON(emptyStr)
			_, err := d.UPlayurl(ctx, cid)
			cx.So(err, convey.ShouldNotBeNil)
		})
	})
}
