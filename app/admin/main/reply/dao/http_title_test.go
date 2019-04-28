package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoNoticeTitle(t *testing.T) {
	convey.Convey("NoticeTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.NoticeTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBanTitle(t *testing.T) {
	convey.Convey("BanTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.BanTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCreditTitle(t *testing.T) {
	convey.Convey("CreditTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.CreditTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(link, convey.ShouldNotBeNil)
				ctx.So(title, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLiveVideoTitle(t *testing.T) {
	convey.Convey("LiveVideoTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.LiveVideoTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldBeBlank)
				ctx.So(title, convey.ShouldBeBlank)
			})
		})
	})
}

func TestDaoLiveActivityTitle(t *testing.T) {
	convey.Convey("LiveActivityTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.LiveActivityTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldBeBlank)
				ctx.So(title, convey.ShouldBeBlank)
			})
		})
	})
}

func TestDaoLiveNoticeTitle(t *testing.T) {
	convey.Convey("LiveNoticeTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.LiveNoticeTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldBeBlank)
				ctx.So(title, convey.ShouldBeBlank)
			})
		})
	})
}

func TestDaoLivePictureTitle(t *testing.T) {
	convey.Convey("LivePictureTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.LivePictureTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldBeBlank)
				ctx.So(title, convey.ShouldBeBlank)
			})
		})
	})
}

func TestDaoTopicTitle(t *testing.T) {
	convey.Convey("TopicTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.TopicTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				if err == nil {
					ctx.So(err, convey.ShouldBeNil)
				}
				ctx.So(link, convey.ShouldBeBlank)
				ctx.So(title, convey.ShouldBeBlank)
			})
		})
	})
}

func TestDaoTopicsLink(t *testing.T) {
	convey.Convey("TopicsLink", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			links   map[int64]string
			isTopic bool
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := _d.TopicsLink(c, links, isTopic)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoActivitySub(t *testing.T) {
	convey.Convey("ActivitySub", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.ActivitySub(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldBeBlank)
				ctx.So(title, convey.ShouldBeBlank)
			})
		})
	})
}

func TestDaoDynamicTitle(t *testing.T) {
	convey.Convey("DynamicTitle", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			title, link, err := _d.DynamicTitle(c, oid)
			ctx.Convey("Then err should be nil.title,link should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(link, convey.ShouldBeBlank)
				ctx.So(title, convey.ShouldBeBlank)
			})
		})
	})
}
