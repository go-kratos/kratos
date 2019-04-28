package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/satori/go.uuid"
)

// RunVersions  Run Versions.
func (s *Service) RunVersions(buglyVersionID int64) (rep map[string]interface{}, err error) {
	lock := s.getBatchRunLock(buglyVersionID)
	lock.Lock()
	defer lock.Unlock()

	var (
		buglyRunVersions []*model.BuglyVersion
		uid              string
		isEnableRun      bool
	)

	// 获取可以跑的版本号list
	if buglyRunVersions, err = s.GetBatchRunVersions(); err != nil {
		log.Error("GetBatchRunVersions err(%v) ", err)
		return
	}

	for _, buglyRunVersion := range buglyRunVersions {
		if buglyRunVersion.ID == buglyVersionID {
			isEnableRun = true
			uid = uuid.NewV4().String()

			// 插入batchrun表
			buglyBatchRun := &model.BuglyBatchRun{
				BuglyVersionID: buglyRunVersion.ID,
				Version:        buglyRunVersion.Version,
				BatchID:        uid,
				Status:         model.BuglyBatchRunStatusRunning,
			}

			if err = s.dao.InsertBuglyBatchRun(buglyBatchRun); err != nil {
				return
			}

			s.batchRunCache.Do(context.Background(), func(ctx context.Context) {
				if err = s.runPerVersion(uid, buglyRunVersion); err != nil {
					log.Error("runPerVersion uid(%s) version(%s) err(%v) ", uid, buglyRunVersion.Version, err)
				}

				defer func() {
					// batch run 并且更新batch run表
					batchStatus := model.BuglyBatchRunStatusDone
					errMsg := ""
					endTime := time.Now()

					if err != nil {
						batchStatus = model.BuglyBatchRunStatusFailed
						errMsg = err.Error()
					}

					updateBuglyBatchRun := &model.BuglyBatchRun{
						ID:       buglyBatchRun.ID,
						Status:   batchStatus,
						ErrorMsg: errMsg,
						EndTime:  endTime,
					}

					if err = s.dao.UpdateBuglyBatchRun(updateBuglyBatchRun); err != nil {
						log.Error("runPerVersion UpdateBuglyBatchRun uid(%s) version(%s) err(%v) ", uid, buglyRunVersion.Version, err)
					}
				}()
			})

			break
		}
	}

	rep = make(map[string]interface{})
	rep["uid"] = uid
	rep["enable_to_run"] = isEnableRun

	if !isEnableRun {
		err = ecode.MartheTaskInRunning
	}
	return
}

// BatchRunVersions Batch Run Versions
func (s *Service) BatchRunVersions() (err error) {
	var (
		buglyRunVersions []*model.BuglyVersion
		uid              = uuid.NewV4().String()
	)

	// 获取可以跑的版本号list
	if buglyRunVersions, err = s.dao.FindEnableAndReadyVersions(); err != nil {
		log.Error("GetBatchRunVersions err(%v) ", err)
		return
	}

	for _, buglyRunVersion := range buglyRunVersions {
		if _, err = s.RunVersions(buglyRunVersion.ID); err != nil {
			log.Error("runPerVersion uid(%s) version(%s) err(%v) ", uid, buglyRunVersion.Version, err)
		}
	}
	return
}

func (s *Service) runPerVersion(uid string, buglyRunVersion *model.BuglyVersion) (err error) {
	var (
		buglyRet       *model.BugRet
		buglyCookie    *model.BuglyCookie
		requestPageCnt int
		c              = context.Background()
		lastRunTime    time.Time
		buglyBatchRun  *model.BuglyBatchRun
		buglyProject   *model.BuglyProject
	)

	log.Info("start run version: [%s] batchId: [%s]", buglyRunVersion.Version, uid)

	if buglyProject, err = s.dao.QueryBuglyProject(buglyRunVersion.BuglyProjectID); err != nil {
		return
	}

	if buglyProject.ID == 0 {
		err = ecode.NothingFound
		return
	}

	//get issue total count
	bugIssueRequest := &model.BugIssueRequest{
		ProjectID:     buglyProject.ProjectID,
		PlatformID:    buglyProject.PlatformID,
		Version:       buglyRunVersion.Version,
		ExceptionType: buglyProject.ExceptionType,
		StartNum:      0,
		Rows:          1,
	}

	//get last run time
	if buglyBatchRun, err = s.dao.QueryLastSuccessBatchRunByVersion(buglyRunVersion.Version); err != nil {
		return
	}
	if buglyBatchRun.ID > 0 {
		lastRunTime = buglyBatchRun.CTime
	} else {
		loc, _ := time.LoadLocation("Local")
		if lastRunTime, err = time.ParseInLocation(model.TimeLayout, model.TimeLayout, loc); err != nil {
			return
		}
	}

	//get enable cookie
	if buglyCookie, err = s.GetEnableCookie(); err != nil {
		return
	}

	// if cookie, update cookie as expired
	defer func() {
		if err != nil && err == ecode.MartheCookieExpired {
			s.DisableCookie(c, buglyCookie.ID)
		}
	}()

	if buglyRet, err = s.dao.BuglyIssueAndRetry(c, buglyCookie, bugIssueRequest); err != nil || len(buglyRet.BugIssues) < 1 {
		return
	}

	//获取issue count 和 page 上限
	requestPageCnt = int(buglyRet.NumFound/s.c.Bugly.IssuePageSize) + 1
	requestPageCntMax := s.c.Bugly.IssueCountUpper / s.c.Bugly.IssuePageSize
	if requestPageCnt > requestPageCntMax {
		requestPageCnt = requestPageCntMax
	}

	// update or add issue
	for i := 0; i < requestPageCnt; i++ {
		innerBreak := false

		bugIssueRequest.StartNum = s.c.Bugly.IssuePageSize * i
		bugIssueRequest.Rows = s.c.Bugly.IssuePageSize

		var ret *model.BugRet
		if ret, err = s.dao.BuglyIssueAndRetry(c, buglyCookie, bugIssueRequest); err != nil {
			return
		}

		loc, _ := time.LoadLocation("Local")
		issueLink := "/v2/crash-reporting/errors/%s/%s/report?pid=%s&searchType=detail&version=%s&start=0&date=all"
		for _, issueDto := range ret.BugIssues {
			var (
				issueTime      time.Time
				bugIssueDetail *model.BugIssueDetail
				tagStr         string
				bugDetail      = "no detail"
			)

			tmpTime := []rune(issueDto.LastTime)
			issueTime, _ = time.ParseInLocation(model.TimeLayout, string(tmpTime[:len(tmpTime)-4]), loc)

			//issue时间早于库里面最新时间的，跳出
			if issueTime.Before(lastRunTime) {
				innerBreak = true
				break
			}

			var tmpBuglyIssue *model.BuglyIssue

			if tmpBuglyIssue, err = s.dao.GetBuglyIssue(issueDto.IssueID, issueDto.Version); err != nil {
				log.Error("d.GetSaveIssues url(%s) err(%v)", "GetSaveIssues", err)
				continue
			}

			for _, bugTag := range issueDto.Tags {
				tagStr = tagStr + bugTag.TagName + ","
			}

			if tmpBuglyIssue.ID != 0 {
				//update
				issueRecord := &model.BuglyIssue{
					IssueNo:     issueDto.IssueID,
					Title:       issueDto.Title,
					LastTime:    issueTime,
					HappenTimes: issueDto.Count,
					UserTimes:   issueDto.UserCount,
					Version:     issueDto.Version,
					Tags:        tagStr,
				}
				s.dao.UpdateBuglyIssue(issueRecord)
			} else {
				//create
				if bugIssueDetail, err = s.dao.BuglyIssueDetailAndRetry(c, buglyCookie, buglyProject.ProjectID, buglyProject.PlatformID, issueDto.IssueID); err == nil {
					bugDetail = bugIssueDetail.CallStack
				}

				issueURL := s.c.Bugly.Host + fmt.Sprintf(issueLink, buglyProject.ProjectID, issueDto.IssueID, buglyProject.PlatformID, buglyRunVersion.Version)
				issueRecord := &model.BuglyIssue{
					IssueNo:      issueDto.IssueID,
					Title:        issueDto.Title,
					LastTime:     issueTime,
					HappenTimes:  issueDto.Count,
					UserTimes:    issueDto.UserCount,
					Version:      issueDto.Version,
					Tags:         tagStr,
					Detail:       bugDetail,
					ExceptionMsg: issueDto.ExceptionMsg,
					KeyStack:     issueDto.KeyStack,
					IssueLink:    issueURL,
					ProjectID:    buglyProject.ProjectID,
				}
				s.dao.InsertBuglyIssue(issueRecord)
			}
		}

		if innerBreak {
			break
		}
	}

	log.Info("end run version: [%s] batchId: [%s]", buglyRunVersion.Version, uid)
	return

}

// GetBatchRunVersions Get Batch Run Versions.
func (s *Service) GetBatchRunVersions() (buglyRunVersions []*model.BuglyVersion, err error) {
	var (
		buglyVersions  []*model.BuglyVersion
		buglyBatchRuns []*model.BuglyBatchRun
	)
	if buglyVersions, err = s.dao.FindEnableAndReadyVersions(); err != nil {
		return
	}

	if buglyBatchRuns, err = s.dao.QueryBuglyBatchRunsByStatus(model.BuglyBatchRunStatusRunning); err != nil {
		return
	}

	// 排除正在执行的版本
	for _, buglyVersion := range buglyVersions {
		var isVersionRun bool
		for _, buglyBatchRun := range buglyBatchRuns {
			if buglyBatchRun.Version == buglyVersion.Version {
				isVersionRun = true
				break
			}
		}

		if !isVersionRun {
			buglyRunVersions = append(buglyRunVersions, buglyVersion)
		}
	}
	return
}

// DisableBatchRunOverTime Disable Batch Run OverTime.
func (s *Service) DisableBatchRunOverTime() (err error) {
	var (
		buglyBatchRuns []*model.BuglyBatchRun
		tapdBugRecords []*model.TapdBugRecord
		timeNow        = time.Now()
	)

	//清 未完成 disable batch run
	if buglyBatchRuns, err = s.dao.QueryBuglyBatchRunsByStatus(model.BuglyBatchRunStatusRunning); err != nil {
		return
	}

	for _, buglyBatchRun := range buglyBatchRuns {
		if timeNow.Sub(buglyBatchRun.CTime).Hours() > float64(s.c.Scheduler.BatchRunOverHourTime) {
			updateBuglyBatchRun := &model.BuglyBatchRun{
				ID:       buglyBatchRun.ID,
				Status:   model.BuglyBatchRunStatusFailed,
				ErrorMsg: "over time",
				EndTime:  timeNow,
			}

			if err = s.dao.UpdateBuglyBatchRun(updateBuglyBatchRun); err != nil {
				continue
			}
		}
	}

	// 清 未完成 disable insert tapd bug
	if tapdBugRecords, err = s.dao.QueryTapdBugRecordByStatus(model.InsertBugStatusRunning); err != nil {
		return
	}

	for _, tapdBugRecord := range tapdBugRecords {
		if timeNow.Sub(tapdBugRecord.CTime).Hours() > float64(s.c.Scheduler.BatchRunOverHourTime) {
			tapdBugRecord.Status = model.InsertBugStatusFailed
			if err = s.dao.UpdateTapdBugRecord(tapdBugRecord); err != nil {
				continue
			}
		}
	}

	return
}

func (s *Service) getBatchRunLock(buglyVersionId int64) (batchRunLock *sync.Mutex) {
	s.syncBatchRunLock.Lock()
	defer s.syncBatchRunLock.Unlock()

	var ok bool
	if batchRunLock, ok = s.mapBatchRunLocks[buglyVersionId]; !ok {
		batchRunLock = new(sync.Mutex)
		s.mapBatchRunLocks[buglyVersionId] = batchRunLock
	}
	return
}
