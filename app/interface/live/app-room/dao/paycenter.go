package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"go-common/app/interface/live/app-room/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	customID = "10005"
)

//PayCenterWallet Wallet get user bcoin  doc:http://info.bilibili.co/pages/viewpage.action?pageId=7559096
func (d *Dao) PayCenterWallet(c context.Context, mid int64, platform string) (wallet *model.Wallet, err error) {
	var plat int
	if platform == "ios" {
		plat = 1
	} else if platform == "android" {
		plat = 2
	} else if platform == "pc" {
		plat = 3
	} else {
		err = ecode.ParamInvalid
		return
	}
	platStr := strconv.Itoa(plat)
	params := url.Values{}
	params.Set("customerId", customID)
	params.Set("platformType", platStr)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("traceId", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	params.Set("signType", "MD5")
	params.Set("token", d.c.HTTPClient.PayCenter.Secret)
	type pJSON struct {
		CustomerID   string `json:"customerId"`
		PlatformType int    `json:"platformType"`
		Mid          int64  `json:"mid"`
		TraceID      string `json:"traceId"`
		Timestamp    string `json:"timestamp"`
		SignType     string `json:"signType"`
		Sign         string `json:"sign"`
	}
	tmp := encode(params)
	mh := md5.Sum([]byte(tmp))
	sign := hex.EncodeToString(mh[:])
	p := &pJSON{
		CustomerID:   customID,
		PlatformType: plat,
		Mid:          mid,
		TraceID:      params.Get("traceId"),
		Timestamp:    params.Get("timestamp"),
		SignType:     params.Get("signType"),
		Sign:         sign,
	}
	bs, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", d.payCenterWalletURL, strings.NewReader(string(bs)))
	req.Header.Set("Content-Type", "application/json")
	var res struct {
		Code int    `json:"errno"`
		Msg  string `json:"msg"`
		Data struct {
			// DefaultBp     float32 `json:"defaultBp"`
			CouponBalance float32 `json:"couponBalance"`
			AvailableBp   float32 `json:"availableBp"`
			DefaultBp     float32 `json:"defaultBp"`
			IosBp         float32 `json:"iosBp"`
		} `json:"data"`
	}
	if err = d.payCenterClient.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.Wrap(ecode.Int(res.Code), d.payCenterWalletURL+"?"+string(bs))
		log.Error("account pay url(%s) error(%v) %s", d.payCenterWalletURL+"?"+string(bs), res.Code, res.Msg)
		return
	}
	wallet = &model.Wallet{
		Mid:           mid,
		BcoinBalance:  res.Data.AvailableBp,
		CouponBalance: res.Data.CouponBalance,
		DefaultBp:     res.Data.DefaultBp,
		IosBp:         res.Data.IosBp,
	}
	return
}

func encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	ht := false
	for k := range v {
		if k == "token" {
			ht = true
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if ht { // 需要把token放在最后 支付中心的规则
		keys = append(keys, "token")
	}
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
