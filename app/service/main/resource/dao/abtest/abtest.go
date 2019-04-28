package abtest

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

// AbTest get abtest data from data-platform.
func (d *Dao) AbTest(c context.Context, names, ipaddr string) (adr []*model.AbTest, err error) {
	params := url.Values{}
	params.Set("groupNames", names)
	var res struct {
		Code int             `json:"code"`
		Data []*model.AbTest `json:"expItems"`
		Msg  string          `json:"msg"`
	}
	if err = d.httpClient.Get(c, d.testURL, ipaddr, params, &res); err != nil {
		log.Error("AbTest url(%s) error(%v)", d.testURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("AbTest api failed(%d)", res.Code)
		log.Error("CpmsApp url(%s) res code(%d) or res.data(%v)", d.testURL+"?"+params.Encode(), res.Code, res.Data)
		return
	}
	adr = res.Data
	return
}
