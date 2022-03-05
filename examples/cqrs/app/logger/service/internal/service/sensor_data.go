package service

import (
	"context"
	v1 "kratos-cqrs/api/logger/service/v1"
)

func (s *LoggerService) ListSensorData(ctx context.Context, req *v1.ListSensorDataReq) (*v1.ListSensorDataReply, error) {
	return s.sensorData.List(ctx, req)
}

func (s *LoggerService) GetSensorAvgData(ctx context.Context, req *v1.GetSensorAvgDataReq) (*v1.GetSensorAvgDataReply, error) {
	return s.sensorData.Avg(ctx, req)
}
func (s *LoggerService) GetSensorAvgAndLatestData(ctx context.Context, req *v1.GetSensorAvgAndLatestDataReq) (*v1.GetSensorAvgAndLatestDataReply, error) {
	return s.sensorData.AvgAndLatest(ctx, req)
}
