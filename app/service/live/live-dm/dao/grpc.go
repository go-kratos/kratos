package dao

import (
	"context"

	broadcasrtService "go-common/app/service/live/broadcast-proxy/api/v1"
	"go-common/app/service/live/live-dm/conf"
	xuserService "go-common/app/service/live/xuser/api/grpc/v1"
	acctountService "go-common/app/service/main/account/api"
	filterService "go-common/app/service/main/filter/api/grpc/v1"
	locationService "go-common/app/service/main/location/api"
	spyService "go-common/app/service/main/spy/api"
	"go-common/library/net/rpc/warden"

	"google.golang.org/grpc"
)

//LocationAppID location服务注册ID
const locationAppID = "location.service"
const vipAppID = "live.xuser"

var (
	ac     acctountService.AccountClient
	vipCli xuserService.VipClient
	//FilterClient 屏蔽词过滤
	FilterClient filterService.FilterClient
	//LcClient 地理区域信息
	LcClient locationService.LocationClient
	//SpyClient 用户真实分
	SpyClient spyService.SpyClient
	//BcastClient  弹幕推送
	BcastClient *broadcasrtService.Client
	//UserExp 用户等级
	userExp *xuserService.Client
	//isAdmin 房管
	isAdmin xuserService.RoomAdminClient
)

//InitGrpc 初始化grpcclient
func InitGrpc(c *conf.Config) {
	var err error
	ac, err = acctountService.NewClient(c.AccClient)
	if err != nil {
		panic(err)
	}
	FilterClient, err = filterService.NewClient(c.FilterClient)
	if err != nil {
		panic(err)
	}
	LcClient, err = newLocationClient(c.LocationClient)
	if err != nil {
		panic(err)
	}
	vipCli, err = newVipService(c.XuserClent)
	if err != nil {
		panic(err)
	}
	SpyClient, err = spyService.NewClient(c.SpyClient)
	if err != nil {
		panic(err)
	}
	BcastClient, err = broadcasrtService.NewClient(c.BcastClient)
	if err != nil {
		panic(err)
	}
	userExp, err = xuserService.NewClient(c.UExpClient)
	if err != nil {
		panic(err)
	}
	isAdmin, err = xuserService.NewXuserRoomAdminClient(c.IsAdminClient)
	if err != nil {
		panic(err)
	}
}

//创建Location服务client
func newLocationClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (locationService.LocationClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+locationAppID)
	if err != nil {
		return nil, err
	}
	return locationService.NewLocationClient(conn), nil
}

func newVipService(cfg *warden.ClientConfig, opts ...grpc.DialOption) (xuserService.VipClient, error) {
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), "discovery://default/"+vipAppID)
	if err != nil {
		return nil, err
	}
	return xuserService.NewVipClient(conn), nil
}
