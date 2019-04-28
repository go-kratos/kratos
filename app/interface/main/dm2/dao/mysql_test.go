package dao

import (
	"context"
	xtime "go-common/library/time"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSubject(t *testing.T) {
	convey.Convey("Subject", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.Subject(c, tp, oid)
		})
	})
}

func TestDaoSubjects(t *testing.T) {
	convey.Convey("Subjects", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.Subjects(c, tp, oids)
		})
	})
}

func TestDaoUptSubAttr(t *testing.T) {
	convey.Convey("UptSubAttr", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oid  = int64(0)
			attr = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.UptSubAttr(c, tp, oid, attr)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIncrSubMoveCount(t *testing.T) {
	convey.Convey("IncrSubMoveCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			count = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.IncrSubMoveCount(c, tp, oid, count)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpSubjectMCount(t *testing.T) {
	convey.Convey("UpSubjectMCount", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
			cnt = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.UpSubjectMCount(c, tp, oid, cnt)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpSubjectPool(t *testing.T) {
	convey.Convey("UpSubjectPool", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tp        = int32(0)
			oid       = int64(0)
			childpool = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.UpSubjectPool(c, tp, oid, childpool)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIncrSubjectCount(t *testing.T) {
	convey.Convey("IncrSubjectCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			count = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.IncrSubjectCount(c, tp, oid, count)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDMIDs(t *testing.T) {
	convey.Convey("DMIDs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			ps    = int64(0)
			pe    = int64(0)
			limit = int64(0)
			pool  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.DMIDs(c, tp, oid, ps, pe, limit, pool)
		})
	})
}

func TestDaoIndexs(t *testing.T) {
	convey.Convey("Indexs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.Indexs(c, tp, oid)
		})
	})
}

func TestDaoIndexByid(t *testing.T) {
	convey.Convey("IndexByid", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int8(0)
			oid  = int64(0)
			dmid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.IndexByid(c, tp, oid, dmid)
		})
	})
}

func TestDaoIndexsByid(t *testing.T) {
	convey.Convey("IndexsByid", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.IndexsByid(c, tp, oid, dmids)
		})
	})
}

func TestDaoJudgeIndex(t *testing.T) {
	convey.Convey("JudgeIndex", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tp     = int8(0)
			oid    = int64(0)
			ctime1 = xtime.Time(time.Now().Unix())
			ctime2 = xtime.Time(time.Now().Unix())
			prog1  = int32(0)
			prog2  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.JudgeIndex(c, tp, oid, ctime1, ctime2, prog1, prog2)
		})
	})
}

func TestDaoContent(t *testing.T) {
	convey.Convey("Content", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			oid  = int64(0)
			dmid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ct, err := testDao.Content(c, oid, dmid)
			ctx.Convey("Then err should be nil.ct should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ct, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoContents(t *testing.T) {
	convey.Convey("Contents", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			oid   = int64(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ctsMap, err := testDao.Contents(c, oid, dmids)
			ctx.Convey("Then err should be nil.ctsMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ctsMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoContentsSpecial(t *testing.T) {
	convey.Convey("ContentsSpecial", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			dmids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.ContentsSpecial(c, dmids)
		})
	})
}

func TestDaoUpdateDMStat(t *testing.T) {
	convey.Convey("UpdateDMStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			state = int32(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpdateDMStat(c, tp, oid, state, dmids)
		})
	})
}

func TestDaoUpdateUserDMStat(t *testing.T) {
	convey.Convey("UpdateUserDMStat", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			mid   = int64(0)
			state = int32(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpdateUserDMStat(c, tp, oid, mid, state, dmids)
		})
	})
}

func TestDaoUpdateDMPool(t *testing.T) {
	convey.Convey("UpdateDMPool", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tp    = int32(0)
			oid   = int64(0)
			pool  = int32(0)
			dmids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.UpdateDMPool(c, tp, oid, pool, dmids)
		})
	})
}

func TestDaoUpdateDMAttr(t *testing.T) {
	convey.Convey("UpdateDMAttr", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tp   = int32(0)
			oid  = int64(0)
			dmid = int64(0)
			attr = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.UpdateDMAttr(c, tp, oid, dmid, attr)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDMCount(t *testing.T) {
	convey.Convey("DMCount", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			typ    = int32(0)
			oid    = int64(0)
			states = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.DMCount(c, typ, oid, states)
		})
	})
}

func TestDaoSpecialDmLocation(t *testing.T) {
	convey.Convey("SpecialDmLocation", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.SpecialDmLocation(c, tp, oid)
		})
	})
}

func TestDaoSpecalDMs(t *testing.T) {
	convey.Convey("SpecalDMs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tp  = int32(0)
			oid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			testDao.SpecalDMs(c, tp, oid)
		})
	})
}

func TestDaoAddUpperConfig(t *testing.T) {
	convey.Convey("AddUpperConfig", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(0)
			advPermit = int8(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			affect, err := testDao.AddUpperConfig(c, mid, advPermit)
			ctx.Convey("Then err should be nil.affect should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affect, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpperConfig(t *testing.T) {
	convey.Convey("UpperConfig", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			advPermit, err := testDao.UpperConfig(c, mid)
			ctx.Convey("Then err should be nil.advPermit should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(advPermit, convey.ShouldNotBeNil)
			})
		})
	})
}
