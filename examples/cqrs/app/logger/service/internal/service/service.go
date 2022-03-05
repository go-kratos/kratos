package service

import (
	v1 "kratos-cqrs/api/logger/service/v1"
	"kratos-cqrs/app/logger/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewLoggerService)

type LoggerService struct {
	v1.UnimplementedLoggerServer

	sensorData *biz.SensorDataUseCase
	sensor     *biz.SensorUseCase
	log        *log.Helper
}

func NewLoggerService(sensorData *biz.SensorDataUseCase, sensor *biz.SensorUseCase, logger log.Logger) *LoggerService {
	return &LoggerService{
		sensorData: sensorData,
		sensor:     sensor,
		log:        log.NewHelper(log.With(logger, "module", "service/logger-service"))}
}
