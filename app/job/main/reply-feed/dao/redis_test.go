package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/reply-feed/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyUV(t *testing.T) {
	convey.Convey("keyUV", t, func(ctx convey.C) {
		var (
			action = ""
			hour   = 0
			slot   = 0
			kind   = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyUV(action, hour, slot, kind)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUV(t *testing.T) {
	convey.Convey("addUV+countUV", t, func(ctx convey.C) {
		var (
			action = "test"
			hour   = 0
			slot   = 0
			kind   = "test"
			mid    = int64(0)
			err    error
			counts []int64
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err = d.AddUV(context.Background(), action, hour, slot, mid, kind)
			ctx.So(err, convey.ShouldBeNil)
			keys := []string{keyUV(action, hour, slot, kind)}
			counts, err = d.CountUV(context.Background(), keys)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(counts), convey.ShouldEqual, 1)
		})
	})
}

func TestDaokeyRefreshChecker(t *testing.T) {
	convey.Convey("keyRefreshChecker", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRefreshChecker(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyReplyZSet(t *testing.T) {
	convey.Convey("keyReplyZSet", t, func(ctx convey.C) {
		var (
			name = ""
			oid  = int64(0)
			tp   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyReplyZSet(name, oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokeyReplySet(t *testing.T) {
	convey.Convey("keyReplySet", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyReplySet(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPingRedis(t *testing.T) {
	convey.Convey("PingRedis", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PingRedis(context.Background())
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExpireCheckerRds(t *testing.T) {
	convey.Convey("ExpireCheckerRds", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireCheckerRds(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireReplyZSetRds(t *testing.T) {
	convey.Convey("ExpireReplyZSetRds", t, func(ctx convey.C) {
		var (
			name = ""
			oid  = int64(0)
			tp   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireReplyZSetRds(context.Background(), name, oid, tp)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoExpireReplySetRds(t *testing.T) {
	convey.Convey("ExpireReplySetRds", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ok, err := d.ExpireReplySetRds(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetReplySetRds(t *testing.T) {
	convey.Convey("ReplySetRds", t, func(ctx convey.C) {
		var (
			oid   = int64(-1)
			tp    = int(-1)
			idMap = make(map[int64]struct{})
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			for i := 0; i < 10; i++ {
				if err := d.AddReplySetRds(context.Background(), oid, tp, int64(i)); err != nil {
					return
				}
				idMap[int64(i)] = struct{}{}
			}
			d.ExpireReplySetRds(context.Background(), oid, tp)
			rpIDs, err := d.ReplySetRds(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(rpIDs), convey.ShouldEqual, 10)
				for _, rpID := range rpIDs {
					if _, ok := idMap[rpID]; !ok {
						t.Fatal("id not match")
					}
				}
			})
		})
	})
}

func TestDaoSetReplySetRds(t *testing.T) {
	convey.Convey("SetReplySetRds", t, func(ctx convey.C) {
		var (
			oid   = int64(0)
			tp    = int(0)
			rpIDs = []int64{}
			c     = context.Background()
		)
		for i := 0; i < 10000; i++ {
			rpIDs = append(rpIDs, int64(i))
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetReplySetRds(c, oid, tp, rpIDs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			receivedRpIDs, err := d.ReplySetRds(c, oid, tp)
			if err != nil {
				t.Fatal(err)
			}
			ctx.So(len(receivedRpIDs), convey.ShouldEqual, len(rpIDs))
		})
		d.DelReplySetRds(c, oid, tp)
	})
}

func TestDaoRemReplySetRds(t *testing.T) {
	convey.Convey("RemReplySetRds", t, func(ctx convey.C) {
		var (
			oid  = int64(0)
			rpID = int64(0)
			tp   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemReplySetRds(context.Background(), oid, rpID, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelReplySetRds(t *testing.T) {
	convey.Convey("DelReplySetRds", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelReplySetRds(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddReplySetRds(t *testing.T) {
	convey.Convey("AddReplySetRds", t, func(ctx convey.C) {
		var (
			oid  = int64(0)
			tp   = int(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddReplySetRds(context.Background(), oid, tp, rpID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetReplyZSetRds(t *testing.T) {
	convey.Convey("ReplyZSetRds", t, func(ctx convey.C) {
		var (
			name  = ""
			oid   = int64(0)
			tp    = int(0)
			start = int(0)
			end   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rpIDs, err := d.ReplyZSetRds(context.Background(), name, oid, tp, start, end)
			ctx.Convey("Then err should be nil.rpIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rpIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetReplyZSetRds(t *testing.T) {
	convey.Convey("SetReplyZSetRds", t, func(ctx convey.C) {
		var (
			name = "test"
			oid  = int64(0)
			tp   = int(0)
			c    = context.Background()
			rs   = []*model.ReplyScore{}
		)
		for i := 0; i < 10000; i++ {
			rs = append(rs, &model.ReplyScore{RpID: int64(i), Score: float64(i)})
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetReplyZSetRds(c, name, oid, tp, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			rpIDs, err := d.ReplyZSetRds(c, name, oid, tp, 0, -1)
			if err != nil {
				t.Fatal(err)
			}
			ctx.So(len(rpIDs), convey.ShouldEqual, len(rs))
		})
		d.DelReplyZSetRds(c, []string{name}, oid, tp)
	})
}

func TestDaoAddReplyZSetRds(t *testing.T) {
	convey.Convey("AddReplyZSetRds", t, func(ctx convey.C) {
		var (
			name = ""
			oid  = int64(0)
			tp   = int(0)
			rs   = &model.ReplyScore{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddReplyZSetRds(context.Background(), name, oid, tp, rs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRemReplyZSetRds(t *testing.T) {
	convey.Convey("RemReplyZSetRds", t, func(ctx convey.C) {
		var (
			name = ""
			oid  = int64(0)
			tp   = int(0)
			rpID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.RemReplyZSetRds(context.Background(), name, oid, tp, rpID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelReplyZSetRds(t *testing.T) {
	convey.Convey("DelReplyZSetRds", t, func(ctx convey.C) {
		var (
			names = []string{}
			oid   = int64(0)
			tp    = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelReplyZSetRds(context.Background(), names, oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetCheckerTsRds(t *testing.T) {
	convey.Convey("CheckerTsRds", t, func(ctx convey.C) {
		var (
			oid = int64(-1)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ts, err := d.CheckerTsRds(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.ts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ts, convey.ShouldEqual, 0)
			})
		})
	})
}

func TestDaoSetCheckerTsRds(t *testing.T) {
	convey.Convey("SetCheckerTsRds", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetCheckerTsRds(context.Background(), oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoRangeReplyZSetRds(t *testing.T) {
	convey.Convey("RangeReplyZSetRds", t, func(ctx convey.C) {
		var (
			oid  = int64(0)
			tp   = int(0)
			name = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rpIDs, err := d.RangeReplyZSetRds(context.Background(), name, oid, tp, 0, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(rpIDs), convey.ShouldEqual, 0)
			})
		})
	})
}
