package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	model "go-common/app/interface/main/credit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// SendSysMsg send msg.
func (dao *Dao) SendSysMsg(c context.Context, mid int64, title string, context string) (err error) {
	params := url.Values{}
	params.Set("mc", "2_1_13")
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", context)
	params.Set("mid_list", fmt.Sprintf("%d", mid))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Status int8   `json:"status"`
			Remark string `json:"remark"`
		} `json:"data"`
	}
	if err = dao.client.Post(c, dao.sendMsgURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("SendSysMsg(%s) error(%v)", dao.sendMsgURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("SendSysMsg(%s) error(%v)", dao.sendMsgURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	return
}

// GetQS get question from big data.
func (dao *Dao) GetQS(c context.Context, mid int64) (qs *model.AIQsID, err error) {
	qs = &model.AIQsID{}
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int           `json:"code"`
		Data *model.AIQsID `json:"data"`
	}
	if err = dao.client.Get(c, dao.getQSURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("GetQS(%s) error(%v)", dao.getQSURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("GetQS(%s) error(%v)", dao.getQSURL+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	qs = res.Data
	return
}

// ReplysCount  get reply count.
func (dao *Dao) ReplysCount(c context.Context, oid []int64) (counts map[string]int64, err error) {
	params := url.Values{}
	params.Set("oid", xstr.JoinInts(oid))
	params.Set("type", "6")
	var res struct {
		Data map[string]int64 `json:"data"`
	}
	if err = dao.client.Get(c, dao.replyCountURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("ReplysCount(%s) err(%v)", dao.replyCountURL+"?"+params.Encode(), err)
		return
	}
	counts = res.Data
	return
}

// SendMedal send mdal.
func (dao *Dao) SendMedal(c context.Context, mid int64, nid int64) (err error) {
	params := url.Values{}
	params.Set("nid", strconv.FormatInt(nid, 10))
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = dao.client.Post(c, dao.sendMedalURL, "", params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = errors.WithStack(ecode.Int(res.Code))
	}
	return
}
