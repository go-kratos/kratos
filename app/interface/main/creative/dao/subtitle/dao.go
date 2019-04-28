package subtitle

import (
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/dm2/rpc/client"
)

// Dao fn
type Dao struct {
	c   *conf.Config
	sub *client.Service
}

// New fn
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:   c,
		sub: client.New(c.SubRPC),
	}
	return
}
