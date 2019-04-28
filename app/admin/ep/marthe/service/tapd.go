package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// UpdateTapdBugTpl Update Tapd Bug Tpl.
func (s *Service) UpdateTapdBugTpl(c context.Context, username string, req *model.UpdateTapdBugTplRequest) (rep map[string]interface{}, err error) {
	//check if access to workspace
	if !s.AccessToWorkspace(req.WorkspaceID, username) {
		err = ecode.AccessDenied
		return
	}

	//check sql
	if _, err = s.CheckTapdBugTplSQL(c, req.IssueFilterSQL); err != nil {
		return
	}

	tapdBugTemplate := &model.TapdBugTemplate{
		ID:             req.ID,
		WorkspaceID:    req.WorkspaceID,
		BuglyProjectId: req.BuglyProjectId,
		IssueFilterSQL: req.IssueFilterSQL,
		SeverityKey:    req.SeverityKey,
		UpdateBy:       username,

		TapdProperty: req.TapdProperty,
	}

	if req.ID > 0 {
		err = s.dao.UpdateTapdBugTemplate(tapdBugTemplate)
	} else {
		// add new
		err = s.dao.InsertTapdBugTemplate(tapdBugTemplate)
	}

	rep = make(map[string]interface{})
	rep["template_id"] = tapdBugTemplate.ID
	return
}

// QueryTapdBugTpl Query Tapd Bug Tpl
func (s *Service) QueryTapdBugTpl(c context.Context, req *model.QueryTapdBugTemplateRequest) (rep *model.PaginateTapdBugTemplates, err error) {
	var (
		total                           int64
		tapdBugTemplates                []*model.TapdBugTemplate
		tapdBugTemplateWithProjectNames []*model.TapdBugTemplateWithProjectName
	)
	if total, tapdBugTemplates, err = s.dao.FindTapdBugTemplates(req); err != nil {
		return
	}

	for _, tapdBugTemplate := range tapdBugTemplates {
		var buglyProject *model.BuglyProject
		if buglyProject, err = s.dao.QueryBuglyProject(tapdBugTemplate.BuglyProjectId); err != nil {
			return
		}

		if req.ProjectName == "" || strings.Contains(buglyProject.ProjectName, req.ProjectName) {
			tapdBugTemplateWithProjectName := &model.TapdBugTemplateWithProjectName{
				TapdBugTemplate: tapdBugTemplate,
				ProjectName:     buglyProject.ProjectName,
			}
			tapdBugTemplateWithProjectNames = append(tapdBugTemplateWithProjectNames, tapdBugTemplateWithProjectName)
		}
	}

	rep = &model.PaginateTapdBugTemplates{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		TapdBugTemplateWithProjectNames: tapdBugTemplateWithProjectNames,
	}
	return

}

// QueryAllTapdBugTpl Query All Tapd Bug Tpl
func (s *Service) QueryAllTapdBugTpl(c context.Context) (rep []*model.TapdBugTemplateShortResponse, err error) {
	var tapdBugTemplates []*model.TapdBugTemplate

	if tapdBugTemplates, err = s.dao.QueryAllTapdBugTemplates(); err != nil {
		return
	}

	for _, tapdBugTemplate := range tapdBugTemplates {
		var buglyProject *model.BuglyProject
		if buglyProject, err = s.dao.QueryBuglyProject(tapdBugTemplate.BuglyProjectId); err != nil {
			return
		}

		ret := &model.TapdBugTemplateShortResponse{
			ID:               tapdBugTemplate.ID,
			WorkspaceID:      tapdBugTemplate.WorkspaceID,
			BuglyProjectId:   tapdBugTemplate.BuglyProjectId,
			BuglyProjectName: buglyProject.ProjectName,
		}
		rep = append(rep, ret)
	}

	return
}

// UpdateTapdBugVersionTpl Update Tapd Bug Version Tpl.
func (s *Service) UpdateTapdBugVersionTpl(c context.Context, username string, req *model.UpdateTapdBugVersionTplRequest) (rep map[string]interface{}, err error) {

	// check sql
	if req.IssueFilterSQL != "" {
		if _, err = s.CheckTapdBugTplSQL(c, req.IssueFilterSQL); err != nil {
			return
		}
	}

	//check project id
	var tmp *model.TapdBugTemplate

	if tmp, err = s.dao.QueryTapdBugTemplate(req.ProjectTemplateID); err != nil {
		return
	}

	if tmp.ID == 0 {
		err = ecode.NothingFound
		return
	}

	//check if access to workspace
	if !s.AccessToWorkspace(tmp.WorkspaceID, username) {
		err = ecode.AccessDenied
		return
	}

	tapdBugVersionTemplate := &model.TapdBugVersionTemplate{
		ID:                req.ID,
		Version:           req.Version,
		ProjectTemplateID: req.ProjectTemplateID,

		IssueFilterSQL: req.IssueFilterSQL,
		SeverityKey:    req.SeverityKey,
		UpdateBy:       username,

		TapdProperty: req.TapdProperty,
	}

	if req.ID > 0 {
		err = s.dao.UpdateTapdBugVersionTemplate(tapdBugVersionTemplate)
	} else {
		// add new
		err = s.dao.InsertTapdBugVersionTemplate(tapdBugVersionTemplate)
	}

	rep = make(map[string]interface{})
	rep["version_template_id"] = tapdBugVersionTemplate.ID
	return
}

// QueryTapdBugVersionTpl Query Tapd Bug Version Tpl
func (s *Service) QueryTapdBugVersionTpl(c context.Context, req *model.QueryTapdBugVersionTemplateRequest) (rep *model.PaginateTapdBugVersionTemplates, err error) {
	var (
		total                   int64
		tapdBugVersionTemplates []*model.TapdBugVersionTemplate
	)
	if total, tapdBugVersionTemplates, err = s.dao.FindTapdBugVersionTemplates(req); err != nil {
		return
	}
	rep = &model.PaginateTapdBugVersionTemplates{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		TapdBugVersionTemplates: tapdBugVersionTemplates,
	}
	return

}

// CheckTapdBugTplSQL Check Tapd Bug Tpl SQL.
func (s *Service) CheckTapdBugTplSQL(c context.Context, sql string) (rep map[string]interface{}, err error) {
	var buglyIssues []*model.BuglyIssue
	if buglyIssues, err = s.getBuglyIssuesByFilterSQLWithNoTapdBug(sql); err != nil {
		err = ecode.MartheFilterSqlError
		return
	}
	rep = make(map[string]interface{})
	rep["count"] = len(buglyIssues)
	return
}

// BugInsertTapdWithProject Bug Insert Tapd With Project.
func (s *Service) BugInsertTapdWithProject(c context.Context, id int64, username string) (rep map[string]interface{}, err error) {
	var count int
	count, err = s.bugInsertTapd(id, 0, username)
	rep = make(map[string]interface{})
	rep["tapd_bug_count"] = count
	return

}

// BugInsertTapdWithVersion Bug Insert Tapd With Version.
func (s *Service) BugInsertTapdWithVersion(c context.Context, id int64, username string) (rep map[string]interface{}, err error) {
	var count int
	count, err = s.bugInsertTapd(0, id, username)
	rep = make(map[string]interface{})
	rep["tapd_bug_count"] = count
	return
}

func (s *Service) bugInsertTapd(projectTemplateId, versionTemplateId int64, username string) (count int, err error) {

	var (
		tapdBugTemplate        *model.TapdBugTemplate
		tapdBugVersionTemplate *model.TapdBugVersionTemplate
		buglyIssues            []*model.BuglyIssue
		tapdBugRecords         []*model.TapdBugRecord
	)

	if (projectTemplateId == 0 && versionTemplateId == 0) || (projectTemplateId > 0 && versionTemplateId > 0) {
		return
	}

	if projectTemplateId > 0 {
		if tapdBugTemplate, err = s.dao.QueryTapdBugTemplate(projectTemplateId); err != nil {
			return
		}

		if tapdBugTemplate.ID == 0 {
			err = ecode.NothingFound
			return
		}

	} else if versionTemplateId > 0 {
		if tapdBugVersionTemplate, err = s.dao.QueryTapdBugVersionTemplate(versionTemplateId); err != nil {
			return
		}

		if tapdBugVersionTemplate.ID == 0 {
			err = ecode.NothingFound
			return
		}

		projectTemplateId = tapdBugVersionTemplate.ProjectTemplateID

		if tapdBugTemplate, err = s.dao.QueryTapdBugTemplate(projectTemplateId); err != nil {
			return
		}

		if tapdBugTemplate.ID == 0 {
			err = ecode.NothingFound
			return
		}

		if err = s.UpdateTplAsVersionTpl(tapdBugVersionTemplate, tapdBugTemplate); err != nil {
			return
		}
	}

	// get map lock, every project tpl 串行
	insertLock := s.getTapdBugInsertLock(projectTemplateId)
	insertLock.Lock()
	defer insertLock.Unlock()

	//check if access to workspace
	if !s.AccessToWorkspace(tapdBugTemplate.WorkspaceID, username) {
		err = ecode.AccessDenied
		return
	}

	if tapdBugRecords, err = s.dao.QueryTapdBugRecordByProjectIDAndStatus(projectTemplateId, model.InsertBugStatusRunning); err != nil {
		return
	}

	if len(tapdBugRecords) > 0 {
		err = ecode.MartheTaskInRunning
		return
	}

	if buglyIssues, err = s.getBuglyIssuesByFilterSQLWithNoTapdBug(tapdBugTemplate.IssueFilterSQL); err != nil {
		return
	}

	count = len(buglyIssues)

	if count > 0 {
		tapdBugInsertLog := &model.TapdBugRecord{
			ProjectTemplateID: projectTemplateId,
			VersionTemplateID: versionTemplateId,
			Operator:          username,
			Count:             count,
			IssueFilterSQL:    tapdBugTemplate.IssueFilterSQL,
			Status:            model.InsertBugStatusRunning,
		}

		if err = s.dao.InsertTapdBugRecord(tapdBugInsertLog); err != nil {
			return
		}

		s.tapdBugCache.Do(context.TODO(), func(ctx context.Context) {

			defer func() {
				if err != nil {
					tapdBugInsertLog.Status = model.InsertBugStatusFailed
				} else {
					tapdBugInsertLog.Status = model.InsertBugStatusDone
				}
				err = s.dao.UpdateTapdBugRecord(tapdBugInsertLog)
			}()

			for _, buglyIssue := range buglyIssues {
				var (
					bug   *model.Bug
					bugID string
				)

				if bug, err = s.getBugModel(buglyIssue, tapdBugTemplate); err != nil {
					log.Error("getLiveIOSBugModel error (%v)", err)
					continue
				}

				if bugID, err = s.dao.CreateBug(bug); err != nil && bugID == "" {
					log.Error("CreateBug error (%v)", err)
					continue
				}
				log.Info("insert issue no: [%s], bug number [%s]", buglyIssue.IssueNo, bugID)

				if err = s.dao.UpdateBuglyIssueTapdBugID(buglyIssue.ID, bugID); err != nil {
					log.Error("UpdateIssueRecordTapdBugID error (%v)", err)
					continue
				}

				buglyIssue.TapdBugID = bugID

				if err = s.updateBugInTapd(tapdBugTemplate, buglyIssue, time.Now()); err != nil {
					continue
				}

				log.Info("update issue no [%s], bug number [%s]", buglyIssue.IssueNo, buglyIssue.TapdBugID)
			}
		})
	}

	return
}

func (s *Service) getBugModel(bugIssue *model.BuglyIssue, bugTemplate *model.TapdBugTemplate) (bug *model.Bug, err error) {
	title := fmt.Sprintf(bugTemplate.Title, bugIssue.IssueNo, bugIssue.Version, bugIssue.Title, strconv.Itoa(bugIssue.HappenTimes), strconv.Itoa(bugIssue.UserTimes))
	description := fmt.Sprintf(bugTemplate.Description, title, bugIssue.KeyStack,
		bugIssue.IssueLink, bugIssue.IssueLink, bugIssue.IssueLink,
		bugIssue.Detail)

	bug = &model.Bug{
		Title:            title,
		Description:      description,
		Priority:         bugTemplate.Priority,
		Severity:         bugTemplate.Severity,
		Module:           bugTemplate.Module,
		Status:           bugTemplate.Status,
		Reporter:         bugTemplate.Reporter,
		BugType:          bugTemplate.BugType,
		CurrentOwner:     bugTemplate.CurrentOwner,
		Source:           bugTemplate.Source,
		OriginPhase:      bugTemplate.OriginPhase,
		Platform:         bugTemplate.Platform,
		ReleaseID:        bugTemplate.ReleaseID,
		CustomFieldThree: bugTemplate.CustomFieldThree,
		CustomFieldFour:  bugTemplate.CustomFieldFour,
		WorkspaceID:      bugTemplate.WorkspaceID,
		IterationID:      bugTemplate.IterationID,
	}
	return
}

// UpdateTplAsVersionTpl Update Tpl As Version Tpl.
func (s *Service) UpdateTplAsVersionTpl(tapdBugVersionTemplate *model.TapdBugVersionTemplate, tapdBugTemplate *model.TapdBugTemplate) (err error) {

	if tapdBugVersionTemplate.ProjectTemplateID != tapdBugTemplate.ID {
		log.Error("projectID should be [%s],actual [%s]", tapdBugVersionTemplate.ProjectTemplateID, tapdBugTemplate.ID)
		return
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Title) != "" {
		tapdBugTemplate.Title = tapdBugVersionTemplate.Title
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Description) != "" {
		tapdBugTemplate.Description = tapdBugVersionTemplate.Description
	}

	if strings.TrimSpace(tapdBugVersionTemplate.CurrentOwner) != "" {
		tapdBugTemplate.CurrentOwner = tapdBugVersionTemplate.CurrentOwner
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Platform) != "" {
		tapdBugTemplate.Platform = tapdBugVersionTemplate.Platform
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Module) != "" {
		tapdBugTemplate.Module = tapdBugVersionTemplate.Module
	}

	if strings.TrimSpace(tapdBugVersionTemplate.IterationID) != "" {
		tapdBugTemplate.IterationID = tapdBugVersionTemplate.IterationID
	}

	if strings.TrimSpace(tapdBugVersionTemplate.ReleaseID) != "" {
		tapdBugTemplate.ReleaseID = tapdBugVersionTemplate.ReleaseID
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Priority) != "" {
		tapdBugTemplate.Priority = tapdBugVersionTemplate.Priority
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Severity) != "" {
		tapdBugTemplate.Severity = tapdBugVersionTemplate.Severity
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Source) != "" {
		tapdBugTemplate.Source = tapdBugVersionTemplate.Source
	}

	if strings.TrimSpace(tapdBugVersionTemplate.CustomFieldFour) != "" {
		tapdBugTemplate.CustomFieldFour = tapdBugVersionTemplate.CustomFieldFour
	}

	if strings.TrimSpace(tapdBugVersionTemplate.BugType) != "" {
		tapdBugTemplate.BugType = tapdBugVersionTemplate.BugType
	}

	if strings.TrimSpace(tapdBugVersionTemplate.OriginPhase) != "" {
		tapdBugTemplate.OriginPhase = tapdBugVersionTemplate.OriginPhase
	}

	if strings.TrimSpace(tapdBugVersionTemplate.CustomFieldThree) != "" {
		tapdBugTemplate.CustomFieldThree = tapdBugVersionTemplate.CustomFieldThree
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Reporter) != "" {
		tapdBugTemplate.Reporter = tapdBugVersionTemplate.Reporter
	}

	if strings.TrimSpace(tapdBugVersionTemplate.Status) != "" {
		tapdBugTemplate.Status = tapdBugVersionTemplate.Status
	}

	if strings.TrimSpace(tapdBugVersionTemplate.IssueFilterSQL) != "" {
		tapdBugTemplate.IssueFilterSQL = tapdBugVersionTemplate.IssueFilterSQL
	}

	if strings.TrimSpace(tapdBugVersionTemplate.SeverityKey) != "" {
		tapdBugTemplate.SeverityKey = tapdBugVersionTemplate.SeverityKey
	}

	return
}

// BatchRunUpdateBugInTapd Batch Run Update Bug In Tapd.
func (s *Service) BatchRunUpdateBugInTapd() (err error) {
	var buglyIssues []*model.BuglyIssue
	if buglyIssues, err = s.dao.GetBuglyIssuesHasInTapd(); err != nil {
		return
	}

	for _, buglyIssue := range buglyIssues {
		var (
			tapdBugTemplate        *model.TapdBugTemplate
			tapdBugVersionTemplate *model.TapdBugVersionTemplate
		)
		if tapdBugTemplate, err = s.dao.QueryTapdBugTemplateByProjectID(buglyIssue.ProjectID); err != nil {
			continue
		}

		if tapdBugTemplate.ID == 0 {
			continue
		}

		if tapdBugVersionTemplate, err = s.dao.QueryTapdBugVersionTemplateByVersion(buglyIssue.Version); err != nil {
			continue
		}

		if tapdBugVersionTemplate.ID > 0 {
			if err = s.UpdateTplAsVersionTpl(tapdBugVersionTemplate, tapdBugTemplate); err != nil {
				continue
			}
		}
		err = s.updateBugInTapd(tapdBugTemplate, buglyIssue, time.Now())
	}

	return
}

func (s *Service) updateBugInTapd(tapdBugTemplate *model.TapdBugTemplate, bugIssue *model.BuglyIssue, timeNow time.Time) (err error) {
	var bug *model.Bug

	// update title
	title := fmt.Sprintf(tapdBugTemplate.Title, bugIssue.IssueNo, bugIssue.Version, bugIssue.Title, strconv.Itoa(bugIssue.HappenTimes), strconv.Itoa(bugIssue.UserTimes))

	if bug, err = s.dao.BugPre(tapdBugTemplate.WorkspaceID, bugIssue.TapdBugID); err != nil {
		log.Error("BugPre projectId %s, error (%v)", bugIssue.ProjectID, err)
		return
	}

	bug.Title = title

	// update priority
	var (
		tapdBugPriorityConfs []*model.TapdBugPriorityConf
	)
	if tapdBugPriorityConfs, err = s.dao.QueryTapdBugPriorityConfsByProjectTemplateIdAndStatus(tapdBugTemplate.ID, model.TapdBugPriorityConfEnable); err != nil {
		return
	}

	for _, tmpTapdBugPriorityConf := range tapdBugPriorityConfs {
		if timeNow.After(tmpTapdBugPriorityConf.StartTime) && timeNow.Before(tmpTapdBugPriorityConf.EndTime) {
			if bugIssue.UserTimes >= tmpTapdBugPriorityConf.Urgent || bugIssue.HappenTimes >= tmpTapdBugPriorityConf.Urgent {
				bug.Priority = "urgent"
			} else if bugIssue.UserTimes >= tmpTapdBugPriorityConf.High || bugIssue.HappenTimes >= tmpTapdBugPriorityConf.High {
				bug.Priority = "high"
			} else if bugIssue.UserTimes >= tmpTapdBugPriorityConf.Medium || bugIssue.HappenTimes >= tmpTapdBugPriorityConf.Medium {
				bug.Priority = "medium"
			}
			break
		}
	}

	// update serious
	if strings.TrimSpace(tapdBugTemplate.SeverityKey) != "" {
		keys := strings.Split(tapdBugTemplate.SeverityKey, ",")
		for _, key := range keys {
			if strings.Contains(bug.Description, key) {
				bug.Severity = "serious"
				break
			}
		}
	}

	updateBug := &model.UpdateBug{
		Bug:         bug,
		CurrentUser: bug.CurrentOwner,
	}

	if err := s.dao.UpdateBug(updateBug); err != nil {
		log.Error("UpdateBug bugid %s, error (%v)", bug.ID, err)
	}
	return
}

// QueryBugRecords Query Bug Records
func (s *Service) QueryBugRecords(c context.Context, req *model.QueryBugRecordsRequest) (rep *model.PaginateBugRecords, err error) {
	var (
		total          int64
		tapdBugRecords []*model.TapdBugRecord
	)
	if total, tapdBugRecords, err = s.dao.FindBugRecords(req); err != nil {
		return
	}
	rep = &model.PaginateBugRecords{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		TapdBugRecords: tapdBugRecords,
	}
	return

}

// QueryTapdBugPriorityConfsRequest Query Tapd Bug Priority Confs Request
func (s *Service) QueryTapdBugPriorityConfsRequest(c context.Context, req *model.QueryTapdBugPriorityConfsRequest) (rep *model.PaginateTapdBugPriorityConfs, err error) {
	var (
		total                int64
		tapdBugPriorityConfs []*model.TapdBugPriorityConf
	)
	if total, tapdBugPriorityConfs, err = s.dao.FindTapdBugPriorityConfs(req); err != nil {
		return
	}
	rep = &model.PaginateTapdBugPriorityConfs{
		PaginationRep: model.PaginationRep{
			Total:    total,
			PageSize: req.PageSize,
			PageNum:  req.PageNum,
		},
		TapdBugPriorityConfs: tapdBugPriorityConfs,
	}
	return

}

// UpdateTapdBugPriorityConf Update Tapd Bug Priority Conf.
func (s *Service) UpdateTapdBugPriorityConf(c context.Context, username string, req *model.UpdateTapdBugPriorityConfRequest) (rep map[string]interface{}, err error) {
	var (
		startTime                time.Time
		endTime                  time.Time
		tapdBugPriorityConfsInDB []*model.TapdBugPriorityConf
		tapdBugTemplate          *model.TapdBugTemplate
	)

	//check project id
	if tapdBugTemplate, err = s.dao.QueryTapdBugTemplate(req.ProjectTemplateID); err != nil {
		return
	}

	if tapdBugTemplate.ID == 0 {
		err = ecode.NothingFound
		return
	}

	//check if access to workspace
	if !s.AccessToWorkspace(tapdBugTemplate.WorkspaceID, username) {
		err = ecode.AccessDenied
		return
	}

	if startTime, err = time.ParseInLocation(model.TimeLayout, req.StartTime, time.Local); err != nil {
		return
	}

	if endTime, err = time.ParseInLocation(model.TimeLayout, req.EndTime, time.Local); err != nil {
		return
	}

	// when status == enable, check time if conflict
	if req.Status == model.TapdBugPriorityConfEnable {
		if tapdBugPriorityConfsInDB, err = s.dao.QueryTapdBugPriorityConfsByProjectTemplateIdAndStatus(req.ProjectTemplateID, model.TapdBugPriorityConfEnable); err != nil {
			return
		}

		for _, tapdBugPriorityConfInDB := range tapdBugPriorityConfsInDB {
			isStartTimeInDuration := tapdBugPriorityConfInDB.StartTime.After(startTime) && tapdBugPriorityConfInDB.StartTime.Before(endTime)
			isEndTimeInDuration := tapdBugPriorityConfInDB.EndTime.After(startTime) && tapdBugPriorityConfInDB.EndTime.Before(endTime)
			if isStartTimeInDuration || isEndTimeInDuration {
				err = ecode.MartheTimeConflictError
				return
			}
		}

	}

	tapdBugPriorityConf := &model.TapdBugPriorityConf{
		ID:                req.ID,
		ProjectTemplateID: req.ProjectTemplateID,
		Urgent:            req.Urgent,
		High:              req.High,
		Medium:            req.Medium,
		UpdateBy:          username,
		Status:            req.Status,
		StartTime:         startTime,
		EndTime:           endTime,
	}

	if req.ID > 0 {
		err = s.dao.UpdateTapdBugPriorityConf(tapdBugPriorityConf)
	} else {
		// add new
		err = s.dao.InsertTapdBugPriorityConf(tapdBugPriorityConf)
	}

	rep = make(map[string]interface{})
	rep["tapd_bug_priority_conf"] = tapdBugPriorityConf.ID
	return
}

// AccessToWorkspace Access To Workspace.
func (s *Service) AccessToWorkspace(workspaceID, username string) (isAccess bool) {
	if !s.c.Tapd.BugOperateAuth {
		return true
	}

	var (
		usernames   []string
		contactInfo *model.ContactInfo
		err         error
	)

	if contactInfo, err = s.dao.QueryContactInfoByUsername(username); err != nil {
		return
	}

	if contactInfo.ID == 0 {
		err = ecode.NothingFound
		return
	}
	if usernames, err = s.dao.WorkspaceUser(workspaceID); err != nil {
		return
	}

	for _, u := range usernames {
		if u == contactInfo.NickName {
			return true
		}
	}
	return
}

// HttpAccessToWorkspace Access To Workspace.
func (s *Service) HttpAccessToWorkspace(c context.Context, workspaceID, username string) (rep map[string]interface{}, err error) {
	var (
		usernames   []string
		contactInfo *model.ContactInfo
		isAccess    bool
	)

	if contactInfo, err = s.dao.QueryContactInfoByUsername(username); err != nil {
		return
	}

	if contactInfo.ID == 0 {
		err = ecode.NothingFound
		return
	}
	if usernames, err = s.dao.WorkspaceUser(workspaceID); err != nil {
		return
	}

	for _, u := range usernames {
		if u == contactInfo.NickName {
			isAccess = true
			break
		}
	}

	rep = make(map[string]interface{})
	rep["is_access"] = isAccess
	rep["names"] = usernames
	return
}

func (s *Service) getBuglyIssuesByFilterSQLWithNoTapdBug(issueFilterSQL string) (buglyIssues []*model.BuglyIssue, err error) {
	var tmpBuglyIssues []*model.BuglyIssue

	if tmpBuglyIssues, err = s.dao.GetBuglyIssuesByFilterSQL(issueFilterSQL); err != nil {
		return
	}

	for _, buglyIssue := range tmpBuglyIssues {
		if buglyIssue.TapdBugID == "" {
			buglyIssues = append(buglyIssues, buglyIssue)
		}
	}
	return
}

func (s *Service) getTapdBugInsertLock(projectTplId int64) (tapdBugInsertLock *sync.Mutex) {
	s.syncTapdBugInsertLock.Lock()
	defer s.syncTapdBugInsertLock.Unlock()

	var ok bool
	if tapdBugInsertLock, ok = s.mapTapdBugInsertLocks[projectTplId]; !ok {
		tapdBugInsertLock = new(sync.Mutex)
		s.mapTapdBugInsertLocks[projectTplId] = tapdBugInsertLock
	}
	return
}
