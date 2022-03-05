package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	svcV1 "kratos-cqrs/api/logger/service/v1"
)

type SensorDataRepo interface {
	InsertSensorData(ctx context.Context, req *svcV1.SensorData) error
	BatchInsertSensorData(ctx context.Context, sd []*svcV1.SensorData) error
}

type SensorDataUseCase struct {
	repo SensorDataRepo
	log  *log.Helper
}

func NewSensorDataUseCase(repo SensorDataRepo, logger log.Logger) *SensorDataUseCase {
	l := log.NewHelper(log.With(logger, "module", "sensor-data/usecase/logger-service"))
	return &SensorDataUseCase{repo: repo, log: l}
}

func (uc *SensorDataUseCase) InsertSensorData(ctx context.Context, req *svcV1.SensorData) error {
	return uc.repo.InsertSensorData(ctx, req)
}

func (uc *SensorDataUseCase) BatchInsertSensorData(ctx context.Context, req []*svcV1.SensorData) error {
	return uc.repo.BatchInsertSensorData(ctx, req)
}
