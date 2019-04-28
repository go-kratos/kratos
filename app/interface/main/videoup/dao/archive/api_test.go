package archive

import (
	"context"
	"go-common/app/interface/main/videoup/model/archive"
	upapi "go-common/app/service/main/up/api/v1"
	"go-common/library/ecode"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestPing(t *testing.T) {
	Convey("Ping", t, func(ctx C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			err := d.Ping(c)
			ctx.Convey("Then err should be nil.a,vs should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestArchiveView(t *testing.T) {
	Convey("View", t, func(ctx C) {
		var (
			c   = context.Background()
			aid = int64(10110826)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			defer gock.OffAll()
			httpMock("GET", d.viewURI).Reply(200).JSON(`{"code":20001,"data":""}`)
			a, vs, err := d.View(c, aid, ip)
			ctx.Convey("Then err should be nil.a,vs should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
				ctx.So(vs, ShouldBeNil)
				ctx.So(a, ShouldBeNil)
			})
		})
	})
}

func TestArchiveAdd(t *testing.T) {
	Convey("Add", t, func(ctx C) {
		var (
			c  = context.Background()
			ap = &archive.ArcParam{}
			ip = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			defer gock.OffAll()
			httpMock("Post", d.addURI).Reply(200).JSON(`{"code":20001,"data":""}`)
			aid, err := d.Add(c, ap, ip)
			ctx.Convey("Then err should be nil.aid should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
				ctx.So(aid, ShouldNotBeNil)
			})
		})
	})
}

func TestArchiveEdit(t *testing.T) {
	Convey("Edit", t, func(ctx C) {
		var (
			c  = context.Background()
			ap = &archive.ArcParam{}
			ip = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			defer gock.OffAll()
			httpMock("Post", d.editURI).Reply(200).JSON(`{"code":20001,"data":""}`)
			err := d.Edit(c, ap, ip)
			ctx.Convey("Then err should be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestArchiveDescFormat(t *testing.T) {
	Convey("DescFormat", t, func(ctx C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			descFormats, err := d.DescFormat(c)
			ctx.Convey("Then err should be nil.descFormats should not be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(descFormats, ShouldNotBeNil)
			})
		})
	})
}

func TestArchiveTagUp(t *testing.T) {
	Convey("TagUp", t, func(ctx C) {
		var (
			c   = context.Background()
			aid = int64(10110826)
			tag = "iamatag"
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			err := d.TagUp(c, aid, tag, ip)
			ctx.Convey("Then err should be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
			})
		})
	})
}

func TestArchivePorderCfgList(t *testing.T) {
	Convey("PorderCfgList", t, func(ctx C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			cfgs, err := d.PorderCfgList(c)
			ctx.Convey("Then err should be nil.cfgs should not be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(cfgs, ShouldNotBeNil)
			})
		})
	})
}

func TestArchiveGameList(t *testing.T) {
	Convey("GameList", t, func(ctx C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx C) {
			defer gock.OffAll()
			httpMock("GET", d.viewURI).Reply(200).JSON(`{"code":20051,"data":""}`)
			gameMap, err := d.GameList(c)
			ctx.Convey("Then err should be nil.gameMap should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
				ctx.So(gameMap, ShouldNotBeNil)
			})
		})
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

func TestDao_ApplyStaffs(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(10110826)
		ip  = "127.0.0.1"
		err error
	)
	Convey("ApplyStaffs", t, func(ctx C) {
		httpMock("GET", d.applyStaffs).Reply(200).JSON(`{"code":0}`)
		_, err = d.ApplyStaffs(c, aid, ip)
		So(err, ShouldBeNil)
	})
}

func WithMock(t *testing.T, f func(mock *gomock.Controller)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		f(mockCtrl)
	}
}

func TestDao_StaffUps(t *testing.T) {
	Convey("StaffUps", t, WithMock(t, func(mockCtrl *gomock.Controller) {
		var (
			c   = context.Background()
			err error
			ups map[int64]int64
		)
		mock := upapi.NewMockUpClient(mockCtrl)
		d.UpClient = mock
		mockReq := &upapi.UpGroupMidsReq{
			GroupID: 1,
			Pn:      1,
			Ps:      1,
		}
		mock.EXPECT().UpGroupMids(gomock.Any(), mockReq).Return(nil, ecode.CreativeAccServiceErr)
		ups, err = d.StaffUps(c)
		So(err, ShouldNotBeNil)
		So(ups, ShouldBeNil)
	}))
}

func TestDao_StaffTypeConfig(t *testing.T) {
	var (
		c   = context.Background()
		err error
	)
	Convey("StaffTypeConfig", t, func(ctx C) {
		httpMock("GET", d.staffConfigURI).Reply(200).JSON(`{"code":0,"data":{"is_gray":true,"typelist":[{"typeid":22,"max_staff":6}]}}`)
		_, _, err = d.StaffTypeConfig(c)
		So(err, ShouldBeNil)
	})
}
