package http

import (
	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/ut"
	bm "go-common/library/net/http/blademaster"
)

// @params PCurveReq
// @router get /x/admin/apm/ut/dashboard/pcurve
// @response PCurveResp
func utDashCurve(c *bm.Context) {
	var (
		curve []*ut.PCurveResp
		err   error
	)
	v := new(ut.PCurveReq)
	if err = c.Bind(v); err != nil {
		return
	}
	if curve, err = apmSvc.DashCurveGraph(c, name(c), v); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(curve, nil)
}

// @params PCurveReq
// @router get /x/admin/apm/ut/dashboard/histogram
// @response PCurveDetailResp
func utDashHistogram(c *bm.Context) {
	var (
		histogram []*ut.PCurveDetailResp
		err       error
	)
	v := new(ut.PCurveReq)
	if err = c.Bind(v); err != nil {
		return
	}
	if histogram, err = apmSvc.DashGraphDetail(c, name(c), v); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(histogram, nil)
}

// @params PCurveReq
// @router get /x/admin/apm/ut/dashboard/user/detail
// @response PCurveDetailResp
func utDashUserDetail(c *bm.Context) {
	var (
		detail []*ut.PCurveDetailResp
		err    error
	)
	v := new(ut.PCurveReq)
	if err = c.Bind(v); err != nil {
		return
	}
	if detail, err = apmSvc.DashGraphDetailSingle(c, name(c), v); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(detail, nil)
}

// @params QATrendReq
// @router get /x/admin/apm/ut/quality/trend
// @response QATrendResp
func utQATrend(c *bm.Context) {
	var (
		trend *ut.QATrendResp
		err   error
	)
	v := new(ut.QATrendReq)
	if err = c.Bind(v); err != nil {
		return
	}
	if trend, err = apmSvc.QATrend(c, v); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(trend, nil)
}

// @params commits
// @router get /x/admin/apm/ut/commits
// @response CommitInfo
func utGeneralCommit(c *bm.Context) {
	var (
		cmInfos []*ut.CommitInfo
		err     error
	)
	v := new(struct {
		Commits string `form:"commits"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if cmInfos, err = apmSvc.UTGernalCommit(c, v.Commits); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(cmInfos, nil)
}

// @params pkg
// @router get /x/admin/apm/ut/dashboard/pkgs
// @response []*ut.PkgAnls
func utDashPkgsTree(c *bm.Context) {
	var (
		err      error
		pkgs     []*ut.PkgAnls
		username = name(c)
		req      = new(struct {
			PKG string `form:"pkg"`
		})
	)
	if err = c.Bind(req); err != nil {
		c.JSON(nil, err)
		return
	}
	if pkgs, err = apmSvc.DashPkgsTree(c, req.PKG, username); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(pkgs, nil)
}

// @params project_id,merge_id,commit_id
// @router get /x/admin/apm/ut/git/report
// @response EmptyResp
func utGitReport(c *bm.Context) {
	var (
		err error
		req = new(struct {
			ProjectID int    `form:"project_id" validate:"required"`
			MergeID   int    `form:"merge_id" validate:"required"`
			CommitID  string `form:"commit_id" validate:"required"`
		})
	)
	if err = c.Bind(req); err != nil {
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.GitReport(c, req.ProjectID, req.MergeID, req.CommitID); err != nil {
		c.JSON(nil, err)
		return
	}
}

// @params username,times
// @router get /x/admin/apm/ut/dashboard/history/commit
// @response []*ut.PkgAnls
func dashHistoryCommit(c *bm.Context) {
	var (
		err  error
		pkgs = make([]*ut.PkgAnls, 0)
		req  = new(struct {
			UserName string `form:"user_name" default:""`
			Times    int64  `form:"times" default:"10"`
		})
	)
	if err = c.Bind(req); err != nil {
		c.JSON(nil, err)
		return
	}
	if req.UserName == "" {
		req.UserName = name(c)
	}
	if pkgs, err = apmSvc.CommitHistory(c, req.UserName, req.Times); err != nil {
		c.JSON(nil, err)
		return
	}
	data := new(struct {
		Pkgs     []*ut.PkgAnls `json:"pkgs"`
		BaseLine struct {
			Coverage int `json:"coverage"`
			PassRate int `json:"pass_rate"`
		} `json:"base_line"`
	})
	data.Pkgs = pkgs
	data.BaseLine.Coverage = conf.Conf.UTBaseLine.Coverage
	data.BaseLine.PassRate = conf.Conf.UTBaseLine.Passrate
	c.JSON(data, nil)
}
