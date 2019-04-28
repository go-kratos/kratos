package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_insertTopicVideo = "insert ignore into topic_video (`topic_id`,`svid`,`state`) values %s"
	_updateScore      = "update topic_video set score=? where svid=%d"
	_updateState      = "update topic_video set state=? where svid=%d"
	_selectTopicVideo = "select svid from topic_video where topic_id=%d and state=0 %s order by score desc limit %d,%d"
	_selectVideoTopic = "select topic_id, score, state from topic_video where svid=%d"
)

// InsertTopicVideo 插入topic video表
func (d *Dao) InsertTopicVideo(ctx context.Context, svid int64, topicIDs []int64) (rowsAffected int64, err error) {
	var str string
	for _, topicID := range topicIDs {
		if len(str) != 0 {
			str += ","
		}
		str += fmt.Sprintf("(%d,%d,%d)", topicID, svid, api.TopicVideoStateUnAvailable)
	}

	insertSQL := fmt.Sprintf(_insertTopicVideo, str)
	res, err := d.db.Exec(ctx, insertSQL)
	if err != nil {
		log.Errorw(ctx, "log", "insert topic_video fail", "svid", svid, "topic_ids", topicIDs)
		return
	}

	rowsAffected, tmpErr := res.RowsAffected()
	if tmpErr != nil {
		log.Errorw(ctx, "log", "get rows affected fail", "svid", svid, "topic_ids", topicIDs)
	}
	log.V(1).Infow(ctx, "log", "insert one video topics", "svid", svid, "topics", topicIDs)
	return
}

// UpdateVideoScore 更新视频的score
//  @param topicID: 携带的时候会修改指定的topicID的video的score，否则会全部修改
func (d *Dao) UpdateVideoScore(ctx context.Context, svid int64, score float64) (err error) {
	updateSQL := fmt.Sprintf(_updateScore, svid)
	_, err = d.db.Exec(ctx, updateSQL, score)
	if err != nil {
		log.Errorw(ctx, "log", "update topic video score fail", "svid", svid, "score", score)
		return
	}

	return
}

// UpdateVideoState 更新视频的state
func (d *Dao) UpdateVideoState(ctx context.Context, svid int64, state int32) (err error) {
	updateSQL := fmt.Sprintf(_updateState, svid)
	_, err = d.db.Exec(ctx, updateSQL, state)
	if err != nil {
		log.Errorw(ctx, "log", "update topic video score fail", "svid", svid, "state", state)
		return
	}

	return
}

// ListTopicVideos 获取话题下排序的视频列表
// 按道理来说，Dao层不应该有那么多的复杂逻辑的，但是redis、db等操作在业务本身就是耦合在一起的，因此移到dao层，简化逻辑操作
// TODO: 这里把置顶的数据放在了redis里，所以导致排序问题过于复杂，待修正
func (d *Dao) ListTopicVideos(ctx context.Context, topicID int64, cursorPrev, cursorNext string, size int) (res []*api.VideoItem, hasMore bool, err error) {
	hasMore = true
	// 0. check
	if topicID == 0 {
		log.Errorw(ctx, "log", "topic_id=0")
		return
	}
	// 0.1 获取cursor和direction
	cursor, directionNext, err := parseCursor(ctx, cursorPrev, cursorNext)
	if err != nil {
		log.Warnw(ctx, "log", "parse cursor fail", "prev", cursorPrev, "next", cursorNext)
		return
	}

	// 1. 获取置顶视频
	stickSvid, err := d.GetStickTopicVideo(ctx, topicID)
	if err != nil {
		log.Warnw(ctx, "log", "get stick topic video fail")
		// 获取置顶视频失败后，属于可失败事件，继续往下走
	}
	stickMap := make(map[int64]bool)
	additionalConditionSQL := ""
	if len(stickSvid) > 0 {
		additionalConditionSQL = fmt.Sprintf("and svid not in (%s)", xstr.JoinInts(stickSvid))
		for _, svid := range stickSvid {
			stickMap[svid] = true
		}
	}

	// 2. 查询db
	var svids []int64
	dbOffset := cursor.Offset
	limit := size
	var rows *sql.Rows
	// 有两种情况才需要请求db：1、directionNext；2、directionPrev && stickRank==0
	if directionNext || cursor.StickRank == 0 {
		if !directionNext {
			dbOffset = cursor.Offset - 1 - size
			if dbOffset < 0 {
				dbOffset = 0
				limit = cursor.Offset - 1
			}
		}
		querySQL := fmt.Sprintf(_selectTopicVideo, topicID, additionalConditionSQL, dbOffset, limit)
		log.V(1).Infow(ctx, "log", "select topic video", "sql", querySQL)
		rows, err = d.db.Query(ctx, querySQL)
		if err != nil {
			log.Errorw(ctx, "log", "get topic video error", "err", err, "sql", querySQL)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var svid int64
			if err = rows.Scan(&svid); err != nil {
				log.Errorw(ctx, "log", "get topic from mysql fail", "sql", querySQL)
				return
			}
			svids = append(svids, svid)
		}
		log.V(1).Infow(ctx, "log", "get topic video", "svid", svids)
	}

	// 3. 组装回包
	if directionNext {
		if dbOffset == 0 {
			index := 0
			if cursor.StickRank != 0 {
				index = cursor.StickRank
			}
			for ; index < len(stickSvid); index++ {
				data, _ := json.Marshal(model.CursorValue{StickRank: index + 1})
				res = append(res, &api.VideoItem{Svid: stickSvid[index], CursorValue: string(data)})
			}
		}
		for index, svid := range svids {
			data, _ := json.Marshal(model.CursorValue{Offset: dbOffset + 1 + index})
			res = append(res, &api.VideoItem{Svid: svid, CursorValue: string(data)})
		}
		// TODO：为了避免db查询量过大，这里做限制
		if len(svids) != limit || dbOffset > model.MaxTopicVideoOffset {
			hasMore = false
		}
	} else {
		for index := len(svids) - 1; index >= 0; index-- {
			data, _ := json.Marshal(model.CursorValue{Offset: dbOffset + 1 + index})
			res = append(res, &api.VideoItem{Svid: svids[index], CursorValue: string(data)})
		}
		// 如果dbOffset==0，我们会把stick的视频页附上
		if dbOffset == 0 {
			index := len(stickSvid) - 1
			if cursor.StickRank != 0 {
				index = cursor.StickRank - 2
			}
			for ; index >= 0; index-- {
				data, _ := json.Marshal(model.CursorValue{StickRank: index + 1})
				res = append(res, &api.VideoItem{Svid: stickSvid[index], CursorValue: string(data)})
			}
			hasMore = false
		}
	}

	// 4. 添加hot_type结果
	for _, videoItem := range res {
		if _, exists := stickMap[videoItem.Svid]; exists {
			videoItem.HotType = api.TopicHotTypeStick
		}
	}

	return
}

// GetVideoTopic 获取视频的话题列表
func (d *Dao) GetVideoTopic(ctx context.Context, svid int64) (list []*api.TopicVideoItem, err error) {
	querySQL := fmt.Sprintf(_selectVideoTopic, svid)

	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorw(ctx, "log", "get video topic_id fail", "svid", svid)
		return
	}
	defer rows.Close()

	for rows.Next() {
		item := new(api.TopicVideoItem)
		item.Svid = svid
		if err = rows.Scan(&item.TopicId, &item.Score, &item.State); err != nil {
			log.Errorw(ctx, "log", "get topic video fail", "svid", svid)
			return
		}
		list = append(list, item)
	}
	log.V(1).Infow(ctx, "log", "get video topic", "sql", querySQL)
	return
}

// GetStickTopicVideo 获取置顶视频
func (d *Dao) GetStickTopicVideo(ctx context.Context, topicID int64) (list []int64, err error) {
	return d.getRedisList(ctx, fmt.Sprintf(model.ReidsStickTopicVideoKey, topicID))
}

// SetStickTopicVideo 设置置顶视频
func (d *Dao) SetStickTopicVideo(ctx context.Context, topicID int64, list []int64) (err error) {
	return d.setRedisList(ctx, fmt.Sprintf(model.ReidsStickTopicVideoKey, topicID), list)
}

// StickTopicVideo 操作置顶视频
func (d *Dao) StickTopicVideo(ctx context.Context, opTopicID, opSvid, op int64) (err error) {
	// 0. check
	topicVideoItems, err := d.GetVideoTopic(ctx, opSvid)
	if err != nil {
		log.Warnw(ctx, "log", "get svid topic topicVideoItems fail", "topic_id", opTopicID)
		return
	}
	var topicVideoItem *api.TopicVideoItem
	for _, item := range topicVideoItems {
		if item.TopicId == opTopicID {
			topicVideoItem = item
			break
		}
	}
	if topicVideoItem == nil {
		log.Errorw(ctx, "log", "stick topic fail due to error topic_id", "topic_id", opTopicID)
		err = ecode.TopicIDNotFound
		return
	}
	if topicVideoItem.State != api.TopicVideoStateAvailable {
		log.Errorw(ctx, "log", "topic video state unavailable to do sticking", "state", topicVideoItem.State, "topic_id", opTopicID)
		err = ecode.TopicVideoStateErr
		return
	}

	// 1. 获取stick topic video
	stickList, err := d.GetStickTopicVideo(ctx, opTopicID)
	if err != nil {
		log.Warnw(ctx, "log", "get stick topic video fail")
		return
	}

	// 2. 操作stick topic video
	var newStickList []int64
	if op != 0 {
		newStickList = append(newStickList, opSvid)
	}
	for _, stickSvid := range stickList {
		if stickSvid != opSvid {
			newStickList = append(newStickList, stickSvid)
		}
	}
	if len(newStickList) > model.MaxStickTopicVideoNum {
		newStickList = newStickList[:model.MaxStickTopicVideoNum]
	}

	// 3. 更新stick topic video
	err = d.SetStickTopicVideo(ctx, opTopicID, newStickList)
	if err != nil {
		log.Warnw(ctx, "update stick topic video fail")
		return
	}

	return
}
