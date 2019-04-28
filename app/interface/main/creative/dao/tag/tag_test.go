package tag

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTagTagList(t *testing.T) {
	var (
		c   = context.TODO()
		ids = []int64{}
	)
	Convey("TagList", t, func(ctx C) {
		_, err := d.TagList(c, ids)
		ctx.Convey("Then err should be nil.tgs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldNotBeNil)
		})
	})
}

func TestTagTagCheck(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(0)
		tagName = ""
	)
	Convey("TagCheck", t, func(ctx C) {
		_, err := d.TagCheck(c, mid, tagName)
		ctx.Convey("Then err should be nil.no should not be nil.", func(ctx C) {
			ctx.So(err, ShouldNotBeNil)
		})
	})
}

func TestTagAppealTag(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(123)
		ip  = ""
	)
	Convey("AppealTag", t, func(ctx C) {
		tid, err := d.AppealTag(c, aid, ip)
		ctx.Convey("Then err should be nil.tid should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(tid, ShouldNotBeNil)
		})
	})
}

func TestDao_StaffTitleList(t *testing.T) {
	var c = context.Background()
	Convey("StaffTitleList1", t, func() {
		httpMock("GET", d.mngTagListURI).Reply(200).JSON(`{"code":0,"data":{"data":[{"id":1,"name":"测试"}]}}`)
		_, err := d.StaffTitleList(c)
		So(err, ShouldBeNil)
	})
}
