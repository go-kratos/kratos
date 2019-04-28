package search

import (
	"context"
	"testing"

	"go-common/app/interface/main/creative/model/search"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchArchivesES(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(2089809)
		tid     = int16(160)
		keyword = ""
		order   = "click"
		class   = ""
		ip      = "127.0.0.1"
		pn      = int(1)
		ps      = int(10)
	)
	convey.Convey("ArchivesES", t, func(ctx convey.C) {
		sres, err := d.ArchivesES(c, mid, tid, keyword, order, class, ip, pn, ps, 1)
		ctx.Convey("Then err should be nil.sres should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(sres, convey.ShouldNotBeNil)
		})
	})
}

func TestSearchReplyES(t *testing.T) {
	var (
		c = context.TODO()
	)
	p := &search.ReplyParam{
		Ak:          "1",
		Ck:          "1",
		OMID:        1,
		OID:         1,
		Pn:          1,
		Ps:          1,
		IP:          "",
		IsReport:    1,
		Type:        1,
		FilterCtime: "",
		Kw:          "",
		Order:       "",
	}
	convey.Convey("ArchivesES", t, func(ctx convey.C) {
		_, err := d.ReplyES(c, p)
		ctx.Convey("Then err should be nil.sres should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestSearchStaffES(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(2089809)
		tid     = int16(160)
		keyword = ""
		pn      = int(1)
		ps      = int(10)
	)
	convey.Convey("ArchivesES", t, func(ctx convey.C) {
		_, err := d.ArchivesStaffES(c, mid, tid, keyword, "0", pn, ps)
		ctx.Convey("Then err should be nil.sres should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
