package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

// RegReply del reply.
func (d *Dao) RegReply(c context.Context, id int64, tid int8) (err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(id, 10))
	params.Set("mid", "0")
	params.Set("type", strconv.FormatInt(int64(tid), 10))
	var res struct {
		Code int `json:"code"`
	}
	for i := 0; i <= 10; i++ {
		if err = d.client.Post(c, d.regReplyURL, "", params, &res); err != nil {
			log.Error("d.regReplyURL url(%s) res(%v) error(%v)", d.regReplyURL+"?"+params.Encode(), res, err)
			continue
		}
		if res.Code != ecode.OK.Code() && res.Code != ecode.ReplySubjectExist.Code() {
			log.Error("d.regReplyURL code(%v) url(%s)", res.Code, d.regReplyURL+"?"+params.Encode())
			continue
		}
		log.Info("d.regReplyURL url(%s) res(%v)", d.regReplyURL+"?"+params.Encode(), res)
		return
	}
	return
}

// DelReply del reply.
func (d *Dao) DelReply(c context.Context, rpid, tp, oid string) (err error) {
	params := url.Values{}
	params.Set("oid", oid)
	params.Set("rpid", rpid)
	params.Set("type", tp)
	params.Set("adid", "-1")
	params.Set("notify", "false")
	params.Set("moral", "0")
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.delReplyURL, "", params, &res); err != nil {
		log.Error("d.delReplyURL url(%s) res(%v) error(%v)", d.delReplyURL+"?"+params.Encode(), res, err)
		return
	}
	if res.Code != 0 {
		log.Error("d.delReplyURL url(%s) code(%d)", d.delReplyURL+"?"+params.Encode(), res.Code)
	}
	log.Info("d.delReplyURL url(%s) res(%v)", d.delReplyURL+"?"+params.Encode(), res)
	return
}

// UpReplyState update reply state.
func (d *Dao) UpReplyState(c context.Context, oid, rpid int64, tp, state int8) (err error) {
	params := url.Values{}
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("rpid", strconv.FormatInt(rpid, 10))
	params.Set("type", strconv.FormatInt(int64(tp), 10))
	params.Set("state", strconv.FormatInt(int64(state), 10))
	params.Set("adid", "0")
	var res struct {
		Code int `json:"code"`
	}
	for i := 0; i <= 10; i++ {
		if err = d.client.Post(c, d.upReplyStateURL, "", params, &res); err != nil {
			log.Error("d.upReplyStateURL url(%s) res(%v) error(%v)", d.upReplyStateURL+"?"+params.Encode(), res, err)
			continue
		}
		if res.Code != 0 {
			log.Error("d.upReplyStateURL url(%s) code(%d)", d.upReplyStateURL+"?"+params.Encode(), res.Code)
			continue
		}
		log.Info("d.upReplyStateURL url(%s) res(%v)", d.upReplyStateURL+"?"+params.Encode(), res)
		return
	}
	return
}
