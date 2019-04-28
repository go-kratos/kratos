package http

import (
	epgrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

//getPgcSeason GetSeason get season from pgc with grpc
func getPgcSeason(c *bm.Context) {
	var (
		err         error
		seasonCards map[int32]*seasongrpc.CardInfoProto
	)
	res := map[string]interface{}{}
	param := &struct {
		ID int32 `form:"id" validate:"required"`
	}{}
	if err = c.Bind(param); err != nil {
		return
	}
	v := []int32{param.ID}
	if seasonCards, err = pgcSvr.GetSeason(c, v); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(seasonCards, nil)
}

//getPgcSeasons GetSeasons get season from pgc with grpc
func getPgcSeasons(c *bm.Context) {
	var (
		err         error
		seasonCards map[int32]*seasongrpc.CardInfoProto
	)
	res := map[string]interface{}{}
	param := &struct {
		IDs []int32 `form:"ids,split" validate:"required,dive,gt=0"`
	}{}
	if err = c.Bind(param); err != nil {
		return
	}
	if seasonCards, err = pgcSvr.GetSeason(c, param.IDs); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(seasonCards, nil)
}

//getPgcEp GetSeasons get ep from pgc with grpc
func getPgcEp(c *bm.Context) {
	var (
		err     error
		epCards map[int32]*epgrpc.EpisodeCardsProto
	)
	res := map[string]interface{}{}
	param := &struct {
		IDs []int32 `form:"ids,split" validate:"required,dive,gt=0"`
	}{}
	if err = c.Bind(param); err != nil {
		return
	}
	if epCards, err = pgcSvr.GetEp(c, param.IDs); err != nil {
		res["message"] = "获取失败:" + err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	c.JSON(epCards, nil)
}
