package dao

import (
	"context"
	"testing"
	xtime "time"

	"go-common/app/interface/main/web/model"
	"go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoHelpList(t *testing.T) {
	convey.Convey("HelpList", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			pTypeID = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.helpListURL).Reply(200).JSON(`{"retCode":"000000"}`)
			data, err := d.HelpList(c, pTypeID)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}

func TestDaokeyHl(t *testing.T) {
	convey.Convey("keyHl", t, func(ctx convey.C) {
		var (
			pTypeID = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyHl(pTypeID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyHd(t *testing.T) {
	convey.Convey("keyHd", t, func(ctx convey.C) {
		var (
			qTypeID = ""
			keyFlag = int(0)
			pn      = int(0)
			ps      = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyHd(qTypeID, keyFlag, pn, ps)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetHlCache(t *testing.T) {
	convey.Convey("SetHlCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			pTypeID = ""
			Hl      = []*model.HelpList{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetHlCache(c, pTypeID, Hl)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoHlCache(t *testing.T) {
	convey.Convey("HlCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			pTypeID = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.HlCache(c, pTypeID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", res)
			})
		})
	})
}

func TestDaoHelpDetail(t *testing.T) {
	convey.Convey("HelpDetail", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			qTypeID = ""
			keyFlag = int(0)
			pn      = int(0)
			ps      = int(0)
			ip      = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.helpSearchURL).Reply(200).JSON(`{"retCode":"000000", "total":1}`)
			data, total, err := d.HelpDetail(c, qTypeID, keyFlag, pn, ps, ip)
			ctx.Convey("Then err should be nil.data,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}

func TestDaoHelpSearch(t *testing.T) {
	convey.Convey("HelpSearch", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			pTypeID  = ""
			keyWords = ""
			keyFlag  = int(0)
			pn       = int(0)
			ps       = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.helpSearchURL).Reply(200).JSON(`{"retCode":"000000", "total":1}`)
			data, total, err := d.HelpSearch(c, pTypeID, keyWords, keyFlag, pn, ps)
			ctx.Convey("Then err should be nil.data,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}

func TestDaoSetDetailCache(t *testing.T) {
	convey.Convey("SetDetailCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			qTypeID = ""
			keyFlag = int(0)
			pn      = int(0)
			ps      = int(0)
			total   = int(0)
			data    = []*model.HelpDeatil{{AllTypeName: "1111"}, {AllTypeName: "2222"}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetDetailCache(c, qTypeID, keyFlag, pn, ps, total, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDetailCache(t *testing.T) {
	convey.Convey("DetailCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			qTypeID = ""
			keyFlag = int(0)
			pn      = int(0)
			ps      = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, count, err := d.DetailCache(c, qTypeID, keyFlag, pn, ps)
			ctx.Convey("Then err should be nil.res,count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
				convey.Printf("%+v", res)
			})
		})
	})
}

func TestDaofromHd(t *testing.T) {
	convey.Convey("fromHd", t, func(ctx convey.C) {
		var (
			i = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := fromHd(i)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocombineHd(t *testing.T) {
	convey.Convey("combineHd", t, func(ctx convey.C) {
		var (
			create = time.Time(xtime.Now().Unix())
			count  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := combineHd(create, count)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
