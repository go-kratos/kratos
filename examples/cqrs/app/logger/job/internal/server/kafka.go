package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-transport/broker"
	"github.com/tx7do/kratos-transport/transport/kafka"

	"kratos-cqrs/app/logger/job/internal/conf"
	"kratos-cqrs/app/logger/job/internal/service"
)

// NewKafkaServer create a kafka server.
func NewKafkaServer(c *conf.Server, _ log.Logger, s *service.LoggerJobService) *kafka.Server {
	ctx := context.Background()

	srv := kafka.NewServer(
		broker.Addrs(c.Kafka.Addrs...),
		broker.OptionContext(ctx),
	)

	_ = srv.RegisterSubscriber("logger.sensor.ts",
		s.InsertSensorData,
		broker.SubscribeContext(ctx),
		broker.Queue("sensor_logger"),
		//broker.DisableAutoAck(),
	)

	_ = srv.RegisterSubscriber("logger.sensor.instance",
		s.InsertSensor,
		broker.SubscribeContext(ctx),
		broker.Queue("sensor"),
		//broker.DisableAutoAck(),
	)

	return srv
}
