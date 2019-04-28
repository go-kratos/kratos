package service

import (
	"context"
	"go-common/app/service/live/xuser/conf"
	expv1 "go-common/app/service/live/xuser/service/exp/v1"
	guardv1 "go-common/app/service/live/xuser/service/guard/v1"
	roomadminv1 "go-common/app/service/live/xuser/service/roomAdmin/v1"
	vipv1 "go-common/app/service/live/xuser/service/vip/v1"
)

// Service struct
type Service struct {
	c              *conf.Config
	vipv1svc       *vipv1.VipService
	guardv1svc     *guardv1.GuardService
	expv1svc       *expv1.UserExpService
	roomadminv1svc *roomadminv1.RoomAdminService
}

// New init
func New(c *conf.Config) (s *Service) {
	// init vip.v1 service
	s = &Service{
		c:              c,
		vipv1svc:       vipv1.New(c),
		guardv1svc:     guardv1.New(c),
		expv1svc:       expv1.NewUserExpService(c),
		roomadminv1svc: roomadminv1.NewRoomAdminService(c),
	}
	return s
}

// VipV1Svc return vip v1 service
func (s *Service) VipV1Svc() *vipv1.VipService {
	return s.vipv1svc
}

// GuardV1Svc return guard v1 service
func (s *Service) GuardV1Svc() *guardv1.GuardService {
	return s.guardv1svc
}

// ExpV1Svc return exp v1 service
func (s *Service) ExpV1Svc() *expv1.UserExpService {
	return s.expv1svc
}

// RoomAdminV1Svc return roomadmin v1 service
func (s *Service) RoomAdminV1Svc() *roomadminv1.RoomAdminService {
	return s.roomadminv1svc
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return nil
}

// Close Service
func (s *Service) Close() {
	s.expv1svc.Close()
	s.vipv1svc.Close()
}
