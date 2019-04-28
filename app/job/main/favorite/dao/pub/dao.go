package pub

import (
	"go-common/app/job/main/favorite/conf"
	"go-common/library/queue/databus"
)

// Dao stat dao.
type Dao struct {
	databus2     *databus.Databus
	consumersMap map[int8]string
}

// New new a stat dao and return.
func New(c *conf.Config) *Dao {
	consumersMap := make(map[int8]string)
	for name, typ := range c.StatFavDatabus.Consumers {
		consumersMap[typ] = name
	}
	return &Dao{
		databus2:     databus.New(c.StatFavDatabus.Config),
		consumersMap: consumersMap,
	}
}
