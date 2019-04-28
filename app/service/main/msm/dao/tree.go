package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-common/app/service/main/msm/model"
	"go-common/library/conf/env"
	"go-common/library/log"
)

const (
	treeUsername  = "msm"
	treeAuthURL   = "/v1/token"
	allAppAuthURL = "%s/v1/node/app/secretinfo/%s"
)

func (d *Dao) treeToken(c context.Context) (token string, err error) {
	var (
		jsonBytes []byte
		url       = d.treeHost + treeAuthURL
	)
	body := &struct {
		Username   string `json:"user_name"`
		PlatformID string `json:"platform_id"`
	}{
		Username:   treeUsername,
		PlatformID: d.platformID,
	}
	if jsonBytes, err = json.Marshal(body); err != nil {
		log.Error("json.Marshal(body) error(%v)", err)
		return
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Error("http.NewRequest failed", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	res := &struct {
		Code int64 `json:"code"`
		Data struct {
			Token    string `json:"token"`
			Username string `json:"user_name"`
			Secret   string `json:"secret"`
			Expired  int64  `json:"expired"`
		} `json:"data"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{}
	if err = d.client.Do(c, req, res); err != nil {
		log.Error("service-tree client Do failed", err)
		return
	}
	if res.Code != 90000 {
		log.Error("service-tree client Do failed", err)
		return
	}
	token = res.Data.Token
	return
}

// TreeAppInfo TreeAppInfo.
func (d *Dao) TreeAppInfo(c context.Context) (appInfo map[int64]*model.AppToken, err error) {
	var (
		token string
		url   = fmt.Sprintf(allAppAuthURL, d.treeHost, env.DeployEnv)
	)
	appInfo = make(map[int64]*model.AppToken)
	if token, err = d.treeToken(c); err != nil {
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("http.NewRequest failed", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)
	res := &struct {
		Code    int               `json:"code"`
		Data    []*model.AppToken `json:"data"`
		Message string            `json:"message"`
		Status  int               `json:"status"`
	}{}
	if err = d.client.Do(c, req, res); err != nil {
		log.Error("service-tree client Do failed", err)
		return
	}
	if res.Code != 90000 {
		log.Error("service-tree client Do failed", err)
		return
	}
	for _, auth := range res.Data {
		appInfo[auth.AppTreeID] = auth
	}
	return
}
