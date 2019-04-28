package datadao

import (
	"context"
	"go-common/app/interface/main/mcn/tool/datacenter"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/mcn/model/datamodel"
)

func TestDatadaocallDataAPI(t *testing.T) {
	convey.Convey("callDataAPI", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			api   = APIMcnSummary
			query = &datacenter.Query{}
			res   []*datamodel.DmConMcnArchiveD
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.callDataAPI(c, api, query, res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnSummary(t *testing.T) {
	convey.Convey("GetMcnSummary", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnSummary(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnSummaryCache(t *testing.T) {
	convey.Convey("GetMcnSummaryCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnSummaryCache(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetIndexInc(t *testing.T) {
	convey.Convey("GetIndexInc", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetIndexInc(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetIndexIncCache(t *testing.T) {
	convey.Convey("GetIndexIncCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetIndexIncCache(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetIndexSource(t *testing.T) {
	convey.Convey("GetIndexSource", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetIndexSource(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetIndexSourceCache(t *testing.T) {
	convey.Convey("GetIndexSourceCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetIndexSourceCache(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetPlaySource(t *testing.T) {
	convey.Convey("GetPlaySource", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetPlaySource(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetPlaySourceCache(t *testing.T) {
	convey.Convey("GetPlaySourceCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetPlaySourceCache(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFans(t *testing.T) {
	convey.Convey("GetMcnFans", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFans(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFansCache(t *testing.T) {
	convey.Convey("GetMcnFansCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFansCache(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFansInc(t *testing.T) {
	convey.Convey("GetMcnFansInc", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFansInc(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFansIncCache(t *testing.T) {
	convey.Convey("GetMcnFansIncCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFansIncCache(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFansDec(t *testing.T) {
	convey.Convey("GetMcnFansDec", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFansDec(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFansDecCache(t *testing.T) {
	convey.Convey("GetMcnFansDecCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFansDecCache(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFansAttentionWay(t *testing.T) {
	convey.Convey("GetMcnFansAttentionWay", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFansAttentionWay(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetMcnFansAttentionWayCache(t *testing.T) {
	convey.Convey("GetMcnFansAttentionWayCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetMcnFansAttentionWayCache(c, signID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansBaseFansAttr(t *testing.T) {
	convey.Convey("GetFansBaseFansAttr", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansBaseFansAttr(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansBaseFansAttrCache(t *testing.T) {
	convey.Convey("GetFansBaseFansAttrCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansBaseFansAttrCache(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansArea(t *testing.T) {
	convey.Convey("GetFansArea", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansArea(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansAreaCache(t *testing.T) {
	convey.Convey("GetFansAreaCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansAreaCache(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansType(t *testing.T) {
	convey.Convey("GetFansType", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansType(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansTypeCache(t *testing.T) {
	convey.Convey("GetFansTypeCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansTypeCache(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansTag(t *testing.T) {
	convey.Convey("GetFansTag", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansTag(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDatadaoGetFansTagCache(t *testing.T) {
	convey.Convey("GetFansTagCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			signID = int64(0)
			date   = time.Now()
			tp     = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.GetFansTagCache(c, signID, date, tp)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
