package tag

import (
	"context"
	"fmt"
)

const (
	// delete
	_delArchiveRatioSQL = "DELETE FROM av_charge_ratio WHERE ctype = ? LIMIT ?"
	_delUpRatioSQL      = "DELETE FROM up_charge_ratio WHERE ctype = ? LIMIT ?"

	// insert
	_addAvRatioSQL = "INSERT INTO av_charge_ratio(tag_id, av_id, ratio, adjust_type, ctype) VALUES %s"
	_addUpRatioSQL = "INSERT INTO up_charge_ratio(tag_id, mid, ratio, adjust_type, ctype) VALUES %s"
)

// DelArchiveRatio delete all av_charge_ratio by type
func (d *Dao) DelArchiveRatio(c context.Context, ctype int, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delArchiveRatioSQL, ctype, limit)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// DelUpRatio delete all up_charge_ratio
func (d *Dao) DelUpRatio(c context.Context, ctype int, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delUpRatioSQL, ctype, limit)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// InsertAvRatio insert av charge ratio
func (d *Dao) InsertAvRatio(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addAvRatioSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// InsertUpRatio insert av charge ratio
func (d *Dao) InsertUpRatio(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addUpRatioSQL, values))
	if err != nil {
		return
	}
	return res.RowsAffected()
}
