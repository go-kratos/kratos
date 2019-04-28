package http

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/up/model"
	"go-common/app/service/main/up/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/siddontang/go-mysql/mysql"
)

func register(ctx *bm.Context) {
	param := ctx.Request.Form
	midStr := param.Get("mid")
	fromStr := param.Get("from")
	isAuthorStr := param.Get("is_author")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 {
		log.Error("register error mid (%d)", mid)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fromStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	isAuthor, err := strconv.Atoi(isAuthorStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", isAuthorStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	result := 0
	row, err := Svc.Edit(ctx, mid, isAuthor, uint8(from))
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	if row > 0 {
		result = 1
	}
	ctx.JSON(map[string]interface{}{
		"result": result,
	}, nil)
}

func info(ctx *bm.Context) {
	param := ctx.Request.Form
	midStr := param.Get("mid")
	fromStr := param.Get("from")
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 {
		log.Error("info error mid (%d)", mid)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	from, err := strconv.Atoi(fromStr)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", fromStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	isAuthor, err := Svc.Info(ctx, mid, uint8(from))
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(map[string]interface{}{
		"is_author": isAuthor,
	}, nil)
}

func all(ctx *bm.Context) {
	param := ctx.Request.Form
	midStr := param.Get("mid")
	ip := metadata.String(ctx, metadata.RemoteIP)
	// check params
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if mid <= 0 {
		log.Error("up all error mid (%d)", mid)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	identifyAll, err := Svc.IdentifyAll(ctx, mid, ip)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(map[string]interface{}{
		"identify": identifyAll,
	}, nil)
}

// flows get specialUps list
func specialUps(c *bm.Context) {
	var res interface{}
	params := c.Request.Form
	groupStr := params.Get("group_id")
	var err error
	// check params
	//default all groups
	groupID := int64(0)
	if groupStr != "" {
		groupID, err = strconv.ParseInt(groupStr, 10, 64)
		if err != nil {
			log.Error("strconv.ParseInt(groupStr(%s)) error(%v)", groupStr, err)
			err = ecode.RequestErr
			c.JSON(res, err)
			return
		}
	}
	res = Svc.UpsByGroup(c, groupID)
	c.JSON(res, err)
}

func specialDel(c *bm.Context) {
	var res interface{}
	var err error
	var r = new(struct {
		Ids string `form:"ids" validate:"required"`
	})
	if err = c.Bind(r); err != nil {
		log.Error("params error, err=%v", err)
		err = ecode.RequestErr
		c.JSON(res, err)
		return
	}
	var idstr = strings.Split(r.Ids, ",")
	for _, s := range idstr {
		var affectedRow int64
		var id, e = strconv.ParseInt(s, 10, 64)
		err = e
		if err != nil {
			log.Error("id is not integer, id=%s", s)
			continue
		}
		affectedRow, err = Svc.SpecialDel(c, id)
		log.Info("delete special up, id=%d, affected=%d, err=%v", id, affectedRow, err)
	}
	c.JSON(res, err)
}

func specialAdd(c *bm.Context) {
	var res interface{}
	var err error
	var errMsg string
	var r = new(struct {
		GroupIds string `form:"group_ids" validate:"required"` // 支持多个group id，用,分隔
		MidStr   string `form:"mids" validate:"required"`      // 支持多个mid，用,分隔
		Note     string `form:"note" default:""`
	})
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("params error, err=%v", err)
			err = ecode.RequestErr
			errMsg = "params error"
			break
		}
		var groupIds = strings.Split(r.GroupIds, ",")
		// 检查是否有特殊权限
		for _, groupIDStr := range groupIds {
			var groupID, _ = strconv.ParseInt(groupIDStr, 10, 64)
			var e = Svc.SpecialGroupPermit(c, groupID)
			if e != nil {
				err = e
				break
			}
		}
		var midStrArray = strings.Split(r.MidStr, ",")
		if len(midStrArray) == 0 {
			log.Error("params error, no mid got, mid=%s", r.MidStr)
			err = ecode.RequestErr
			errMsg = "params error, no mid got"
			break
		}
		var mids []int64
		for _, v := range midStrArray {
			mid, e := strconv.ParseInt(v, 10, 64)
			if e != nil {
				continue
			}
			mids = append(mids, int64(mid))
		}
		if len(mids) == 0 {
			log.Error("params error, wrong mid got, mid=%s", r.MidStr)
			err = ecode.RequestErr
			errMsg = "params error, wrong mid got"
			break
		}

		uidtemp, ok := c.Get("uid")
		var uid int64
		if ok {
			uid = uidtemp.(int64)
		}

		if len(groupIds) == 0 {
			log.Error("params error, no group id got, groupid=%s", r.GroupIds)
			err = ecode.RequestErr
			errMsg = "params error, wrong mid got"
			break
		}
		var uname, _ = bmGetStringOrDefault(c, "username", "unkown")
		for _, groupIDStr := range groupIds {
			var groupID, err2 = strconv.ParseInt(groupIDStr, 10, 64)
			if err2 != nil {
				log.Warn("group id convert fail, group id=%s", groupIDStr)
				continue
			}
			var data = &model.UpSpecial{
				GroupID: groupID,
				Note:    r.Note,
				UID:     uid,
			}
			go func() {
				const step = 100
				for start := 0; start < len(mids); start += step {
					var end = start + step
					if end > len(mids) {
						end = len(mids)
					}
					_, err = Svc.SpecialAdd(context.Background(), uname, data, mids[start:end]...)
					if err != nil {
						log.Error("add special fail, mid=%s, err=%v", r.MidStr, err)
						errMsg = "server error"
						return
					}
					log.Info("add special, id=%v, group id=%d", mids, groupID)
				}

			}()

		}

	}

	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}

func specialEdit(c *bm.Context) {
	var res interface{}
	var err error
	var errMsg string
	var r = new(struct {
		GroupID int64  `form:"group_id" validate:"required"`
		Mid     int64  `form:"mid"`
		ID      int64  `form:"id" validate:"required"`
		Note    string `form:"note" default:""`
	})
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("request argument bind fail, err=%v", err)
			errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
			err = ecode.RequestErr
			break
		}
		// 检查是否有特殊权限
		err = Svc.SpecialGroupPermit(c, r.GroupID)
		if err != nil {
			break
		}

		uidtemp, ok := c.Get("uid")
		var uid int64
		if ok {
			uid = uidtemp.(int64)
		}
		var data = &model.UpSpecial{
			GroupID: r.GroupID,
			Note:    r.Note,
			UID:     uid,
			ID:      r.ID,
		}

		_, err = Svc.SpecialEdit(c, data, r.ID)
		if err != nil {
			log.Error("fail edit ups, err=%v, info=%v", err, data)
			errMsg = err.Error()
			break
		}
		log.Info("edit ups successful, info=%v", data)
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(res, err)
	}
}

// 支持更多的条件类型
func specialGet(c *bm.Context) {
	var data interface{}
	var err error
	var errMsg string
	var r = new(model.GetSpecialArg)
	var ups []*model.UpSpecialWithName
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("request argument bind fail, err=%v", err)
			errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
			err = ecode.RequestErr
			break
		}

		if r.Pn == 0 {
			r.Pn = 1
		}
		if r.Ps == 0 {
			r.Ps = 20
		}
		if r.Ps > 100 {
			r.Ps = 100
		}
		if r.Ps < 5 {
			r.Ps = 5
		}
		uidtemp, ok := c.Get("uid")
		var uid int64
		if ok {
			uid, _ = uidtemp.(int64)
		}

		var total int
		ups, total, err = Svc.SpecialGet(c, r)
		if err != nil {
			log.Error("fail to get special, err=%v, arg=%+v", err, r)
			errMsg = err.Error()
			err = ecode.ServerErr
			break
		}
		log.Info("get special successful, arg=%+v, result count=%d", len(ups), uid)
		data = model.UpsPage{
			Items: ups,
			Pager: &model.Pager{Num: r.Pn, Size: r.Ps, Total: total},
		}
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {

		if r.Export == "csv" {
			c.Writer.Header().Set("Content-Type", "application/csv")
			c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"up_special_%s.csv\"", time.Now().Format(mysql.TimeFormat)))

			var buf = &bytes.Buffer{}
			var csvWriter = csv.NewWriter(buf)
			csvWriter.Write((&model.UpSpecialWithName{}).GetTitleFields())
			for _, v := range ups {
				csvWriter.Write(v.ToStringFields())
			}
			csvWriter.Flush()
			c.Writer.Write(buf.Bytes())
		} else {
			c.JSON(data, err)
		}
	}
}

func specialGetByMid(c *bm.Context) {
	var data interface{}
	var err error
	var errMsg string
	var r = new(model.GetSpecialByMidArg)
	switch {
	default:
		if err = c.Bind(r); err != nil {
			log.Error("request argument bind fail, err=%v", err)
			errMsg = fmt.Sprintf("wrong argument, %s", err.Error())
			err = ecode.RequestErr
			break
		}

		uidtemp, ok := c.Get("uid")
		var uid int64
		if ok {
			uid = uidtemp.(int64)
		}
		var ups []*model.UpSpecial
		ups, err = Svc.SpecialGetByMid(c, r)
		if err != nil {
			log.Error("fail to get special, err=%v, arg=%+v", err, r)
			errMsg = err.Error()
			err = ecode.ServerErr
			break
		}
		log.Info("get special successful, arg=%+v, result count=%d", len(ups), uid)
		data = ups
	}
	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
	} else {
		c.JSON(data, err)
	}
}

func listUp(c *bm.Context) {
	var (
		mids      []int64
		newLastID int64
		err       error
		errMsg    string
		arg       = new(model.ListUpBaseArg)
	)
	switch {
	default:
		if err = c.Bind(arg); err != nil {
			log.Error("request argument bind fail, err=%v", err)
			errMsg = err.Error()
			err = ecode.RequestErr
			break
		}
		if !arg.Validate() {
			errMsg, err = "illegal argument", ecode.RequestErr
			break
		}
		if arg.Size == 0 {
			arg.Size = 100
		}

		mids, newLastID, err = Svc.ListUpBase(c, arg.Size, arg.LastID, arg.Activity)

		if err != nil {
			log.Error("fail to get special, err=%v, arg=%v", err, arg)
			errMsg = err.Error()
			err = ecode.ServerErr
			break
		}
	}

	if err != nil {
		service.BmHTTPErrorWithMsg(c, err, errMsg)
		return
	}

	c.JSON(map[string]interface{}{
		"result":  mids,
		"last_id": newLastID,
	}, nil)
}

func active(ctx *bm.Context) {
	var (
		err    error
		errMsg string
		req    = new(model.UpInfoActiveReq)
	)

	if err = ctx.Bind(req); err != nil {
		log.Error("request param error: req: %+v, err: %+v", req, err)
		return
	}
	if req.Mid <= 0 {
		log.Error("param error mid: %d", req.Mid)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	log.Info("req: %+v", req)
	res, err := Svc.GetUpInfoActive(ctx, req)
	if err != nil {
		log.Error("Svc.GetUpInfoActive error: %+v", err)
		service.BmHTTPErrorWithMsg(ctx, err, errMsg)
		return
	}

	ctx.JSON(res, nil)
}

func actives(ctx *bm.Context) {
	var (
		err    error
		errMsg string
		req    = new(model.UpsInfoActiveReq)
	)

	if err = ctx.Bind(req); err != nil {
		log.Error("request param error: req: %+v, err: %+v", req, err)
		return
	}
	if len(req.Mids) <= 0 {
		log.Error("param error mids: %+v", req.Mids)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	log.Info("req: %+v", req)
	res, err := Svc.GetUpsInfoActive(ctx, req)
	if err != nil {
		log.Error("Svc.GetUpsInfoActive error: %+v", err)
		service.BmHTTPErrorWithMsg(ctx, err, errMsg)
		return
	}

	ctx.JSON(res, nil)
}
