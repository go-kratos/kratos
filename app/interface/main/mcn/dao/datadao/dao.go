package datadao

import (
	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/dao/global"
	"go-common/app/interface/main/mcn/tool/cache"
	"go-common/app/interface/main/mcn/tool/datacenter"
	"go-common/library/cache/memcache"
	bm "go-common/library/net/http/blademaster"
)

//Dao data dao
type Dao struct {
	Client    *datacenter.HttpClient
	Conf      *conf.Config
	mc        *memcache.Pool
	McWrapper *cache.MCWrapper
	bmClient  *bm.Client
}

//New .
func New(c *conf.Config) *Dao {
	return &Dao{
		Client:    datacenter.New(c.DataClientConf),
		Conf:      c,
		mc:        global.GetMc(),
		McWrapper: cache.New(global.GetMc()),
		bmClient:  global.GetBMClient(),
	}
}
