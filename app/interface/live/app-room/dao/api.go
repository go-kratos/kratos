package dao

import (
	"go-common/app/interface/live/app-room/conf"
	userextApi "go-common/app/service/live/userext/api/liverpc"
	"go-common/library/net/rpc/liverpc"
)

// InitAPI init all service APIs
func InitAPI(dao *Dao) {
	dao.UserExtAPI = userextApi.New(getConf("userext"))
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}
