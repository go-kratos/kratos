package favorite

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/model"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is account dao.
type Dao struct {
	favRPC *favrpc.Service // rpc
	conf   *conf.Config
	client *bm.Client
}

// New account dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		favRPC: favrpc.New2(c.FavoriteRPC),
		conf:   c,
		client: bm.NewClient(c.HTTPClient),
	}
	return
}

const (
	_FavBusiness = 2
	_DefaultFav  = 0
)

// FavoriteV3 picks favorite info from rpc
func (d *Dao) FavoriteV3(ctx context.Context, mid int64, pn int) (res *favmdl.Favorites, err error) {
	var ip = metadata.String(ctx, metadata.RemoteIP)
	arg := &favmdl.ArgFavs{
		Type:   _FavBusiness,
		Mid:    mid,
		Fid:    _DefaultFav,
		Tv:     1,
		Pn:     pn,
		Ps:     d.conf.Cfg.FavPs,
		RealIP: ip,
	}
	if res, err = d.favRPC.Favorites(ctx, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// favAct adds/deletes favorite into/from the default folder
func (d *Dao) favAct(ctx context.Context, mid int64, aid int64, host string) (err error) {
	var (
		ip     = metadata.String(ctx, metadata.RemoteIP)
		params = url.Values{}
		res    = model.RespFavAct{}
	)
	params.Set("mid", fmt.Sprintf("%d", mid))
	params.Set("aid", fmt.Sprintf("%d", aid))
	if err = d.client.Post(ctx, host, ip, params, &res); err != nil {
		log.Error("FavAdd Aid %d, Mid %d, Err %v", aid, mid, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), fmt.Sprintf("Fav AID %d, Mid %d, API Error %s", aid, mid, res.Message))
		log.Error("FavAdd ERROR:%v, URL: %s", err, host+"?"+params.Encode())
		return
	}
	return
}

// FavAdd def.
func (d *Dao) FavAdd(ctx context.Context, mid, aid int64) (err error) {
	host := d.conf.Host.FavAdd
	return d.favAct(ctx, mid, aid, host)
}

// FavDel deletes favorite from the default folder
func (d *Dao) FavDel(ctx context.Context, mid int64, aid int64) (err error) {
	host := d.conf.Host.FavDel
	return d.favAct(ctx, mid, aid, host)
}

// InDefault returns whether the aid is in Default of Mid
func (d *Dao) InDefault(ctx context.Context, mid int64, aid int64) (bool, error) {
	var ip = metadata.String(ctx, metadata.RemoteIP)
	arg := &favmdl.ArgInDefaultFolder{
		Type:   _FavBusiness,
		Mid:    mid,
		RealIP: ip,
		Oid:    aid,
	}
	return d.favRPC.InDefault(ctx, arg)
}
