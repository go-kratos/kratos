package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func batchInfo(c *bm.Context) {
	var (
		err error
		vib *model.VipResourceBatch
	)
	arg := new(struct {
		BatchID int64  `form:"batchId" validate:"required,min=1,gte=1"`
		Appkey  string `form:"appkey"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("batchInfo Bind error(%v)", err)
		return
	}
	if vib, _, err = vipSvc.BatchInfo(c, arg.BatchID, arg.Appkey); err != nil {
		log.Error("vipSvc.BatchInfo(%d) error(%v)", arg.BatchID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(vib, nil)
}

func useBatchInfo(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		BatchID int64  `form:"batchId" validate:"required,min=1,gte=1"`
		Appkey  string `form:"appkey"`
		Mid     int64  `form:"mid" validate:"required,min=1,gte=1"`
		OrderNo string `form:"orderNo" validate:"required"`
		Remark  string `form:"remark" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("useBatchInfo  Bind error(%v)", err)
		return
	}
	c.JSON(nil, vipSvc.ResourceBatchOpenVip(c, &model.ArgUseBatch{
		Mid:     arg.Mid,
		OrderNo: arg.OrderNo,
		Remark:  arg.Remark,
		Appkey:  arg.Appkey,
		BatchID: arg.BatchID,
	}))
}
