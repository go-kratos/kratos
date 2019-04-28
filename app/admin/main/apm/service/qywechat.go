package service

import (
	"context"
	"net/url"

	"go-common/library/log"
)

// Ret .
type Ret struct {
	ReqID    string   `json:"ReqId"`
	Action   string   `json:"Action"`
	RetCode  int      `json:"RetCode"`
	Data     []string `json:"Data"`
	Response struct {
		Status int `json:"status"`
	} `json:"Response"`
}

// SendWeChat send message to WeChat
// users: zhangsan,lisi,wangwu
func (s *Service) SendWeChat(c context.Context, title, msg, treeID, users string) (err error) {
	var (
		params = url.Values{}
		ret    = &Ret{}
	)
	params.Add("Action", "CreateWechatMessage")
	params.Add("PublicKey", s.c.Prometheus.Key)
	params.Add("Signature", "1")
	params.Add("UserName", users)
	params.Add("Title", title)
	params.Add("Content", title+"\n"+msg)
	params.Add("TreeId", "bilibili."+treeID)
	if err = s.PrometheusProxy(context.Background(), params, ret); err != nil {
		log.Error("s.SendWeChat error(%v)", err)
	}
	return
}
