package service

import (
	"context"

	"go-common/app/service/bbq/user/internal/conf"
	"go-common/app/service/bbq/user/internal/dao"
	"go-common/library/queue/databus"
)

// Service struct
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	userFaceSub *databus.Databus
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}

	if databusConfig, exists := c.Databus["bfs"]; exists {
		s.userFaceSub = databus.New(databusConfig)
		// 监听BFS databus
		go s.subBfsUserFace()
	}

	return s
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}
