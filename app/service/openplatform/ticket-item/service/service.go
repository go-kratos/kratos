package service

import (
	"context"

	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/app/service/openplatform/ticket-item/dao"

	validator "gopkg.in/go-playground/validator.v9"
)

var (
	v = validator.New()
)

// ItemService Service http service
type ItemService struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *ItemService) {
	s = &ItemService{
		c:   c,
		dao: dao.New(c),
	}
	return
}

// Ping check server ok
func (s *ItemService) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close dao
func (s *ItemService) Close() {
	s.dao.Close()
}
