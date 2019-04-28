package http

import (
	pb "go-common/app/service/main/history/api/grpc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func add(c *bm.Context) {
	arg := new(pb.AddHistoryReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.AddHistory(c, arg))
}

func addHistories(c *bm.Context) {
	arg := new(pb.AddHistoriesReq)
	if err := c.BindWith(arg, binding.JSON); err != nil {
		return
	}
	c.JSON(srv.AddHistories(c, arg))
}

func del(c *bm.Context) {
	arg := new(pb.DelHistoriesReq)
	if err := c.BindWith(arg, binding.JSON); err != nil {
		return
	}
	c.JSON(srv.DelHistories(c, arg))
}

func clear(c *bm.Context) {
	arg := new(pb.ClearHistoryReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.ClearHistory(c, arg))
}
func userHistories(c *bm.Context) {
	arg := new(pb.UserHistoriesReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.UserHistories(c, arg))
}

func histories(c *bm.Context) {
	arg := new(pb.HistoriesReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.Histories(c, arg))
}

func userHide(c *bm.Context) {
	arg := new(pb.UserHideReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.UserHide(c, arg))
}

func updateHide(c *bm.Context) {
	arg := new(pb.UpdateUserHideReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(srv.UpdateUserHide(c, arg))
}
