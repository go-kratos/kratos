package dao

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaogetStreamNamekey(t *testing.T) {
	convey.Convey("getStreamNamekey", t, func(ctx convey.C) {
		var (
			streamName = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getStreamNamekey(streamName)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetRoomIDKey(t *testing.T) {
	convey.Convey("getRoomIDKey", t, func(ctx convey.C) {
		var (
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getRoomIDKey(rid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetRoomFieldDefaultKey(t *testing.T) {
	convey.Convey("getRoomFieldDefaultKey", t, func(ctx convey.C) {
		var (
			streamName = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getRoomFieldDefaultKey(streamName)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetRoomFieldOriginKey(t *testing.T) {
	convey.Convey("getRoomFieldOriginKey", t, func(ctx convey.C) {
		var (
			streamName = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getRoomFieldOriginKey(streamName)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetRoomFieldForwardKey(t *testing.T) {
	convey.Convey("getRoomFieldForwardKey", t, func(ctx convey.C) {
		var (
			streamName = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getRoomFieldForwardKey(streamName)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetRoomFieldSecretKey(t *testing.T) {
	convey.Convey("getRoomFieldSecretKey", t, func(ctx convey.C) {
		var (
			streamName = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getRoomFieldSecretKey(streamName)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetLastCDNKey(t *testing.T) {
	convey.Convey("getLastCDNKey", t, func(ctx convey.C) {
		var (
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getLastCDNKey(rid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetChangeSrcKey(t *testing.T) {
	convey.Convey("getChangeSrcKey", t, func(ctx convey.C) {
		var (
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.getChangeSrcKey(rid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCacheStreamFullInfo(t *testing.T) {
	convey.Convey("CacheStreamFullInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int64(11891462)
			sname = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheStreamFullInfo(c, rid, sname)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheStreamFullInfo(t *testing.T) {
	convey.Convey("AddCacheStreamFullInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			id     = int64(11891462)
			stream = &model.StreamFullInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheStreamFullInfo(c, id, stream)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheStreamRIDByName(t *testing.T) {
	convey.Convey("CacheStreamRIDByName", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sname = "live_19148701_6447624"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheStreamRIDByName(c, sname)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheStreamRIDByName(t *testing.T) {
	convey.Convey("AddCacheStreamRIDByName", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			sname  = "live_19148701_6447624"
			stream = &model.StreamFullInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheStreamRIDByName(c, sname, stream)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddCacheMultiStreamInfo(t *testing.T) {
	convey.Convey("AddCacheMultiStreamInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			res map[int64]*model.StreamFullInfo
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheMultiStreamInfo(c, res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCacheMultiStreamInfo(t *testing.T) {
	convey.Convey("CacheMultiStreamInfo", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rids = []int64{11891462}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheMultiStreamInfo(c, rids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateLastCDNCache(t *testing.T) {
	convey.Convey("UpdateLastCDNCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			rid    = int64(11891462)
			origin = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateLastCDNCache(c, rid, origin)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateChangeSrcCache(t *testing.T) {
	convey.Convey("UpdateChangeSrcCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			rid    = int64(11891462)
			origin = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateChangeSrcCache(c, rid, origin)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetLastCDNFromCache(t *testing.T) {
	convey.Convey("GetLastCDNFromCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.GetLastCDNFromCache(c, rid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetChangeSrcFromCache(t *testing.T) {
	convey.Convey("GetChangeSrcFromCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1, err := d.GetChangeSrcFromCache(c, rid)
			ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDeleteStreamByRIDFromCache(t *testing.T) {
	convey.Convey("DeleteStreamByRIDFromCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteStreamByRIDFromCache(c, rid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDeleteLastCDNFromCache(t *testing.T) {
	convey.Convey("DeleteLastCDNFromCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			rid = int64(11891462)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DeleteLastCDNFromCache(c, rid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
