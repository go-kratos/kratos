package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/live/dao-anchor/model"

	v1pb "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

//房間狀態常量
const (
	LIVE_OPEN  = 1
	LIVE_CLOSE = 0
	LIVE_ROUND = 2
)

//attr表相关常量定义
const (
	ATTRID_POPULARITY      = 1
	ATTRID_REVENUE         = 2
	ATTRID_DANMU           = 3
	ATTRID_RANK_LIST       = 4
	ATTRID_VALID_LIVE_DAYS = 5
)

const (
	//ATTRSUBID_RANK_HOUR 小时榜
	ATTRSUBID_RANK_HOUR = 1
)
const (
	//ATTRSUBID_POPULARITY_REALTIME 实时人气值
	ATTRSUBID_POPULARITY_REALTIME = 1
	//ATTRSUBID_POPULARITY_MAX_TO_ARG_7 7日人气峰值的均值
	ATTRSUBID_POPULARITY_MAX_TO_ARG_7 = 2
	//ATTRSUBID_POPULARITY_MAX_TO_ARG_30 30日人气峰值的均值
	ATTRSUBID_POPULARITY_MAX_TO_ARG_30 = 3
)
const (
	//ATTRSUBID_REVENUE_MINUTE_NUM_15  15分钟营收
	ATTRSUBID_REVENUE_MINUTE_NUM_15 = 1
	//ATTRSUBID_REVENUE_MINUTE_NUM_30  30分钟营收
	ATTRSUBID_REVENUE_MINUTE_NUM_30 = 2
	//ATTRSUBID_REVENUE_MINUTE_NUM_45  45分钟营收
	ATTRSUBID_REVENUE_MINUTE_NUM_45 = 3
	//ATTRSUBID_REVENUE_MINUTE_NUM_60  60分钟营收
	ATTRSUBID_REVENUE_MINUTE_NUM_60 = 4
)
const (
	//ATTRSUBID_DANMU_MINUTE_NUM_15  15分钟弹幕
	ATTRSUBID_DANMU_MINUTE_NUM_15 = 1
	//ATTRSUBID_DANMU_MINUTE_NUM_30  30分钟弹幕
	ATTRSUBID_DANMU_MINUTE_NUM_30 = 2
	//ATTRSUBID_DANMU_MINUTE_NUM_45  45分钟弹幕
	ATTRSUBID_DANMU_MINUTE_NUM_45 = 3
	//ATTRSUBID_DANMU_MINUTE_NUM_60  60分钟弹幕
	ATTRSUBID_DANMU_MINUTE_NUM_60 = 4
)

//有效开播天数
const (
	VALID_LIVE_DAYS_TYPE_1 = 1 //一次开播大于5分钟
	VALID_LIVE_DAYS_TYPE_2 = 2 // 累加大于等于120分钟
)
const (
	//ATTRSUBID_VALID_LIVE_DAYS_TYPE_1_DAY_7 近7天的有效天数,(有效天数：一次开播大于5分钟)
	ATTRSUBID_VALID_LIVE_DAYS_TYPE_1_DAY_7 = 1
	//ATTRSUBID_VALID_LIVE_DAYS_TYPE_1_DAY_14 近14天的有效天数,(有效天数：一次开播大于5分钟)
	ATTRSUBID_VALID_LIVE_DAYS_TYPE_1_DAY_14 = 2
	//ATTRSUBID_VALID_LIVE_DAYS_TYPE_2_DAY_7 近7天的有效天数,(有效天数：累加大于等于120分钟)
	ATTRSUBID_VALID_LIVE_DAYS_TYPE_2_DAY_7 = 3
	//ATTRSUBID_VALID_LIVE_DAYS_TYPE_2_DAY_30 近30天的有效天数,(有效天数：累加大于等于120分钟)
	ATTRSUBID_VALID_LIVE_DAYS_TYPE_2_DAY_30 = 4
)

//tag表相关常量定义
const (
	TAGID_PK          = 1
	TAGID_LOTTERY     = 2
	ROOM_EXT_SHARDING = 10

	EXP_2_SCORE_RATE = 100
)

const (
	_roomTable          = "room"
	_roomExtTablePrefix = "room_extend"
	_anchorTable        = "anchor"
	_tagTable           = "tag"
	_attrTable          = "attr"

	_shortTable   = "ap_short_room"
	_subAreaTable = "ap_room_area_v2"

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
	// tag create info
	_tagCreateInfo = "insert ignore into `%s` (`room_id`,`tag_id`,`tag_sub_id`,`tag_value`,`tag_ext`,`tag_expire_at`) values (?,?,?,?,?,?) on duplicate key update `tag_value`=?,`tag_ext`=?,`tag_expire_at`=?"
	// attr create info
	_attrCreateInfo = "insert ignore into `%s` (`room_id`,`attr_id`,`attr_sub_id`,`attr_value`,`attr_ext`) values (?,?,?,?,?) on duplicate key update `attr_value`=?,`attr_ext`=?"
	// attr set ex info
	_attrSetRoomId = "update `%s` set `room_id`=? where `attr_id`=? and `attr_sub_id`=? and `attr_value`=?"
	// attr set ex info
	_attrSetValue = "update `%s` set `attr_value`=? where `room_id`=? and `attr_id`=? and `attr_sub_id`=?"
	// attr select room_id
	_attrSelectRoomId = "select `room_id` from `%s` where `attr_id`=? and `attr_sub_id`=? and `attr_value`=?"
	//
	_attrSelectValue = "select `attr_value` from `%s` where `room_id`=? and `attr_id`=? and `attr_sub_id`=?"
	// query online room info
	_queryOnlineRoomInfo = "select `room_id`,`uid`,`title`,`description`,`tags`,`background`,`cover`,`lock_status`,`lock_time`,`hidden_time`,`record_switch`,`round_switch`,`live_start_time`,`live_screen_type`,`live_area_id`,`live_area_parent_id`,`live_type` from `room` where `live_start_time`!=0 order by `room_id` limit ?,?"
	// TODO 在播房间是否要处理轮播场景
	// query online room by area info
	_queryOnlineRoomByAreaInfo = "select `room_id` from `room` where `live_start_time`!=0 %s order by `live_start_time` desc"
	_queryOnlineRoomByAreaCond = "and (`live_area_id`=%d or `live_area_parent_id`=%d)"

	// query room info
	_queryRoomInfo = "select `room_id`,`uid`,`title`,`description`,`tags`,`background`,`cover`,`lock_status`,`lock_time`,`hidden_time`,`record_switch`,`round_switch`,`live_start_time`,`live_screen_type`,`live_area_id`,`live_area_parent_id`,`live_type` from `room` where `room_id` in (%s)"
	// query anchor info
	_queryAnchorInfo = "select `room_id`,`san_score`,`profile_type`,`round_status`,`record_status`,`exp` from `%s` where `uid` in (%s)"
	// query tag info
	_queryTagInfo = "select `room_id`,`tag_id`,`tag_sub_id`,`tag_value`,`tag_ext`,`tag_expire_at` from `%s` where `room_id` in (%s) and `tag_expire_at`>?"
	// query room ext info
	_queryRoomExtInfo = "select `room_id`,`keyframe`,`popularity_count` from `%s_%d` where `room_id` in (%s)"

	// get parent area id
	_queryParentAreaID = "select `parent_id` from `%s` where `id`=?"
	// get short id
	_queryShortID = "select `short_id`,`roomid` from `%s` where `roomid` in (%s)"

	// filter out short room-id and its corresponding room-id
	_filterShortID = "select `short_id`,`roomid` from `%s` where `short_id` in (%s) and `status`=1"

	// query attr info
	_queryAttrInfo = "select `room_id`,`attr_value` from `%s` where `attr_id`=? and `attr_sub_id`=? and `room_id` in (%s)"

	// query attr info
	_queryAttrInfo2 = "select `room_id`,`attr_sub_id`,`attr_value` from `%s` where `attr_id`=? order by `attr_sub_id`"

	// delete attr info
	_deleteAttrInfo = "delete from `%s` where `attr_id`=? and `attr_sub_id`=?"

	// query area info
	_queryAreaInfo = "select `id`,`name`,`parent_id` from `%s` where `id`=?"

	// query sub-areas for a given area
	_querySubAreaInfo = "select `id`,`name` from `%s` where `parent_id`=?"
)

// Turns short-id, if any, into room-id
func (d *Dao) dbNormalizeRoomIDs(ctx context.Context, roomIDs []int64) (resp []int64, err error) {
	if len(roomIDs) <= 0 {
		return
	}

	resp = make([]int64, len(roomIDs))

	// Resort to cache first, and for miss roomid, we save pair (roomid, index) to preserve order
	// of our final result
	agnosticRoomIds := make(map[int64]int)
	for i, roomID := range roomIDs {
		xid, ok := d.shortIDMapping.Get(fmt.Sprintf("%d", roomID))
		if ok {
			resp[i] = xid.(int64)
		} else {
			agnosticRoomIds[roomID] = i
		}
	}

	// We are done.
	if len(agnosticRoomIds) <= 0 {
		return
	}

	var queryValues string
	for roomID := range agnosticRoomIds {
		if len(queryValues) > 0 {
			queryValues += ","
		}
		queryValues += fmt.Sprintf("%d", roomID)
	}

	sql := fmt.Sprintf(_filterShortID, _shortTable, queryValues)
	rows, err := d.dbLiveApp.Query(ctx, sql)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbNormalizeRoomIDs] query short id record error(%v), sql(%s)", err, sql)
		return nil, err
	}

	defer rows.Close()

	shortIDMapping := make(map[int64]int64)
	for rows.Next() {
		var shortID int64
		var roomID int64
		err = rows.Scan(&shortID, &roomID)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbNormalizeRoomIDs] scan short id record error(%v), roomIDs(%v)",
				err, roomIDs)
			return nil, err
		}

		shortIDMapping[shortID] = roomID
	}

	for roomID, index := range agnosticRoomIds {
		xid, ok := shortIDMapping[roomID]
		if ok {
			resp[index] = xid
			d.shortIDMapping.Put(fmt.Sprintf("%d", roomID), xid)
		} else {
			resp[index] = roomID
			d.shortIDMapping.Put(fmt.Sprintf("%d", roomID), roomID)
		}
	}

	return
}

func (d *Dao) dbDealWithStatus(ctx context.Context, data *v1pb.RoomData) (err error) {
	if data == nil {
		return
	}

	// 处理开播状态
	if data.LiveStartTime > 0 {
		// 如果开播
		data.LiveStatus = LIVE_OPEN
	} else if data.AnchorRoundSwitch == 1 {
		// 如果轮播
		data.LiveStatus = LIVE_ROUND
	} else {
		data.LiveStatus = LIVE_CLOSE
	}

	// 处理隐藏状态
	if data.HiddenTime > time.Now().Unix() {
		data.HiddenStatus = 1
	}

	// 获取分区名称
	areaInfo := &model.AreaInfo{}
	if data.AreaId > 0 {
		areaInfo, err = d.dbFetchAreaInfo(ctx, data.AreaId)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbDealWithStatus] fetch area info error(%v), data(%v), areaid(%v)", err, data, data.AreaId)
			return
		}
		data.AreaName = areaInfo.AreaName
	}

	if areaInfo.ParentAreaID > 0 {
		areaInfo, err = d.dbFetchAreaInfo(ctx, areaInfo.ParentAreaID)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbDealWithStatus] fetch area info error(%v), data(%v), parent areaid(%v)", err, data, areaInfo.ParentAreaID)
			return
		}
		data.ParentAreaName = areaInfo.AreaName
	}
	return
}

// dbFetchRoomIDByUID implementation
// dbFetchRoomIDByUID 查询主播房间号
func (d *Dao) dbFetchRoomIDByUID(ctx context.Context, uid int64) (roomID int64) {
	uids := []int64{uid}
	res := make(map[int64]*v1pb.RoomData)

	err := d.dbFetchAnchorInfo(ctx, uids, res, false)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbFetchRoomIDByUID] get room ID error(%v), uid(%v)", err, uid)
		return
	}

	if len(res) <= 1 {
		return
	}

	for _, v := range res {
		if v.Uid == uid {
			roomID = v.RoomId
		}
	}

	return
}

// dbFetchAreaInfo implementation
// dbFetchAreaInfo 查询分区信息
func (d *Dao) dbFetchAreaInfo(ctx context.Context, areaID int64) (info *model.AreaInfo, err error) {
	if areaID <= 0 {
		return
	}
	key := fmt.Sprintf("%d", areaID)
	res, ok := d.areaInfoMapping.Get(key)
	if ok {
		return res.(*model.AreaInfo), nil
	}
	info = &model.AreaInfo{}
	sql := fmt.Sprintf(_queryAreaInfo, _subAreaTable)
	err = d.dbLiveApp.QueryRow(ctx, sql, areaID).Scan(&info.AreaID, &info.AreaName, &info.ParentAreaID)
	if err == nil {
		d.areaInfoMapping.Put(key, info)
	} else {
		// 尝试从配置中心拿到一级分区名称
		if confInfo, ok := d.c.FirstAreas[key]; ok {
			info.AreaID = areaID
			info.AreaName = confInfo.Name
			info.ParentAreaID = 0
			err = nil
		}
	}
	return
}

func (d *Dao) dbFetchAnchorInfo(ctx context.Context, uids []int64, resp map[int64]*v1pb.RoomData, overwrite bool) (err error) {
	if len(uids) <= 0 {
		return
	}

	valueList := ""
	for i, id := range uids {
		if i > 0 {
			valueList += ","
		}
		valueList += fmt.Sprintf("%d", id)
	}

	sql := fmt.Sprintf(_queryAnchorInfo, _anchorTable, valueList)
	rows, err := d.db.Query(ctx, sql)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] get anchor record error(%v), uids(%v)", err, uids)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		data := &v1pb.RoomData{
			AnchorLevel: new(v1pb.AnchorLevel),
		}
		var exp int64
		err = rows.Scan(&data.RoomId, &data.AnchorSan, &data.AnchorProfileType, &data.AnchorRoundStatus, &data.AnchorRecordStatus, &exp)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] scan anchor record error(%v), uids(%v)", err, uids)
			return err
		}

		data.AnchorLevel.Score = exp / EXP_2_SCORE_RATE

		data.AnchorLevel.Level, err = model.GetAnchorLevel(data.AnchorLevel.Score)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] Failed to get anchor level (%v), level(%d)", err, data.AnchorLevel.Score)
			return err
		}

		data.AnchorLevel.Left, data.AnchorLevel.Right, err = model.GetLevelScoreInfo(data.AnchorLevel.Level)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] Failed to get anchor level score (%v), level(%d)", err, data.AnchorLevel.Level)
			return err
		}

		data.AnchorLevel.Color, err = model.GetAnchorLevelColor(data.AnchorLevel.Level)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] Failed to get anchor level color (%v), level(%d)", err, data.AnchorLevel.Level)
			return err
		}

		data.AnchorLevel.MaxLevel = model.MaxAnchorLevel

		if overwrite {
			d := resp[data.RoomId]
			d.AnchorSan = data.AnchorSan
			d.AnchorProfileType = data.AnchorProfileType
			d.AnchorRoundStatus = data.AnchorRoundStatus
			d.AnchorRecordStatus = data.AnchorRecordStatus
			d.AnchorLevel = data.AnchorLevel
		} else {
			resp[data.RoomId] = data
		}
	}
	return
}

func (d *Dao) dbFetchRoomInfo(ctx context.Context, roomIDs []int64, resp map[int64]*v1pb.RoomData, overwrite bool) (err error) {
	if len(roomIDs) <= 0 {
		return
	}

	valueList := ""
	for i, id := range roomIDs {
		if i > 0 {
			valueList += ","
		}
		valueList += fmt.Sprintf("%d", id)
	}

	sql := fmt.Sprintf(_queryRoomInfo, valueList)
	rows, err := d.db.Query(ctx, sql)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbFetchRoomInfo] get room record error(%v), roomIDs(%v)", err, roomIDs)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		data := &v1pb.RoomData{}
		err = rows.Scan(&data.RoomId, &data.Uid, &data.Title, &data.Description, &data.Tags, &data.Background, &data.Cover, &data.LockStatus, &data.LockTime, &data.HiddenTime, &data.AnchorRecordSwitch, &data.AnchorRoundSwitch, &data.LiveStartTime, &data.LiveScreenType, &data.AreaId, &data.ParentAreaId, &data.LiveType)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomInfo] scan room record error(%v), roomIDs(%v)", err, roomIDs)
			return err
		}
		if overwrite {
			d := resp[data.RoomId]
			d.Uid = data.Uid
			d.Title = data.Title
			d.Description = data.Description
			d.Tags = data.Tags
			d.Background = data.Background
			d.Cover = data.Cover
			d.LockStatus = data.LockStatus
			d.LockTime = data.LockTime
			d.HiddenTime = data.HiddenTime
			d.AnchorRecordSwitch = data.AnchorRecordSwitch
			d.AnchorRoundSwitch = data.AnchorRoundSwitch
			d.LiveStartTime = data.LiveStartTime
			d.LiveScreenType = data.LiveScreenType
			d.AreaId = data.AreaId
			d.ParentAreaId = data.ParentAreaId
			d.LiveType = data.LiveType
		} else {
			resp[data.RoomId] = data
		}

		d.dbDealWithStatus(ctx, resp[data.RoomId])
	}
	return
}

func (d *Dao) dbFetchTagInfo(ctx context.Context, roomIDs []int64, resp map[int64]*v1pb.RoomData) (err error) {
	if len(roomIDs) <= 0 {
		return
	}

	valueList := ""
	for i, id := range roomIDs {
		if i > 0 {
			valueList += ","
		}
		valueList += fmt.Sprintf("%d", id)
	}

	sql := fmt.Sprintf(_queryTagInfo, _tagTable, valueList)
	rows, err := d.db.Query(ctx, sql, time.Now().Unix())
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbFetchTagInfo] get room record error(%v), roomIDs(%v)", err, roomIDs)
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var roomID int64
		data := &v1pb.TagData{}
		err = rows.Scan(&roomID, &data.TagId, &data.TagSubId, &data.TagValue, &data.TagExt, &data.TagExpireAt)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchTagInfo] scan room record error(%v), roomIDs(%v)", err, roomIDs)
			return err
		}

		if resp[roomID] == nil {
			resp[roomID] = &v1pb.RoomData{}
		}

		if resp[roomID].TagList == nil {
			resp[roomID].TagList = make([]*v1pb.TagData, 0)
		}

		resp[roomID].TagList = append(resp[roomID].TagList, data)
	}
	return
}

func (d *Dao) dbFetchExtInfo(ctx context.Context, roomIDs []int64, resp map[int64]*v1pb.RoomData) (err error) {
	if len(roomIDs) <= 0 {
		return
	}

	sharding := make(map[int64][]int64)
	for i := 0; i < 10; i++ {
		sharding[int64(i)] = make([]int64, 0, len(roomIDs))
	}
	for _, id := range roomIDs {
		sharding[id%ROOM_EXT_SHARDING] = append(sharding[id%ROOM_EXT_SHARDING], id)
	}

	for shard, ids := range sharding {
		if len(ids) <= 0 {
			continue
		}

		valueList := ""
		for i, id := range ids {
			if i > 0 {
				valueList += ","
			}
			valueList += fmt.Sprintf("%d", id)
		}

		sql := fmt.Sprintf(_queryRoomExtInfo, _roomExtTablePrefix, shard, valueList)
		rows, err := d.db.Query(ctx, sql)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchExtInfo] get room ext record error(%v), roomIDs(%v)", err, ids)
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var roomID int64
			var keyframe string
			var popularityCount int64
			err = rows.Scan(&roomID, &keyframe, &popularityCount)
			if err != nil {
				log.Error("[dao.dao-anchor.mysql|dbFetchExtInfo] scan room ext record error(%v), roomIDs(%v)", err, ids)
				return err
			}

			resp[roomID].Keyframe = keyframe
			resp[roomID].PopularityCount = popularityCount
		}
	}
	return
}

// dbFetchRoomByIDs implementation
// dbFetchRoomByIDs 查询房间信息
func (d *Dao) dbFetchRoomByIDs(ctx context.Context, req *v1pb.RoomByIDsReq) (resp *v1pb.RoomByIDsResp, err error) {
	resp = &v1pb.RoomByIDsResp{
		RoomDataSet: make(map[int64]*v1pb.RoomData),
	}

	roomIDs := make([]int64, 0, len(req.Uids))
	if len(req.RoomIds) > 0 {
		for _, id := range req.RoomIds {
			roomIDs = append(roomIDs, id)
		}

		// 先查room表
		err = d.dbFetchRoomInfo(ctx, roomIDs, resp.RoomDataSet, false)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] get room record error(%v), req(%v)", err, req)
			return nil, err
		}

		uids := make([]int64, 0, len(req.RoomIds))
		for _, v := range resp.RoomDataSet {
			uids = append(uids, v.Uid)
		}

		err = d.dbFetchAnchorInfo(ctx, uids, resp.RoomDataSet, true)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] get anchor record error(%v), req(%v)", err, req)
			return nil, err
		}
	} else if len(req.Uids) > 0 {
		// 先查anchor表
		uids := make([]int64, 0, len(req.Uids))
		for _, id := range req.Uids {
			uids = append(uids, id)
		}

		err = d.dbFetchAnchorInfo(ctx, uids, resp.RoomDataSet, false)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] get anchor record error(%v), req(%v)", err, req)
			return nil, err
		}

		for _, v := range resp.RoomDataSet {
			roomIDs = append(roomIDs, v.RoomId)
		}

		err = d.dbFetchRoomInfo(ctx, roomIDs, resp.RoomDataSet, true)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] get room record error(%v), req(%v)", err, req)
			return nil, err
		}
	}

	// Ext Info
	err = d.dbFetchExtInfo(ctx, roomIDs, resp.RoomDataSet)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] get room ext record error(%v), req(%v)", err, req)
		return nil, err
	}

	// Tag List
	err = d.dbFetchTagInfo(ctx, roomIDs, resp.RoomDataSet)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbFetchRoomByIDs] get tag record error(%v), req(%v)", err, req)
		return nil, err
	}

	// TODO： @wangyao需要处理short_id
	return
}

// dbOnlineListByArea implementation
// dbOnlineListByArea 分区在线房间列表
func (d *Dao) dbOnlineListByArea(ctx context.Context, areaId int64) (roomIDs []int64, err error) {
	cond := ""
	if areaId > 0 {
		cond = fmt.Sprintf(_queryOnlineRoomByAreaCond, areaId, areaId)
	}

	sql := fmt.Sprintf(_queryOnlineRoomByAreaInfo, cond)
	rows, err := d.db.Query(ctx, sql)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|dbOnlineListByArea] get online room list by area error(%v), areaId(%v)", err, areaId)
		return nil, err
	}

	defer rows.Close()

	roomIDs = make([]int64, 0)

	for rows.Next() {
		var roomid int64
		err = rows.Scan(&roomid)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|dbOnlineListByArea] scan online room list by area error(%v), areaId(%v)", err, areaId)
			return nil, err
		}
		roomIDs = append(roomIDs, roomid)
	}
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
				log.Error("[dao.dao-anchor.mysql|roomCreate] create room record rollback error(%v), req(%v)", e, req)
			}
			// unique key exists error
			log.Error("[dao.dao-anchor.mysql|roomCreate] create room record error(%v), req(%v)", err, req)
			return resp, err
		}
	} else {
		sql := fmt.Sprintf(_addRoomInfo1, _roomTable)
		res, err := tx.Exec(sql, req.Uid)
		if err != nil {
			if e := tx.Rollback(); e != nil {
				log.Error("[dao.dao-anchor.mysql|roomCreate] create room record rollback error(%v), req(%v)", e, req)
			}
			// unique key exists error
			log.Error("[dao.dao-anchor.mysql|roomCreate] create room record error(%v), req(%v)", err, req)
			return resp, err
		}
		if resp.RoomId, err = res.LastInsertId(); err != nil {
			err = errors.WithStack(err)
			log.Error("[dao.dao-anchor.mysql|AddGuard] get last insert id error(%v), req(%v)", err, req)
		}
	}

	sql := fmt.Sprintf(_addRoomExtInfo, _roomExtTablePrefix, resp.RoomId%ROOM_EXT_SHARDING)
	if _, err = tx.Exec(sql, resp.RoomId); err != nil {
		if e := tx.Rollback(); e != nil {
			log.Error("[dao.dao-anchor.mysql|roomCreate] create room extend record rollback error(%v), req(%v)", e, req)
		}
		log.Error("[dao.dao-anchor.mysql|roomCreate] create room extend record error(%v), req(%v)", err, req)
		return resp, err
	}

	sql = fmt.Sprintf(_addAnchorInfo, _anchorTable)
	if _, err = tx.Exec(sql, req.Uid, resp.RoomId); err != nil {
		if e := tx.Rollback(); e != nil {
			log.Error("[dao.dao-anchor.mysql|roomCreate] create anchor record rollback error(%v), req(%v)", e, req)
		}
		log.Error("[dao.dao-anchor.mysql|roomCreate] create anchor record error(%v), req(%v)", err, req)
		return resp, err
	}

	if err = tx.Commit(); err != nil {
		log.Error("[dao.dao-anchor.mysql|roomCreate] commit error(%v), req(%v)", err, req)
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
		case "live_type":
			args[i] = req.LiveType
		case "lock_status":
			args[i] = req.LockStatus
		case "lock_time":
			args[i] = req.LockTime
		case "hidden_time":
			args[i] = req.HiddenTime
		case "area_id":
			f = "live_area_id"
			args[i] = req.AreaId
			areaInfo := &model.AreaInfo{}
			// 审核后台会设置成无分区
			if req.AreaId > 0 {
				areaInfo, err = d.dbFetchAreaInfo(ctx, req.AreaId)
				if err != nil {
					log.Error("[dao.dao-anchor.mysql|roomUpdate] fetch area info error(%v), req(%v), areaid(%v)", err, req, req.AreaId)
					return
				}
			}
			if i > 0 {
				updateSub += ","
			}
			updateSub += fmt.Sprintf("`live_area_parent_id`=%d", areaInfo.ParentAreaID)
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
			log.Error("[dao.dao-anchor.mysql|roomUpdate] unsupported field(%v), req(%v)", f, req)
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
		log.Error("[dao.dao-anchor.mysql|roomUpdate] update room record error(%v), req(%v)", err, req)
		return
	}

	resp.AffectedRows, err = res.RowsAffected()
	return
}

// roomExtendUpdate implementation
// roomExtendUpdate 房间扩展信息更新
func (d *Dao) roomExtendUpdate(ctx context.Context, req *v1pb.RoomExtendUpdateReq) (resp *v1pb.UpdateResp, err error) {
	if len(req.Fields) <= 0 {
		log.Error("[dao.dao-anchor.mysql|roomExtendUpdate] no fields, req(%v)", req)
		err = ecode.InvalidParam
		return
	}

	updateSub := ""
	args := make([]interface{}, len(req.Fields)+1)
	args[len(req.Fields)] = req.RoomId

	for i, f := range req.Fields {
		switch f {
		case "keyframe":
			args[i] = req.Keyframe
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
			log.Error("[dao.dao-anchor.mysql|roomExtendUpdate] unsupported field(%v), req(%v)", f, req)
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
		log.Error("[dao.dao-anchor.mysql|roomExtendUpdate] update room extend record error(%v), req(%v)", err, req)
		return
	}

	resp.AffectedRows, err = res.RowsAffected()
	return
}

// roomExtendIncre implementation
// roomExtendIncre 房间扩展信息增量更新
func (d *Dao) roomExtendIncre(ctx context.Context, req *v1pb.RoomExtendIncreReq) (resp *v1pb.UpdateResp, err error) {
	if len(req.Fields) <= 0 {
		log.Error("[dao.dao-anchor.mysql|roomExtendIncre] no fields, req(%v)", req)
		err = ecode.InvalidParam
		return
	}

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
			log.Error("[dao.dao-anchor.mysql|roomExtendIncre] unsupported field(%v), req(%v)", f, req)
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
		log.Error("[dao.dao-anchor.mysql|roomExtendIncre] update room extend increment record error(%v), req(%v)", err, req)
		return
	}

	resp.AffectedRows, err = res.RowsAffected()
	return
}

// roomTagCreate implementation
// roomTagCreate 房间Tag创建
func (d *Dao) roomTagCreate(ctx context.Context, req *v1pb.RoomTagCreateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}

	sql := fmt.Sprintf(_tagCreateInfo, _tagTable)
	_, err = d.db.Exec(ctx, sql, req.RoomId, req.TagId, req.TagSubId, req.TagValue, req.TagExt, req.TagExpireAt, req.TagValue, req.TagExt, req.TagExpireAt)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|roomTagCreate] create room tag error(%v), req(%v)", err, req)
		return
	}

	resp.AffectedRows = 1
	return
}

// roomAttrCreate implementation
// roomAttrCreate 房间Attr创建
func (d *Dao) roomAttrCreate(ctx context.Context, req *v1pb.RoomAttrCreateReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}

	sql := fmt.Sprintf(_attrCreateInfo, _attrTable)
	_, err = d.db.Exec(ctx, sql, req.RoomId, req.AttrId, req.AttrSubId, req.AttrValue, req.AttrExt, req.AttrValue, req.AttrExt)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|roomAttrCreate] create room attr error(%v), req(%v)", err, req)
		return
	}

	resp.AffectedRows = 1
	return
}

// roomAttrSetEx implementation
// roomAttrSetEx 房间Attr更新/插入
func (d *Dao) roomAttrSetEx(ctx context.Context, req *v1pb.RoomAttrSetExReq) (resp *v1pb.UpdateResp, err error) {
	needSet, err := RoomAttrNeedSet(d, ctx, req)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|roomAttrSetEx] set room attr error(%v), req(%v)", err, req)
		return
	}
	if needSet <= 0 {
		return
	}
	if needSet == NEED_INSERT {
		reqCreate := &v1pb.RoomAttrCreateReq{}
		reqCreate.RoomId = req.RoomId
		reqCreate.AttrId = req.AttrId
		reqCreate.AttrSubId = req.AttrSubId
		reqCreate.AttrValue = req.AttrValue
		reqCreate.AttrExt = req.AttrExt
		resp, err = d.roomAttrCreate(ctx, reqCreate)
	} else {
		if req.AttrId == ATTRID_RANK_LIST {
			resp, err = d.roomAttrSetRoomId(ctx, req)
		} else {
			resp, err = d.roomAttrSetValue(ctx, req)
		}
	}
	return
}

const (
	NEED_INSERT = 1
	NEED_UPDATE = 2
)

// roomAttrNeedSet  内部函数 0 不需要set 1 需要insert 2 需要update
func RoomAttrNeedSet(d *Dao, ctx context.Context, req *v1pb.RoomAttrSetExReq) (resp int, err error) {
	attrId := req.AttrId
	//排行榜设计:更新roomID
	resp = 0
	if attrId == ATTRID_RANK_LIST {
		roomID := 0
		sql := fmt.Sprintf(_attrSelectRoomId, _attrTable)
		err = d.db.QueryRow(ctx, sql, req.AttrId, req.AttrSubId, req.AttrValue).Scan(&roomID)
		if err != nil && err.Error() == "sql: no rows in result set" {
			err = nil
			resp = NEED_INSERT
			return
		}
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|RoomAttrNeedSet] set room attr_rank error(%v), req(%v)", err, req)
			return
		}
		if int64(roomID) != req.RoomId {
			resp = NEED_UPDATE
		}
	} else {
		value := 0
		sql := fmt.Sprintf(_attrSelectValue, _attrTable)
		err = d.db.QueryRow(ctx, sql, req.AttrId, req.AttrSubId, req.AttrValue).Scan(value)
		if err != nil && err.Error() == "sql: no rows in result set" {
			err = nil
			resp = NEED_INSERT
			return
		}
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|RoomAttrNeedSet] set room attr error(%v), req(%v)", err, req)
			return
		}
		if int64(value) != req.AttrValue {
			resp = NEED_UPDATE
		}
	}

	return
}

//roomAttrSet  更新value
func (d *Dao) roomAttrSetValue(ctx context.Context, req *v1pb.RoomAttrSetExReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_attrSetValue, _attrTable)
	res, err := d.db.Exec(ctx, sql, req.AttrValue, req.RoomId, req.AttrId, req.AttrSubId)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|roomAttrSetEx] set room attr error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

//roomAttrSetRoomId 更新roomID
func (d *Dao) roomAttrSetRoomId(ctx context.Context, req *v1pb.RoomAttrSetExReq) (resp *v1pb.UpdateResp, err error) {
	resp = &v1pb.UpdateResp{}
	sql := fmt.Sprintf(_attrSetRoomId, _attrTable)
	res, err := d.db.Exec(ctx, sql, req.RoomId, req.AttrId, req.AttrSubId, req.AttrValue)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|roomAttrSetEx] set room attr error(%v), req(%v)", err, req)
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
			log.Error("[dao.dao-anchor.mysql|anchorUpdate] unsupported field(%v), req(%v)", f, req)
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
		log.Error("[dao.dao-anchor.mysql|anchorUpdate] update anchor record error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
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
			log.Error("[dao.dao-anchor.mysql|anchorIncre] unsupported field(%v), req(%v)", f, req)
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
		log.Error("[dao.dao-anchor.mysql|anchorUpdate] update anchor increment record error(%v), req(%v)", err, req)
		return
	}
	resp.AffectedRows, err = res.RowsAffected()
	return
}

// fetchAreas implementation
// fetchAreas 根据父分区号查询子分区
// If the request area-id does not exist, the function returns with the `err` being set.
func (d *Dao) fetchAreas(ctx context.Context, req *v1pb.FetchAreasReq) (resp *v1pb.FetchAreasResp, err error) {
	// Query parent area info first and fail fast in case the area doesn't exist.
	areaInfo, err := d.dbFetchAreaInfo(ctx, req.AreaId)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|fetchAreas] fetch main area info error(%v), req area(%d)", err, req.AreaId)
		return nil, err
	}

	resp = &v1pb.FetchAreasResp{
		Info: &v1pb.AreaInfo{
			AreaId:   req.AreaId,
			AreaName: areaInfo.AreaName,
		},
	}

	sql := fmt.Sprintf(_querySubAreaInfo, _subAreaTable)
	rows, err := d.dbLiveApp.Query(ctx, sql, req.AreaId)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|fetchAreas] fetch area records error(%v), req area(%d)", err, req.AreaId)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var subAreaID int64
		var subAreaName string
		if err = rows.Scan(&subAreaID, &subAreaName); err != nil {
			log.Error("[dao.dao-anchor.mysql|fetchAreas] fetch subarea info error(%v), req area(%d)", err, req.AreaId)
			return nil, err
		}
		resp.Areas = append(resp.Areas, &v1pb.AreaInfo{
			AreaId:   subAreaID,
			AreaName: subAreaName,
		})
	}

	return
}

// fetchAttrByIDs implementation
// fetchAttrByIDs 批量根据房间号查询指标
func (d *Dao) fetchAttrByIDs(ctx context.Context, req *v1pb.FetchAttrByIDsReq) (resp *v1pb.FetchAttrByIDsResp, err error) {
	if len(req.RoomIds) <= 0 {
		return
	}

	resp = &v1pb.FetchAttrByIDsResp{
		Attrs: make(map[int64]*v1pb.AttrData),
	}

	valueList := ""
	for i, id := range req.RoomIds {
		if i > 0 {
			valueList += ","
		}
		valueList += fmt.Sprintf("%d", id)
	}

	sql := fmt.Sprintf(_queryAttrInfo, _attrTable, valueList)
	rows, err := d.db.Query(ctx, sql, req.AttrId, req.AttrSubId)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|fetchAttrByIDs] get attr record error(%v), req(%v)", err, req)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		data := &v1pb.AttrData{
			AttrId:    req.AttrId,
			AttrSubId: req.AttrSubId,
		}
		err = rows.Scan(&data.RoomId, &data.AttrValue)
		if err != nil {
			log.Error("[dao.dao-anchor.mysql|fetchAttrByIDs] scan attr record error(%v), req(%v)", err, req)
			return nil, err
		}
		resp.Attrs[data.RoomId] = data
	}
	return
}

// deleteAttr implementation
// deleteAttr 删除一个指标
func (d *Dao) deleteAttr(ctx context.Context, req *v1pb.DeleteAttrReq) (resp *v1pb.UpdateResp, err error) {
	sql := fmt.Sprintf(_deleteAttrInfo, _attrTable)
	res, err := d.db.Exec(ctx, sql, req.AttrId, req.AttrSubId)
	if err != nil {
		log.Error("[dao.dao-anchor.mysql|roomExtendIncre] update room extend increment record error(%v), req(%v)", err, req)
		return
	}

	resp = &v1pb.UpdateResp{}
	resp.AffectedRows, err = res.RowsAffected()
	return
}
