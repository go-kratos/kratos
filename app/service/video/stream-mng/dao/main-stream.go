package dao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"

	"github.com/pkg/errors"
)

/*
全新推流结构

所需功能:
1. 创建流 （从老表搬运） done
2. 校验流/读取 done
3. 开关播回调
4. 切上行
5. 清理互推标记
*/

const (
	// 创建流
	_insertMainStream = "INSERT INTO `main_stream` (room_id, stream_name, `key`, default_vendor, options) VALUES (?, ?, ?, ?, ?);"
	// 读取流
	_getMainStreamWithoutConds = "SELECT `room_id`, `stream_name`, `key`, `default_vendor`, `origin_upstream`, `streaming`, `last_stream_time`, `options` from `main_stream` WHERE "
	_getMultiMainStreamByRID   = "SELECT `room_id`, `stream_name`, `key`, `default_vendor`, `origin_upstream`, `streaming`, `last_stream_time`, `options` from `main_stream` WHERE room_id = %d"
	// 切上行
	_changeDefaultVendor = "UPDATE `main_stream` SET `default_vendor` = ? WHERE `room_id` = ? AND status = 1"
	// 切options
	_changeOptions = "UPDATE `main_stream` SET `options` = ? WHERE `room_id` = ? AND status = 1 AND `options` = ?"
	// 清理互推
	_clearAllStreaming = "UPDATE `main_stream` SET `origin_upstream` = 0, `streaming` = 0, `options` = ? WHERE `room_id` = ? AND `options` = ? AND status = 1"

	// 开关回调
	_notifyMainStreamOrigin      = "UPDATE `main_stream` SET `origin_upstream` = ?, `streaming` = ? WHERE `room_id` = ? and `streaming` = ? and status = 1 limit 1"
	_notifyMainStreamOriginClose = "UPDATE `main_stream` SET `options` = ?,`origin_upstream` = 0, `streaming` = 0, `last_stream_time` = CURRENT_TIMESTAMP WHERE `room_id` = ? AND `options` = ? AND `status` = 1 limit 1"
	_notifyMainStreamForward     = "UPDATE `main_stream` SET `streaming` = ? WHERE `room_id` = ? and `streaming` = ? and status = 1 limit 1"
)

// CreateNewStream used to create new Stream record
func (d *Dao) CreateNewStream(c context.Context, stream *model.MainStream) (*model.MainStream, error) {
	if stream.RoomID <= 0 {
		return stream, fmt.Errorf("room id can not be empty")
	}
	if stream.StreamName == "" {
		return stream, fmt.Errorf("stream name can not be empty")
	}
	if stream.Key == "" {
		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%s%d", stream.StreamName, time.Now().Nanosecond())))
		stream.Key = hex.EncodeToString(h.Sum(nil))
	}
	if stream.DefaultVendor == 0 {
		stream.DefaultVendor = 1
	}
	res, err := d.stmtMainStreamCreate.Exec(c, stream.RoomID, stream.StreamName, stream.Key, stream.DefaultVendor, stream.Options)
	if err != nil {
		return stream, err
	}
	stream.ID, err = res.LastInsertId()
	return stream, nil
}

// GetMainStreamFromDB 从DB中读取流信息
// roomID 和 streamName 可以只传一个，传哪个就用哪个查询，否则必须两者对应
func (d *Dao) GetMainStreamFromDB(c context.Context, roomID int64, streamName string) (*model.MainStream, error) {
	if roomID <= 0 && streamName == "" {
		return nil, errors.New("roomID and streamName cannot be empty at SAME time")
	}
	var row *sql.Row
	if roomID > 0 && streamName != "" {
		q := fmt.Sprintf("%s `room_id` = ? AND `stream_name` = ? AND status = 1", _getMainStreamWithoutConds)
		row = d.db.QueryRow(c, q, roomID, streamName)
	} else if roomID > 0 && streamName == "" {
		q := fmt.Sprintf("%s `room_id` = ? AND status = 1", _getMainStreamWithoutConds)
		row = d.db.QueryRow(c, q, roomID)
	} else if roomID <= 0 && streamName != "" {
		q := fmt.Sprintf("%s `stream_name` = ? AND status = 1", _getMainStreamWithoutConds)
		row = d.db.QueryRow(c, q, streamName)
	}

	stream := new(model.MainStream)
	err := row.Scan(&stream.RoomID, &stream.StreamName, &stream.Key,
		&stream.DefaultVendor, &stream.OriginUpstream, &stream.Streaming,
		&stream.LastStreamTime, &stream.Options)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// GetMultiMainStreamFromDB 批量从main-stream读取
func (d *Dao) GetMultiMainStreamFromDB(c context.Context, rids []int64) (mainStream []*model.MainStream, err error) {
	len := len(rids)
	muSql := ""
	for i := 0; i < len; i++ {
		ss := fmt.Sprintf(_getMultiMainStreamByRID, rids[i])
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
		stream := new(model.MainStream)
		if err = rows.Scan(&stream.RoomID, &stream.StreamName, &stream.Key,
			&stream.DefaultVendor, &stream.OriginUpstream, &stream.Streaming,
			&stream.LastStreamTime, &stream.Options); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			err = errors.WithStack(err)
			return
		}
		mainStream = append(mainStream, stream)
	}

	err = rows.Err()
	return
}

// ChangeDefaultVendor 切换默认上行
func (d *Dao) ChangeDefaultVendor(c context.Context, roomID int64, newVendor int64) error {
	if roomID <= 0 {
		return errors.New("invalid roomID")
	}
	if _, ok := common.BitwiseMapName[newVendor]; !ok {
		return errors.New("invalid vendor")
	}
	_, err := d.stmtMainStreamChangeDefaultVendor.Exec(c, newVendor, roomID)
	return err
}

// ChangeMainStreamOptions 切换Options
func (d *Dao) ChangeMainStreamOptions(c context.Context, roomID int64, newOptions int64, options int64) error {
	if roomID <= 0 {
		return errors.New("invalid roomID")
	}
	_, err := d.stmtMainStreamChangeOptions.Exec(c, newOptions, roomID, options)
	return err
}

// ClearMainStreaming 清理互推标记
func (d *Dao) ClearMainStreaming(c context.Context, roomID int64, newoptions int64, options int64) error {
	if roomID <= 0 {
		return errors.New("invalid roomID")
	}
	_, err := d.stmtMainStreamClearAllStreaming.Exec(c, newoptions, roomID, options)
	return err
}

// MainStreamNotify 开关播回调
// @param roomID 房间号
// @param vendor 上行 CDN 位
// @param isOpen 是否是开播 true 开播 false 关播
// @param isOrigin 是否是原始上行 true 是 false 转推
func (d *Dao) MainStreamNotify(c context.Context, roomID, vendor int64, isOpen bool, isOrigin bool, options int64, newoptions int64) error {
	if _, ok := common.BitwiseMapName[vendor]; !ok {
		return fmt.Errorf("Unknow vendor %d", vendor)
	}
	log.Infov(c, log.KV("roomID", roomID), log.KV("vendor", vendor), log.KV("isOpen", isOpen), log.KV("isOrigin", isOrigin), log.KV("options", options), log.KV("newoptions", newoptions))
	// "UPDATE `main_stream` SET `origin_upstream` = ?, `streaming` = ? WHERE `room_id` = ? AND `streaming` = ? AND `origin_upstream` = 0 and status = 1 limit 1"
	// "UPDATE `main_stream` SET `streaming` = ? WHERE `room_id` = ? and `streaming` = ? and status = 1 limit 1"

	ms, err := d.GetMainStreamFromDB(c, roomID, "")
	if ms == nil || err != nil {
		return fmt.Errorf("cannot found main stream by roomid (%d) with error：%v", roomID, err)
	}

	// 开播
	if isOpen {
		if isOrigin { // 主推
			_, err := d.db.Exec(c, _notifyMainStreamOrigin, vendor, ms.Streaming|vendor, roomID, ms.Streaming)
			return err
		}
		// 转推
		_, err := d.db.Exec(c, _notifyMainStreamForward, ms.Streaming|vendor, roomID, ms.Streaming)
		return err
	} else {
		log.Infov(c, log.KV("----test----", fmt.Sprintf("---- %v ----- %v ---- %v ---- %v -", _notifyMainStreamOriginClose, newoptions, roomID, options)))
		// 关播的时候， 必须是当前的origin=传递过来的cdn才可以关， 修复开关播的时序性问题
		if isOrigin && ms.OriginUpstream == vendor {
			_, err := d.db.Exec(c, _notifyMainStreamOriginClose, newoptions, roomID, options)
			return err
		}

		// 转推
		_, err := d.db.Exec(c, _notifyMainStreamForward, ms.Streaming&^vendor, roomID, ms.Streaming)
		return err
	}
}
