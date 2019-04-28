package http

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func queryPool(c *bm.Context) {
	var (
		arg   = new(model.ResoucePoolBo)
		pools []*model.VipResourcePool
		count int
		err   error
	)

	if err = c.Bind(arg); err != nil {
		return
	}
	if pools, count, err = vipSvc.QueryPool(c, arg); err != nil {
		c.JSON(nil, err)
		return
	}
	for _, v := range pools {
		var business *model.VipBusinessInfo
		if business, err = vipSvc.BusinessInfo(c, v.BusinessID); err != nil {
			return
		}
		v.BusinessName = business.BusinessName
	}
	dataMap := make(map[string]interface{})
	dataMap["item"] = pools
	dataMap["count"] = count
	c.JSONMap(dataMap, nil)
}

func getPool(c *bm.Context) {
	var (
		err      error
		arg      = new(model.ArgID)
		pool     *model.VipResourcePool
		business *model.VipBusinessInfo
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if pool, err = vipSvc.PoolInfo(c, int(arg.ID)); err != nil {
		c.JSON(nil, err)
		return
	}
	if pool == nil {
		c.JSON(pool, nil)
		return
	}
	if business, err = vipSvc.BusinessInfo(c, pool.BusinessID); err != nil {
		return
	}
	if business != nil {
		pool.BusinessName = business.BusinessName
	}

	c.JSON(pool, nil)

}

func savePool(c *bm.Context) {
	arg := new(model.ResoucePoolBo)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID == 0 {
		c.JSON(nil, vipSvc.AddPool(c, arg))
		return
	}
	c.JSON(nil, vipSvc.UpdatePool(c, arg))
}

func getBatch(c *bm.Context) {
	var (
		err   error
		arg   = new(model.ArgID)
		batch *model.VipResourceBatch
		pool  *model.VipResourcePool
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	if batch, err = vipSvc.BatchInfo(c, int(arg.ID)); err != nil {
		c.JSON(nil, err)
		return
	}
	if batch == nil {
		c.JSON(nil, nil)
		return
	}
	if pool, err = vipSvc.PoolInfo(c, batch.PoolID); err != nil {
		c.JSON(nil, err)
		return
	}
	if pool == nil {
		c.JSON(nil, ecode.VipPoolIDErr)
		return
	}
	vo := new(model.ResouceBatchVo)
	vo.Unit = batch.Unit
	vo.ID = batch.ID
	vo.PoolName = pool.PoolName
	vo.SurplusCount = batch.SurplusCount
	vo.PoolID = batch.PoolID
	vo.Count = batch.Count
	vo.StartTime = batch.StartTime
	vo.EndTime = batch.EndTime
	vo.DirectUseCount = batch.DirectUseCount
	vo.Ver = batch.Ver
	vo.CodeUseCount = batch.CodeUseCount
	c.JSON(vo, nil)
}

func queryBatch(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgPoolID)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.BatchInfoOfPool(c, arg.PoolID))
}

func addBatch(c *bm.Context) {
	var arg = new(model.ResouceBatchBo)

	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.AddBatch(c, arg))
}

func saveBatch(c *bm.Context) {
	var arg = new(model.ArgReSource)

	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(nil, vipSvc.UpdateBatch(c, arg.ID, arg.Increment, arg.StartTime, arg.EndTime))

}

func grantResouce(c *bm.Context) {
	var (
		req      = c.Request
		mids     []int
		username string
		file     multipart.File
		err      error
		failMids []int
		uc       *http.Cookie
	)

	if file, _, err = req.FormFile("uploadFile"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	r := bufio.NewReader(file)
	for {
		var buf []byte
		buf, _, err = r.ReadLine()
		if err == io.EOF {
			break
		}
		midStr := string(buf)
		if len(midStr) > 0 && midStr != "0" {
			mid := 0
			if mid, err = strconv.Atoi(midStr); err != nil {
				c.JSON(nil, ecode.RequestErr)
				return
			} else if err == nil {
				mids = append(mids, mid)
			}
		}
	}

	if uc, err = c.Request.Cookie("username"); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	username = uc.Value
	arg := new(struct {
		Mid     int    `form:"mid"`
		BatchID int64  `form:"batch_id" validate:"required"`
		Remark  string `form:"remark" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	if arg.Mid > 0 {
		mids = append(mids, arg.Mid)
	}

	if failMids, err = vipSvc.GrandResouce(c, arg.Remark, arg.BatchID, mids, username); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(failMids, nil)

}

func saveBatchCode(c *bm.Context) {
	arg := new(model.BatchCode)
	if err := c.Bind(arg); err != nil {
		return
	}
	username, exists := c.Get("username")
	if exists {
		arg.Operator = username.(string)
	}
	c.JSON(nil, vipSvc.SaveBatchCode(c, arg))
}

func batchCodes(c *bm.Context) {

	arg := new(model.ArgBatchCode)

	if err := c.Bind(arg); err != nil {
		return
	}
	pageArg := new(struct {
		PN int `form:"pn"`
		PS int `form:"ps"`
	})
	if err := c.Bind(pageArg); err != nil {
		return
	}

	res, total, err := vipSvc.SelBatchCode(c, arg, pageArg.PN, pageArg.PS)
	result := make(map[string]interface{})
	result["data"] = res
	result["total"] = total
	c.JSON(result, err)
}

func frozenCode(c *bm.Context) {
	arg := new(struct {
		ID     int64 `form:"id" validate:"required"`
		Status int8  `form:"status" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(nil, vipSvc.FrozenCode(c, arg.ID, arg.Status))
}
func frozenBatchCode(c *bm.Context) {
	arg := new(struct {
		ID     int64 `form:"id" validate:"required"`
		Status int8  `form:"status" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, vipSvc.FrozenBatchCode(c, arg.ID, arg.Status))
}
func codes(c *bm.Context) {
	arg := new(model.ArgCode)

	if err := c.Bind(arg); err != nil {
		return
	}

	pageArg := new(struct {
		PS       int    `form:"ps"`
		Cursor   int64  `form:"cursor"`
		Username string `form:"username"`
	})
	if err := c.Bind(pageArg); err != nil {
		return
	}
	username, exists := c.Get("username")
	if exists {
		pageArg.Username = username.(string)
	}
	res, cursor, pre, err := vipSvc.SelCode(c, arg, pageArg.Username, pageArg.Cursor, pageArg.PS)
	result := make(map[string]interface{})
	result["data"] = res
	result["cursor"] = cursor
	result["pre"] = pre
	c.JSON(result, err)
}

func exportCodes(c *bm.Context) {
	var (
		batchCodes []*model.BatchCode
		err        error
		codes      []string
	)
	arg := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	batchIds := []int64{arg.ID}
	if batchCodes, err = vipSvc.SelBatchCodes(c, batchIds); err != nil {
		c.JSON(nil, err)
		return
	}
	if len(batchCodes) == 0 {
		c.JSON(nil, ecode.VipBatchIDErr)
		return
	}

	if codes, err = vipSvc.ExportCode(c, arg.ID); err != nil {
		c.JSON(nil, err)
		return
	}
	writer := c.Writer
	header := writer.Header()
	header.Add("Content-disposition", "attachment; filename="+fmt.Sprintf("%v", batchCodes[0].BatchName)+".txt")
	header.Add("Content-Type", "application/x-download;charset=utf-8")
	for _, v := range codes {
		writer.Write([]byte(fmt.Sprintf("%v\r\n", v)))
	}
}
