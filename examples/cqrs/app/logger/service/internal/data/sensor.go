package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "kratos-cqrs/api/logger/service/v1"
	"kratos-cqrs/app/logger/service/internal/biz"
	"kratos-cqrs/app/logger/service/internal/data/ent"
	"kratos-cqrs/app/logger/service/internal/data/ent/sensor"
	paging "kratos-cqrs/pkg/util/pagination"
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

func (r *sensorRepo) ListSensor(ctx context.Context, req *v1.ListSensorReq) (*v1.ListSensorReply, error) {
	sensors, err := r.data.db.Sensor.Query().
		Offset(paging.GetPageOffset(req.GetPage(), req.GetPageSize())).
		Limit(int(req.GetPageSize())).
		Order(ent.Asc(sensor.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.Sensor, 0)
	for _, u := range sensors {
		item := v1.Sensor{
			Id:       u.ID,
			Type:     u.Type,
			Location: u.Location,
		}
		items = append(items, &item)
	}

	return &v1.ListSensorReply{
		Results: items,
	}, nil
}
