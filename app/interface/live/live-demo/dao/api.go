package dao

import (
	"go-common/app/interface/live/live-demo/conf"
	room_api "go-common/app/service/live/room/api/liverpc"
	"go-common/library/net/rpc/liverpc"
)

// RoomAPI liverpc room-service api
var RoomAPI *room_api.Client

// InitAPI init all service APIs
func InitAPI() {
	RoomAPI = room_api.New(getConf("room"))
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}
