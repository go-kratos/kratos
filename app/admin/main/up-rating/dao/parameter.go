package dao

import (
	"context"
	"fmt"

	"go-common/library/log"
)

var (
	// select
	_allParameterSQL = "SELECT name, value FROM rating_parameter"

	// insert
	_inParameterSQL = "INSERT INTO rating_parameter(name, value, remark) VALUES %s ON DUPLICATE KEY UPDATE value=VALUES(value), remark=VALUES(remark)"
)

// GetAllParameter get all parameter
func (d *Dao) GetAllParameter(c context.Context) (parameters map[string]int64, err error) {
	parameters = make(map[string]int64)
	rows, err := d.db.Query(c, _allParameterSQL)
	if err != nil {
		log.Error("d.db.Query GetAllParameter error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var vaule int64
		err = rows.Scan(&name, &vaule)
		if err != nil {
			log.Error("rows.Scan GetAllParameter error(%v)", err)
			return
		}
		parameters[name] = vaule
	}
	return
}

// InsertParameter insert vals into rating_parameter
func (d *Dao) InsertParameter(c context.Context, val string) (rows int64, err error) {
	if val == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inParameterSQL, val))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
