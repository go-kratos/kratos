package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"go-common/app/job/main/passport-sns/model"
	"go-common/library/log"
)

const (
	_qqGetUnionidUrl = "https://graph.qq.com/oauth2.0/get_unionid"

	_respCodeSuccess             = 0
	_respCodeInvalidOpenid       = 100050
	_respCodeInvalidOpenidLength = 100058
)

// QQUnionID .
func (d *Dao) QQUnionID(c context.Context, openID string) (unionID string, err error) {
	var (
		res    *http.Response
		bs     []byte
		params = url.Values{}
	)
	params.Set("client_id", model.OldAppID)
	params.Set("openid", openID)

	if res, err = http.Get(_qqGetUnionidUrl + "?" + params.Encode()); err != nil {
		log.Error("d.QQUnionID error(%+v) openID(%s)", err, openID)
		return "", err
	}
	defer res.Body.Close()

	if bs, err = ioutil.ReadAll(res.Body); err != nil {
		log.Error("ioutil.ReadAll() error(%+v) openID(%s)", err, openID)
		return "", err
	}
	respStr := string(bs)
	if strings.HasPrefix(respStr, "callback") {
		start := strings.Index(respStr, "{")
		end := strings.Index(respStr, "}")
		respStr = respStr[start : end+1]
	}
	unionRes := new(model.QQUnionIDResp)
	if err = json.Unmarshal([]byte(respStr), unionRes); err != nil {
		log.Error("json.Unmarshal() error(%+v) respStr(%s) openID(%s)", err, respStr, openID)
		return "", err
	}
	if unionRes.Code == _respCodeSuccess && unionRes.UnionID != "" {
		return unionRes.UnionID, nil
	}
	log.Error("request qq get_unionid failed with code(%d) desc(%s) openID(%s)", unionRes.Code, unionRes.Description, openID)
	if unionRes.Code == _respCodeInvalidOpenidLength || unionRes.Code == _respCodeInvalidOpenid {
		return "", nil
	}
	return "", fmt.Errorf("request qq get_unionid failed with code(%d) desc(%s) openID(%s)", unionRes.Code, unionRes.Description, openID)
}
