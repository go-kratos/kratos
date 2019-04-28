package dao

import (
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// GetBuglyIssue Get Issue Record.
func (d *Dao) GetBuglyIssue(issueNo, version string) (buglyIssue *model.BuglyIssue, err error) {
	buglyIssue = &model.BuglyIssue{}
	if err = d.db.Where("issue_no = ? and version = ?", issueNo, version).First(buglyIssue).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// UpdateBuglyIssue Update Issue Record.
func (d *Dao) UpdateBuglyIssue(buglyIssue *model.BuglyIssue) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyIssue{}).Where("issue_no = ? and version = ?", buglyIssue.IssueNo, buglyIssue.Version).UpdateColumn(map[string]interface{}{
		"last_time":    buglyIssue.LastTime,
		"happen_times": buglyIssue.HappenTimes,
		"user_times":   buglyIssue.UserTimes,
	}).Error)
}

// InsertBuglyIssue Insert Issue Record.
func (d *Dao) InsertBuglyIssue(buglyIssue *model.BuglyIssue) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.BuglyIssue{}).Create(buglyIssue).Error)
}

// GetBuglyIssuesByFilterSQL  Get Bugly Issues By Filter SQL.
func (d *Dao) GetBuglyIssuesByFilterSQL(issueFilterSQL string) (buglyIssues []*model.BuglyIssue, err error) {
	if err = d.db.Raw(issueFilterSQL).Order("id asc").Find(&buglyIssues).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// UpdateBuglyIssueTapdBugID Update Issue Record Tapd Bug ID.
func (d *Dao) UpdateBuglyIssueTapdBugID(id int64, tapdBugID string) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyIssue{}).Where("id=?", id).Update("tapd_bug_id", tapdBugID).Error)
}

// FindBuglyIssues Find Bugly Issues.
func (d *Dao) FindBuglyIssues(req *model.QueryBuglyIssueRequest) (total int64, buglyIssues []*model.BuglyIssue, err error) {
	gDB := d.db.Model(&model.BuglyIssue{})

	if req.IssueNo != "" {
		gDB = gDB.Where("issue_no = ?", req.IssueNo)
	}
	if req.Title != "" {
		gDB = gDB.Where("title like ?", _wildcards+req.Title+_wildcards)
	}
	if req.ExceptionMsg != "" {
		gDB = gDB.Where("exception_msg like ?", _wildcards+req.ExceptionMsg+_wildcards)
	}
	if req.KeyStack != "" {
		gDB = gDB.Where("key_stack like ?", _wildcards+req.KeyStack+_wildcards)
	}
	if req.Detail != "" {
		gDB = gDB.Where("detail like ?", _wildcards+req.Detail+_wildcards)
	}
	if req.Tags != "" {
		gDB = gDB.Where("tags like ?", _wildcards+req.Tags+_wildcards)
	}
	if req.Version != "" {
		gDB = gDB.Where("version like ?", _wildcards+req.Version+_wildcards)
	}
	if req.ProjectID != "" {
		gDB = gDB.Where("project_id like ?", _wildcards+req.ProjectID+_wildcards)
	}
	if req.TapdBugID != "" {
		gDB = gDB.Where("tapd_bug_id = ?", req.TapdBugID)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("mtime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&buglyIssues).Error)
	return
}

// GetBuglyIssuesHasInTapd Get Bugly Issues Has In Tapd.
func (d *Dao) GetBuglyIssuesHasInTapd() (buglyIssues []*model.BuglyIssue, err error) {
	if err = d.db.Where("tapd_bug_id<>''").Find(&buglyIssues).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}
