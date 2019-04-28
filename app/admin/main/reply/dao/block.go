package dao

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_blockURL    string = "http://account.bilibili.co/api/member/blockAccountWithTime"
	_transferURL string = "http://api.bilibili.co/x/internal/credit/blocked/case/add"
)

// BlockAccount ban an account.
func (d *Dao) BlockAccount(c context.Context, mid int64, ftime int64, notify bool, freason int32, originTitle string, originContent string, redirectURL string, adname string, remark string) (err error) {

	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))

	if ftime == -1 {
		params.Set("blockTimeLength", "0")
		params.Set("blockForever", "1")
	} else {
		params.Set("blockForever", "0")
		params.Set("blockTimeLength", strconv.FormatInt(ftime, 10))
	}
	params.Set("blockRemark", remark)
	params.Set("operator", adname)
	params.Set("originType", "1")
	params.Set("originContent", originContent)
	params.Set("reasonType", strconv.FormatInt(int64(freason), 10))
	if notify {
		params.Set("isNotify", "1")
	} else {
		params.Set("isNotify", "0")
	}
	params.Set("originTitle", originTitle)
	params.Set("originUrl", redirectURL)

	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, _blockURL, "", params, &res); err != nil {
		log.Error("sendMsg error(%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = model.ErrMsgSend
		log.Error("sendMsg failed(%v) error(%v)", _apiMsgSend+"?"+params.Encode(), res.Code)
	}
	return
}

// TransferData transefer data
type TransferData struct {
	Oid        int64  `json:"oid"`
	Type       int32  `json:"type"`
	Mid        int64  `json:"mid"`
	RpID       int64  `json:"rp_id"`
	Operator   string `json:"operator"`
	OperatorID int64  `json:"oper_id"`
	Content    string `json:"origin_content"`
	Reason     int32  `json:"reason_type"`
	Title      string `json:"origin_title"`
	Link       string `json:"origin_url"`
	Ctime      int64  `json:"business_time"`
	OriginType int32  `json:"origin_type"`
}

// TransferArbitration transfer report to Arbitration.
func (d *Dao) TransferArbitration(c context.Context, rps map[int64]*model.Reply, rpts map[int64]*model.Report, adid int64, adname string, titles map[int64]string, links map[int64]string) (err error) {
	var data []TransferData
	for _, rp := range rps {
		if rpts[rp.ID] == nil {
			continue
		}
		rpt := rpts[rp.ID]
		d := TransferData{
			RpID:       rp.ID,
			Oid:        rp.Oid,
			Type:       rp.Type,
			Mid:        rp.Mid,
			Operator:   adname,
			OperatorID: adid,
			Content:    rp.Content.Message,
			Reason:     rpt.Reason,
			Ctime:      int64(rp.CTime),
			OriginType: 1,
			Title:      titles[rp.ID],
			Link:       links[rp.ID],
		}
		data = append(data, d)
	}
	content, err := json.Marshal(data)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Set("data", string(content))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, _transferURL, "", params, &res); err != nil {
		log.Error("sendMsg error(%v)", err)
		return
	}
	if res.Code != ecode.OK.Code() {
		err = model.ErrMsgSend
		log.Error("sendMsg failed(%v) error(%v)", _transferURL+"?"+params.Encode(), res.Code)
	}
	return
}
