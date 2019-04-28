package dao

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"

	"github.com/pkg/errors"
)

const (
	_getOfficialStreamByName           = "SELECT id, room_id, src, `name`, `key`, up_rank, down_rank, `status` FROM `sv_ls_stream` WHERE name = ?"
	_getOfficialStreamByRoomID         = "SELECT id, room_id, src, `name`, `key`, up_rank, down_rank, `status` FROM `sv_ls_stream` WHERE room_id = ?"
	_getMultiOfficalStreamByRID        = "SELECT id, room_id, src, `name`, `key`, up_rank, down_rank, `status` FROM `sv_ls_stream` WHERE room_id = %d"
	_insertOfficialStream              = "INSERT INTO `sv_ls_stream` (room_id, `name`, `src`,`key`, `status`, up_rank, down_rank, last_status_updated_at, created_at, updated_at)  VALUES (?,?,?,?,?,?,?,?,?,?)"
	_updateUpOfficialStreamStatus      = "UPDATE `sv_ls_stream` SET `up_rank` = 1,`last_status_updated_at` = CURRENT_TIMESTAMP  WHERE  `room_id` = ? and `src` = ?"
	_updateForwardOfficialStreamStatus = "UPDATE `sv_ls_stream` SET `up_rank` = 0,`last_status_updated_at` = CURRENT_TIMESTAMP  WHERE  `room_id` = ? and `src` != ?"
	_updateOfficalStreamUpRankStatus   = "UPDATE `sv_ls_stream` SET `up_rank` = ?,`last_status_updated_at` = CURRENT_TIMESTAMP  WHERE  `room_id` = ? AND `up_rank` = ?;"
)

// GetOfficialStreamByName 根据流名查流信息， 可以获取多条记录
func (d *Dao) GetOfficialStreamByName(c context.Context, name string) (infos []*model.OfficialStream, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getOfficialStreamByName, name); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.OfficialStream)
		if err = rows.Scan(&info.ID, &info.RoomID,
			&info.Src, &info.Name, &info.Key,
			&info.UpRank, &info.DownRank, &info.Status); err != nil {
			log.Warn("sv_ls_stream sql err = %v", err)
			err = errors.WithStack(err)
			infos = nil
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}

// GetOfficialStreamByRoomID 根据roomid查询流信息, 可以获取多条记录
func (d *Dao) GetOfficialStreamByRoomID(c context.Context, rid int64) (infos []*model.OfficialStream, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _getOfficialStreamByRoomID, rid); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.OfficialStream)
		if err = rows.Scan(&info.ID, &info.RoomID,
			&info.Src, &info.Name, &info.Key,
			&info.UpRank, &info.DownRank, &info.Status); err != nil {
			log.Warn("sv_ls_stream sql err = %v", err)
			err = errors.WithStack(err)
			infos = nil
			return
		}
		infos = append(infos, info)
	}
	err = rows.Err()

	// 查询多个数据，不会报错， 只能判断为空
	if err == nil && len(infos) == 0 && rid < 10000 {
		infos = append(infos, &model.OfficialStream{
			RoomID: rid,
			Name:   "miss",
			Src:    32,
			UpRank: 1,
			Key:    "miss",
		})
	}
	return
}

// GetMultiOfficalStreamByRID 批量获取
func (d *Dao) GetMultiOfficalStreamByRID(c context.Context, rids []int64) (infos []*model.OfficialStream, err error) {
	len := len(rids)
	muSql := ""
	for i := 0; i < len; i++ {
		ss := fmt.Sprintf(_getMultiOfficalStreamByRID, rids[i])
		if i == 0 {
			muSql = fmt.Sprintf("%s%s", muSql, ss)
		} else {
			muSql = fmt.Sprintf("%s UNION %s", muSql, ss)
		}
	}

	var rows *sql.Rows
	if rows, err = d.db.Query(c, muSql); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		info := new(model.OfficialStream)
		if err = rows.Scan(&info.ID, &info.RoomID,
			&info.Src, &info.Name, &info.Key,
			&info.UpRank, &info.DownRank, &info.Status); err != nil {
			log.Warn("sv_ls_stream sql err = %v", err)
			if err == sql.ErrNoRows {
				continue
			} else {
				err = errors.WithStack(err)
				infos = nil
				return infos, err
			}
		}
		infos = append(infos, info)
	}
	err = rows.Err()
	return
}

// CreateOfficialStream 创建正式流
func (d *Dao) CreateOfficialStream(c context.Context, infos []*model.OfficialStream) (err error) {
	tx, err := d.db.Begin(c)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	ts := time.Now().Format("2006-01-02 15:04:05")
	for _, v := range infos {
		if _, err = d.stmtLegacyStreamCreate.Exec(c, v.RoomID, v.Name, v.Src, v.Key, v.Status, v.UpRank, v.DownRank, ts, ts, ts); err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

// UpdateOfficialStreamStatus 切换cdn时，更新流状态
func (d *Dao) UpdateOfficialStreamStatus(c context.Context, rid int64, src int8) (err error) {
	// 事务操作， 同时操作多条记录，任何一条失败均回滚
	tx, err := d.db.Begin(c)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if _, err = d.stmtLegacyStreamEnableNewUpRank.Exec(c, rid, src); err != nil {
		return err
	}

	if _, err = d.stmtLegacyStreamDisableUpRank.Exec(c, rid, src); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

// UpdateOfficialUpRankStatus 清理互推标准，更新up_rank
func (d *Dao) UpdateOfficialUpRankStatus(c context.Context, rid int64, whereSrc int8, toSrc int8) error {
	res, err := d.stmtLegacyStreamClearStreamFoward.Exec(c, toSrc, rid, whereSrc)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()

	return err
}
