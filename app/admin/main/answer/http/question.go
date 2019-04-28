package http

import (
	"context"
	"go-common/app/admin/main/answer/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"mime/multipart"
)

func queDisable(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		Qids     []int64 `form:"id,split"`
		Operator string  `form:"operator"`
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	if len(arg.Qids) > model.MaxCount {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = username.(string)
	if err = answerSvc.BatchUpdateState(c, arg.Qids, model.StageDisable, arg.Operator); err != nil {
		c.JSON(nil, err)
		return
	}
	// if arg.State == 1 {
	// 	answerSvc.CreateBFSImg(c, arg.Qids)
	// }
	c.JSON(nil, nil)
}

func quesList(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgQue)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(answerSvc.QuestionList(c, arg))
}

func types(c *bm.Context) {
	c.JSON(answerSvc.Types(c))
}

func uploadQsts(c *bm.Context) {
	var (
		f   multipart.File
		h   *multipart.FileHeader
		err error
	)
	f, h, err = c.Request.FormFile("file")
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	c.JSON(answerSvc.UploadQsts(c, f, h, username.(string)))
}

func queEdit(c *bm.Context) {
	var (
		err error
		arg = new(model.QuestionDB)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	arg.Operator = username.(string)
	c.JSON(answerSvc.QuestionEdit(c, arg))
}

func loadImg(c *bm.Context) {
	c.JSON(nil, answerSvc.LoadImg(context.Background()))
}

func queHistory(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgHistory)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(answerSvc.QueHistory(c, arg))
}

func history(c *bm.Context) {
	var (
		err error
		arg = new(model.ArgHistory)
	)
	if err = c.Bind(arg); err != nil {
		return
	}
	c.JSON(answerSvc.History(c, arg))
}
