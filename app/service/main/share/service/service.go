package service

import (
	"time"

	"go-common/app/service/main/share/conf"
	"go-common/app/service/main/share/dao"
	"go-common/app/service/main/share/model"
)

// Service is service.
type Service struct {
	c          *conf.Config
	dao        *dao.Dao
	allowType  map[int]string
	pubMsgType map[int]string
}

// New new a Service and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	s.allowType = make(map[int]string, len(s.c.AllowType))
	s.pubMsgType = map[int]string{
		model.BangumiTyp:  model.BangumiMsgTyp,
		model.ComicTyp:    model.ComicMsgTyp,
		model.ArchiveTyp:  model.ArchiveMsgTyp,
		model.PlaylistTyp: model.PlaylistMsgTyp,
	}
	for _, tp := range s.c.AllowType {
		s.allowType[tp] = s.pubMsgType[tp]
	}
	go s.loadConf()
	return
}

func (s *Service) loadConf() {
	for {
		time.Sleep(10 * time.Second)
		allowType := make(map[int]string, len(s.c.AllowType))
		for _, tp := range s.c.AllowType {
			allowType[tp] = s.pubMsgType[tp]
		}
		s.allowType = allowType
	}
}

// Ping ping
func (s *Service) Ping() (err error) {
	return s.dao.Ping()
}

// Close close
func (s *Service) Close() (err error) {
	return s.dao.Ping()
}
