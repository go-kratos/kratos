package dao

import (
	"context"

	"go-common/app/job/live/push-search/conf"
	"go-common/library/queue/databus"
	"go-common/library/database/hbase.v2"
)

// Dao dao
type Dao struct {
	c                 *conf.Config
	RoomInfoDataBus   *databus.Databus
	AttentionDataBus  *databus.Databus
	UserNameDataBus   *databus.Databus
	PushSearchDataBus *databus.Databus
	SearchHBase       *hbase.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                 c,
		RoomInfoDataBus:   databus.New(c.DataBus.RoomInfo),
		AttentionDataBus:  databus.New(c.DataBus.Attention),
		UserNameDataBus:   databus.New(c.DataBus.UserName),
		PushSearchDataBus: databus.New(c.DataBus.PushSearch),
		SearchHBase:       hbase.NewClient(&c.SearchHBase.Config),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.RoomInfoDataBus.Close()
	d.AttentionDataBus.Close()
	d.UserNameDataBus.Close()
	return
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return nil
}
