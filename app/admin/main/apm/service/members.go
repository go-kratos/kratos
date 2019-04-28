package service

import (
	"context"
	"net/http"

	"go-common/app/admin/main/apm/model/monitor"
	"go-common/library/log"
)

const chatURL = "https://api.bilibili.com/x/web-interface/online"

// ChatResult .
type ChatResult struct {
	Code int `json:"code"`
	Data struct {
		WebOnline int64 `json:"web_online"`
	} `json:"data"`
}

// ChatProxy chat proxy
func (s *Service) ChatProxy(c context.Context) (ret *ChatResult, err error) {
	var (
		req *http.Request
	)
	if req, err = s.client.NewRequest(http.MethodGet, chatURL, "", nil); err != nil {
		log.Error("s.ChatProxy.client.NewRequest err(%v)", err)
		return
	}
	if err = s.client.Do(c, req, &ret); err != nil {
		log.Error("s.ChatProxy.client.DO err(%v)", err)
		return
	}
	if ret.Code != 0 {
		log.Error("s.ChatProxy.client http_status(%d)", ret.Code)
		return
	}
	return
}

// Members get current online members
func (s *Service) Members(c context.Context) (mt *monitor.Monitor, err error) {
	var (
		ret = &ChatResult{}
	)
	mt = &monitor.Monitor{}
	if ret, err = s.ChatProxy(c); err != nil {
		return nil, err
	}
	mt.AppID = "online"
	mt.Interface = "online-count"
	mt = &monitor.Monitor{
		AppID:     "online",
		Interface: "online-count",
		Count:     ret.Data.WebOnline,
	}
	return
}
