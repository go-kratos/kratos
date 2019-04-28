package service

import (
	"context"
	"go-common/app/service/live/xuserex/service/v1"

	"go-common/app/service/live/xuserex/conf"
)

// Service struct
type Service struct {
	c               *conf.Config
	roomNotivev1svc *v1.RoomNoticeService
}

// New init
func New(c *conf.Config) (s *Service) {
	// init vip.v1 service
	s = &Service{
		c:               c,
		roomNotivev1svc: v1.NewRoomNoticeService(conf.Conf),
	}
	return s
}

// RoomNoticeV1Svc return roomadmin v1 service
func (s *Service) RoomNoticeV1Svc() *v1.RoomNoticeService {
	return s.roomNotivev1svc
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return nil
}

// Close Service
func (s *Service) Close() {
}
