package like

import (
	"context"
	"database/sql"

	"go-common/app/job/main/activity/model/like"

	"github.com/pkg/errors"
)

const (
	_webDataCntSQL  = "SELECT COUNT(1) AS cnt FROM act_web_data WHERE state = 1 AND vid = ?"
	_webDataListSQL = "SELECT id,vid,data FROM act_web_data WHERE state= 1 AND vid = ? ORDER BY id LIMIT ?,?"
)

// WebDataCnt get web data count.
func (d *Dao) WebDataCnt(c context.Context, vid int64) (count int, err error) {
	row := d.db.QueryRow(c, _webDataCntSQL, vid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			err = errors.Wrapf(err, "WebDataCnt:QueryRow(%d)", vid)
		}
	}
	return
}

// WebDataList get web data list by vid.
func (d *Dao) WebDataList(c context.Context, vid int64, offset, limit int) (list []*like.WebData, err error) {
	rows, err := d.db.Query(c, _webDataListSQL, vid, offset, limit)
	if err != nil {
		err = errors.Wrapf(err, "WebDataList:d.db.Query(%d,%d,%d)", vid, offset, limit)
		return
	}
	defer rows.Close()
	for rows.Next() {
		n := new(like.WebData)
		if err = rows.Scan(&n.ID, &n.Vid, &n.Data); err != nil {
			err = errors.Wrapf(err, "WebDataList:row.Scan row (%d,%d,%d)", vid, offset, limit)
			return
		}
		list = append(list, n)
	}
	if err = rows.Err(); err != nil {
		err = errors.Wrapf(err, "LikeList:rowsErr(%d,%d,%d)", vid, offset, limit)
	}
	return
}
