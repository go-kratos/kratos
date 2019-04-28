package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

//DoPtest do post
func (s *Service) DoPtest(c context.Context, ptestParam model.DoPtestParam) (resp model.DoPtestResp, err error) {
	var (
		testNameNick  string
		testNameNicks []string
	)
	if ptestParam.TestNameNick == "" {
		tim := strconv.FormatInt(time.Now().Unix(), 10)
		for _, testName := range ptestParam.TestNames {
			testNameNick = testName + tim
			testNameNicks = append(testNameNicks, testNameNick)
		}
	}
	return s.DoPtestByJmeter(c, ptestParam, testNameNicks)
}

//StopPtest stop test
func (s *Service) StopPtest(c context.Context, ptestJob model.PtestJob) (err error) {
	var (
		jobNames  []string
		ptestJobs []*model.PtestJob
	)
	ptestJob.Active = 1
	if ptestJobs, err = s.dao.QueryPtestJob(&ptestJob); err != nil {
		log.Error("s.dao.QueryPtestJob err :(%v)", err)
		return
	}
	for _, ptestJob := range ptestJobs {
		jobNames = append(jobNames, ptestJob.JobName)
	}
	for _, jobName := range jobNames {
		s.DeleteJob(context.TODO(), jobName)
	}
	return
}

//ReducePtest reduce ptest
func (s *Service) ReducePtest(c context.Context, reducePtest model.ReducePtest) (message string, err error) {
	if reducePtest.JobName == "" {
		message = "请输入 jobName"
		return
	}
	if _, err = s.DeleteJob(c, reducePtest.JobName); err != nil {
		message = "调用删除 job 接口失败"
		log.Error("s.DeleteJob err :(%v)", err)
		return
	}
	if err = s.dao.DeletePtestJob(reducePtest.ID); err != nil {
		message = "sql 执行失败"
		log.Error("s.dao.DeletePtestJob err :(%v)", err)
		return
	}
	message = "success"
	return
}

//QueryAllJobFree Query AllJob Free
func (s *Service) QueryAllJobFree(c context.Context, ptesJob *model.PtestJob) ([]*model.PtestJob, error) {
	return s.dao.QueryPtestJob(ptesJob)
}

//QueryOrStopAllPtestByJobName query or stop all test by job name
func (s *Service) QueryOrStopAllPtestByJobName(c context.Context, reportSuID int, IsDelete bool, testStatus int) (ptestjobsd []*model.PtestJob, err error) {

	// 删除所有正在允许此接口的容器
	ptestJobd := model.PtestJob{ReportSuID: reportSuID, Active: 1}
	if ptestjobsd, err = s.dao.QueryPtestJob(&ptestJobd); err != nil {
		log.Error("s.dao.QueryPtestJob err (%v)", err)
		return
	}

	if IsDelete {
		for _, ptestJob := range ptestjobsd {
			if _, err = s.DeleteJob(context.TODO(), ptestJob.JobName); err != nil {
				log.Error("s.DeleteJob err :(%v)", err)
				return
			}
			if err = s.dao.DeletePtestJob(ptestJob.ID); err != nil {
				log.Error("s.dao.DeletePtestJob err (%v)", err)
				return
			}
			if err = s.dao.UpdateReportStatusByID(reportSuID, testStatus); err != nil {
				log.Error("s.UpdateReportStatusByID err :(%v)", err)
				return
			}
		}
	}
	return
}

//DoAddPtest doadd test stress
func (s *Service) DoAddPtest(c context.Context, ptestAdd model.PtestAdd) (status string, err error) {
	status = "fail"
	var (
		resp        model.DoPtestResp
		grpc        *model.GRPC
		scripts     []*model.Script
		resportSumm []*model.ReportSummary
	)
	ptestParam := model.DoPtestParam{}
	// 此处 ScriptType 取自 reportSummary , 0-http  1-grpc  2-scene
	if ptestAdd.ScriptType == model.PROTOCOL_GRPC {
		if grpc, err = s.QueryGrpcById(ptestAdd.ScriptID); err != nil {
			log.Error("Query GRPC err: (%v)", err)
			return status, err
		}
		ptestParam = model.DoPtestParam{
			UserName:   grpc.UpdateBy,                // 用户名
			LoadTime:   grpc.LoadTime,                //运行时间
			TestNames:  StringToSlice(grpc.TaskName), //接口名转数组
			FileName:   grpc.JmxPath,                 // jmx文件
			ResJtl:     ptestAdd.ResJtl,              // jtl时间戳
			JmeterLog:  ptestAdd.JmeterLog,           // jmeterlog时间戳
			Department: grpc.Department,
			Project:    grpc.Project,
			IsDebug:    false,
			APP:        grpc.APP,
			ScriptID:   grpc.ID,
			Cookie:     "",               // 用不到
			URL:        grpc.ServiceName, // 用于微信通知
			LabelIDs:   nil,
			Domain:     grpc.HostName, // 微信通知Domain
			FileSplit:  false,         // 文件切割
			SplitNum:   0,             // 切割数量
			JarPath:    grpc.JarPath,
			Type:       model.PROTOCOL_GRPC, //grpc
			AddPtest:   true,                // 加压
		}
	}
	if ptestAdd.ScriptType == model.PROTOCOL_SCENE {
		ptestScene := model.DoPtestSceneParam{
			SceneID:  ptestAdd.SceneId,
			UserName: ptestAdd.UserName,
		}
		if resp, err = s.DoScenePtest(c, ptestScene, true, ""); err != nil {
			log.Error("s.DoScenePtest err :(%v)", err)
			return status, err
		}
	}
	if ptestAdd.ScriptType == model.PROTOCOL_HTTP {
		script := model.Script{ID: ptestAdd.ScriptID}
		if scripts, err = s.QueryScripts(&script, 1, 5); err != nil {
			log.Error("QueryScripts err :(%v)", err)
			return status, err
		}
		if len(scripts) > 0 {
			// http 模板
			ptestParam = model.DoPtestParam{
				UserName:    script.UpdateBy,
				TestNames:   QueryTestNamesByJmfile(scripts[0].SavePath),
				FileName:    scripts[0].SavePath,
				LoadTime:    1800,
				Upload:      scripts[0].Upload,
				ProjectName: scripts[0].ProjectName,
				ResJtl:      ptestAdd.ResJtl,
				JmeterLog:   ptestAdd.JmeterLog,
				Department:  scripts[0].Department,
				Project:     scripts[0].Project,
				APP:         scripts[0].App,
				ScriptID:    scripts[0].ID,
				AddPtest:    true,
				Domain:      scripts[0].Domain,
				FileSplit:   scripts[0].FileSplit,
				SplitNum:    scripts[0].SplitNum,
				DockerSum:   ptestAdd.DockerSum,
				Type:        model.PROTOCOL_HTTP,
				APIHeader:   scripts[0].APIHeader,
			}
		}
	}
	if ptestAdd.ScriptType != model.PROTOCOL_SCENE {
		if resp, err = s.DoPtestByJmeter(c, ptestParam, nil); err != nil {
			log.Error("DoPtestByJmeter err: (%v)", err)
			return
		}
		log.Info("add ---jobName:(%s)", resp.JobName)
	}
	status = "success"

	// 更新reportSummary数据表
	reportSu := model.ReportSummary{ID: ptestAdd.ReportSuID}
	if resportSumm, err = s.dao.QueryReportSurys(&reportSu); err != nil {
		return
	}
	ptestJob := model.PtestJob{JobName: resp.JobName, ScriptID: ptestAdd.ScriptID, ReportSuID: ptestAdd.ReportSuID, Active: 1, ExecuteID: ptestAdd.ExecuteID, HostIP: resp.HostIP}
	log.Info("add ---jobName:(%s)", resp.JobName)
	if _, err = s.dao.AddPtestJob(&ptestJob); err != nil {
		log.Error("s.dao.AddPtestJob err (%v)", err)
		return
	}
	DockerSum := resportSumm[0].DockerSum + 1
	if err = s.dao.UpdateReportDockByID(resportSumm[0].ID, DockerSum); err != nil {
		log.Error("s.dao.UpdateReportDockByID err (%v)", err)
		return
	}
	return
}

//AddPtest add ptest
func (s *Service) AddPtest(c context.Context, ptestAdd model.PtestAdd) (err error) {
	go s.DoAddPtestWithSleep(c, ptestAdd)
	return
}

//DoAddPtestWithSleep Do AddPtestWith Sleep
func (s *Service) DoAddPtestWithSleep(c context.Context, ptestAdd model.PtestAdd) (err error) {
	var ptestJobs []*model.PtestJob
	if ptestAdd.DockerNum == 0 {
		ptestAdd.DockerNum = 1
	}
	if ptestAdd.SleepTime == 0 {
		ptestAdd.SleepTime = 1
	}
	for i := 1; i <= ptestAdd.DockerNum; i++ {
		time.Sleep(time.Duration(ptestAdd.SleepTime) * time.Second)

		// 判断主容器状态，主容器被删除，则不再添加新的容器
		ptestJob := model.PtestJob{JobName: ptestAdd.JobName}
		if ptestJobs, err = s.dao.QueryPtestJob(&ptestJob); err != nil {
			log.Error("s.dao.QueryPtestJob err :(%v)", err)
		}
		if len(ptestJobs) > 0 {
			if ptestJobs[0].Active != -1 {
				go s.DoAddPtest(context.TODO(), ptestAdd)
			}
		}
	}
	return
}

//DoPtestByFile do test by file
func (s *Service) DoPtestByFile(c context.Context, script model.Script, resJtl, jmeterLog string) (resp model.DoPtestResp, err error) {
	var scripts []*model.Script
	if scripts, err = s.QueryScripts(&script, 1, 5); err != nil {
		log.Error("QueryScripts err :(%v)", err)
		return
	}
	scriptID := scripts[0].ID
	testNames := QueryTestNamesByJmfile(scripts[0].SavePath)
	ptestParam := model.DoPtestParam{
		UserName:            script.UpdateBy,
		TestNames:           testNames,
		FileName:            scripts[0].SavePath,
		LoadTime:            1800,
		Upload:              scripts[0].Upload,
		ProjectName:         scripts[0].ProjectName,
		ResJtl:              resJtl,
		JmeterLog:           jmeterLog,
		Department:          scripts[0].Department,
		Project:             scripts[0].Project,
		APP:                 scripts[0].App,
		ScriptID:            scripts[0].ID,
		Fusing:              scripts[0].Fusing,
		UseBusinessStop:     scripts[0].UseBusinessStop,
		BusinessStopPercent: scripts[0].BusinessStopPercent,
	}
	if resp, err = s.DoPtest(context.TODO(), ptestParam); err != nil {
		log.Error("DoPtest error :(%v)", err)
		return
	}
	resp.ScriptID = scriptID
	return
}

//DoPtestArr QueryScriptsInID
func (s *Service) DoPtestArr(c context.Context, PtestBatch model.PtestBatch, cookie string) (status string, err error) {
	var scripts []*model.Script
	if scripts, err = s.dao.QueryScriptsInID(PtestBatch.IDArr); err != nil {
		status = "fail"
		log.Error("QueryScriptsInID err --- :(%v) ", err)
		return status, err
	}

	for _, script := range scripts {
		JmeterLog := script.JmeterLog + strconv.FormatInt(time.Now().Unix(), 10)
		ResJtl := script.ResJtl + strconv.FormatInt(time.Now().Unix(), 10)

		ptestParam := model.DoPtestParam{
			UserName:    PtestBatch.UserName,
			FileName:    script.SavePath,
			ProjectName: script.ProjectName,
			Upload:      script.Upload,
			ResJtl:      ResJtl,
			JmeterLog:   JmeterLog,
			Department:  script.Department,
			Project:     script.Project,
			APP:         script.App,
			ScriptID:    script.ID,
			URL:         script.URL,
			Cookie:      cookie,
			Domain:      script.Domain,
			FileSplit:   script.FileSplit,
			SplitNum:    script.SplitNum,
			Fusing:      script.Fusing,
			APIHeader:   script.APIHeader,
		}
		if script.Upload {
			ptestParam.LoadTime = 1800
			ptestParam.TestNames = QueryTestNamesByJmfile(script.SavePath)
		} else {
			ptestParam.LoadTime = script.LoadTime
			ptestParam.TestNames = StringToSlice(script.TestName)
		}
		go s.DoPtest(context.TODO(), ptestParam)
	}
	status = "success"
	return
}

//DoPtestByScriptId do ptest by scriptid
func (s *Service) DoPtestByScriptId(c context.Context, script *model.Script, cookie, executor string) (resp model.DoPtestResp, err error) {
	var scripts []*model.Script
	if scripts, err = s.QueryScripts(script, 1, 5); err != nil {
		log.Error("QueryScripts err :(%v)", err)
		return
	}
	ptestParam := model.DoPtestParam{
		TestNames:           StringToSlice(scripts[0].TestName),
		ProjectName:         scripts[0].ProjectName,
		FileName:            scripts[0].SavePath,
		LoadTime:            scripts[0].LoadTime,
		ResJtl:              scripts[0].ResJtl,
		JmeterLog:           scripts[0].JmeterLog,
		UserName:            executor,
		Department:          scripts[0].Department,
		Project:             scripts[0].Project,
		APP:                 scripts[0].App,
		ScriptID:            scripts[0].ID,
		URL:                 scripts[0].URL,
		Cookie:              cookie,
		Domain:              scripts[0].Domain,
		FileSplit:           scripts[0].FileSplit,
		SplitNum:            scripts[0].SplitNum,
		Fusing:              scripts[0].Fusing,
		APIHeader:           scripts[0].APIHeader,
		Upload:              scripts[0].Upload,
		UseBusinessStop:     scripts[0].UseBusinessStop,
		BusinessStopPercent: scripts[0].BusinessStopPercent,
	}

	if ptestParam.Upload {
		ptestParam.TestNames = QueryTestNamesByJmfile(scripts[0].SavePath)
	}
	return s.DoPtest(c, ptestParam)
}

//DoPtestByJmeter do ptest by jmeter
func (s *Service) DoPtestByJmeter(c context.Context, ptestParam model.DoPtestParam, testNameNicks []string) (resp model.DoPtestResp, err error) {
	var (
		du            string
		Debug         int
		CPUCore       int
		reportSuID    int
		scriptSnapIDs []int
		beginTime     string
		token         string
		Index         = 0
		btm           time.Time
		dus           time.Duration
		tempRes       *model.PaasJobQueryStatus
		hostIP        string
		command       string
	)
	lay := "2006-01-02 15:04:05"
	executeID := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	if !ptestParam.AddPtest {
		// 不加压，增加新的 log
		ptestParam.ResJtl = fmt.Sprintf("%s%s", ptestParam.ResJtl, executeID)
		ptestParam.JmeterLog = fmt.Sprintf("%s%s", ptestParam.JmeterLog, executeID)
	}
	jobNamed := executeID + "run"
	jobName := string([]rune(jobNamed)[10:])

	// add filesplit messages to env variables
	if !ptestParam.AddPtest {
		Index = 1
	} else {
		Index = ptestParam.DockerSum + 1
	}
	ptestParam.EnvInfo = "FileSplit:" + strconv.FormatBool(ptestParam.FileSplit) + ",SplitNum:" + strconv.Itoa(ptestParam.SplitNum) +
		",Index:" + strconv.Itoa(Index)

	// 为压测容器绑定host
	if ptestParam.Upload {
		HostInfoByUploadSc(ptestParam.Domain)
	} else {
		if ptestParam.Type == model.PROTOCOL_HTTP {
			HostInfo(ptestParam.Domain, ptestParam.APIHeader)
		} else if ptestParam.Type == model.PROTOCOL_SCENE {
			HostInfoList(ptestParam.Scripts)
		}
	}

	if _, _, err = CreateResJtlAndJmeterLog(ptestParam); err != nil {
		log.Info("create resjtl or jmeterlog err :(%v)", err)
		return
	}
	//获取 Debug,CPUCore,command
	Debug, CPUCore, command = s.CreateCommand(ptestParam)

	jobInfo := model.Job{FileName: ptestParam.FileName, CPU: CPUCore, Memory: 8096, Parallelism: 1, Name: jobName, ResJtl: ptestParam.ResJtl,
		JmeterLog: ptestParam.JmeterLog, EnvInfo: ptestParam.EnvInfo, JarPath: ptestParam.JarPath, Command: command}

	beginTime = time.Now().Format(lay)
	if btm, err = time.Parse(lay, beginTime); err != nil {
		return
	}
	resd := int(btm.Unix()) % 5
	du = "-" + strconv.Itoa(resd) + "s"
	if dus, err = time.ParseDuration(du); err != nil {
		return
	}

	beginTime = btm.Add(dus).Format(lay)
	res := make(chan interface{}, 10000)
	resHostIp := make(chan interface{}, 10000)
	tm := time.Now()
	for i := 0; i <= ptestParam.LoadTime+10; i += 5 {
		tm = tm.Add(time.Second * 5)
	}
	timeout := time.After(time.Until(tm))
	log.Info("job set timeout util --- :(%s) ", time.Until(tm))

	if ptestParam.ExecuDockerSum == 0 {
		ptestParam.ExecuDockerSum = 1
	}
	if _, err = s.AddJob(context.TODO(), jobInfo); err != nil {
		log.Error("add ptest job err:(%v)", err)
		return
	}
	if token, err = s.RmToken(c); err != nil {
		log.Error("get token err:(%v)", err)
		return
	}
	//获取物理机ip
	if tempRes, err = s.Job(c, jobName); err != nil {
		hostIP = ""
	}
	if len(tempRes.Data.Pods) > 0 {
		hostIP = tempRes.Data.Pods[0].HostIP
	}

	// 如果不是增加容器，生成一份快照
	if !ptestParam.AddPtest {
		if scriptSnapIDs, err = s.AddSnap(c, ptestParam, executeID, jobName, jobNamed); err != nil || len(scriptSnapIDs) == 0 {
			s.DeleteJob(context.TODO(), jobName)
			return
		}
		// 从页面上直接输入参数生成的单场景脚本，testNames 和 testNameNicks 长度都是 1，reportSummary 测试计划为空
		if reportSuID, err = s.AddReSummaryByPtest(ptestParam, jobName, jobNamed, testNameNicks, scriptSnapIDs[0], Debug); err != nil {
			log.Error("s.dao.AddReportSummary err :(%v)", err)
			s.DeleteJob(context.TODO(), jobName)
			return
		}
		//容器写入 ptest_job 表
		ptestJob := model.PtestJob{ScriptID: ptestParam.ScriptID, ReportSuID: reportSuID, JobName: jobName, Active: 1, ExecuteID: jobNamed, HostIP: hostIP}
		if _, err = s.dao.AddPtestJob(&ptestJob); err != nil {
			log.Error("s.dao.AddPtestJob err :(%v)", err)
			s.DeleteJob(context.TODO(), jobName)
			return
		}
		//获取场景中每个接口的熔断成功率的值
		fusingList := make([]int, len(ptestParam.Scripts))
		useBusiStopList := make([]bool, len(ptestParam.Scripts))
		busiStopPercentList := make([]int, len(ptestParam.Scripts))
		for index, script := range ptestParam.Scripts {
			fusingList[index] = script.Fusing
			useBusiStopList[index] = script.UseBusinessStop
			busiStopPercentList[index] = script.BusinessStopPercent
		}
		AddReGraphTimer := model.AddReGraphTimer{
			ScriptID:            ptestParam.ScriptID,
			JobName:             jobName,
			BeginTime:           beginTime,
			Token:               token,
			TestNames:           ptestParam.TestNames,
			TestNameNicks:       testNameNicks,
			ReportSuID:          reportSuID,
			Fusing:              ptestParam.Fusing,
			FusingList:          fusingList,
			TestType:            ptestParam.Type,
			UseBusinessStop:     ptestParam.UseBusinessStop,
			BusinessStopPercent: ptestParam.BusinessStopPercent,
			UseBusiStopList:     useBusiStopList,
			BusiStopPercentList: busiStopPercentList,
		}
		//场景脚本的逻辑，可能存在多个test_name
		if ptestParam.Type == model.PROTOCOL_SCENE && ptestParam.TestNameNick != "" {
			AddReGraphTimer.TestNameNicks = ptestParam.TestNameNicks
			AddReGraphTimer.TestNames = ptestParam.TestNames
		}
		//数据同步job
		log.Info("-------开始执行数据同步job,loadTime :(%v)", ptestParam.LoadTime)
		go s.addReGraphTimer(context.TODO(), AddReGraphTimer, res, timeout)

		//查询容器hostIp
		go s.addPtestJobHostIp(context.TODO(), AddReGraphTimer, resHostIp, timeout)
		//增加label
		if len(ptestParam.LabelIDs) > 0 && !ptestParam.IsDebug {
			for _, lableID := range ptestParam.LabelIDs {
				lrl := model.LabelRelation{LabelID: int64(lableID), Type: model.ScriptType, TargetID: int64(ptestParam.ScriptID)}
				s.AddLabelRelation(&lrl)
				lrll := model.LabelRelation{LabelID: int64(lableID), Type: model.ReportType, TargetID: int64(reportSuID)}
				s.AddLabelRelation(&lrll)
			}
		}
		//如果不是debug ，且配置发生通知为true，则发送微信通知
		if !ptestParam.IsDebug && conf.Conf.Wechat.SendMessage {
			serviceName := ptestParam.Department + "." + ptestParam.Project + "." + ptestParam.APP
			var userService map[string][]string
			if userService, err = s.QueryDependServiceAdmins(c, serviceName, s.getSessionInCookie(ptestParam.Cookie)); err != nil {
				log.Error("query depend service admin error(%v)", err)
				return
			}
			// 发送群通知
			content := AddWechatContent(ptestParam, reportSuID, jobName, userService)
			go s.AddWechatSend(context.TODO(), ptestParam.Cookie, content)

			// 发送依赖服务通知
			for user := range userService {
				toUser := StringToSlice(user)
				content = AddWechatDependServiceContent(ptestParam, userService, reportSuID, user)
				go s.dao.PushWechatMsgToPerson(context.TODO(), ptestParam.Cookie, toUser, content)
			}
		}
	}
	resp = model.DoPtestResp{BeginTime: beginTime, ReportSuID: reportSuID, JobName: jobName, ScriptSnapIDs: scriptSnapIDs,
		JmeterLog: ptestParam.JmeterLog, JtlLog: ptestParam.ResJtl, JmxFile: ptestParam.FileName, LoadTime: ptestParam.LoadTime, HostIP: hostIP}
	return
}

//QueryTestNameNick query test nick name
func QueryTestNameNick(TestName string, testNameNicks []string) (testNameNick string) {
	for _, testNameNickd := range testNameNicks {
		if strings.Contains(testNameNickd, TestName) {
			testNameNick = testNameNickd
			break
		}
	}
	return testNameNick
}

//IsUniqObject check whether object is unique
func IsUniqObject(aa []string) (t bool) {
	for _, a := range aa {
		if a == aa[0] {
			t = true
			continue
		}

		if a != aa[0] {
			t = false
			break
		}
	}
	return t
}

//StringToSlice convert string to string slice
func StringToSlice(str string) (strs []string) {
	strs = append(strs, str)
	return
}

//RemoveRepByMap 去除切片钟的重复元素
func RemoveRepByMap(slc []string) (result []string) {
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

//SliceToString join slice elements to string
func SliceToString(strs []string, inter string) (result string) {
	for _, str := range strs {
		result = result + inter + str
	}
	return string([]rune(result)[1:])
}

//QueryTestNamesByJmfile query test names by jmeter file
func QueryTestNamesByJmfile(fileName string) (testNames []string) {
	var (
		buff []byte
		err  error
	)
	if buff, err = ioutil.ReadFile(fileName); err != nil {
		log.Error("open script file failed! error %v", err)
		return
	}
	reg := regexp.MustCompile("SampleGui\"(.*)\" enabled=\"true")
	reg2 := regexp.MustCompile("testname=\"(.*)\" enabled=\"true")
	results := reg.FindAllString(string(buff), -1)
	for _, res := range results {
		resd := reg2.FindString(res)
		resdd := strings.Replace(resd, "\" enabled=\"true", "", -1)
		resddd := strings.Replace(resdd, "testname=\"", "", -1)
		testNames = append(testNames, resddd)
	}
	return
}

//ExistsInSlice check whether the string in the slice
func ExistsInSlice(str string, strs []string) bool {
	for _, st := range strs {
		if st == str {
			return true
		}
	}
	return false
}

//PingDomain ping domain
func PingDomain(ptestParam model.DoPtestParam) (pingString string) {
	var pingStrings []string
	if ptestParam.Type == model.PROTOCOL_HTTP {
		pingString = " ping -c 1 " + ptestParam.Domain + " |tee -a " + ptestParam.JmeterLog
	} else if ptestParam.Type == model.PROTOCOL_SCENE {
		for _, script := range ptestParam.Scripts {
			pingStrings = append(pingStrings, " ping -c 1 "+script.Domain+" |tee -a "+ptestParam.JmeterLog)
		}
		pingString = SliceToString(pingStrings, "&")
	}
	return
}

//CreateResJtlAndJmeterLog Create ResJt lAnd JmeterLog
func CreateResJtlAndJmeterLog(ptestParam model.DoPtestParam) (ResJtlFile, JmeterLogFile *os.File, err error) {
	if ptestParam.ResJtl != "" {
		if ResJtlFile, err = os.Create(ptestParam.ResJtl); err != nil {
			log.Error("create ResJtl error :(%v)", err)
			return
		}
		ResJtlFile.WriteString("此处显示 err 日志")
		defer ResJtlFile.Close()
	}
	if ptestParam.JmeterLog != "" {
		if JmeterLogFile, err = os.Create(ptestParam.JmeterLog); err != nil {
			log.Error("create JmeterLog error :(%v)", err)
			return
		}
		JmeterLogFile.WriteString("此处显示启动日志，执行debug 则显示debug 日志")
		defer JmeterLogFile.Close()
	}
	return
}

//CreateCommand Create Command
func (s *Service) CreateCommand(ptestParam model.DoPtestParam) (debug, CPUCore int, command string) {
	// ping 所有接口的 domain
	pingString := PingDomain(ptestParam)

	cpJar := ""
	if ptestParam.JarPath != "" {
		cpJar = fmt.Sprintf("cp %s %s & ", ptestParam.JarPath, s.c.Jmeter.JmeterExtLibPathContainer)
	}

	// 调试逻辑
	if ptestParam.IsDebug {
		CPUCore = s.c.Paas.CPUCoreDebug
		debug = 1
		command = cpJar + " mkdir /data/jmeter-log & jmeter -n -t " + ptestParam.FileName + " -j " + ptestParam.JmeterLog + " -l " + ptestParam.ResJtl + " -F ;" + pingString
	} else {
		CPUCore = s.c.Paas.CPUCore
		debug = -1
		command = cpJar + " mkdir /data/jmeter-log & jmeter -n -t " + ptestParam.FileName + " -j " + ptestParam.JmeterLog + " -l " + ptestParam.ResJtl + " -F"
	}
	return
}
