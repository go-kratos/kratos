package like

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-common/app/interface/main/activity/model/like"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeipRequestKey(t *testing.T) {
	convey.Convey("ipRequestKey", t, func(ctx convey.C) {
		var (
			ip = "10.256.8.3"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := ipRequestKey(ip)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeListCtimeKey(t *testing.T) {
	convey.Convey("likeListCtimeKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeListCtimeKey(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeListRandomKey(t *testing.T) {
	convey.Convey("likeListCtimeKey", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := likeListRandomKey(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikelikeListTypeCtimeKey(t *testing.T) {
	convey.Convey("likeListTypeCtimeKey", t, func(ctx convey.C) {
		var (
			types = int(1)
			sid   = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := likeListTypeCtimeKey(types, sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyLikeTag(t *testing.T) {
	convey.Convey("keyLikeTag", t, func(ctx convey.C) {
		var (
			sid   = int64(10256)
			tagID = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyLikeTag(sid, tagID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyLikeTagCounts(t *testing.T) {
	convey.Convey("keyLikeTagCounts", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyLikeTagCounts(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyLikeRegion(t *testing.T) {
	convey.Convey("keyLikeRegion", t, func(ctx convey.C) {
		var (
			sid      = int64(10256)
			regionID = int16(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyLikeRegion(sid, regionID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyStoryLikeKey(t *testing.T) {
	convey.Convey("keyStoryLikeKey", t, func(ctx convey.C) {
		var (
			sid   = int64(10256)
			mid   = int64(1)
			daily = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyStoryLikeKey(sid, mid, daily)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyStoryEachLike(t *testing.T) {
	convey.Convey("keyStoryEachLike", t, func(ctx convey.C) {
		var (
			sid   = int64(10256)
			mid   = int64(1)
			lid   = int64(1)
			daily = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyStoryEachLike(sid, mid, lid, daily)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeLikeList(t *testing.T) {
	convey.Convey("LikeList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ns, err := d.LikeList(c, sid)
			ctx.Convey("Then err should be nil.ns should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ns, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeLikeTagCache(t *testing.T) {
	convey.Convey("LikeTagCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			tagID = int64(1)
			start = int(1)
			end   = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			likes, err := d.LikeTagCache(c, sid, tagID, start, end)
			ctx.Convey("Then err should be nil.likes should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", likes)
			})
		})
	})
}

func TestLikeLikeTagCnt(t *testing.T) {
	convey.Convey("LikeTagCnt", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			tagID = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.LikeTagCnt(c, sid, tagID)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeLikeRegionCache(t *testing.T) {
	convey.Convey("LikeRegionCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			sid      = int64(10256)
			regionID = int16(1)
			start    = int(1)
			end      = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			likes, err := d.LikeRegionCache(c, sid, regionID, start, end)
			ctx.Convey("Then err should be nil.likes should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", likes)
			})
		})
	})
}

func TestLikeLikeRegionCnt(t *testing.T) {
	convey.Convey("LikeRegionCnt", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			sid      = int64(10256)
			regionID = int16(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.LikeRegionCnt(c, sid, regionID)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSetLikeRegionCache(t *testing.T) {
	convey.Convey("SetLikeRegionCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			sid      = int64(10256)
			regionID = int16(1)
			likes    = []*like.Item{{Sid: 10256, Wid: 1, Mid: 44}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetLikeRegionCache(c, sid, regionID, likes)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeSetTagLikeCountsCache(t *testing.T) {
	convey.Convey("SetTagLikeCountsCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			sid    = int64(10256)
			counts = map[int64]int32{1: 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetTagLikeCountsCache(c, sid, counts)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeTagLikeCountsCache(t *testing.T) {
	convey.Convey("TagLikeCountsCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			sid    = int64(10256)
			tagIDs = []int64{1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			counts, err := d.TagLikeCountsCache(c, sid, tagIDs)
			ctx.Convey("Then err should be nil.counts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(counts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeRawLike(t *testing.T) {
	convey.Convey("RawLike", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawLike(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeLikeListMoreLid(t *testing.T) {
	convey.Convey("LikeListMoreLid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			lid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LikeListMoreLid(c, lid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeLikesBySid(t *testing.T) {
	convey.Convey("LikesBySid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			lid = int64(77)
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LikesBySid(c, lid, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeIPReqquestCheck(t *testing.T) {
	convey.Convey("IPReqquestCheck", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ip = "10.248.56.23"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			val, err := d.IPReqquestCheck(c, ip)
			ctx.Convey("Then err should be nil.val should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(val, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSetIPRequest(t *testing.T) {
	convey.Convey("SetIPRequest", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ip = "10.248.56.23"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetIPRequest(c, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeLikeListCtime(t *testing.T) {
	convey.Convey("LikeListCtime", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			items = []*like.Item{{Sid: 10256, Wid: 55, Mid: 234}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.LikeListCtime(c, sid, items)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeDelLikeListCtime(t *testing.T) {
	convey.Convey("DelLikeListCtime", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			items = []*like.Item{{Sid: 10256, Wid: 55, Mid: 234}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelLikeListCtime(c, sid, items)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeLikeMaxID(t *testing.T) {
	convey.Convey("LikeMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LikeMaxID(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeStoryLikeSum(t *testing.T) {
	convey.Convey("StoryLikeSum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.StoryLikeSum(c, sid, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeIncrStoryLikeSum(t *testing.T) {
	convey.Convey("IncrStoryLikeSum", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			mid   = int64(77)
			score = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.IncrStoryLikeSum(c, sid, mid, score)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSetLikeSum(t *testing.T) {
	convey.Convey("SetLikeSum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			mid = int64(77)
			sum = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetLikeSum(c, sid, mid, sum)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				fmt.Printf("%v", err)
			})
		})
	})
}

func TestLikeStoryEachLikeSum(t *testing.T) {
	convey.Convey("StoryEachLikeSum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10296)
			mid = int64(216761)
			lid = int64(13538)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.StoryEachLikeSum(c, sid, mid, lid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%d", res)
			})
		})
	})
}

func TestLikeIncrStoryEachLikeAct(t *testing.T) {
	convey.Convey("IncrStoryEachLikeAct", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			mid   = int64(77)
			lid   = int64(77)
			score = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.IncrStoryEachLikeAct(c, sid, mid, lid, score)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSetEachLikeSum(t *testing.T) {
	convey.Convey("SetEachLikeSum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			mid = int64(77)
			lid = int64(77)
			sum = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetEachLikeSum(c, sid, mid, lid, sum)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				fmt.Printf("%v", err)
			})
		})
	})
}

func TestLikeCtime(t *testing.T) {
	convey.Convey("SetEachLikeSum", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10365)
			start = 1
			end   = 100
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LikeCtime(c, sid, start, end)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%v", res)
			})
		})
	})
}

func TestLikeRandom(t *testing.T) {
	convey.Convey("SetEachLikeSum", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10365)
			start = 1
			end   = 100
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.LikeRandom(c, sid, start, end)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%v", res)
			})
		})
	})
}

func TestLikeRandomCount(t *testing.T) {
	convey.Convey("SetEachLikeSum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10365)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.LikeRandomCount(c, sid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%v", res)
			})
		})
	})
}

func TestSetLikeRandom(t *testing.T) {
	convey.Convey("SetEachLikeSum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10365)
			ids = []int64{2354, 2355}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetLikeRandom(c, sid, ids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeCount(t *testing.T) {
	convey.Convey("SetEachLikeSum", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10365)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LikeCount(c, sid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%v", res)
			})
		})
	})
}

func TestDao_SourceItemData(t *testing.T) {
	convey.Convey("test group item data", t, func(ctx convey.C) {
		sid := int64(37)
		data, err := d.SourceItemData(context.Background(), sid)
		convey.So(err, convey.ShouldBeNil)
		str, _ := json.Marshal(data)
		convey.Printf("%+v", string(str))
	})
}

func TestDao_ListFromES(t *testing.T) {
	convey.Convey("test group item data", t, func(ctx convey.C) {
		sid := int64(1)
		ps := 100
		pn := 1
		data, err := d.ListFromES(context.Background(), sid, "", ps, pn, time.Now().Unix())
		convey.So(err, convey.ShouldBeNil)
		for _, v := range data.List {
			convey.Printf(" %+v ", v.Item)
		}
	})
}

func TestDao_MultiTags(t *testing.T) {
	convey.Convey("test group item data", t, func(ctx convey.C) {
		wids := []int64{10109984}
		data, err := d.MultiTags(context.Background(), wids)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
	})
}

func TestDao_OidInfoFromES(t *testing.T) {
	convey.Convey("test group item data", t, func(ctx convey.C) {
		oids := []int64{11, 21}
		stype := 1
		data, err := d.OidInfoFromES(context.Background(), oids, stype)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%+v", data)
	})
}
