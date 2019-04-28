package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/utils"
	"go-common/app/admin/ep/saga/service/wechat"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

const (
	_gitHTTP            = "http://git.bilibili.co/"
	_gitSSH             = "git@git-test.bilibili.co:"
	_gitSSHTail         = ".git"
	_manualJob          = "manual"
	_androidScheduleJob = "daily branch check"
	_iosScheduleJob     = "daily:on-schedule"
)

// QueryTeamPipeline query pipeline info according to team.
func (s *Service) QueryTeamPipeline(c context.Context, req *model.TeamDataRequest) (resp *model.PipelineDataResp, err error) {
	var (
		projectInfo []*model.ProjectInfo
		reqProject  = &model.ProjectInfoRequest{}

		data        []*model.PipelineDataTime
		queryDes    string
		total       int
		succNum     int
		key         string
		keyNotExist bool
	)

	if len(req.Department) <= 0 && len(req.Business) <= 0 {
		log.Warn("query department and business are empty!")
		return
	}

	//get pipeline info from mc
	key = "saga_admin_" + req.Department + "_" + req.Business + "_" + model.KeyTypeConst[3]
	if resp, err = s.dao.GetPipeline(c, key); err != nil {
		if err == memcache.ErrNotFound {
			keyNotExist = true
		} else {
			return
		}
	} else {
		return
	}

	log.Info("sync team pipeline start => type= %d, Department= %s, Business= %s", req.QueryType, req.Department, req.Business)

	//query team projects
	reqProject.Department = req.Department
	reqProject.Business = req.Business
	if _, projectInfo, err = s.dao.QueryProjectInfo(false, reqProject); err != nil {
		return
	}
	if len(projectInfo) <= 0 {
		log.Warn("Found no project!")
		return
	}

	if data, total, succNum, err = s.QueryTeamPipelineByTime(projectInfo, model.LastWeekPerDay); err != nil {
		return
	}

	successScale := succNum * 100 / total
	queryDes = req.Department + " " + req.Business + " " + "pipeline上一周每天数量"
	resp = &model.PipelineDataResp{
		Department:   req.Department,
		Business:     req.Business,
		QueryDes:     queryDes,
		Total:        total,
		SuccessNum:   succNum,
		SuccessScale: successScale,
		Data:         data,
	}

	//set pipeline info to mc
	if keyNotExist {
		if err = s.dao.SetPipeline(c, key, resp); err != nil {
			return
		}
	}

	log.Info("sync team pipeline end")
	return
}

// QueryTeamPipelineByTime ...
func (s *Service) QueryTeamPipelineByTime(projectInfo []*model.ProjectInfo, queryType int) (resp []*model.PipelineDataTime, allNum, succNum int, err error) {
	var (
		layout = "2006-01-02"
		since  time.Time
		until  time.Time

		total   int
		success int
		count   int
	)

	if queryType == model.LastWeekPerDay {
		count = model.DayNumPerWeek
	} else {
		log.Warn("Query Type is not in range!")
		return
	}

	year, month, day := time.Now().Date()
	weekDay := (int)(time.Now().Weekday())
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	for i := 0; i < count; i++ {
		since = today.AddDate(0, 0, -weekDay-i)
		until = today.AddDate(0, 0, -weekDay-i+1)

		totalAll := 0
		successAll := 0
		//log.Info("== start query from: %v, to: %v", since, until)
		for _, project := range projectInfo {
			if total, success, err = s.QueryProjectPipeline(project.ProjectID, "success", since, until); err != nil {
				return
			}
			totalAll = totalAll + total
			successAll = successAll + success
		}

		perData := &model.PipelineDataTime{
			TotalItem:   totalAll,
			SuccessItem: successAll,
			StartTime:   since.Format(layout),
			EndTime:     until.Format(layout),
		}
		resp = append(resp, perData)
		allNum = allNum + totalAll
		succNum = succNum + successAll
	}

	return
}

// QueryProjectPipeline query pipeline info according to project id.
func (s *Service) QueryProjectPipeline(projectID int, state string, since, until time.Time) (totalNum, stateNum int, err error) {
	var (
		pipelineList gitlab.PipelineList
		pipeline     *gitlab.Pipeline
		resp         *gitlab.Response
		startQuery   bool
	)

	if _, resp, err = s.gitlab.ListProjectPipelines(1, projectID, ""); err != nil {
		return
	}

	page := 1
	for page <= resp.TotalPages {

		if !startQuery {
			if pipelineList, _, err = s.gitlab.ListProjectPipelines(page, projectID, ""); err != nil {
				return
			}
			if page == 1 && len(pipelineList) <= 0 {
				return
			}
			if pipeline, _, err = s.gitlab.GetPipeline(projectID, pipelineList[0].ID); err != nil {
				return
			}

			if pipeline.CreatedAt.After(until) {
				page++
				continue
			} else {
				startQuery = true
				page--
				continue
			}
		}

		if pipelineList, _, err = s.gitlab.ListProjectPipelines(page, projectID, ""); err != nil {
			return
		}

		for _, v := range pipelineList {
			if pipeline, _, err = s.gitlab.GetPipeline(projectID, v.ID); err != nil {
				return
			}

			createTime := pipeline.CreatedAt
			//year, month, day := createTime.Date()
			//log.Info("index: %d createTime: %d, month: %d, day: %d", k, year, month, day)

			if createTime.After(since) && createTime.Before(until) {
				totalNum = totalNum + 1
				if pipeline.Status == state {
					stateNum = stateNum + 1
				}
			}

			if createTime.Before(since) {
				return
			}
		}
		page++
	}

	return
}

// QueryProjectPipelineNew ...
func (s *Service) QueryProjectPipelineNew(c context.Context, req *model.PipelineDataReq) (resp *model.PipelineDataAvgResp, err error) {
	var (
		data            []*model.PipelineDataAvg
		queryDes        string
		total           int
		totalStatus     int
		avgDurationTime float64
		avgPendingTime  float64
		avgRunningTime  float64
	)

	log.Info("QuerySingleProjectData Type: %d", req.Type)
	switch req.Type {
	case model.LastYearPerMonth:
		queryDes = model.LastYearPerMonthNote
	case model.LastMonthPerDay:
		queryDes = model.LastMonthPerDayNote
	case model.LastYearPerDay:
		queryDes = model.LastYearPerDayNote
	default:
		log.Warn("QueryProjectCommit Type is not in range")
		return
	}
	queryDes = req.ProjectName + " pipeline " + req.State + " " + queryDes

	if data, total, totalStatus, avgDurationTime, avgPendingTime, avgRunningTime, err = s.QueryProjectByTimeNew(c, req, req.Type); err != nil {
		return
	}

	resp = &model.PipelineDataAvgResp{
		ProjectName:     req.ProjectName,
		QueryDes:        queryDes,
		Status:          req.State,
		Total:           total,
		TotalStatus:     totalStatus,
		AvgDurationTime: avgDurationTime,
		AvgPendingTime:  avgPendingTime,
		AvgRunningTime:  avgRunningTime,
		Data:            data,
	}

	return
}

// QueryProjectByTimeNew ...
func (s *Service) QueryProjectByTimeNew(c context.Context, req *model.PipelineDataReq, queryType int) (resp []*model.PipelineDataAvg, allNum, allStatusNum int, avgDurationTime, avgPendingTime, avgRunningTime float64, err error) {
	var (
		layout              = "2006-01-02"
		since               time.Time
		until               time.Time
		count               int
		pendingTimeListAll  []float64
		runningTimeListAll  []float64
		durationTimeListAll []float64
		pipelineTime        *model.PipelineTime
		perData             *model.PipelineDataAvg
		avgTotalTime        float64
		totalNum            int
		statusNum           int
	)

	year, month, day := time.Now().Date()
	thisMonth := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	if queryType == model.LastYearPerMonth {
		count = model.MonthNumPerYear
	} else if queryType == model.LastMonthPerDay {
		//_, _, count = thisMonth.AddDate(0, 0, -1).Date()
		count = model.DayNumPerMonth
	} else if queryType == model.LastYearPerDay {
		count = model.DayNumPerYear
	}

	for i := count; i >= 1; i-- {

		if queryType == model.LastYearPerMonth {
			since = thisMonth.AddDate(0, -i, 0)
			until = thisMonth.AddDate(0, -i+1, 0)
		} else if queryType == model.LastMonthPerDay {
			since = thisMonth.AddDate(0, 0, -i)
			until = thisMonth.AddDate(0, 0, -i+1)
		} else if queryType == model.LastYearPerDay {
			since = thisMonth.AddDate(0, 0, -i)
			until = thisMonth.AddDate(0, 0, -i+1)
		}

		/*if totalNum, statusNum, pipelineTime, err = s.QueryProjectPipelines(c, req, since, until); err != nil {
			log.Error("QueryProjectPipelines err：%+v", err)
			return
		}*/
		if totalNum, statusNum, pipelineTime, err = s.QueryPipelinesFromDB(c, req, since, until); err != nil {
			log.Error("QueryPipelinesFromDB err：%+v", err)
			return
		}

		avgTotalTime = utils.CalAverageTime(req.StatisticsType, pipelineTime.DurationList)
		avgPendingTime = utils.CalAverageTime(req.StatisticsType, pipelineTime.PendingList)
		avgRunningTime = utils.CalAverageTime(req.StatisticsType, pipelineTime.RunningList)

		perData = &model.PipelineDataAvg{
			TotalItem:       totalNum,
			TotalStatusItem: statusNum,
			AvgDurationTime: avgTotalTime,
			AvgPendingTime:  avgPendingTime,
			AvgRunningTime:  avgRunningTime,
			MaxDurationTime: pipelineTime.DurationMax,
			MinDurationTime: pipelineTime.DurationMin,
			MaxPendingTime:  pipelineTime.PendingMax,
			MinPendingTime:  pipelineTime.PendingMin,
			MaxRunningTime:  pipelineTime.RunningMax,
			MinRunningTime:  pipelineTime.RunningMin,
			StartTime:       since.Format(layout),
			EndTime:         until.Format(layout),
		}
		resp = append(resp, perData)
		allNum = allNum + totalNum
		allStatusNum = allStatusNum + statusNum

		pendingTimeListAll = utils.CombineSlice(pendingTimeListAll, pipelineTime.PendingList)
		runningTimeListAll = utils.CombineSlice(runningTimeListAll, pipelineTime.RunningList)
		durationTimeListAll = utils.CombineSlice(durationTimeListAll, pipelineTime.DurationList)
	}

	avgDurationTime = utils.CalAverageTime(req.StatisticsType, durationTimeListAll)
	avgPendingTime = utils.CalAverageTime(req.StatisticsType, pendingTimeListAll)
	avgRunningTime = utils.CalAverageTime(req.StatisticsType, runningTimeListAll)
	log.Info("avgDurationTime: %v, avgPendingTime: %v, avgRunningTime: %v", avgDurationTime, avgPendingTime, avgRunningTime)

	return
}

// QueryProjectPipelines ...
func (s *Service) QueryProjectPipelines(c context.Context, req *model.PipelineDataReq, since, until time.Time) (totalNum, statusNum int, pipelineTime *model.PipelineTime, err error) {
	var (
		pipelineList gitlab.PipelineList
		pipeline     *gitlab.Pipeline
		resp         *gitlab.Response
		startQuery   bool
		meetTime     bool
		projectID    = req.ProjectID
		pendingTime  float64
		runningTime  float64
		totalTime    float64
	)

	pipelineTime = &model.PipelineTime{}

	opt := &gitlab.ListProjectPipelinesOptions{}
	if _, resp, err = s.gitlab.ListProjectPipelines(1, projectID, ""); err != nil {
		log.Error("ListProjectPipelines err: %+v", err)
		return
	}

	page := 1
	for page <= resp.TotalPages {
		opt.ListOptions.Page = page
		if !startQuery && (!since.IsZero() || !until.IsZero()) {
			if pipelineList, _, err = s.gitlab.ListProjectPipelines(page, projectID, ""); err != nil {
				log.Error("ListProjectPipelines err: %+v", err)
				return
			}
			if page == 1 && len(pipelineList) <= 0 {
				return
			}
			if pipeline, _, err = s.gitlab.GetPipeline(projectID, pipelineList[0].ID); err != nil {
				return
			}

			if pipeline.CreatedAt.After(until) {
				page++
				continue
			} else {
				startQuery = true
				page--
				continue
			}
		}

		// start query
		if pipelineList, _, err = s.gitlab.ListProjectPipelines(page, projectID, ""); err != nil {
			return
		}

		meetTime = true
		for _, v := range pipelineList {
			if pipeline, _, err = s.gitlab.GetPipeline(projectID, v.ID); err != nil {
				return
			}

			createTime := pipeline.CreatedAt
			if !since.IsZero() || !until.IsZero() {
				meetTime = createTime.After(since) && createTime.Before(until)
			}

			//the pipeline is we need
			if meetTime {
				totalNum = totalNum + 1

				if req.Branch != "" && req.Branch != pipeline.Ref {
					continue
				} else if req.User != "" && req.User != pipeline.User.Name {
					continue
				} else if req.State != "" && req.State != pipeline.Status {
					continue
				}

				statusNum = statusNum + 1

				if pipeline.Status != "cancel" {
					if pipeline.StartedAt == nil {
						pendingTime = 0
						runningTime = pipeline.FinishedAt.Sub(*pipeline.CreatedAt).Seconds()
					} else {
						pendingTime = pipeline.StartedAt.Sub(*pipeline.CreatedAt).Seconds()
						runningTime = pipeline.FinishedAt.Sub(*pipeline.StartedAt).Seconds()
					}
					totalTime = pipeline.FinishedAt.Sub(*pipeline.CreatedAt).Seconds()

					pipelineTime.PendingMax, pipelineTime.PendingMin = utils.CalSizeTime(pendingTime, pipelineTime.PendingMax, pipelineTime.PendingMin)
					pipelineTime.RunningMax, pipelineTime.RunningMin = utils.CalSizeTime(runningTime, pipelineTime.RunningMax, pipelineTime.RunningMin)
					pipelineTime.DurationMax, pipelineTime.DurationMin = utils.CalSizeTime(totalTime, pipelineTime.DurationMax, pipelineTime.DurationMin)

					pipelineTime.PendingList = append(pipelineTime.PendingList, pendingTime)
					pipelineTime.RunningList = append(pipelineTime.RunningList, runningTime)
					pipelineTime.DurationList = append(pipelineTime.DurationList, totalTime)
				}
			}

			// time is over, so return
			if (!since.IsZero() || !until.IsZero()) && createTime.Before(since) {
				return
			}
		}

		page++
	}
	return
}

// QueryPipelinesFromDB ...
func (s *Service) QueryPipelinesFromDB(c context.Context, req *model.PipelineDataReq, since, until time.Time) (totalNum, statusNum int, pipelineTime *model.PipelineTime, err error) {
	var (
		fmtLayout   = `%d-%d-%d 00:00:00`
		pipelines   []*model.StatisticsPipeline
		projectID   = req.ProjectID
		pendingTime float64
		runningTime float64
		totalTime   float64
	)
	pipelineTime = &model.PipelineTime{}

	sinceStr := fmt.Sprintf(fmtLayout, since.Year(), since.Month(), since.Day())
	untilStr := fmt.Sprintf(fmtLayout, until.Year(), until.Month(), until.Day())
	if totalNum, statusNum, pipelines, err = s.dao.QueryPipelinesByTime(projectID, req, sinceStr, untilStr); err != nil {
		return
	}

	for _, pipeline := range pipelines {
		if pipeline.Status == model.StatusCancel {
			continue
		}
		if pipeline.FinishedAt == nil {
			continue
		}

		if pipeline.StartedAt == nil {
			pendingTime = 0
			runningTime = pipeline.FinishedAt.Sub(*pipeline.CreatedAt).Seconds()
		} else {
			pendingTime = pipeline.StartedAt.Sub(*pipeline.CreatedAt).Seconds()
			runningTime = pipeline.FinishedAt.Sub(*pipeline.StartedAt).Seconds()
		}
		totalTime = pipeline.FinishedAt.Sub(*pipeline.CreatedAt).Seconds()

		pipelineTime.PendingMax, pipelineTime.PendingMin = utils.CalSizeTime(pendingTime, pipelineTime.PendingMax, pipelineTime.PendingMin)
		pipelineTime.RunningMax, pipelineTime.RunningMin = utils.CalSizeTime(runningTime, pipelineTime.RunningMax, pipelineTime.RunningMin)
		pipelineTime.DurationMax, pipelineTime.DurationMin = utils.CalSizeTime(totalTime, pipelineTime.DurationMax, pipelineTime.DurationMin)

		pipelineTime.PendingList = append(pipelineTime.PendingList, pendingTime)
		pipelineTime.RunningList = append(pipelineTime.RunningList, runningTime)
		pipelineTime.DurationList = append(pipelineTime.DurationList, totalTime)
	}
	return
}

//alertProjectPipelineProc cron func
func (s *Service) alertProjectPipelineProc() {
	for _, alert := range conf.Conf.Property.Git.AlertPipeline {
		projectId := alert.ProjectID
		runningTimeout := alert.RunningTimeout
		runningRate := alert.RunningRate
		runningThreshold := alert.RunningThreshold
		pendingTimeout := alert.PendingTimeout
		pendingThreshold := alert.PendingThreshold
		go func() {
			var err error
			if err = s.PipelineAlert(context.TODO(), projectId, runningTimeout, runningThreshold, runningRate, gitlab.Running); err != nil {
				log.Error("PipelineAlert Running (%+v)", err)
			}
			if err = s.PipelineAlert(context.TODO(), projectId, pendingTimeout, pendingThreshold, 0, gitlab.Pending); err != nil {
				log.Error("PipelineAlert Pending (%+v)", err)
			}
		}()
	}
}

//PipelineAlert ...
func (s *Service) PipelineAlert(c context.Context, projectID, timeout, threshold, rate int, status gitlab.BuildStateValue) (err error) {
	var (
		layout          = "2006-01-02 15:04:05"
		pipeline        *gitlab.Pipeline
		timeoutNum      int
		message         string
		pipelineurl     string
		durationTime    float64
		pipelineSUM     int
		timeoutPipeline string
		pipelineList    gitlab.PipelineList
		resp            *gitlab.Response
		projectInfo     *model.ProjectInfo
		userlist        = conf.Conf.Property.Git.UserList
		w               = wechat.New(s.dao)
		sendMessage     = false
	)
	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}
	repo := projectInfo.Repo
	if len(repo) > len(_gitSSH) {
		repo = repo[len(_gitSSH) : len(repo)-len(_gitSSHTail)]
		repo = _gitHTTP + repo
	}
	timeNow := time.Now().Format(layout)
	message = fmt.Sprintf("[SAGA]Pipeline 告警   %v\n项目：%s\n", timeNow, repo)
	for page := 1; ; page++ {
		if pipelineList, resp, err = s.git.ListProjectPipelines(page, projectID, status); err != nil {
			return
		}
		for _, item := range pipelineList {
			pipelineSUM += 1
			if pipeline, _, err = s.git.GetPipeline(projectID, item.ID); err != nil {
				return
			}
			if status == gitlab.Pending {
				durationTime = pipeline.UpdatedAt.Sub(*pipeline.CreatedAt).Minutes()
			} else if status == gitlab.Running {
				//此处时间计算换成job
				if durationTime, err = s.PipelineRunningTime(projectID, item.ID); err != nil {
					return
				}
			}
			if int(durationTime) >= timeout {
				timeoutNum += 1
				pipelineurl = fmt.Sprintf("%d. %s/pipelines/%d (%vmin)\n", timeoutNum, repo, pipeline.ID, int(durationTime))
				timeoutPipeline += pipelineurl
			}
		}
		if resp.NextPage == 0 {
			break
		}
	}

	if timeoutPipeline != "" {
		message += fmt.Sprintf(`列表(url|%s时间):%s%s`, status, "\n", timeoutPipeline)
	}
	if pipelineSUM > 0 {
		message += fmt.Sprintf(`状态：%s 总数为%d个`, status, pipelineSUM)
	}
	if status == gitlab.Pending {
		var alertMessage string
		message += fmt.Sprintf(`%s告警：`, "\n")
		if pipelineSUM >= threshold {
			alertMessage = fmt.Sprintf(`[ 数量（%d）>=警戒值（%d） ]`, pipelineSUM, threshold)
			sendMessage = true
		}
		message += alertMessage
		if timeoutNum >= 1 {
			if alertMessage != "" {
				message = message[:strings.LastIndex(message, " ]")] + fmt.Sprintf(`，%s时间>=警戒值（%d） ]`, status, timeout)
			} else {
				message += fmt.Sprintf(`[ %s时间>=警戒值（%d） ]`, status, timeout)
			}
			sendMessage = true
		}
	}
	if status == gitlab.Running && timeoutNum >= threshold {
		sendMessage = true
		message += fmt.Sprintf(`，时间>%dmin为%d个%s告警：[ 数量（%d）>=警戒值（%d) ]`, timeout, timeoutNum, "\n", timeoutNum, threshold)
		if timeoutNum*100/pipelineSUM >= rate {
			message = message[:strings.LastIndex(message, " ]")] + fmt.Sprintf(`，比例（%v%s）>=警戒值%d%s ]`, timeoutNum*100/pipelineSUM, "%", rate, "%")
		}
	}
	if sendMessage {
		return w.PushMsg(c, userlist, message)
	}
	return
}

//PipelineRunningTime ...
func (s *Service) PipelineRunningTime(projectID, pipelineID int) (durationTime float64, err error) {
	var jobList []*gitlab.Job
	if jobList, _, err = s.git.ListPipelineJobs(nil, projectID, pipelineID); err != nil {
		return
	}
	for _, job := range jobList {
		if job.Status != _manualJob && job.Name != _androidScheduleJob && job.Name != _iosScheduleJob {
			if job.FinishedAt != nil && job.StartedAt != nil {
				durationTime += job.FinishedAt.Sub(*job.StartedAt).Minutes()
			} else if job.StartedAt != nil {
				durationTime += time.Since(*job.StartedAt).Minutes()
			}
		}
	}
	return
}

/*-------------------------------------- sync pipeline ----------------------------------------*/

// SyncProjectPipelines ...
func (s *Service) SyncProjectPipelines(projectID int) (result *model.SyncResult, err error) {
	var (
		//syncAllTime = conf.Conf.Property.SyncData.SyncAllTime
		syncAllTime = false
		since       *time.Time
		until       *time.Time
		projectInfo *model.ProjectInfo
	)
	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}
	if !syncAllTime {

		since, until = utils.CalSyncTime()
		log.Info("sync project id(%d), name(%s) pipeline, time since: %v, until: %v", projectID, projectInfo.Name, since, until)
		if result, err = s.SyncProjectPipelinesByTime(projectID, projectInfo.Name, *since, *until); err != nil {
			return
		}
	} else {
		log.Info("sync project id(%d), name(%s) pipeline", projectID, projectInfo.Name)
		if result, err = s.SyncProjectAllPipelines(projectID, projectInfo.Name); err != nil {
			return
		}
	}

	return
}

// SyncProjectPipelinesByTime ...
func (s *Service) SyncProjectPipelinesByTime(projectID int, projectName string, since, until time.Time) (result *model.SyncResult, err error) {
	var (
		pipelines  gitlab.PipelineList
		pipeline   *gitlab.Pipeline
		resp       *gitlab.Response
		startQuery bool
	)
	result = &model.SyncResult{}

	if _, resp, err = s.gitlab.ListProjectPipelines(1, projectID, ""); err != nil {
		return
	}

	page := 1
	for page <= resp.TotalPages {
		result.TotalPage++

		if !startQuery {
			if pipelines, resp, err = s.gitlab.ListProjectPipelines(page, projectID, ""); err != nil {
				return
			}
			if page == 1 && len(pipelines) <= 0 {
				return
			}

			if pipeline, _, err = s.gitlab.GetPipeline(projectID, pipelines[0].ID); err != nil {
				return
			}

			if pipeline.CreatedAt.After(until) {
				page++
				continue
			} else {
				startQuery = true
				page--
				continue
			}
		}

		// start query
		if pipelines, _, err = s.gitlab.ListProjectPipelines(page, projectID, ""); err != nil {
			return
		}

		for _, v := range pipelines {

			if pipeline, _, err = s.gitlab.GetPipeline(projectID, v.ID); err != nil {
				return
			}

			createTime := pipeline.CreatedAt
			if createTime.After(since) && createTime.Before(until) {

				if err = s.structureDBPipeline(projectID, projectName, pipeline); err != nil {
					log.Error("pipeline Save Database err: projectID(%d), PipelineID(%d)", projectID, pipeline.ID)
					err = nil

					errData := &model.FailData{
						ChildID: pipeline.ID,
					}
					result.FailData = append(result.FailData, errData)

					continue
				}
				result.TotalNum++
			}

			if createTime.Before(since) {
				return
			}
		}
		page++
	}

	return
}

// SyncProjectAllPipelines ...
func (s *Service) SyncProjectAllPipelines(projectID int, projectName string) (result *model.SyncResult, err error) {
	var (
		pipelines gitlab.PipelineList
		pipeline  *gitlab.Pipeline
		resp      *gitlab.Response
	)
	result = &model.SyncResult{}

	for page := 1; ; page++ {
		result.TotalPage++
		if pipelines, resp, err = s.gitlab.ListProjectPipelines(page, projectID, ""); err != nil {
			return
		}

		for _, v := range pipelines {
			if pipeline, _, err = s.gitlab.GetPipeline(projectID, v.ID); err != nil {
				return
			}

			if err = s.structureDBPipeline(projectID, projectName, pipeline); err != nil {
				log.Error("pipeline Save Database err: projectID(%d), PipelineID(%d)", projectID, pipeline.ID)
				err = nil

				errData := &model.FailData{
					ChildID: pipeline.ID,
				}
				result.FailData = append(result.FailData, errData)

				continue
			}
			result.TotalNum++
		}

		if resp.NextPage == 0 {
			break
		}
	}
	return
}

// structureDBPipeline ...
func (s *Service) structureDBPipeline(projectID int, projectName string, pipeline *gitlab.Pipeline) (err error) {

	statisticPipeline := &model.StatisticsPipeline{
		PipelineID:   pipeline.ID,
		ProjectID:    projectID,
		ProjectName:  projectName,
		Status:       pipeline.Status,
		Ref:          pipeline.Ref,
		Tag:          pipeline.Tag,
		User:         pipeline.User.Name,
		CreatedAt:    pipeline.CreatedAt,
		UpdatedAt:    pipeline.UpdatedAt,
		StartedAt:    pipeline.StartedAt,
		FinishedAt:   pipeline.FinishedAt,
		CommittedAt:  pipeline.CommittedAt,
		Coverage:     pipeline.Coverage,
		Duration:     pipeline.Duration,
		DurationTime: 0,
	}

	return s.SaveDatabasePipeline(statisticPipeline)
}

// SaveDatabasePipeline ...
func (s *Service) SaveDatabasePipeline(pipelineDB *model.StatisticsPipeline) (err error) {
	var total int

	if total, err = s.dao.HasPipeline(pipelineDB.ProjectID, pipelineDB.PipelineID); err != nil {
		log.Error("SaveDatabasePipeline HasPipeline(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		if err = s.dao.UpdatePipeline(pipelineDB.ProjectID, pipelineDB.PipelineID, pipelineDB); err != nil {
			log.Error("SaveDatabasePipeline UpdatePipeline err(%+v)", err)
			return
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabasePipeline pipeline has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreatePipeline(pipelineDB); err != nil {
		log.Error("SaveDatabasePipeline CreatePipeline err(%+v)", err)
		return
	}

	return
}
