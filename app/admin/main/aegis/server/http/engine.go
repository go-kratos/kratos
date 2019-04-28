package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/business"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/middleware"
	"go-common/app/admin/main/aegis/model/task"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/xstr"
)

func listBizFlow(c *bm.Context) {
	bizID := getAccessBiz(c)
	flowID := getAccessFlow(c)
	opt := &business.OptList{}
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.ListBizFlow(c, opt.TP, bizID, flowID))
}

func getBizFlow(c *bm.Context) {
	opt := new(struct {
		BusinessID int64 `form:"business_id" validate:"required"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.ListBizFlow(c, 0, []int64{opt.BusinessID}, nil))
}

func next(c *bm.Context) {
	opt := &task.NextOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}
	if opt.BusinessID == 0 || opt.FlowID == 0 {
		httpCode(c, "缺少business_id或flow_id", ecode.RequestErr)
		return
	}
	opt.BisLeader = opt.Role == task.TaskRoleLeader

	c.JSON(srv.Next(c, opt))
}

func infoByTask(c *bm.Context) {
	opt := new(struct {
		TaskID int64 `form:"task_id" validate:"required"`
		common.BaseOptions
	})
	if err := c.Bind(opt); err != nil {
		return
	}
	if opt.BusinessID == 0 || opt.FlowID == 0 {
		httpCode(c, "缺少business_id或flow_id", ecode.RequestErr)
		return
	}

	c.JSON(srv.InfoTask(c, &opt.BaseOptions, opt.TaskID))
}

func listByTask(c *bm.Context) {
	opt := &task.ListOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}
	if opt.BusinessID == 0 || opt.FlowID == 0 {
		httpCode(c, "缺少business_id或flow_id", ecode.RequestErr)
		return
	}

	opt.BisLeader = opt.Role == task.TaskRoleLeader
	infos, err := srv.ListByTask(c, opt)
	c.JSONMap(map[string]interface{}{
		"data":  infos,
		"pager": &opt.Pager,
	}, err)
}

func infoByResource(c *bm.Context) {
	opt := new(common.BaseOptions)
	if err := c.Bind(opt); err != nil {
		return
	}
	if opt.BusinessID == 0 || opt.OID == "" {
		httpCode(c, "缺少business_id或oid", ecode.RequestErr)
		return
	}

	info, err := srv.InfoResource(c, opt)
	middleware.Response(
		info,
		err,
		c.JSON,
		&middleware.MiddleAggregate{
			Cfg:    srv.GetMiddlewareCache(opt.BusinessID),
			Encode: true,
		})
}

func listByResource(c *bm.Context) {
	opt := new(model.SearchParams)
	if err := c.Bind(opt); err != nil {
		return
	}

	middleware.Request(opt, &middleware.MiddleAggregate{
		Cfg:    srv.GetMiddlewareCache(opt.BusinessID),
		Encode: false,
	})
	columns, resources, operation, err := srv.ListByResource(c, opt)
	c.JSONMap(map[string]interface{}{
		"columns":    columns,
		"data":       resources,
		"operations": operation,
		"pager":      &opt.Pager,
	}, err)
}

func submit(c *bm.Context) {
	opt, err := parseOptions(c)
	if err != nil {
		if err == ecode.AegisBusinessCfgErr {
			c.JSON(nil, err)
		}
		return
	}

	if opt.BusinessID == 0 || opt.FlowID == 0 || opt.OID == "" || opt.Result == nil || opt.Binds == nil || len(opt.Binds) == 0 {
		httpCode(c, "business_id,flow_id,oid,result,binds不能为空", ecode.RequestErr)
		return
	}

	if err = srv.Submit(c, opt); err == ecode.AegisNotRunInFlow {
		c.JSONMap(map[string]interface{}{
			"tips": "资源已被流传,本次提交无效",
		}, nil)
		return
	}
	c.JSON(nil, err)
}

func listforjump(c *bm.Context) {
	opt := &common.BaseOptions{}
	if err := c.Bind(opt); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.FetchJumpFlowInfo(c, opt.FlowID))
}

func jump(c *bm.Context) {
	opt, err := parseOptions(c)
	if err != nil {
		if err == ecode.AegisBusinessCfgErr {
			c.JSON(nil, err)
		}
		return
	}
	if opt.BusinessID == 0 || opt.FlowID == 0 || opt.OID == "" || ((opt.Binds == nil || len(opt.Binds) == 0) && opt.NewFlowID == 0) {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	c.JSON(nil, srv.JumpFlow(c, opt))
}

func batchSubmit(c *bm.Context) {
	opt := &model.BatchOption{}
	if err := c.Bind(opt); err != nil {
		return
	}

	tips, err := srv.BatchSubmit(c, opt)
	if err != nil {
		c.JSON(tips, err)
		return
	}
	if tips != nil && len(tips.Fail) > 0 {
		log.Error("批量操作有错误 %+v", tips.Fail)
		msg, _ := json.Marshal(tips.Fail)
		httpCode(c, string(msg), ecode.ServerErr)
		return
	}
	c.JSON(tips, err)
}

func add(c *bm.Context) {
	opt := &model.AddOption{}
	if err := c.Bind(opt); err != nil {
		return
	}

	if opt.BusinessID == 0 || opt.NetID == 0 || opt.OID == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.Add(c, opt))
}

func cancel(c *bm.Context) {
	opt := new(model.CancelOption)
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(nil, srv.Cancel(c, opt.BusinessID, opt.Oids, 399, "业务方"))
}

func cancelByOper(c *bm.Context) {
	opt := new(struct {
		BissinessID int64    `form:"business_id" validate:"required"`
		Oids        []string `form:"oids,split" validate:"required"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	uid := uid(c)
	username := uname(c)
	c.JSON(nil, srv.CancelByOper(c, opt.BissinessID, opt.Oids, uid, username))
}

func update(c *bm.Context) {
	opt := new(model.UpdateOption)
	if err := c.BindWith(opt, binding.Form); err != nil {
		return
	}
	upParams := c.Request.Form.Get("update")
	if err := json.Unmarshal([]byte(upParams), &opt.Update); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.Update(c, opt))
}

func upload(c *bm.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error("FormFile err(%v)", err)
		httpCode(c, fmt.Sprintf("File Upload FormFile Error:(%v)", err), ecode.RequestErr)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(" uploadFile.ReadAll error(%v)", err)
		httpCode(c, fmt.Sprintf("File ioutil.ReadAll Error:(%v)", err), ecode.RequestErr)
		return
	}
	filetype := http.DetectContentType(content)
	switch filetype {
	case "image/jpeg", "image/jpg":
	case "image/gif":
	case "image/png":
	default:
		httpCode(c, fmt.Sprintf("not allow filetype(%s)", filetype), ecode.RequestErr)
		log.Warn("not allow filetype(%s) ", filetype)
		return
	}
	local, err := srv.Upload(c, "", filetype, time.Now().Unix(), content)
	if err != nil {
		log.Error("svc.Upload error(%v)", err)
		httpCode(c, fmt.Sprintf("svc.Upload error:(%v)", err), ecode.RequestErr)
		return
	}
	c.JSON(local, nil)
}

func parseOptions(c *bm.Context) (opt *model.SubmitOptions, err error) {
	opt = &model.SubmitOptions{}
	if err = c.BindWith(opt, binding.JSON); err != nil {
		log.Error("parseOptions err(%v)", err)
		return
	}

	if opt.Binds, err = xstr.SplitInts(opt.BindStr); err != nil {
		log.Error("parseOptions binds err(%v)", err)
		err = ecode.RequestErr
		return
	}

	if uidi, ok := c.Get("uid"); ok {
		opt.UID, _ = uidi.(int64)
	}
	if unamei, ok := c.Get("username"); ok {
		opt.Uname, _ = unamei.(string)
	}

	if opt.Result != nil {
		var cfg map[string]uint
		if cfg, err = srv.AttributeCFG(c, opt.BusinessID); err != nil {
			err = ecode.AegisBusinessCfgErr
			return
		} else if len(cfg) > 0 {
			opt.Result.AttrSet(cfg)
		}

		if opt.ExtraData != nil {
			if re, ok := opt.ExtraData["reason_extend"]; ok && len(fmt.Sprint(re)) > 0 {
				opt.Result.RejectReason = fmt.Sprintf("[%v]%s", re, opt.Result.RejectReason)
			}
		}
	}
	return
}

func httpCode(c *bm.Context, msg string, err error) {
	if c.IsAborted() {
		return
	}
	c.Error = err
	bcode := ecode.Cause(err)
	if msg == "" {
		msg = err.Error()
	}
	c.Render(http.StatusOK, render.JSON{
		Code:    bcode.Code(),
		Message: msg,
		Data:    nil,
	})
}

func track(c *bm.Context) {
	pm := new(model.TrackParam)
	if err := c.Bind(pm); err != nil {
		return
	}
	if pm.Pn > 1 && pm.LastPageTime == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	data, pager, err := srv.TrackResource(c, pm)
	c.JSONMap(map[string]interface{}{
		"data":  data,
		"pager": pager,
	}, err)
}

func auditLog(c *bm.Context) {
	pm := new(model.SearchAuditLogParam)
	if err := c.Bind(pm); err != nil {
		return
	}

	data, pger, err := srv.SearchAuditLog(c, pm)
	c.JSONMap(map[string]interface{}{
		"data":  data,
		"pager": pger,
	}, err)
}

func auditLogCSV(c *bm.Context) {
	pm := new(model.SearchAuditLogParam)
	csv, err := srv.SearchAuditLogCSV(c, pm)
	if err != nil {
		log.Error("auditLogCSV error(%v) pm(%v)", err, pm)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Title:   "操作日志",
		Content: FormatCSV(csv),
	})
}

func auth(c *bm.Context) {
	uid := uid(c)
	auth, err := srv.Auth(c, uid)
	c.JSONMap(map[string]interface{}{
		"data": auth,
	}, err)
}
