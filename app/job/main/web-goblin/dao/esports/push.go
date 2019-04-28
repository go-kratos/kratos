package esports

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	mdlesp "go-common/app/job/main/web-goblin/model/esports"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_pinkVersion = 1
	_linkLive    = 3
)

type _response struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

// NoticeUser pushs the notification to users.
func (d *Dao) NoticeUser(mids []int64, body string, contest *mdlesp.Contest) (err error) {
	var strMids string
	if d.c.Push.OnlyMids == "" {
		strMids = xstr.JoinInts(mids)
	} else {
		strMids = d.c.Push.OnlyMids
	}
	uuid := d.getUUID(strMids, contest)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField("app_id", strconv.Itoa(_pinkVersion))
	w.WriteField("business_id", strconv.Itoa(d.c.Push.BusinessID))
	w.WriteField("alert_title", d.c.Push.Title)
	w.WriteField("alert_body", body)
	w.WriteField("mids", strMids)
	w.WriteField("link_type", strconv.Itoa(_linkLive))
	w.WriteField("link_value", strconv.FormatInt(contest.LiveRoom, 10))
	w.WriteField("uuid", uuid)
	w.Close()
	//签名
	query := map[string]string{
		"ts":     strconv.FormatInt(time.Now().Unix(), 10),
		"appkey": d.c.App.Key,
	}
	query["sign"] = d.signature(query, d.c.App.Secret)
	url := fmt.Sprintf("%s?ts=%s&appkey=%s&sign=%s", d.pushURL, query["ts"], query["appkey"], query["sign"])

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v)", url, err)
		return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("token=%s", d.c.Push.BusinessToken))
	res := &_response{}
	if err = d.http.Do(context.TODO(), req, &res); err != nil {
		log.Error("httpClient.Do() error(%v)", err)
		return
	}
	if res.Code != 0 || res.Data == 0 {
		log.Error("push failed mids_total(%d) body(%s) response(%+v)", len(mids), body, res)
	} else {
		log.Info("push success mids_total(%d) body(%s) response(%+v)", len(mids), body, res)
	}
	return
}

//signature 加密算法
func (d *Dao) signature(params map[string]string, secret string) string {
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := bytes.Buffer{}
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k) + "=")
		buf.WriteString(url.QueryEscape(params[k]))
	}
	//加密
	h := md5.New()
	io.WriteString(h, buf.String()+secret)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (d *Dao) getUUID(mids string, contest *mdlesp.Contest) string {
	var b bytes.Buffer
	b.WriteString(strconv.Itoa(d.c.Push.BusinessID))
	b.WriteString(strconv.FormatInt(contest.ID, 10))
	b.WriteString(strconv.FormatInt(contest.Stime, 10))
	b.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
	b.WriteString(mids)
	mh := md5.Sum(b.Bytes())
	uuid := hex.EncodeToString(mh[:])
	return uuid
}
