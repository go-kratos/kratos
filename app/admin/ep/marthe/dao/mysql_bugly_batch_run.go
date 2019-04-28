package dao

import (
	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertBuglyBatchRuns Insert Bugly Batch Runs.
func (d *Dao) InsertBuglyBatchRuns(buglyBatchRuns []*model.BuglyBatchRun) (err error) {
	tx := d.db.Begin()

	defer func() {
		if err != nil {
			if err == ecode.NothingFound {
				err = nil
			} else {
				err = pkgerr.WithStack(err)
			}
		}
	}()

	if err = tx.Error; err != nil {
		return
	}

	for _, buglyBatchRun := range buglyBatchRuns {
		if err = d.db.Create(buglyBatchRun).Error; err != nil {
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

// InsertBuglyBatchRun Insert Bugly Batch Run.
func (d *Dao) InsertBuglyBatchRun(buglyBatchRun *model.BuglyBatchRun) error {
	return pkgerr.WithStack(d.db.Create(buglyBatchRun).Error)
}

// UpdateBuglyBatchRun Update Bugly Batch Run.
func (d *Dao) UpdateBuglyBatchRun(buglyBatchRun *model.BuglyBatchRun) error {
	return pkgerr.WithStack(d.db.Model(&model.BuglyBatchRun{}).Where("id=?", buglyBatchRun.ID).Updates(buglyBatchRun).Error)
}

// FindBuglyBatchRuns Find Bugly Batch Runs.
func (d *Dao) FindBuglyBatchRuns(req *model.QueryBuglyBatchRunsRequest) (total int64, buglyBatchRuns []*model.BuglyBatchRun, err error) {
	gDB := d.db.Model(&model.BuglyBatchRun{})

	if req.Version != "" {
		gDB = gDB.Where("version=?", req.Version)
	}

	if req.Status != 0 {
		gDB = gDB.Where("status=?", req.Status)
	}

	if req.BatchID != "" {
		gDB = gDB.Where("batch_id=?", req.BatchID)
	}

	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}

	err = pkgerr.WithStack(gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&buglyBatchRuns).Error)
	return
}

// QueryBuglyBatchRunsByStatus Find Bugly Batch Runs By Status.
func (d *Dao) QueryBuglyBatchRunsByStatus(status int) (buglyBatchRuns []*model.BuglyBatchRun, err error) {
	err = pkgerr.WithStack(d.db.Where("status = ?", status).Find(&buglyBatchRuns).Error)
	return
}

// QueryLastSuccessBatchRunByVersion Find Last Success Batch Run By Version.
func (d *Dao) QueryLastSuccessBatchRunByVersion(version string) (buglyBatchRun *model.BuglyBatchRun, err error) {
	buglyBatchRun = &model.BuglyBatchRun{}
	if err = d.db.Where("version = ? and status = ?", version, model.BuglyBatchRunStatusDone).Order("id desc").First(&buglyBatchRun).Error; err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			err = pkgerr.WithStack(err)
		}
	}
	return
}
