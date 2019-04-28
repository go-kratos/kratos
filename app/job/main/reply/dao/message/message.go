package message

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/reply/conf"
	"go-common/library/ecode"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	// 1 main site related
	// 1_1 reply related
	// 1_2 at related
	// 1_3 report related
	//_codeReplyGet    = "1_1_1"
	_codeReplyDelete = "1_1_2"
	_codeReplyLike   = "1_1_3"
	_codeAt          = "1_2_1"
	_codeReport      = "1_3_1"

	_dataTypeReply  = 1
	_dataTypeAt     = 2
	_dataTypeLike   = 3
	_dataTypeSystem = 4

	_notifyTypeCnt = 2
)

// Dao message dao.
type Dao struct {
	httpCli *xhttp.Client
	apiURL  string
}

// NewMessageDao new a message dao and return.
func NewMessageDao(c *conf.Config) *Dao {
	return &Dao{
		httpCli: xhttp.NewClient(c.HTTPClient),
		apiURL:  c.Host.Message + "/api/notify/send.user.notify.do",
	}
}

// Like send a like message.
func (dao *Dao) Like(c context.Context, mid, tomid int64, title, msg, extraInfo string, now time.Time) (err error) {
	return dao.send(c, _codeReplyLike, "", title, msg, _dataTypeLike, mid, []int64{tomid}, extraInfo, now.Unix())
}

// Reply send a reply message.
func (dao *Dao) Reply(c context.Context, mc, resID string, mid, tomid int64, title, msg, extraInfo string, now time.Time) (err error) {
	return dao.send(c, mc, resID, title, msg, _dataTypeReply, mid, []int64{tomid}, extraInfo, now.Unix())
}

// DeleteReply send delete reply message.
func (dao *Dao) DeleteReply(c context.Context, mid int64, title, msg string, now time.Time) (err error) {
	return dao.send(c, _codeReplyDelete, "", title, msg, _dataTypeSystem, 0, []int64{mid}, "", now.Unix())
}

// At send a at message.
func (dao *Dao) At(c context.Context, mid int64, mids []int64, title, msg, extraInfo string, now time.Time) (err error) {
	if len(mids) == 0 {
		return
	}
	return dao.send(c, _codeAt, "", title, msg, _dataTypeAt, mid, mids, extraInfo, now.Unix())
}

// AcceptReport send accept report message.
func (dao *Dao) AcceptReport(c context.Context, mid int64, title, msg string, now time.Time) (err error) {
	return dao.send(c, _codeReport, "", title, msg, _dataTypeSystem, 0, []int64{mid}, "", now.Unix())
}

// System send a system message.
func (dao *Dao) System(c context.Context, mc, resID string, mid int64, title, msg, info string, now time.Time) (err error) {
	return dao.send(c, mc, resID, title, msg, _dataTypeSystem, 0, []int64{mid}, info, now.Unix())
}

func (dao *Dao) send(c context.Context, mc, resID, title, msg string, tp int, pub int64, mids []int64, info string, ts int64) (err error) {
	params := url.Values{}
	params.Set("type", "json")
	params.Set("source", "1")
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("data_type", strconv.Itoa(tp))
	params.Set("context", msg)
	params.Set("mid_list", xstr.JoinInts(mids))
	params.Set("publisher", strconv.FormatInt(pub, 10))
	params.Set("ext_info", info)
	if resID != "" {
		params.Set("notify_type", fmt.Sprint(_notifyTypeCnt))
		params.Set("res_id", resID)
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = dao.httpCli.Post(c, dao.apiURL, "", params, &res); err != nil {
		log.Error("message url(%s) error(%v)", dao.apiURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("message url(%s) error(%v)", dao.apiURL+"?"+params.Encode(), res.Code)
		err = fmt.Errorf("message send failed")
		return
	}
	log.Info("sendmessage success:%v;code:%d", params, res.Code)

	if tp != _dataTypeSystem {
		params.Set("mobi_app", "android_i")
		if tp == _dataTypeAt {
			params.Set("title", converAt(title))
		} else if tp == _dataTypeLike {
			params.Set("context", convertMsg(msg))
		} else if tp == _dataTypeReply {
			params.Set("title", convertMsg(title))
			params.Set("context", convertMsg(msg))
		}
		var res1 struct {
			Code int `json:"code"`
		}
		if err = dao.httpCli.Post(c, dao.apiURL, "", params, &res1); err != nil {
			log.Error("message url(%s) error(%v)", dao.apiURL+"?"+params.Encode(), err)
			return
		}
		if res1.Code != ecode.OK.Code() {
			log.Error("message url(%s) error(%v)", dao.apiURL+"?"+params.Encode(), res1.Code)
			err = fmt.Errorf("message send failed")
			return
		}
		log.Info("send international message success:%v;code:%d", params, res1.Code)
	}
	return
}

func converAt(title string) string {
	return strings.Replace(title, "评论中@了你", "評論中@了你", -1)
}

func convertMsg(msg string) string {
	rmsg := []rune(msg)
	for i, c := range rmsg {
		switch c {
		case '评':
			rmsg[i] = '評'
		case '论':
			rmsg[i] = '論'
		case '赞':
			rmsg[i] = '讚'
		case '条':
			rmsg[i] = '條'
		case '专':
			rmsg[i] = '專'
		case '栏':
			rmsg[i] = '欄'
		case '数':
			rmsg[i] = '數'
		case '达':
			rmsg[i] = '達'
		case '应':
			rmsg[i] = '應'
		case '点':
			rmsg[i] = '點'
		case '击':
			rmsg[i] = '擊'
		default:
			continue
		}
	}
	return string(rmsg)
}
