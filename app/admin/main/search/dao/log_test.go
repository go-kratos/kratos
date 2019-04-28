package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/search/model"
)

func TestDaoNewLog(t *testing.T) {
	convey.Convey("NewLog", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.NewLog()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDaoGetLogInfo(t *testing.T) {
	convey.Convey("GetLogInfo", t, func(ctx convey.C) {
		var (
			appID = ""
			id    = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			business, ok := d.GetLogInfo(appID, id)
			ctx.Convey("Then business,ok should not be nil.", func(ctx convey.C) {
				ctx.So(ok, convey.ShouldNotBeNil)
				ctx.So(business, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoinitMapping(t *testing.T) {
	convey.Convey("initMapping", t, func(ctx convey.C) {
		var (
			appID = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			business, err := d.initMapping(appID)
			ctx.Convey("Then business should not be nil.", func(ctx convey.C) {
				ctx.So(business, convey.ShouldNotBeNil)
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaologIndexName(t *testing.T) {
	convey.Convey("logIndexName", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &model.LogParams{
				CTimeFrom: "2010-01-01 00:00:00",
				CTimeTo:   "2020-01-01 00:00:00",
			}
			business = &model.Business{
				ID:    0,
				AppID: "log_audit",
			}
		)
		ctx.Convey("2006", func(ctx convey.C) {
			business.IndexFormat = "2006"
			res, err := d.logIndexName(c, p, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("2006-01", func(ctx convey.C) {
			business.IndexFormat = "2006-01"
			res, err := d.logIndexName(c, p, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("2006-01-week", func(ctx convey.C) {
			business.IndexFormat = "2006-01-week"
			res, err := d.logIndexName(c, p, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("2006-01-02", func(ctx convey.C) {
			business.IndexFormat = "2006-01-02"
			res, err := d.logIndexName(c, p, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("all", func(ctx convey.C) {
			business.IndexFormat = "all"
			res, err := d.logIndexName(c, p, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetLogAuditIndexName(t *testing.T) {
	convey.Convey("getLogAuditIndexName", t, func(ctx convey.C) {
		var (
			business  = int(0)
			indexName = ""
			format    = ""
			time      = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			index := getLogAuditIndexName(business, indexName, format, time)
			ctx.Convey("Then index should not be nil.", func(ctx convey.C) {
				ctx.So(index, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaogetQuery(t *testing.T) {
	convey.Convey("getQuery", t, func(ctx convey.C) {
		var (
			pr           map[string][]interface{}
			indexMapping map[string]string
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			query := d.getQuery(pr, indexMapping)
			ctx.Convey("Then query should not be nil.", func(ctx convey.C) {
				ctx.So(query, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogAudit(t *testing.T) {
	convey.Convey("LogAudit", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			pr       map[string][]interface{}
			sp       = &model.LogParams{}
			business = &model.Business{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LogAudit(c, pr, sp, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogAuditGroupBy(t *testing.T) {
	convey.Convey("LogAuditGroupBy", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			pr = map[string][]interface{}{
				"group": {"group"},
			}
			sp       = &model.LogParams{}
			business = &model.Business{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LogAuditGroupBy(c, pr, sp, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogAuditDelete(t *testing.T) {
	convey.Convey("LogAuditDelete", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			pr       map[string][]interface{}
			sp       = &model.LogParams{}
			business = &model.Business{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LogAuditDelete(c, pr, sp, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogUserAction(t *testing.T) {
	convey.Convey("LogUserAction", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			pr       map[string][]interface{}
			sp       = &model.LogParams{}
			business = &model.Business{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LogUserAction(c, pr, sp, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogUserActionDelete(t *testing.T) {
	convey.Convey("LogUserActionDelete", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			pr       map[string][]interface{}
			sp       = &model.LogParams{}
			business = &model.Business{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.LogUserActionDelete(c, pr, sp, business)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUDepTs(t *testing.T) {
	convey.Convey("UDepTs", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			uids = []string{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.UDepTs(c, uids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIP(t *testing.T) {
	convey.Convey("IP", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ip := []string{
			"127.0.0.1",
		}
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//res, err :=
			d.IP(c, ip)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLogCount(t *testing.T) {
	convey.Convey("LogCount", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			name     = ""
			business = int(0)
			uid      = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.LogCount(c, name, business, uid)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
