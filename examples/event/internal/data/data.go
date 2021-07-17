package data

import (
	"github.com/go-kratos/kratos/examples/event/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	Event Event
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	event, err := NewKafkaEvent(WithKafkaAddr(c.Event.Addr),WithKafkaLogger(logger))
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = event.Close()
		_ = logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{Event:event}, cleanup, nil
}
