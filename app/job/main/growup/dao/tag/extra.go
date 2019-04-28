package tag

import (
	"context"
	"fmt"

	"go-common/library/log"
)

const (
	_inUpTagYearSQL     = "INSERT INTO up_tag_year(mid, %s) VALUES %s ON DUPLICATE KEY UPDATE %s = %s + VALUES(%s)"
	_inUpYearAccountSQL = "INSERT INTO up_tag_year(mid, total_income) VALUES %s ON DUPLICATE KEY UPDATE total_income = VALUES(total_income)"
)

// InsertUpTagYear insert up_tag_year
func (d *Dao) InsertUpTagYear(c context.Context, vals string, col string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpTagYearSQL, col, vals, col, col, col))
	if err != nil {
		log.Error("dao.InsertUpTagYear exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// InsertUpYearAccount .
func (d *Dao) InsertUpYearAccount(c context.Context, vals string) (rows int64, err error) {
	if vals == "" {
		return
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_inUpYearAccountSQL, vals))
	if err != nil {
		log.Error("InsertUpYearAccount exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}
