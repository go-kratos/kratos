package archive

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-common/library/ecode"
	"go-common/library/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	_nsMd5     = "/videoup/ns/md5"
	_notifyURL = "https://wax.gtloadbalance.cn:38080/InterfaceInfo/GetResult"
)

// AddNetSafeMd5 fn
func (d *Dao) AddNetSafeMd5(c context.Context, nid int64, md5 string) (err error) {
	params := url.Values{}
	params.Set("nid", strconv.FormatInt(nid, 10))
	params.Set("md5", md5)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.nsMd5, "", params, &res); err != nil {
		log.Error("d.client.Post(%s,%s) err(%v)", d.nsMd5, params.Encode(), err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("d.client.res.Code (%s,%s) err(%v)", d.nsMd5, params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// NotifyNetSafe fn
func (d *Dao) NotifyNetSafe(c context.Context, nid int64) (err error) {
	req, err := http.NewRequest("POST", _notifyURL, strings.NewReader(fmt.Sprintf(`nid=%d&companyId=%d`, nid, 2)))
	if err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s)", err, _notifyURL)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")
	httpClient := &http.Client{
		Timeout: time.Duration(time.Duration(5) * time.Second),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				Deadline:  time.Now().Add(5 * time.Second),
				KeepAlive: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
	var resp *http.Response
	if resp, err = httpClient.Do(req); err != nil {
		log.Error("NotifyNetSafe d.client.Do error(%v) | uri(%s)", err, _notifyURL)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Warn("NotifyNetSafe ret body (%+v), nid(%+v)", string(body), nid)
	return
}
