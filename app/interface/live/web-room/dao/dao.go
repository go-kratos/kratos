package dao

import (
	"context"
	accrpc "go-common/app/service/main/account/rpc/client"

	"go-common/app/interface/live/web-room/conf"
	"go-common/app/service/live/xuser/api/grpc/v1"
)

// Dao dao
type Dao struct {
	c            *conf.Config
	RoomAdminAPI v1.RoomAdminClient
	acc          *accrpc.Service3
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	roomAdminClient, err := v1.NewXuserRoomAdminClient(conf.Conf.XRoomAdmin)
	if err != nil {
		panic(err)
	}
	dao = &Dao{
		c:            c,
		RoomAdminAPI: roomAdminClient,
		acc:          accrpc.New3(c.AccountRPC),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	return nil
}
