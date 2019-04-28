package http

import (
	"fmt"
	"strconv"

	"go-common/app/admin/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// medal.
func medalList(c *bm.Context) {
	var (
		err error
		res []*model.MedalInfo
	)
	if res, err = svc.Medal(c); err != nil {
		httpCode(c, err)
		return
	}
	httpData(c, res, nil)
}

// medalView .
func medalView(c *bm.Context) {
	arg := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	var (
		err error
		np  *model.MedalInfo
	)
	if np, err = svc.MedalView(c, arg.ID); err != nil {
		httpCode(c, err)
		return
	}
	np.Image = "http://i0.hdslb.com" + np.Image
	np.ImageSmall = "http://i0.hdslb.com" + np.ImageSmall
	httpData(c, np, nil)
}

// medalAdd add medal .
func medalAdd(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.Medal)
	if err = c.Bind(arg); err != nil {
		return
	}
	if err = svc.AddMedal(c, arg); err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func medalEdit(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.Medal)
	if err = c.Bind(arg); err != nil {
		return
	}
	if err = svc.UpMedal(c, arg.ID, arg); err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func medalGroup(c *bm.Context) {
	var (
		err error
		res []*model.MedalGroup
	)
	if res, err = svc.MedalGroupInfo(c); err != nil {
		httpCode(c, err)
		return
	}
	httpData(c, res, nil)
}

func medalGroupView(c *bm.Context) {
	arg := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	var (
		err error
		res *model.MedalGroup
	)
	if res, err = svc.MedalGroupByGid(c, arg.ID); err != nil {
		httpCode(c, err)
		return
	}
	httpData(c, res, nil)
}

func medalGroupParent(c *bm.Context) {
	var (
		err error
		res []*model.MedalGroup
	)
	if res, err = svc.MedalGroupParent(c); err != nil {
		httpCode(c, err)
		return
	}
	httpData(c, res, nil)
}

func medalGroupAdd(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.MedalGroup)
	if err = c.Bind(arg); err != nil {
		return
	}
	if err = svc.MedalGroupAdd(c, arg); err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func medalGroupEdit(c *bm.Context) {
	var (
		err error
	)
	arg := new(model.MedalGroup)
	if err = c.Bind(arg); err != nil {
		return
	}
	if err = svc.MedalGroupUp(c, arg.ID, arg); err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func medalMemberMID(c *bm.Context) {
	arg := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	var (
		err error
		res []*model.MedalMemberMID
	)
	if res, err = svc.MedalOwner(c, arg.MID); err != nil {
		httpCode(c, err)
		return
	}
	httpData(c, res, nil)
}

func medalOwnerUpActivated(c *bm.Context) {
	arg := new(struct {
		ID  int64 `form:"id" validate:"required"`
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	var (
		err error
	)
	if err = svc.MedalOwnerUpActivated(c, arg.MID, arg.ID); err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func medalMemberAddList(c *bm.Context) {
	arg := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	var (
		err error
		res []*model.MedalMemberAddList
	)
	if res, err = svc.MedalOwnerAddList(c, arg.MID); err != nil {
		httpCode(c, err)
		return
	}
	httpData(c, res, nil)
}

func medalMemberAdd(c *bm.Context) {
	arg := new(struct {
		MID     int64  `form:"mid" validate:"required"`
		NID     int64  `form:"nid" validate:"required"`
		Title   string `form:"title"`
		Message string `form:"message"`
		OID     int64  `form:"oper_id"  validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	var (
		err error
	)
	if err = svc.MedalOwnerAdd(c, arg.MID, arg.NID, arg.Title, arg.Message, arg.OID); err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func medalMemberDel(c *bm.Context) {
	arg := new(struct {
		MID     int64  `form:"mid" validate:"required"`
		NID     int64  `form:"nid" validate:"required"`
		IsDel   int8   `form:"is_del"`
		Title   string `form:"title"`
		Message string `form:"message"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	var (
		err error
	)
	if err = svc.MedalOwnerDel(c, arg.MID, arg.NID, arg.IsDel, arg.Title, arg.Message); err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func medalBatchAdd(c *bm.Context) {
	var (
		err error
		nid int64
	)
	f, h, err := c.Request.FormFile("file")
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}
	defer f.Close()
	params := c.Request.Form
	nidStr := params.Get("nid")
	nid, err = strconv.ParseInt(nidStr, 10, 64)
	if err != nil || nid <= 0 {
		fmt.Printf("nid:%+v\n", nid)
		httpCode(c, ecode.RequestErr)
		return
	}
	msg, err := svc.BatchAdd(c, nid, f, h)
	if err != nil || msg != "" {
		log.Error("svc.BatchAdd error(%v), msg(%v)", err, msg)
		httpCode(c, ecode.ServerErr)
		return
	}
	res := new(struct {
		Message string `form:"message"`
	})
	res.Message = msg
	httpData(c, res, nil)
}

func medalOperlog(c *bm.Context) {
	arg := new(struct {
		PN  int   `form:"pn"`
		PS  int   `form:"ps"`
		MID int64 `form:"mid" validate:"required"`
	})
	arg.PN, arg.PS = 1, 20
	if err := c.Bind(arg); err != nil {
		return
	}
	opers, pager, err := svc.MedalOperlog(c, arg.MID, arg.PN, arg.PS)
	if err != nil {
		log.Error("svc.MedalOperlog(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	if len(opers) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, opers, pager)
}
