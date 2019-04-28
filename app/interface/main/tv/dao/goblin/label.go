package goblin

import (
	"context"

	gbmdl "go-common/app/interface/main/tv/model/goblin"
)

const (
	_labelSQL = "SELECT id,name,param, param_name,value FROM tv_label WHERE category = ? AND cat_type = ? AND valid = 1 AND deleted = 0" +
		" ORDER BY position,id ASC "
)

// Label picks one category's label
func (d *Dao) Label(c context.Context, category, catType int) (res []*gbmdl.Label, err error) {
	rows, err := d.db.Query(c, _labelSQL, category, catType)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		li := &gbmdl.Label{}
		if err = rows.Scan(&li.ID, &li.Name, &li.Param, &li.ParamName, &li.Value); err != nil {
			return
		}
		res = append(res, li)
	}
	err = rows.Err()
	return
}
