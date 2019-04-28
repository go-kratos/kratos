package data

import (
	"go-common/app/service/main/up/conf"
	"go-common/library/database/hbase.v2"
	"time"
)

//Dao hbase dao
type Dao struct {
	c            *conf.Config
	hbase        *hbase.Client
	hbaseTimeOut time.Duration
}

//New create dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		hbase:        hbase.NewClient(&c.HBase.Config),
		hbaseTimeOut: time.Millisecond * 500,
	}
	return
}
