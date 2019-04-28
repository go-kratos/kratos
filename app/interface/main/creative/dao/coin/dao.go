package coin

import (
	"go-common/app/interface/main/creative/conf"
	coinclient "go-common/app/service/main/coin/api"
)

// Dao str
type Dao struct {
	c          *conf.Config
	coinClient coinclient.CoinClient
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
	}
	var err error
	if d.coinClient, err = coinclient.NewClient(c.CoinClient); err != nil {
		panic(err)
	}
	return
}
