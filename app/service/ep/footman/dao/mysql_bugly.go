package dao

import (
	"fmt"

	"go-common/app/service/ep/footman/model"
	"go-common/library/ecode"
	"go-common/library/log"

	pkgerr "github.com/pkg/errors"
)

const (
	_issueRecords   = "issue_records"
	_issueLastTimes = "issue_last_times"
)

// UpdateIssueRecord Update Issue Record.
func (d *Dao) UpdateIssueRecord(issueRecord *model.IssueRecord) (err error) {
	err = d.db.Table(_issueRecords).Where("issue_no = ? and version = ?", issueRecord.IssueNo, issueRecord.Version).UpdateColumn(map[string]interface{}{
		"last_time":    issueRecord.LastTime,
		"happen_times": issueRecord.HappenTimes,
		"user_times":   issueRecord.UserTimes,
	}).Error
	log.Info("update issue record: %s", issueRecord.IssueNo)
	fmt.Print("update issue record: " + issueRecord.IssueNo)
	return
}

// InsertIssueRecord Insert Issue Record.
func (d *Dao) InsertIssueRecord(issueRecord *model.IssueRecord) (err error) {
	err = d.db.Table(_issueRecords).Create(issueRecord).Error
	log.Info("insert issue record: %s", issueRecord.IssueNo)
	fmt.Println("insert issue record: " + issueRecord.IssueNo)
	return
}

// GetIssueLastTime Get Issue LastTime.
func (d *Dao) GetIssueLastTime(version string) (issueLastTime *model.IssueLastTime, err error) {
	issueLastTime = &model.IssueLastTime{}
	if err = d.db.Where("version = ?", version).First(issueLastTime).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// UpdateTaskStatus Update Task Status.
func (d *Dao) UpdateTaskStatus(issueLastTime *model.IssueLastTime) (err error) {
	err = pkgerr.WithStack(d.db.Table(_issueLastTimes).Where("version=?", issueLastTime.Version).Update("task_status", issueLastTime.TaskStatus).Error)
	return
}

// GetIssueRecord Get Issue Record.
func (d *Dao) GetIssueRecord(issueNo, version string) (issueRecord *model.IssueRecord, err error) {
	issueRecord = &model.IssueRecord{}
	if err = d.db.Where("issue_no = ? and version = ?", issueNo, version).First(issueRecord).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// GetIssueRecordNotInTapd Get Issue Record Not in tapd.
func (d *Dao) GetIssueRecordNotInTapd(issueFilterSQL string) (issueRecords []*model.IssueRecord, err error) {
	if err = d.db.Raw(issueFilterSQL).Order("id asc").Find(&issueRecords).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// GetIssueRecordHasInTapd Get Issue Record in tapd.
func (d *Dao) GetIssueRecordHasInTapd(projectID string) (issueRecords []*model.IssueRecord, err error) {
	if err = d.db.Table(_issueRecords).Where("project_id = ? and tapd_bug_id<>''", projectID).Find(&issueRecords).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}

// UpdateIssueRecordTapdBugID Update Issue Record Tapd Bug ID.
func (d *Dao) UpdateIssueRecordTapdBugID(id int64, tapdBugID string) (err error) {
	return d.db.Table(_issueRecords).Where("id=?", id).Update("tapd_bug_id", tapdBugID).Error
}

// UpdateLastIssueTime Update Last Issue Time.
func (d *Dao) UpdateLastIssueTime(issueLastTime *model.IssueLastTime) (err error) {
	return d.db.Table(_issueLastTimes).Where("version=?", issueLastTime.Version).Update("last_time", issueLastTime.LastTime).Error
}

// InsertIssueLastTime Insert Issue Last Time.
func (d *Dao) InsertIssueLastTime(issueLastTime *model.IssueLastTime) (err error) {
	return d.db.Table(_issueLastTimes).Create(issueLastTime).Error
}

// UpdateLastIssue Update Last Issue.
func (d *Dao) UpdateLastIssue(issueLastTime *model.IssueLastTime) (err error) {
	err = pkgerr.WithStack(d.db.Table(_issueLastTimes).Where("version=?", issueLastTime.Version).Update("last_issue", issueLastTime.LastIssue).Error)
	return
}

// UpdateVersionRecord Update Version Record.
func (d *Dao) UpdateVersionRecord(issueLastTime *model.IssueLastTime) (err error) {
	err = pkgerr.WithStack(d.db.Table(_issueLastTimes).Where("version=?", issueLastTime.Version).Updates(issueLastTime).Error)
	return
}
