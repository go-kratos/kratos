package data

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

// OpenReply change subject state to open
func (d *Dao) OpenReply(c context.Context, aid int64, mid int64) (err error) {
	params := url.Values{}
	// guanguan admin id
	params.Set("adid", "399")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("type", "1")
	params.Set("state", "0")
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.clientWriter.Post(c, d.replyChangeURL, "", params, &res); err != nil {
		log.Error("OpenReply url(%s) error(%v)", d.replyChangeURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("OpenReply url(%s) code(%d) msg(%s)", d.replyChangeURL+"?"+params.Encode(), res.Code, res.Message)
		if res.Code == ecode.ReplySubjectExist.Code() || res.Code == ecode.ReplySubjectFrozen.Code() || res.Code == ecode.ReplyIllegalSubState.Code() || res.Code == ecode.ReplyIllegalSubType.Code() {
			return
		}
		err = fmt.Errorf("OpenReply call failed")
	}
	return
}

// CloseReply change subject state to close
func (d *Dao) CloseReply(c context.Context, aid int64, mid int64) (err error) {
	params := url.Values{}
	// guanguan admin id
	params.Set("adid", "399")
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("type", "1")
	params.Set("state", "1")
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.clientWriter.Post(c, d.replyChangeURL, "", params, &res); err != nil {
		log.Error("CloseReply url(%s) error(%v)", d.replyChangeURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("CloseReply url(%s) code(%d) msg(%s)", d.replyChangeURL+"?"+params.Encode(), res.Code, res.Message)
		if res.Code == ecode.ReplySubjectExist.Code() || res.Code == ecode.ReplySubjectFrozen.Code() || res.Code == ecode.ReplyIllegalSubState.Code() || res.Code == ecode.ReplyIllegalSubType.Code() {
			return
		}
		err = fmt.Errorf("CloseReply call failed")
	}
	return
}

// CheckReply get subject state
func (d *Dao) CheckReply(c context.Context, aid int64) (replyState int64, err error) {
	params := url.Values{}
	// guanguan admin id
	params.Set("adid", "399")
	params.Set("oid", strconv.FormatInt(aid, 10))
	params.Set("type", "1")
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Oid   int64 `json:"oid"`
			Mid   int64 `json:"mid"`
			State int8  `json:"state"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err = d.client.Get(c, d.replyInfoURL, "", params, &res); err != nil {
		replyState = archive.ReplyDefault
		log.Error("CheckReply url(%s) error(%v)", d.replyInfoURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		replyState = archive.ReplyDefault
		log.Info("CheckReply url(%s) code(%d)", d.replyInfoURL+"?"+params.Encode(), res.Code)
		return
	}
	if res.Data == nil {
		replyState = archive.ReplyDefault
		log.Info("CheckReply url(%s) code(%d) data(%v)", d.replyInfoURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	replyState = int64(res.Data.State)
	return
}
