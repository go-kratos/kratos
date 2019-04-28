package message

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go-common/library/log"
	"go-common/library/xstr"
)

// Send message to upper
func (d *Dao) Send(c context.Context, mc, title, msg string, mids []int64, ts int64) (err error) {
	params := url.Values{}
	source := strings.Split(mc, "_")
	params.Set("type", "json")
	params.Set("source", source[0])
	params.Set("data_type", "4")
	params.Set("mc", mc)
	params.Set("title", title)
	params.Set("context", msg)
	var midList string
	for _, mid := range mids {
		midList += strconv.FormatInt(mid, 10)
		midList += ","
	}
	midList = strings.TrimSuffix(midList, ",")
	params.Set("mid_list", midList)
	var res struct {
		Code int `json:"code"`
	}
	log.Info("params:%v", params)
	if err = d.client.Post(c, d.uri, "", params, &res); err != nil {
		log.Error("growup-admin message url(%s) error(%v)", d.uri+"?"+params.Encode(), err)
		return
	}
	log.Info("message res code:%d", res.Code)
	if res.Code != 0 {
		log.Error("growup-admin message url(%s) error(%v)", d.uri+"?"+params.Encode(), err)
		err = fmt.Errorf("message send failed")
	}
	return
}

// NotifyTask notify task finish
func (d *Dao) NotifyTask(c context.Context, mids []int64) (err error) {
	params := url.Values{}
	params.Set("mids", xstr.JoinInts(mids))
	var res struct {
		Code int `json:"code"`
	}

	log.Info("creative notify task params:%v", params)
	if err = d.client.Post(c, d.creativeURL, "", params, &res); err != nil {
		log.Error("growup-admin creative notify task  url(%s) error(%v)", d.creativeURL+"?"+params.Encode(), err)
		return
	}
	log.Info("creative notify task res code:%d", res.Code)
	if res.Code != 0 {
		log.Error("growup-admin creative notify task  url(%s) error(%v)", d.creativeURL+"?"+params.Encode(), err)
		err = fmt.Errorf("creative notify task send failed")
	}
	return
}
