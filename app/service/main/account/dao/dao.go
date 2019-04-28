package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	v1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/conf"
	member "go-common/app/service/main/member/api/gorpc"
	mmodel "go-common/app/service/main/member/model"
	usersuit "go-common/app/service/main/usersuit/rpc/client"
	"go-common/library/cache/memcache"
	"go-common/library/database/elastic"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/sync/pipeline/fanout"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	Info(c context.Context, key int64) (*v1.Info, error)
	//cache: -batch=50 -max_group=10 -batch_err=continue
	Infos(c context.Context, keys []int64) (map[int64]*v1.Info, error)
	Card(c context.Context, key int64) (*v1.Card, error)
	//cache: -batch=50 -max_group=10 -batch_err=continue
	Cards(c context.Context, keys []int64) (map[int64]*v1.Card, error)
	Vip(c context.Context, key int64) (*v1.VipInfo, error)
	//cache: -batch=50 -max_group=10 -batch_err=continue
	Vips(c context.Context, keys []int64) (map[int64]*v1.VipInfo, error)
	Profile(c context.Context, key int64) (*v1.Profile, error)
}

const (
	_nameURL           = "/api/member/getInfoByName"
	_vipInfoURL        = "/internal/v1/user/%d"
	_vipMultiInfoURL   = "/internal/v1/user/list"
	_passportDetailURL = "/intranet/acc/detail"
	_passportProfile   = "/intranet/acc/queryByMid"
)

// Dao dao.
type Dao struct {
	// memcache
	mc       *memcache.Pool
	mcExpire int32
	// cache async save
	cache *fanout.Fanout
	// rpc
	mRPC    *member.Service
	suitRPC *usersuit.Service2
	// http
	httpR *bm.Client
	httpW *bm.Client
	httpP *bm.Client
	// api
	detailURI  string
	profileURI string
	nameURI    string
	// vip api
	vipInfoURI      string
	vipMultiInfoURI string
	//es
	es *elastic.Elastic
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// account memcache
		mc:       memcache.NewPool(c.Memcache.Account),
		mcExpire: int32(time.Duration(c.Memcache.AccountExpire) / time.Second),
		// cache chan
		cache: fanout.New("accountServiceCache", fanout.Worker(1), fanout.Buffer(1024)),
		// rpc
		// mRPC:    member.New(c.MemberRPC),
		mRPC:    member.New(c.MemberRPC),
		suitRPC: usersuit.New(c.SuitRPC),
		// http read&write client
		httpR: bm.NewClient(c.HTTPClient.Read),
		httpW: bm.NewClient(c.HTTPClient.Write),
		httpP: bm.NewClient(c.HTTPClient.Privacy),
		es:    elastic.NewElastic(c.Elastic),
		// api
		nameURI: c.Host.AccountURI + _nameURL,
		// vip api
		vipInfoURI:      c.Host.VipURI + _vipInfoURL,
		vipMultiInfoURI: c.Host.VipURI + _vipMultiInfoURL,
		//passport
		detailURI:  c.Host.PassportURI + _passportDetailURL,
		profileURI: c.Host.PassportURI + _passportProfile,
	}
	return
}

// LevelExp get member level exp.
func (d *Dao) LevelExp(c context.Context, mid int64) (lexp *mmodel.LevelInfo, err error) {
	lexp, err = d.mRPC.Exp(c, &mmodel.ArgMid2{Mid: mid})
	return
}

// AddMoral add moral.
func (d *Dao) AddMoral(c context.Context, arg *mmodel.ArgUpdateMoral) (err error) {
	return d.mRPC.AddMoral(c, arg)
}

// UpdateExp update exp.
func (d *Dao) UpdateExp(c context.Context, arg *mmodel.ArgAddExp) error {
	return d.mRPC.UpdateExp(c, arg)
}

// Ping check connection success.
func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.mc.Get(c)
	err = conn.Set(&memcache.Item{
		Key:   "ping",
		Value: []byte("pong"),
	})
	conn.Close()
	return
}

// Close close memcache resource.
func (d *Dao) Close() {
	if d.mc != nil {
		d.mc.Close()
	}
}

func fullImage(mid int64, image string) string {
	if len(image) == 0 {
		return ""
	}
	if strings.HasPrefix(image, "http://") {
		return image
	}
	return fmt.Sprintf("http://i%d.hdslb.com%s", mid%3, image)
}
