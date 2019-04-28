package service

import (
	"context"

	"go-common/app/interface/main/report-click/model"
)

// ErrReport reports the failures of calling the api "heartbeat/mobile"
func (s *Service) ErrReport(ctx context.Context, req *model.ErrReport) {
	s.promErr.Incr(req.ToProm())
}

// SuccReport reports the success of calling the api "heartbeat/mobile"
func (s *Service) SuccReport(ctx context.Context, req *model.SuccReport) {
	s.promInfo.Incr(req.ToProm())
}
