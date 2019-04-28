package like

import (
	"context"

	match "go-common/app/interface/main/activity/model/like"

	xtime "go-common/library/time"
	"testing"

	"time"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikekeyMatch(t *testing.T) {
	convey.Convey("keyMatch", t, func(ctx convey.C) {
		var (
			id = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyMatch(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyActMatch(t *testing.T) {
	convey.Convey("keyActMatch", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyActMatch(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyObject(t *testing.T) {
	convey.Convey("keyObject", t, func(ctx convey.C) {
		var (
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyObject(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyObjects(t *testing.T) {
	convey.Convey("keyObjects", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyObjects(sid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyUserLog(t *testing.T) {
	convey.Convey("keyUserLog", t, func(ctx convey.C) {
		var (
			sid = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyUserLog(sid, mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikekeyMatchFollow(t *testing.T) {
	convey.Convey("keyMatchFollow", t, func(ctx convey.C) {
		var (
			mid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyMatchFollow(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeMatchCache(t *testing.T) {
	convey.Convey("MatchCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mat, err := d.MatchCache(c, id)
			ctx.Convey("Then err should be nil.mat should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", mat)
			})
		})
	})
}

func TestLikeSetMatchCache(t *testing.T) {
	convey.Convey("SetMatchCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			id  = int64(10256)
			mat = &match.Match{Sid: 10256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetMatchCache(c, id, mat)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeActMatchCache(t *testing.T) {
	convey.Convey("ActMatchCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ActMatchCache(c, sid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeSetActMatchCache(t *testing.T) {
	convey.Convey("SetActMatchCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			sid    = int64(10256)
			matchs = []*match.Match{{Sid: 10256}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetActMatchCache(c, sid, matchs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeObjectCache(t *testing.T) {
	convey.Convey("ObjectCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mat, err := d.ObjectCache(c, id)
			ctx.Convey("Then err should be nil.mat should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", mat)
			})
		})
	})
}

func TestLikeCacheMatchSubjects(t *testing.T) {
	convey.Convey("CacheMatchSubjects", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{10256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.CacheMatchSubjects(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSetObjectCache(t *testing.T) {
	convey.Convey("SetObjectCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			id     = int64(10256)
			object = &match.Object{Sid: 10256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetObjectCache(c, id, object)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeAddCacheMatchSubjects(t *testing.T) {
	convey.Convey("AddCacheMatchSubjects", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = map[int64]*match.Object{1: {Sid: 1}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheMatchSubjects(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeObjectsCache(t *testing.T) {
	convey.Convey("ObjectsCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			start = int(1)
			end   = int(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, total, err := d.ObjectsCache(c, sid, start, end)
			ctx.Convey("Then err should be nil.res,total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v,%+v", res, total)
			})
		})
	})
}

func TestLikeSetObjectsCache(t *testing.T) {
	convey.Convey("SetObjectsCache", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			sid     = int64(10256)
			objects = []*match.Object{{Sid: 10256}}
			total   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetObjectsCache(c, sid, objects, total)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeSetUserLogCache(t *testing.T) {
	convey.Convey("SetUserLogCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			sid      = int64(10256)
			mid      = int64(7)
			userLogs = []*match.UserLog{{Sid: 10256}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUserLogCache(c, sid, mid, userLogs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeDelUserLogCache(t *testing.T) {
	convey.Convey("DelUserLogCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			sid = int64(10256)
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelUserLogCache(c, sid, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeDelActMatchCache(t *testing.T) {
	convey.Convey("DelActMatchCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(10256)
			matID = int64(7)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelActMatchCache(c, sid, matID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeDelObjectCache(t *testing.T) {
	convey.Convey("DelObjectCache", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			objID = int64(1)
			sid   = int64(10256)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelObjectCache(c, objID, sid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeAddFollow(t *testing.T) {
	convey.Convey("AddFollow", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(77)
			teams = []string{"qwe"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddFollow(c, mid, teams)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestLikeFollow(t *testing.T) {
	convey.Convey("Follow", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Follow(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikefrom(t *testing.T) {
	convey.Convey("from", t, func(ctx convey.C) {
		var (
			i = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := from(i)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikecombine(t *testing.T) {
	convey.Convey("combine", t, func(ctx convey.C) {
		var (
			ctime = xtime.Time(time.Now().Unix())
			count = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := combine(ctime, count)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
