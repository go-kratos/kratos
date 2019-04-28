package http

import (
	"strconv"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func list(ctx *bm.Context) {
	params := new(model.ListParams)
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Limit > 0 && (params.Limit <= (params.Pn-1)*params.Ps || params.Seed == "") {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	//params 默认值
	tformat := "2006-01-02 15:04:05"
	if params.CTimeFrom == "" && params.CTimeTo == "" {
		params.CTimeFrom = time.Now().AddDate(0, 0, -7).Format(tformat)
		params.CTimeTo = time.Now().Format(tformat)
	}
	if params.FTimeFrom != "" || params.FTimeTo != "" {
		params.State = model.QAStateFinish
	}
	if params.State != 0 && params.State != model.QAStateFinish {
		params.State = model.QAStateWait
	}

	list, err := srv.GetVideoList(ctx, params)
	ctx.JSON(list, err)
}

func detail(ctx *bm.Context) {
	idStr := ctx.Request.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	//任务详情
	detail, err := srv.GetDetail(ctx, id)
	ctx.JSON(detail, err)
}

func add(ctx *bm.Context) {
	//veri params
	params := new(model.AddVideoParams)
	if err := ctx.BindWith(params, binding.JSON); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	//insert
	taskID, err := srv.AddQATaskVideo(ctx, params)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	if taskID <= 0 {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	ctx.JSON(taskID, nil)
}

func submit(ctx *bm.Context) {
	uid, username := getUIDName(ctx)
	params := new(model.QASubmitParams)
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if _, exist := model.QAAuditStatus[params.AuditStatus]; !exist {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	if params.AuditStatus == model.VideoStatusRecycle && (params.TagID <= 0 || params.Reason == "") {
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	ctx.JSON(nil, srv.QAVideoSubmit(ctx, username, uid, params))
}

func upTaskUTime(ctx *bm.Context) {
	params := new(struct {
		TaskID int64 `form:"task_id" validate:"required,gt=0"`
		AID    int64 `form:"aid" validate:"required,gt=0"`
		CID    int64 `form:"cid" validate:"required,gt=0"`
		UTime  int64 `form:"utime"`
	})
	if err := ctx.Bind(params); err != nil {
		log.Error("upTaskUTime ctx.Bind error(%v) params(%+v)", err, ctx.Request.PostForm)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}

	ctx.JSON(nil, srv.UpVideoUTime(ctx, params.AID, params.CID, params.TaskID, params.UTime))
}
