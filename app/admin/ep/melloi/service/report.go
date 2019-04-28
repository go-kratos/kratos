package service

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

// QueryReportSummarys query report summary
func (s *Service) QueryReportSummarys(c context.Context, sessionID string, qrsr *model.QueryReportSuRequest) (qrs *model.QueryReportSuResponse, err error) {
	// 获取服务树节点
	var treeNodes, treeNodesd []string
	if treeNodesd, err = s.QueryUserRoleNode(c, sessionID); err != nil {
		log.Error("QueryUserRoleNode  (%v): ", err)
		return
	}
	treeNodes = append(treeNodesd, "")
	if ExistsInSlice(qrsr.Executor, conf.Conf.Melloi.Executor) {
		if qrs, err = s.dao.QueryReportSummarysWhiteName(&qrsr.ReportSummary, qrsr.SearchAll, qrsr.PageNum, qrsr.PageSize); err != nil {
			return
		}
	} else {
		if qrs, err = s.dao.QueryReportSummarys(&qrsr.ReportSummary, qrsr.SearchAll, qrsr.PageNum, qrsr.PageSize, treeNodes); err != nil {
			return
		}
	}

	// 获取label
	for _, report := range qrs.ReportSummarys {
		lr := model.LabelRelation{Type: model.ReportType, TargetID: int64(report.ID)}
		if report.Labels, err = s.dao.QueryLabelRelation(&lr); err != nil {
			return
		}
	}
	return
}

// QueryReportByID Query Report By ID
func (s *Service) QueryReportByID(id int) (summary *model.ReportSummary, err error) {
	return s.dao.QueryReportSuryByID(id)
}

// CountQueryReportSummarys count query report summarys
func (s *Service) CountQueryReportSummarys(reportSummary *model.ReportSummary) (int, error) {
	return s.dao.CountQueryReportSummarys(reportSummary)
}

// UpdateReportSummary update report summary
func (s *Service) UpdateReportSummary(reportSummary *model.ReportSummary) (status string, err error) {
	if err = s.dao.UpdateReportSummary(reportSummary); err != nil {
		status = "update fail"
		return
	}
	status = "success"
	return
}

// UpdateReportStatus update report summary status
func (s *Service) UpdateReportStatus(status int) error {
	return s.dao.UpdateReportStatus(status)
}

// AddReSummaryByPtest Add ReSummary By Ptest
func (s *Service) AddReSummaryByPtest(ptestParam model.DoPtestParam, jobName, executeID string, testNameNicks []string, snapID, Debug int) (reportSuID int, err error) {
	reportSummary := model.ReportSummary{
		JobName:      jobName,
		TestStatus:   2,
		UserName:     ptestParam.UserName,
		ResJtl:       ptestParam.ResJtl,
		JmeterLog:    ptestParam.JmeterLog,
		Department:   ptestParam.Department,
		Project:      ptestParam.Project,
		APP:          ptestParam.APP,
		ScriptID:     ptestParam.ScriptID,
		DockerSum:    1,
		Active:       1,
		Debug:        Debug,
		ScriptSnapID: snapID,
		ExecuteID:    executeID,
		SceneID:      ptestParam.SceneID,
		Type:         ptestParam.Type,
		LoadTime:     ptestParam.LoadTime,
	}
	//Upload==true，这里表示从前端上传的脚本，testNames 和 testNameNicks 长度可能大于1，需要遍历，写入 reportSummary
	if ptestParam.Upload {
		reportSummary.TestNameNick = SliceToString(testNameNicks, ",")
		reportSummary.TestName = ptestParam.ProjectName
		reportSummary.Type = model.PROTOCOL_SCENE
	}

	// 从页面上直接输入参数生成的单场景脚本，testNames 和 testNameNicks 长度都是 1
	// 生成的场景脚本，ptestParam.TestNameNick 不为空
	if !ptestParam.Upload {
		for _, testName := range ptestParam.TestNames {
			reportSummary.TestNameNick = QueryTestNameNick(testName, testNameNicks)
			reportSummary.TestName = testName
		}
		//场景脚本的 reportSummary 逻辑
		if ptestParam.Type == model.PROTOCOL_SCENE && ptestParam.TestNameNick != "" {
			reportSummary.TestNameNick = ptestParam.TestNameNick
			reportSummary.TestName = ptestParam.SceneName
		}
	}
	if reportSuID, err = s.dao.AddReportSummary(&reportSummary); err != nil {
		log.Error("s.dao.AddReportSummary err :(%v)", err)
		return
	}
	return
}

// AddReportSummary add report summary
func (s *Service) AddReportSummary(reportSummary *model.ReportSummary) (status string, reportSuID int, err error) {
	var total int
	if total, err = s.dao.CountQueryReportSummarys(reportSummary); err != nil {
		log.Error("CountQueryReportSummarys error:(%v)", err)
		status = "fail"
		return
	}
	if (total == 1) || (reportSummary.TestName == "") {
		status = "数据已存在或者无接口名"
	} else {
		if reportSuID, err = s.dao.AddReportSummary(reportSummary); err != nil {
			status = "fail"
			return
		}
		status = "success"
	}
	return
}

// QueryReGraph query reGraph
func (s *Service) QueryReGraph(testNameNicks []string) (reportGraphssd [][]model.ReportGraph, err error) {
	var (
		TestNames     []string
		reportGraphss []model.ReportGraph
		reportGraphs  []model.ReportGraph
	)
	if reportGraphs, err = s.dao.QueryReportGraph(testNameNicks); err != nil {
		log.Error("dao.QueryReportGraph error:(%v)", err)
		return
	}
	if len(reportGraphs) > 0 {
		for _, reportGraph := range reportGraphs {
			TestNames = append(TestNames, reportGraph.TestName)
		}
		TestNamesd := RemoveRepByMap(TestNames)
		for _, testName := range TestNamesd {
			for _, reportGraph := range reportGraphs {
				if reportGraph.TestName == testName {
					reportGraphss = append(reportGraphss, reportGraph)
				}
			}
			reportGraphssd = append(reportGraphssd, reportGraphss)
			reportGraphss = []model.ReportGraph{}
		}
	}
	return
}

// QueryClientMoni query client moni
func (s *Service) QueryClientMoni(cli *model.ClientMoni) (clientMonisd [][]*model.ClientMoni, err error) {
	var (
		jobNames     []string
		clientMoniss []*model.ClientMoni
		clientMonis  []*model.ClientMoni
	)
	if clientMonis, err = s.dao.QueryClientMoni(cli); err != nil {
		log.Error("dao.QueryClientMoni error:(%v)", err)
		return
	}
	if len(clientMonis) > 0 {
		for _, clientMoni := range clientMonis {
			jobNames = append(jobNames, clientMoni.JobName)
		}
		jobNamesd := RemoveRepByMap(jobNames)
		for _, jobNamed := range jobNamesd {
			for _, clientMoni := range clientMonis {
				if clientMoni.JobName == jobNamed {
					clientMoniss = append(clientMoniss, clientMoni)
				}
			}
			clientMonisd = append(clientMonisd, clientMoniss)
			clientMoniss = []*model.ClientMoni{}
		}
	}
	return
}

//QueryReGraphAvg query reGraph
func (s *Service) QueryReGraphAvg(testNameNicks []string) (reportGraphAvgs []model.ReportGraph, err error) {
	var reportGraphssd [][]model.ReportGraph
	if reportGraphssd, err = s.QueryReGraph(testNameNicks); err != nil {
		log.Error("QueryReGraph error :(%s)", err)
		return
	}
	for _, reportGraphsd := range reportGraphssd {
		reportGraph := ReportGraphAvg(reportGraphsd)
		reportGraphAvgs = append(reportGraphAvgs, reportGraph)
	}
	return
}

//ReportGraphAvg report graph avg
func ReportGraphAvg(reportGraphs []model.ReportGraph) (reportGraph model.ReportGraph) {
	var sumTime, sumQPS, sumNetIo, min, max, sumCount, sumError, sumCodeEll, sumCodeWll, sumCodeWly, sumCodeWle, sumCodeWls,
		sumCodeSll, sumCodeSly, sumCodeSls, sumCodeKong, sumCodeNonHTTP, sumCodeOthers, fiftyTime, ninetyTime, ninetyFiveTime,
		ninetyNineTime, sumCode301, sumCode302, timeCount, sumQPSRecent, QPSRecent int
	var mins, maxs []int

	if len(reportGraphs) > 0 {
		for _, reportGraph := range reportGraphs {
			mins = append(mins, reportGraph.Min)
			maxs = append(maxs, reportGraph.Max)
			sumCount = sumCount + reportGraph.Count
			timeCount = reportGraph.Count * reportGraph.AvgTime
			sumTime = sumTime + timeCount
			sumError = sumError + reportGraph.Error
			sumQPS = sumQPS + reportGraph.QPS
			sumNetIo = sumNetIo + reportGraph.NetIo
			sumCodeEll = sumCodeEll + reportGraph.CodeEll
			sumCodeWll = sumCodeWll + reportGraph.CodeWll
			sumCodeWly = sumCodeWly + reportGraph.CodeWly
			sumCodeWle = sumCodeWle + reportGraph.CodeWle
			sumCodeWls = sumCodeWls + reportGraph.CodeWls
			sumCodeSll = sumCodeSll + reportGraph.CodeSll
			sumCodeSly = sumCodeSly + reportGraph.CodeSly
			sumCodeSls = sumCodeSls + reportGraph.CodeSls
			sumCodeKong = sumCodeKong + reportGraph.CodeKong
			sumCode301 = sumCode301 + reportGraph.Code301
			sumCode302 = sumCode302 + reportGraph.Code302
			sumCodeNonHTTP = sumCodeNonHTTP + reportGraph.CodeNonHTTP
			sumCodeOthers = sumCodeOthers + reportGraph.CodeOthers
			fiftyTime = fiftyTime + reportGraph.FiftyTime
			ninetyTime = ninetyTime + reportGraph.NinetyTime
			ninetyFiveTime = ninetyFiveTime + reportGraph.NinetyFiveTime
			ninetyNineTime = ninetyNineTime + reportGraph.NinetyNineTime
		}
		// 取最后 recent次统计的平均 qps，作为最近一段时间的qps
		recent := conf.Conf.Melloi.Recent
		if len(reportGraphs) > recent {
			for i := 1; i <= recent; i++ {
				sumQPSRecent += reportGraphs[len(reportGraphs)-i].QPS
				QPSRecent = sumQPSRecent / recent
			}
		} else {
			QPSRecent = sumQPS / len(reportGraphs)
		}
		if len(mins) != 0 {
			sort.Ints(mins)
			min = mins[0]
			sort.Ints(maxs)
			max = maxs[len(maxs)-1]
		}
		failPercentd := float64(sumError) * 100 / float64(sumCount)
		n10 := math.Pow10(3)
		fail := math.Trunc((failPercentd+0.5/n10)*n10) / n10
		failStr := strconv.FormatFloat(fail, 'f', -1, 64)
		failPercent := failStr + "%"
		num := len(reportGraphs)
		reportGraph = model.ReportGraph{
			TestName: reportGraphs[0].TestName, TestNameNick: reportGraphs[0].TestNameNick, Max: max, Min: min, Count: sumCount, Error: sumError,
			NetIo: sumNetIo / num, QPS: sumQPS / num, AvgTime: sumTime / sumCount, Ctime: reportGraphs[len(reportGraphs)-1].Ctime,
			CodeEll: sumCodeEll, CodeSll: sumCodeSll, CodeSly: sumCodeSly, CodeSls: sumCodeSls, CodeWll: sumCodeWll, CodeWly: sumCodeWly,
			CodeWle: sumCodeWle, CodeWls: sumCodeWls, CodeNonHTTP: sumCodeNonHTTP, CodeKong: sumCodeKong, CodeOthers: sumCodeOthers,
			FailPercent: failPercent, FiftyTime: fiftyTime / num, NinetyTime: ninetyTime / num, NinetyFiveTime: ninetyFiveTime / num,
			NinetyNineTime: ninetyNineTime / num, Code301: sumCode301, Code302: sumCode302, BeginTime: reportGraphs[0].Ctime, QpsRecent: QPSRecent,
		}
	}
	return
}

//addReGraphTimer add regraph timer
func (s *Service) addReGraphTimer(c context.Context, addReGrapht model.AddReGraphTimer, res chan interface{}, timeout <-chan time.Time) {
	var (
		status              = 1
		reportSummary       model.ReportSummary
		testNamesd          []string
		podNames            []string
		beginTimed          string
		afterTimed          string
		beginTime           = addReGrapht.BeginTime
		testNames           = addReGrapht.TestNames
		jobName             = addReGrapht.JobName
		token               = addReGrapht.Token
		testNameNicks       = addReGrapht.TestNameNicks
		RetMap              *model.PaasJobQueryStatus
		scriptID            = addReGrapht.ScriptID
		reportSuID          = addReGrapht.ReportSuID
		ptestJobs           []*model.PtestJob
		err                 error
		btm                 time.Time
		firstRetMap         *model.PaasJobQueryStatus
		reportGraphs        []model.ReportGraph
		JSON                []byte
		fusingList          = addReGrapht.FusingList
		useBusiStopList     = addReGrapht.UseBusiStopList
		busiStopPercentList = addReGrapht.BusiStopPercentList
	)
	timer := time.NewTicker(5 * time.Second)
	done := make(chan bool, 1)
	num := 1
	go func() {
		defer close(res)
		a := 1
		for {
			if btm, err = time.Parse("2006-01-02 15:04:05", beginTime); err != nil {
				timeout = time.After(time.Until(time.Now().Add(1 * time.Second)))
				return
			}
			select {
			case <-timer.C:
				elapsedTime := (num - 1) * 5
				//实时当前任务的所有容器
				if ptestJobs, err = s.QueryOrStopAllPtestByJobName(context.TODO(), reportSuID, false, 3); err != nil {
					log.Error("get all Job err (%v)", err)
				}
				//遍历容器，并查询每个容器的cpu
				if len(ptestJobs) > 0 {
					for _, ptestJob := range ptestJobs {
						if RetMap, err = s.dao.Job(c, token, ptestJob.JobName); err != nil {
							log.Error("get job info err (%v)", err)
						}
						if RetMap != nil && len(RetMap.Data.Pods) > 0 {
							for _, pod := range RetMap.Data.Pods {
								podNames = append(podNames, pod.Name)
								log.Info("containerID :(%s)", pod.ContainerID)
								go s.addJobCPU(context.TODO(), pod, scriptID, reportSuID, elapsedTime, ptestJob.JobName)
							}
						}
					}
				}
				log.Info("podName :(%s)", podNames)
				if firstRetMap, err = s.dao.Job(c, token, jobName); err != nil {
					log.Error("查询容器状态 error-------(%v)", err)
					//接口报错，则删除所有压测容器
					timeout = time.After(time.Until(time.Now().Add(1 * time.Second)))
					return
				}
				if JSON, err = json.Marshal(firstRetMap); err != nil {
					timeout = time.After(time.Until(time.Now().Add(1 * time.Second)))
				}
				result := string(JSON)
				log.Info("查询容器状态:(%s)", result)
				if firstRetMap.Status == 400 {
					status = 3
					timeout = time.After(time.Until(time.Now().Add(1 * time.Second)))
				}
				if firstRetMap.Status == 200 && firstRetMap.Data.ActiveNum == 1 {
					log.Info("第(%d)次同步数据", a)
					if num == 1 {
						log.Info("job 执行参数 testName:(%s), beginTime :(%s),elapsedTime :(%d), num: (%d)", testNames, beginTime, elapsedTime, num)
						for index, testName := range testNames {
							testNameNick := QueryTestNameNick(testName, testNameNicks)
							reportGraphAdd := model.ReportGraphAdd{
								JobName:             jobName,
								TestName:            testName,
								BeginTime:           beginTime,
								AfterTime:           beginTime,
								TestNameNick:        testNameNick,
								PodNames:            podNames,
								ElapsedTime:         elapsedTime,
								ReportSuID:          reportSuID,
								UseBusinessStop:     addReGrapht.UseBusinessStop,
								BusinessStopPercent: addReGrapht.BusinessStopPercent,
							}
							if addReGrapht.TestType == model.PROTOCOL_SCENE {
								reportGraphAdd.Fusing = fusingList[index]
								reportGraphAdd.UseBusinessStop = useBusiStopList[index]
								reportGraphAdd.BusinessStopPercent = busiStopPercentList[index]
							} else {
								reportGraphAdd.Fusing = addReGrapht.Fusing
							}
							go s.addReportGraph(context.TODO(), reportGraphAdd)
						}
					} else {
						for i := 1; i <= num-1; i++ {
							btm = btm.Add(time.Second * 5)
							aft := btm.Add(time.Second * 2)
							beginTimed = btm.Format("2006-01-02 15:04:05")
							afterTimed = aft.Format("2006-01-02 15:04:05")
						}
						log.Info("job 执行参数 testName:(%s),beginTime :(%s),elapsedTime :(%d)", testNames, beginTimed, elapsedTime)
						for index, testName := range testNames {
							testNameNick := QueryTestNameNick(testName, testNameNicks)
							reportGraphAdd := model.ReportGraphAdd{
								JobName:             jobName,
								TestName:            testName,
								BeginTime:           beginTimed,
								AfterTime:           afterTimed,
								TestNameNick:        testNameNick,
								PodNames:            podNames,
								ElapsedTime:         elapsedTime,
								ReportSuID:          reportSuID,
								UseBusinessStop:     addReGrapht.UseBusinessStop,
								BusinessStopPercent: addReGrapht.BusinessStopPercent,
							}
							if addReGrapht.TestType == model.PROTOCOL_SCENE {
								reportGraphAdd.Fusing = fusingList[index]
								reportGraphAdd.UseBusinessStop = useBusiStopList[index]
								reportGraphAdd.BusinessStopPercent = busiStopPercentList[index]
							} else {
								reportGraphAdd.Fusing = addReGrapht.Fusing
							}
							go s.addReportGraph(context.TODO(), reportGraphAdd)
						}
					}

				}
				podNames = []string{}
				if firstRetMap.Status == 200 && firstRetMap.Data.ActiveNum == 0 {
					timeout = time.After(time.Until(time.Now().Add(1 * time.Second)))
				}
				num++
				res <- a
				a++
			case <-timeout:
				close(done)
				log.Info("real stop time (%s)", time.Now().Format("2006-01-02 15:04:05"))
				// 关闭所有压测该项目或者接口的容器
				go s.QueryOrStopAllPtestByJobName(context.TODO(), reportSuID, true, status)
				if reportGraphs, err = s.dao.QueryReportGraph(testNameNicks); err != nil {
					log.Error("query reportGraphs error :(%v) ", err)
					return
				}
				for _, reportGraph := range reportGraphs {
					testNamesd = append(testNames, reportGraph.TestName)
				}
				// 如果reportGraphs 的成员都是同一个testName ，则执行如下代码，表示进行结果合并，生成 reportSummary
				if IsUniqObject(testNamesd) {
					if len(reportGraphs) > 0 {
						reportGraph := ReportGraphAvg(reportGraphs)
						reportSummary = model.ReportSummary{
							ID:           reportSuID,
							TestName:     reportGraph.TestName,
							TestNameNick: reportGraph.TestNameNick,
							Count:        reportGraph.Count,
							Error:        reportGraph.Error,
							QPS:          reportGraph.QPS,
							AvgTime:      reportGraph.AvgTime,
							NetIo:        reportGraph.NetIo,
							Min:          reportGraph.Min,
							Max:          reportGraph.Max,
						}
						reportSummary.TestStatus = status
					} else {
						log.Info("sorry,your test duration is too short ..... ")
						reportSummary = model.ReportSummary{ID: reportSuID, TestStatus: 3}
					}
					s.dao.UpdateReportSummary(&reportSummary)
				}
				//场景压测、多个接口走该逻辑
				if !IsUniqObject(testNamesd) && len(testNameNicks) > 0 {
					if len(reportGraphs) > 0 {
						reportSummary = model.ReportSummary{ID: reportSuID, TestStatus: status}
					} else {
						reportSummary = model.ReportSummary{ID: reportSuID, TestStatus: 3}
					}
				}
				s.dao.UpdateReportSummary(&reportSummary)
				return
			}
		}
	}()
	<-done
}

//addPtestJobHostIp add regraph timer
func (s *Service) addPtestJobHostIp(c context.Context, addReGrapht model.AddReGraphTimer, res chan interface{}, timeout <-chan time.Time) {
	var (
		RetMap      *model.PaasJobQueryStatus
		token       = addReGrapht.Token
		err         error
		ptestJobs   []*model.PtestJob
		ptestJobsd  []*model.PtestJob
		firstRetMap *model.PaasJobQueryStatus
	)
	timer := time.NewTicker(10 * time.Second)
	done := make(chan bool, 1)
	go func() {
		defer close(res)
		a := 1
		for {
			select {
			case <-timer.C:
				if firstRetMap, err = s.dao.Job(c, token, addReGrapht.JobName); err != nil {
					log.Error("查询容器状态 error-------(%v)", err)
					//接口报错，则删除所有压测容器
					return
				}
				if firstRetMap == nil || firstRetMap.Status == 400 || (firstRetMap.Status == 200 && firstRetMap.Data.ActiveNum == 0) {
					return
				}
				ptestJob := model.PtestJob{ReportSuID: addReGrapht.ReportSuID, Active: 1}
				if ptestJobs, err = s.dao.QueryPtestJob(&ptestJob); err != nil {
					return
				}
				for _, job := range ptestJobs {
					if job.HostIP == "" || job.JobIP == "" || job.JobID == "" {
						ptestJobsd = append(ptestJobsd, job)
					}
				}
				if len(ptestJobsd) > 0 {
					for _, ptestJob := range ptestJobsd {
						//查询容器状态
						if RetMap, err = s.dao.Job(c, token, ptestJob.JobName); err != nil {
							log.Error("get job info err (%v)", err)
						}
						if RetMap != nil && len(RetMap.Data.Pods) > 0 {
							for _, pod := range RetMap.Data.Pods {
								ptestjo := model.PtestJob{ID: ptestJob.ID, HostIP: pod.HostIP, JobID: pod.ContainerID, JobIP: pod.IP}
								//更新容器hostIp
								s.dao.UpdatePtestJob(ptestjo)
							}
						}
					}
				}
				res <- a
				a++
			case <-timeout:
				close(done)
				return
			}
		}
	}()
	<-done
}

//addJobCPU job cpu
func (s *Service) addJobCPU(c context.Context, pod model.PodInfo, scriptID, reportSuID, elapsedTime int, jobName string) (err error) {

	var (
		clientMoni = model.ClientMoni{
			ScriptID:   scriptID,
			ReportSuID: reportSuID,
			JobName:    jobName,
			JobNameAll: pod.Name,
			ElapsdTime: elapsedTime,
		}
		dockerStats *model.DockerStats
		CPUResult   *model.PaasQueryJobCPUResult
	)
	//首先查我们自己的容器cpu 获取服务
	if pod.ContainerID != "" {
		if dockerStats, err = s.QueryJobCPUByEP(context.TODO(), pod.ContainerID, pod.HostIP); err != nil {
			log.Error("query CPU err ...(%s)", err)
		}
		if dockerStats != nil && dockerStats.CPU != "" {
			cpu := strings.Replace(dockerStats.CPU, "%", "", 1)
			cpufloat, _ := strconv.ParseFloat(cpu, 32)
			cpuFinal := cpufloat
			cpuString := strconv.FormatFloat(cpuFinal, 'f', -1, 32)
			clientMoni.CPUUsed = cpuString
		}
		if dockerStats == nil || dockerStats.CPU == "" {
			//如果我们的服务挂了，就查 moni
			if CPUResult, err = s.QueryJobCPU(context.TODO(), pod.Name); err != nil {
				log.Error("query client cpu err(%v)", err)
				return
			}
			if len(CPUResult.Data) > 0 {
				for _, data := range CPUResult.Data {
					if len(data.Value) > 1 {
						clientMoni.CPUUsed = data.Value[1].(string)
					}
				}
			}
		}
	}

	//cpu 信息记录到数据库
	if _, err = s.AddClientMoni(&clientMoni); err != nil {
		log.Error("add ClientMoni err(%v)", err)
		return
	}
	return
}

//addReportGraph add reportGraph
func (s *Service) addReportGraph(c context.Context, reportGraphAdd model.ReportGraphAdd) {
	var (
		sumTime, sumQPS, sumNetIo, sumCodeEll, sumCodeWll, sumCodeWly, sumCodeWle, sumCodeWls, sumCodeSll, sumCodeSly, sumCodeSls,
		sumCodeKong, sumCodeNonHTTP, sumCodeOthers, sumCount, sumError, sumThreads, fiftyTime, ninetyTime, ninetyFiveTime,
		timeCount, ninetyNineTime, sumCode301, sumCode302 int
		mins, maxs    []int
		reportTimelys []*model.ReportTimely
		err           error
	)
	busiStopPerc := reportGraphAdd.BusinessStopPercent
	if reportTimelys, err = s.dao.QueryReTimely(reportGraphAdd.TestName, reportGraphAdd.BeginTime, reportGraphAdd.AfterTime, reportGraphAdd.PodNames); err != nil {
		return
	}
	log.Info("------reportTimelys:(%v)", reportTimelys)
	if len(reportTimelys) > 0 {
		for _, reportTimely := range reportTimelys {
			mins = append(mins, reportTimely.Min)
			maxs = append(maxs, reportTimely.Max)
			sumCount = sumCount + reportTimely.Count
			//计算一个容器的接口响应总时间
			timeCount = reportTimely.Count * reportTimely.AvgTime
			//计算所有容器的接口响应时间
			sumTime = sumTime + timeCount
			sumQPS = sumQPS + reportTimely.QPS
			sumError = sumError + reportTimely.Error
			sumNetIo = sumNetIo + reportTimely.NetIo
			sumCodeEll = sumCodeEll + reportTimely.CodeEll
			sumCodeWll = sumCodeWll + reportTimely.CodeWll
			sumCodeWly = sumCodeWly + reportTimely.CodeWly
			sumCodeWle = sumCodeWle + reportTimely.CodeWle
			sumCodeWls = sumCodeWls + reportTimely.CodeWls
			sumCodeSll = sumCodeSll + reportTimely.CodeSll
			sumCodeSly = sumCodeSly + reportTimely.CodeSly
			sumCodeSls = sumCodeSls + reportTimely.CodeSls
			sumCodeKong = sumCodeKong + reportTimely.CodeKong
			sumCodeNonHTTP = sumCodeNonHTTP + reportTimely.CodeNonHTTP
			sumCodeOthers = sumCodeOthers + reportTimely.CodeOthers
			sumCode301 = sumCode301 + reportTimely.Code301
			sumCode302 = sumCode302 + reportTimely.Code302
			sumThreads = sumThreads + reportTimely.ThreadsSum
			fiftyTime = fiftyTime + reportTimely.FiftyTime
			ninetyTime = ninetyTime + reportTimely.NinetyTime
			ninetyFiveTime = ninetyFiveTime + reportTimely.NinetyFiveTime
			ninetyNineTime = ninetyNineTime + reportTimely.NinetyNineTime
		}
		sort.Ints(mins)
		min := mins[0]
		sort.Ints(maxs)
		max := maxs[len(maxs)-1]
		log.Info("test_nick_name ----- :(%s)", reportGraphAdd.TestNameNick)
		reportGraph := model.ReportGraph{
			TestName: reportGraphAdd.TestName, TestNameNick: reportGraphAdd.TestNameNick,
			Count:          sumCount,
			QPS:            sumQPS,
			Error:          sumError,
			AvgTime:        sumTime / sumCount,
			NetIo:          sumNetIo,
			Min:            min,
			Max:            max,
			ElapsdTime:     reportGraphAdd.ElapsedTime,
			CodeEll:        sumCodeEll,
			CodeWll:        sumCodeWll,
			CodeWly:        sumCodeWly,
			CodeWle:        sumCodeWle,
			CodeWls:        sumCodeWls,
			CodeSll:        sumCodeSll,
			CodeSly:        sumCodeSly,
			CodeSls:        sumCodeSls,
			CodeKong:       sumCodeKong,
			CodeNonHTTP:    sumCodeNonHTTP,
			CodeOthers:     sumCodeOthers,
			ThreadsSum:     sumThreads,
			FiftyTime:      fiftyTime / len(reportTimelys),
			NinetyTime:     ninetyTime / len(reportTimelys),
			NinetyFiveTime: ninetyFiveTime / len(reportTimelys),
			NinetyNineTime: ninetyNineTime / len(reportTimelys),
			Code301:        sumCode301,
			Code302:        sumCode302,
		}
		suCodeRate, suBusinessRate := CalcuSuccess(reportGraph)
		if reportGraphAdd.Fusing == 0 {
			reportGraphAdd.Fusing = conf.Conf.Melloi.DefaultFusing
		}
		//压测熔断
		//使用业务熔断，只要http code 和 业务断言错误满足一条阈值，则熔断
		if reportGraphAdd.UseBusinessStop {
			if busiStopPerc == 0 {
				busiStopPerc = conf.Conf.Melloi.DefaultBusinessRate
			}
			if suCodeRate < reportGraphAdd.Fusing || suBusinessRate < busiStopPerc {
				log.Info("接口失败率超过设置阈值，执行自动熔断,jobName ：(%s)", reportGraphAdd.JobName)
				if _, err = s.DeleteJob(c, reportGraphAdd.JobName); err != nil {
					return
				}
				reportsu := model.ReportSummary{ID: reportGraphAdd.ReportSuID, IsFusing: true, BusinessValue: busiStopPerc, UseBusinessStop: reportGraphAdd.UseBusinessStop,
					FusingTestName: reportGraphAdd.TestName, SuccessCodeRate: suCodeRate, SuccessBusinessRate: suBusinessRate, FusingValue: reportGraphAdd.Fusing}
				s.dao.UpdateReportSummary(&reportsu)
			}
		} else { //不使用业务熔断
			if suCodeRate < reportGraphAdd.Fusing {
				log.Info("接口失败率超过设置阈值，执行自动熔断,jobName ：(%s)", reportGraphAdd.JobName)
				if _, err = s.DeleteJob(c, reportGraphAdd.JobName); err != nil {
					return
				}
				reportsu := model.ReportSummary{ID: reportGraphAdd.ReportSuID, IsFusing: true, BusinessValue: busiStopPerc,
					FusingTestName: reportGraphAdd.TestName, SuccessCodeRate: suCodeRate, SuccessBusinessRate: suBusinessRate, FusingValue: reportGraphAdd.Fusing}
				s.dao.UpdateReportSummary(&reportsu)
			}
		}

		s.dao.AddReportGraph(&reportGraph)
		return
	}
}

//CalcuSuccess calcu success
func CalcuSuccess(reportGraph model.ReportGraph) (successCodeRate, successBusinessRate int) {
	success := reportGraph.CodeEll + reportGraph.CodeKong + reportGraph.Code301 + reportGraph.Code302
	successCodeRate = success * 100 / reportGraph.Count
	successBusinessRate = (reportGraph.Count - reportGraph.Error) * 100 / reportGraph.Count
	return
}

// DelReportSummary delete report summary
func (s *Service) DelReportSummary(id int) error {
	return s.dao.DelReportSummary(id)
}
