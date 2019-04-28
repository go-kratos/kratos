package favorite

import (
	"go-common/app/interface/main/app-interface/conf"
	artdao "go-common/app/interface/main/app-interface/dao/article"
	audiodao "go-common/app/interface/main/app-interface/dao/audio"
	bangumidao "go-common/app/interface/main/app-interface/dao/bangumi"
	bplusdao "go-common/app/interface/main/app-interface/dao/bplus"
	favdao "go-common/app/interface/main/app-interface/dao/favorite"
	malldao "go-common/app/interface/main/app-interface/dao/mall"
	spdao "go-common/app/interface/main/app-interface/dao/sp"
	ticketdao "go-common/app/interface/main/app-interface/dao/ticket"
	topicdao "go-common/app/interface/main/app-interface/dao/topic"
)

// Service is favorite.
type Service struct {
	c *conf.Config
	// dao
	favDao     *favdao.Dao
	artDao     *artdao.Dao
	spDao      *spdao.Dao
	topicDao   *topicdao.Dao
	bplusDao   *bplusdao.Dao
	audioDao   *audiodao.Dao
	bangumiDao *bangumidao.Dao
	ticketDao  *ticketdao.Dao
	mallDao    *malldao.Dao
}

// New new favorite。
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		// dao
		favDao:     favdao.New(c),
		topicDao:   topicdao.New(c),
		artDao:     artdao.New(c),
		spDao:      spdao.New(c),
		bplusDao:   bplusdao.New(c),
		audioDao:   audiodao.New(c),
		bangumiDao: bangumidao.New(c),
		ticketDao:  ticketdao.New(c),
		mallDao:    malldao.New(c),
	}
	return s
}
