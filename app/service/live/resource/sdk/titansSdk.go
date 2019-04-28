package titansSdk

import (
	"context"
	"go-common/app/service/live/resource/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
	"sync/atomic"
	"time"
)

// titansConfigSdk 初始化时 单实例对象
var titansConfigSdk = &titansSdk{}

// Config 服务可配置项
type Config struct {
	// 服务的tree_id
	TreeId int64
	// 缓存的更新间隔，单位为s，不配置则为5s
	Expire int64
}

// titansSdk sdk
type titansSdk struct {
	client v1.TitansClient
	cache  atomic.Value
}

// Init 业务初始化sdk接口
func Init(titansConfig *Config) {
	conf := &warden.ClientConfig{}
	client := warden.NewClient(conf)
	conn, err := client.Dial(context.Background(), "discovery://default/"+v1.AppID)
	if err != nil {
		panic("依赖TitansConfig, 但是Titans的Client没有创建成功！")
	}
	if titansConfig.TreeId == 0 {
		panic("依赖了TitansConfig, 但是titansConfig的treeId配置为空！")
	}
	titansConfigSdk = &titansSdk{}
	titansConfigSdk.client = v1.NewTitansClient(conn)
	titansConfigSdk.cache.Store(make(map[string]string))
	load(titansConfigSdk, titansConfig.TreeId)
	go update(titansConfigSdk, titansConfig.TreeId, titansConfig.Expire)
}

// update 定时更新
func update(sdk *titansSdk, treeId int64, expire int64) {
	if expire < 1 {
		expire = 5
	}
	for {
		load(sdk, treeId)
		time.Sleep(time.Duration(expire) * time.Second)
	}
}

// load grpc接口
func load(sdk *titansSdk, treeId int64) {
	for {
		resp, err := sdk.client.GetByTreeId(context.Background(), &v1.TreeIdReq{TreeId: treeId})
		if err != nil {
			log.Error("[SyncTitansConfig][call resource][error], err:%+v", err)
		} else {
			sdk.cache.Store(resp.List)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// Get 获取配置接口
func Get(keyword string) (res string, err error) {
	res = ""
	cache, ok := titansConfigSdk.cache.Load().(map[string]string)
	if !ok {
		log.Error("[GetTitansConfig][cache content exception][error] content assert failed")
		err = ecode.GetConfAdminErr
		return
	}
	res, ok = cache[keyword]
	if !ok || res == "" {
		log.Error("[GetTitansConfig][cache content empty][Warn], keyword:%s, cache:%v", keyword, cache)
		err = ecode.GetConfAdminErr
		return
	}
	return
}
