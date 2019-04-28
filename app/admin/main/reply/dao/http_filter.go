package dao

import (
	"context"
	"fmt"

	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_selFilteredReply = "select rpid,message FROM reply_filtered WHERE rpid in (%s)"

	_filterSearch = "http://api.bilibili.co/x/admin/filter/origins"
)

// FilterContents get filtered contents from db
func (d *Dao) FilterContents(ctx context.Context, rpMaps map[int64]string) error {
	if len(rpMaps) == 0 {
		return nil
	}
	var rpids []int64
	for k := range rpMaps {
		rpids = append(rpids, k)
	}
	rows, err := d.db.Query(ctx, fmt.Sprintf(_selFilteredReply, xstr.JoinInts(rpids)))
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var message string
		if err = rows.Scan(&id, &message); err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				log.Error("row.Scan error(%v)", err)
				return err
			}
		}
		rpMaps[id] = message
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return err
	}
	return nil
}
