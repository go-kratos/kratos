package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/account/model"
	mc "go-common/library/cache/memcache"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyInfo(t *testing.T) {
	var (
		mid = int64(2205)
	)
	convey.Convey("Generate info-key", t, func(ctx convey.C) {
		p1 := keyInfo(mid)
		ctx.Convey("Then info-key should contains info prefix.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldContainSubstring, _prefixInfo)
		})
	})
}

func TestDaokeyCard(t *testing.T) {
	var (
		mid = int64(2205)
	)
	convey.Convey("Generate card-info-key", t, func(ctx convey.C) {
		p1 := keyCard(mid)
		ctx.Convey("Then card-info-key should contains card prefix.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldContainSubstring, _prefixCard)
		})
	})
}

func TestDaokeyVip(t *testing.T) {
	var (
		mid = int64(2205)
	)
	convey.Convey("Generate vip-info-key", t, func(ctx convey.C) {
		p1 := keyVip(mid)
		ctx.Convey("Then vip-info-key should contains vip prefix.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldContainSubstring, _prefixVip)
		})
	})
}

func TestDaokeyProfile(t *testing.T) {
	var (
		mid = int64(2205)
	)
	convey.Convey("Generate profile-key", t, func(ctx convey.C) {
		p1 := keyProfile(mid)
		ctx.Convey("Then profile-key should contains profile prefix.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldContainSubstring, _prefixProfile)
		})
	})
}

func TestDaoCacheInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2205)
	)
	convey.Convey("Get member base-info from cache", t, func(ctx convey.C) {
		_, err := d.CacheInfo(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddCacheInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2205)
		v   = &model.Info{
			Mid:  2205,
			Name: "Haha",
			Sex:  "男",
			Face: "http://i1.hdslb.com/bfs/face/4b12a3e65d344e31a11e6425767863019738c7bc.jpg",
			Sign: "来电只是",
			Rank: 500,
		}
	)
	convey.Convey("Add member base-info to cache", t, func(ctx convey.C) {
		err := d.AddCacheInfo(c, mid, v)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCacheInfos(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{2205, 2805}
	)
	convey.Convey("Batch get members' base-info", t, func(ctx convey.C) {
		res, err := d.CacheInfos(c, mids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddCacheInfos(t *testing.T) {
	var (
		c  = context.TODO()
		im = map[int64]*model.Info{
			2205: {
				Mid:  2205,
				Name: "板桥真菜",
				Sex:  "2",
				Face: "/bfs/face/e93098c3aa8c18b24001740e707ebe2df180f5f7.jpg",
				Sign: "没有",
				Rank: 10000,
			},
			3305: {
				Mid:  3305,
				Name: "FGNB",
				Sex:  "1",
				Face: "/bfs/face/e93098c3aa8c18b24001740e707ebe2df180f5f7.jpg",
				Sign: "啦啦",
				Rank: 5000,
			},
		}
	)
	convey.Convey("Batch set members' base-info to cache", t, func(ctx convey.C) {
		err := d.AddCacheInfos(c, im)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCacheCard(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2805)
	)
	convey.Convey("Get card-info from cache", t, func(ctx convey.C) {
		_, err := d.CacheCard(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddCacheCard(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2805)

		v = &model.Card{
			Mid:     10920044,
			Name:    "冠冠爱看书",
			Sex:     "男",
			Face:    "http://i1.hdslb.com/bfs/face/4b12a3e65d344e31a11e6425767863019738c7bc.jpg",
			Sign:    "来点字",
			Rank:    10000,
			Level:   5, //等级
			Silence: 0,
			Vip: model.VipInfo{
				Type:    2,
				Status:  1,
				DueDate: 162930240,
			},
			Pendant: model.PendantInfo{
				Pid:    159,
				Name:   "2018拜年祭",
				Image:  "http://i2.hdslb.com/bfs/face/aace621fa64a698f2ca94d13645a26e9a7a99ed2.png",
				Expire: 1566367231,
			},
			Nameplate: model.NameplateInfo{
				Nid:        7,
				Name:       "见习搬运工",
				Image:      "http://i1.hdslb.com/bfs/face/8478fb7c54026cd47f09daa493a1b1683113a90d.png",
				ImageSmall: "http://i0.hdslb.com/bfs/face/50eef47c3a30a75659d3cc298cfb09031d1a2ce5.png",
				Level:      "普通勋章",
				Condition:  "转载视频",
			},
		}
	)
	convey.Convey("Add card-info to cache", t, func(ctx convey.C) {
		err := d.AddCacheCard(c, mid, v)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCacheCards(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{110017381, 110019061, 110020081}
	)
	convey.Convey("Batch get card-info from cache", t, func(ctx convey.C) {
		res, err := d.CacheCards(c, mids)
		ctx.Convey("Then err should be nil and res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddCacheCards(t *testing.T) {
	var (
		c     = context.TODO()
		card1 = &model.Card{
			Mid:     10920044,
			Name:    "冠冠爱看书",
			Sex:     "男",
			Face:    "http://i1.hdslb.com/bfs/face/4b12a3e65d344e31a11e6425767863019738c7bc.jpg",
			Sign:    "来点字",
			Rank:    10000,
			Level:   5, //等级
			Silence: 0,
			Vip: model.VipInfo{
				Type:    2,
				Status:  1,
				DueDate: 162930240,
			},
			Pendant: model.PendantInfo{
				Pid:    159,
				Name:   "2018拜年祭",
				Image:  "http://i2.hdslb.com/bfs/face/aace621fa64a698f2ca94d13645a26e9a7a99ed2.png",
				Expire: 1566367231,
			},
			Nameplate: model.NameplateInfo{
				Nid:        7,
				Name:       "见习搬运工",
				Image:      "http://i1.hdslb.com/bfs/face/8478fb7c54026cd47f09daa493a1b1683113a90d.png",
				ImageSmall: "http://i0.hdslb.com/bfs/face/50eef47c3a30a75659d3cc298cfb09031d1a2ce5.png",
				Level:      "普通勋章",
				Condition:  "转载视频",
			},
		}
		cm = map[int64]*model.Card{
			card1.Mid: card1,
		}
	)
	convey.Convey("Batch set card-info to cache", t, func(ctx convey.C) {
		err := d.AddCacheCards(c, cm)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCacheVip(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110003731)
	)
	convey.Convey("Get vip-info from cache", t, func(ctx convey.C) {
		_, err := d.CacheVip(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddCacheVip(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110003731)
		v   = &model.VipInfo{
			Type:    2,
			Status:  1,
			DueDate: 162930240,
		}
	)
	convey.Convey("Set vip-cache to cache", t, func(ctx convey.C) {
		err := d.AddCacheVip(c, mid, v)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCacheVips(t *testing.T) {
	var (
		c    = context.TODO()
		mids = []int64{110002741, 110004601, 110006251}
	)
	convey.Convey("Batch get vip-infos from cache", t, func(ctx convey.C) {
		res, err := d.CacheVips(c, mids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAddCacheVips(t *testing.T) {
	var (
		c  = context.TODO()
		vm = map[int64]*model.VipInfo{
			110007391: {
				Type:    2,
				Status:  1,
				DueDate: 162930240,
			},
			110010271: {
				Type:    2,
				Status:  1,
				DueDate: 162930240,
			},
		}
	)
	convey.Convey("Batch set vip-infos to cache", t, func(ctx convey.C) {
		err := d.AddCacheVips(c, vm)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCacheProfile(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110011831)
	)
	convey.Convey("Get profile-info from cache", t, func(ctx convey.C) {
		_, err := d.CacheProfile(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddCacheProfile(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110011951)
		v   = &model.Profile{
			Mid:            10920044,
			Name:           "冠冠爱看书",
			Sex:            "男",
			Face:           "http://i1.hdslb.com/bfs/face/4b12a3e65d344e31a11e6425767863019738c7bc.jpg",
			Sign:           "来点字",
			Rank:           10000,
			Level:          5,
			JoinTime:       1503296503,
			Moral:          71,
			Silence:        0,
			EmailStatus:    1,
			TelStatus:      1,
			Identification: 0,
			Vip: model.VipInfo{
				Type:    2,
				Status:  1,
				DueDate: 1629302400000,
			},
			Pendant: model.PendantInfo{
				Pid:    159,
				Name:   "2018拜年祭",
				Image:  "http://i2.hdslb.com/bfs/face/aace621fa64a698f2ca94d13645a26e9a7a99ed2.png",
				Expire: 1551413548,
			},
			Nameplate: model.NameplateInfo{
				Nid:        7,
				Name:       "见习搬运工",
				Image:      "http://i1.hdslb.com/bfs/face/8478fb7c54026cd47f09daa493a1b1683113a90d.png",
				ImageSmall: "http://i0.hdslb.com/bfs/face/50eef47c3a30a75659d3cc298cfb09031d1a2ce5.png",
				Level:      "普通勋章",
				Condition:  "转载视频投稿通过总数>=10",
			},
		}
	)
	convey.Convey("Set profile-info to cache", t, func(ctx convey.C) {
		err := d.AddCacheProfile(c, mid, v)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(110014081)
	)
	convey.Convey("Delete member's cache", t, func(ctx convey.C) {
		errs := d.DelCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			for _, e := range errs {
				if e != mc.ErrNotFound {
					ctx.So(e, convey.ShouldBeNil)
				}
			}
		})
	})
}
