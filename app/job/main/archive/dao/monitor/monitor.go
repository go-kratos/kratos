package monitor

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

	"go-common/library/log"
)

type wxParams struct {
	Username  string `json:"username"`
	Content   string `json:"content"`
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"signature"`
}
type resp struct {
	Status int64  `json:"status"`
	Msg    string `json:"msg"`
}

// Send send message to phone
func (d *Dao) Send(c context.Context, users, msg, token, secret string) (err error) {
	params := url.Values{}
	params.Set("username", users)
	params.Set("content", msg)
	params.Set("token", token)
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	mh := md5.Sum([]byte(params.Encode() + secret))
	params.Set("signature", hex.EncodeToString(mh[:]))
	p := &wxParams{
		Username: params.Get("username"),
		Content:  params.Get("content"),
		Token:    params.Get("token"),
		Sign:     params.Get("signature"),
	}
	p.Timestamp, _ = strconv.ParseInt(params.Get("timestamp"), 10, 64)
	bs, _ := json.Marshal(p)
	payload := strings.NewReader(string(bs))
	req, _ := http.NewRequest("POST", "http://bap.bilibili.co/api/v1/message/add", payload)
	req.Header.Add("content-type", "application/json; charset=utf-8")
	v := &resp{}
	if err = d.client.Do(context.TODO(), req, v); err != nil {
		log.Error("s.client.Do error(%v)", err)
		return
	}
	return
}
