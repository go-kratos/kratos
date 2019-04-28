package http

import (
	"strconv"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func updateReportSummary(c *bm.Context) {
	reportSummary := model.ReportSummary{}
	if err := c.BindWith(&reportSummary, binding.JSON); nil != err {
		c.JSON(nil, err)
		return
	}
	var ResultMap = make(map[string]string)
	status, err := srv.UpdateReportSummary(&reportSummary)
	ResultMap["status"] = status
	c.JSON(ResultMap, err)
}

func queryReportSummarys(c *bm.Context) {
	qrsr := model.QueryReportSuRequest{}
	if err := c.BindWith(&qrsr, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	if err := qrsr.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	sessionID, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}

	res, err := srv.QueryReportSummarys(c, sessionID.Value, &qrsr)
	if err != nil {
		log.Error("queryScripts Error", err)
		return
	}
	c.JSON(res, err)
}

func queryReportByID(c *bm.Context) {
	v := new(struct {
		ID int `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(srv.QueryReportByID(v.ID))
}

func queryReGraph(c *bm.Context) {
	reGraphParam := model.QueryReGraphParam{}
	if err := c.BindWith(&reGraphParam, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	log.Info("TestNameNicks:(%s)", reGraphParam.TestNameNicks)
	reportGraphs, err := srv.QueryReGraph(reGraphParam.TestNameNicks)
	var resultGraphMap = make(map[string][][]model.ReportGraph)
	resultGraphMap["reportGraphs"] = reportGraphs
	c.JSON(resultGraphMap, err)
}

func queryReGraphAvg(c *bm.Context) {
	reGraphParam := model.QueryReGraphParam{}
	if err := c.BindWith(&reGraphParam, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	log.Info("TestNameNicks:(%s)", reGraphParam.TestNameNicks)
	reportGraphs, err := srv.QueryReGraphAvg(reGraphParam.TestNameNicks)
	var resultGraphMap = make(map[string][]model.ReportGraph)
	resultGraphMap["reportGraphs"] = reportGraphs
	c.JSON(resultGraphMap, err)
}

func updateReportStatus(c *bm.Context) {
	testStatus := c.Request.Form.Get("test_status")
	status, err := strconv.Atoi(testStatus)
	if err != nil {
		log.Error("test_status 输入错误，(%s)", err)
		return
	}
	c.JSON(srv.UpdateReportStatus(status), nil)
}

func delReport(c *bm.Context) {
	v := new(struct {
		ID int `form:"id"`
	})
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, srv.DelReportSummary(v.ID))

}
