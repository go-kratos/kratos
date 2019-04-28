package dao

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/database/sql"
	"time"

	"github.com/pkg/errors"
)

const (
	_insertStreamChangeLog = "INSERT INTO %s (room_id, from_origin,to_origin, source,operate_name,reason) VALUES (?,?,?,?,?,?)"
	_selectStreamChangelog = "SELECT room_id, from_origin, to_origin, source, operate_name,reason, ctime FROM %s where room_id = ? ORDER BY mtime DESC LIMIT ?"
)

// InsertChangeLog 插入日志
func (d *Dao) InsertChangeLog(c context.Context, change *model.StreamChangeLog) error {
	// 判断当前年和月
	now := time.Now().Format("200601")
	tableName := fmt.Sprintf("stream_change_log_%s", now)
	_, err := d.db.Exec(c, fmt.Sprintf(_insertStreamChangeLog, tableName), change.RoomID, change.FromOrigin, change.ToOrigin, change.Source, change.OperateName, change.Reason)

	return err
}

// GetChangeLogByRoomID 查询日志
func (d *Dao) GetChangeLogByRoomID(c context.Context, rid int64, limit int64) (infos []*model.StreamChangeLog, err error) {
	now := time.Now().Format("200601")
	tableName := fmt.Sprintf("stream_change_log_%s", now)

	var rows *sql.Rows
	if rows, err = d.db.Query(c, fmt.Sprintf(_selectStreamChangelog, tableName), rid, limit); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.StreamChangeLog)
		if err = rows.Scan(&info.RoomID, &info.FromOrigin, &info.ToOrigin, &info.Source, &info.OperateName, &info.Reason, &info.CTime); err != nil {
			err = errors.WithStack(err)
			infos = nil
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}
