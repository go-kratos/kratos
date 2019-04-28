package relation

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/service/main/videoup/conf"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

const (
	_add_black    = "/x/internal/relation/black/add"
	_get_relation = "/x/internal/relation"
)

// Dao is redis dao.
type Dao struct {
	c              *conf.Config
	httpW          *bm.Client
	addBlackURL    string
	getRelationURL string
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:              c,
		httpW:          bm.NewClient(c.HTTPClient.Write),
		addBlackURL:    c.Host.APICO + _add_black,
		getRelationURL: c.Host.APICO + _get_relation,
	}
	return d
}

// Ping ping cpdb
func (d *Dao) Ping(c context.Context) (err error) {
	return
}

//拉黑 http://info.bilibili.co/pages/viewpage.action?pageId=1742202#id-%E5%85%B3%E7%B3%BB%E6%9C%8D%E5%8A%A1%E5%86%85%E7%BD%91%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3-%E8%8E%B7%E5%8F%96%E9%BB%91%E5%90%8D%E5%8D%95%E5%88%97%E8%A1%A8
//网关层 两者关系 http://info.bilibili.co/pages/viewpage.action?pageId=1742202#id-%E5%85%B3%E7%B3%BB%E6%9C%8D%E5%8A%A1%E5%86%85%E7%BD%91%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3-%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7%E4%B8%8E%E5%85%B6%E4%BB%96%E7%94%A8%E6%88%B7%E5%85%B3%E7%B3%BB

// Bind aid,sid,cid bind in one
func (d *Dao) AddBalck(c context.Context, mid, fid, aid int64) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	params.Set("src", strconv.Itoa(221))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpW.Post(c, d.addBlackURL, "", params, &res); err != nil {
		log.Error("d.httpW.Post(%s) error(%v)", d.addBlackURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) code(%d)", d.addBlackURL+"?"+params.Encode(), res.Code)
	}
	log.Info("aid(%d) AddBalck mid(%d) url(%s) code(%d)", aid, mid, d.addBlackURL+"?"+params.Encode(), res.Code)
	return
}

// Relation aid,sid,cid bind in one
func (d *Dao) Relation(c context.Context, mid, fid, aid int64) (attribute int64, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("fid", strconv.FormatInt(fid, 10))
	var res struct {
		Code int `json:"code"`
		Data struct {
			MID int64 `json:"mid"`
			//大于等于128表示拉黑
			Attribute int64 `json:"attribute"`
			Mtime     int64 `json:"mtime"`
		}
	}
	if err = d.httpW.Get(c, d.getRelationURL, "", params, &res); err != nil {
		log.Error("d.httpW.Get(%s) error(%v)", d.getRelationURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) code(%d)", d.getRelationURL+"?"+params.Encode(), res.Code)
	}
	attribute = res.Data.Attribute
	log.Info("aid(%d) Relation mid(%d) url(%s) code(%d) data(%+v)", aid, mid, d.getRelationURL+"?"+params.Encode(), res.Code, res)
	return
}
