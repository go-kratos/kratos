package v2

import (
	"database/sql"

	"go-common/app/infra/config/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// UpdateConfValue update config state/
func (d *Dao) UpdateConfValue(ID int64, value string) (err error) {
	err = d.DB.Model(&model.Config{ID: ID}).Where("state=?", model.ConfigIng).Update("comment", value).Error
	return
}

// UpdateConfState update config state/
func (d *Dao) UpdateConfState(ID int64, state int8) (err error) {
	err = d.DB.Model(&model.Config{ID: ID}).Update("state", state).Error
	return
}

// ConfigsByIDs get Config by IDs.
func (d *Dao) ConfigsByIDs(ids []int64) (confs []*model.Value, err error) {
	var rows *sql.Rows
	if rows, err = d.DB.Where(ids).Select("id,name,comment").Where("state = ?", model.ConfigEnd).Model(&model.Config{}).Rows(); err != nil {
		log.Error("ConfigsByIDs(%v) error(%v)", ids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		v := new(model.Value)
		if err = rows.Scan(&v.ConfigID, &v.Name, &v.Config); err != nil {
			log.Error("ConfigsByIDs(%v) error(%v)", ids, err)
			return
		}
		confs = append(confs, v)
	}
	if len(confs) == 0 {
		err = ecode.NothingFound
	}
	return
}
