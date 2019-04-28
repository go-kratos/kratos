package http

import (
	"encoding/json"

	"go-common/app/admin/ep/melloi/model"
	"go-common/app/admin/ep/melloi/service"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

//excute performTest
func doPtest(c *bm.Context) {
	var (
		script     model.Script
		ptestParam model.DoPtestParam
	)
	userName := c.Request.Form.Get("userName")
	param := c.Request.Form.Get("param")
	resJtl := c.Request.Form.Get("resJtl")
	jmeterLog := c.Request.Form.Get("jmeterLog")
	cookie := c.Request.Header.Get("Cookie")
	if err := json.Unmarshal([]byte(param), &script); err != nil {
		log.Error("sorry,script err(%v)", err)
		c.JSON(nil, err)
		return
	}
	log.Info("接口名：(%s),压测脚本路径：(%s)", script.TestName, script.SavePath)
	if script.Upload {
		testNames := service.QueryTestNamesByJmfile(script.SavePath)
		ptestParam = model.DoPtestParam{
			TestNames:   testNames,
			FileName:    script.SavePath,
			LoadTime:    1800,
			UserName:    userName,
			Upload:      true,
			ProjectName: script.ProjectName,
			FileSplit:   script.FileSplit,
			SplitNum:    script.SplitNum,
			APIHeader:   script.APIHeader,
		}
	} else {
		ptestParam = model.DoPtestParam{
			TestNames: service.StringToSlice(script.TestName),
			FileName:  script.SavePath,
			LoadTime:  script.LoadTime,
			UserName:  userName,
			URL:       script.URL,
			FileSplit: script.FileSplit,
			SplitNum:  script.SplitNum,
			APIHeader: script.APIHeader,
		}
	}
	ptestParam.Department = script.Department
	ptestParam.Project = script.Project
	ptestParam.APP = script.App
	ptestParam.Domain = script.Domain
	ptestParam.Cookie = cookie
	ptestParam.JmeterLog = jmeterLog
	ptestParam.ResJtl = resJtl
	ptestParam.ScriptID = script.ID
	ptestParam.UseBusinessStop = script.UseBusinessStop
	ptestParam.BusinessStopPercent = script.BusinessStopPercent
	ptestParam.Fusing = script.Fusing
	c.JSON(srv.DoPtest(c, ptestParam))
}

//addPtest  加压
func addPtest(c *bm.Context) {
	var ptestAdd model.PtestAdd
	if err := c.BindWith(&ptestAdd, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	log.Info("AddPtest ----")
	c.JSON(nil, srv.AddPtest(c, ptestAdd))
}

//stopPtest
func stopPtest(c *bm.Context) {
	var ptestJob model.PtestJob
	if err := c.BindWith(&ptestJob, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, srv.StopPtest(c, ptestJob))
}

//stopAllPtest
func stopAllPtest(c *bm.Context) {
	var allPtestStop model.AllPtestStop
	if err := c.BindWith(&allPtestStop, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryOrStopAllPtestByJobName(c, allPtestStop.ReportSuID, true, 3))
}

//reducePtest
func reducePtest(c *bm.Context) {
	var (
		reducePtest model.ReducePtest
		resMap      = make(map[string]interface{})
	)
	if err := c.BindWith(&reducePtest, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	message, err := srv.ReducePtest(c, reducePtest)
	if err != nil {
		return
	}
	resMap["message"] = message
	c.JSON(resMap, err)
}

//queryAllJob
func queryAllJob(c *bm.Context) {
	var (
		ptestJob model.PtestJob
		resMap   = make(map[string]interface{})
	)
	if err := c.BindWith(&ptestJob, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	ptestJobs, err := srv.QueryOrStopAllPtestByJobName(c, ptestJob.ReportSuID, false, 3)
	if err != nil {
		log.Error("query ptestJob err", err)
		return
	}
	resMap["ptest_jobs"] = ptestJobs
	c.JSON(resMap, err)
}

//queryAllJobFree queryAll JobFree
func queryAllJobFree(c *bm.Context) {
	var (
		ptestJob model.PtestJob
		resMap   = make(map[string]interface{})
	)
	if err := c.BindWith(&ptestJob, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	ptestJobs, err := srv.QueryAllJobFree(c, &ptestJob)
	if err != nil {
		log.Error("query ptestJob err", err)
		return
	}
	resMap["ptest_jobs"] = ptestJobs
	c.JSON(resMap, err)
}

//doDebug
func doDebug(c *bm.Context) {
	var script model.Script
	var ptestParam model.DoPtestParam
	param := c.Request.Form.Get("param")
	userName := c.Request.Form.Get("userName")
	if err := json.Unmarshal([]byte(param), &script); err != nil {
		log.Error("sorry,script err(%v)", err)
		c.JSON(nil, err)
		return
	}
	log.Info("接口名：(%s),压测脚本路径：(%s)", script.TestName, script.SavePath)
	if script.Upload {
		ptestParam = model.DoPtestParam{
			TestNames:   service.QueryTestNamesByJmfile(script.SavePath),
			FileName:    script.SavePath,
			LoadTime:    1800,
			UserName:    userName,
			Upload:      true,
			ProjectName: script.ProjectName,
			ResJtl:      script.ResJtl,
			JmeterLog:   script.JmeterLog,
			Department:  script.Department,
			Project:     script.Project,
			APP:         script.App,
			ScriptID:    script.ID,
		}
	} else {
		ptestParam = model.DoPtestParam{
			TestNames:  service.StringToSlice(script.TestName),
			FileName:   script.SavePath,
			LoadTime:   script.LoadTime,
			UserName:   userName,
			ResJtl:     script.ResJtl,
			JmeterLog:  script.JmeterLog,
			Department: script.Department,
			Project:    script.Project,
			APP:        script.App,
			ScriptID:   script.ID,
			Domain:     script.Domain,
		}
	}
	c.JSON(srv.DoPtest(c, ptestParam))
}

func doPtestByFile(c *bm.Context) {
	var script model.Script
	if err := c.BindWith(&script, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	resJtl := script.ResJtl
	jmeterLog := script.JmeterLog
	script.ResJtl = ""
	script.JmeterLog = ""
	resp, err := srv.DoPtestByFile(c, script, resJtl, jmeterLog)
	if err != nil {
		c.JSON(err, nil)
		return
	}
	c.JSON(resp, err)
}

//excute batch ptest
func doPtestBatch(c *bm.Context) {
	var PtestBatch model.PtestBatch
	if err := c.BindWith(&PtestBatch, binding.JSON); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	status, _ := srv.DoPtestArr(c, PtestBatch, cookie)
	var doPostMap = make(map[string]string)
	doPostMap["status"] = status
	c.JSON(doPostMap, nil)
}

//doPtestByScriptId
func doPtestByScriptId(c *bm.Context) {
	var script model.Script
	if err := c.BindWith(&script, binding.Form); err != nil {
		c.JSON(nil, err)
		return
	}
	cookie := c.Request.Header.Get("Cookie")
	executor := c.Request.Form.Get("executor")
	c.JSON(srv.DoPtestByScriptId(c, &script, cookie, executor))
}
