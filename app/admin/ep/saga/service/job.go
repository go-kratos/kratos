package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/utils"
	"go-common/library/cache/redis"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

// QueryProjectJob ...
func (s *Service) QueryProjectJob(c context.Context, req *model.ProjectJobRequest) (resp *model.ProjectJobResp, err error) {
	var (
		layout        = "2006-01-02"
		queryCacheKey string
		jobs          []*model.ProjectJob
		since         time.Time
		util          time.Time
	)

	resp = &model.ProjectJobResp{ProjectID: req.ProjectID, QueryDescription: "最近一月的Jobs日常", State: req.Scope, DataInfo: []*model.DateJobInfo{}}

	year, month, day := time.Now().Date()
	util = time.Date(year, month, day-1, 0, 0, 0, 0, time.Local)
	since = util.AddDate(0, -1, 0)

	//query from redis first
	queryCacheKey = fmt.Sprintf("saga_admin_job_%d_%s_%s_%s_%d", req.ProjectID, req.Branch, req.Scope, req.Machine, req.StatisticsType)
	if err = s.dao.ItemRedis(c, queryCacheKey, &resp); err != redis.ErrNil {
		return
	}

	if resp.TotalItem, jobs, err = s.queryProjectJobByTime(c, req.ProjectID, since, util); err != nil {
		return
	}

	//init map key
	pendingTime := make(map[string][]float64)
	runningTime := make(map[string][]float64)
	for i := 1; ; i++ {
		day := since.AddDate(0, 0, i)
		if day.After(util) {
			break
		}
		dayStr := day.Format(layout)
		pendingTime[dayStr] = []float64{}
		runningTime[dayStr] = []float64{}

		resp.DataInfo = append(resp.DataInfo, &model.DateJobInfo{Date: dayStr, SlowestPendingJob: []*model.ProjectJob{}})
	}

	//根据查询条件进行过滤
	newJobs := jobs[:0]
	if req.Branch != "" {
		for _, job := range jobs {
			if job.Branch == req.Branch {
				newJobs = append(newJobs, job)
			}
		}
		jobs = newJobs
		newJobs = jobs[:0]
	}
	if req.User != "" {
		for _, job := range jobs {
			if job.User == req.User {
				newJobs = append(newJobs, job)
			}
		}
		jobs = newJobs
		newJobs = jobs[:0]
	}
	if req.Machine != "" {
		for _, job := range jobs {
			if job.Machine == req.Machine {
				newJobs = append(newJobs, job)
			}
		}
		jobs = newJobs
	}

	//统计pending running 时间
	for _, job := range jobs {
		var (
			jobInfo *model.DateJobInfo
		)
		jobDate := job.CreatedAt.Format(layout)
		for _, j := range resp.DataInfo {
			if j.Date == jobDate {
				jobInfo = j
				break
			}
		}
		jobInfo.JobTotal++
		if job.Status == req.Scope {
			jobInfo.StatusNum++
			if job.StartedAt != nil && job.CreatedAt != nil {
				pending := job.StartedAt.Sub(*job.CreatedAt).Seconds()
				pendingTime[jobDate] = append(pendingTime[jobDate], pending)
				if pending >= 300 {
					jobInfo.SlowestPendingJob = append(jobInfo.SlowestPendingJob, job)
				}
			}
			if job.Status == "success" {
				running := job.FinishedAt.Sub(*job.StartedAt).Seconds()
				runningTime[jobDate] = append(runningTime[jobDate], running)
			}
		}
	}

	for k, v := range runningTime {
		var (
			jobInfo *model.DateJobInfo
		)
		for _, j := range resp.DataInfo {
			if j.Date == k {
				jobInfo = j
				break
			}
		}

		jobInfo.PendingTime = utils.CalAverageTime(req.StatisticsType, v)
		jobInfo.RunningTime = utils.CalAverageTime(req.StatisticsType, pendingTime[k])
	}

	// set data to redis
	err = s.dao.SetItemRedis(c, queryCacheKey, resp, model.ExpiredOneDay)

	return
}

// queryProjectJobByTime ...
func (s *Service) queryProjectJobByTime(c context.Context, projectID int, since, util time.Time) (count int, result []*model.ProjectJob, err error) {
	var (
		resp     *gitlab.Response
		jobs     []gitlab.Job
		overTime bool
	)

	if _, resp, err = s.gitlab.ListProjectJobs(projectID, 1); err != nil {
		return
	}
	if resp.TotalItems <= 0 {
		return
	}

	for page := 1; ; page++ {
		if jobs, resp, err = s.gitlab.ListProjectJobs(projectID, page); err != nil {
			return
		}
		for _, job := range jobs {
			if job.CreatedAt == nil {
				continue
			}
			if job.CreatedAt.After(since) && job.CreatedAt.Before(util) {
				count++

				jobInfo := &model.ProjectJob{
					Status:     job.Status,
					Branch:     job.Ref,
					Machine:    job.Runner.Description,
					User:       job.User.Name,
					CreatedAt:  job.CreatedAt,
					StartedAt:  job.StartedAt,
					FinishedAt: job.FinishedAt}

				result = append(result, jobInfo)
			}
			if job.CreatedAt.Before(since) {
				overTime = true
			}
		}
		if overTime {
			break
		}
		if resp.NextPage == 0 {
			break
		}
	}

	return
}

// QueryProjectJobNew ...
func (s *Service) QueryProjectJobNew(c context.Context, req *model.ProjectJobRequest) (resp *model.ProjectJobResp, err error) {
	var (
		layout    = "2006-01-02"
		fmtLayout = `%d-%d-%d 00:00:00`
		jobs      []*model.StatisticsJobs
	)

	resp = &model.ProjectJobResp{ProjectID: req.ProjectID, QueryDescription: "最近一月的Jobs日常", State: req.Scope, DataInfo: []*model.DateJobInfo{}}

	year, month, day := time.Now().Date()
	until := time.Date(year, month, day-1, 0, 0, 0, 0, time.Local)
	since := until.AddDate(0, -1, 0)

	sinceStr := fmt.Sprintf(fmtLayout, since.Year(), since.Month(), since.Day())
	untilStr := fmt.Sprintf(fmtLayout, until.Year(), until.Month(), until.Day())
	if resp.TotalItem, jobs, err = s.dao.QueryJobsByTime(req.ProjectID, req, sinceStr, untilStr); err != nil {
		return
	}

	//init map key
	pendingTime := make(map[string][]float64)
	runningTime := make(map[string][]float64)
	for i := 1; ; i++ {
		day := since.AddDate(0, 0, i)
		if day.After(until) {
			break
		}
		dayStr := day.Format(layout)
		pendingTime[dayStr] = []float64{}
		runningTime[dayStr] = []float64{}

		resp.DataInfo = append(resp.DataInfo, &model.DateJobInfo{Date: dayStr, SlowestPendingJob: []*model.ProjectJob{}})
	}

	//统计pending running 时间
	for _, job := range jobs {
		var (
			jobInfo *model.DateJobInfo
		)
		jobDate := job.CreatedAt.Format(layout)
		for _, j := range resp.DataInfo {
			if j.Date == jobDate {
				jobInfo = j
				break
			}
		}
		if jobInfo == nil {
			continue
		}
		jobInfo.JobTotal++
		if job.Status == req.Scope {
			jobInfo.StatusNum++
			if job.StartedAt != nil && job.CreatedAt != nil {
				pending := job.StartedAt.Sub(*job.CreatedAt).Seconds()
				pendingTime[jobDate] = append(pendingTime[jobDate], pending)
				if pending >= 300 {
					jo := &model.ProjectJob{
						Status:     job.Status,
						User:       job.UserName,
						Branch:     job.Ref,
						Machine:    job.RunnerDescription,
						CreatedAt:  job.CreatedAt,
						StartedAt:  job.StartedAt,
						FinishedAt: job.FinishedAt,
					}
					jobInfo.SlowestPendingJob = append(jobInfo.SlowestPendingJob, jo)
				}
			}
			if job.Status == "success" {
				running := job.FinishedAt.Sub(*job.StartedAt).Seconds()
				runningTime[jobDate] = append(runningTime[jobDate], running)
			}
		}
	}

	for k, v := range runningTime {
		var (
			jobInfo *model.DateJobInfo
		)
		for _, j := range resp.DataInfo {
			if j.Date == k {
				jobInfo = j
				break
			}
		}

		jobInfo.PendingTime = utils.CalAverageTime(req.StatisticsType, v)
		jobInfo.RunningTime = utils.CalAverageTime(req.StatisticsType, pendingTime[k])
	}

	return
}

/*-------------------------------------- sync job ----------------------------------------*/

// SyncProjectJobs ...
func (s *Service) SyncProjectJobs(projectID int) (result *model.SyncResult, err error) {
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
		log.Info("sync project(%d) job time since: %v, until: %v", projectID, since, until)
		if result, err = s.SyncProjectJobsByTime(projectID, projectInfo.Name, *since, *until); err != nil {
			return
		}
	} else {
		if result, err = s.SyncProjectJobsNormal(projectID, projectInfo.Name); err != nil {
			return
		}
	}

	return
}

// SyncProjectJobsNormal ...
func (s *Service) SyncProjectJobsNormal(projectID int, projectName string) (result *model.SyncResult, err error) {
	var (
		jobs []gitlab.Job
		resp *gitlab.Response
	)
	result = &model.SyncResult{}

	for page := 1; ; page++ {
		result.TotalPage++
		if jobs, resp, err = s.gitlab.ListProjectJobs(projectID, page); err != nil {
			return
		}

		for _, job := range jobs {
			if err = s.structureDatabasejob(projectID, projectName, job); err != nil {
				log.Error("job Save Database err: projectID(%d), JobID(%d)", projectID, job.ID)
				err = nil

				errData := &model.FailData{
					ChildID: job.ID,
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

// SyncProjectJobsByTime ...
func (s *Service) SyncProjectJobsByTime(projectID int, projectName string, since, until time.Time) (result *model.SyncResult, err error) {
	var (
		jobs       []gitlab.Job
		resp       *gitlab.Response
		startQuery bool
	)
	result = &model.SyncResult{}

	if _, resp, err = s.gitlab.ListProjectJobs(projectID, 1); err != nil {
		return
	}

	page := 1
	for page <= resp.TotalPages {
		result.TotalPage++

		if !startQuery {
			if jobs, _, err = s.gitlab.ListProjectJobs(projectID, page); err != nil {
				return
			}
			if page == 1 && len(jobs) <= 0 {
				return
			}

			if jobs[0].CreatedAt.After(until) {
				page++
				continue
			} else {
				startQuery = true
				page--
				continue
			}
		}

		if jobs, _, err = s.gitlab.ListProjectJobs(projectID, page); err != nil {
			return
		}

		for _, job := range jobs {

			createTime := job.CreatedAt
			if createTime.After(since) && createTime.Before(until) {

				if err = s.structureDatabasejob(projectID, projectName, job); err != nil {
					log.Error("job Save Database err: projectID(%d), JobID(%d)", projectID, job.ID)
					err = nil

					errData := &model.FailData{
						ChildID: job.ID,
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

// structureDatabasejob ...
func (s *Service) structureDatabasejob(projectID int, projectName string, job gitlab.Job) (err error) {
	var (
		jobArtifactsFile string
		jobCommitID      string
	)

	jobArtifactsFileByte, _ := json.Marshal(job.ArtifactsFile)
	jobArtifactsFile = string(jobArtifactsFileByte)
	if job.Commit != nil {
		jobCommitID = job.Commit.ID
	}
	jobDB := &model.StatisticsJobs{
		ProjectID:         projectID,
		ProjectName:       projectName,
		CommitID:          jobCommitID,
		CreatedAt:         job.CreatedAt,
		Coverage:          job.Coverage,
		ArtifactsFile:     jobArtifactsFile,
		FinishedAt:        job.FinishedAt,
		JobID:             job.ID,
		Name:              job.Name,
		Ref:               job.Ref,
		RunnerID:          job.Runner.ID,
		RunnerDescription: job.Runner.Description,
		Stage:             job.Stage,
		StartedAt:         job.StartedAt,
		Status:            job.Status,
		Tag:               job.Tag,
		UserID:            job.User.ID,
		UserName:          job.User.Name,
		WebURL:            job.WebURL,
	}

	return s.SaveDatabasejob(jobDB)
}

// SaveDatabasejob ...
func (s *Service) SaveDatabasejob(jobDB *model.StatisticsJobs) (err error) {
	var total int

	if total, err = s.dao.HasJob(jobDB.ProjectID, jobDB.JobID); err != nil {
		log.Error("SaveDatabaseJob HasJob(%+v)", err)
		return
	}

	// found only one, so update
	if total == 1 {
		if err = s.dao.UpdateJob(jobDB.ProjectID, jobDB.JobID, jobDB); err != nil {
			log.Error("SaveDatabaseJob UpdateJob(%+v)", err)
			return
		}

		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabasejob job has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateJob(jobDB); err != nil {
		log.Error("SaveDatabaseJob CreateJob(%+v)", err)
		return
	}

	return
}
