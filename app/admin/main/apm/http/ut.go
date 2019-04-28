package http

import (
	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/ut"
	saga "go-common/app/tool/saga/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// @params ListReq
// @router get /x/admin/apm/ut/info/list
// @response Paper
func utList(c *bm.Context) {
	var (
		mrInfs []*ut.Merge
		data   *Paper
		err    error
		count  int
	)
	v := new(ut.MergeReq)
	if err = c.Bind(v); err != nil {
		return
	}
	if mrInfs, count, err = apmSvc.UtList(c, v); err != nil {
		c.JSON(nil, err)
		return
	}
	data = &Paper{
		Total: count,
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: mrInfs,
	}
	c.JSON(data, nil)
}

// @params DetailReq
// @router get /x/admin/apm/ut/detail/list
// @response PkgAnls
func utDetail(c *bm.Context) {
	var (
		utpkgs []*ut.PkgAnls
		err    error
	)
	v := new(ut.DetailReq)
	if err = c.Bind(v); err != nil {
		return
	}
	if utpkgs, err = apmSvc.UtDetailList(c, v); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(utpkgs, nil)
}

// @params HistoryCommitReq
// @router get /x/admin/apm/ut/history/commit
// @response Commit
func utHistoryCommit(c *bm.Context) {
	var (
		utcmts []*ut.Commit
		count  int
		err    error
		data   *Paper
	)
	v := new(ut.HistoryCommitReq)
	if err = c.Bind(v); err != nil {
		return
	}
	if utcmts, count, err = apmSvc.UtHistoryCommit(c, v); err != nil {
		c.JSON(nil, err)
		return
	}
	data = &Paper{
		Total: count,
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: utcmts,
	}
	c.JSON(data, nil)
}

func utBaseline(c *bm.Context) {
	data := map[string]int{
		"coverage": conf.Conf.UTBaseLine.Coverage,
		"passrate": conf.Conf.UTBaseLine.Passrate,
	}
	c.JSON(data, nil)
}

// @params commit_id
// @router get /x/admin/apm/ut/check
// @response Tyrant
func check(c *bm.Context) {
	var (
		err error
		ty  *ut.Tyrant
		res = new(struct {
			CommitID string `form:"commit_id" validate:"required"`
		})
	)
	if err = c.Bind(res); err != nil {
		c.JSON(nil, err)
		return
	}
	if ty, err = apmSvc.CheckUT(c, res.CommitID); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(ty, nil)
}

// @params merge_id
// @router get /x/admin/apm/ut/merge/set
// @response message
func utSetMerged(c *bm.Context) {
	var (
		err    error
		hookMR = &saga.HookMR{}
	)
	if err = c.BindWith(hookMR, binding.JSON); err != nil {
		return
	}
	if hookMR.ObjectAttributes.State != "merged" {
		c.JSON(nil, nil)
		return
	}
	if err = apmSvc.SetMerged(c, hookMR.ObjectAttributes.IID); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.WechatReport(c, hookMR.ObjectAttributes.IID, hookMR.ObjectAttributes.LastCommit.ID, hookMR.ObjectAttributes.SourceBranch, hookMR.ObjectAttributes.TargetBranch); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"message": "单元测试is_merged更新成功",
	}, nil)
}
