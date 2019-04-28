package v2

import (
	"database/sql"

	"go-common/app/infra/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// BuildsByAppID get builds by app id.
func (d *Dao) BuildsByAppID(appID int64) (builds []string, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Select("name").Model(&model.Build{}).Where("app_id = ? ", appID).Rows(); err != nil {
		log.Error("BuildsByAppID(%v) error(%v)", appID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var build string
		if err = rows.Scan(&build); err != nil {
			log.Error("BuildsByAppID(%v) error(%v)", appID, err)
			return
		}
		builds = append(builds, build)
	}
	if len(builds) == 0 {
		err = ecode.NothingFound
	}
	return
}

// BuildsByAppIDs get builds by app id.
func (d *Dao) BuildsByAppIDs(appIDs []int64) (builds []string, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Select("name").Model(&model.Build{}).Where("app_id in (?) ", appIDs).Rows(); err != nil {
		log.Error("BuildsByAppIDs(%v) error(%v)", appIDs, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var build string
		if err = rows.Scan(&build); err != nil {
			log.Error("BuildsByAppIDs(%v) error(%v)", appIDs, err)
			return
		}
		builds = append(builds, build)
	}
	if len(builds) == 0 {
		err = ecode.NothingFound
	}
	return
}

// TagID get TagID by ID.
func (d *Dao) TagID(appID int64, build string) (tagID int64, err error) {
	row := d.DB.Select("tag_id").Where("app_id =? and name= ?", appID, build).Model(&model.Build{}).Row()
	if err = row.Scan(&tagID); err != nil {
		log.Error("TagID(%v) error(%v)", build, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

// BuildID get build by ID.
func (d *Dao) BuildID(appID int64, build string) (buildID int64, err error) {
	row := d.DB.Select("id").Where("app_id =? and name= ?", appID, build).Model(&model.Build{}).Row()
	if err = row.Scan(&buildID); err != nil {
		log.Error("buildID(%v) error(%v)", buildID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}
