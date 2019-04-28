package http

import (
	"io/ioutil"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

func members(ctx *bm.Context) {
	arg := &model.ArgList{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.PN <= 0 {
		arg.PN = 1
	}
	if arg.PS <= 0 {
		arg.PS = 10
	}
	ctx.JSON(svc.Members(ctx, arg))
}

func memberProfile(ctx *bm.Context) {
	arg := &model.ArgMid{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.MemberProfile(ctx, arg.Mid))
}

func delSign(ctx *bm.Context) {
	arg := &model.ArgMids{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, svc.DelSign(ctx, arg))
}

// pubExpMsg is.
func pubExpMsg(ctx *bm.Context) {
	arg := &model.ArgPubExpMsg{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.IP == "" {
		arg.IP = metadata.String(ctx, metadata.RemoteIP)
	}
	if arg.Ts == 0 {
		arg.Ts = time.Now().Unix()
	}
	ctx.JSON(nil, svc.PubExpMsg(ctx, arg))
}

func expSet(ctx *bm.Context) {
	arg := &model.ArgExpSet{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(nil, svc.SetExp(ctx, arg))
}

func moralSet(ctx *bm.Context) {
	arg := &model.ArgMoralSet{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(nil, svc.SetMoral(ctx, arg))
}

func coinSet(ctx *bm.Context) {
	arg := &model.ArgCoinSet{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(nil, svc.SetCoin(ctx, arg))
}

func rankSet(ctx *bm.Context) {
	arg := &model.ArgRankSet{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	arg.IP = metadata.String(ctx, metadata.RemoteIP)
	ctx.JSON(nil, svc.SetRank(ctx, arg))
}

func additRemarkSet(ctx *bm.Context) {
	arg := &model.ArgAdditRemarkSet{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(nil, svc.SetAdditRemark(ctx, arg))
}

func baseReview(ctx *bm.Context) {
	arg := &model.ArgBaseReview{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.BaseReview(ctx, arg))
}

func clearFace(ctx *bm.Context) {
	arg := &model.ArgMids{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.ClearFace(ctx, arg), nil)
}

func clearSign(ctx *bm.Context) {
	arg := &model.ArgMids{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.ClearSign(ctx, arg), nil)
}

func clearName(ctx *bm.Context) {
	arg := &model.ArgMids{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.ClearName(ctx, arg), nil)
}

func expLog(ctx *bm.Context) {
	arg := &model.ArgMid{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.ExpLog(ctx, arg.Mid))
}

func faceHistory(ctx *bm.Context) {
	arg := &model.ArgFaceHistory{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if ctx.Request.Form.Get("status") == "" {
		// 0:队列中 1:待审核 2:通过 3:驳回
		arg.Status = []int8{0, 1, 2, 3}
	}
	if arg.PN <= 0 {
		arg.PN = 1
	}
	if arg.PS <= 0 {
		arg.PS = 50
	}
	if arg.ETime <= 0 {
		arg.ETime = xtime.Time(time.Now().Unix())
	}
	ctx.JSON(svc.FaceHistory(ctx, arg))
}

func moralLog(ctx *bm.Context) {
	arg := &model.ArgMid{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	ctx.JSON(svc.MoralLog(ctx, arg.Mid))
}

func batchFormal(ctx *bm.Context) {
	arg := &model.ArgBatchFormal{}
	if err := ctx.Bind(arg); err != nil {
		return
	}
	defer ctx.Request.Form.Del("file") // 防止日志不出现
	ctx.Request.ParseMultipartForm(32 << 20)
	fd, _, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Warn("Failed to parse form file: %+v", err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	defer fd.Close()
	file, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Warn("Failed to read form file: %+v", err)
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	log.Info("Succeeded to parse file data: file-length: %d", len(file))
	arg.FileData = file
	ctx.JSON(nil, svc.BatchFormal(ctx, arg))
}
