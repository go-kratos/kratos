package http

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

const (
	_stateup   = "1"
	_statedown = "0"
)

func reglist(c *bm.Context) {
	var (
		err error
		v   = &model.Param{}
	)
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(tvSrv.RegList(c, v))
}

func saveReg(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			IndexType string `form:"index_type"`
			IndexTid  string `form:"index_tid"`
			Rank      string `form:"rank"`
			Title     string `form:"title" validate:"required"`
			PageId    string `form:"page_id"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.PageId == "" {
		c.JSON(nil, tvSrv.AddReg(c, v.Title, v.IndexType, v.IndexTid, v.Rank))
	} else {
		c.JSON(nil, tvSrv.EditReg(c, v.PageId, v.Title, v.IndexType, v.IndexTid))
	}
}

func upState(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			Pids  []int  `form:"page_id,split" validate:"required,min=1,dive,gt=0"`
			Valid string `form:"valid" validate:"required"`
		})
	)
	if err = c.Bind(v); err != nil {
		return
	}
	if !(v.Valid == _statedown || v.Valid == _stateup) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, tvSrv.UpState(c, v.Pids, v.Valid))
}

func regSort(c *bm.Context) {
	param := new(struct {
		IDs []int `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	})
	if err := c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, tvSrv.RegSort(c, param.IDs))
}
