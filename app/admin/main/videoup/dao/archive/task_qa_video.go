package archive

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

//SendQAVideoAdd http request to add qa video task
func (d *Dao) SendQAVideoAdd(c context.Context, task []byte) (err error) {
	ctx, cancel := context.WithTimeout(c, time.Millisecond*500)
	defer cancel()

	res := new(struct {
		Code int   `json:"code"`
		Data int64 `json:"data"`
	})

	val := url.Values{}
	val.Set("appkey", d.c.HTTPClient.Write.App.Key)
	val.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	val.Set("sign", sign(val, d.c.HTTPClient.Write.App.Key, d.c.HTTPClient.Write.App.Secret, true))
	host := fmt.Sprintf("%s?%s", d.addQAVideoURL, val.Encode())
	req, err := http.NewRequest(http.MethodPost, host, bytes.NewBuffer(task))
	if err != nil {
		log.Error("SendQAVideoAdd http.NewRequest error(%v), params(%s)", err, string(task))
		return
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	if err = d.clientW.Do(c, req, res); err != nil {
		log.Error("SendQAVideoAdd d.clientW.Do error(%v)", err)
		return
	}

	if res == nil || res.Code != 0 {
		log.Error("SendQAVideoAdd request failed, response(%+v)", res)
		return
	}

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
