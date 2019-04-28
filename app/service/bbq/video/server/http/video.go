package http

import (
	"go-common/app/service/bbq/video/api/grpc/v1"
	grpc "go-common/app/service/bbq/video/api/grpc/v1"
	httpV1 "go-common/app/service/bbq/video/api/http/v1"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// example for http request handler
func bvcTransBack(c *bm.Context) {
	arg := new(v1.BVCTransBackRequset)
	err := bindJSON(c, arg)
	if err == nil {
		err = srv.BVCTransRes(c, arg)
	}
	c.JSON(nil, err)
}

func bvcTransCommit(c *bm.Context) {
	arg := new(v1.BVideoTransRequset)
	err := c.Bind(arg)
	if err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.BVCTransCommit(c, arg))
}

func createID(c *bm.Context) {
	arg := new(v1.CreateIDRequest)
	err := c.Bind(arg)
	if err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.CreateID(c, arg))
}

func videoViewsAdd(c *bm.Context) {
	arg := new(httpV1.ViewsAddRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.VideoViewsAdd(c, arg))
}

func videoStat(c *bm.Context) {
	arg := new(grpc.SvStatisticsInfoReq)
	if err := bindJSON(c, arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.SvStatisticsInfo(c, arg))
}

func limitsModify(c *bm.Context) {
	arg := new(grpc.ModifyLimitsRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.ModifyLimits(c, arg))
}

func svPlays(c *bm.Context) {
	arg := new(grpc.PlayInfoRequest)
	if err := bindJSON(c, arg); err != nil {
		log.Error("param err:%v", err)
		return
	}
	c.JSON(srv.PlayInfo(c, arg))
}
