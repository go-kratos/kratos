package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	svcV1 "kratos-cqrs/api/logger/service/v1"
	"kratos-cqrs/app/logger/job/internal/biz"
)

var _ biz.SensorRepo = (*sensorRepo)(nil)

type sensorRepo struct {
	data *Data
	log  *log.Helper
}

func NewSensorRepo(data *Data, logger log.Logger) biz.SensorRepo {
	return &sensorRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "sensor/repo/logger-service")),
	}
}

func (r *sensorRepo) CreateSensor(ctx context.Context, req *svcV1.Sensor) error {
	return r.data.db.Sensor.Create().
		SetType(req.Type).
		SetLocation(req.Location).
		Exec(ctx)
}
