package order

import (
	"context"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func TestOrderExecuteOrders(t *testing.T) {
	convey.Convey("ExecuteOrders", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.executeOrdersURI).Reply(200).JSON(`{"code": 0,"message": "0", "ttl": 1,"data":"{}"}`)
			_, err := d.ExecuteOrders(c, mid, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestOrderUps(t *testing.T) {
	convey.Convey("Ups", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.upsURI).Reply(200).JSON(`
			{
				"code": 0,
				"message": "ok",
				"requestId": "46aab0b0157111e9bca70a36f8bfcc85",
				"ts": 1547191187259,
				"data": [
				  2089809
				]
			  }`)
			_, err := d.Ups(c)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestOrderOrderByAid(t *testing.T) {
	convey.Convey("OrderByAid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(10111835)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.getOrderByAidURI).Reply(200).JSON(`
			{
				"code": 0,
				"message": "ok",
				"requestId": "cbf35e10157211e9bca70a36f8bfcc85",
				"ts": 1547191840369,
				"data": {
				  "execute_order_id": 13,
				  "business_order_id": 6,
				  "business_order_name": "这是一个测试项目",
				  "id_code": 10002824,
				  "game_base_id": 0,
				  "game_name": ""
				}
			  }`)
			_, _, _, err := d.OrderByAid(c, aid)
			ctx.Convey("Then err should be nil.orderID,orderName,gameBaseID should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestOrderUnbind(t *testing.T) {
	convey.Convey("Unbind", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			aid = int64(10111835)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.archiveStatusURI).Reply(200).JSON(`{"code": 0,"message": "0", "ttl": 1,"data":"{}"}`)
			err := d.Unbind(c, mid, aid, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestOrderOasis(t *testing.T) {
	convey.Convey("Oasis", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.oasisURI).Reply(200).JSON(`
			{
				"code": 0,
				"message": "ok",
				"requestId": "14909c20157111e9bca70a36f8bfcc85",
				"ts": 1547191103202,
				"data": {
					"state": 1,
					"running_execute_order_count": 0,
					"total_execute_order_count": 0
				}
			}`)
			_, err := d.Oasis(c, mid, ip)
			ctx.Convey("Then err should be nil.oa should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestOrderLaunchTime(t *testing.T) {
	convey.Convey("LaunchTime", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			orderID = int64(15)
			ip      = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.launchTimeURI).Reply(200).JSON(`
			{
				"code": 0,
				"message": "ok",
				"requestId": "75839970157011e9bca70a36f8bfcc85",
				"ts": 1547190836359,
				"data": {
					"begin_date": 1547127000
				}
			}
			`)
			_, err := d.LaunchTime(c, orderID, ip)
			ctx.Convey("Then err should be nil.beginDate,endDate should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestOrderUpValidate(t *testing.T) {
	convey.Convey("UpValidate", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			ip  = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.upValidateURI).Reply(200).JSON(`{"code": 0,"message": "0", "ttl": 1,"data":"{}"}`)
			_, err := d.UpValidate(c, mid, ip)
			ctx.Convey("Then err should be nil.uv should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestOrderGrowAccountState(t *testing.T) {
	convey.Convey("GrowAccountState", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(2089809)
			ty  = int(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.accountStateURI).Reply(200).JSON(`{"code": 0,"message": "0", "ttl": 1,"data":"{}"}`)
			_, err := d.GrowAccountState(c, mid, ty)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, err)
			})
		})
	})
}
