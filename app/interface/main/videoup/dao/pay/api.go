package pay

import (
	"context"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"net/url"
	"strconv"
)

const (
	_assRegURI = "/x/internal/ugcpay/asset/register"
	_assURI    = "/x/internal/ugcpay/asset"
)

// AssReg 注册付费内容
func (d *Dao) AssReg(c context.Context, mid, aid int64, bp int, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("price", strconv.Itoa(bp*100))
	params.Set("otype", "archive")
	params.Set("platform", "web")
	params.Set("currency", "bp")
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.assRegURI, ip, params, &res); err != nil {
		log.Error("d.client.Do uri(%s) aid(%d) mid(%d) bp(%d) code(%d) error(%v)", d.assRegURI+"?"+params.Encode(), mid, aid, bp, res.Code, err)
		err = ecode.VideoupPayAPIErr
		return
	}
	log.Info("UgcPay AssReg url(%s)", d.assRegURI+"?"+params.Encode())
	if res.Code != 0 {
		log.Error("UgcPay asset register url(%s) res(%+v); mid(%d),aid(%d),bp(%d),ip(%s),code(%d),error(%v)", d.assRegURI, res, mid, aid, bp, ip, res.Code, err)
		err = ecode.VideoupPayAPIErr
		return
	}
	return
}

// Ass 查看付费信息
// UGCPayAssetInvalid         = New(88001) // ugcpay 内容无效
func (d *Dao) Ass(c context.Context, aid int64, ip string) (assert *archive.PayAsset, registed bool, err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("otype", "archive")
	params.Set("currency", "bp")
	var res struct {
		Code int               `json:"code"`
		Data *archive.PayAsset `json:"data"`
	}
	if err = d.client.Get(c, d.assURI, ip, params, &res); err != nil {
		log.Error("d.client.Do uri(%s) aid(%d) code(%d) error(%v)", d.assURI+"?"+params.Encode(), aid, res.Code, err)
		err = ecode.VideoupPayAPIErr
		return
	}
	log.Info("UgcPay AssView url(%s)", d.assURI+"?"+params.Encode())
	if res.Code != 0 {
		if res.Code == ecode.UGCPayAssetInvalid.Code() {
			log.Warn("UgcPay asset UGCPayAssetInvalid url(%s) res(%+v); aid(%d),ip(%s),error(%v)", d.assURI, res, aid, ip, err)
			return nil, false, nil
		}
		err = ecode.VideoupPayAPIErr
		log.Error("VideoupPayAPIErr AssView url(%s) res(%+v); aid(%d),ip(%s),code(%d),error(%v)", d.assURI, res, aid, ip, res.Code, err)
		return
	}
	if res.Data != nil {
		assert = res.Data
		assert.Price = res.Data.Price / 100
		registed = true
	}
	return
}

// UserAcceptProtocol fn: 判断当前的协议是否已经同意过,前端必须传递当前的投稿协议ID
func (d *Dao) UserAcceptProtocol(c context.Context, protocolID string, mid int64) (accept bool, err error) {
	type Res struct {
		Page *struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Total int `json:"total"`
		} `json:"page"`
		Result []*report.UserActionLog `json:"result"`
	}
	res := &Res{}
	r := d.es.NewRequest("log_user_action")
	r.Index("log_user_action_83_all").Pn(1).Ps(2000).OrderScoreFirst(true)
	r.WhereEq("str_0", protocolID).WhereEq("mid", mid)
	r.Order("ctime", "desc")
	log.Info("UserAcceptProtocol params(%s)", r.Params())
	if err = r.Scan(c, res); err != nil {
		log.Error("UserAcceptProtocol r.Scan params(%s)|error(%v)", r.Params(), err)
		return
	}
	if res.Page.Total == 0 {
		accept = false
	} else {
		accept = true
	}
	return
}
