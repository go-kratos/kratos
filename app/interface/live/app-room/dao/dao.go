package dao

import (
	"context"

	"go-common/app/interface/live/app-room/conf"
	userextApi "go-common/app/service/live/userext/api/liverpc"
	xUserEx "go-common/app/service/live/xuserex/api/grpc/v1"
	"go-common/library/cache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_payCenterWalletURL = "/wallet-int/wallet/getUserWalletInfo"
	_liveWalletURL      = "/x/internal/livewallet/wallet/getAll"
)

// Dao dao
type Dao struct {
	c                  *conf.Config
	payCenterWalletURL string
	payCenterClient    *bm.Client
	liveWalletURL      string
	liveWalletClient   *bm.Client
	UserExtAPI         *userextApi.Client
	giftCache          *cache.Cache
	XuserexAPI         *xUserEx.Client
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:                  c,
		payCenterWalletURL: c.Host.PayCenter + _payCenterWalletURL,
		payCenterClient:    bm.NewClient(c.HTTPClient.PayCenter),
		liveWalletURL:      c.Host.LiveRpc + _liveWalletURL,
		liveWalletClient:   bm.NewClient(c.HTTPClient.LiveRpc),
		giftCache:          cache.New(1, 1024),
	}
	xUserexApi, err := xUserEx.NewClient(c.Warden)
	if err != nil {
		log.Error("init xuserex error(%v)", err)
	}
	dao.XuserexAPI = xUserexApi

	InitAPI(dao)
	return
}

// Close close the resource.
func (d *Dao) Close() {
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) error {
	// TODO: if you need use mc,redis, please add
	// check
	return nil
}
