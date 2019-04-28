package http

import (
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/net/http/blademaster"
)

// archives handler
func archives(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.ArchivesReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "archives"))
}

// archiveHistoryList handler
func archiveHistoryList(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.ArchiveHistoryListReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "archiveHistoryList"))
}

// archiveVideos handler
func archiveVideos(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.ArchiveVideosReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "archiveVideos"))
}

// dataArchive handler
func dataArchive(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataArchiveReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataArchive"))
}

// dataVideoQuit handler
func dataVideoQuit(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataVideoQuitReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataVideoQuit"))
}

// danmuDistri handler
func danmuDistri(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DanmuDistriReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "danmuDistri"))
}

// dataBase handler
func dataBase(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataBaseReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataBase"))
}

// dataTrend handler
func dataTrend(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataTrendReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataTrend"))
}

// dataAction handler
func dataAction(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataActionReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataAction"))
}

// dataFan handler
func dataFan(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataFanReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataFan"))
}

// dataPandect handler
func dataPandect(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataPandectReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataPandect"))
}

// dataSurvey handler
func dataSurvey(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataSurveyReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataSurvey"))
}

// dataPlaySource handler
func dataPlaySource(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataPlaySourceReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataPlaySource"))
}

// dataPlayAnalysis handler
func dataPlayAnalysis(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataPlayAnalysisReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataPlayAnalysis"))
}

// dataArticleRank handler
func dataArticleRank(c *blademaster.Context) {
	var err error
	v := new(mcnmodel.DataArticleRankReq)
	if err = c.Bind(v); err != nil {
		return
	}
	mcnMid, _ := c.Get("mid")
	common := &mcnmodel.CreativeCommonReq{
		UpMid:  v.UpMid,
		McnMid: mcnMid.(int64),
	}
	c.JSON(srv.CreativeHandle(c, common, c.Request.Form, "dataArticleRank"))
}
