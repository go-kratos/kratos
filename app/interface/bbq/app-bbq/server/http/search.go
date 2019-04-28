package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	qn1080 = 80
)

//hotWord .
func hotWord(c *bm.Context) {
	arg := new(v1.HotWordRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.HotWord(c, arg))
}

func videoSearch(c *bm.Context) {
	arg := new(v1.BaseSearchReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	if arg.Qn == 0 {
		arg.Qn = qn1080
	}
	if arg.PageSize == 0 || arg.PageSize > 20 {
		arg.PageSize = 20
	}
	res, err := srv.VideoSearch(c, arg)
	c.JSON(res, err)

	// 埋点
	if err != nil {
		return
	}
	svidList := make([]int64, len(res.List))
	for i, v := range res.List {
		svidList[i] = v.SVID
	}
	ext := struct {
		Request *v1.BaseSearchReq
		Svid    []int64
	}{
		Request: arg,
		Svid:    svidList,
	}
	uiLog(c, model.ActionVideoSearch, ext)
}

func userSearch(c *bm.Context) {
	arg := new(v1.BaseSearchReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	mid, _ := c.Get("mid")
	if mid == nil {
		mid = int64(0)
	}
	if arg.Qn == 0 {
		arg.Qn = qn1080
	}
	if arg.PageSize == 0 || arg.PageSize > 20 {
		arg.PageSize = 20
	}
	res, err := srv.UserSearch(c, mid.(int64), arg)
	c.JSON(res, err)

	// 埋点
	if err != nil {
		return
	}
	midList := make([]int64, len(res.List))
	for i, v := range res.List {
		midList[i] = v.Mid
	}
	ext := struct {
		Request *v1.BaseSearchReq
		MID     []int64
	}{
		Request: arg,
		MID:     midList,
	}
	uiLog(c, model.ActionUserSearch, ext)
}

func sug(c *bm.Context) {
	arg := new(v1.SugReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	if arg.PageSize == 0 || arg.PageSize > 20 {
		arg.PageSize = 20
	}
	c.JSON(srv.BBQSug(c, arg))
}
