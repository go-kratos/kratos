package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1pb "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
)

//实时消费缓存设计，异步落地

const VALUE = "value"
const DATE = "date" //是否最新消息，是为1(需要刷新到DB) 否为0(不需要刷新到DB)
const DATE_1 = "1"

type ListIntValueInfo struct {
	Value int64 `json:"value"`
	Time  int64 `json:"time"`
}

//Set 设置实时数据
func (d *Dao) Set(ctx context.Context, redisKey string, value string, timeOut int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	if _, err = conn.Do("HMSET", redisKey, VALUE, value, DATE, DATE_1); err != nil {
		log.Error("redis_set_err:key=%s;value=%s;err=%v", redisKey, value, err)
		return
	}
	conn.Do("EXPIRE", redisKey, timeOut)
	return
}

//Incr  设置增加数据
func (d *Dao) Incr(ctx context.Context, redisKey string, value int64, timeOut int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	if err = conn.Send("HINCRBY", redisKey, VALUE, value); err != nil {
		log.Error("redis_incr_err:key=%s;value=%d;err=%v", redisKey, value, err)
		return
	}
	if err = conn.Send("HSET", redisKey, DATE, DATE_1); err != nil {
		log.Error("redis_hset_err:key=%s;date_value=%s;err=%v", redisKey, DATE_1, err)
		return
	}
	conn.Send("EXPIRE", redisKey, timeOut)

	if err = conn.Flush(); err != nil {
		log.Error("redisIncreRoomInfo conn.Flush error(%v)", err)
		return
	}

	_, err = conn.Receive()
	_, err = conn.Receive()
	_, err = conn.Receive()
	return
}
func (d *Dao) HGet(ctx context.Context, redisKey string) (resp int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	resp, err = redis.Int64(conn.Do("HGET", redisKey, VALUE))
	if err != nil {
		log.Error("redis_incr_err:key=%s;reply=%d;err=%v", redisKey, resp, err)
		return
	}
	return
}

func (d *Dao) SetList(ctx context.Context, redisKey string, value string, timeOut int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	if _, err = conn.Do("LPUSH", redisKey, value); err != nil {
		log.Error("redis_setList_error:key=%s;value=%s;err=%v", redisKey, value, err)
		return
	}
	conn.Do("EXPIRE", redisKey, timeOut)
	return
}

func (d *Dao) GetList(ctx context.Context, redisKey string, start int, end int) (resp []*ListIntValueInfo, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	res, err := redis.Values(conn.Do("LRANGE", redisKey, start, end))
	if err != nil {
		log.Error("redis_getList_error:key=%s;err=%v", redisKey, err)
		return
	}
	for _, sList := range res {
		list := &ListIntValueInfo{}
		if err = json.Unmarshal(sList.([]byte), &list); err != nil {
			log.Error("GetList_json_error")
			continue
		}
		resp = append(resp, list)
	}

	return
}

// GetRoomRecordsCurrent return a list of records corresponding to `content`.
func (d *Dao) GetRoomRecordsCurrent(ctx context.Context, content string, roomIds []int64) (list []int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	for _, roomId := range roomIds {
		key := d.SRoomRecordCurrent(content, roomId)
		if err = conn.Send("HGET", key.RedisKey, VALUE); err != nil {
			log.Error("GetRoomRecordsCurrent conn.Send(HGET, %s, %s) error(%v)", key.RedisKey, VALUE, err)
			return nil, err
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("GetRoomRecordsCurrent conn.Flush error(%v)", err)
		return nil, err
	}

	for i := 0; i < len(roomIds); i++ {
		var data int64
		if data, err = redis.Int64(conn.Receive()); err != nil {
			if err != redis.ErrNil {
				log.Error("GetRoomRecordsCurrent conn.Receive() %d error(%v)", i, err)
				return nil, err
			}
		}

		list = append(list, data)
	}

	return
}
func (d *Dao) DelRoomRecordsCurrent(ctx context.Context, content string, roomIds []int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	args := make([]interface{}, len(roomIds))
	for i, roomId := range roomIds {
		args[i] = d.SRoomRecordCurrent(content, roomId)
	}
	if _, err = conn.Do("DEL", args...); err != nil {
		log.Error("DelRoomRecordsCurrent_del_error:%v;roomIds=%v", err, roomIds)
		return
	}
	return
}

func (d *Dao) SetRoomRecordsList(ctx context.Context, roomIds []int64, keys map[int64]interface{}, values map[int64]string) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	for _, roomId := range roomIds {
		keyInfo := keys[roomId].(*redisKeyResp)
		value := values[roomId]
		if err = conn.Send("LPUSH", keyInfo.RedisKey, value); err != nil {
			log.Error("SetRoomRecordsList conn.Send(LPUSH, %s, %s) error(%v)", keyInfo.RedisKey, value, err)
			return err
		}

		if err = conn.Send("EXPIRE", keyInfo.RedisKey, keyInfo.TimeOut); err != nil {
			log.Error("SetRoomRecordsList conn.Send(EXPIRE, %s, %d) error(%v)", keyInfo.RedisKey, keyInfo.TimeOut, err)
			return err
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("SetRoomRecordsList conn.Flush() error(%v)", err)
		return err
	}

	for range roomIds {
		conn.Receive()
		conn.Receive()
	}

	return
}

// GetRoomLiveRecordsRange can partially succeed, and in this case, err is still nil.
func (d *Dao) GetRoomLiveRecordsRange(ctx context.Context, content string, roomIds []int64, liveTime int64, start, end int) (resp map[int64][]*ListIntValueInfo, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	okRoomIds := make([]int64, 0, len(roomIds))
	for _, roomId := range roomIds {
		keyInfo := d.LRoomLiveRecordList(content, roomId, liveTime)
		if err = conn.Send("LRANGE", keyInfo.RedisKey, start, end); err != nil {
			log.Error("GetRoomLiveRecordsRange conn.Send(LRANGE, %s, %d, %d) error(%v)", keyInfo.RedisKey, start, end, err)
			continue
		}

		okRoomIds = append(okRoomIds, roomId)
	}

	if err = conn.Flush(); err != nil {
		log.Error("GetRoomLiveRecordsRange conn.Flush() error(%v)", err)
		return nil, err
	}

	resp = make(map[int64][]*ListIntValueInfo)

	for i := 0; i < len(okRoomIds); i++ {
		values, err := redis.Values(conn.Receive())
		if err != nil {
			log.Error("GetRoomLiveRecordsRange redis.Values(conn.Receive()) error(%v)", err)
			continue
		}

		roomId := okRoomIds[i]

		for _, info := range values {
			valueInfo := &ListIntValueInfo{}
			if err = json.Unmarshal(info.([]byte), &valueInfo); err != nil {
				log.Error("GetRoomLiveRecordsRange json unmarshall error(%v)", err)
				continue
			}
			resp[roomId] = append(resp[roomId], valueInfo)
		}
	}
	return resp, nil
}

const (
	_roomInfoKey   = "room_info_v3:%d"
	_anchorInfoKey = "anchor_info:%d"
	_onlineListKey = "online_list_v3:%d"
	_tagListKey    = "tag_list_v3:%d"

	_onlineListAllArea = 0
)

var (
	_allRoomInfoFields = []string{"uid", "title", "cover", "tags", "background", "description", "live_start_time", "live_status", "live_screen_type", "live_type", "lock_status", "lock_time", "hidden_time", "hidden_status", "area_id", "parent_area_id", "anchor_profile_type", "anchor_round_switch", "anchor_record_switch", "anchor_exp", "popularity_count", "keyframe"}
)

func (d *Dao) filterOutTagList(fields []string) (resp []string, needTagList bool) {
	resp = make([]string, 0, len(fields))
	for _, f := range fields {
		if f == "tag_list" {
			needTagList = true
		} else {
			resp = append(resp, f)
		}
	}
	return
}

func (d *Dao) redisGetRoomInfo(ctx context.Context, roomID int64, fields []string) (data *v1pb.RoomData, err error) {
	if len(fields) <= 0 {
		return
	}

	conn := d.redis.Get(ctx)
	defer conn.Close()

	roomKey := fmt.Sprintf(_roomInfoKey, roomID)

	ok, err := redis.Bool(conn.Do("EXPIRE", roomKey, d.c.Common.ExpireTime))
	if err != nil && err != redis.ErrNil {
		log.Error("redisGetRoomInfo conn.Do(EXPIRE, %s) error(%v)", roomKey, err)
		return
	}

	if !ok {
		err = ecode.RoomNotFound
		return
	}

	for _, f := range fields {
		if f == "room_id" {
			continue
		}
		if err = conn.Send("HGET", roomKey, f); err != nil {
			log.Error("redisGetRoomInfo conn.Send(HGET, %s, %s) error(%v)", roomKey, f, err)
			return
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("redisGetRoomInfo conn.Flush error(%v)", err)
		return
	}

	data = &v1pb.RoomData{
		RoomId:      roomID,
		AnchorLevel: new(v1pb.AnchorLevel),
	}
	for _, f := range fields {
		if f == "room_id" {
			continue
		}
		reply, err := conn.Receive()
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") != -1 {
			return data, err
		}
		switch f {
		case "uid":
			data.Uid, err = redis.Int64(reply, err)
		case "title":
			data.Title, err = redis.String(reply, err)
		case "cover":
			data.Cover, err = redis.String(reply, err)
		case "tags":
			data.Tags, err = redis.String(reply, err)
		case "background":
			data.Background, err = redis.String(reply, err)
		case "description":
			data.Description, err = redis.String(reply, err)
		case "live_start_time":
			data.LiveStartTime, err = redis.Int64(reply, err)
		case "live_status":
			data.LiveStatus, err = redis.Int64(reply, err)
		case "live_screen_type":
			data.LiveScreenType, err = redis.Int64(reply, err)
		case "live_type":
			data.LiveType, err = redis.Int64(reply, err)
		case "lock_status":
			data.LockStatus, err = redis.Int64(reply, err)
		case "lock_time":
			data.LockTime, err = redis.Int64(reply, err)
		case "hidden_time":
			data.HiddenTime, err = redis.Int64(reply, err)
		case "hidden_status":
			data.HiddenStatus, err = redis.Int64(reply, err)
		case "area_id":
			data.AreaId, err = redis.Int64(reply, err)
		case "parent_area_id":
			data.ParentAreaId, err = redis.Int64(reply, err)
		case "anchor_san":
			data.AnchorSan, err = redis.Int64(reply, err)
		case "anchor_profile_type":
			data.AnchorProfileType, err = redis.Int64(reply, err)
		case "anchor_round_switch":
			data.AnchorRoundSwitch, err = redis.Int64(reply, err)
		case "anchor_record_switch":
			data.AnchorRecordSwitch, err = redis.Int64(reply, err)
		case "anchor_exp":
			data.AnchorLevel.Score, err = redis.Int64(reply, err)
		case "popularity_count":
			data.PopularityCount, err = redis.Int64(reply, err)
		case "keyframe":
			data.Keyframe, err = redis.String(reply, err)
		default:
			log.Error("redisGetRoomInfo unsupported field(%v), roomID(%d)", f, roomID)
			err = ecode.InvalidParam
			return nil, err
		}
		if err != nil {
			log.Warn("redisGetRoomInfo conn.Receive() field(%v), error(%v)", f, err)
			return nil, err
		}
	}

	return
}

func (d *Dao) redisSetRoomInfo(ctx context.Context, roomID int64, fields []string, data *v1pb.RoomData, fastfail bool) (err error) {
	if len(fields) <= 0 {
		return
	}

	conn := d.redis.Get(ctx)
	defer conn.Close()

	roomKey := fmt.Sprintf(_roomInfoKey, roomID)

	ok, err := redis.Bool(conn.Do("EXPIRE", roomKey, d.c.Common.ExpireTime))
	if err != nil && err != redis.ErrNil {
		log.Error("redisSetRoomInfo conn.Do(EXPIRE, %s) error(%v)", roomKey, err)
		return
	}

	if !ok && fastfail {
		return
	}

	args := make([]interface{}, len(fields)*2+1)
	args[0] = roomKey

	for i, f := range fields {
		var v interface{}
		switch f {
		case "roomid":
			v = data.RoomId
		case "uid":
			v = data.Uid
		case "title":
			v = data.Title
		case "cover":
			v = data.Cover
		case "tags":
			v = data.Tags
		case "background":
			v = data.Background
		case "description":
			v = data.Description
		case "live_start_time":
			v = data.LiveStartTime
		case "live_status":
			v = data.LiveStatus
		case "live_screen_type":
			v = data.LiveScreenType
		case "live_type":
			v = data.LiveType
		case "lock_status":
			v = data.LockStatus
		case "lock_time":
			v = data.LockTime
		case "hidden_time":
			v = data.HiddenTime
		case "hidden_status":
			v = data.HiddenStatus
		case "area_id":
			v = data.AreaId
		case "parent_area_id":
			v = data.ParentAreaId
		case "anchor_san":
			v = data.AnchorSan
		case "anchor_profile_type":
			v = data.AnchorProfileType
		case "anchor_round_switch":
			v = data.AnchorRoundSwitch
		case "anchor_record_switch":
			v = data.AnchorRecordSwitch
		case "anchor_exp":
			v = data.AnchorRecordSwitch
		case "popularity_count":
			v = data.PopularityCount
		case "keyframe":
			v = data.Keyframe
		default:
			log.Error("redisSetRoomInfo unsupported field(%v), roomID(%d)", f, roomID)
			return ecode.InvalidParam
		}

		args[i*2+1] = f
		args[i*2+2] = v
	}

	if _, err = conn.Do("HMSET", args...); err != nil {
		log.Error("redisSetRoomInfo conn.Do(HMSET, %v) error(%v)", args, err)
		return
	}

	return
}

func (d *Dao) redisIncreRoomInfo(ctx context.Context, roomID int64, fields []string, data *v1pb.RoomData) (err error) {
	if len(fields) <= 0 {
		return
	}

	conn := d.redis.Get(ctx)
	defer conn.Close()

	roomKey := fmt.Sprintf(_roomInfoKey, roomID)

	ok, err := redis.Bool(conn.Do("EXPIRE", roomKey, d.c.Common.ExpireTime))
	if err != nil && err != redis.ErrNil {
		log.Error("redisIncreRoomInfo conn.Do(EXPIRE, %s) error(%v)", roomKey, err)
		return
	}

	// fast fail if key not exists
	if !ok {
		return
	}

	for _, f := range fields {
		var v int64
		switch f {
		case "anchor_san":
			v = data.AnchorSan
		case "anchor_exp":
			v = data.AnchorRecordSwitch
		case "popularity_count":
			v = data.PopularityCount
		default:
			log.Error("redisIncreRoomInfo unsupported field(%v), roomID(%d)", f, roomID)
			return ecode.InvalidParam
		}
		if err = conn.Send("HINCRBY", roomKey, f, v); err != nil {
			log.Error("redisIncreRoomInfo conn.Send(HINCRBY, %s, %s, %d) error(%v)", roomKey, f, v, err)
			return
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("redisIncreRoomInfo conn.Flush error(%v)", err)
		return
	}

	for _, f := range fields {
		_, err = conn.Receive()
		if err != nil && err != redis.ErrNil {
			log.Error("redisIncreRoomInfo conn.Receive() field(%v), error(%v)", f, err)
			return
		}
	}

	return
}

func (d *Dao) redisGetOnlineList(ctx context.Context, areaID int64) (list []int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	key := fmt.Sprintf(_onlineListKey, areaID)

	ok, err := redis.Bool(conn.Do("EXPIRE", key, d.c.Common.ExpireTime))
	if err != nil && err != redis.ErrNil {
		log.Error("redisGetOnlineList conn.Do(EXPIRE, %s) error(%v)", key, err)
		return
	}

	if !ok {
		// 不存在或者在播列表为空都会重新去DB获取
		return
	}

	list = make([]int64, 0)

	roomids, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			log.Error("redisGetOnlineList conn.Do(SMEMBERS, %s) error(%v)", key, err)
			return
		}
		return list, nil
	}

	for _, id := range roomids {
		roomid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Warn("redisGetOnlineList ParseInt(%d) error(%v)", roomid, err)
			return nil, err
		}
		list = append(list, roomid)
	}

	return
}

func (d *Dao) redisSetOnlineList(ctx context.Context, areaID int64, list []int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	key := fmt.Sprintf(_onlineListKey, areaID)

	if len(list) <= 0 {
		// 设置哨兵
		if _, err = conn.Do("SETEX", key, d.c.Common.ExpireTime, "emptylist"); err != nil {
			log.Error("redisSetOnlineList conn.Do(SETEX, %s, %v) error(%v)", key, list, err)
		}
		return
	}

	if _, err = conn.Do("DEL", key); err != nil && err != redis.ErrNil {
		log.Error("redisSetOnlineList conn.Do(DEL, %s) error(%v)", key, err)
		return
	}

	args := make([]interface{}, len(list)+1)
	args[0] = key

	for i, id := range list {
		args[i+1] = id
	}

	if _, err = conn.Do("SADD", args...); err != nil {
		log.Error("redisSetOnlineList conn.Do(SADD, %s, %v) error(%v)", key, list, err)
		return
	}

	return
}

func (d *Dao) redisAddOnlineList(ctx context.Context, areaID int64, roomID int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	key := fmt.Sprintf(_onlineListKey, areaID)

	typ, err := redis.String(conn.Do("TYPE", key))
	if err != nil && err != redis.ErrNil {
		log.Error("redisAddOnlineList conn.Do(TYPE, %s) error(%v)", key, err)
		return
	}

	if strings.ToLower(typ) == "none" {
		// 不存在就不处理
		return
	} else if strings.ToLower(typ) != "set" {
		if _, err = conn.Do("DEL", key); err != nil && err != redis.ErrNil {
			log.Error("redisAddOnlineList conn.Do(DEL, %s) error(%v)", key, err)
			return
		}
	}

	if _, err = conn.Do("SADD", key, roomID); err != nil {
		log.Error("redisAddOnlineList conn.Do(SADD, %s, %d) error(%v)", key, roomID, err)
		return
	}

	return
}

func (d *Dao) redisDelOnlineList(ctx context.Context, areaID int64, roomID int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	key := fmt.Sprintf(_onlineListKey, areaID)

	typ, err := redis.String(conn.Do("TYPE", key))
	if err != nil && err != redis.ErrNil {
		log.Error("redisDelOnlineList conn.Do(TYPE, %s) error(%v)", key, err)
		return
	}

	if strings.ToLower(typ) != "set" {
		// 不存在就不处理
		return
	}

	if _, err = conn.Do("SREM", key, roomID); err != nil {
		log.Error("redisDelOnlineList conn.Do(SADD, %s, %d) error(%v)", key, roomID, err)
		return
	}

	return
}

func (d *Dao) redisGetTagList(ctx context.Context, roomID int64) (list []*v1pb.TagData, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	key := fmt.Sprintf(_tagListKey, roomID)

	ok, err := redis.Bool(conn.Do("EXPIRE", key, d.c.Common.ExpireTime))
	if err != nil && err != redis.ErrNil {
		log.Error("redisGetTagList conn.Do(EXPIRE, %s) error(%v)", key, err)
		return nil, err
	}

	if !ok {
		err = ecode.RoomNotFound
		return
	}

	list = make([]*v1pb.TagData, 0)

	tags, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil && err != redis.ErrNil {
		if e, ok := err.(*redis.Error); ok && strings.Index(e.Error(), "WRONGTYPE") == -1 {
			log.Error("redisGetTagList conn.Do(SMEMBERS, %s) error(%v)", key, err)
			return
		}
		return list, nil
	}

	for _, tag := range tags {
		seg := strings.Split(tag, ":")
		if len(seg) < 5 {
			log.Error("redisGetTagList Split(%s) error(%v)", tag, err)
			return nil, err
		}
		data := &v1pb.TagData{
			TagExt: strings.Join(seg[4:], ":"),
		}
		data.TagExpireAt, _ = strconv.ParseInt(seg[3], 10, 64)
		if data.TagExpireAt > time.Now().Unix() {
			data.TagId, _ = strconv.ParseInt(seg[0], 10, 64)
			data.TagSubId, _ = strconv.ParseInt(seg[1], 10, 64)
			data.TagValue, _ = strconv.ParseInt(seg[2], 10, 64)
			list = append(list, data)
		}
	}

	return
}

func (d *Dao) redisAddTag(ctx context.Context, roomID int64, tag *v1pb.TagData) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	key := fmt.Sprintf(_tagListKey, roomID)

	typ, err := redis.String(conn.Do("TYPE", key))
	if err != nil && err != redis.ErrNil {
		log.Error("redisAddTag conn.Do(TYPE, %s) error(%v)", key, err)
		return
	}

	if strings.ToLower(typ) == "none" {
		// 不存在就不处理
		return
	} else if strings.ToLower(typ) != "set" {
		if _, err = conn.Do("DEL", key); err != nil && err != redis.ErrNil {
			log.Error("redisAddTag conn.Do(DEL, %s) error(%v)", key, err)
			return
		}
	}

	tagVal := fmt.Sprintf("%d:%d:%d:%d:%s", tag.TagId, tag.TagSubId, tag.TagValue, tag.TagExpireAt, tag.TagExt)
	if _, err = conn.Do("SADD", key, tagVal); err != nil {
		log.Error("redisAddTag conn.Do(SADD, %s, %s) error(%v)", key, tagVal, err)
		return
	}

	return
}

func (d *Dao) redisSetTagList(ctx context.Context, roomID int64, list []*v1pb.TagData) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	key := fmt.Sprintf(_tagListKey, roomID)

	if len(list) <= 0 {
		// 设置哨兵
		if _, err = conn.Do("SETEX", key, d.c.Common.ExpireTime, "emptylist"); err != nil {
			log.Error("redisSetTagList conn.Do(SETEX, %s, %v) error(%v)", key, list, err)
		}
		return
	}

	if _, err = conn.Do("DEL", key); err != nil && err != redis.ErrNil {
		log.Error("redisSetTagList conn.Do(DEL, %s) error(%v)", key, err)
		return
	}

	args := make([]interface{}, len(list)+1)
	args[0] = key

	for i, tag := range list {
		args[i+1] = fmt.Sprintf("%d:%d:%d:%d:%s", tag.TagId, tag.TagSubId, tag.TagValue, tag.TagExpireAt, tag.TagExt)
	}

	if _, err = conn.Do("SADD", args...); err != nil {
		log.Error("redisSetTagList conn.Do(SADD, %s, %v) error(%v)", key, list, err)
		return
	}

	return
}
