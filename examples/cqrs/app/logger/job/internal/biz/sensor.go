package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	svcV1 "kratos-cqrs/api/logger/service/v1"
)

type SensorRepo interface {
	CreateSensor(ctx context.Context, req *svcV1.Sensor) error
}

type SensorUseCase struct {
	repo SensorRepo
	log  *log.Helper
}

func NewSensorUseCase(repo SensorRepo, logger log.Logger) *SensorUseCase {
	l := log.NewHelper(log.With(logger, "module", "sensor/usecase/logger-service"))
	return &SensorUseCase{repo: repo, log: l}
}

func (uc *SensorUseCase) Create(ctx context.Context, req *svcV1.Sensor) error {
	return uc.repo.CreateSensor(ctx, req)
}
