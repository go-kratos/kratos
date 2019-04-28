package space

import (
	"context"
	"encoding/json"
	"go-common/library/cache/redis"
	"net/url"
	"strconv"

	article "go-common/app/interface/openplatform/article/model"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	"go-common/app/job/main/app/conf"
	"go-common/app/job/main/app/model/space"
	"go-common/app/service/main/archive/api"
	arcrpc "go-common/app/service/main/archive/api/gorpc"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const _contributeUp = "/x/v2/space/upContribute"

// Dao is favorite dao
type Dao struct {
	client     *httpx.Client
	clientAsyn *httpx.Client
	clipList   string
	albumList  string
	audioList  string
	// app
	contributeUp string
	// rpc
	arcRPC *arcrpc.Service2
	artRPC *artrpc.Service
	// redis
	redis *redis.Pool
}

// New initial favorite dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client:     httpx.NewClient(c.HTTPClient),
		clientAsyn: httpx.NewClient(c.HTTPClientAsyn),
		clipList:   c.Host.VC + _clipList,
		albumList:  c.Host.VC + _albumList,
		audioList:  c.Host.APICo + _audioList,
		// app
		contributeUp: c.Host.APP + _contributeUp,
		// rpc
		arcRPC: arcrpc.New2(c.ArchiveRPC),
		artRPC: artrpc.New(c.ArticleRPC),
		// redis
		redis: redis.NewPool(c.Redis.Contribute.Config),
	}
	return
}

func (d *Dao) UpContributeCache(c context.Context, vmid int64, attrs *space.Attrs, items []*space.Item) (err error) {
	var b []byte
	params := url.Values{}
	params.Set("vmid", strconv.FormatInt(vmid, 10))
	if b, err = json.Marshal(attrs); err != nil {
		err = errors.Wrapf(err, "%v", attrs)
		return
	}
	params.Set("attrs", string(b))
	if b, err = json.Marshal(items); err != nil {
		err = errors.Wrapf(err, "%v", items)
		return
	}
	params.Set("items", string(b))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.contributeUp, "", params, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		if res.Code == ecode.RequestErr.Code() {
			return
		}
		err = errors.Wrap(err, d.contributeUp+"?"+params.Encode())
	}
	return
}

// UpArchives get upper archives
func (d *Dao) UpArchives(c context.Context, mid int64, pn, ps int, ip string) (as []*api.Arc, err error) {
	arg := &archive.ArgUpArcs2{Mid: mid, Pn: pn, Ps: ps, RealIP: ip}
	if as, err = d.arcRPC.UpArcs3(c, arg); err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// UpArcCnt get upper count.
func (d *Dao) UpArcCnt(c context.Context, mid int64, ip string) (cnt int, err error) {
	arg := &archive.ArgUpCount2{Mid: mid}
	if cnt, err = d.arcRPC.UpCount2(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}

// UpArticles get article data from api.
func (d *Dao) UpArticles(c context.Context, mid int64, pn, ps int) (ts []*article.Meta, count int, err error) {
	var res *article.UpArtMetas
	arg := &article.ArgUpArts{Mid: mid, Pn: pn, Ps: ps}
	if res, err = d.artRPC.UpArtMetas(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	if res != nil {
		ts = res.Articles
		count = res.Count
	}
	return
}
