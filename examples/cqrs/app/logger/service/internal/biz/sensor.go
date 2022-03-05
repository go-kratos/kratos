package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "kratos-cqrs/api/logger/service/v1"
)

type SensorRepo interface {
	ListSensor(ctx context.Context, req *v1.ListSensorReq) (*v1.ListSensorReply, error)
}

type SensorUseCase struct {
	repo SensorRepo
	log  *log.Helper
}

func NewSensorUseCase(repo SensorRepo, logger log.Logger) *SensorUseCase {
	l := log.NewHelper(log.With(logger, "module", "sensor/usecase/logger-service"))
	return &SensorUseCase{repo: repo, log: l}
}

func (uc *SensorUseCase) List(ctx context.Context, req *v1.ListSensorReq) (*v1.ListSensorReply, error) {
	return uc.repo.ListSensor(ctx, req)
}
