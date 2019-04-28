package http

import (
	"strconv"
	"time"

	pb "go-common/app/service/main/coin/api"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// addCoin
func addCoin(c *bm.Context) {
	var (
		tp   int64
		upid int64
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	tpStr := params.Get("avtype")
	multiplyStr := params.Get("multiply")
	tpidStr := params.Get("typeid")
	maxStr := params.Get("max")
	upidStr := params.Get("upid")
	mid, _ := strconv.ParseInt(params.Get("mid"), 10, 64)
	if mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	multiply, err := strconv.ParseInt(multiplyStr, 10, 64)
	if err != nil || multiply <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if business, ok := c.Get("business"); ok {
		tp = business.(int64)
	} else {
		if tpStr != "" {
			if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil || tp < 1 || tp > 3 {
				c.JSON(nil, ecode.RequestErr)
				return
			}
		} else {
			tp = 1
		}
	}
	if upidStr != "" {
		if upid, err = strconv.ParseInt(upidStr, 10, 64); err != nil || upid <= 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	typeid, _ := strconv.ParseInt(tpidStr, 10, 64)
	max, _ := strconv.ParseInt(maxStr, 10, 8)
	c.JSON(nil, coinSvc.WebAddCoin(c, mid, upid, max, aid, tp, multiply, int16(typeid)))
}

func list(c *bm.Context) {
	params := c.Request.Form
	midStr := params.Get("mid")
	tpStr := params.Get("tp")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var tp int64
	if business, ok := c.Get("business"); ok {
		tp = business.(int64)
	} else {
		if tp, err = strconv.ParseInt(tpStr, 10, 64); err != nil {
			tp = 1
		}
	}
	b, err := coinSvc.GetBusinessName(tp)
	if err != nil {
		c.JSON(nil, err)
	}
	arg := &pb.ListReq{Mid: mid, Business: b, Ts: time.Now().Unix()}
	c.JSON(coinSvc.List(c, arg))
}

func todayexp(c *bm.Context) {
	v := new(pb.TodayExpReq)
	if err := c.Bind(v); err != nil {
		return
	}
	res, err := coinSvc.TodayExp(c, v)
	c.JSONMap(map[string]interface{}{
		"number": res.Exp,
	}, err)
}

func updateSettle(c *bm.Context) {
	form := c.Request.Form
	aidStr := form.Get("aid")
	tpStr := form.Get("avtype")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	expSubStr := form.Get("exp_sub")
	expSub, err := strconv.ParseInt(expSubStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var tp int64
	if business, ok := c.Get("business"); ok {
		tp = business.(int64)
	} else {
		var err error
		tp, err = strconv.ParseInt(tpStr, 10, 64)
		if err != nil || tp <= 0 {
			tp = 1
		}
	}
	c.JSON(nil, coinSvc.UpdateSettle(c, aid, tp, expSub, form.Get("describe")))
}

func coins(c *bm.Context) {
	form := c.Request.Form
	midStr := form.Get("mid")
	upMidStr := form.Get("up_mid")
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	upMid, err := strconv.ParseInt(upMidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(coinSvc.AddedCoins(c, int64(mid), int64(upMid)))
}

func amend(c *bm.Context) {
	form := c.Request.Form
	aidStr := form.Get("aid")
	tpStr := form.Get("avtype")
	coinsStr := form.Get("coins")
	aid, err := strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var tp int64
	if business, ok := c.Get("business"); ok {
		tp = business.(int64)
	} else {
		var err1 error
		tp, err1 = strconv.ParseInt(tpStr, 10, 64)
		if err1 != nil || tp <= 0 {
			tp = 1
		}
	}
	coins, err := strconv.ParseInt(coinsStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, coinSvc.UpdateItemCoins(c, aid, tp, int64(coins)))
}

func ccounts(c *bm.Context) {
	params := new(struct {
		Aid    int64  `json:"aid" form:"aid" validate:"required,min=1"`
		Avtype int64  `json:"avtype" form:"avtype"`
		IP     string `json:"ip" `
	})
	if err := c.Bind(params); err != nil {
		return
	}
	if business, ok := c.Get("business"); ok {
		params.Avtype = business.(int64)
	}
	if params.Avtype == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	count, err := coinSvc.ItemCoin(c, params.Aid, params.Avtype)
	c.JSON(map[string]interface{}{
		"count": count,
	}, err)
}

// @params AddCoinReq
// @router get /x/internal/v1/coin/add
// @response AddCoinReply
func internalAddCoin(c *bm.Context) {
	v := new(pb.AddCoinReq)
	if err := c.Bind(v); err != nil {
		return
	}
	v.IP = metadata.String(c, metadata.RemoteIP)
	c.JSON(coinSvc.AddCoin(c, v))
}

// @params ItemUserCoinsReq
// @router get /x/internal/v1/coin/item/coins
// @response ItemUserCoinsReply
func itemCoins(c *bm.Context) {
	v := new(pb.ItemUserCoinsReq)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(coinSvc.ItemUserCoins(c, v))
}
