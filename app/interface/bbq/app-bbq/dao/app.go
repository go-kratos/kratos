package dao

import (
	"context"
	"time"

	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/database/sql"
)

const (
	_queryNewAppVersion = "select `id`, `platform`, `ver_name`, `ver_code`, `title`, `content`, `download`, `md5`, `size`, `force`, `status` from `app_package` where platform = ? and ver_code > ? and status = 1 order by ver_code desc;"
	_queryAppPackage    = "select `id`, `platform`, `ver_name`, `ver_code`, `title`, `content`, `download`, `md5`, `size`, `force`, `status`, `ctime` from `app_package` where status>0;"
	_queryAppResource   = "select `id`, `platform`, `name`, `version`, `md5`, `download`, `status`, `start_time`, `end_time` from `app_resource` where `platform` in (0, ?) and `version` = ? and `status` = ? and `end_time` > ?;"
)

// FetchNewAppVersion .
func (d *Dao) FetchNewAppVersion(c context.Context, platform int, vCode int) (result *model.AppVersion, err error) {
	result = &model.AppVersion{}
	err = d.db.QueryRow(c, _queryNewAppVersion, platform, vCode).Scan(&result.ID, &result.Platform, &result.Name, &result.Code, &result.Title, &result.Content, &result.Download, &result.MD5, &result.Size, &result.Force, &result.Status)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// FetchAppPackage .
func (d *Dao) FetchAppPackage(c context.Context) (result []*v1.AppPackage, err error) {
	rows, err := d.db.Query(c, _queryAppPackage)
	for rows.Next() {
		tmp := &v1.AppPackage{}
		err = rows.Scan(&tmp.ID, &tmp.Platform, &tmp.VersionName, &tmp.VersionCode, &tmp.Title, &tmp.Content, &tmp.Download, &tmp.MD5, &tmp.Size, &tmp.Force, &tmp.Status, &tmp.CTime)
		if err != nil {
			continue
		}
		result = append(result, tmp)
	}
	return
}

// FetchAppResource .
func (d *Dao) FetchAppResource(c context.Context, plat int, ver int) (result []*model.AppResource, err error) {
	result = make([]*model.AppResource, 0)
	rows, err := d.db.Query(c, _queryAppResource, plat, ver, 1, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return
	}
	for rows.Next() {
		tmp := new(model.AppResource)
		err = rows.Scan(&tmp.ID, &tmp.Platform, &tmp.Name, &tmp.Code, &tmp.MD5, &tmp.Download, &tmp.Status, &tmp.StartTime, &tmp.EndTime)
		if err != nil {
			continue
		}
		result = append(result, tmp)
	}
	return
}
