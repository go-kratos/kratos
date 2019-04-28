package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

// Upimg upload image
func (s *Service) Upimg(ctx context.Context, file multipart.File) (url string, err error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	w.WriteField("bucket", "push")
	dest, err := w.CreateFormFile("file", "push-service-upimg")
	if err != nil {
		log.Error("Upimg CreateFormFile error(%v)", err)
		return
	}
	if _, err = io.Copy(dest, file); err != nil {
		log.Error("Upimg Copy error(%v)", err)
		return
	}
	w.Close()
	query := map[string]string{
		"ts":     strconv.FormatInt(time.Now().Unix(), 10),
		"appkey": s.c.HTTPClient.App.Key,
	}
	query["sign"] = sign(query, s.c.HTTPClient.App.Secret)
	u := fmt.Sprintf("%s?ts=%s&appkey=%s&sign=%s", s.c.Push.UpimgURL, query["ts"], query["appkey"], query["sign"])
	req, err := http.NewRequest(http.MethodPost, u, buf)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v)", u, err)
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res := &struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Etag     string `json:"etag"`
			Location string `json:"location"`
		} `json:"data"`
	}{}
	if err = s.httpclient.Do(context.TODO(), req, &res); err != nil {
		log.Error("httpClient.Do() error(%v)", err)
		return
	}
	if !ecode.Int(res.Code).Equal(ecode.OK) {
		log.Error("upimg error code(%d)", res.Code)
		err = ecode.Int(res.Code)
		return
	}
	url = strings.Replace(res.Data.Location, "http://", "https://", 1)
	return
}

func sign(params map[string]string, secret string) string {
	var keys []string
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
	h := md5.New()
	io.WriteString(h, buf.String()+secret)
	return fmt.Sprintf("%x", h.Sum(nil))
}
