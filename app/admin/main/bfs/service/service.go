package service

import (
	"context"

	"go-common/app/admin/main/bfs/conf"
	"go-common/app/admin/main/bfs/dao"
)

// Service struct
type Service struct {
	c *conf.Config
	d *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		d: dao.New(c),
	}
	return s
}

// Clusters .
func (s *Service) Clusters(c context.Context) (clusters []string) {
	for name := range s.c.Zookeepers {
		clusters = append(clusters, name)
	}
	return
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	return s.d.Ping(c)
}

// Close .
func (s *Service) Close() {
	s.d.Close()
}
