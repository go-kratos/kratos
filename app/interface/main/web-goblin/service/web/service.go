package web

import (
	"context"
	"time"

	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/app/interface/main/web-goblin/conf"
	"go-common/app/interface/main/web-goblin/dao/web"
	webmdl "go-common/app/interface/main/web-goblin/model/web"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/library/log"
)

const _chCardTypeAv = "av"

// Service struct .
type Service struct {
	c            *conf.Config
	dao          *web.Dao
	arc          *arcrpc.Service2
	tag          *tagrpc.Service
	maxAid       int64
	channelCards map[int64][]*webmdl.ChCard
}

// New init .
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: web.New(c),
		arc: arcrpc.New2(c.ArchiveRPC),
		tag: tagrpc.New2(c.TagRPC),
	}
	go s.justAID()
	go s.chCardproc()
	return s
}

// Ping Service .
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service .
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) chCardproc() {
	for {
		now := time.Now()
		cardMap, err := s.dao.ChCard(context.Background(), now)
		if err != nil {
			log.Error("chCardproc s.dao.ChCard() error(%v)", err)
			time.Sleep(time.Second)
		}
		l := len(cardMap)
		if l == 0 {
			time.Sleep(time.Duration(s.c.Rule.ChCardInterval))
			continue
		}
		tmp := make(map[int64][]*webmdl.ChCard, l)
		for channelID, card := range cardMap {
			for _, v := range card {
				if v.Type == _chCardTypeAv {
					tmp[channelID] = append(tmp[channelID], v)
				}
			}
		}
		s.channelCards = tmp
		time.Sleep(time.Duration(s.c.Rule.ChCardInterval))
	}
}
