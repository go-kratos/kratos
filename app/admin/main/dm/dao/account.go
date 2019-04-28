package dao

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"

	"go-common/app/admin/main/dm/model"
	"go-common/library/log"
)

const (
	_addMoral     = "/api/moral/add"
	_blockUser    = "/x/internal/block/block"
	_blockInfoAdd = "/x/internal/credit/blocked/info/add"

	_blockArea      = "2"
	_blockSource    = "1"
	_blockForever   = "2"
	_blockTimeLimit = "1"
)

// ReduceMoral change moral
func (d *Dao) ReduceMoral(c context.Context, moral *model.ReduceMoral) (err error) {
	var (
		res = &struct {
			Code   int64             `json:"code"`
			Morals map[int64]float64 `json:"morals"`
		}{}
	)
	params := url.Values{}
	params.Set("mid", fmt.Sprint(moral.UID))
	params.Set("addMoral", fmt.Sprint(-math.Abs(float64(moral.Moral))))
	params.Set("origin", "2")
	params.Set("reason", model.AdminRptReason[moral.Reason])
	params.Set("reason_type", "1")
	params.Set("operater", moral.Operator)
	params.Set("is_notify", fmt.Sprint(moral.IsNotify))
	params.Set("remark", moral.Remark)
	err = d.httpCli.Get(c, d.addMoralURI, "", params, res)
	if err != nil {
		log.Error("d.httpCli.Get(%s) error(%v)", d.addMoralURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpCli.Get(%s) error(%v)", d.addMoralURI+"?"+params.Encode(), err)
	}
	return
}

// BlockUser block user
func (d *Dao) BlockUser(c context.Context, blu *model.BlockUser) (err error) {
	if err = d.blockUser(c, blu); err != nil {
		return
	}
	if err = d.blockInfoAdd(c, blu); err != nil {
		return
	}
	return
}

func (d *Dao) blockUser(c context.Context, blu *model.BlockUser) (err error) {
	var (
		res = new(struct {
			Code int `json:"data"`
		})
		params = url.Values{}
	)
	params.Set("mid", fmt.Sprint(blu.UID))
	params.Set("source", _blockSource)
	params.Set("area", _blockArea)
	if blu.BlockForever == 1 {
		params.Set("action", _blockForever)
	} else {
		params.Set("action", _blockTimeLimit)
	}
	params.Set("duration", fmt.Sprint(blu.BlockTimeLength*24*3600))
	params.Set("start_time", fmt.Sprint(time.Now().Unix()))
	params.Set("operator", blu.Operator)
	params.Set("reason", fmt.Sprint(blu.ReasonType))
	params.Set("comment", blu.BlockRemark)
	params.Set("notify", "0")
	err = d.httpCli.Post(c, d.blockUserURI, "", params, res)
	if err != nil {
		log.Error("d.httpCli.Post(%s) error(%v)", d.blockUserURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpCli.Post(%s) error(%v)", d.blockUserURI+"?"+params.Encode(), err)
	}
	return
}

func (d *Dao) blockInfoAdd(c context.Context, blu *model.BlockUser) (err error) {
	var (
		res = new(struct {
			Code int `json:"data"`
		})
		params = url.Values{}
	)
	params.Set("mid", fmt.Sprint(blu.UID))
	if blu.BlockForever == 1 {
		params.Set("blocked_forever", "1")
		params.Set("punish_type", "3")
	} else {
		params.Set("blocked_forever", "0")
		params.Set("punish_type", "2")
		if blu.BlockTimeLength == 0 {
			params.Set("punish_type", "1")
		}
	}
	params.Set("blocked_days", fmt.Sprint(blu.BlockTimeLength))
	params.Set("blocked_remark", blu.BlockRemark)
	params.Set("moral_num", fmt.Sprint(blu.Moral))
	params.Set("origin_content", fmt.Sprint(blu.OriginContent))
	params.Set("origin_title", fmt.Sprint(blu.OriginTitle))
	params.Set("origin_type", fmt.Sprint(blu.OriginType))
	params.Set("origin_url", fmt.Sprint(blu.OriginURL))
	params.Set("punish_time", fmt.Sprint(time.Now().Unix()))
	params.Set("reason_type", fmt.Sprint(blu.ReasonType))
	params.Set("operator_name", blu.Operator)
	err = d.httpCli.Post(c, d.blockInfoAddURI, "", params, res)
	if err != nil {
		log.Error("d.httpCli.Post(%s) error(%v)", d.blockInfoAddURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("return code:%d", res.Code)
		log.Error("d.httpCli.Post(%s) error(%v)", d.blockInfoAddURI+"?"+params.Encode(), err)
	}
	return
}
