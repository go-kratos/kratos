package gorpc

import (
	"context"
	"go-common/app/service/main/member/model"
)

const (
	_AddUserMonitor    = "RPC.AddUserMonitor"
	_IsInMonitor       = "RPC.IsInMonitor"
	_AddPropertyReview = "RPC.AddPropertyReview"
)

// AddUserMonitor is
func (s *Service) AddUserMonitor(c context.Context, arg *model.ArgAddUserMonitor) error {
	return s.client.Call(c, _AddUserMonitor, arg, &_noRes)
}

// IsInMonitor is
func (s *Service) IsInMonitor(c context.Context, arg *model.ArgMid) (bool, error) {
	isInMonitor := false
	if err := s.client.Call(c, _IsInMonitor, arg, &isInMonitor); err != nil {
		return false, err
	}
	return isInMonitor, nil
}

// AddPropertyReview is
func (s *Service) AddPropertyReview(c context.Context, arg *model.ArgAddPropertyReview) error {
	return s.client.Call(c, _AddPropertyReview, arg, &_noRes)
}
