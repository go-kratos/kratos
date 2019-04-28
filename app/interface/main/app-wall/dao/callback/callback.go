package callback

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/model"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_iis = "/iis?clkid=%s"
)

type Dao struct {
	client *httpx.Client
	iisURL string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(conf.Conf.HTTPActive),
		iisURL: c.Host.Dotin + _iis,
	}
	return
}

func (d *Dao) GdtCallback(c context.Context, appID, appType, aderID string, idfa, cb string, now time.Time) (err error) {
	key, ok := model.ChannelGdt[aderID]
	if !ok {
		return
	}
	encrypt := []byte(key.Encrypt)
	signKey := key.Sign
	uri := fmt.Sprintf("http://t.gdt.qq.com/conv/app/%s/conv?", appID)
	// sign v
	queryS := fmt.Sprintf("muid=%s&conv_time=%d&click_id=%s", idfa, now.Unix(), cb)
	page := signKey + "&GET&" + url.QueryEscape(uri+queryS)
	bs := md5.Sum([]byte(page))
	sign := hex.EncodeToString(bs[:])
	queryS = queryS + "&sign=" + sign
	queryBs := []byte(queryS)
	i := 0
	bss := []byte{}
	for _, b := range queryBs {
		bss = append(bss, b^encrypt[i])
		i = i + 1
		i = i % len(encrypt)
	}
	baseS := base64.StdEncoding.EncodeToString(bss)
	baseS = strings.Replace(baseS, "\n", "", -1)
	// finish uri
	furi := uri + "v=" + url.QueryEscape(baseS) + fmt.Sprintf("&conv_type=MOBILEAPP_ACTIVITE&app_type=%s&advertiser_id=%s", appType, aderID)
	var res struct {
		Ret int    `json:"ret"`
		Msg string `json:"msg"`
	}
	for i := 0; i < 5; i++ {
		if err = d.client.Get(c, furi, "", nil, &res); err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	if !ecode.Int(res.Ret).Equal(ecode.OK) {
		err = errors.Wrapf(ecode.Int(res.Ret), furi)
		if res.Ret == -1 {
			log.Error("%+v", err)
			err = nil
		}
		return
	}
	log.Info("callback gdt furi(%s) idfa(%s) cb(%s) success ret(%d) msg(%s)", furi, idfa, cb, res.Ret, res.Msg)
	return
}

func (d *Dao) ShikeCallback(c context.Context, idfa, cb string, now time.Time) (err error) {
	var res struct {
		Success string `json:"success"`
		Message string `json:"message"`
	}
	if err = d.client.Get(c, cb, "", nil, &res); err != nil {
		return
	}
	log.Info("callback shike idfa(%s) cb(%s) success ret(%s) msg(%s)", idfa, cb, res.Success, res.Message)
	return
}

func (d *Dao) DontinCallback(c context.Context, idfa, clickid string) (err error) {
	urlStr := fmt.Sprintf(d.iisURL, clickid)
	if err = d.client.Get(c, urlStr, "", nil, nil); err != nil {
		return
	}
	log.Info("callback dontin idfa(%s) clickid(%s) success", idfa, clickid)
	return
}

func (d *Dao) ToutiaoCallback(c context.Context, cb string, eventType string) (err error) {
	if cb == "" {
		return
	}
	cbURL := strings.TrimSpace(cb + "&event_type=" + eventType)
	var res struct {
		Ret int `json:"ret"`
	}
	if err = d.client.Get(c, cbURL, "", nil, &res); err != nil {
		return
	}
	if !ecode.Int(res.Ret).Equal(ecode.OK) {
		err = errors.Wrap(ecode.Int(res.Ret), cbURL)
		if res.Ret == -1 {
			log.Error("%+v", err)
			err = nil
		}
		return
	}
	log.Info("callback toutiao cb(%s) eventType(%s) success", cb, eventType)
	return
}
