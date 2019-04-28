package dao

import (
	"context"

	"go-common/library/log"
)

var (
	_allParamterSQL = "SELECT name, value FROM rating_parameter"
)

// GetAllParamter get all paramter
func (d *Dao) GetAllParamter(c context.Context) (paramters map[string]int64, err error) {
	paramters = make(map[string]int64)
	rows, err := d.db.Query(c, _allParamterSQL)
	if err != nil {
		log.Error("d.db.Query GetAllParamter error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var vaule int64
		err = rows.Scan(&name, &vaule)
		if err != nil {
			log.Error("rows.Scan GetAllParamter error(%v)", err)
			return
		}
		paramters[name] = vaule
	}
	return
}
