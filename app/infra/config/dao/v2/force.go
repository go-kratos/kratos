package v2

import (
	"database/sql"

	"go-common/app/infra/config/model"
	"go-common/library/log"
)

// Force get force by ID.
func (d *Dao) Force(appID int64, hostname string) (version int64, err error) {
	row := d.DB.Select("version").Where("app_id = ? and hostname = ?", appID, hostname).Model(&model.Force{}).Row()
	if err = row.Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			version = 0
		} else {
			log.Error("version(%v) error(%v)", version, err)
		}
	}
	return
}
