package dao

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/json-iterator/go"

	"go-common/app/job/bbq/recall/internal/conf"
	"go-common/app/job/bbq/recall/internal/model"
)

func (d *Dao) queryHDFS(c context.Context, api string, key *conf.BerserkerKey, suffix string) (result *[]byte, err error) {
	dt := time.Now().Format("2006-01-02 15:04:05")
	sign := d.berserkerSign(key.AppKey, key.Secret, dt, "1.0")

	params := &url.Values{}
	params.Set("appKey", key.AppKey)
	params.Set("timestamp", dt)
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	params.Set("sign", sign)

	fileSuffix := struct {
		FileSuffix string `json:"fileSuffix"`
	}{
		FileSuffix: suffix,
	}
	j, err := jsoniter.Marshal(fileSuffix)
	params.Set("query", string(j))

	for retry := 0; retry < 3; retry++ {
		resp, err := http.DefaultClient.Get(api + "?" + params.Encode())
		if err != nil {
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil && len(b) > 0 {
			result = &b
			break
		}
		// sleep 5s berserker limit
		time.Sleep(5 * time.Second)
	}

	return
}

func (d *Dao) scanHDFSPath(c context.Context, api string, key *conf.BerserkerKey, suffix string) (result *model.HDFSResult, err error) {
	b, err := d.queryHDFS(c, api, key, suffix)
	if err != nil {
		return
	}
	result = &model.HDFSResult{}
	err = jsoniter.Unmarshal(*b, result)
	return
}

func (d *Dao) loadHDFSFile(c context.Context, api string, key *conf.BerserkerKey, suffix string) (result *[]byte, err error) {
	return d.queryHDFS(c, api, key, suffix)
}
