package medal

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
)

// SendMsg send message.
func (d *Dao) SendMsg(c context.Context, mid int64, title string, context string) (err error) {
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
	if err = d.client.Post(c, d.sendMsgURI, "", params, &res); err != nil {
		log.Error("sendMsgURL(%s) error(%v)", d.sendMsgURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("sendMsgURL(%s) res(%d)", d.sendMsgURI+"?"+params.Encode(), res.Code)
	}
	log.Info("d.sendMsgURL url(%s) res(%d)", d.sendMsgURI+"?"+params.Encode(), res.Code)
	return
}

// GetWearedfansMedal get weared fans medals from live.
func (d *Dao) GetWearedfansMedal(c context.Context, mid int64, source int8) (isLove bool, err error) {
	params := url.Values{}
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("source", strconv.FormatInt(int64(source), 10))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Max       int8 `json:"max"`
			Cnt       int8 `json:"cnt"`
			MasterMax int8 `json:"master_max"`
		} `json:"data"`
	}
	if err = d.client.Post(c, d.getWaredFansMedalURI, "", params, &res); err != nil {
		log.Error("GetWearedfansMedal(%s) error(%v)", d.getWaredFansMedalURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("GetWearedfansMedal(%s) res(%d)", d.getWaredFansMedalURI+"?"+params.Encode(), res.Code)
	}
	log.Info("GetWearedfansMedal(%s) res(%+v)", d.getWaredFansMedalURI+"?"+params.Encode(), res)
	if res.Code == 0 {
		if res.Data.Cnt > 0 {
			isLove = true
			return
		}
	}
	return
}
