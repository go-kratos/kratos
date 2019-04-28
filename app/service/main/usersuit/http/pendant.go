package http

import (
	"strconv"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	_paySuccess = "1"
)

// groupInfo get group info
func groupInfo(c *bm.Context) {
	c.JSON(usersuitSvc.GroupInfo(c))
}

// groupInfoByID get group info
func groupInfoByID(c *bm.Context) {
	var (
		err    error
		gid    int64
		params = c.Request.Form
		gidstr = params.Get("gid")
	)
	if gidstr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if gid, err = strconv.ParseInt(gidstr, 10, 64); err != nil || gid < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(usersuitSvc.GroupInfoByID(c, gid))
}

// vipGroup return vip group info
func vipGroup(c *bm.Context) {
	c.JSON(usersuitSvc.GroupInfoByID(c, 31))
}

// entryGroup return vip group info
func entryGroup(c *bm.Context) {
	c.JSON(usersuitSvc.GroupInfoByID(c, 30))
}

// pointRecommend get pendant which has point pay type
func pointRecommend(c *bm.Context) {
	var (
		err    error
		mid    int64
		params = c.Request.Form
		midstr = params.Get("mid")
	)
	if midstr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midstr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(usersuitSvc.PendantPoint(c, mid))
}

// pendantByID return pendant info
func pendantByID(c *bm.Context) {
	var (
		err    error
		pid    int64
		params = c.Request.Form
		pidstr = params.Get("pid")
	)
	if pidstr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pid, err = strconv.ParseInt(pidstr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(usersuitSvc.PendantInfo(c, pid))
}

// orderHistory query order history
func orderHistory(c *bm.Context) {
	var (
		err       error
		t         int64
		history   *model.ArgOrderHistory
		data      []*model.PendantOrderInfo
		count     map[string]int64
		req       = c.Request
		params    = req.Form
		midStr    = params.Get("mid")
		pidStr    = params.Get("pid")
		orderID   = params.Get("orderID")
		payID     = params.Get("payID")
		payType   = params.Get("payType")
		status    = params.Get("status")
		startTime = params.Get("start_time")
		endTime   = params.Get("end_time")
		page      = params.Get("page")
	)
	if midStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	history = new(model.ArgOrderHistory)
	if history.Mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || history.Mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pidStr != "" {
		if t, err = strconv.ParseInt(pidStr, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		history.Pid = t
	}
	if orderID != "" {
		history.OrderID = orderID
	}
	if payID != "" {
		history.PayID = payID
	}
	if payType != "" {
		if t, err = strconv.ParseInt(payType, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		history.PayType = int32(t)
	}
	if status != "" {
		if t, err = strconv.ParseInt(status, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		history.Status = int32(t)
	}
	if startTime != "" {
		if t, err = strconv.ParseInt(startTime, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		history.StartTime = t
	}
	if endTime != "" {
		if t, err = strconv.ParseInt(endTime, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		history.EndTime = t
	}
	if page != "" {
		if t, err = strconv.ParseInt(page, 10, 64); err != nil {
			history.Page = 1
		}
		history.Page = t
	} else {
		history.Page = 1
	}

	if data, count, err = usersuitSvc.OrderHistory(c, history); err != nil {
		log.Error("usersuitSvc.History error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"data":  data,
		"count": count,
	}, nil)
}

// packageInfo get user package
func packageInfo(c *bm.Context) {
	var (
		err    error
		mid    int64
		req    = c.Request
		midStr = req.Form.Get("mid")
	)
	if midStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(usersuitSvc.PackageInfo(c, mid))
}

// equipment get current equipped pendant
func equipment(c *bm.Context) {
	var (
		err    error
		mid    int64
		req    = c.Request
		params = req.Form
		midStr = params.Get("mid")
	)
	if midStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(usersuitSvc.Equipment(c, mid))
}

// order order pendant
func order(c *bm.Context) {
	var (
		err                   error
		mid, pid, expires, tp int64
		orderData             *model.PayInfo
		req                   = c.Request
		params                = req.Form
		midStr                = params.Get("mid")
		pidStr                = params.Get("pid")
		expiresStr            = params.Get("expires")
		typeStr               = params.Get("moneytype")
	)

	if midStr == "" || pidStr == "" || expiresStr == "" || typeStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if expires, err = strconv.ParseInt(expiresStr, 10, 64); err != nil || expires < 1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if tp, err = strconv.ParseInt(typeStr, 10, 64); err != nil || tp < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if orderData, err = usersuitSvc.OrderPendant(c, mid, pid, expires, tp); err != nil {
		c.JSON(nil, err)
		return
	}
	if orderData != nil && orderData.PayURL != "" {
		c.JSON(orderData, nil)
	}
}

// multiGrant batch grant pendant
func multiGrantByMid(c *bm.Context) {
	var (
		err         error
		pid, expire int64
		mids        []int64
		req         = c.Request
		params      = req.Form
		midStr      = params.Get("mids")
		pidStr      = params.Get("pid")
		expireStr   = params.Get("expire")
	)

	if midStr == "" || pidStr == "" || expireStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mids, err = xstr.SplitInts(midStr); err != nil || len(mids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if expire, err = strconv.ParseInt(expireStr, 10, 64); err != nil || expire <= 0 || expire > 3650 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, usersuitSvc.BatchGrantPendantByMid(c, pid, expire, mids))
}

// multiGrant batch grant pendant
func multiGrantByPid(c *bm.Context) {
	var (
		err           error
		mid           int64
		pids, expires []int64
		req           = c.Request
		params        = req.Form
		midStr        = params.Get("mid")
		pidStr        = params.Get("pids")
		expireStr     = params.Get("expires")
	)

	if midStr == "" || pidStr == "" || expireStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if pids, err = xstr.SplitInts(pidStr); err != nil || len(pids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if expires, err = xstr.SplitInts(expireStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, usersuitSvc.BatchGrantPendantByPid(c, mid, expires, pids))
}

// equip euqip pendant
func equip(c *bm.Context) {
	var (
		err                      error
		mid, pid, status, source int64
		req                      = c.Request
		params                   = req.Form
		midStr                   = params.Get("mid")
		pidStr                   = params.Get("pid")
		statusStr                = params.Get("status")
		sourceStr                = params.Get("source") // 挂件来源（可选）：默认0，0 未知，1背包，2大会员
	)

	if midStr == "" || pidStr == "" || statusStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if pid, err = strconv.ParseInt(pidStr, 10, 64); err != nil || pid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if status, err = strconv.ParseInt(statusStr, 10, 64); err != nil || (status != 1 && status != 2) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	source = model.ParseSource(sourceStr)

	c.JSON(nil, usersuitSvc.EquipPendant(c, mid, pid, int8(status), source))
}

// pendantCallback pay callback func
func pendantCallback(c *bm.Context) {
	var (
		err     error
		sFlag   int64
		arg     *model.PendantOrderInfo
		req     = c.Request
		params  = req.Form
		success = params.Get("isSuccess") // NOTE : how to  fixed 1
		orderID = params.Get("out_trade_no")
		payID   = params.Get("trade_no")
		status  = params.Get("trade_status")
	)

	if success == "" || orderID == "" || payID == "" || status == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if status != _paySuccess {
		c.JSON(nil, nil)
		return
	}

	arg = &model.PendantOrderInfo{OrderID: orderID, PayID: payID}
	if sFlag, err = strconv.ParseInt(status, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arg.Stauts = int32(sFlag)
	arg.IsCallback = 1

	c.JSON(nil, usersuitSvc.PendantCallback(c, arg))
}
