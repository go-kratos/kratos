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
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

// Upload add mids file.
func (s *Service) Upload(c context.Context, dir, path string, data []byte) (err error) {
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			log.Error("os.IsNotExist(%s) error(%v)", dir, err)
			return
		}
		if err = os.MkdirAll(dir, 0777); err != nil {
			log.Error("os.MkdirAll(%s) error(%v)", dir, err)
			return
		}
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("s.Upload(%s) error(%v)", path, err)
		return
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		log.Error("f.Write() error(%v)", err)
	}
	return
}

// Upimg upload image
func (s *Service) Upimg(ctx context.Context, file multipart.File, header *multipart.FileHeader) (url string, err error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	dest, err := w.CreateFormFile("file", header.Filename)
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
	u := fmt.Sprintf("%s?ts=%s&appkey=%s&sign=%s", s.c.Cfg.UpimgURL, query["ts"], query["appkey"], query["sign"])
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
			URL string `json:"url"`
		} `json:"data"`
	}{}
	if err = s.httpClient.Do(context.TODO(), req, &res); err != nil {
		log.Error("httpClient.Do() error(%v)", err)
		return
	}
	if !ecode.Int(res.Code).Equal(ecode.OK) {
		log.Error("upimg error code(%d)", res.Code)
		err = ecode.Int(res.Code)
		return
	}
	url = strings.Replace(res.Data.URL, "http://", "https://", 1)
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
