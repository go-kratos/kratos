package dao

import (
	"context"
	"fmt"

	v1pb "go-common/app/service/live/xanchor/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_roomTable          = "room"
	_roomExtTablePrefix = "room_extend"

	_anchorTable = "anchor"
	_tagTable    = "tag"

	_areaTable = "ap_room_area_v2"

	// add room info
	_addRoomInfo1 = "insert into `%s` (`uid`) values (?)"
	// add room info
	_addRoomInfo2 = "insert into `%s` (`uid`,`room_id`) values (?,?)"

	// add room extend info
	_addRoomExtInfo = "insert into `%s_%d` (`room_id`) values (?)"

	// add anchor info
	_addAnchorInfo = "insert into `%s` (`uid`,`room_id`,`san_score`) values (?,?,12)"

	// update room info
	_updateRoomInfo = "update `%s` set %s where `room_id`=?"
	// update room extend info
	_updateRoomExtInfo = "update `%s_%d` set %s where `room_id`=?"
	// update anchor info
	_updateAnchorInfo = "update `%s` set %s where `uid`=?"
	// tag set info
	_tagSetInfo = "insert into `%s` set `target_type`=%d,%s on duplicate key update %s"
	// query room info
	_queryRoomInfo = "select `room`.`room_id`,`room`.`uid`,`room`.`title`,`room`.`description`,`room`.`tags`,`room`.`background`,`room`.`cover`,`room`.`lock_status`,`room`.`lock_time`,`room`.`hidden_time`,`room`.`record_switch`,`room`.`round_switch`,`room`.`live_start_time`,`room`.`live_screen_type`,`room`.`live_area_id`,`room`.`live_area_parent_id`,`room`.`live_type`,`anchor`.`san_score`,`anchor`.`profile_type`,`anchor`.`round_status`,`anchor`.`record_status`,`anchor`.`exp` from `room`,`anchor` where `room`.`uid`=`anchor`.`uid` and `%s` in (%s)"
	// query online room info
	_queryOnlineRoomInfo = "select `room`.`room_id`,`room`.`uid`,`room`.`title`,`room`.`description`,`room`.`tags`,`room`.`background`,`room`.`cover`,`room`.`lock_status`,`room`.`lock_time`,`room`.`hidden_time`,`room`.`record_switch`,`room`.`round_switch`,`room`.`live_start_time`,`room`.`live_screen_type`,`room`.`live_area_id`,`room`.`live_area_parent_id`,`room`.`live_type`,`anchor`.`san_score`,`anchor`.`profile_type`,`anchor`.`round_status`,`anchor`.`record_status`,`anchor`.`exp` from `room`,`anchor` where `room`.`uid`=`anchor`.`uid` and `live_start_time`!=0 order by id limit ?,?"

	// get parent area id
	_queryParentAreaID = "select `parent_id` from `%s` where `id`=?"
)

// fetchParentAreaID implementation
// fetchParentAreaID 查询房间信息
func (d *Dao) fetchParentAreaID(ctx context.Context, areaID int64) (parentAreaID int64, err error) {
	sql := fmt.Sprintf(_queryParentAreaID, _areaTable)
	err = d.dbLiveApp.QueryRow(ctx, sql, areaID).Scan(&parentAreaID)
	return
}

// fetchRoomByIDs implementation
// fetchRoomByIDs 查询房间信息
func (d *Dao) fetchRoomByIDs(ctx context.Context, req *v1pb.RoomByIDsReq) (resp *v1pb.RoomByIDsResp, err error) {
	condField := ""
	valueList := ""
	if len(req.RoomIds) > 0 {
		condField = "room`.`room_id"
		for i, id := range req.RoomIds {
			if i > 0 {
				valueList += ","
			}
			valueList += fmt.Sprintf("%d", id)
		}
	} else if len(req.Uids) > 0 {
		condField = "room`.`uid"
		for i, id := range req.Uids {
			if i > 0 {
				valueList += ","
			}
			valueList += fmt.Sprintf("%d", id)
		}
	}

	sql := fmt.Sprintf(_queryRoomInfo, condField, valueList)
	rows, err := d.db.Query(ctx, sql)
	if err != nil {
		log.Error("[dao.xanchor.mysql|fetchRoomByIDs] get room record error(%v), req(%v)", err, req)
		return nil, err
	}
	defer rows.Close()
	resp = &v1pb.RoomByIDsResp{}
	for rows.Next() {
		var data v1pb.RoomData
		err = rows.Scan(&data.RoomId, &data.Uid, &data.Title, &data.Description, &data.Tags, &data.Background, &data.Cover, &data.LockStatus, &data.LockTime, &data.HiddenTime, &data.AnchorRecordSwitch, &data.AnchorRoundSwitch, &data.LiveStartTime, &data.LiveScreenType, &data.AreaId, &data.ParentAreaId, &data.AnchorSan, &data.AnchorProfileType, &data.AnchorRoundStatus, &data.AnchorRecordStatus, &data.AnchorExp)
		// TODO area name
		if err != nil {
			log.Error("[dao.xanchor.mysql|fetchRoomByIDs] scan room record error(%v), req(%v)", err, req)
			return nil, err
		}
		resp.RoomDataSet[data.RoomId] = &data
		// TODO short-id
	}
	// TODO Tag List
	return
}

// roomOnlineList implementation
// roomOnlineList 在线房间列表
func (d *Dao) roomOnlineList(ctx context.Context, req *v1pb.RoomOnlineListReq) (resp *v1pb.RoomOnlineListResp, err error) {
	rows, err := d.db.Query(ctx, _queryOnlineRoomInfo, req.Page*req.PageSize, req.PageSize)
	if err != nil {
		log.Error("[dao.xanchor.mysql|roomOnlineList] get online room list error(%v), req(%v)", err, req)
		return nil, err
	}
	defer rows.Close()
	resp = &v1pb.RoomOnlineListResp{}
	for rows.Next() {
		var data v1pb.RoomData
		err = rows.Scan(&data.RoomId, &data.Uid, &data.Title, &data.Description, &data.Tags, &data.Background, &data.Cover, &data.LockStatus, &data.LockTime, &data.HiddenTime, &data.AnchorRecordSwitch, &data.AnchorRoundSwitch, &data.LiveStartTime, &data.LiveScreenType, &data.AreaId, &data.ParentAreaId, &data.AnchorSan, &data.AnchorProfileType, &data.AnchorRoundStatus, &data.AnchorRecordStatus, &data.AnchorExp)
		// TODO area name
		if err != nil {
			log.Error("[dao.xanchor.mysql|roomOnlineList] scan online room list error(%v), req(%v)", err, req)
			return nil, err
		}
		resp.RoomDataList[data.RoomId] = &data
		// TODO short-id
	}
	// TODO Tag List
	return
}

// roomCreate implementation
// roomCreate 房间创建
func (d *Dao) roomCreate(ctx context.Context, req *v1pb.RoomCreateReq) (resp *v1pb.RoomCreateResp, err error) {
	resp = &v1pb.RoomCreateResp{}
	tx, err := d.db.Begin(ctx)
	if err != nil {
		err = errors.WithStack(err)
		return resp, err
	}
	if req.RoomId != 0 {
		resp.RoomId = req.RoomId
		sql := fmt.Sprintf(_addRoomInfo2, _roomTable)
		_, err = tx.Exec(sql, req.Uid, req.RoomId)
		if err != nil {
			if e := tx.Rollback(); e != nil {
				log.Error("[dao.xanchor.mysql|roomCreate] create room record rollback error(%v), req(%v)", e, req)
			}
			// unique key exists error
			log.Error("[dao.xanchor.mysql|roomCreate] create room record error(%v), req(%v)", err, req)
			return resp, err
		}
	} else {
		sql := fmt.Sprintf(_addRoomInfo1, _roomTable)
		res, err := tx.Exec(sql, req.Uid)
		if err != nil {
			if e := tx.Rollback(); e != nil {
				log.Error("[dao.xanchor.mysql|roomCreate] create room record rollback error(%v), req(%v)", e, req)
			}
			// unique key exists error
			log.Error("[dao.xanchor.mysql|roomCreate] create room record error(%v), req(%v)", err, req)
			return resp, err
		}
		if resp.RoomId, err = res.LastInsertId(); err != nil {
			err = errors.WithStack(err)
			log.Error("[dao.xanchor.mysql|AddGuard] get last insert id error(%v), req(%v)", err, req)
		}
	}
	sql := fmt.Sprintf(_addRoomExtInfo, _roomExtTablePrefix, resp.RoomId%10)
	if _, err = tx.Exec(sql, resp.RoomId); err != nil {
		if e := tx.Rollback(); e != nil {
			log.Error("[dao.xanchor.mysql|roomCreate] create room extend record rollback error(%v), req(%v)", e, req)
		}
		log.Error("[dao.xanchor.mysql|roomCreate] create room extend record error(%v), req(%v)", err, req)
		return resp, err
	}
	sql = fmt.Sprintf(_addAnchorInfo, _anchorTable)
	if _, err = tx.Exec(sql, req.Uid, resp.RoomId); err != nil {
		if e := tx.Rollback(); e != nil {
			log.Error("[dao.xanchor.mysql|roomCreate] create anchor record rollback error(%v), req(%v)", e, req)
		}
		log.Error("[dao.xanchor.mysql|roomCreate] create anchor record error(%v), req(%v)", err, req)
		return resp, err
	}
	if err = tx.Commit(); err != nil {
		log.Error("[dao.xanchor.mysql|roomCreate] commit error(%v), req(%v)", err, req)
		return resp, err
	}
	return
}

// roomUpdate implementation
// roomUpdate 房间信息更新
func (d *Dao) roomUpdate(ctx context.Context, req *v1pb.RoomUpdateReq) (resp *v1pb.UpdateResp, err error) {
	updateSub := ""
	args := make([]interface{}, len(req.Fields)+1)
	args[len(req.Fields)] = req.RoomId
	for i, f := range req.Fields {
		switch f {
		case "title":
			args[i] = req.Title
		case "cover":
			args[i] = req.Cover
		case "tags":
			args[i] = req.Tags
		case "background":
			args[i] = req.Background
		case "description":
			args[i] = req.Description
		case "live_start_time":
			args[i] = req.LiveStartTime
			if req.LiveStartTime > 0 {
				if i > 0 {
					updateSub += ","
				}
				updateSub += "`live_mark`=1"
				if i == 0 {
					updateSub += ","
				}
			}
		case "live_screen_type":
			args[i] = req.LiveScreenType
		case "lock_status":
			args[i] = req.LockStatus
		case "lock_time":
			args[i] = req.LockTime
		case "hidden_time":
			args[i] = req.HiddenTime
		case "area_id":
			f = "live_area_id"
			args[i] = req.AreaId
			var parentAreaID int64
			parentAreaID, err = d.fetchParentAreaID(ctx, req.AreaId)
			if err != nil {
				log.Error("[dao.xanchor.mysql|roomUpdate] fetch parent area ID error(%v), req(%v)", err, req)
				err = ecode.InvalidParam
				return
			}
			if i > 0 {
				updateSub += ","
			}
			updateSub += fmt.Sprintf("`live_area_parent_id`=%d", parentAreaID)
			if i == 0 {
				updateSub += ","
			}
		case "anchor_round_switch":
			f = "round_switch"
			args[i] = req.AnchorRoundSwitch
		case "anchor_record_switch":
			f = "record_switch"
			args[i] = req.AnchorRecordSwitch
		default:
			log.Error("[dao.xanchor.mysql|roomUpdate] unsupported field(%v), req(%s)", f, req)
			err = ecode.InvalidParam
			return
		}
		if i > 0 {
			updateSub += ","
		}
		updateSub += fmt.Sprintf("`%s`=?", f)
	}
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_updateRoomInfo, _roomTable, updateSub)
	res, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("[dao.xanchor.mysql|roomUpdate] update room record error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

// roomBatchUpdate implementation
// roomBatchUpdate 房间信息批量更新
func (d *Dao) roomBatchUpdate(ctx context.Context, req *v1pb.RoomBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}
	for _, r := range req.Reqs {
		res, err := d.roomUpdate(ctx, r)
		if err != nil {
			log.Error("[dao.xanchor.mysql|roomBatchUpdate] update room record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// roomExtendUpdate implementation
// roomExtendUpdate 房间扩展信息更新
func (d *Dao) roomExtendUpdate(ctx context.Context, req *v1pb.RoomExtendUpdateReq) (resp *v1pb.UpdateResp, err error) {
	updateSub := ""
	args := make([]interface{}, len(req.Fields)+1)
	args[len(req.Fields)] = req.RoomId
	for i, f := range req.Fields {
		switch f {
		case "key_frame":
			args[i] = req.KeyFrame
		case "danmu_count":
			args[i] = req.DanmuCount
		case "popularity_count":
			args[i] = req.PopularityCount
		case "audience_count":
			args[i] = req.AudienceCount
		case "gift_count":
			args[i] = req.GiftCount
		case "gift_gold_amount":
			args[i] = req.GiftGoldAmount
		case "gift_gold_count":
			args[i] = req.GiftGoldCount
		default:
			log.Error("[dao.xanchor.mysql|roomExtendUpdate] unsupported field(%v), req(%s)", f, req)
			err = ecode.InvalidParam
			return
		}
		if i > 0 {
			updateSub += ","
		}
		updateSub += fmt.Sprintf("`%s`=?", f)
	}
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_updateRoomExtInfo, _roomExtTablePrefix, req.RoomId%10, updateSub)
	res, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("[dao.xanchor.mysql|roomExtendUpdate] update room extend record error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

// roomExtendBatchUpdate implementation
// roomExtendBatchUpdate 房间扩展信息批量更新
func (d *Dao) roomExtendBatchUpdate(ctx context.Context, req *v1pb.RoomExtendBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}
	for _, r := range req.Reqs {
		res, err := d.roomExtendUpdate(ctx, r)
		if err != nil {
			log.Error("[dao.xanchor.mysql|roomExtendBatchUpdate] update room extend record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// roomExtendIncre implementation
// roomExtendIncre 房间扩展信息增量更新
func (d *Dao) roomExtendIncre(ctx context.Context, req *v1pb.RoomExtendIncreReq) (resp *v1pb.UpdateResp, err error) {
	// TODO: req_id
	updateSub := ""
	args := make([]interface{}, len(req.Fields)+1)
	args[len(req.Fields)] = req.RoomId
	for i, f := range req.Fields {
		switch f {
		case "danmu_count":
			args[i] = req.DanmuCount
		case "popularity_count":
			args[i] = req.PopularityCount
		case "audience_count":
			args[i] = req.AudienceCount
		case "gift_count":
			args[i] = req.GiftCount
		case "gift_gold_amount":
			args[i] = req.GiftGoldAmount
		case "gift_gold_count":
			args[i] = req.GiftGoldCount
		default:
			log.Error("[dao.xanchor.mysql|roomExtendIncre] unsupported field(%v), req(%s)", f, req)
			err = ecode.InvalidParam
			return
		}
		if i > 0 {
			updateSub += ","
		}
		updateSub += fmt.Sprintf("`%s`=`%s`+(?)", f, f)
	}
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_updateRoomExtInfo, _roomExtTablePrefix, req.RoomId%10, updateSub)
	res, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("[dao.xanchor.mysql|roomExtendIncre] update room extend increment record error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

// roomExtendBatchIncre implementation
// roomExtendBatchIncre 房间扩展信息批量更新
func (d *Dao) roomExtendBatchIncre(ctx context.Context, req *v1pb.RoomExtendBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}
	for _, r := range req.Reqs {
		res, err := d.roomExtendIncre(ctx, r)
		if err != nil {
			log.Error("[dao.xanchor.mysql|roomExtendBatchIncre] update room extend increment record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// roomTagSet implementation
// roomTagSet 房间Tag更新
func (d *Dao) roomTagSet(ctx context.Context, req *v1pb.RoomTagSetReq) (resp *v1pb.UpdateResp, err error) {
	updateSub := ""
	args := make([]interface{}, len(req.Fields)*2+1)
	args[len(req.Fields)] = req.RoomId
	for i, f := range req.Fields {
		switch f {
		case "tag_value":
			args[i] = req.TagValue
			args[len(req.Fields)+i] = req.TagValue
		case "tag_attribute":
			args[i] = req.TagAttribute
			args[len(req.Fields)+i] = req.TagAttribute
		case "tag_expire_at":
			args[i] = req.TagExpireAt
			args[len(req.Fields)+i] = req.TagExpireAt
		default:
			log.Error("[dao.xanchor.mysql|roomTagSet] unsupported field(%v), req(%s)", f, req)
			err = ecode.InvalidParam
			return
		}
		if i > 0 {
			updateSub += ","
		}
		updateSub += fmt.Sprintf("`%s`=?", f)
	}
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_tagSetInfo, _tagTable, 2, updateSub, updateSub)
	res, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("[dao.xanchor.mysql|roomTagSet] set room tag error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

// anchorUpdate implementation
// anchorUpdate 主播信息更新
func (d *Dao) anchorUpdate(ctx context.Context, req *v1pb.AnchorUpdateReq) (resp *v1pb.UpdateResp, err error) {
	updateSub := ""
	args := make([]interface{}, len(req.Fields)+1)
	args[len(req.Fields)] = req.Uid
	for i, f := range req.Fields {
		switch f {
		case "profile_type":
			args[i] = req.ProfileType
		case "san_score":
			args[i] = req.SanScore
		case "round_status":
			args[i] = req.RoundStatus
		case "record_status":
			args[i] = req.RecordStatus
		case "exp":
			args[i] = req.Exp
		default:
			log.Error("[dao.xanchor.mysql|anchorUpdate] unsupported field(%v), req(%s)", f, req)
			err = ecode.InvalidParam
			return
		}
		if i > 0 {
			updateSub += ","
		}
		updateSub += fmt.Sprintf("`%s`=?", f)
	}
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_updateAnchorInfo, _anchorTable, updateSub)
	res, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("[dao.xanchor.mysql|anchorUpdate] update anchor record error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

// anchorBatchUpdate implementation
// anchorBatchUpdate 主播信息批量更新
func (d *Dao) anchorBatchUpdate(ctx context.Context, req *v1pb.AnchorBatchUpdateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}
	for _, r := range req.Reqs {
		res, err := d.anchorUpdate(ctx, r)
		if err != nil {
			log.Error("[dao.xanchor.mysql|anchorBatchUpdate] update anchor record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// anchorIncre implementation
// anchorIncre 主播信息增量更新
func (d *Dao) anchorIncre(ctx context.Context, req *v1pb.AnchorIncreReq) (resp *v1pb.UpdateResp, err error) {
	// TODO: req_id
	updateSub := ""
	args := make([]interface{}, len(req.Fields)+1)
	args[len(req.Fields)] = req.Uid
	for i, f := range req.Fields {
		switch f {
		case "san_score":
			args[i] = req.SanScore
		case "exp":
			args[i] = req.Exp
		default:
			log.Error("[dao.xanchor.mysql|anchorIncre] unsupported field(%v), req(%s)", f, req)
			err = ecode.InvalidParam
			return
		}
		if i > 0 {
			updateSub += ","
		}
		updateSub += fmt.Sprintf("`%s`=`%s`+(?)", f, f)
	}
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_updateAnchorInfo, _anchorTable, updateSub)
	res, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("[dao.xanchor.mysql|anchorUpdate] update anchor increment record error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

// anchorBatchIncre implementation
// anchorBatchIncre 主播信息批量增量更新
func (d *Dao) anchorBatchIncre(ctx context.Context, req *v1pb.AnchorBatchIncreReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}
	for _, r := range req.Reqs {
		res, err := d.anchorIncre(ctx, r)
		if err != nil {
			log.Error("[dao.xanchor.mysql|anchorBatchIncre] update anchor increment record error(%v), req(%v)", err, r)
			return nil, err
		}
		resp.AffectedRows += res.AffectedRows
	}
	return
}

// anchorTagSet implementation
// anchorTagSet 主播Tag更新
func (d *Dao) anchorTagSet(ctx context.Context, req *v1pb.AnchorTagSetReq) (resp *v1pb.UpdateResp, err error) {
	updateSub := ""
	args := make([]interface{}, len(req.Fields)*2+1)
	args[len(req.Fields)] = req.AnchorId
	for i, f := range req.Fields {
		switch f {
		case "tag_value":
			args[i] = req.TagValue
			args[len(req.Fields)+i] = req.TagValue
		case "tag_attribute":
			args[i] = req.TagAttribute
			args[len(req.Fields)+i] = req.TagAttribute
		case "tag_expire_at":
			args[i] = req.TagExpireAt
			args[len(req.Fields)+i] = req.TagExpireAt
		default:
			log.Error("[dao.xanchor.mysql|roomTagSet] unsupported field(%v), req(%s)", f, req)
			err = ecode.InvalidParam
			return
		}
		if i > 0 {
			updateSub += ","
		}
		updateSub += fmt.Sprintf("`%s`=?", f)
	}
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_tagSetInfo, _tagTable, 1, updateSub, updateSub)
	res, err := d.db.Exec(ctx, sql, args...)
	if err != nil {
		log.Error("[dao.xanchor.mysql|roomTagSet] set anchor tag error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}
