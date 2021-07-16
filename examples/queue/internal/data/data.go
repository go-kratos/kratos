package data

import (
	"github.com/go-kratos/kratos/examples/queue/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO warpped database client
	Kafka  *Kafka
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	kafka,err := NewKafka(c.Kafka.Addr, logger)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		kafka.close()
		logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{Kafka: kafka}, cleanup, nil
}
