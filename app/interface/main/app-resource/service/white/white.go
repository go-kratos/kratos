package white

import (
	"go-common/app/interface/main/app-resource/conf"
)

// Service white service.
type Service struct {
	c         *conf.Config
	listCache map[string][]string
}

// New new a interest service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		listCache: c.White.List,
	}
	return
}

// List white list
func (s *Service) List() (res map[string][]string) {
	res = s.listCache
	return
}
