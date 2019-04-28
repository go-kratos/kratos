package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// AllCallbacks return all callbacks in database
func (d *Dao) AllCallbacks(c context.Context) (cbs map[int32]*model.Callback, err error) {
	cbs = make(map[int32]*model.Callback)

	cblist := make([]*model.Callback, 0)
	err = d.ReadORM.Table("workflow_callback").Find(&cblist).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	for _, cb := range cblist {
		cbs[cb.CbID] = cb
	}
	return
}

// SendCallback send callback to pre configured server
func (d *Dao) SendCallback(c context.Context, cb *model.Callback, payload *model.Payload) (err error) {
	var (
		req   *http.Request
		pdata []byte
	)

	if pdata, err = json.Marshal(payload); err != nil {
		return
	}
	// TODO:(zhoujiahui): with sign?
	uv := url.Values{}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	uv.Set("ts", ts)
	uv.Set("appkey", d.writeConf.Key)
	sign := sign(uv, d.writeConf.Key, d.writeConf.Secret, true)

	if req, err = http.NewRequest(http.MethodPost, cb.URL+"?ts="+ts+"&appkey="+d.writeConf.Key+"&sign="+sign, bytes.NewReader(pdata)); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	res := &model.CommonResponse{}
	if err = d.httpWrite.Do(c, req, &res); err != nil {
		log.Error("d.httpWrite.Do(%+v) error(%v)", req, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("callback occur code error url(%s) body(%s) error code(%v)", req.URL, string(pdata), ecode.Int(res.Code))
		return
	}
	log.Info("send callback ok, req(%+v) body(%s) callback(%+v) ", req, string(pdata), cb)
	return
}

// sign is used to sign form params by given condition.
func sign(params url.Values, appkey string, secret string, lower bool) (hexdigest string) {
	data := params.Encode()
	if strings.IndexByte(data, '+') > -1 {
		data = strings.Replace(data, "+", "%20", -1)
	}
	if lower {
		data = strings.ToLower(data)
	}
	digest := md5.Sum([]byte(data + secret))
	hexdigest = hex.EncodeToString(digest[:])
	return
}
