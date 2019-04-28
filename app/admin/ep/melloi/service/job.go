package service

import (
	"context"
	"encoding/json"
	"net/http"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

//AddJob add perf job
func (s *Service) AddJob(c context.Context, job model.Job) (firstRetMap *model.PaasJobResponse, err error) {
	var (
		token      string
		clusterMap *model.PaasJobResponse
		JSON       []byte
	)
	if JSON, err = json.Marshal(job); err != nil {
		log.Error("json.Marshal err: (%v)", err)
		return
	}
	log.Info("JOBINFO: (%s)", string(JSON))
	if token, err = s.RmToken(c); err != nil {
		log.Error("token err :(%v)", err)
		return
	}
	// 调度逻辑：先从ptest_job 表里面查看
	log.Info("创建job的 token:(%s) ", token)
	if clusterMap, err = s.dao.AddJob(c, token, &job); err != nil {
		log.Error("s.dao.AddJob  err :(%v)", err)
		return
	}
	if clusterMap == nil {
		return
	}
	firstRetMap = clusterMap
	return
}

//DeleteJob force stop perf job
func (s *Service) DeleteJob(c context.Context, name string) (firstRetMap *model.PaasJobResponse, err error) {
	var token string
	if token, err = s.RmToken(c); err != nil {
		return
	}
	log.Info("开始执行删除容器的 dao 层******* ")
	if firstRetMap, err = s.dao.DeleteJob(c, token, name); err != nil {
		err = ecode.MelloiPaasRequestErr
		return
	}
	return
}

// CleanJob  clean None running job
func (s *Service) CleanJob(c context.Context) (err error) {
	var dockers = []*model.CleanableDocker{}
	if dockers, err = s.dao.QueryCleanableDocker(c); err != nil {
		log.Error("s.job query docker error(%v)", err)
		return
	}
	for _, docker := range dockers {
		if _, err := s.DeleteJob(c, docker.Name); err != nil {
			continue
		}
	}
	return
}

// QueryCleanableDocker query clearable Docker
func (s *Service) QueryCleanableDocker(c context.Context) (dockers []*model.CleanableDocker, err error) {
	return s.dao.QueryCleanableDocker(c)
}

//DeleteJobBatch force stop perf job batch
func (s *Service) DeleteJobBatch(c context.Context, JobBatch model.JobBatch) (status string) {
	for _, reportSuID := range JobBatch.ReportSuIDs {
		go s.QueryOrStopAllPtestByJobName(context.TODO(), reportSuID, true, 3)
	}
	status = "success"
	return
}

//QueryJobCPU query job cpu
func (s *Service) QueryJobCPU(c context.Context, jobNamed string) (responseBody *model.PaasQueryJobCPUResult, err error) {
	var token string
	if token, err = s.RmToken(c); err != nil {
		return
	}
	log.Info("执行查询容器实时cpu 接口 ")
	if responseBody, err = s.dao.QueryJobCPU(c, token, jobNamed); err != nil {
		log.Error("query job cpu err (%v)", err)
		err = ecode.MelloiPaasRequestErr
		return
	}
	return
}

//QueryJobCPUByEP query job cpu
func (s *Service) QueryJobCPUByEP(c context.Context, id, hostIP string) (dockerStats *model.DockerStats, err error) {
	if dockerStats, err = s.dao.QueryJobCPUByEP(c, id, hostIP); err != nil {
		log.Error("query QueryJobCpuByEP  err (%v)", err)
		return
	}
	return
}

//Job get job info
func (s *Service) Job(c context.Context, name string) (firstRetMap *model.PaasJobQueryStatus, err error) {
	var token string
	if token, err = s.RmToken(c); err != nil {
		err = ecode.MelloiPaasRequestErr
		return
	}
	if firstRetMap, err = s.dao.Job(c, token, name); err != nil {
		err = ecode.MelloiPaasRequestErr
		return
	}
	return
}

//ForceDeleteJob Force Delete Job
func (s *Service) ForceDeleteJob(c context.Context, reportSuId int) (jobList *model.JobInfoList, err error) {
	log.Info("开始执行强制删除容器******* ")
	if jobList, err = s.dao.ForceDeleteJob(reportSuId); err != nil {
		return
	}
	for _, job := range jobList.JobList {
		url := "http://" + job.HostIp + ":8999" + "/ep/docker/remove?job_name=" + job.JobName
		if _, delErr := http.Get(url); delErr != nil {
			//log.Error("强制删除失败，请重试！")
			s.QueryOrStopAllPtestByJobName(c, reportSuId, true, 3)
			return
		}
	}
	if upErr := s.dao.UpdateJobStatus(reportSuId); upErr != nil {
		return
	}
	return
}
