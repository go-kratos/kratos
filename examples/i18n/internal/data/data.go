package data

import (
	"github.com/SeeMusic/kratos/examples/i18n/internal/conf"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct { // TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		_ = logger.Log(log.LevelInfo, "closing the data resources")
	}
	return &Data{}, cleanup, nil
}
