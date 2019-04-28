package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"crypto/md5"
	"encoding/hex"
	"go-common/app/admin/main/appstatic/conf"
	"go-common/app/admin/main/appstatic/model"
	"go-common/library/log"
)

// Sign fn
func Sign(params url.Values) (query string, err error) {
	if len(params) == 0 {
		return
	}
	if params.Get("appkey") == "" {
		err = fmt.Errorf("utils http get must have parameter appkey")
		return
	}
	if params.Get("appsecret") == "" {
		err = fmt.Errorf("utils http get must have parameter appsecret")
		return
	}
	if params.Get("sign") != "" {
		err = fmt.Errorf("utils http get must have not parameter sign")
		return
	}
	// sign
	secret := params.Get("appsecret")
	params.Del("appsecret")
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp + secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	query = params.Encode()
	return
}

// get sign for NAS storage
func getSign(nas *conf.Bfs) (uri string, err error) {
	var (
		params = url.Values{}
		query  string
	)
	params.Set("appkey", nas.Key)
	params.Set("appsecret", nas.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	if query, err = Sign(params); err != nil {
		log.Error("UpNAS getSign Error (%s)-(%v)-(%v)", nas, err)
		return
	}
	uri = nas.Host + "?" + query
	return
}

// UploadNas can upload the file into Nas Storage
func (d *Dao) UploadNas(c context.Context, fileName string, data []byte, nas *conf.Bfs) (location string, err error) {
	var (
		req    *http.Request
		resp   *http.Response
		client = &http.Client{Timeout: time.Duration(nas.Timeout) * time.Millisecond}
		url    string
		res    = model.ResponseNas{}
	)
	// get sign
	if url, err = getSign(nas); err != nil {
		log.Error("UpNAS getSign Error (%s)-(%v)-(%v)", nas, err)
		return
	}
	// prepare the data of the file and init the request
	buf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(buf)
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		log.Error("UpNAS fileWriter Error (%v)-(%v)", nas, err)
		return
	}
	if _, err = io.Copy(fileWriter, bytes.NewReader(data)); err != nil {
		log.Error("UpNAS fileWriter Copy Error (%v)-(%v)", nas, err)
		return
	}
	bodyWriter.Close()
	// request setting
	if req, err = http.NewRequest(_methodNas, url, buf); err != nil {
		log.Error("http.NewRequest() Upload(%v) error(%v)", url, err)
		return
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	resp, err = client.Do(req)
	// response treatment
	if err != nil {
		log.Error("Nas client.Do(%s) error(%v)", url, err)
		return
	}
	defer resp.Body.Close()
	log.Info("NasAPI returns (%v)", resp)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Nas status code error:%v", resp.StatusCode)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(respBody, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(respBody), err)
		return
	}
	log.Info("NasAPI res struct (%v)", res)
	location = res.Data
	// workaround solution for Macross Upload URL issue
	if d.c.Nas.NewURL != "" {
		location = strings.Replace(location, d.c.Nas.OldURL, d.c.Nas.NewURL, -1)
		log.Error("NasURL replace [%s] to [%s]", res.Data, location)
	}
	return
}
