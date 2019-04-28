package reply

import (
	"context"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/article"
	"go-common/app/interface/main/creative/dao/elec"
	"go-common/app/interface/main/creative/dao/music"
	"go-common/app/interface/main/creative/dao/reply"
	"go-common/app/interface/main/creative/dao/search"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
)

// Service reply.
type Service struct {
	c     *conf.Config
	sear  *search.Dao
	acc   *account.Dao
	elec  *elec.Dao
	reply *reply.Dao
	arc   *archive.Dao
	art   *article.Dao
	mus   *music.Dao
}

// New get reply service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:     c,
		sear:  search.New(c),
		acc:   rpcdaos.Acc,
		elec:  elec.New(c),
		reply: reply.New(c),
		arc:   rpcdaos.Arc,
		art:   article.New(c),
		mus:   music.New(c),
	}
	return s
}

// Ping service.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.sear.Ping(c); err != nil {
		log.Error("s.archive.Dao.PingDb err(%v)", err)
	}
	return
}
