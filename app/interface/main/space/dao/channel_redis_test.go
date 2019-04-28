package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/space/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyCl(t *testing.T) {
	convey.Convey("keyCl", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyCl(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyClArc(t *testing.T) {
	convey.Convey("keyClArc", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
			cid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyClArc(mid, cid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyClArcSort(t *testing.T) {
	convey.Convey("keyClArcSort", t, func(ctx convey.C) {
		var (
			mid = int64(2222)
			cid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyClArcSort(mid, cid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetChannelListCache(t *testing.T) {
	convey.Convey("SetChannelListCache", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			mid         = int64(2222)
			channelList = []*model.Channel{{Cid: 2222, Mid: 2222, Name: "2222"}, {Cid: 3333, Mid: 2222, Name: "3333"}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetChannelListCache(c, mid, channelList)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetChannelCache(t *testing.T) {
	convey.Convey("SetChannelCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(2222)
			cid     = int64(2222)
			channel = &model.Channel{Cid: 2222, Mid: 2222}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetChannelCache(c, mid, cid, channel)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoChannelCache(t *testing.T) {
	convey.Convey("ChannelCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
			cid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			channel, err := d.ChannelCache(c, mid, cid)
			ctx.Convey("Then err should be nil.channel should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(channel, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelChannelCache(t *testing.T) {
	convey.Convey("DelChannelCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelChannelCache(c, mid, cid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoChannelListCache(t *testing.T) {
	convey.Convey("ChannelListCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			channels, err := d.ChannelListCache(c, mid)
			ctx.Convey("Then err should be nil.channels should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", channels)
			})
		})
	})
}

func TestDaoChannelArcsCache(t *testing.T) {
	convey.Convey("ChannelArcsCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(2222)
			cid   = int64(2222)
			start = int(0)
			end   = int(1)
			order bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			arcs, err := d.ChannelArcsCache(c, mid, cid, start, end, order)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", arcs)
			})
		})
	})
}

func TestDaoAddChannelArcCache(t *testing.T) {
	convey.Convey("AddChannelArcCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(2222)
			cid  = int64(2222)
			arcs = []*model.ChannelArc{{ID: 2222, Mid: 2222, Cid: 2222, Aid: 2222}, {ID: 3333, Mid: 3333, Cid: 3333, Aid: 3333}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddChannelArcCache(c, mid, cid, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetChannelArcSortCache(t *testing.T) {
	convey.Convey("SetChannelArcSortCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			cid  = int64(0)
			sort = []*model.ChannelArcSort{{Aid: 4444, OrderNum: 100}, {Aid: 5555, OrderNum: 200}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetChannelArcSortCache(c, mid, cid, sort)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelChannelArcCache(t *testing.T) {
	convey.Convey("DelChannelArcCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			cid = int64(0)
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelChannelArcCache(c, mid, cid, aid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelChannelArcsCache(t *testing.T) {
	convey.Convey("DelChannelArcsCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2222)
			cid = int64(2222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelChannelArcsCache(c, mid, cid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetChannelArcsCache(t *testing.T) {
	convey.Convey("SetChannelArcsCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(2222)
			cid  = int64(2222)
			arcs = []*model.ChannelArc{{ID: 2222, Mid: 2222, Cid: 2222, Aid: 2222}, {ID: 2222, Mid: 2222, Cid: 2222, Aid: 3333}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetChannelArcsCache(c, mid, cid, arcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
