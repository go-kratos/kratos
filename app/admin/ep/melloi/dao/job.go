package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_rmJob            = "/api/rm/job"
	_jobNodeURI       = "/api/rm/jobs"
	statsPort         = "8999"
	statsPath         = "/ep/docker/stats"
	_clearableDockers = "/ep/docker/get"
	_dockerKeyword    = "run"
)

//AddJob add perf job
func (d *Dao) AddJob(c context.Context, token string, jobInfo *model.Job) (cluster *model.PaasJobResponse, err error) {
	var (
		url        = d.c.ServiceCluster.TestHost + _jobNodeURI
		req        *http.Request
		passDetail []byte
	)

	paasJobDetail := &model.PaasJobDetail{
		BusinessUnit:   conf.Conf.Paas.BusinessUnit,
		Project:        conf.Conf.Paas.Project,
		App:            conf.Conf.Paas.App,
		Env:            conf.Conf.Paas.Env,
		Name:           jobInfo.Name,
		Image:          conf.Conf.Paas.Image,
		ImageVersion:   conf.Conf.Paas.ImageVersion,
		Volumes:        conf.Conf.Paas.Volumes,
		CPURequest:     jobInfo.CPU / 4,
		CPULimit:       jobInfo.CPU,
		MemoryRequest:  jobInfo.Memory,
		Command:        jobInfo.Command,
		ResourcePoolID: conf.Conf.Paas.ResourcePoolID,
		Parallelism:    jobInfo.Parallelism,
		Completions:    conf.Conf.Paas.Completions,
		RetriesLimit:   conf.Conf.Paas.RetriesLimit,
		NetworkID:      conf.Conf.Paas.NetworkID,
		ClusterID:      conf.Conf.Paas.ClusterID,
		TreeID:         conf.Conf.Paas.TreeID,
		HostInfo:       conf.Conf.Paas.HostInfo,
		EnvInfo:        jobInfo.EnvInfo,
	}

	log.Info("Command:(%s)", paasJobDetail.Command)

	if passDetail, err = json.Marshal(paasJobDetail); err != nil {
		return
	}
	log.Info("passDetail:(%s)", string(passDetail))

	if req, err = d.newRequest(http.MethodPost, url, paasJobDetail); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)
	if err = d.httpClient.Do(c, req, &cluster); err != nil {
		log.Error("d.AddJob url(%s) res($s) error(%v)", url, err)
		return
	}
	return
}

//DeleteJob force stop perf job
func (d *Dao) DeleteJob(c context.Context, token, name string) (cluster *model.PaasJobResponse, err error) {
	log.Info("删除job -------- token:(%s), job_name(%s)： ", token, name)
	var (
		url = d.c.ServiceCluster.TestHost + _rmJob
		req *http.Request
	)

	if req, err = d.newRequest(http.MethodDelete, url, nil); err != nil {
		return
	}
	req.Header.Set("X-Authorization-Token", token)
	q := req.URL.Query()
	q.Add("business_unit", "test")
	q.Add("project", "ep")
	q.Add("app", "melloi-launch")
	q.Add("env", "dev")
	q.Add("name", name)
	q.Add("cluster_id", "1")
	q.Add("tree_id", "39400")
	req.URL.RawQuery = q.Encode()
	if err = d.httpClient.Do(c, req, &cluster); err != nil {
		log.Error("d.Job url(%s) res($s) error(%v)", url, err)
		err = ecode.MelloiPaasRequestErr
		return
	}
	// 此处不return
	if err = d.DeletePtestJobByJobName(name); err != nil {
		log.Error("d.DeletePtestJobByJobName uerror(%v)", name)
	}
	return
}

// QueryCleanableDocker query clearable docker list
func (d *Dao) QueryCleanableDocker(c context.Context) (dockers []*model.CleanableDocker, err error) {
	var (
		url    = fmt.Sprintf("%s:%d%s", d.c.DockerStatus.Host, d.c.DockerStatus.Port, _clearableDockers)
		req    *http.Request
		bs     []byte
		result []string
	)

	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	if bs, err = d.httpClient.Raw(c, req, ""); err != nil {
		log.Error("d.job query clearable docker url(%s) error (%v)", url, err)
		return
	}
	result = strings.Split(string(bs), "\n")
	for _, item := range result {
		var docker = &model.CleanableDocker{}
		if strings.Contains(item, _dockerKeyword) {
			docker.Name = item
			dockers = append(dockers, docker)
		}
	}
	return
}

//QueryJobCPU query job cpu status
func (d *Dao) QueryJobCPU(c context.Context, token, jobName string) (responseBody *model.PaasQueryJobCPUResult, err error) {

	Query := strings.Replace(conf.Conf.Paas.Query, "JobNamed", jobName, -1)
	var (
		url                       = d.c.ServiceCluster.QueryJobCPUHost
		req                       *http.Request
		paasQueryJobPostCPUDetail = &model.PaasQueryJobCPUPostDetail{
			Action:     conf.Conf.Paas.Action,
			PublicKey:  conf.Conf.Paas.PublicKey,
			Signature:  conf.Conf.Paas.Signature,
			DataSource: conf.Conf.Paas.DataSource,
			Query:      Query,
		}
	)
	if req, err = d.newRequest(http.MethodPost, url, paasQueryJobPostCPUDetail); err != nil {
		return
	}
	req.Header.Set("X-Authorization-Token", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)
	if err = d.httpClient.Do(c, req, &responseBody); err != nil {
		log.Error("d.Job url(%s) res($s) error(%v)", url, err)
		err = ecode.MelloiPaasRequestErr
		return
	}
	return
}

//QueryJobCPUByEP query job cup by ep
func (d *Dao) QueryJobCPUByEP(c context.Context, id, hostIP string) (dokcerStats *model.DockerStats, err error) {
	var (
		req *http.Request
	)
	url := "http://" + hostIP + ":" + statsPort + statsPath
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("id", id)
	req.URL.RawQuery = q.Encode()
	if err = d.httpClient.Do(c, req, &dokcerStats); err != nil {
		log.Error("d.queryJobCPU By EP err : url(%s)  error(%v)", url, err)
		return
	}
	return
}

//Job get job info
func (d *Dao) Job(c context.Context, token, name string) (responseBody *model.PaasJobQueryStatus, err error) {
	var (
		url = d.c.ServiceCluster.TestHost + _rmJob
		req *http.Request
	)
	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("business_unit", "test")
	q.Add("project", "ep")
	q.Add("app", "melloi-launch")
	q.Add("env", "dev")
	q.Add("name", name)
	q.Add("cluster_id", "1")
	q.Add("tree_id", "39400")
	req.Header.Set("X-Authorization-Token", token)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)
	if err = d.httpClient.Do(c, req, &responseBody); err != nil {
		log.Error("d.Job url(%s) res($s) error(%v)", url, err)
		err = ecode.MelloiPaasRequestErr
		return
	}
	return
}

//ForceDeleteJob Force Delete Job
func (d *Dao) ForceDeleteJob(reportSuId int) (jobList *model.JobInfoList, err error) {
	jobList = &model.JobInfoList{}
	err = d.DB.Table("ptest_job").Select("host_ip, job_name").Where("report_su_id = ?", reportSuId).Find(&jobList.JobList).Error
	return
}

//UpdateJobStatus Update Job Status
func (d *Dao) UpdateJobStatus(reportSuId int) (err error) {
	if err = d.DB.Model(&model.PtestJob{}).Where("report_su_id = ?", reportSuId).Update("active", -1).Error; err != nil {
		return
	}
	return d.DB.Table("report_summary").Where("id = ?", reportSuId).Update("test_status", 3).Error
}
