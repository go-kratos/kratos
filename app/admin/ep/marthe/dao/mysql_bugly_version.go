package dao

import (
	"database/sql"
	"fmt"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

const (
	_versionInnerJoinProjectSql      = "select a.id,a.version,a.bugly_project_id,a.action,a.task_status,a.update_by,a.ctime,a.mtime,b.project_name,b.exception_type from bugly_versions as a inner join bugly_projects as b on a.bugly_project_id = b.id"
	_versionInnerJoinProjectSqlCount = "select count(-1) as totalCount from bugly_versions as a inner join bugly_projects as b on a.bugly_project_id = b.id"
	_where                           = "WHERE"
	_and                             = "AND"
)

// InsertBuglyVersion Insert Bugly Version.
func (d *Dao) InsertBuglyVersion(buglyVersion *model.BuglyVersion) error {
	return pkgerr.WithStack(d.db.Create(buglyVersion).Error)
}

// UpdateBuglyVersion Update Bugly Version.
func (d *Dao) UpdateBuglyVersion(buglyVersion *model.BuglyVersion) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyVersion{}).Updates(buglyVersion).Error)
}

// QueryBuglyVersionByVersion Query Bugly Version By Version.
func (d *Dao) QueryBuglyVersionByVersion(version string) (buglyVersion *model.BuglyVersion, err error) {
	buglyVersion = &model.BuglyVersion{}
	if err = d.db.Where("version = ?", version).First(buglyVersion).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// QueryBuglyVersion Query Bugly Version .
func (d *Dao) QueryBuglyVersion(id int64) (buglyVersion *model.BuglyVersion, err error) {
	buglyVersion = &model.BuglyVersion{}
	if err = d.db.Where("id = ?", id).First(buglyVersion).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}

// QueryBuglyVersionList Query Bugly Version List.
func (d *Dao) QueryBuglyVersionList() (versionList []string, err error) {
	var (
		rows *sql.Rows
	)
	sql := "select DISTINCT version from bugly_versions"
	if rows, err = d.db.Raw(sql).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var ver string
		if err = rows.Scan(&ver); err != nil {
			return
		}
		versionList = append(versionList, ver)
	}
	return
}

// FindBuglyProjectVersions Find Bugly Project Versions.
func (d *Dao) FindBuglyProjectVersions(req *model.QueryBuglyVersionRequest) (total int64, buglyProjectVersions []*model.BuglyProjectVersion, err error) {
	var (
		qSQL = _versionInnerJoinProjectSql
		cSQL = _versionInnerJoinProjectSqlCount
		rows *sql.Rows
	)

	if req.UpdateBy != "" || req.ProjectName != "" || req.Action > 0 || req.Version != "" {
		var (
			partSql     string
			logicalWord = _where
		)

		if req.UpdateBy != "" {
			partSql = fmt.Sprintf("%s %s a.update_by = '%s'", partSql, logicalWord, req.UpdateBy)
			logicalWord = _and
		}

		if req.ProjectName != "" {
			partSql = fmt.Sprintf("%s %s b.project_name like '%s'", partSql, logicalWord, _wildcards+req.ProjectName+_wildcards)
			logicalWord = _and
		}

		if req.Action > 0 {
			partSql = fmt.Sprintf("%s %s a.action = %d", partSql, logicalWord, req.Action)
			logicalWord = _and
		}

		if req.Version != "" {
			partSql = fmt.Sprintf("%s %s a.version like '%s'", partSql, logicalWord, _wildcards+req.Version+_wildcards)
			logicalWord = _and
		}

		qSQL = qSQL + partSql
		cSQL = cSQL + partSql
	}

	cDB := d.db.Raw(cSQL)
	if err = pkgerr.WithStack(cDB.Count(&total).Error); err != nil {
		return
	}
	gDB := d.db.Raw(qSQL)
	if rows, err = gDB.Order("a.ctime DESC").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Rows(); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		pv := &model.BuglyProjectVersion{}
		if err = rows.Scan(&pv.ID, &pv.Version, &pv.BuglyProjectID, &pv.Action, &pv.TaskStatus, &pv.UpdateBy, &pv.CTime, &pv.MTime, &pv.ProjectName, &pv.ExceptionType); err != nil {
			return
		}
		buglyProjectVersions = append(buglyProjectVersions, pv)
	}

	return
}

// FindEnableAndReadyVersions Find Enable And Ready Versions.
func (d *Dao) FindEnableAndReadyVersions() (buglyVersions []*model.BuglyVersion, err error) {
	err = pkgerr.WithStack(d.db.Where("action = ? and task_status = ?", model.BuglyVersionActionEnable, model.BuglyVersionTaskStatusReady).Find(&buglyVersions).Error)
	return
}
