package data

import (
	"context"
	"errors"
	"go-common/library/log"
	"net/url"
	"strconv"
)

const UpProfitStateSigned = 3 //激励计划签约状态

// UpProfitState 获取UP主激励计划状态
// 返回State：
// 1: 未申请; 2: 待审核; 3: 已签约; 4.已驳回; 5.主动退出; 6:被动退出; 7:封禁
func (d *Dao) UpProfitState(c context.Context, mid int64) (state int8, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("type", "0")

	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data *struct {
			Mid   int64 `json:"mid"`
			State int8  `json:"state"`
		} `json:"data"`
	}
	if err = d.client.Get(c, d.upProfitState, "", params, &res); err != nil {
		log.Error("UpProfitState(%d) error(%v)", mid, err)
		return
	}
	if res.Data == nil {
		err = errors.New("UP主激励计划状态获取失败")
		log.Error("UpProfitState(%d) nil response(%v)", res)
		return
	}
	state = res.Data.State
	return
}
