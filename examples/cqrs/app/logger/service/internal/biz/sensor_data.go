package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "kratos-cqrs/api/logger/service/v1"
)

type SensorDataRepo interface {
	ListSensorData(ctx context.Context, req *v1.ListSensorDataReq) (*v1.ListSensorDataReply, error)
	Avg(ctx context.Context, _ *v1.GetSensorAvgDataReq) (*v1.GetSensorAvgDataReply, error)
	AvgAndLatest(ctx context.Context, _ *v1.GetSensorAvgAndLatestDataReq) (*v1.GetSensorAvgAndLatestDataReply, error)
}

type SensorDataUseCase struct {
	repo SensorDataRepo
	log  *log.Helper
}

func NewSensorDataUseCase(repo SensorDataRepo, logger log.Logger) *SensorDataUseCase {
	l := log.NewHelper(log.With(logger, "module", "sensor-data/usecase/logger-service"))
	return &SensorDataUseCase{repo: repo, log: l}
}

func (uc *SensorDataUseCase) List(ctx context.Context, req *v1.ListSensorDataReq) (*v1.ListSensorDataReply, error) {
	return uc.repo.ListSensorData(ctx, req)
}

func (uc *SensorDataUseCase) Avg(ctx context.Context, req *v1.GetSensorAvgDataReq) (*v1.GetSensorAvgDataReply, error) {
	return uc.repo.Avg(ctx, req)
}

func (uc *SensorDataUseCase) AvgAndLatest(ctx context.Context, req *v1.GetSensorAvgAndLatestDataReq) (*v1.GetSensorAvgAndLatestDataReply, error) {
	return uc.repo.AvgAndLatest(ctx, req)
}
