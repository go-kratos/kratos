package pay

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/model/pay"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

type Dao struct {
	c      *conf.Config
	client *httpx.Client
	wallet string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		client: httpx.NewClient(c.HTTPClient),
		wallet: c.Host.Pay + "/wallet-int/wallet/getUserWalletInfo",
	}
	return
}

// UserWalletInfo get user bcoin  doc:http://info.bilibili.co/pages/viewpage.action?pageId=7559096
func (d *Dao) UserWalletInfo(c context.Context, mid int64, platform string) (availableBp float64, err error) {
	var plat int
	if platform == "ios" {
		plat = 1
	} else if platform == "android" {
		plat = 2
	} else {
		err = fmt.Errorf("platform(%s) error", platform)
		return
	}
	params := url.Values{}
	params.Set("customerId", "10006")
	params.Set("platformType", strconv.Itoa(plat))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("traceId", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/1000, 10))
	params.Set("signType", "MD5")
	params.Set("appkey", d.c.HTTPClient.Key)
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
	tmp := params.Encode() + d.c.HTTPClient.Secret
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp))
	sign := hex.EncodeToString(mh[:])
	p := &pJSON{
		CustomerID:   "10006",
		PlatformType: plat,
		Mid:          mid,
		TraceID:      params.Get("traceId"),
		Timestamp:    params.Get("timestamp"),
		SignType:     params.Get("signType"),
		Appkey:       params.Get("appkey"),
		Sign:         sign,
	}
	bs, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", d.wallet, strings.NewReader(string(bs)))
	req.Header.Set("Content-Type", "application/json")
	var wallet *pay.UserWallet
	if err = d.client.Do(c, req, &wallet); err != nil {
		return
	}
	if wallet.Code != 0 {
		err = errors.Wrap(ecode.Int(wallet.Code), d.wallet+"?"+params.Encode())
		return
	}
	availableBp = wallet.Data.AvailableBp
	return
}
