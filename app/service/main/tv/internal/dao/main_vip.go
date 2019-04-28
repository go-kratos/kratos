package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/service/main/tv/internal/model"
	mvip "go-common/app/service/main/vipinfo/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_mvipGiftRemark = "电视大会员赠送"
)

// MainVip returns main vip info.
func (d *Dao) MainVip(c context.Context, mid int64) (mv *model.MainVip, err error) {
	var (
		res *mvip.InfoReply
	)
	res, err = d.mvipCli.Info(c, &mvip.InfoReq{Mid: int64(mid)})
	if err != nil {
		log.Error("d.MainVip(%d) err(%v)", mid, err)
		return
	}
	mv = &model.MainVip{
		Mid:        int64(mid),
		VipType:    int8(res.Res.Type),
		VipStatus:  int8(res.Res.Status),
		VipDueDate: res.Res.DueDate,
	}
	log.Info("d.MainVip(%d) res(%+v)", mid, res)
	return
}

// GiveMVipGift gives bilibili vip to user.
func (d *Dao) GiveMVipGift(c context.Context, mid int64, batchId int, orderNo string) error {
	var (
		err error
	)
	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		TTL     int    `json:"ttl"`
	})
	params := url.Values{}
	params.Set("batchId", strconv.Itoa(batchId))
	params.Set("mid", strconv.Itoa(int(mid)))
	params.Set("orderNo", orderNo)
	params.Set("remark", _mvipGiftRemark)
	url := d.c.MVIP.BatchUserInfoUrl
	if err = d.mvipHttpCli.Post(c, url, "", params, res); err != nil {
		log.Error("d.mvipHttpCli.Post(%s, %+v) err(%+v)", url, params, err)
		return err
	}
	if res.Code != 0 {
		log.Error("d.mvipHttpCli.Post(%s, %+v) res(%+v)", url, params, res)
		return ecode.TVIPGiveMVipFailed
	}
	log.Info("d.GiveMVipGift(%d, %d, %s) res(%+v)", mid, batchId, orderNo, res)
	return nil
}
