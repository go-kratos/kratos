package dao

import (
	"context"
	"net/http"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

const (
	_wechatGroup  = "/ep/admin/saga/v2/wechat/appchat/send"
	_wechatPerson = "/ep/admin/saga/v2/wechat/message/send"
)

//AddWechatSend send msg to group
func (d *Dao) AddWechatSend(c context.Context, cookie, content string) (msgSendRes *model.MsgSendRes, err error) {
	var (
		url        = conf.Conf.Wechat.Host + _wechatGroup
		req        *http.Request
		msgSendReq = &model.MsgSendReq{
			ChatID:  conf.Conf.Wechat.Chatid,
			MsgType: conf.Conf.Wechat.Msgtype,
			Text:    model.MsgSendReqText{Content: content},
			Safe:    conf.Conf.Wechat.Safe,
		}
	)
	if req, err = d.newRequest(http.MethodPost, url, msgSendReq); err != nil {
		return
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")

	if err = d.httpClient.Do(c, req, &msgSendRes); err != nil {
		log.Error("d.AddWechatSend url(%s) res($s) error(%v)", url, msgSendRes, err)
		return
	}
	return
}

// PushWechatMsgToPerson send msg to users
func (d *Dao) PushWechatMsgToPerson(c context.Context, cookie string, users []string, msg string) (msgSendRes *model.MsgSendRes, err error) {
	var (
		url        = conf.Conf.Wechat.Host + _wechatPerson
		req        *http.Request
		msgSendReq = &model.MsgSendPersonReq{
			Users:   users,
			Content: msg,
		}
	)

	if req, err = d.newRequest(http.MethodPost, url, msgSendReq); err != nil {
		return
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")

	if err = d.httpClient.Do(c, req, &msgSendRes); err != nil {
		log.Error("d.WeChatPerson url(%s) res($s) error(%v)", url, msgSendRes, err)
		return
	}
	return
}
