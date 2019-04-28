package service

import (
	"context"

	"go-common/app/interface/openplatform/seo/conf"
	"go-common/app/interface/openplatform/seo/dao"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return s
}

// Ping .
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close .
func (s *Service) Close() {
	s.dao.Close()
}

// Sitemap 生成站点地图
func (s *Service) Sitemap(c context.Context, host string) ([]byte, error) {
	return s.dao.Sitemap(c, host)
}
