package service

import (
	"context"
	v1 "kratos-cqrs/api/logger/service/v1"
)

func (s *LoggerService) ListSensor(ctx context.Context, req *v1.ListSensorReq) (*v1.ListSensorReply, error) {
	return s.sensor.List(ctx, req)
}
