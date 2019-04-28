package dao

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"go-common/app/service/main/spy/model"
	"go-common/library/log"
)

func (d *Dao) hmacsha1(key, text string) (h string) {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(text))
	h = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return
}

func (d *Dao) makeURL(method string, action string, region string, secretID string, secretKey string,
	args url.Values, charset string, URL string) (req string) {
	args.Set("Nonce", fmt.Sprintf("%d", d.r.Uint32()))
	args.Set("Action", action)
	args.Set("Region", region)
	args.Set("SecretId", secretID)
	args.Set("Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	args.Set("Signature", d.hmacsha1(secretKey, fmt.Sprintf("%s%s?%s", method, URL, d.makeQueryString(args))))
	req = args.Encode()
	return
}

func (d *Dao) makeQueryString(v url.Values) (str string) {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteString("&")
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}

// RegisterProtection register protection.
func (d *Dao) RegisterProtection(c context.Context, args url.Values, ip string) (level int8, err error) {
	query := d.makeURL("GET", d.c.Qcloud.Path, d.c.Qcloud.Region, d.c.Qcloud.SecretID,
		d.c.Qcloud.SecretKey, args, d.c.Qcloud.Charset, d.c.Qcloud.BaseURL)
	req, err := http.NewRequest("GET", "https://"+d.c.Qcloud.BaseURL+"?"+query, nil)
	if err != nil {
		log.Error("d.RegisterProtection uri(%s) error(%v)", query, err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := &model.QcloudRegProResp{}
	if err = d.httpClient.Do(c, req, res); err != nil {
		log.Error("d.client.Do error(%v) | uri(%s)) res(%v)", err, query, res)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("GET RegisterProtection req faild query(%s) resp(%v)", d.c.Qcloud.BaseURL+"?"+query, res)
		log.Error(" RegisterProtection fail res(%v)", res)
		return
	}
	level = res.Level
	log.Info("GET RegisterProtection suc query(%s) resp(%v)", query, res)
	return
}
