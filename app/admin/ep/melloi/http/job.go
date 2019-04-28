package http

import (
	"go-common/app/admin/ep/melloi/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func addJob(c *bm.Context) {
	job := model.Job{}
	if err := c.BindWith(&job, binding.JSON); nil != err {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.AddJob(c, job))
}
func deleteJob(c *bm.Context) {
	name := c.Request.Form.Get("name")
	c.JSON(srv.DeleteJob(c, name))
}

func deleteJobBatch(c *bm.Context) {
	JobBatch := model.JobBatch{}
	if err := c.BindWith(&JobBatch, binding.JSON); nil != err {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.DeleteJobBatch(c, JobBatch), nil)
}

//Job get job info
func Job(c *bm.Context) {
	name := c.Request.Form.Get("name")
	c.JSON(srv.Job(c, name))
}

func forceDelete(c *bm.Context) {
	ptestJob := model.PtestJob{}
	if err := c.BindWith(&ptestJob, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.ForceDeleteJob(c, ptestJob.ReportSuID))
}

func queryClearnableDocker(c *bm.Context) {
	c.JSON(srv.QueryCleanableDocker(c))
}

func cleanNotRunningJob(c *bm.Context) {
	c.JSON(nil, srv.CleanJob(c))
}
