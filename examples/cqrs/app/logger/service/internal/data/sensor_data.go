package data

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	v1 "kratos-cqrs/api/logger/service/v1"
	"kratos-cqrs/app/logger/service/internal/biz"
	"kratos-cqrs/app/logger/service/internal/data/ent"
	"kratos-cqrs/app/logger/service/internal/data/ent/sensordata"
	paging "kratos-cqrs/pkg/util/pagination"
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

func (r *sensorDataRepo) ListSensorData(ctx context.Context, req *v1.ListSensorDataReq) (*v1.ListSensorDataReply, error) {
	sensors, err := r.data.db.SensorData.Query().
		Offset(paging.GetPageOffset(req.GetPage(), req.GetPageSize())).
		Limit(int(req.GetPageSize())).
		Order(ent.Asc(sensordata.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.SensorData, 0)
	for _, u := range sensors {
		item := v1.SensorData{
			SensorId:    int64(u.ID),
			Ts:          *u.Time,
			Temperature: u.Temperature,
			Cpu:         u.CPU,
		}
		items = append(items, &item)
	}

	return &v1.ListSensorDataReply{
		Results: items,
	}, nil
}

func (r *sensorDataRepo) Avg(ctx context.Context, _ *v1.GetSensorAvgDataReq) (*v1.GetSensorAvgDataReply, error) {
	var v []*v1.SensorAvgData
	err := r.data.db.SensorData.Query().
		Modify(func(s *sql.Selector) {
			s.Select(
				sql.As("time_bucket(1800000, time)", "period"),
				sql.As(sql.Avg(sensordata.FieldTemperature), "avg_temp"),
				sql.As(sql.Avg(sensordata.FieldCPU), "avg_cpu"),
			).
				GroupBy("period")
		}).
		Scan(ctx, &v)
	var results v1.GetSensorAvgDataReply
	results.Results = v
	return &results, err
}

func (r *sensorDataRepo) AvgAndLatest(ctx context.Context, _ *v1.GetSensorAvgAndLatestDataReq) (*v1.GetSensorAvgAndLatestDataReply, error) {
	var v []*v1.SensorAvgAndLatestData
	err := r.data.db.SensorData.Query().
		Modify(func(s *sql.Selector) {
			s.Select(
				sql.As("time_bucket(1800000, time)", "period"),
				sql.As(sql.Avg(sensordata.FieldTemperature), "avg_temp"),
				sql.As(sql.Avg(sensordata.FieldCPU), "avg_cpu"),
				sql.As("last(temperature, time)", "last_temp"),
			).
				GroupBy("period")
		}).
		Scan(ctx, &v)
	var results v1.GetSensorAvgAndLatestDataReply
	results.Results = v
	return &results, err
}
