package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoupPrompt(t *testing.T) {
	var (
		mid = int64(1)
		ts  = time.Now().Unix()
	)
	convey.Convey("upPrompt", t, func(cv convey.C) {
		p1 := d.upPrompt(mid, ts)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaobuPrompt(t *testing.T) {
	var (
		btype = int8(1)
		mid   = int64(1)
		ts    = time.Now().Unix()
	)
	convey.Convey("buPrompt", t, func(cv convey.C) {
		p1 := d.buPrompt(btype, mid, ts)
		cv.Convey("Then p1 should not be nil.", func(cv convey.C) {
			cv.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoIncrPromptCount(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(1)
		fid   = int64(2)
		ts    = time.Now().Unix()
		btype = int8(1)
	)
	convey.Convey("IncrPromptCount", t, func(cv convey.C) {
		ucount, bcount, err := d.IncrPromptCount(c, mid, fid, ts, btype)
		cv.Convey("Then err should be nil.ucount,bcount should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(bcount, convey.ShouldNotBeNil)
			cv.So(ucount, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoClosePrompt(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(1)
		fid   = int64(2)
		ts    = time.Now().Unix()
		btype = int8(1)
	)
	convey.Convey("ClosePrompt", t, func(cv convey.C) {
		err := d.ClosePrompt(c, mid, fid, ts, btype)
		cv.Convey("Then err should be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoUpCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(1)
		fid = int64(2)
		ts  = time.Now().Unix()
	)
	convey.Convey("UpCount", t, func(cv convey.C) {
		count, err := d.UpCount(c, mid, fid, ts)
		cv.Convey("Then err should be nil.count should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBCount(t *testing.T) {
	var (
		c     = context.Background()
		mid   = int64(1)
		btype = int8(1)
		ts    = time.Now().Unix()
	)
	convey.Convey("BCount", t, func(cv convey.C) {
		count, err := d.BCount(c, mid, ts, btype)
		cv.Convey("Then err should be nil.count should not be nil.", func(cv convey.C) {
			cv.So(err, convey.ShouldBeNil)
			cv.So(count, convey.ShouldNotBeNil)
		})
	})
}
