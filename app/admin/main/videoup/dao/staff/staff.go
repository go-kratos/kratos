package staff

import (
	"context"

	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const ()

// Staffs fn
func (d *Dao) Staffs(c context.Context, aid int64) (data []*archive.Staff, err error) {
	params := url.Values{}
	params.Set("aid", strconv.FormatInt(aid, 10))
	var res struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Data    []*archive.Staff `json:"data"`
	}
	if err = d.httpClient.Get(c, d.staffURI, "", params, &res); err != nil {
		log.Error("archive.Staffs url(%s) error(%v)", d.staffURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("Staffs api url(%s) res(%v) code(%d)", d.staffURI, res, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}

// StaffApplyBatchSubmit add .
func (d *Dao) StaffApplyBatchSubmit(c context.Context, ap *archive.StaffBatchParam) (err error) {
	params := url.Values{}
	params.Set("appkey", d.c.App.Key)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + d.c.App.Secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	var (
		uri = d.submitUrl + "?" + params.Encode()
	)
	bs, err := json.Marshal(ap)
	if err != nil {
		log.Error("json.Marshal error(%v) | ap(%v) ", err, ap)
		return
	}
	req, err := http.NewRequest("POST", uri, bytes.NewReader(bs))
	if err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s)", err, uri)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.StaffApplyBatchSubmit error(%v) | uri(%s) ap(%+v)", err, uri, ap)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		err = ecode.Error(ecode.Int(res.Code), res.Message)
		log.Error("d.StaffApplyBatchSubmit nq zero (%v)|(%v)|(%v)|(%v)|uri(%s),ap(%+v)", res.Code, res.Message, res, err, uri, ap)
		return
	}
	log.Info("d.StaffApplyBatchSubmit (%s)|res.Data.Aid aid(%d) res(%+v) ip(%s) ", string(bs), ap.AID, res, ip)
	return
}
