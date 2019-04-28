package http

import (
	"net/http"

	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func genMachines(c *bm.Context) {
	var (
		gmr      = &model.GenMachinesRequest{}
		err      error
		username string
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(gmr, binding.JSON); err != nil {
		return
	}
	if err = gmr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, svc.GenMachines(c, gmr, username))
}

func delMachine(c *bm.Context) {
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
	c.JSON(svc.DelMachineWhenCanBeDel(c, v.ID, username))
}

func queryMachineDetail(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.QueryMachineDetail(c, v.ID))
}

func queryMachinePackages(c *bm.Context) {
	c.JSON(svc.QueryMachinePackages(c))
}

func queryMachines(c *bm.Context) {
	var (
		qmr     = &model.QueryMachineRequest{}
		session *http.Cookie
		err     error
	)
	if err = c.Bind(qmr); err != nil {
		return
	}
	if err = qmr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	if session, err = c.Request.Cookie(_sessIDKey); err != nil {
		return
	}

	if qmr.Requester, err = getUsername(c); err != nil {
		return
	}

	c.JSON(svc.QueryMachines(c, session.Value, qmr))
}

func queryMachineLogs(c *bm.Context) {
	v := &model.QueryMachineLogRequest{}
	if err := c.Bind(v); err != nil {
		return
	}

	if err := v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(svc.QueryMachineLogs(c, v))
}

func queryMachineStatus(c *bm.Context) {
	v := new(struct {
		MachineID int64 `form:"machine_id"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.QueryMachineStatus(c, v.MachineID))
}

func transferMachine(c *bm.Context) {
	var (
		v = new(struct {
			MachineID int64  `form:"machine_id"`
			Receiver  string `form:"receiver"`
		})
		username string
		err      error
	)
	if err = c.Bind(v); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(svc.TransferMachine(c, v.MachineID, username, v.Receiver))
}
