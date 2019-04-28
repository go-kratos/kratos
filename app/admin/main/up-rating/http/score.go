package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/app/admin/main/up-rating/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func scoreList(ctx *blademaster.Context) {
	arg := new(model.RatingListArg)
	if err := ctx.Bind(arg); err != nil {
		log.Error("scoreList ctx.Bind error(%v)", err)
		return
	}
	date := time.Now()
	if arg.ScoreDate != "" {
		var err error
		if date, err = time.ParseInLocation("2006-01", arg.ScoreDate, time.Local); err != nil {
			log.Error("date(%s) parse error", arg.ScoreDate)
			err = ecode.RequestErr
			return
		}
	}
	res, total, err := svr.RatingList(ctx, arg, date)
	if err != nil {
		log.Error("svr.RatingList error(%v) arg(%v)", err, arg)
		ctx.JSON(nil, err)
		return
	}
	ctx.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    res,
		"paging": &model.Paging{
			Ps:    arg.Limit,
			Total: total,
		},
	}))
}

func scoreCurrent(ctx *blademaster.Context) {
	arg := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err := ctx.Bind(arg); err != nil {
		log.Error("ctx.Bind error(%v)", err)
		return
	}
	res, err := svr.ScoreCurrent(ctx, arg.MID)
	if err != nil {
		log.Error("svr.ScoreCurrent error(%v) arg(%v)", err, arg)
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(res, nil)
}

func scoreHistory(ctx *blademaster.Context) {
	arg := new(model.UpRatingHistoryArg)
	if err := ctx.Bind(arg); err != nil {
		log.Error("ctx.Bind error(%v)", err)
		return
	}
	types := []model.ScoreType{
		model.Creativity,
		model.Influence,
		model.Credit,
	}
	if arg.ScoreType != model.Magnetic {
		types = []model.ScoreType{arg.ScoreType}
	}
	res, err := svr.ScoreHistory(ctx, types, arg.Mid, arg.Month == 0, arg.Month)
	if err != nil {
		log.Error("svr.RatingList error(%v) arg(%v)", err, arg)
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(&model.UpRatingHistoryResp{Data: res}, nil)
}

func scoreExport(ctx *blademaster.Context) {
	arg := new(model.RatingListArg)
	if err := ctx.Bind(arg); err != nil {
		log.Error("scoreList ctx.Bind error(%v)", err)
		return
	}
	date := time.Now()
	if arg.ScoreDate != "" {
		var err error
		if date, err = time.ParseInLocation("2006-01", arg.ScoreDate, time.Local); err != nil {
			log.Error("date(%s) parse error", arg.ScoreDate)
			err = ecode.RequestErr
			return
		}
	}
	content, err := svr.ExportScores(ctx, arg, date)
	if err != nil {
		ctx.JSON(nil, err)
		log.Error("up-rating svr.ExportScores error(%v)", err)
		return
	}
	ctx.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "scores"),
	})
}
