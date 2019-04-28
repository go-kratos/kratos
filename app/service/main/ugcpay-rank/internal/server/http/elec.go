package http

import (
	"net/http"

	api "go-common/app/service/main/ugcpay-rank/api/http"
	"go-common/app/service/main/ugcpay-rank/internal/conf"
	"go-common/app/service/main/ugcpay-rank/internal/model"
	bm "go-common/library/net/http/blademaster"

	"github.com/json-iterator/go"
)

const (
	_contentTypeJSON = "application/json; charset=utf-8"
)

func elecMonthUP(ctx *bm.Context) {
	var (
		err  error
		arg  = &api.ArgRankElecMonthUP{}
		resp = &api.RetRankElecMonthUP{
			Data: &api.RespRankElecMonthUP{},
		}
		upRank *model.RankElecUPProto
		bytes  []byte
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.RankSize <= 0 || arg.RankSize > conf.Conf.Biz.ElecUPRankSize {
		arg.RankSize = conf.Conf.Biz.ElecUPRankSize
	}
	if upRank, err = svc.ElecMonthRankUP(ctx, arg.UPMID, arg.RankSize); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp.Data.Parse(upRank)

	if bytes, err = jsoniter.ConfigFastest.Marshal(resp); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, _contentTypeJSON, bytes)
}

func elecMonth(ctx *bm.Context) {
	var (
		err  error
		arg  = &api.ArgRankElecMonth{}
		resp = &api.RetRankElecMonth{
			Data: &api.RespRankElecMonth{},
		}
		upRank *model.RankElecUPProto
		avRank *model.RankElecAVProto
		bytes  []byte
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.RankSize <= 0 || arg.RankSize > conf.Conf.Biz.ElecAVRankSize {
		arg.RankSize = conf.Conf.Biz.ElecAVRankSize
	}
	if upRank, err = svc.ElecMonthRankUP(ctx, arg.UPMID, arg.RankSize); err != nil {
		ctx.JSON(nil, err)
		return
	}
	if avRank, err = svc.ElecMonthRankAV(ctx, arg.UPMID, arg.AVID, arg.RankSize); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp.Data.Parse(avRank, upRank)

	if bytes, err = jsoniter.ConfigFastest.Marshal(resp); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, _contentTypeJSON, bytes)
}

func elecAllAV(ctx *bm.Context) {
	var (
		err  error
		arg  = &api.ArgRankElecMonth{}
		resp = &api.RetRankElecAllAV{
			Data: &api.RespRankElecAllAV{},
		}
		avRank *model.RankElecAVProto
		bytes  []byte
	)
	if err = ctx.Bind(arg); err != nil {
		return
	}
	if arg.RankSize <= 0 || arg.RankSize > conf.Conf.Biz.ElecAVRankSize {
		arg.RankSize = conf.Conf.Biz.ElecAVRankSize
	}
	if avRank, err = svc.ElecTotalRankAV(ctx, arg.UPMID, arg.AVID, arg.RankSize); err != nil {
		ctx.JSON(nil, err)
		return
	}
	resp.Data.Parse(avRank)

	if bytes, err = jsoniter.ConfigFastest.Marshal(resp); err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.Bytes(http.StatusOK, _contentTypeJSON, bytes)
}
