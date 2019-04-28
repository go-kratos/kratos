package http

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/admin/main/push/model"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func addDPTask(ctx *bm.Context) {
	task := &model.DPTask{}
	if err := ctx.Bind(task); err != nil {
		return
	}
	if err := parseDpTask(task); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(nil, pushSrv.AddDPTask(context.Background(), task))
}

func parseDpTask(task *model.DPTask) (err error) {
	if task.ActivePeriodStr != "" {
		if err = json.Unmarshal([]byte(task.ActivePeriodStr), &task.ActivePeriods); err != nil {
			log.Error("parse ActivePeriod(%s) error(%v)", task.ActivePeriodStr, err)
			return
		}
	}
	if task.VipExpireStr != "" {
		if err = json.Unmarshal([]byte(task.VipExpireStr), &task.VipExpires); err != nil {
			log.Error("parse VipExpire(%s) error(%v)", task.VipExpireStr, err)
			return
		}
	}
	if task.AttentionStr != "" {
		if err = json.Unmarshal([]byte(task.AttentionStr), &task.Attentions); err != nil {
			log.Error("parse Attention(%s) error(%v)", task.AttentionStr, err)
			return
		}
	}
	task.PushTime = time.Unix(task.PushTimeUnix, 0)
	task.ExpireTime = time.Unix(task.ExpireTimeUnix, 0)
	task.Job = strconv.FormatInt(pushmdl.JobName(time.Now().UnixNano(), task.Summary, task.LinkValue, task.Group), 10)
	task.Status = int(pushmdl.TaskStatusPending)
	extra, _ := json.Marshal(pushmdl.TaskExtra{Group: task.Group})
	task.Extra = string(extra)
	return
}

func dpTaskInfo(ctx *bm.Context) {
	task := new(model.Task)
	if err := ctx.Bind(task); err != nil {
		return
	}
	id, _ := strconv.ParseInt(task.ID, 10, 64)
	if id < 1 {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	ctx.JSON(pushSrv.DpTaskInfo(ctx, id, task.Job))
}

func checkDpData(ctx *bm.Context) {
	ctx.JSON(nil, pushSrv.CheckDpData(ctx))
}
