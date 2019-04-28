package service

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"go-common/library/log"
)

// HTTPClient http client handle
func (s *Service) HTTPClient(method, url string, params map[string]string, nowTime int64) (body []byte, err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Error("http.NewRequest error(%v)", err)
		return
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	q.Add("appkey", s.conf.KeySecret.Key)
	q.Add("ts", strconv.FormatInt(nowTime, 10))

	sign := q.Encode() + s.conf.KeySecret.Secret
	q.Add("sign", fmt.Sprintf("%x", md5.Sum([]byte(sign))))

	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("http.DefaultClient.Do error(%v)", err)
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("ioutil.ReadAll error(%v)", err)
	}

	return
}
