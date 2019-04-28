package history

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-interface/conf"
	model "go-common/app/interface/main/app-interface/model/history"
	hismodle "go-common/app/interface/main/history/model"
	hisrpc "go-common/app/interface/main/history/rpc/client"
	artmodle "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	arcmodle "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	_apiPGC = "/internal_api/get_eps_v2"
)

// Dao is history dao
type Dao struct {
	client     *bm.Client
	historyRPC *hisrpc.Service
	arcRPC     *arcrpc.Service2
	artRPC     *artrpc.Service
	pgcAPI     string
}

// New initial history dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     bm.NewClient(c.HTTPClient),
		historyRPC: hisrpc.New(c.HistoryRPC),
		arcRPC:     arcrpc.New2(c.ArchiveRPC),
		artRPC:     artrpc.New(c.ArticleRPC),
		pgcAPI:     c.Host.Bangumi + _apiPGC,
	}
	return
}

// History get history
func (d *Dao) History(c context.Context, mid int64, pn, ps int) (res []*hismodle.Resource, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &hismodle.ArgHistories{Mid: mid, Pn: pn, Ps: ps, RealIP: ip}
	if res, err = d.historyRPC.History(c, arg); err != nil {
		err = errors.Wrapf(err, "d.historyRPC.History(%+v)", arg)
	}
	return
}

// Archive get archive info
func (d *Dao) Archive(c context.Context, aids []int64) (info map[int64]*arcmodle.View3, err error) {
	arg := &arcmodle.ArgAids2{Aids: aids}
	if info, err = d.arcRPC.Views3(c, arg); err != nil {
		err = errors.Wrapf(err, "d.arcRPC.Views3(%+v)", arg)
	}
	return
}

// PGC get PGC info
func (d *Dao) PGC(c context.Context, epid, platform string, build, mid int64) (info map[int64]*model.PGCRes, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("ep_ids", epid)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("build", strconv.FormatInt(build, 10))
	params.Set("platform", platform)
	var res struct {
		Code int             `json:"code"`
		Data []*model.PGCRes `json:"result"`
	}
	if err = d.client.Get(c, d.pgcAPI, ip, params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(res.Code), d.pgcAPI+"?"+params.Encode())
		return
	}
	info = make(map[int64]*model.PGCRes, len(res.Data))
	for _, v := range res.Data {
		v.Title = v.Season.Title
		info[v.EpID] = v
	}
	return
}

// Article get articl info
func (d *Dao) Article(c context.Context, articleIDs []int64) (info map[int64]*artmodle.Meta, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &artmodle.ArgAids{Aids: articleIDs, RealIP: ip}
	if info, err = d.artRPC.ArticleMetas(c, arg); err != nil {
		err = errors.Wrapf(err, "d.artRPC.ArticleMetas(%+v)", arg)
	}
	return
}

// HistoryByTP histroy by tp
func (d *Dao) HistoryByTP(c context.Context, mid int64, pn, ps int, tp int8) (res []*hismodle.Resource, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &hismodle.ArgHistories{Mid: mid, Pn: pn, Ps: ps, RealIP: ip, TP: tp}
	if res, err = d.historyRPC.History(c, arg); err != nil {
		err = errors.Wrapf(err, "d.historyRPC.History(%+v)", arg)
	}
	return
}

// Cursor 5.28游标由MaxOid+MaxTP唯一确定 改为 由ViewAt唯一确定（防止客户端改动对客户端仍用max字段）
func (d *Dao) Cursor(c context.Context, mid, max int64, ps int, tp int8, businesses []string) (res []*hismodle.Resource, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &hismodle.ArgCursor{Mid: mid, Max: max, Ps: ps, RealIP: ip, TP: tp, ViewAt: max, Businesses: businesses}
	if res, err = d.historyRPC.HistoryCursor(c, arg); err != nil {
		err = errors.Wrapf(err, "d.historyRPC.HistoryCursor(%+v)", arg)
	}
	return
}

// Del for history
func (d *Dao) Del(c context.Context, mid int64, hisRes []*hismodle.Resource) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &hismodle.ArgDelete{Mid: mid, RealIP: ip, Resources: hisRes}
	if err = d.historyRPC.Delete(c, arg); err != nil {
		err = errors.Wrapf(err, "d.historyRPC.Delete(%+v)", arg)
	}
	return
}

// Clear for history
func (d *Dao) Clear(c context.Context, mid int64, businesses []string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &hismodle.ArgClear{Mid: mid, RealIP: ip, Businesses: businesses}
	if err = d.historyRPC.Clear(c, arg); err != nil {
		err = errors.Wrapf(err, "d.historyRPC.Clear(%+v)", arg)
	}
	return
}
