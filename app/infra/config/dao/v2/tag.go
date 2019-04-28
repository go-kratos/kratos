package v2

import (
	"database/sql"
	"go-common/app/infra/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

//Tags get builds by app id.
func (d *Dao) Tags(appID int64) (tags []*model.ReVer, err error) {
	rows, err := d.DB.Select("id,mark").Where("app_id = ? ", appID).Model(&model.DBTag{}).Rows()
	if err != nil {
		log.Error("Tags(%v) error(%v)", appID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		reVer := &model.ReVer{}
		if err = rows.Scan(&reVer.Version, &reVer.Remark); err != nil {
			log.Error("Tags(%v) error(%v)", appID, err)
			return
		}
		tags = append(tags, reVer)
	}
	if len(tags) == 0 {
		err = ecode.NothingFound
	}
	return
}

//ConfIDs get tag by id.
func (d *Dao) ConfIDs(ID int64) (ids []int64, err error) {
	tag := &model.DBTag{}
	if err = d.DB.First(tag, ID).Error; err != nil {
		log.Error("ConfIDs(%v) error(%v)", ID, err)
		return
	}
	ids, _ = xstr.SplitInts(tag.ConfigIDs)
	if len(ids) == 0 {
		err = ecode.NothingFound
	}
	return
}

// TagForce get force by tag.
func (d *Dao) TagForce(ID int64) (force int8, err error) {
	row := d.DB.Select("`force`").Where("id = ?", ID).Model(&model.DBTag{}).Row()
	if err = row.Scan(&force); err != nil {
		log.Error("tagID(%v) error(%v)", ID, err)
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		}
	}
	return
}

// LastForce ...
func (d *Dao) LastForce(appID, buildID int64) (lastForce int64, err error) {
	row := d.DB.Select("id").Where("`app_id` = ? and `build_id` = ? and `force` = 1", appID, buildID).Model(&model.DBTag{}).Row()
	if err = row.Scan(&lastForce); err != nil {
		log.Error("lastForce(%v) error(%v)", lastForce, err)
		if err == sql.ErrNoRows {
			err = nil
			lastForce = 0
		}
	}
	return
}

//TagAll ...
func (d *Dao) TagAll(tagID int64) (tag *model.DBTag, err error) {
	tag = &model.DBTag{}
	if err = d.DB.First(tag, tagID).Error; err != nil {
		log.Error("tagID(%v) error(%v)", tagID, err)
	}
	return
}
