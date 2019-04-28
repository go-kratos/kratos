package dao

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go-common/library/conf/env"
	"net/http"
	"net/url"
)

// NewRequst http 请求
func (d *Dao) NewRequst(c context.Context, method string, url string, query url.Values, body []byte, headers map[string]string, resp interface{}) error {
	var req *http.Request
	if body != nil && len(body) > 0 {
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if err := d.httpClient.Do(c, req, &resp); err != nil {
		err = errors.WithStack(err)
		return err
	}

	return nil
}

// getLiveStreamUrl 对接live-stream.bilibili.co的相关业务
func (d *Dao) getLiveStreamUrl(path string) string {
	url := ""
	if env.DeployEnv == env.DeployEnvProd {
		url = fmt.Sprintf("%s%s", "http://prod-live-stream.bilibili.co", path)
	} else {
		url = fmt.Sprintf("%s%s", "http://live-stream.bilibili.co", path)
	}
	return url
}
