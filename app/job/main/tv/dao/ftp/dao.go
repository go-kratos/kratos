package ftp

import (
	"go-common/app/job/main/tv/conf"
)

// Dao dao.
type Dao struct {
	conf *conf.Config
}

// New create a instance of Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		conf: c,
	}
	return
}
