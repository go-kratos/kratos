package http

import (
	"sync"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

var (
	mapLock sync.Mutex
	// Lock withdraw lock
	Lock map[int64]*sync.Mutex
)

func getWithdraw(c *bm.Context) {
	v := new(struct {
		DateVersion string `form:"dateVersion" validate:"required"`
		Page        int    `form:"pageNo" default:"1" validate:"min=1"`
		Size        int    `form:"size" default:"15" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	from := (v.Page - 1) * v.Size
	count, data, err := svc.GetWithdraw(c, v.DateVersion, from, v.Size)
	if err != nil {
		log.Error("growup svc.upWithdraw error(%v)", err)
		c.JSON(nil, err)
		return
	}

	Lock = make(map[int64]*sync.Mutex)
	c.JSON(map[string]interface{}{
		"page":        v.Page,
		"total_count": count,
		"data":        data,
	}, nil)
}

func withdrawSuccess(c *bm.Context) {
	v := new(struct {
		OrderNo     int64 `form:"order_no" validate:"required"`
		TradeStatus int   `form:"trade_status" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	lock := getLock(v.OrderNo)
	lock.Lock()
	defer lock.Unlock()
	err := svc.WithdrawSuccess(c, v.OrderNo, v.TradeStatus)
	if err != nil {
		log.Error("growup svc.WithdrawSuccess error(%v)", err)
	}
	c.JSON(nil, err)
}

func getLock(orderNo int64) *sync.Mutex {
	mapLock.Lock()
	defer mapLock.Unlock()

	if Lock == nil {
		Lock = make(map[int64]*sync.Mutex)
	}

	if _, ok := Lock[orderNo]; !ok {
		Lock[orderNo] = new(sync.Mutex)
	}
	return Lock[orderNo]
}

func withdrawDetail(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.WithdrawDetail(c, mid)
	if err != nil {
		log.Error("growup svc.WithdrawDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
