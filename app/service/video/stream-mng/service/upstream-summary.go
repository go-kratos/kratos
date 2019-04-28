package service

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
)

// GetSummaryUpStreamRtmp 统计上行调度信息
func (s *Service) GetSummaryUpStreamRtmp(c context.Context, start int64, end int64) ([]*model.SummaryUpStreamRtmp, error) {
	return s.dao.GetSummaryUpStreamRtmp(c, start, end)
}

// GetSummaryUpStreamISP 统计ISP调度信息
func (s *Service) GetSummaryUpStreamISP(c context.Context, start int64, end int64) ([]*model.SummaryUpStreamRtmp, error) {
	return s.dao.GetSummaryUpStreamISP(c, start, end)
}

// GetSummaryUpStreamCountry 统计Country调度信息
func (s *Service) GetSummaryUpStreamCountry(c context.Context, start int64, end int64) ([]*model.SummaryUpStreamRtmp, error) {
	return s.dao.GetSummaryUpStreamCountry(c, start, end)
}

// GetSummaryUpStreamCountry 统计Platform调度信息
func (s *Service) GetSummaryUpStreamPlatform(c context.Context, start int64, end int64) ([]*model.SummaryUpStreamRtmp, error) {
	return s.dao.GetSummaryUpStreamPlatform(c, start, end)
}

// GetSummaryUpStreamCity 统计City调度信息
func (s *Service) GetSummaryUpStreamCity(c context.Context, start int64, end int64) ([]*model.SummaryUpStreamRtmp, error) {
	return s.dao.GetSummaryUpStreamCity(c, start, end)
}
