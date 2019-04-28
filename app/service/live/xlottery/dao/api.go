package dao

import (
	giftApi "go-common/app/service/live/gift/api/liverpc"
	rcApi "go-common/app/service/live/rc/api/liverpc"
	userExtApi "go-common/app/service/live/userext/api/liverpc"
	"go-common/app/service/live/xlottery/conf"
	account "go-common/app/service/main/account/rpc/client"
	"go-common/library/net/rpc/liverpc"
)

// AccountApi liverpc user api
var AccountApi *account.Service3

// GiftApi liverpc gift api
var GiftApi *giftApi.Client

// RcApi rc api
var RcApi *rcApi.Client

// UserExtApi userext api
var UserExtApi *userExtApi.Client

// InitAPI init all service APIs
func InitAPI() {
	AccountApi = account.New3(nil)
	GiftApi = giftApi.New(getConf("gift"))
	RcApi = rcApi.New(getConf("rc"))
	UserExtApi = userExtApi.New(getConf("userext"))
}

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc
	if c != nil {
		return c[appName]
	}
	return nil
}
