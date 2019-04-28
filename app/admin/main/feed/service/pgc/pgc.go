package pgc

import (
	"context"

	"go-common/app/admin/main/feed/conf"
	pgcdao "go-common/app/admin/main/feed/dao/pgc"
	epgrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
)

// Service is egg service
type Service struct {
	pgc *pgcdao.Dao
}

// New new a egg service
func New(c *conf.Config) (s *Service) {
	var (
		b   *pgcdao.Dao
		err error
	)
	if b, err = pgcdao.New(c); err != nil {
		log.Error("pgcdao.New error(%v)", err)
		return
	}
	s = &Service{
		pgc: b,
	}
	return
}

//GetSeason get season from pgc
func (s *Service) GetSeason(c context.Context, seasonIDs []int32) (seasonCards map[int32]*seasongrpc.CardInfoProto, err error) {
	if seasonCards, err = s.pgc.CardsInfoReply(c, seasonIDs); err != nil {
		log.Error("%+v", err)
	}
	return
}

//GetEp get ep from pgc
func (s *Service) GetEp(c context.Context, epIds []int32) (res map[int32]*epgrpc.EpisodeCardsProto, err error) {
	if res, err = s.pgc.CardsEpInfoReply(c, epIds); err != nil {
		log.Error("%+v", err)
	}
	return
}
