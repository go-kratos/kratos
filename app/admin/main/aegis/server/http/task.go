package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/aegis/model/business"
	"go-common/app/admin/main/aegis/model/common"
	taskmod "go-common/app/admin/main/aegis/model/task"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	libtime "go-common/library/time"
)

func taskDelay(c *bm.Context) {
	opt := &taskmod.DelayOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}
	err := srv.Delay(c, opt)
	if err == ecode.AccessDenied || err == ecode.NothingFound {
		c.JSONMap(map[string]interface{}{"tips": "任务已被他人处理，毋需延迟"}, nil)
		return
	}
	c.JSON(nil, err)
}

func taskRelease(c *bm.Context) {
	opt := &common.BaseOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.Release(c, opt, false))
}

func taskUnDo(c *bm.Context) {
	opt := &common.BaseOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.UnDoStat(c, opt))
}

func taskStat(c *bm.Context) {
	opt := &common.BaseOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.TaskStat(c, opt))
}

func consumerOn(c *bm.Context) {
	opt := &common.BaseOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(nil, srv.On(c, opt))
}

func consumerOff(c *bm.Context) {
	opt := new(struct {
		common.BaseOptions
		Delay bool `form:"delay" default:"true"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}
	if err := srv.Off(c, &opt.BaseOptions); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)

	srv.Release(c, &opt.BaseOptions, opt.Delay)
}

func kickOut(c *bm.Context) {
	opt := new(struct {
		common.BaseOptions
		KickUID int64 `form:"kick_uid" validate:"required"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}
	if opt.BaseOptions.Role != taskmod.TaskRoleLeader {
		httpCode(c, "组长才有权限踢人", ecode.RequestErr)
		return
	}

	if err := srv.KickOut(c, &opt.BaseOptions, opt.KickUID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)

	srv.Release(c, &opt.BaseOptions, false)
}

func consumerWatcher(c *bm.Context) {
	opt := new(struct {
		BizID  int64 `form:"business_id" validate:"required"`
		FlowID int64 `form:"flow_id" validate:"required"`
		Role   int8  `form:"role"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.Watcher(c, opt.BizID, opt.FlowID, opt.Role))
}

func configEdit(c *bm.Context) {
	twc, _, err := readConfig(c)
	if err != nil || twc == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if twc.ID <= 0 {
		httpCode(c, "id错误", ecode.RequestErr)
		return
	}

	c.JSON(nil, srv.UpdateConfig(c, twc))
}

func configAdd(c *bm.Context) {
	twc, confJSON, err := readConfig(c)
	if err != nil || twc == nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.AddConfig(c, twc, confJSON))
}

func readConfig(c *bm.Context) (tc *taskmod.Config, confJSON interface{}, err error) {
	opt := &taskmod.ConfigOption{}
	if err = c.Bind(opt); err != nil {
		return nil, nil, err
	}

	switch opt.ConfType {
	case taskmod.TaskConfigAssign:
		type Assign struct {
			Mids string `json:"mids"`
			Uids string `json:"uids"`
		}
		confJSON = new(Assign)
	case taskmod.TaskConfigRangeWeight:
		confJSON = new(taskmod.RangeWeightConfig)
	case taskmod.TaskConfigEqualWeight:
		confJSON = new(taskmod.EqualWeightConfig)
	default:
		return nil, nil, ecode.RequestErr
	}

	tc = &taskmod.Config{
		ID:          opt.ID,
		ConfType:    opt.ConfType,
		BusinessID:  opt.BusinessID,
		FlowID:      opt.FlowID,
		Description: opt.Description,
		UID:         opt.UID,
		Uname:       opt.Uname,
	}
	bt, _ := time.ParseInLocation(common.TimeFormat, opt.Btime, time.Local)
	et, _ := time.ParseInLocation(common.TimeFormat, opt.Etime, time.Local)
	if !bt.IsZero() {
		tc.Btime = libtime.Time(bt.Unix())
	}
	if !et.IsZero() {
		tc.Etime = libtime.Time(et.Unix())
	}

	err = jsonhelp(confJSON, []byte(opt.ConfJSON), tc)
	return
}

func jsonhelp(jsonObject interface{}, bs []byte, tc *taskmod.Config) (err error) {
	if err = json.Unmarshal(bs, jsonObject); err != nil {
		log.Error("jsonhelp error(%v)", string(bs), err)
		return
	}

	if bs, err = json.Marshal(jsonObject); err != nil {
		log.Error("jsonhelp json.Marshal(%s) error(%v)", err)
		return
	}
	tc.ConfJSON = string(bs)
	return
}

func configList(c *bm.Context) {
	qp := &taskmod.QueryParams{}
	if err := c.Bind(qp); err != nil {
		return
	}

	configs, count, err := srv.QueryConfigs(c, qp)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	qp.Pager.Total = int(count)

	c.JSONMap(map[string]interface{}{
		"pager": qp.Pager,
		"data":  configs,
	}, nil)
}

func configDelete(c *bm.Context) {
	opt := new(struct {
		ID int64 `form:"id" validate:"required"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(nil, srv.DeleteConfig(c, opt.ID))
}

func configSet(c *bm.Context) {
	opt := new(struct {
		ID    int64 `form:"id" validate:"required"`
		State int8  `form:"state"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(nil, srv.SetStateConfig(c, opt.ID, opt.State))
}

func maxWeight(c *bm.Context) {
	opt := &common.BaseOptions{}
	if err := c.Bind(opt); err != nil {
		return
	}

	c.JSON(srv.MaxWeight(c, opt))
}

func weightlog(c *bm.Context) {
	opt := new(struct {
		TaskID int64 `form:"task_id" validate:"required"`
		Pn     int   `form:"pn" default:"1"`
		Ps     int   `form:"ps" default:"20"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	ls, count, err := srv.WeightLog(c, opt.TaskID, opt.Pn, opt.Ps)
	c.JSONMap(map[string]interface{}{
		"data": ls,
		"pager": common.Pager{
			Pn:    opt.Pn,
			Ps:    opt.Ps,
			Total: count,
		},
	}, err)
}

func role(c *bm.Context) {
	opt := &common.BaseOptions{}
	if err := c.Bind(opt); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	opt.UID = uid(c)
	opt.Uname = uname(c)
	_, role, err := srv.GetRole(c, opt)
	c.JSON(role, err)
}

func roleFlush(c *bm.Context) {
	opt := new(struct {
		Uids   []int64 `form:"uids,split" validate:"required"`
		BizID  int64   `form:"business_id" validate:"required"`
		FlowID int64   `form:"flow_id" validate:"required"`
	})
	if err := c.Bind(opt); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.FlushRole(c, opt.BizID, opt.FlowID, opt.Uids))
}

func checkTaskRole() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		opt := &common.BaseOptions{}
		if err := ctx.Bind(opt); err != nil {
			ctx.Abort()
			return
		}

		if opt.BusinessID <= 0 || opt.FlowID <= 0 {
			httpCode(ctx, "缺少business_id或flow_id", ecode.RequestErr)
			ctx.Abort()
			return
		}

		if srv.Debug() == "local" {
			ctx.Request.Form.Set("role", strconv.Itoa(int(taskmod.TaskRoleLeader)))
			return
		}
		user := uid(ctx)
		if srv.IsAdmin(user) {
			ctx.Request.Form.Set("role", strconv.Itoa(int(taskmod.TaskRoleLeader)))
			log.Info("checkTaskRole uid(%d) is admin", user)
			return
		}

		_, role, err := srv.GetRole(ctx, opt)
		if err != nil || role == 0 {
			ctx.JSON(nil, ecode.AccessDenied)
			ctx.Abort()
			return
		}

		log.Info("checkTaskRole opt(%+v) role(%d)", role)
		ctx.Request.Form.Set("role", fmt.Sprint(role))
	}
}

//判断业务用户权限，无授权业务则报错
func checkBizRole(role string, noAdmin bool, bizID bool) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		if srv.Debug() == "local" {
			return
		}

		user := uid(ctx)
		if srv.IsAdmin(user) {
			log.Info("checkBizRole uid(%d) is admin", user)
			return
		}

		businessID, err := srv.GetRoleBiz(ctx, user, role, noAdmin)
		if err != nil || len(businessID) == 0 {
			ctx.JSON(nil, ecode.AccessDenied)
			ctx.Abort()
			return
		}
		ctx.Set(business.AccessBiz, businessID)

		//request business in accessed biz range
		if bizID {
			var biz int64
			if strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
				var body []byte
				if body, err = ioutil.ReadAll(ctx.Request.Body); err != nil {
					log.Error("checkBizRole ioutil.ReadAll error(%+v)", err)
					return
				}
				ctx.Request.Body.Close()
				ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				pm := new(struct {
					BusinessID int64 `json:"business_id" form:"business_id"`
				})
				if err = json.Unmarshal(body, pm); err != nil {
					log.Error("checkBizRole json.Unmarshal error(%+v) body(%s)", err, string(body))
					ctx.JSON(nil, ecode.RequestErr)
					ctx.Abort()
					return
				}
				biz = pm.BusinessID
			} else {
				if biz, err = strconv.ParseInt(ctx.Request.Form.Get("business_id"), 10, 64); err != nil {
					log.Error("checkBizRole strconv.ParseInt(%s) error(%v)", ctx.Request.Form.Get("business_id"), err)
					ctx.JSON(nil, ecode.RequestErr)
					ctx.Abort()
					return
				}
			}

			exist := false
			for _, item := range businessID {
				if item == biz {
					exist = true
					break
				}
			}
			if !exist || biz <= 0 {
				ctx.JSON(nil, ecode.AccessDenied)
				ctx.Abort()
				return
			}
		}
	}
}

//判断任务用户权限，无授权业务则报错
func checkAccessTask() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		if srv.Debug() == "local" {
			return
		}
		user := uid(ctx)
		if srv.IsAdmin(user) {
			log.Info("checkrole uid(%d) is admin", user)
			return
		}
		businessID, flowID, err := srv.GetTaskBizFlows(ctx, user)
		if err != nil || len(businessID) == 0 || len(flowID) == 0 {
			//ctx.JSON(nil, ecode.AccessDenied)
			ctx.JSON(nil, nil)
			ctx.Abort()
			return
		}
		log.Info("checkAccessTask uid(%d) can see business(%+v) flow(%+v)", user, businessID, flowID)
		ctx.Set(business.AccessBiz, businessID)
		ctx.Set(business.AccessFlow, flowID)
	}
}

//是否为指定业务的管理员角色
func checkBizID() bm.HandlerFunc {
	return checkBizRole(business.BizBIDAdmin, false, true)
}

//为哪些业务的管理员角色
func checkBizAdmin() bm.HandlerFunc {
	return checkBizRole(business.BizBIDAdmin, false, false)
}

//为哪些业务的非管理员角色
func checkBizBID() bm.HandlerFunc {
	return checkBizRole("", true, false)
}

func checkBizLeader() bm.HandlerFunc {
	return checkBizRole("leader", true, true)
}

func checkBizOper() bm.HandlerFunc {
	return checkBizRole("oper", true, true)
}

//为指定业务的非管理员角色
func checkBizBIDBiz() bm.HandlerFunc {
	return checkBizRole("", true, true)
}

func getAccessBiz(c *bm.Context) (biz []int64) {
	biz = []int64{}
	ids, _ := c.Get(business.AccessBiz)
	if ids != nil {
		biz = ids.([]int64)
	}
	return
}

func getAccessFlow(c *bm.Context) (flowID []int64) {
	flowID = []int64{}
	ids, _ := c.Get(business.AccessFlow)
	if ids != nil {
		flowID = ids.([]int64)
	}
	return
}

func checkon() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		if srv.Debug() == "local" {
			return
		}
		opt := &common.BaseOptions{}
		if err := ctx.Bind(opt); err != nil {
			ctx.Abort()
			return
		}

		if !srv.IsOn(ctx, opt) {
			if err := srv.On(ctx, opt); err != nil {
				ctx.JSON(nil, err)
				ctx.Abort()
			}
			return
		}
	}
}

func preHandlerUser() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		if srv.Debug() == "local" {
			return
		}
		uidS, ok := ctx.Get("uid")
		if !ok {
			ctx.JSON(nil, ecode.NoLogin)
			ctx.Abort()
			return
		}
		unamei, ok := ctx.Get("username")
		if !ok {
			ctx.JSON(nil, ecode.NoLogin)
			ctx.Abort()
			return
		}

		ctx.Request.Form.Set("uid", fmt.Sprint(uidS))
		ctx.Request.Form.Set("uname", unamei.(string))
	}
}
