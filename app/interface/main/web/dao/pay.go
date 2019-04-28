package dao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Wallet get user bcoin  doc:http://info.bilibili.co/pages/viewpage.action?pageId=7559096
func (d *Dao) Wallet(c context.Context, mid int64) (wallet *model.Wallet, err error) {
	const (
		plat     = 3
		platStr  = "3"
		customID = "10008"
	)
	params := url.Values{}
	params.Set("customerId", customID)
	params.Set("platformType", platStr)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("traceId", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/1000, 10))
	params.Set("signType", "MD5")
	params.Set("appkey", d.c.HTTPClient.Pay.Key)
	type pJSON struct {
		CustomerID   string `json:"customerId"`
		PlatformType int    `json:"platformType"`
		Mid          int64  `json:"mid"`
		TraceID      string `json:"traceId"`
		Timestamp    string `json:"timestamp"`
		SignType     string `json:"signType"`
		Appkey       string `json:"appkey"`
		Sign         string `json:"sign"`
	}
	tmp := params.Encode() + d.c.HTTPClient.Pay.Secret
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp))
	sign := hex.EncodeToString(mh[:])
	p := &pJSON{
		CustomerID:   customID,
		PlatformType: plat,
		Mid:          mid,
		TraceID:      params.Get("traceId"),
		Timestamp:    params.Get("timestamp"),
		SignType:     params.Get("signType"),
		Appkey:       params.Get("appkey"),
		Sign:         sign,
	}
	bs, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", d.walletURL, strings.NewReader(string(bs)))
	req.Header.Set("Content-Type", "application/json")
	var res struct {
		Code int `json:"code"`
		Data struct {
			// DefaultBp     float32 `json:"defaultBp"`
			CouponBalance float32 `json:"couponBalance"`
			AvailableBp   float32 `json:"availableBp"`
		} `json:"data"`
	}
	if err = d.httpPay.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrap(ecode.Int(res.Code), d.walletURL+"?"+params.Encode())
		log.Error("account pay url(%s) error(%v)", d.walletURL+"?"+params.Encode(), res.Code)
		return
	}
	wallet = &model.Wallet{
		Mid:           mid,
		BcoinBalance:  res.Data.AvailableBp,
		CouponBalance: res.Data.CouponBalance,
	}
	return
}

// OldWallet get wallet info
func (d *Dao) OldWallet(c context.Context, mid int64) (w *model.Wallet, err error) {
	var (
		params   = url.Values{}
		remoteIP = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int           `json:"code"`
		Data *model.Wallet `json:"data"`
	}
	err = d.httpR.Get(c, d.walletOldURL, remoteIP, params, &res)
	if err != nil {
		log.Error("account pay url(%s) error(%v)", d.walletOldURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("account pay url(%s) error(%v)", d.walletOldURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	w = res.Data
	return
}
