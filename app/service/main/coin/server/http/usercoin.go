package http

import (
	"time"

	pb "go-common/app/service/main/coin/api"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// @params ModifyCoinsReq
// @router post /x/internal/v1/coin/user/modify
// @response ModifyCoinsReply
func modify(c *bm.Context) {
	arg := new(pb.ModifyCoinsReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(c, metadata.RemoteIP)
	arg.Ts = time.Now().Unix()
	c.JSON(coinSvc.ModifyCoins(c, arg))
}

// @params UserCoinsReq
// @router get /x/internal/v1/coin/user/count
// @response UserCoinsReply
func userCoins(c *bm.Context) {
	v := new(pb.UserCoinsReq)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(coinSvc.UserCoins(c, v))
}

// @params CoinsLogReq
// @router get /x/internal/v1/coin/user/log
// @response CoinsLogReply
func coinLog(c *bm.Context) {
	v := new(pb.CoinsLogReq)
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(coinSvc.CoinsLog(c, v))
}
