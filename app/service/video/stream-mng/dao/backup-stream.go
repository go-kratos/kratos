package dao

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
	"math/rand"
	"strings"
	"time"

	"go-common/library/database/sql"

	"github.com/pkg/errors"
	"go-common/app/service/video/stream-mng/common"
)

const (
	_maxRetryTimes = 5

	_vendorBVC = 1
	_vendorKS  = 2
	_vendorQN  = 4
	_vendorTC  = 8
	_vendorWS  = 16

	_getBackupStreamByRoomID     = "SELECT `room_id`,`stream_name`, `key`,`default_vendor`, `origin_upstream`,`streaming`,`last_stream_time`, `expires_at`, `options` from `backup_stream` WHERE `room_id` = ? and status = 1"
	_getBackupStreamByStreamName = "SELECT `room_id`,`stream_name`, `key`, `default_vendor`, `origin_upstream`, `streaming`, `last_stream_time`, `expires_at`, `options` from `backup_stream` WHERE `stream_name` = ? and status = 1;"
	_getMultiBackupStreamByRID   = "SELECT `room_id`,`stream_name`, `key`, `default_vendor`, `origin_upstream`, `streaming`, `last_stream_time`, `expires_at`, `options` from `backup_stream` WHERE `room_id` = %d and status = 1"
	_getBackupRoom               = "SELECT distinct room_id from `backup_stream`"
	_insertBackupStream          = "INSERT INTO `backup_stream` (room_id, stream_name, `key`, default_vendor, expires_at, options) VALUES (?, ?, ?, ?, ?, ?);"
)

// GetBackupStreamByRoomID 根据roomid获取备用流信息
func (d *Dao) GetBackupStreamByRoomID(ctx context.Context, rid int64) (infos []*model.BackupStream, err error) {
	var rows *sql.Rows

	if rows, err = d.db.Query(ctx, _getBackupStreamByRoomID, rid); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		bs := new(model.BackupStream)
		if err = rows.Scan(&bs.RoomID, &bs.StreamName, &bs.Key, &bs.DefaultVendor, &bs.OriginUpstream, &bs.Streaming, &bs.LastStreamTime, &bs.ExpiresAt, &bs.Options); err != nil {
			err = errors.WithStack(err)
			return
		}
		infos = append(infos, bs)
	}

	err = rows.Err()
	return
}

// GetMultiBackupStreamByRID 批量查询备用流
func (d *Dao) GetMultiBackupStreamByRID(c context.Context, rids []int64) (infos []*model.BackupStream, err error) {
	len := len(rids)
	muSql := ""
	for i := 0; i < len; i++ {
		ss := fmt.Sprintf(_getMultiBackupStreamByRID, rids[i])
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
		bs := new(model.BackupStream)
		if err = rows.Scan(&bs.RoomID, &bs.StreamName, &bs.Key, &bs.DefaultVendor, &bs.OriginUpstream, &bs.Streaming, &bs.LastStreamTime, &bs.ExpiresAt, &bs.Options); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			err = errors.WithStack(err)
			return
		}
		infos = append(infos, bs)
	}

	err = rows.Err()
	return
}

// GetBackupRoom 临时获取所有的房间号
func (d *Dao) GetBackupRoom(ctx context.Context) (res map[int64]int64, err error) {
	res = map[int64]int64{}
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _getBackupRoom); err != nil {
		err = errors.WithStack(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		bs := new(model.BackupStream)
		if err = rows.Scan(&bs.RoomID); err != nil {
			err = errors.WithStack(err)
			return
		}
		res[bs.RoomID] = bs.RoomID
	}

	err = rows.Err()
	return
}

// CreateBackupStream 创建备用流
func (d *Dao) CreateBackupStream(ctx context.Context, bs *model.BackupStream) (*model.BackupStream, error) {
	if bs.StreamName == "" {
		bs.StreamName = fmt.Sprintf("live_%d_bs_%d", bs.RoomID, rand.Intn(9899999)+100000)
	}

	if bs.Key == "" {
		h := md5.New()
		h.Write([]byte(fmt.Sprintf("%s%d", bs.StreamName, time.Now().Nanosecond())))
		bs.Key = hex.EncodeToString(h.Sum(nil))
	}

	if bs.DefaultVendor == 0 {
		bs.DefaultVendor = 1
	}

	// 当传入的默认上行不是五家cdn
	if _, ok := common.BitwiseMapSrc[bs.DefaultVendor]; !ok {
		bs.DefaultVendor = 1
	}

	if bs.ExpiresAt.Before(time.Now()) {
		bs.ExpiresAt = time.Now().Add(time.Hour * 336) // 14 * 24
	}

	res, err := d.stmtBackupStreamCreate.Exec(ctx, bs.RoomID, bs.StreamName, bs.Key, bs.DefaultVendor, bs.ExpiresAt.Format("2006-01-02 15:04:05"), bs.Options)
	if err != nil {
		return bs, err
	}
	bs.ID, err = res.LastInsertId()

	return bs, err
}

// GetBackupStreamByStreamName 根据流名查询备用流
func (d *Dao) GetBackupStreamByStreamName(c context.Context, sn string) (*model.BackupStream, error) {
	row := d.db.QueryRow(c, _getBackupStreamByStreamName, sn)
	bs := &model.BackupStream{}
	err := row.Scan(&bs.RoomID, &bs.StreamName, &bs.Key,
		&bs.DefaultVendor, &bs.OriginUpstream, &bs.Streaming,
		&bs.LastStreamTime, &bs.ExpiresAt, &bs.Options)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

const (
	_setOriginUpstream  = "UPDATE `backup_stream` SET `origin_upstream` = ?, `streaming` = ? WHERE `stream_name` = ? and `origin_upstream` = 0;"
	_setForwardUpstream = "UPDATE `backup_stream` SET `streaming` = ? WHERE `stream_name` = ? and `streaming` = ?;"

	_setOriginUpstreamOnClose  = "UPDATE `backup_stream` SET `origin_upstream` = 0, `streaming` = 0, `last_stream_time` = CURRENT_TIMESTAMP WHERE `stream_name` = ? and `origin_upstream` != 0;"
	_setForwardUpstreamOnClose = "UPDATE `backup_stream` SET `streaming` = ? WHERE `stream_name` = ? and `streaming` = ?;"
)

var cdnBitwiseMap = map[string]int64{
	"bvc": _vendorBVC,
	"ks":  _vendorKS,
	"js":  _vendorKS, // alias
	"qn":  _vendorQN,
	"tc":  _vendorTC,
	"tx":  _vendorTC, // alias
	"txy": _vendorTC, // alias
	"ws":  _vendorWS,
}

// SetBackupStreamStreamingStatus
func (d *Dao) SetBackupStreamStreamingStatus(c context.Context, p *model.StreamingNotifyParam, bs *model.BackupStream, open bool) (*model.BackupStream, error) {
	bitwise, ok := cdnBitwiseMap[strings.ToLower(p.SRC)]
	if !ok {
		return nil, errors.New("unknown src:" + p.SRC)
	}
	for i := 1; i <= _maxRetryTimes; i++ {
		if open { // 开播
			if bs.Streaming&bitwise == bitwise {
				return bs, nil
			}

			if p.Type.String() == "0" { // 主推
				if bs.OriginUpstream == 0 { // 只有当前没有原始上行时才去尝试更新主推记录
					res, err := d.db.Exec(c, _setOriginUpstream, bitwise, bitwise, bs.StreamName)
					if err != nil {
						log.Errorw(c, "backup_stream_update_origin_record", err)
					} else {
						ra, err := res.RowsAffected()
						if err != nil {
							log.Errorw(c, "backup_stream_update_origin_record_rows_affected", err)
						}
						if ra == 1 { // 成功
							bs.Streaming = bitwise
							bs.OriginUpstream = bitwise
							return bs, nil
						}
						// 影响行数为 0，可能是发生了错误
						// 也可能是因为原始数据已经发生变更，等待后面重新读取DB中的数据。
					}
				} else { // 目前已经有上行了。在这里处理。
					if bitwise != bs.OriginUpstream {
						return bs, errors.New("origin upstream already exists")
					}
				}
			} else { // 转推
				if bs.OriginUpstream == 0 {
					return bs, errors.New("origin upstream not exists")
				}

				res, err := d.db.Exec(c, _setForwardUpstream, bs.Streaming|bitwise, bs.StreamName, bs.Streaming)
				if err != nil {
					log.Errorw(c, "backup_stream_update_forward_record", err)
				} else {
					ra, err := res.RowsAffected()
					if err != nil {
						log.Errorw(c, "backup_stream_update_forward_record_rows_affected", err)
					}
					if ra == 1 { // 成功
						bs.Streaming = bs.Streaming | bitwise
						return bs, nil
					}
				}
			}
		} else { // 关播
			if p.Type.String() == "0" { // 主推
				if bs.OriginUpstream != bitwise { // 如果不是当前主推，直接拒绝
					return bs, errors.New("permission denied")
				}

				res, err := d.db.Exec(c, _setOriginUpstreamOnClose, bs.StreamName)
				if err != nil {
					log.Errorw(c, "backup_stream_update_onclose_origin_record", err)
				} else {
					ra, err := res.RowsAffected()
					if err != nil {
						log.Errorw(c, "backup_stream_update_origin_onclose_record_rows_affected", err)
					}
					if ra == 1 { // 成功
						bs.Streaming = 0
						bs.OriginUpstream = 0
						return bs, nil
					}
				}
				// 影响行数为 0，可能是发生了错误
				// 也可能是因为原始数据已经发生变更，等待后面重新读取DB中的数据。
			} else { // 转推
				if bs.OriginUpstream == bitwise {
					return bs, errors.New("invalid params. you are origin upstream.")
				}

				res, err := d.db.Exec(c, _setForwardUpstreamOnClose, bs.Streaming&^bitwise, bs.StreamName, bs.Streaming)
				if err != nil {
					log.Errorw(c, "backup_stream_update_onclose_forward_record", err)
				} else {
					ra, err := res.RowsAffected()
					if err != nil {
						log.Errorw(c, "backup_stream_update_forward_onclose_record_rows_affected", err)
					}
					if ra == 1 { // 成功
						bs.Streaming = bs.Streaming &^ bitwise
						return bs, nil
					}
				}
			}
		}

		time.Sleep(time.Millisecond * 100)
		bs, err := d.GetBackupStreamByStreamName(c, bs.StreamName)
		if err != nil {
			log.Errorw(c, "backup_stream_refresh_record_row", err)
			return bs, errors.New("system busy")
		}
	}
	return bs, errors.New("update backup stream failed")
}
