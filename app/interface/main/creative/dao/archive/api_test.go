package archive

import (
	"context"
	"encoding/json"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/ecode"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestFlowJudge(t *testing.T) {
	var (
		c        = context.TODO()
		err      error
		business = int64(1)
		groupID  = int64(2)
		oids     = []int64{1, 2, 3}
		hitOids  []int64
	)
	convey.Convey("FlowJudge", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.flowJudge).Reply(200).JSON(`{"code":20001}`)
		hitOids, err = d.FlowJudge(c, business, groupID, oids)
		ctx.Convey("Then err should be nil.hitOids should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(hitOids, convey.ShouldBeNil)
		})
	})
}

func TestArchiveSimpleArchive(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110817)
		ip  = "127.0.0.1"
		err error
		sa  *archive.SpArchive
		res struct {
			Code int                `json:"code"`
			Data *archive.SpArchive `json:"data"`
		}
	)
	convey.Convey("4", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.simpleArchive).Reply(-502)
		sa, err = d.SimpleArchive(c, aid, ip)
		ctx.Convey("Then err should be nil.sa should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.CreativeArchiveAPIErr)
			ctx.So(sa, convey.ShouldBeNil)
		})
	})
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.simpleArchive).Reply(200).JSON(`{"code":20001}`)
		sa, err = d.SimpleArchive(c, aid, ip)
		ctx.Convey("Then err should be nil.sa should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.CreativeArchiveAPIErr)
			ctx.So(sa, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res.Code = 0
		res.Data = &archive.SpArchive{
			Aid:   aid,
			Title: "iamtitle",
			Mid:   2089809,
		}
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("GET", d.simpleArchive).Reply(200).JSON(string(js))
		sa, err = d.SimpleArchive(c, aid, ip)
		ctx.Convey("Then err should be nil.sa should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(sa, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		res.Code = 0
		res.Data = nil
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("GET", d.simpleArchive).Reply(200).JSON(string(js))
		sa, err = d.SimpleArchive(c, aid, ip)
		ctx.Convey("Then err should be nil.sa should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(sa, convey.ShouldBeNil)
		})
	})
}

func TestArchiveSimpleVideos(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10110817)
		ip  = "127.0.0.1"
		res struct {
			Code int                `json:"code"`
			Data []*archive.SpVideo `json:"data"`
		}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		js, _ := json.Marshal(res)
		httpMock("GET", d.simpleVideos).Reply(200).JSON(string(js))
		vs, err := d.SimpleVideos(c, aid, ip)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vs, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 20001
		res.Data = nil
		js, _ := json.Marshal(res)
		httpMock("GET", d.simpleVideos).Reply(200).JSON(string(js))
		vs, err := d.SimpleVideos(c, aid, ip)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(vs, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 0
		res.Data = append(res.Data, &archive.SpVideo{
			Cid:   1,
			Index: 1,
			Title: "1title",
		}, &archive.SpVideo{
			Cid:   2,
			Index: 2,
			Title: "2title",
		})
		js, _ := json.Marshal(res)
		httpMock("GET", d.simpleVideos).Reply(200).JSON(string(js))
		vs, err := d.SimpleVideos(c, aid, ip)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vs, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveView(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		aid = int64(10110817)
		ip  = "127.0.0.1"
	)
	convey.Convey("View", t, func(ctx convey.C) {
		av, err := d.View(c, mid, aid, ip, 0, 0)
		ctx.Convey("Then err should be nil.av should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(av, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveViews(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(2089809)
		aids = []int64{10110817, 10110816}
		ip   = "127.0.0.1"
	)
	convey.Convey("Views", t, func(ctx convey.C) {
		avm, err := d.Views(c, mid, aids, ip)
		ctx.Convey("Then err should be nil.avm should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(avm, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveDel(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		aid = int64(10110817)
		ip  = "127.0.0.1"
	)
	convey.Convey("Del", t, func(ctx convey.C) {
		err := d.Del(c, mid, aid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.ArchiveAlreadyDel)
		})
	})
}

func TestArchiveVideoByCid(t *testing.T) {
	var (
		c   = context.TODO()
		cid = int64(10134702)
		ip  = "127.0.0.1"
	)
	convey.Convey("VideoByCid", t, func(ctx convey.C) {
		v, err := d.VideoByCid(c, cid, ip)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(v, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUpArchives(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(2089809)
		pn    = int64(1)
		ps    = int64(10)
		group = int64(0)
		ip    = "127.0.0.1"
	)
	convey.Convey("UpArchives", t, func(ctx convey.C) {
		aids, count, err := d.UpArchives(c, mid, pn, ps, group, ip)
		ctx.Convey("Then err should be nil.aids,count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
			ctx.So(aids, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveDescFormat(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("DescFormat", t, func(ctx convey.C) {
		descs, err := d.DescFormat(c)
		ctx.Convey("Then err should be nil.descs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(descs, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveVideoJam(t *testing.T) {
	var (
		c  = context.TODO()
		ip = "127.0.0.1"
	)
	convey.Convey("VideoJam", t, func(ctx convey.C) {
		level, err := d.VideoJam(c, ip)
		ctx.Convey("Then err should be nil.level should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(level, convey.ShouldNotBeNil)
		})
	})
}
