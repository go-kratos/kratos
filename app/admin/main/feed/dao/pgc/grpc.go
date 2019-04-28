package pgc

import (
	"context"

	epgrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
)

//CardsInfoReply pgc grpc
func (d *Dao) CardsInfoReply(c context.Context, seasonIds []int32) (res map[int32]*seasongrpc.CardInfoProto, err error) {
	arg := &seasongrpc.SeasonInfoReq{
		SeasonIds: seasonIds,
	}
	info, err := d.rpcClient.Cards(c, arg)
	if err != nil {
		log.Error("d.rpcClient.Cards error(%v)", err)
		return nil, err
	}
	res = info.Cards
	return
}

//CardsEpInfoReply get pgc ep cards values by epid
func (d *Dao) CardsEpInfoReply(c context.Context, epIds []int32) (res map[int32]*epgrpc.EpisodeCardsProto, err error) {
	var epInfo *epgrpc.EpisodeCardsReply
	arg := &epgrpc.EpReq{
		Epids: epIds,
	}
	epInfo, err = d.epClient.Cards(c, arg)
	if err != nil {
		log.Error("d.rpcClient.Cards error(%v)", err)
		return nil, err
	}
	res = epInfo.Cards
	return
}
