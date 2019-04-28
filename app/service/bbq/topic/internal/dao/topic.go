package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup.v2"
	"go-common/library/xstr"
	"strings"
)

const (
	_selectTopic            = "select `id`, `name`, `desc`, `state` from topic where id in (%s)"
	_insertUpdateTopic      = "insert into topic (`name`,`score`,`state`,`video_num`) values %s on duplicate key update `video_num`=`video_num`+1"
	_selectTopicID          = "select id, name from topic where name in (%s)"
	_selectDiscoveryTopic   = "select id from topic where state=0 %s order by score desc, id desc limit %d, %d"
	_selectUnavailabelTopic = "select id from topic where state=1 limit %d,%d"
	_updateTopicField       = "update topic set `%s` = ? where `id` = ?"
)

const (
	_topicKey = "topic:%d"
)

// RawTopicInfo 从mysql获取topic info
func (d *Dao) RawTopicInfo(ctx context.Context, topicIDs []int64) (res map[int64]*api.TopicInfo, err error) {
	res = make(map[int64]*api.TopicInfo)
	if len(topicIDs) == 0 {
		return
	}

	querySQL := fmt.Sprintf(_selectTopic, xstr.JoinInts(topicIDs))
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorw(ctx, "log", "get topic error", "err", err, "sql", querySQL)
		return
	}
	defer rows.Close()
	for rows.Next() {
		topicInfo := new(api.TopicInfo)
		if err = rows.Scan(&topicInfo.TopicId, &topicInfo.Name, &topicInfo.Desc, &topicInfo.State); err != nil {
			log.Errorw(ctx, "log", "get topic from mysql fail", "sql", querySQL)
			return
		}
		topicInfo.CoverUrl = "http://i0.hdslb.com/bfs/bbq/video-image/userface/155886860_1547729941"
		res[topicInfo.TopicId] = topicInfo
	}
	log.V(1).Infow(ctx, "log", "get topic", "req", topicIDs, "rsp_size", len(res))
	return
}

// CacheTopicInfo 从缓存获取topic info
func (d *Dao) CacheTopicInfo(ctx context.Context, topicIDs []int64) (res map[int64]*api.TopicInfo, err error) {
	res = make(map[int64]*api.TopicInfo)

	keys := make([]string, 0, len(topicIDs))
	keyMidMap := make(map[int64]bool, len(topicIDs))
	for _, topicID := range topicIDs {
		key := fmt.Sprintf(_topicKey, topicID)
		if _, exist := keyMidMap[topicID]; !exist {
			// duplicate mid
			keyMidMap[topicID] = true
			keys = append(keys, key)
		}
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	for _, key := range keys {
		conn.Send("GET", key)
	}
	conn.Flush()
	var data []byte
	for i := 0; i < len(keys); i++ {
		if data, err = redis.Bytes(conn.Receive()); err != nil {
			if err == redis.ErrNil {
				err = nil
			} else {
				log.Errorv(ctx, log.KV("event", "redis_get"), log.KV("key", keys[i]))
			}
			continue
		}
		topicInfo := new(api.TopicInfo)
		json.Unmarshal(data, topicInfo)
		res[topicInfo.TopicId] = topicInfo
	}
	log.Infov(ctx, log.KV("event", "redis_get"), log.KV("row_num", len(res)))
	return
}

// AddCacheTopicInfo 添加topic info缓存
func (d *Dao) AddCacheTopicInfo(ctx context.Context, topicInfos map[int64]*api.TopicInfo) (err error) {
	keyValueMap := make(map[string][]byte, len(topicInfos))
	for topicID, topicInfo := range topicInfos {
		key := fmt.Sprintf(_topicKey, topicID)
		if _, exist := keyValueMap[key]; !exist {
			data, _ := json.Marshal(topicInfo)
			keyValueMap[key] = data
		}
	}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	for key, value := range keyValueMap {
		conn.Send("SET", key, value, "EX", d.topicExpire)
	}
	conn.Flush()
	for i := 0; i < len(keyValueMap); i++ {
		conn.Receive()
	}
	log.Infov(ctx, log.KV("event", "redis_set"), log.KV("row_num", len(topicInfos)))
	return
}

// DelCacheTopicInfo 删除topic info缓存
func (d *Dao) DelCacheTopicInfo(ctx context.Context, topicID int64) {
	var key = fmt.Sprintf(_topicKey, topicID)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	conn.Do("DEL", key)
}

// InsertTopics 插入话题
func (d *Dao) InsertTopics(ctx context.Context, topics map[string]*api.TopicInfo) (newTopics map[string]*api.TopicInfo, err error) {
	//func (d *Dao) InsertTopics(ctx context.Context, topics map[string]int64) (err error) {
	newTopics = make(map[string]*api.TopicInfo)
	// 0. check
	if len(topics) == 0 {
		return
	}
	if len(topics) > model.MaxBatchLen {
		err = ecode.TopicNumTooManyErr
		return
	}
	// 长度校验
	for _, item := range topics {
		if strings.Count(item.Name, "")-1 > model.MaxTopicNameLen {
			err = ecode.TopicNameLenErr
			log.Errorw(ctx, "log", "topic name len too long", "name", item.Name)
			return
		}
	}

	// 1. 插入更新
	group := errgroup.WithCancel(ctx)
	group.GOMAXPROCS(5)
	var groupInsertTopic = func(topicInfo *api.TopicInfo) {
		group.Go(func(ctx context.Context) (err error) {
			topicID, err := d.insertTopic(ctx, topicInfo)
			if err != nil {
				log.Warnw(ctx, "log", "get topic videos fail", "topic_name", topicInfo.Name)
				return
			}
			if topicID == 0 {
				log.Errorw(ctx, "log", "get error topic_id", "name", topicInfo.Name)
				err = ecode.TopicInsertErr
				return
			}
			topicInfo.TopicId = topicID
			return
		})
	}
	for _, topic := range topics {
		groupInsertTopic(topic)
	}
	err = group.Wait()
	if err != nil {
		log.Warnw(ctx, "log", "do group insert topic fail")
		return
	}
	// 由于insert的时候会返回ID，所以直接赋值返回
	newTopics = topics
	return
}

// insertTopic 插入话题
func (d *Dao) insertTopic(ctx context.Context, topicInfo *api.TopicInfo) (topicID int64, err error) {
	//func (d *Dao) InsertTopics(ctx context.Context, topics map[string]int64) (err error) {
	// 0. check
	// 长度校验
	if strings.Count(topicInfo.Name, "")-1 > model.MaxTopicNameLen {
		err = ecode.TopicNameLenErr
		log.Errorw(ctx, "log", "topic name len too long", "name", topicInfo.Name)
		return
	}

	var str string
	// 1. 插入更新
	str += fmt.Sprintf("('%s',%f,%d,1)", topicInfo.Name, topicInfo.Score, topicInfo.State)
	insertSQL := fmt.Sprintf(_insertUpdateTopic, str)
	log.V(1).Infow(ctx, "sql", insertSQL)
	res, err := d.db.Exec(ctx, insertSQL)
	if err != nil {
		log.Errorw(ctx, "log", "insert topic fail", "topic_name", topicInfo.Name)
		return
	}
	topicID, err = res.LastInsertId()
	if err != nil {
		log.Errorw(ctx, "log", "insert topic fail", "topic_name", topicInfo.Name)
		return
	}
	return
}

// UpdateTopic 更新话题，当前有简介和状态
// 这个函数把操作权其实已经交给上层了，设计上不是个好设计，但是在于避免重复代码
func (d *Dao) UpdateTopic(ctx context.Context, topicID int64, field string, value interface{}) (err error) {
	if field != "desc" && field != "state" {
		return ecode.ReqParamErr
	}
	querySQL := fmt.Sprintf(_updateTopicField, field)
	_, err = d.db.Exec(ctx, querySQL, value, topicID)
	if err != nil {
		log.Errorw(ctx, "log", "update topic field fail", "field", field, "value", value, "topic_id", topicID)
		return
	}
	d.DelCacheTopicInfo(ctx, topicID)
	return
}

// TopicID 通过话题name获取话题ID
// 话题ID结果存在topics中
func (d *Dao) TopicID(ctx context.Context, names []string) (topics map[string]int64, err error) {
	topics = make(map[string]int64)
	if len(names) == 0 {
		return
	}
	if len(names) > model.MaxBatchLen {
		err = ecode.TopicNumTooManyErr
		return
	}

	querySQL := fmt.Sprintf(_selectTopicID, "\""+strings.Join(names, "\",\"")+"\"")
	log.V(1).Infow(ctx, "log", "select topic id", "sql", querySQL)
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorw(ctx, "log", "get topic id error", "err", err, "sql", querySQL)
		return
	}
	defer rows.Close()
	var topicID int64
	var name string
	for rows.Next() {
		if err = rows.Scan(&topicID, &name); err != nil {
			log.Errorw(ctx, "log", "scan topic id error", "err", err, "sql", querySQL)
			return
		}
		topics[name] = topicID
	}
	log.V(1).Infow(ctx, "log", "get topic id", "req", names, "rsp", topics)
	return
}

// ListUnAvailableTopics .
func (d *Dao) ListUnAvailableTopics(ctx context.Context, page int32, size int32) (list []int64, hasMore bool, err error) {
	hasMore = true
	// 0. check
	if page < 1 {
		err = ecode.TopicReqParamErr
		return
	}
	if page > model.MaxDiscoveryTopicPage {
		hasMore = false
		return
	}

	// 2. get list
	offset := (page - 1) * size
	querySQL := fmt.Sprintf(_selectUnavailabelTopic, offset, size)
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorw(ctx, "log", "get topic video error", "err", err, "sql", querySQL)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var topicID int64
		if err = rows.Scan(&topicID); err != nil {
			log.Errorw(ctx, "log", "get topic from mysql fail", "sql", querySQL)
			return
		}
		list = append(list, topicID)
	}

	// 3. 判断has_more
	if len(list) < int(size) {
		hasMore = false
	}

	return
}

// ListRankTopics 获取推荐的话题列表
// TODO: 把置顶逻辑移上去
func (d *Dao) ListRankTopics(ctx context.Context, page int32, size int32) (list []int64, hasMore bool, err error) {
	hasMore = true
	// 0. check
	if page < 1 {
		err = ecode.TopicReqParamErr
		return
	}
	if page > model.MaxDiscoveryTopicPage {
		hasMore = false
		return
	}

	// 1. 获取置顶数据s
	additionalConditionSQL := ""
	stickList, err := d.GetStickTopic(ctx)
	if err != nil {
		log.Warnw(ctx, "log", "get stick topic fail")
	} else if len(stickList) > 0 {
		additionalConditionSQL = fmt.Sprintf("and id not in (%s)", xstr.JoinInts(stickList))
	}
	// 2. 若page=1，则获取推荐
	if page == 1 {
		list = stickList
	}

	// 3. 根据page获取话题列表
	offset := (page - 1) * size
	querySQL := fmt.Sprintf(_selectDiscoveryTopic, additionalConditionSQL, offset, size)
	log.Infow(ctx, "sql", querySQL, "page", page, "size", size)
	rows, err := d.db.Query(ctx, querySQL)
	if err != nil {
		log.Errorw(ctx, "log", "get topic video error", "err", err, "sql", querySQL)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var topicID int64
		if err = rows.Scan(&topicID); err != nil {
			log.Errorw(ctx, "log", "get topic from mysql fail", "sql", querySQL)
			return
		}
		list = append(list, topicID)
	}

	// 4. 判断has_more
	if len(list) < int(size) {
		hasMore = false
	}

	return
}

// GetStickTopic 获取置顶视频
// TODO: 这个方式是临时之计，当qps增大时，会导致热点的产生
func (d *Dao) GetStickTopic(ctx context.Context) (list []int64, err error) {
	return d.getRedisList(ctx, model.RedisStickTopicKey)
}

func (d *Dao) setStickTopic(ctx context.Context, list []int64) (err error) {
	return d.setRedisList(ctx, model.RedisStickTopicKey, list)
}

// StickTopic .
func (d *Dao) StickTopic(ctx context.Context, opTopicID, op int64) (err error) {
	// 0. check
	info, err := d.TopicInfo(ctx, []int64{opTopicID})
	if err != nil {
		log.Warnw(ctx, "log", "get topic info fail", "topic_id", opTopicID)
		return
	}
	topicInfo, exists := info[opTopicID]
	if !exists {
		log.Errorw(ctx, "log", "stick topic fail due to error topic_id", "topic_id", opTopicID)
		err = ecode.TopicIDNotFound
		return
	}
	if topicInfo.State != api.TopicStateAvailable {
		log.Errorw(ctx, "log", "topic state unavailable to do sticking", "state", topicInfo.State, "topic_id", opTopicID)
		err = ecode.TopicStateErr
		return
	}

	// 1. 获取stick topic
	stickList, err := d.GetStickTopic(ctx)
	if err != nil {
		log.Warnw(ctx, "log", "get stick topic fail")
		return
	}

	// 2. 操作stick topic
	var newStickList []int64
	if op != 0 {
		newStickList = append(newStickList, opTopicID)
	}
	for _, stickTopicID := range stickList {
		if stickTopicID != opTopicID {
			newStickList = append(newStickList, stickTopicID)
		}
	}
	if len(newStickList) > model.MaxStickTopicNum {
		newStickList = newStickList[:model.MaxStickTopicNum]
	}

	// 3. 更新stick topic
	err = d.setStickTopic(ctx, newStickList)
	if err != nil {
		log.Warnw(ctx, "update stick topic fail")
		return
	}

	return
}
