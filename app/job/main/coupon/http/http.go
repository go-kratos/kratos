package http

import (
	"net/http"

	"go-common/app/job/main/coupon/conf"
	"go-common/app/job/main/coupon/model"
	"go-common/app/job/main/coupon/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	svc *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	svc = s
	// init router
	engineInner := bm.DefaultServer(c.BM)
	innerRouter(engineInner)
	if err := engineInner.Start(); err != nil {
		log.Error("bm.DefaultServer error(%v)", err)
		panic(err)
	}
}

// innerRouter init inner router api path.
func innerRouter(e *bm.Engine) {
	//init api
	e.GET("/monitor/ping", ping)
	e.GET("/deliver/success", cartoonDeliverSuccess)
}

// ping check server ok.
func ping(c *bm.Context) {
	if err := svc.Ping(c); err != nil {
		log.Error("coupon http service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func cartoonDeliverSuccess(c *bm.Context) {
	var (
		err error
		o   *model.CouponOrder
	)
	arg := new(struct {
		OrderNo string `form:"order_no"`
		State   int8   `form:"state"`
		Ver     int64  `form:"ver"`
		IsPaid  int8   `form:"is_paid"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("add coupon page %+v", err)
		return
	}
	if o, err = svc.ByOrderNo(c, arg.OrderNo); err != nil {
		c.JSON(nil, err)
		return
	}
	if o == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.UpdateOrderState(c, o, arg.State, &model.CallBackRet{Ver: arg.Ver, IsPaid: arg.IsPaid}); err != nil {
		log.Error("svc.UpdateOrderState(%s) %+v", arg.OrderNo, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
