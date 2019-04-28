package weeklyhonor

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_title          = "叮！你有一份荣誉周报待查收！"
	_newYearTitle   = "嘭~~~啪！你有一份荣誉周报待查收！"
	_content        = "这周的成就已达成，快进来瞅瞅吧～ #{周报传送门}{\"https://member.bilibili.com/studio/annyroal/upper-honor-weekly/my\"}"
	_newYearContent = "这周的成就已达成，元气满满过新年～ #{周报传送门}{\"https://member.bilibili.com/studio/annyroal/upper-honor-weekly/my\"}"
	_notifyDataType = "4" // 消息类型：1、回复我的2、@我3、收到的爱4、业务通知
	_notifyURL      = "/api/notify/send.user.notify.do"
	_notifyMC       = "1_17_1"
	_upMidsURL      = "/x/internal/uper/list_up"
)

type newYear struct {
	start, end time.Time
}

var (
	ny = newYear{
		start: time.Date(2019, 2, 3, 0, 0, 0, 0, time.Local),
		end:   time.Date(2019, 2, 24, 0, 0, 0, 0, time.Local),
	}
)

type msgReply struct {
	Code int   `json:"code"`
	Data *data `json:"data"`
}

type data struct {
	TotalCount   int     `json:"total_count"`
	ErrorCount   int     `json:"error_count"`
	ErrorMidList []int64 `json:"error_mid_list"`
}

// SendNotify 发送站内信
func (d *Dao) SendNotify(c context.Context, mids []int64) (errMids []int64, err error) {
	title, content := getTitleAndContent()
	res := msgReply{}
	params := url.Values{}
	params.Set("mc", _notifyMC)
	params.Set("title", title)
	params.Set("data_type", _notifyDataType)
	params.Set("context", content)
	params.Set("mid_list", xstr.JoinInts(mids))
	notifyURI := d.c.Host.Message + _notifyURL
	if err = d.httpClient.Post(c, notifyURI, "", params, &res); err != nil {
		log.Error("d.httpClient.Post(%s,%v,%d)", notifyURI, params, err)
		return
	}
	if res.Code != 0 {
		err = errors.New("code != 0")
		log.Error("d.httpClient.Post(%s,%v,%v,%d)", notifyURI, params, err, res.Code)
	}
	if res.Data != nil {
		errMids = res.Data.ErrorMidList
		log.Info("SendNotify log total_count(%d) error_count(%d) error_mid_list(%v)", res.Data.TotalCount, res.Data.ErrorCount, res.Data.ErrorMidList)
	}
	return
}

func getTitleAndContent() (string, string) {
	now := time.Now()
	if now.After(ny.end) || now.Before(ny.start) {
		return _title, _content
	}
	return _newYearTitle, _newYearContent
}

// Deprecated: use UpActivesList instead.
func (d *Dao) UpMids(c context.Context, size int, lastid int64, activeOnly bool) (mids []int64, newid int64, err error) {
	res := &struct {
		Code int `json:"code"`
		Data *struct {
			Result []int64 `json:"result"`
			LastID int64   `json:"last_id"`
		} `json:"data"`
	}{}
	params := url.Values{}
	params.Set("size", fmt.Sprintf("%d", size))
	params.Set("last_id", fmt.Sprintf("%d", lastid))
	if activeOnly {
		// filter by having archive within 180 days
		params.Set("activity", "1,2,3")
	}
	midsURI := d.c.Host.API + _upMidsURL
	if err = d.httpClient.Get(c, midsURI, "", params, &res); err != nil {
		log.Error("d.httpClient.Get(%s,%v,%d)", midsURI+"?"+params.Encode(), params, err)
		return
	}
	if res.Code != 0 {
		err = errors.New("code != 0")
		log.Error("d.httpClient.Get(%s,%v,%v,%d)", midsURI, params, err, res.Code)
	}
	if res != nil && res.Data != nil {
		return res.Data.Result, res.Data.LastID, nil
	}
	return
}
