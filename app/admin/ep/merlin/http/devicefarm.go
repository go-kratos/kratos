package http

import (
	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func queryMobileDevice(c *bm.Context) {
	var (
		qdfr = &model.QueryMobileDeviceRequest{}
		err  error
	)
	if err = c.BindWith(qdfr, binding.JSON); err != nil {
		return
	}
	if err = qdfr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(svc.QueryMobileDeviceList(c, qdfr))
}

func refreshMobileDeviceDetail(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.RefreshMobileDeviceDetail(c, v.ID))
}

func bindMobileDevice(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.BindMobileDevice(c, v.ID, username))
}

func releaseMobileDevice(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.ReleaseMobileDevice(c, v.ID, username))
}

func lendOutMobileDevice(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.LendOutMobileDevice(c, v.ID, username))
}

func returnMobileDevice(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.ReturnMobileDevice(c, v.ID, username))
}

func startMobileDevice(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.StartMobileDevice(c, v.ID, username))
}

func shutDownMobileDevice(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.ShutdownMobileDevice(c, v.ID, username))
}

func syncMobileDevice(c *bm.Context) {
	c.JSON(svc.SyncMobileDeviceList(c))
}

func queryCategory(c *bm.Context) {
	var (
		v = new(struct {
			IsShowOfflineMB bool `form:"is_show_offline"`
		})
		err error
	)
	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.MobileMachineCategory(c, v.IsShowOfflineMB))
}

func reportMobileDeviceError(c *bm.Context) {
	var (
		v = new(struct {
			SerialName   string `json:"serial"`
			ErrorMessage string `json:"error_message"`
			ErrorCode    int    `json:"error_code"`
		})
		err error
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	c.JSON(nil, svc.ReportMobileDeviceError(c, v.SerialName, v.ErrorMessage, v.ErrorCode))
}

func queryMobileMachineLogs(c *bm.Context) {
	v := &model.QueryMobileMachineLogRequest{}
	if err := c.Bind(v); err != nil {
		return
	}

	if err := v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(svc.QueryMobileMachineLogs(c, v))
}

func queryMobileMachineLendOut(c *bm.Context) {
	var (
		v        = &model.QueryMobileMachineLogRequest{}
		err      error
		username string
	)
	if err = c.Bind(v); err != nil {
		return
	}

	if err = v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(svc.QueryMobileMachineLendOutLogs(c, username, v))
}

func queryMobileMachineErrorLogs(c *bm.Context) {
	v := &model.QueryMobileMachineErrorLogRequest{}
	if err := c.Bind(v); err != nil {
		return
	}

	if err := v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(svc.QueryMobileMachineErrorLogs(c, v))
}

func callbackMobileDeviceError(c *bm.Context) {
	var (
		v = new(struct {
			SerialName   string `json:"serial"`
			ErrorMessage string `json:"error_message"`
			ErrorCode    int    `json:"error_code"`
		})
		err error
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	c.JSON(nil, svc.CallBackMobileDeviceError(c, v.SerialName, v.ErrorMessage, v.ErrorCode))
}

func queryDeviceFarmSuperUser(c *bm.Context) {
	c.JSON(svc.DeviceFarmSuperUser(), nil)
}

func isBindByTheUser(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.IsBindMobileDeviceByTheUser(c, v.ID, username))
}
