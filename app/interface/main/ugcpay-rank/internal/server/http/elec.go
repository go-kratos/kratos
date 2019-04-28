package http

import (
	"net/http"

	api "go-common/app/interface/main/ugcpay-rank/api/http"
	"go-common/app/interface/main/ugcpay-rank/internal/conf"
	bm "go-common/library/net/http/blademaster"

	"github.com/json-iterator/go"
)

const (
	_contentTypeJSON = "application/json; charset=utf-8"
)

func elecMonthUP(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgRankElecMonthUP{}
		resp  *api.RespRankElecMonthUP
		bytes []byte
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.RankSize <= 0 || arg.RankSize > conf.Conf.Biz.ElecUPRankSize {
		arg.RankSize = conf.Conf.Biz.ElecUPRankSize
	}
	if resp, err = svc.RankElecMonthUP(ctx, arg.UPMID, arg.RankSize); err != nil {
		ctx.JSON(nil, err)
		return
	}

	if bytes, err = jsoniter.ConfigFastest.Marshal(resp); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, _contentTypeJSON, bytes)
}

func elecMonth(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgRankElecMonth{}
		resp  *api.RespRankElecMonth
		bytes []byte
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.RankSize <= 0 || arg.RankSize > conf.Conf.Biz.ElecAVRankSize {
		arg.RankSize = conf.Conf.Biz.ElecAVRankSize
	}
	if resp, err = svc.RankElecMonth(ctx, arg.UPMID, arg.AVID, arg.RankSize); err != nil {
		ctx.JSON(nil, err)
		return
	}

	if bytes, err = jsoniter.ConfigFastest.Marshal(resp); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, _contentTypeJSON, bytes)
}

func elecAllAV(ctx *bm.Context) {
	var (
		err   error
		arg   = &api.ArgRankElecMonth{}
		resp  *api.RespRankElecAllAV
		bytes []byte
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.RankSize <= 0 || arg.RankSize > conf.Conf.Biz.ElecAVRankSize {
		arg.RankSize = conf.Conf.Biz.ElecAVRankSize
	}
	if resp, err = svc.RankElecAllAV(ctx, arg.UPMID, arg.AVID, arg.RankSize); err != nil {
		ctx.JSON(nil, err)
		return
	}

	if bytes, err = jsoniter.ConfigFastest.Marshal(resp); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, _contentTypeJSON, bytes)
}
