package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/ep/footman/model"
	"go-common/library/log"
)

const (
	_issueLink = "/v2/crash-reporting/errors/%s/%s/report?pid=%s&searchType=detail&version=%s&start=0&date=all"
)

// Issue Issue.
func (s *Service) Issue(c context.Context, version string) (status int, err error) {
	status = 0

	bugIssueRequest := &model.BugIssueRequest{
		Version:  version,
		StartNum: 0,
		Rows:     100,
	}

	_, err = s.dao.BuglyIssue(c, bugIssueRequest)
	return
}

// UpdateToken Update Token.
func (s *Service) UpdateToken(c context.Context) (status int, err error) {
	status = 0
	if err = s.dao.UpdateToken(); err != nil {
		status = -1
	}

	return
}

// SaveIssue Save Issue.
func (s *Service) SaveIssue(c context.Context, version string) (status int, err error) {
	s.cache.Save(func() {
		s.GetSaveIssues(context.Background(), version)
	})
	return
}

// SaveIssues Save Issues.
func (s *Service) SaveIssues(c context.Context) (status int, err error) {
	s.cache.Save(func() {
		s.GetSaveIssuesWithMultiVersion(context.Background())
	})
	return
}

// TaskGetSaveIssues Task Get Save Issues.
func (s *Service) TaskGetSaveIssues() {
	fmt.Println("task start")
	s.GetSaveIssues(context.Background(), "5.29.1")
}

// GetSaveIssuesWithMultiVersion Get Save Issues With Multi Version.
func (s *Service) GetSaveIssuesWithMultiVersion(c context.Context) (status int, err error) {
	var versions []string
	if versions, err = ReadLine(s.c.Bugly.Version); err != nil {
		log.Error("failed to get version with file %s", s.c.Bugly.Version)
		return
	}

	for _, version := range versions {

		if _, err = s.GetSaveIssues(context.Background(), version); err != nil {
			log.Error("failed to execute function version %s %+v", version, err)
		}

	}

	return
}

// GetSaveIssues Get Save Issues.
func (s *Service) GetSaveIssues(c context.Context, versions string) (status int, err error) {
	versionAry := strings.Split(versions, ",")
	if len(versionAry) != 3 {
		fmt.Print("version 每行参数为3个 项目id，平台id，版本id")
		return
	}

	projectID := strings.TrimSpace(versionAry[0])
	platformID := strings.TrimSpace(versionAry[1])
	version := strings.TrimSpace(versionAry[2])
	versionCnt := 0
	if strings.Contains(version, "*") {
		versionPatten := strings.Replace(version, "*", "", -1)
		var bugVersionList []*model.BugVersion
		if bugVersionList, err = s.dao.BugVersion(c, projectID, platformID); err != nil {
			return
		}
		fmt.Print(len(bugVersionList))
		for _, bugVersion := range bugVersionList {
			if strings.Contains(bugVersion.Name, versionPatten) {
				versionCnt = versionCnt + 1

				s.GetSaveIssuesWithPerVersion(c, projectID, platformID, bugVersion.Name)

			}
		}
		fmt.Print("#########################" + strconv.Itoa(versionCnt))

	} else {
		versionCnt = versionCnt + 1
		s.GetSaveIssuesWithPerVersion(c, projectID, platformID, version)
	}
	fmt.Print("#########################" + strconv.Itoa(versionCnt))
	return
}

// GetSaveIssuesWithPerVersion Get Save Issues With Per Version.
func (s *Service) GetSaveIssuesWithPerVersion(c context.Context, projectID, platformID, version string) (status int, err error) {
	fmt.Print("start GetSaveIssues version: " + version)
	var (
		requestCnt      int
		taskStatus      int
		buglyRet        *model.BugRet
		lastTimeInSQL   time.Time
		currentLastTime time.Time
		timeLayout      = "2006-01-02 15:04:05"
		lastIssue       string
		exceptionType   string
	)

	loc, _ := time.LoadLocation("Local")

	//get last time in sql
	if lastTimeInSQL, taskStatus, err = s.GetVersionIssueLastTime(version); err != nil || taskStatus == 1 {
		return
	}

	if platformID == "1" {
		exceptionType = "Crash,Native"
	} else if platformID == "2" {
		exceptionType = "Crash,ExtensionCrash"
	}

	//get issue total count
	bugIssueRequest := &model.BugIssueRequest{
		ProjectID:     projectID,
		PlatformID:    platformID,
		Version:       version,
		ExceptionType: exceptionType,
		StartNum:      0,
		Rows:          1,
	}
	if buglyRet, err = s.dao.BuglyIssueAndRetry(c, bugIssueRequest); err != nil || len(buglyRet.BugIssues) < 1 {
		return
	}

	requestCnt = int(buglyRet.NumFound/100) + 1
	if requestCnt > 50 {
		requestCnt = 50
	}

	rs := []rune(buglyRet.BugIssues[0].LastTime)

	if currentLastTime, err = time.ParseInLocation(timeLayout, string(rs[:len(rs)-4]), loc); err != nil || currentLastTime.Before(lastTimeInSQL) {
		return
	}

	issueLastTimeTmp := &model.IssueLastTime{
		Version:    version,
		TaskStatus: 1,
	}
	if err = s.dao.UpdateTaskStatus(issueLastTimeTmp); err != nil {
		return
	}

	//get new issue list
	for i := 0; i < requestCnt; i++ {
		innerBreak := false
		/*	tmpRequest := &model.BugIssueRequest{
			ProjectID:  projectId,
			PlatformID: platformId,
			Version:    version,
			StartNum:   100 * i,
			Rows:       100,
		}*/

		bugIssueRequest.StartNum = 100 * i
		bugIssueRequest.Rows = 100

		var ret *model.BugRet
		if ret, err = s.dao.BuglyIssueAndRetry(c, bugIssueRequest); err != nil {
			break
		}

		for _, issueDto := range ret.BugIssues {
			var (
				issueTime      time.Time
				bugIssueDetail *model.BugIssueDetail
				tagStr         string
				bugDetail      = "no detail"
			)
			tmpTime := []rune(issueDto.LastTime)
			issueTime, _ = time.ParseInLocation(timeLayout, string(tmpTime[:len(tmpTime)-4]), loc)

			//issue时间早于库里面最新时间的，跳出
			if issueTime.Before(lastTimeInSQL) {
				innerBreak = true
				break
			}

			for _, bugTag := range issueDto.Tags {
				tagStr = tagStr + bugTag.TagName + ","
			}

			var tmpIssueRecord *model.IssueRecord
			if tmpIssueRecord, err = s.dao.GetIssueRecord(issueDto.IssueID, issueDto.Version); err != nil {
				log.Error("d.GetSaveIssues url(%s) err(%v)", "GetSaveIssues", err)
				continue
			}
			if tmpIssueRecord.ID != 0 {
				//update
				issueRecord := &model.IssueRecord{
					IssueNo:     issueDto.IssueID,
					Title:       issueDto.Title,
					LastTime:    issueTime,
					HappenTimes: issueDto.Count,
					UserTimes:   issueDto.UserCount,
					Version:     issueDto.Version,
					Tags:        tagStr,
				}
				s.dao.UpdateIssueRecord(issueRecord)
			} else {
				//create
				if bugIssueDetail, err = s.dao.BuglyIssueDetailAndRetry(c, projectID, platformID, issueDto.IssueID); err == nil {
					bugDetail = bugIssueDetail.CallStack
				}
				issueURL := s.c.Bugly.Host + fmt.Sprintf(_issueLink, projectID, issueDto.IssueID, platformID, version)
				issueRecord := &model.IssueRecord{
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
					ProjectID:    projectID,
				}
				s.dao.InsertIssueRecord(issueRecord)
			}
			lastIssue = issueDto.IssueID
		}
		if innerBreak {
			break
		}
	}

	issueLastTime := &model.IssueLastTime{
		LastTime:   currentLastTime,
		Version:    version,
		TaskStatus: 0,
		LastIssue:  lastIssue,
	}

	if err = s.dao.UpdateTaskStatus(issueLastTime); err != nil {
		return
	}

	if err = s.dao.UpdateLastIssueTime(issueLastTime); err != nil {
		return
	}

	if err = s.dao.UpdateLastIssue(issueLastTime); err != nil {
		return
	}

	fmt.Print("end GetSaveIssues version:" + version)
	return
}

// GetVersionIssueLastTime Ge tVersion Issue LastTime.
func (s *Service) GetVersionIssueLastTime(version string) (lastTime time.Time, taskStatus int, err error) {
	var issueLastTime *model.IssueLastTime
	if issueLastTime, err = s.dao.GetIssueLastTime(version); err != nil {
		return
	}

	if issueLastTime.ID != 0 {
		lastTime = issueLastTime.LastTime
		taskStatus = issueLastTime.TaskStatus

	} else {
		toBeCharge := "2018-01-01 00:00:00"
		timeLayout := "2006-01-02 15:04:05"
		loc, _ := time.LoadLocation("Local")
		lastTime, err = time.ParseInLocation(timeLayout, toBeCharge, loc)
		taskStatus = 0
		issueLastTime := &model.IssueLastTime{
			LastTime:   lastTime,
			Version:    version,
			TaskStatus: taskStatus,
		}
		s.dao.InsertIssueLastTime(issueLastTime)
	}

	return
}

// ReadLine Read Line.
func ReadLine(fileName string) (lines []string, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	buf := bufio.NewReader(f)
	for {
		var line string
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)
		lines = append(lines, line)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

	}
}
