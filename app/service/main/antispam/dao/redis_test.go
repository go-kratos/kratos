package dao

import (
	"context"
	"go-common/app/service/main/antispam/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaosendersKey(t *testing.T) {
	var (
		keywordID = int64(0)
	)
	convey.Convey("sendersKey", t, func(ctx convey.C) {
		p1 := sendersKey(keywordID)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoareaSendersKey(t *testing.T) {
	var (
		area     = ""
		senderID = int64(0)
	)
	convey.Convey("areaSendersKey", t, func(ctx convey.C) {
		p1 := areaSendersKey(area, senderID)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaototalCountsKey(t *testing.T) {
	var (
		keywordID = int64(0)
	)
	convey.Convey("totalCountsKey", t, func(ctx convey.C) {
		p1 := totalCountsKey(keywordID)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaolocalCountsKey(t *testing.T) {
	var (
		keywordID = int64(0)
		oid       = int64(0)
	)
	convey.Convey("localCountsKey", t, func(ctx convey.C) {
		p1 := localCountsKey(keywordID, oid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoglobalCountsKey(t *testing.T) {
	var (
		keywordID = int64(0)
	)
	convey.Convey("globalCountsKey", t, func(ctx convey.C) {
		p1 := globalCountsKey(keywordID)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaorulesKey(t *testing.T) {
	var (
		area      = ""
		limitType = ""
	)
	convey.Convey("rulesKey", t, func(ctx convey.C) {
		p1 := rulesKey(area, limitType)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaopingRedis(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		err := d.pingRedis(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoCntSendersCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
	)
	convey.Convey("CntSendersCache", t, func(ctx convey.C) {
		cnt, err := d.CntSendersCache(c, keywordID)
		ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(cnt, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoGlobalLocalLimitCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
		oid       = int64(0)
	)
	convey.Convey("GlobalLocalLimitCache", t, func(ctx convey.C) {
		p1, err := d.GlobalLocalLimitCache(c, keywordID, oid)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoIncrGlobalLimitCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
	)
	convey.Convey("IncrGlobalLimitCache", t, func(ctx convey.C) {
		p1, err := d.IncrGlobalLimitCache(c, keywordID)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoIncrLocalLimitCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
		oid       = int64(0)
	)
	convey.Convey("IncrLocalLimitCache", t, func(ctx convey.C) {
		p1, err := d.IncrLocalLimitCache(c, keywordID, oid)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoLocalLimitExpire(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
		oid       = int64(0)
		dur       = int64(0)
	)
	convey.Convey("LocalLimitExpire", t, func(ctx convey.C) {
		err := d.LocalLimitExpire(c, keywordID, oid, dur)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGlobalLimitExpire(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
		dur       = int64(0)
	)
	convey.Convey("GlobalLimitExpire", t, func(ctx convey.C) {
		err := d.GlobalLimitExpire(c, keywordID, dur)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelRegexpCache(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("DelRegexpCache", t, func(ctx convey.C) {
		err := d.DelRegexpCache(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelRulesCache(t *testing.T) {
	var (
		c         = context.TODO()
		area      = ""
		limitType = ""
	)
	convey.Convey("DelRulesCache", t, func(ctx convey.C) {
		err := d.DelRulesCache(c, area, limitType)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAreaSendersExpire(t *testing.T) {
	var (
		c        = context.TODO()
		area     = ""
		senderID = int64(0)
		dur      = int64(0)
	)
	convey.Convey("AreaSendersExpire", t, func(ctx convey.C) {
		err := d.AreaSendersExpire(c, area, senderID, dur)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoIncrAreaSendersCache(t *testing.T) {
	var (
		c        = context.TODO()
		area     = ""
		senderID = int64(0)
	)
	convey.Convey("IncrAreaSendersCache", t, func(ctx convey.C) {
		p1, err := d.IncrAreaSendersCache(c, area, senderID)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAllSendersCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
	)
	convey.Convey("AllSendersCache", t, func(ctx convey.C) {
		p1, err := d.AllSendersCache(c, keywordID)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoSendersCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
		limit     = int64(0)
		offset    = int64(0)
	)
	convey.Convey("SendersCache", t, func(ctx convey.C) {
		p1, err := d.SendersCache(c, keywordID, limit, offset)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoTotalLimitExpire(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
		dur       = int64(0)
	)
	convey.Convey("TotalLimitExpire", t, func(ctx convey.C) {
		err := d.TotalLimitExpire(c, keywordID, dur)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoIncrTotalLimitCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
	)
	convey.Convey("IncrTotalLimitCache", t, func(ctx convey.C) {
		p1, err := d.IncrTotalLimitCache(c, keywordID)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoZaddSendersCache(t *testing.T) {
	var (
		c         = context.TODO()
		keywordID = int64(0)
		score     = int64(0)
		senderID  = int64(0)
	)
	convey.Convey("ZaddSendersCache", t, func(ctx convey.C) {
		p1, err := d.ZaddSendersCache(c, keywordID, score, senderID)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoZremSendersCache(t *testing.T) {
	var (
		c           = context.TODO()
		keywordID   = int64(1)
		senderIDStr = ""
	)
	convey.Convey("ZremSendersCache", t, func(ctx convey.C) {
		p1, err := d.ZremSendersCache(c, keywordID, senderIDStr)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelKeywordRelatedCache(t *testing.T) {
	var (
		c  = context.TODO()
		ks = []*model.Keyword{}
	)
	convey.Convey("DelKeywordRelatedCache", t, func(ctx convey.C) {
		err := d.DelKeywordRelatedCache(c, ks)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDelCountRelatedCache(t *testing.T) {
	var (
		c = context.TODO()
		k = &model.Keyword{}
	)
	convey.Convey("DelCountRelatedCache", t, func(ctx convey.C) {
		err := d.DelCountRelatedCache(c, k)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
