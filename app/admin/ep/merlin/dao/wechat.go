package dao

import (
	"context"
	"encoding/json"
	"net/http"

	"go-common/app/admin/ep/merlin/conf"
	"go-common/app/admin/ep/merlin/model"
	"go-common/library/log"
)

const (
	_wechatGroup = "/ep/admin/saga/v2/wechat/appchat/send"
)

//WeChatSendMessage We Chat Send Message
func (d *Dao) WeChatSendMessage(c context.Context, msgSendReq *model.MsgSendReq) (msgSendRes *model.MsgSendRes, err error) {
	var (
		url = conf.Conf.WeChat.WeChatHost + _wechatGroup
		req *http.Request
		res = &model.MsgSendRes{}
	)
	msgSendRequest, _ := json.Marshal(msgSendReq)
	log.Info("url:(%s)", url)
	log.Info("msgSendRequest:(%s)", string(msgSendRequest))

	if req, err = d.newRequest(http.MethodPost, url, msgSendReq); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.AddWechatSend url(%s) res($s) error(%v)", url, res, err)
		return
	}
	msgSendRes = res
	rsp, _ := json.Marshal(msgSendRes)
	log.Info("wechat send message response :(%s)", string(rsp))
	return
}
