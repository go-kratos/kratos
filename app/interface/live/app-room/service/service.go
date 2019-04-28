package service

import (
	"go-common/app/interface/live/app-room/conf"
)

// Service struct
type Service struct {
	c *conf.Config
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
	}
	return s
}

// Close -
func (s *Service) Close() error {
	return nil
}
