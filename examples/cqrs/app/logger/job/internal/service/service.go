package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "kratos-cqrs/api/logger/job/v1"
	"kratos-cqrs/app/logger/job/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewLoggerJobService)

type LoggerJobService struct {
	v1.UnimplementedLoggerJobServer

	log        *log.Helper
	sensor     *biz.SensorUseCase
	sensorData *biz.SensorDataUseCase
}

func NewLoggerJobService(sensor *biz.SensorUseCase, sensorData *biz.SensorDataUseCase, logger log.Logger) *LoggerJobService {
	return &LoggerJobService{
		sensor:     sensor,
		sensorData: sensorData,
		log:        log.NewHelper(log.With(logger, "module", "service/logger-job"))}
}
