package account

import (
	"go-common/app/interface/main/tv/conf"
	accwar "go-common/app/service/main/account/api"
)

// Dao is account dao.
type Dao struct {
	// rpc
	accClient accwar.AccountClient
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{}
	var err error
	if d.accClient, err = accwar.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	return
}
