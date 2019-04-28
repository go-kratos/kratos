package dao

import (
	"context"

	"go-common/app/interface/live/app-blink/conf"
	fans_medal_api "go-common/app/service/live/fans_medal/api/liverpc"
	relation_api "go-common/app/service/live/relation/api/liverpc"
	resource_cli "go-common/app/service/live/resource/api/grpc/v1"
	room_api "go-common/app/service/live/room/api/liverpc"
	user_api "go-common/app/service/live/user/api/liverpc"
	member_cli "go-common/app/service/main/member/api"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/liverpc"
)

// Dao dao
type Dao struct {
	c            *conf.Config
	RoomApi      *room_api.Client
	UserApi      *user_api.Client
	RelationApi  *relation_api.Client
	FansMedalApi *fans_medal_api.Client
	memberCli    member_cli.MemberClient
	titansCli    resource_cli.TitansClient
	HttpCli      *bm.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:            c,
		RoomApi:      room_api.New(getConf("room")),
		UserApi:      user_api.New(getConf("user")),
		RelationApi:  relation_api.New(getConf("relation")),
		FansMedalApi: fans_medal_api.New(getConf("fans_medal")),
		HttpCli:      bm.NewClient(c.HttpClient),
	}
	MemberCli, err := member_cli.NewClient(c.GrpcCli)
	if err != nil {
		panic(err)
	}
	dao.memberCli = MemberCli
	TitansCli, errTitans := resource_cli.NewClient(c.GrpcCli)
	if errTitans != nil {
		panic(err)
	}
	dao.titansCli = TitansCli
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

func getConf(appName string) *liverpc.ClientConfig {
	c := conf.Conf.LiveRpc

	if c != nil {
		return c[appName]
	}
	return nil
}
