package popular

import (
	"go-common/app/admin/main/feed/conf"
	showdao "go-common/app/admin/main/feed/dao/show"
)

// Service is search service
type Service struct {
	showDao *showdao.Dao
}

// New new a search service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		showDao: showdao.New(c),
	}
	return
}
