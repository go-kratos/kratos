package http

import (
	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func delayMachineEndTime(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			MachineID int64 `form:"machine_id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.DelayMachineEndTime(c, v.MachineID, username))
}

func applyMachineEndTime(c *bm.Context) {
	var (
		applyEndTimeRequest = &model.ApplyEndTimeRequest{}
		err                 error
		username            string
	)
	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(applyEndTimeRequest, binding.JSON); err != nil {
		return
	}

	c.JSON(svc.ApplyDelayMachineEndTime(c, username, applyEndTimeRequest))
}

func cancelMachineEndTime(c *bm.Context) {
	var (
		username string
		v        = new(struct {
			AuditID int64 `form:"audit_id"`
		})
		err error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.CancelMachineEndTime(c, v.AuditID, username))
}

func auditMachineEndTime(c *bm.Context) {
	var (
		username            string
		auditEndTimeRequest = &model.AuditEndTimeRequest{}
		err                 error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(auditEndTimeRequest, binding.JSON); err != nil {
		return
	}

	c.JSON(svc.AuditMachineEndTime(c, auditEndTimeRequest.AuditID, username, auditEndTimeRequest.AuditResult, auditEndTimeRequest.Comment))
}

func queryApplicationRecordsByApplicant(c *bm.Context) {
	var (
		v        = &model.Pagination{}
		err      error
		username string
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	if err = v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(svc.GetApplicationRecordsByApplicant(c, username, v.PageNum, v.PageSize))
}

func queryApplicationRecordsByAuditor(c *bm.Context) {
	var (
		v        = &model.Pagination{}
		err      error
		username string
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}

	if err = v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(svc.GetApplicationRecordsByAuditor(c, username, v.PageNum, v.PageSize))
}

func queryApplicationRecordsByMachineID(c *bm.Context) {
	v := new(struct {
		model.Pagination
		MachineID int64 `form:"machine_id"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if err := v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(svc.GetApplicationRecordsByMachineID(c, v.MachineID, v.PageNum, v.PageSize))
}
