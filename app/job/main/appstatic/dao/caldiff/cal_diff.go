package caldiff

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/job/main/appstatic/model"
	"go-common/library/log"
)

const (
	_calDiffFmt = "SELECT id,`name`,type, md5, size, url, resource_id, file_type, from_ver FROM resource_file " +
		"WHERE file_type = ? AND url = ? AND is_deleted = 0 %s ORDER BY id DESC %s"
	_saveFile  = "UPDATE resource_file SET url = ?, file_type = ?, md5 = ?, size = ?, `name` = ? WHERE id = ?"
	_parseRes  = "SELECT id, `name`, version, pool_id FROM resource WHERE %s"
	_parseFile = "SELECT id, `name`, type, md5, size, url, resource_id, file_type, from_ver FROM resource_file WHERE" +
		" %s AND is_deleted = 0 %s"
	_updateStatus = "UPDATE resource_file SET file_type = ? WHERE id = ?"
	_diffPkg      = 1
)

var (
	_calDiffNew   = fmt.Sprintf(_calDiffFmt, " AND ctime = mtime", "LIMIT 1")
	_parseResID   = fmt.Sprintf(_parseRes, " id = ?")
	_parseResVer  = fmt.Sprintf(_parseRes, " pool_id = ? AND version = ?")
	_getReadyFile = fmt.Sprintf(_parseFile, " resource_id = ? AND file_type = ? AND url != ?", "LIMIT 1")
)

// UpdateStatus updates the file's status
func (d *Dao) UpdateStatus(c context.Context, status int, id int) (err error) {
	if _, err = d.db.Exec(c, _updateStatus, status, id); err != nil {
		log.Error("UpdateStatus ID %d, Err %v", id, err)
	}
	return
}

// ReadyFile takes the already generated file
func (d *Dao) ReadyFile(c context.Context, resID int, ftype int) (file *model.ResourceFile, err error) {
	file = &model.ResourceFile{}
	row := d.db.QueryRow(c, _getReadyFile, resID, ftype, "")
	if err = row.Scan(&file.ID, &file.Name, &file.Type, &file.Md5, &file.Size, &file.URL, &file.ResourceID, &file.FileType, &file.FromVer); err != nil {
		log.Error("db.QueryRow(%s) (%d,%d) error(%v)", _getReadyFile, resID, ftype, err)
	}
	return
}

// ParseResVer takes one resource info
func (d *Dao) ParseResVer(c context.Context, poolID int, version int) (res *model.Resource, err error) {
	res = &model.Resource{}
	row := d.db.QueryRow(c, _parseResVer, poolID, version)
	// "SELECT id, `name`, version, pool_id FROM resource WHERE pool_id = ? AND version = ?"
	if err = row.Scan(&res.ID, &res.Name, &res.Version, &res.PoolID); err != nil {
		log.Error("db.QueryRow(%s) (%d,%d) error(%v)", _parseResVer, poolID, version, err)
	}
	return
}

// ParseResID takes one resource info
func (d *Dao) ParseResID(c context.Context, resID int) (res *model.Resource, err error) {
	res = &model.Resource{}
	row := d.db.QueryRow(c, _parseResID, resID)
	// "SELECT id, `name`, version, pool_id FROM resource WHERE id = ?"
	if err = row.Scan(&res.ID, &res.Name, &res.Version, &res.PoolID); err != nil {
		log.Error("db.QueryRow(%s) (%d) error(%v)", _parseResID, resID, err)
	}
	return
}

// DiffNew picks the recently created diff packages
func (d *Dao) DiffNew(c context.Context) (file *model.ResourceFile, err error) {
	file = &model.ResourceFile{}
	row := d.db.QueryRow(c, _calDiffNew, _diffPkg, "")
	if err = row.Scan(&file.ID, &file.Name, &file.Type, &file.Md5, &file.Size, &file.URL, &file.ResourceID, &file.FileType, &file.FromVer); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			file = nil
		} else {
			log.Error("db.QueryRow(%s) error(%v)", _calDiffNew, err)
			return
		}
	}
	return
}

// DiffRetry picks the recently created diff packages
func (d *Dao) DiffRetry(c context.Context) (file *model.ResourceFile, err error) {
	file = &model.ResourceFile{}
	query := fmt.Sprintf(_calDiffFmt, " AND ctime != mtime AND mtime < date_sub(now(), INTERVAL "+d.c.Cfg.Diff.Retry+")", "LIMIT 1")
	row := d.db.QueryRow(c, query, _diffPkg, "")
	if err = row.Scan(&file.ID, &file.Name, &file.Type, &file.Md5, &file.Size, &file.URL, &file.ResourceID, &file.FileType, &file.FromVer); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			file = nil
		} else {
			log.Error("db.QueryRow(%s) error(%v)", _calDiffNew, err)
			return
		}
	}
	return
}

// SaveFile saves the file info
func (d *Dao) SaveFile(c context.Context, fileID int, file *model.FileInfo) (err error) {
	// "UPDATE resource_file SET url = ?, file_type = ?, md5 = ?, size = ?, `name` = ? WHERE id = ?"
	if _, err = d.db.Exec(c, _saveFile, file.URL, _diffPkg, file.Md5, file.Size, file.Name, fileID); err != nil {
		log.Error("SaveFile ID %d, Err %v", fileID, err)
	}
	return
}
