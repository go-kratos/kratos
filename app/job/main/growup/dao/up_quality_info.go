package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	_upQualitySQL = "SELECT id,mid,quality_value FROM %s WHERE id > ? AND is_deleted = 0 ORDER BY id LIMIT ?"
)

// GetUpQuality get up_quality_info
func (d *Dao) GetUpQuality(c context.Context, table string, id int64, limit int) (up []*model.UpQuality, last int64, err error) {
	up = make([]*model.UpQuality, 0)
	if table == "" {
		err = fmt.Errorf("ERROR: table is null")
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_upQualitySQL, table), id, limit)
	if err != nil {
		log.Error("GetUpQuality d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := &model.UpQuality{}
		err = rows.Scan(&last, &u.MID, &u.Quality)
		if err != nil {
			log.Error("GetUpQuality rows.Scan error(%v)", err)
			return
		}
		up = append(up, u)
	}

	err = rows.Err()
	return
}
