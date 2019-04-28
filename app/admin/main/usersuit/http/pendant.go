package http

import (
	"encoding/csv"
	"io/ioutil"
	"strconv"
	"strings"

	"go-common/app/admin/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/xstr"
)

func pendantInfoList(c *bm.Context) {
	arg := new(model.ArgPendantGroupList)
	arg.PN, arg.PS = 1, 20
	if err := c.Bind(arg); err != nil {
		return
	}
	pis, pager, err := svc.PendantInfoList(c, arg)
	if err != nil {
		log.Error("svc.PendantInfoList(%+v) err(%v)", arg, err)
		return
	}
	if len(pis) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, pis, pager)
}

func pendantInfoID(c *bm.Context) {
	arg := new(struct {
		PID int64 `form:"pid" validate:"required"`
		GID int64 `form:"gid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	pi, err := svc.PendantInfoID(c, arg.PID, arg.GID)
	if err != nil {
		log.Error("svc.PendantInfoID(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpData(c, pi, nil)
}

func pendantGroupID(c *bm.Context) {
	arg := new(struct {
		GID int64 `form:"gid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	pg, err := svc.PendantGroupID(c, arg.GID)
	if err != nil {
		log.Error("svc.PendantGroupID(%+v) err(%v)", arg, err)
		return
	}
	httpData(c, pg, nil)
}

func pendantGroupList(c *bm.Context) {
	arg := new(model.ArgPendantGroupList)
	arg.PN, arg.PS = 1, 20
	if err := c.Bind(arg); err != nil {
		return
	}
	pgs, pager, err := svc.PendantGroupList(c, arg)
	if err != nil {
		log.Error("svc.PendantGroupList(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	if len(pgs) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, pgs, pager)
}

func pendantGroupAll(c *bm.Context) {
	pgs, err := svc.PendantGroupAll(c)
	if err != nil {
		log.Error("svc.PendantGroupAll() err(%v)", err)
		httpCode(c, err)
		return
	}
	httpData(c, pgs, nil)
}

func pendantInfoAllNoPage(c *bm.Context) {
	pis, err := svc.PendantInfoAllNoPage(c)
	if err != nil {
		log.Error("svc.pendantInfoAllNoPage() err(%v)", err)
		httpCode(c, err)
		return
	}
	httpData(c, pis, nil)
}

func addPendantInfo(c *bm.Context) {
	arg := new(model.ArgPendantInfo)
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svc.AddPendantInfo(c, arg)
	if err != nil {
		log.Error("svc.AddPendantInfo(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func upPendantInfo(c *bm.Context) {
	arg := new(model.ArgPendantInfo)
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svc.UpPendantInfo(c, arg)
	if err != nil {
		log.Error("svc.UpPendantInfo(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func upPendantGroupStatus(c *bm.Context) {
	arg := new(struct {
		GID    int64 `form:"gid" validate:"required"`
		Status int8  `form:"status"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svc.UpPendantGroupStatus(c, arg.GID, arg.Status)
	if err != nil {
		log.Error("svc.UpPendantGroupStatus(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func upPendantInfoStatus(c *bm.Context) {
	arg := new(struct {
		PID    int64 `form:"pid" validate:"required"`
		Status int8  `form:"status"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svc.UpPendantInfoStatus(c, arg.PID, arg.Status)
	if err != nil {
		log.Error("svc.UpPendantInfoStatus(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func addPendantGroup(c *bm.Context) {
	arg := new(model.ArgPendantGroup)
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svc.AddPendantGroup(c, arg)
	if err != nil {
		log.Error("svc.AddPendantGroup(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func upPendantGroup(c *bm.Context) {
	arg := new(model.ArgPendantGroup)
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svc.UpPendantGroup(c, arg)
	if err != nil {
		log.Error("svc.UpPendantGroup(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func pendantOrders(c *bm.Context) {
	arg := new(model.ArgPendantOrder)
	arg.PN, arg.PS = 1, 20
	if err := c.Bind(arg); err != nil {
		return
	}
	pos, pager, err := svc.PendantOrders(c, arg)
	if err != nil {
		log.Error("svc.PendantOrders(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	if len(pos) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, pos, pager)
}

func equipPendant(c *bm.Context) {
	arg := new(struct {
		UID int64 `form:"uid" validate:"required"`
		PID int64 `form:"pid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svc.EquipPendant(c, arg.UID, arg.PID)
	if err != nil {
		log.Error("svc.EquipPendant(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func userPendantPKG(c *bm.Context) {
	arg := new(struct {
		UID int64 `form:"uid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	pkg, equip, err := svc.PendantPKG(c, arg.UID)
	if err != nil {
		log.Error("svc.PendantPKG(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpData(c, map[string]interface{}{
		"pkg":   pkg,
		"equip": equip,
	}, nil)
}

func userPKGDetails(c *bm.Context) {
	arg := new(struct {
		UID int64 `form:"uid" validate:"required"`
		PID int64 `form:"pid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	pkgs, err := svc.UserPKGDetails(c, arg.UID, arg.PID)
	if err != nil {
		log.Error("svc.UserPKGDetails(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpData(c, pkgs, nil)
}

func upPendantPKG(c *bm.Context) {
	arg := new(model.ArgPendantPKG)
	if err := c.Bind(arg); err != nil {
		return
	}
	msg := &model.SysMsg{IsMsg: arg.IsMsg, Type: arg.Type, Title: arg.Title, Content: arg.Content, RemoteIP: metadata.String(c, metadata.RemoteIP)}
	err := svc.UpPendantPKG(c, arg.UID, arg.PID, arg.Day, msg, arg.OID)
	if err != nil {
		log.Error("svc.UpPendantPKG(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

// migrate old pedant info to new db.
func mutliSend(c *bm.Context) {
	var (
		err  error
		data []byte
		uids []int64
	)
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("upload err(%v)", err)
		httpCode(c, err)
		return
	}
	defer file.Close()
	data, err = ioutil.ReadAll(file)
	if err != nil {
		log.Error("ioutil.ReadAll err(%v)", err)
		return
	}
	r := csv.NewReader(strings.NewReader(string(data)))
	r.Comma = ','
	records, err := r.ReadAll()
	if err != nil {
		log.Error("r.ReadAll() err(%v)", err)
	}
	var (
		uid   int64
		muids = make(map[int64]struct{}, len(records))
	)
	for _, v := range records {
		if v[0] == "" {
			continue
		}
		if uid, err = strconv.ParseInt(v[0], 10, 64); err != nil {
			log.Error("strconv.ParseInt err(%v)", err)
			continue
		}
		if _, ok := muids[uid]; ok {
			continue
		}
		muids[uid] = struct{}{}
		uids = append(uids, uid)
	}
	if len(uids) == 0 {
		log.Warn("uids is nothing to send pendant")
		httpCode(c, ecode.RequestErr)
		return
	}
	if len(uids) > 1000 {
		httpCode(c, ecode.PendantAboveSendMaxLimit)
		return
	}
	params := c.Request.Form
	pidStr := params.Get("pid")
	dayStr := params.Get("day")
	isMsgStr := params.Get("is_msg")
	title := params.Get("title")
	content := params.Get("content")
	operStr := params.Get("oper_id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}
	day, err := strconv.ParseInt(dayStr, 10, 64)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}
	isMsg, err := strconv.ParseBool(isMsgStr)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}
	oid, err := strconv.ParseInt(operStr, 10, 64)
	if err != nil {
		httpCode(c, ecode.RequestErr)
		return
	}
	msg := &model.SysMsg{IsMsg: isMsg, Type: model.MsgTypeCustom, Title: title, Content: content, RemoteIP: metadata.String(c, metadata.RemoteIP)}
	if err = svc.MutliSendPendant(c, uids, pid, day, msg, oid); err != nil {
		log.Error("svc.MutliSendPendant(%s,%d,%d,%v,%d) err(%v)", xstr.JoinInts(uids), pid, day, msg, oid, err)
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func pendantOperlog(c *bm.Context) {
	arg := new(struct {
		PN int `form:"pn"`
		PS int `form:"ps"`
	})
	arg.PN, arg.PS = 1, 20
	if err := c.Bind(arg); err != nil {
		return
	}
	opers, pager, err := svc.PendantOperlog(c, arg.PN, arg.PS)
	if err != nil {
		log.Error("svc.PendantOperlog(%+v) err(%v)", arg, err)
		httpCode(c, err)
		return
	}
	if len(opers) == 0 {
		httpData(c, nil, pager)
		return
	}
	httpData(c, opers, pager)
}
