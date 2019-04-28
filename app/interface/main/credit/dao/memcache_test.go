package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaonoticeKey(t *testing.T) {
	convey.Convey("noticeKey", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := noticeKey()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoreasonKey(t *testing.T) {
	convey.Convey("reasonKey", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := reasonKey()
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoopinionKey(t *testing.T) {
	convey.Convey("opinionKey", t, func(convCtx convey.C) {
		var (
			opid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := opinionKey(opid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoquestionKey(t *testing.T) {
	convey.Convey("questionKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := questionKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaolabourKey(t *testing.T) {
	convey.Convey("labourKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := labourKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaojuryInfoKey(t *testing.T) {
	convey.Convey("juryInfoKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := juryInfoKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.pingMC(c)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoNoticeInfoCache(t *testing.T) {
	convey.Convey("NoticeInfoCache", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			dt, err := d.NoticeInfoCache(c)
			convCtx.Convey("Then err should be nil.dt should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(dt, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetReasonListCache(t *testing.T) {
	convey.Convey("SetReasonListCache", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			dt = []*model.Reason{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetReasonListCache(c, dt)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetNoticeInfoCache(t *testing.T) {
	convey.Convey("SetNoticeInfoCache", t, func(convCtx convey.C) {
		var (
			c = context.Background()
			n = &model.Notice{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetNoticeInfoCache(c, n)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelNoticeInfoCache(t *testing.T) {
	convey.Convey("DelNoticeInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelNoticeInfoCache(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoReasonListCache(t *testing.T) {
	convey.Convey("ReasonListCache", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			dt, err := d.ReasonListCache(c)
			convCtx.Convey("Then err should be nil.dt should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(dt, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelReasonListCache(t *testing.T) {
	convey.Convey("DelReasonListCache", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelReasonListCache(c)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddOpinionCache(t *testing.T) {
	convey.Convey("AddOpinionCache", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			op = &model.Opinion{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddOpinionCache(c, op)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOpinionsCache(t *testing.T) {
	convey.Convey("OpinionsCache", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			opids = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ops, miss, err := d.OpinionsCache(c, opids)
			convCtx.Convey("Then err should be nil.ops should not be nil,miss should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(miss, convey.ShouldBeNil)
				convCtx.So(ops, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetQsCache(t *testing.T) {
	convey.Convey("SetQsCache", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(0)
			qsid = &model.QsCache{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetQsCache(c, mid, qsid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetQsCache(t *testing.T) {
	convey.Convey("GetQsCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			qsid, err := d.GetQsCache(c, mid)
			convCtx.Convey("Then err should be nil.qsid should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(qsid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelQsCache(t *testing.T) {
	convey.Convey("DelQsCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelQsCache(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetAnswerStateCache(t *testing.T) {
	convey.Convey("SetAnswerStateCache", t, func(convCtx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(0)
			state = int8(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetAnswerStateCache(c, mid, state)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetAnswerStateCache(t *testing.T) {
	convey.Convey("GetAnswerStateCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			state, found, err := d.GetAnswerStateCache(c, mid)
			convCtx.Convey("Then err should be nil.state,found should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(found, convey.ShouldNotBeNil)
				convCtx.So(state, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoJuryInfoCache(t *testing.T) {
	convey.Convey("JuryInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			bj, err := d.JuryInfoCache(c, mid)
			convCtx.Convey("Then err should be nil.bj should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(bj, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetJuryInfoCache(t *testing.T) {
	convey.Convey("SetJuryInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			bj  = &model.BlockedJury{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetJuryInfoCache(c, mid, bj)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelJuryInfoCache(t *testing.T) {
	convey.Convey("DelJuryInfoCache", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DelJuryInfoCache(c, mid)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
