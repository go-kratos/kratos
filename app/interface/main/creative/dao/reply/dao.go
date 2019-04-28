package reply

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/reply"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
	"net/url"
	"strconv"
)

const (
	_replyMinfo   = "/x/internal/v2/reply/minfo"
	_replyRecover = "/x/internal/v2/reply/recover"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	replyMinfoURI   string
	replyRecoverURI string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:               c,
		client:          httpx.NewClient(c.HTTPClient.Slow),
		replyMinfoURI:   c.Host.API + _replyMinfo,
		replyRecoverURI: c.Host.API + _replyRecover,
	}
	return
}

// ReplyRecover recover reply
func (d *Dao) ReplyRecover(c context.Context, mid, oid, rpid int64, ip string) (err error) {
	params := url.Values{}
	params.Set("type", "1")
	params.Set("remark", "up主撤销协管员操作")
	params.Set("adid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("rpid", strconv.FormatInt(rpid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.replyRecoverURI, ip, params, &res); err != nil {
		log.Error("replyRecoverURI url(%s) response(%+v) error(%v)", d.replyRecoverURI+"?"+params.Encode(), res, err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("replyRecoverURI url(%s) error(%v)", d.replyRecoverURI+"?"+params.Encode(), err)
		return
	}
	return
}

// ReplyMinfo get multi reply info
func (d *Dao) ReplyMinfo(c context.Context, ak, ck string, mid, tp int64, DeriveIds, DeriveOids []int64, ip string) (ReplyMinfo map[int64]*reply.Reply, err error) {
	params := url.Values{}
	params.Set("type", strconv.FormatInt(tp, 10))
	params.Set("oid", xstr.JoinInts(DeriveOids))
	params.Set("rpid", xstr.JoinInts(DeriveIds))
	var res struct {
		Code int                    `json:"code"`
		Data map[int64]*reply.Reply `json:"data"`
	}
	if err = d.client.Get(c, d.replyMinfoURI, ip, params, &res); err != nil {
		log.Error("replyMinfoURI url(%s) response(%+v) error(%v)", d.replyMinfoURI+"?"+params.Encode(), res, err)
		return
	}
	if res.Code != 0 {
		log.Error("replyMinfoURI url(%s) res(%v)", d.replyMinfoURI+"?"+params.Encode(), res)
		err = ecode.Int(res.Code)
		return
	}
	ReplyMinfo = res.Data
	return
}
