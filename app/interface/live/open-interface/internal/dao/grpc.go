package dao

import (
	broadcasrtService "go-common/app/service/live/broadcast-proxy/api/v1"
)

var (
	//BcastClient  弹幕服务
	BcastClient *broadcasrtService.Client
)

//InitGrpc 初始化grpcclient
func InitGrpc() {
	var err error
	BcastClient, err = broadcasrtService.NewClient(nil)
	if err != nil {
		panic(err)
	}
}
