package mis

import (
	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/app/service/openplatform/ticket-sales/dao"
)

//Mis http server
type Mis struct {
	c   *conf.Config
	dao *dao.Dao
}

// New for new mis obj
func New(c *conf.Config, d *dao.Dao) *Mis {
	m := &Mis{
		c:   c,
		dao: d,
	}
	return m
}
