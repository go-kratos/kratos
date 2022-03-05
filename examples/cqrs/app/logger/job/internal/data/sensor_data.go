package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	svcV1 "kratos-cqrs/api/logger/service/v1"
	"kratos-cqrs/app/logger/job/internal/biz"
	"kratos-cqrs/app/logger/job/internal/data/ent"
)

var _ biz.SensorDataRepo = (*sensorDataRepo)(nil)

type sensorDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewSensorDataRepo(data *Data, logger log.Logger) biz.SensorDataRepo {
	return &sensorDataRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "sensor-data/repo/logger-service")),
	}
}

func (r *sensorDataRepo) InsertSensorData(ctx context.Context, req *svcV1.SensorData) error {
	return r.data.db.SensorData.Create().
		SetTime(req.Ts).
		SetSensorID(int(req.SensorId)).
		SetCPU(req.Cpu).
		SetTemperature(req.Temperature).
		Exec(ctx)
}

func (r *sensorDataRepo) BatchInsertSensorData(ctx context.Context, req []*svcV1.SensorData) error {
	bulks := make([]*ent.SensorDataCreate, 0)
	for i := 0; i < len(req); i++ {
		s := req[i]
		bulk := r.data.db.SensorData.Create().
			SetTime(s.Ts).
			SetSensorID(int(s.SensorId)).
			SetCPU(s.Cpu).
			SetTemperature(s.Temperature)
		bulks = append(bulks, bulk)
	}
	return r.data.db.SensorData.CreateBulk(bulks...).Exec(ctx)
}
