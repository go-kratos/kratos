package dao

import (
	"go-common/app/job/live/push-search/conf"
	userApi "go-common/app/service/live/user/api/liverpc"
	relationApi "go-common/app/service/live/relation/api/liverpc"
	roomApi "go-common/app/service/live/room/api/liverpc"
	"go-common/library/net/rpc/liverpc"
)

var UserApi *userApi.Client
var RelationApi *relationApi.Client
var RoomApi *roomApi.Client

// InitAPI init all service APIs
func InitAPI() {
	UserApi = userApi.New(getConf("user"))
	RelationApi = relationApi.New(getConf("relation"))
	RoomApi = roomApi.New(getConf("room"))
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}