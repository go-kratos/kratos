package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/service/ep/footman/model"
	"go-common/library/log"
)

// AsyncBuglyInsertTapd Async Bugly Insert Tapd.
func (s *Service) AsyncBuglyInsertTapd(c context.Context) (err error) {
	return s.cache.Save(func() {
		s.BuglyInsertTapd(context.Background())
	})
}

// BuglyInsertTapd Bugly Insert Tapd.
func (s *Service) BuglyInsertTapd(c context.Context) (err error) {
	for _, projectID := range s.c.Bugly2Tapd.ProjectIds {
		var (
			issueRecords []*model.IssueRecord
			bugTemplate  *model.BugTemplate
		)

		if bugTemplate, err = s.dao.FindBugTemplates(projectID); err != nil {
			log.Error("FindBugTemplates projectId %s, error (%v)", projectID, err)
			continue
		}

		if issueRecords, err = s.dao.GetIssueRecordNotInTapd(bugTemplate.IssueFilterSQL); err != nil {
			log.Error("GetIssueRecordNotInTapd projectId %s, error (%v)", projectID, err)
			continue
		}

		for _, issueRecord := range issueRecords {
			var (
				bug   *model.Bug
				bugID string
			)
			if bug, err = s.getBugModel(issueRecord, bugTemplate); err != nil {
				log.Error("getLiveIOSBugModel error (%v)", err)
				continue
			}

			if bugID, err = s.dao.CreateBug(bug); err != nil && bugID == "" {
				log.Error("CreateBug error (%v)", err)
				continue
			}
			fmt.Println(bugID)

			if err = s.dao.UpdateIssueRecordTapdBugID(issueRecord.ID, bugID); err != nil {
				log.Error("UpdateIssueRecordTapdBugID error (%v)", err)
				continue
			}
		}
		log.Info("finish to insert bugly into tapd [%s] status,", projectID)
	}
	return
}

func (s *Service) getBugModel(issueRecord *model.IssueRecord, bugTemplate *model.BugTemplate) (bug *model.Bug, err error) {
	title := fmt.Sprintf(bugTemplate.Title, issueRecord.IssueNo, issueRecord.Version, issueRecord.Title, strconv.FormatInt(issueRecord.HappenTimes, 10), strconv.FormatInt(issueRecord.UserTimes, 10))
	description := fmt.Sprintf(bugTemplate.Description, title, issueRecord.KeyStack,
		issueRecord.IssueLink, issueRecord.IssueLink, issueRecord.IssueLink,
		issueRecord.Detail)

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
	}
	return
}

// AsyncUpdateBuglyStatusInTapd Async Update Bugly Status In Tapd.
func (s *Service) AsyncUpdateBuglyStatusInTapd(c context.Context) (err error) {
	return s.cache.Save(func() {
		s.UpdateBuglyStatusInTapd(context.Background())
	})
}

// UpdateBuglyStatusInTapd Update Bugly Status In Tapd.
func (s *Service) UpdateBuglyStatusInTapd(c context.Context) (err error) {
	for _, projectID := range s.c.Bugly2Tapd.ProjectIds {
		log.Info("start projectID [%s] update status to bugly", projectID)
		var (
			issueRecords []*model.IssueRecord
			bugTemplate  *model.BugTemplate
		)

		if bugTemplate, err = s.dao.FindBugTemplates(projectID); err != nil {
			log.Error("FindBugTemplates projectId %s, error (%v)", projectID, err)
			continue
		}

		if issueRecords, err = s.dao.GetIssueRecordHasInTapd(projectID); err != nil {
			log.Error("GetIssueRecordNotInTapd projectId %s, error (%v)", projectID, err)
			continue
		}

		for _, issueRecord := range issueRecords {
			log.Info("handle projectID [%s] bug number [%s]", projectID, issueRecord.TapdBugID)
			var (
				bugIssueException *model.IssueException
				bug               *model.Bug
			)
			if bugIssueException, err = s.dao.BuglyIssueExceptionList(context.Background(), projectID, bugTemplate.PlatformID, issueRecord.IssueNo); err != nil {
				log.Error("BuglyIssueExceptionList projectId %s, error (%v)", projectID, err)
				continue
			}

			if bug, err = s.dao.BugPre(bugTemplate.WorkspaceID, issueRecord.TapdBugID); err != nil {
				log.Error("BugPre projectId %s, error (%v)", projectID, err)
				continue
			}

			if bugIssueException == nil {
				continue
			}
			switch bugIssueException.Status {
			// bugly 待处理 -> tap 新
			case 0:
				s.updateBugToStatus(bug, "new")
			// bugly 已处理 -> tap 已解决
			case 1:
				s.updateBugToStatus(bug, "resolved")
			// bugly 处理中 -> tapd 接受处理
			case 2:
				s.updateBugToStatus(bug, "in_progress")
			default:
				//do nothing
			}
		}
		log.Info("finish projectID [%s] update status to bugly", projectID)
	}
	return
}

func (s *Service) updateBugToStatus(bug *model.Bug, status string) {
	if bug.Status != status {
		log.Info("bug [%s] update from %s to %s", bug.ID, bug.Status, status)
		bug.Status = status
		updateBug := &model.UpdateBug{
			Bug:         bug,
			CurrentUser: bug.CurrentOwner,
		}
		if err := s.dao.UpdateBug(updateBug); err != nil {
			log.Error("UpdateBug bugid %s, error (%v)", bug.ID, err)
		}
	}
}

// AsyncUpdateBugInTapd Async Update Title In Tapd.
func (s *Service) AsyncUpdateBugInTapd(c context.Context) (err error) {
	return s.cache.Save(func() {
		s.UpdateBugInTapd(context.Background())
	})
}

// UpdateBugInTapd Update Title In Tapd.
func (s *Service) UpdateBugInTapd(c context.Context) (err error) {
	for _, projectID := range s.c.Bugly2Tapd.ProjectIds {
		log.Info("start projectID [%s] update title in tapd", projectID)
		var (
			issueRecords []*model.IssueRecord
			bugTemplate  *model.BugTemplate
		)

		if bugTemplate, err = s.dao.FindBugTemplates(projectID); err != nil {
			log.Error("FindBugTemplates projectId %s, error (%v)", projectID, err)
			continue
		}

		if issueRecords, err = s.dao.GetIssueRecordHasInTapd(projectID); err != nil {
			log.Error("GetIssueRecordNotInTapd projectId %s, error (%v)", projectID, err)
			continue
		}

		for _, issueRecord := range issueRecords {
			log.Info("handle projectID [%s] bug number [%s]", projectID, issueRecord.TapdBugID)
			var bug *model.Bug

			// update title
			title := fmt.Sprintf(bugTemplate.Title, issueRecord.IssueNo, issueRecord.Version, issueRecord.Title, strconv.FormatInt(issueRecord.HappenTimes, 10), strconv.FormatInt(issueRecord.UserTimes, 10))

			if bug, err = s.dao.BugPre(bugTemplate.WorkspaceID, issueRecord.TapdBugID); err != nil {
				log.Error("BugPre projectId %s, error (%v)", projectID, err)
				continue
			}

			bug.Title = title

			// update priority
			if issueRecord.UserTimes >= 20 || issueRecord.HappenTimes >= 20 {
				bug.Priority = "high"
			}

			// update serious
			if strings.TrimSpace(bugTemplate.SeverityKey) != "" {
				keys := strings.Split(bugTemplate.SeverityKey, ",")
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
		}
		log.Info("finish projectID [%s] update title in tapd", projectID)
	}
	return
}
