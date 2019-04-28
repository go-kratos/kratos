package dao

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"strconv"
	"strings"
)

/*
kv结构：
    streamName: room_id
hash结构如下：
    // 存储所有的流名
	room_id + name + name1,name2, name3
    // 默认直推
    room_id + name1:default + src
    // 一个流名下的直推
	room_id + name1:origin + src
    // 一个流名下的转推
	room_id + name1:streaming + src
    // 一个流名下的key
	room_id + name1:key + key
	。。。
*/

const (
	// kv => name: rid
	_streamNameKey = "mng:name:%s"
	// hash key => rid: field :value
	_streamRooIDKey         = "mng:rid:%d"
	_streamRoomFieldAllName = "mng:name:all"
	_streamRoomFieldDefault = "%s:default"
	_streamRoomFieldOrigin  = "%s:origin"
	_streamRoomFieldForward = "%s:forward"
	_streamRoomFieldSecret  = "%s:key"
	_streamRoomFieldOption  = "%s:options"
	_streamRoomFieldHot     = "%d:hot"
	_streamLastCDN          = "last:cdn:%d"
	// 切流记录
	_streamChangeSrc = "change:src:%d"
	//房间冷热流
	_streamRoomHot = "room:hot:%d"
	// 存12小时
	_streamExpireTime = 4 * 3600
	// stream_name map room_id 流名对应房间号 这个是不变的
	_nameExpireTime = 365 * 86400
	// 存一年
	_lastCDNExpireTime   = 365 * 86400
	_changeSrcExpireTime = 365 * 86400
)

func (d *Dao) getStreamNamekey(streamName string) string {
	return fmt.Sprintf(_streamNameKey, streamName)
}

func (d *Dao) getRoomIDKey(rid int64) string {
	return fmt.Sprintf(_streamRooIDKey, rid)
}

func (d *Dao) getRoomFieldDefaultKey(streamName string) string {
	return fmt.Sprintf(_streamRoomFieldDefault, streamName)
}

func (d *Dao) getRoomFieldOriginKey(streamName string) string {
	return fmt.Sprintf(_streamRoomFieldOrigin, streamName)
}

func (d *Dao) getRoomFieldForwardKey(streamName string) string {
	return fmt.Sprintf(_streamRoomFieldForward, streamName)
}

func (d *Dao) getRoomFieldSecretKey(streamName string) string {
	return fmt.Sprintf(_streamRoomFieldSecret, streamName)
}

func (d *Dao) getRoomFieldOption(streamName string) string {
	return fmt.Sprintf(_streamRoomFieldOption, streamName)
}

func (d *Dao) getRoomFieldHotKey(rid int64) string {
	return fmt.Sprintf(_streamRoomFieldHot, rid)
}

func (d *Dao) getRoomHotKey(rid int64) string {
	return fmt.Sprintf(_streamRoomHot, rid)
}

func (d *Dao) getLastCDNKey(rid int64) string {
	return fmt.Sprintf(_streamLastCDN, rid)
}

func (d *Dao) getChangeSrcKey(rid int64) string {
	return fmt.Sprintf(_streamChangeSrc, rid)
}

// CacheStreamFullInfo 从缓存取流信息， 可传入流名， 也可以传入rid
func (d *Dao) CacheStreamFullInfo(c context.Context, rid int64, sname string) (res *model.StreamFullInfo, err error) {
	if sname != "" {
		infos, err := d.CacheStreamRIDByName(c, sname)
		if err != nil {
			return nil, err
		}

		if infos == nil || infos.RoomID <= 0 {
			return nil, fmt.Errorf("can not find any info by sname =%s", sname)
		}

		rid = infos.RoomID
	}

	// 先从本地缓存中取
	res = d.loadStreamInfo(c, rid)
	if res != nil {
		log.Warn("get from local cache")
		return res, nil
	}

	conn := d.redis.Get(c)
	defer conn.Close()

	roomKey := d.getRoomIDKey(rid)
	hotKey := d.getRoomFieldHotKey(rid)
	// 先判断过期时间是否为-1， 为-1则删除
	//ttl, _ := redis.Int(conn.Do("TTL", roomKey))
	//if ttl == -1 {
	//	d.DeleteStreamByRIDFromCache(c, rid)
	//	return nil, nil
	//}

	values, err := redis.StringMap(conn.Do("HGETALL", roomKey))

	log.Warn("%v", values)

	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, nil
	}

	resp := model.StreamFullInfo{}

	hot, _ := strconv.ParseInt(values[hotKey], 10, 64)
	resp.Hot = hot

	allNames := strings.Split(values[_streamRoomFieldAllName], "|")
	for _, n := range allNames {
		defaultUpStream := d.getRoomFieldDefaultKey(n)
		originKey := d.getRoomFieldOriginKey(n)
		forwardKey := d.getRoomFieldForwardKey(n)
		key := d.getRoomFieldSecretKey(n)
		options := d.getRoomFieldOption(n)

		if values[originKey] != "" && values[forwardKey] != "" && values[key] != "" && values[defaultUpStream] != "" {
			resp.RoomID = rid

			base := model.StreamBase{}
			base.StreamName = n

			or, _ := strconv.ParseInt(values[originKey], 10, 64)
			base.Origin = or

			de, _ := strconv.ParseInt(values[defaultUpStream], 10, 64)
			base.DefaultUpStream = de

			if de == 0 {
				return nil, fmt.Errorf("default is 0")
			}

			forward, _ := strconv.ParseInt(values[forwardKey], 10, 64)
			var num int64
			for num = 256; num > 0; num /= 2 {
				if ((forward & num) == num) && (num != or) {
					base.Forward = append(base.Forward, num)
				}
			}

			op, _ := strconv.ParseInt(values[options], 10, 64)
			base.Options = op
			//这里判断是否有wmask mmask的流
			if 4&op == 4 {
				base.Wmask = true
			}
			if 8&op == 8 {
				base.Mmask = true
			}

			base.Key = values[key]

			if strings.Contains(n, "_bs_") {
				base.Type = 2
			} else {
				base.Type = 1
			}
			resp.List = append(resp.List, &base)
		}
	}

	if len(resp.List) > 0 {
		// 存储到local cache
		d.storeStreamInfo(c, &resp)

		return &resp, nil
	}
	return nil, nil
}

// AddCacheStreamFullInfo 修改缓存数据
func (d *Dao) AddCacheStreamFullInfo(c context.Context, id int64, stream *model.StreamFullInfo) error {
	if stream == nil || stream.RoomID <= 0 {
		return nil
	}

	conn := d.redis.Get(c)
	defer func() {
		//conn.Do("EXPIRE", d.getRoomIDKey(id), _streamExpireTime)
		conn.Close()
	}()

	streamExpireTime := _streamExpireTime

	rid := stream.RoomID
	roomKey := d.getRoomIDKey(rid)
	allName := ""
	len := 0
	for _, v := range stream.List {
		len++
		allName = fmt.Sprintf("%s%s|", allName, v.StreamName)

		// kv 设置流名和room_id映射关系
		nameKey := d.getStreamNamekey(v.StreamName)
		if err := conn.Send("SET", nameKey, rid); err != nil {
			return fmt.Errorf("conn.Do(set, %s, %d) error(%v)", nameKey, rid, err)
		}
		if err := conn.Send("EXPIRE", nameKey, _nameExpireTime); err != nil {
			return fmt.Errorf("conn.Do(EXPIRE, %s,%d) error(%v)", nameKey, _nameExpireTime, err)
		}

		// hash 设置room_id下的field 和key
		field := d.getRoomFieldDefaultKey(v.StreamName)
		if v.DefaultUpStream == 0 {
			return fmt.Errorf("rid= %v, default is 0", roomKey)
		}
		if err := conn.Send("HSET", roomKey, field, v.DefaultUpStream); err != nil {
			return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.DefaultUpStream, err)
		}

		if len == 1 {
			if err := conn.Send("EXPIRE", roomKey, streamExpireTime); err != nil {
				log.Infov(c, log.KV("conn.EXPIRE error", err.Error()))
				return fmt.Errorf("conn.Do(EXPIRE, %s, %d) error(%v)", roomKey, streamExpireTime, err)
			}
		}

		field = d.getRoomFieldOriginKey(v.StreamName)
		if err := conn.Send("HSET", roomKey, field, v.Origin); err != nil {
			return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Origin, err)
		}

		field = d.getRoomFieldForwardKey(v.StreamName)
		var num int64
		for _, f := range v.Forward {
			num += f
		}

		if err := conn.Send("HSET", roomKey, field, num); err != nil {
			return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Forward, err)
		}

		field = d.getRoomFieldSecretKey(v.StreamName)
		if err := conn.Send("HSET", roomKey, field, v.Key); err != nil {
			return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Key, err)
		}

		field = d.getRoomFieldOption(v.StreamName)
		if err := conn.Send("HSET", roomKey, field, v.Options); err != nil {
			return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Options, err)
		}
	}
	// 去除最后的|
	allName = strings.Trim(allName, "|")
	//log.Warn("%v", allName)
	if err := conn.Send("HSET", roomKey, _streamRoomFieldAllName, allName); err != nil {
		return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, _streamRoomFieldAllName, allName, err)
	}

	if err := conn.Flush(); err != nil {
		log.Infov(c, log.KV("conn.Flush error(%v)", err.Error()))
		return fmt.Errorf("conn.Flush error(%v)", err)
	}
	for i := 0; i < 7*len+2; i++ {
		if _, err := conn.Receive(); err != nil {
			log.Infov(c, log.KV("conn.Receive error(%v)", err.Error()))
			return fmt.Errorf("conn.Receive error(%v)", err)
		}
	}
	return nil
}

// CacheStreamRIDByName 根据流名查房间号
func (d *Dao) CacheStreamRIDByName(c context.Context, sname string) (res *model.StreamFullInfo, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	nameKey := d.getStreamNamekey(sname)

	rid, err := redis.Int64(conn.Do("GET", nameKey))

	if err != nil {
		return nil, err
	}

	if rid <= 0 {
		return nil, nil
	}

	res = &model.StreamFullInfo{
		RoomID: rid,
	}
	return res, nil
}

// AddCacheStreamRIDByName 增加缓存
func (d *Dao) AddCacheStreamRIDByName(c context.Context, sname string, stream *model.StreamFullInfo) error {
	return d.AddCacheStreamFullInfo(c, 0, stream)
}

// AddCacheMultiStreamInfo 批量增加redis
func (d *Dao) AddCacheMultiStreamInfo(c context.Context, res map[int64]*model.StreamFullInfo) error {
	conn := d.redis.Get(c)
	defer func() {
		//for _, stream := range res {
		//	conn.Do("EXPIRE", d.getRoomIDKey(stream.RoomID), _streamExpireTime)
		//}
		conn.Close()
	}()

	count := 0
	for _, stream := range res {
		streamExpireTime := _streamExpireTime
		rid := stream.RoomID
		roomKey := d.getRoomIDKey(rid)
		allName := ""

		len := 0
		for _, v := range stream.List {
			len++
			allName = fmt.Sprintf("%s%s|", allName, v.StreamName)

			// kv 设置流名和room_id映射关系
			nameKey := d.getStreamNamekey(v.StreamName)
			if err := conn.Send("SET", nameKey, rid); err != nil {
				return fmt.Errorf("conn.Do(set, %s, %d) error(%v)", nameKey, rid, err)
			}
			if err := conn.Send("EXPIRE", nameKey, _nameExpireTime); err != nil {
				return fmt.Errorf("conn.Do(EXPIRE, %s,%d) error(%v)", nameKey, _nameExpireTime, err)
			}

			// hash 设置room_id下的field 和key
			field := d.getRoomFieldDefaultKey(v.StreamName)
			if err := conn.Send("HSET", roomKey, field, v.DefaultUpStream); err != nil {
				return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.DefaultUpStream, err)
			}

			if len == 1 {
				if err := conn.Send("EXPIRE", roomKey, streamExpireTime); err != nil {
					return fmt.Errorf("conn.Do(EXPIRE, %s, %d) error(%v)", roomKey, streamExpireTime, err)
				}
			}

			field = d.getRoomFieldOriginKey(v.StreamName)
			if err := conn.Send("HSET", roomKey, field, v.Origin); err != nil {
				return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Origin, err)
			}

			field = d.getRoomFieldForwardKey(v.StreamName)
			var num int64
			for _, f := range v.Forward {
				num += f
			}

			if err := conn.Send("HSET", roomKey, field, num); err != nil {
				return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Forward, err)
			}

			field = d.getRoomFieldSecretKey(v.StreamName)
			if err := conn.Send("HSET", roomKey, field, v.Key); err != nil {
				return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Key, err)
			}

			field = d.getRoomFieldOption(v.StreamName)
			if err := conn.Send("HSET", roomKey, field, v.Options); err != nil {
				return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, field, v.Options, err)
			}
		}

		// 去除最后的|
		allName = strings.Trim(allName, "|")

		if err := conn.Send("HSET", roomKey, _streamRoomFieldAllName, allName); err != nil {
			return fmt.Errorf("conn.Do(HSET, %s, %s, %v) error(%v)", roomKey, _streamRoomFieldAllName, allName, err)
		}

		count += (7*len + 2)
	}

	if err := conn.Flush(); err != nil {
		return fmt.Errorf("conn.Flush error(%v)", err)
	}

	// 这里要len
	for i := 0; i < count; i++ {
		if _, err := conn.Receive(); err != nil {
			return fmt.Errorf("conn.Receive error(%v)", err)
		}
	}
	return nil
}

// CacheMultiStreamInfo 批量从redis中获取数据
func (d *Dao) CacheMultiStreamInfo(c context.Context, rids []int64) (res map[int64]*model.StreamFullInfo, err error) {
	if len(rids) == 0 {
		return nil, nil
	}

	infos := map[int64]*model.StreamFullInfo{}

	// 先从local cache读取
	localInfos, missRids := d.loadMultiStreamInfo(c, rids)

	//log.Warn("%v=%v", localInfos, missRids)
	// 若全部命中local cache,直接返回
	if len(missRids) == 0 {
		log.Warn("all hit local cache")
		return localInfos, nil
	}

	if len(localInfos) != 0 {
		infos = localInfos
	}

	conn := d.redis.Get(c)
	defer conn.Close()

	for _, id := range missRids {
		key := d.getRoomIDKey(id)
		if err := conn.Send("HGETALL", key); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("redis: conn.Send(HGETALL, %s) error(%v)", key, err)))
			return nil, fmt.Errorf("redis: conn.Send(HGETALL, %s) error(%v)", key, err)
		}
	}

	if err := conn.Flush(); err != nil {
		return nil, fmt.Errorf("redis: conn.Flush error(%v)", err)
	}

	for i := 0; i < len(missRids); i++ {
		if values, err := redis.StringMap(conn.Receive()); err == nil {
			if len(values) == 0 {
				continue
			}

			item := model.StreamFullInfo{}

			hotKey := d.getRoomFieldHotKey(missRids[i])
			hot, _ := strconv.ParseInt(values[hotKey], 10, 64)
			item.Hot = hot

			allNames := strings.Split(values[_streamRoomFieldAllName], "|")
			for _, n := range allNames {
				defaultUpStream := d.getRoomFieldDefaultKey(n)
				originKey := d.getRoomFieldOriginKey(n)
				forwardKey := d.getRoomFieldForwardKey(n)
				key := d.getRoomFieldSecretKey(n)
				options := d.getRoomFieldOption(n)

				if values[originKey] != "" && values[forwardKey] != "" && values[key] != "" {
					item.RoomID = missRids[i]

					base := model.StreamBase{}
					base.StreamName = n

					or, _ := strconv.ParseInt(values[originKey], 10, 64)
					base.Origin = or

					de, _ := strconv.ParseInt(values[defaultUpStream], 10, 64)
					base.DefaultUpStream = de

					forward, _ := strconv.ParseInt(values[forwardKey], 10, 64)
					var num int64
					for num = 256; num > 0; num /= 2 {
						if ((forward & num) == num) && (num != or) {
							base.Forward = append(base.Forward, num)
						}
					}

					op, _ := strconv.ParseInt(values[options], 10, 64)
					base.Options = op
					//这里判断是否有wmask mmask的流
					if 4&op == 4 {
						base.Wmask = true
					}
					if 8&op == 8 {
						base.Mmask = true
					}

					base.Key = values[key]
					if strings.Contains(n, "_bs_") {
						base.Type = 2
					} else {
						base.Type = 1
					}
					item.List = append(item.List, &base)
				}
			}

			if len(item.List) > 0 {
				//log.Warn("miss=%v", missRids[i])
				infos[missRids[i]] = &item
			}
		}
	}

	// 更新local cache
	d.storeMultiStreamInfo(c, infos)
	return infos, nil
}

// UpdateLastCDNCache 设置last cdn
func (d *Dao) UpdateLastCDNCache(c context.Context, rid int64, origin int64) error {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := d.getLastCDNKey(rid)
	if err := conn.Send("SET", key, origin); err != nil {
		return fmt.Errorf("redis: conn.Send(SET, %s, %v) error(%v)", key, origin, err)
	}
	if err := conn.Send("EXPIRE", key, _lastCDNExpireTime); err != nil {
		return fmt.Errorf("redis: conn.Send(EXPIRE key(%s) expire(%d)) error(%v)", key, _lastCDNExpireTime, err)
	}
	if err := conn.Flush(); err != nil {
		return fmt.Errorf("redis: conn.Flush error(%v)", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := conn.Receive(); err != nil {
			return fmt.Errorf("redis: conn.Receive error(%v)", err)
		}
	}

	return nil
}

// UpdateChangeSrcCache 切流
func (d *Dao) UpdateChangeSrcCache(c context.Context, rid int64, origin int64) error {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := d.getChangeSrcKey(rid)
	if err := conn.Send("SET", key, origin); err != nil {
		return fmt.Errorf("redis: conn.Send(SET, %s, %v) error(%v)", key, origin, err)
	}
	if err := conn.Send("EXPIRE", key, _changeSrcExpireTime); err != nil {
		return fmt.Errorf("redis: conn.Send(EXPIRE key(%s) expire(%d)) error(%v)", key, _changeSrcExpireTime, err)
	}
	if err := conn.Flush(); err != nil {
		return fmt.Errorf("redis: conn.Flush error(%v)", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := conn.Receive(); err != nil {
			return fmt.Errorf("redis: conn.Receive error(%v)", err)
		}
	}

	return nil
}

// UpdateStreamForwardStatus 更新forward值
func (d *Dao) UpdateStreamStatusCache(c context.Context, stream *model.StreamStatus) {
	var (
		exist   bool
		err     error
		conn    = d.redis.Get(c)
		allName string
	)
	defer conn.Close()

	// 首先判断是否存在
	key := d.getRoomIDKey(stream.RoomID)
	if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		if err == redis.ErrNil {
			return
		}

		log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
		d.DeleteStreamByRIDFromCache(c, stream.RoomID)
		return
	}

	if !exist {
		return
	}

	// 不传sname 默认为主流
	if stream.StreamName == "" {
		if allName, err = redis.String(conn.Do("HGET", key, _streamRoomFieldAllName)); err != nil {
			if err == redis.ErrNil {
				return
			}

			log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
			d.DeleteStreamByRIDFromCache(c, stream.RoomID)
			return
		}

		names := strings.Split(allName, "|")
		for _, v := range names {
			if !strings.Contains(v, "_bs_") {
				stream.StreamName = v
				break
			}
		}
	}

	count := 0
	// 是否是新增备用流
	if stream.Add {
		var names string

		if names, err = redis.String(conn.Do("HGET", key, _streamRoomFieldAllName)); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
			d.DeleteStreamByRIDFromCache(c, stream.RoomID)
			return
		}

		count++
		if err := conn.Send("HSET", key, _streamRoomFieldAllName, fmt.Sprintf("%s|%s", names, stream.StreamName)); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
			d.DeleteStreamByRIDFromCache(c, stream.RoomID)
			return
		}

		count++
		secret := d.getRoomFieldSecretKey(stream.StreamName)
		if err := conn.Send("HSET", key, secret, stream.Key); err != nil {
			log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
			d.DeleteStreamByRIDFromCache(c, stream.RoomID)
			return
		}

	}

	// 如果origin改变
	if stream.OriginChange {
		count++
		originKey := d.getRoomFieldOriginKey(stream.StreamName)

		if err := conn.Send("HSET", key, originKey, stream.Origin); err != nil {
			// 如果设置失败，则删除key
			log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
			d.DeleteStreamByRIDFromCache(c, stream.RoomID)
			return
		}
	}

	// forward
	if stream.ForwardChange {
		count++
		forwardKey := d.getRoomFieldForwardKey(stream.StreamName)
		if err := conn.Send("HSET", key, forwardKey, stream.Forward); err != nil {
			// 如果设置失败，则删除key
			log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
			d.DeleteStreamByRIDFromCache(c, stream.RoomID)
			return
		}
	}

	// 切上行
	if stream.DefaultChange {
		count++
		defaultUpKey := d.getRoomFieldDefaultKey(stream.StreamName)
		if err := conn.Send("HSET", key, defaultUpKey, stream.DefaultUpStream); err != nil {
			// 如果设置失败，则删除key
			log.Errorv(c, log.KV("log", fmt.Sprintf("update_status_err=%v", err)))
			d.DeleteStreamByRIDFromCache(c, stream.RoomID)
			return
		}
	}

	//切换options
	if stream.OptionsChange {
		if stream.Options >= 0 {
			count++
			optionsKey := d.getRoomFieldOption(stream.StreamName)
			if err := conn.Send("HSET", key, optionsKey, stream.Options); err != nil {
				// 如果设置失败，则删除key
				d.DeleteStreamByRIDFromCache(c, stream.RoomID)
				return
			}
		}
	}

	if err := conn.Flush(); err != nil {
		log.Infov(c, log.KV("conn.Flush error(%v)", err.Error()))
		return
	}

	for i := 0; i < count; i++ {
		if _, err := conn.Receive(); err != nil {
			log.Infov(c, log.KV("conn.Receive error(%v)", err.Error()))
			return
		}
	}
}

// GetLastCDNFromCache 查询上一次cdn
func (d *Dao) GetLastCDNFromCache(c context.Context, rid int64) (int64, error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := d.getLastCDNKey(rid)
	origin, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		if err != redis.ErrNil {
			return 0, fmt.Errorf("redis: conn.Do(GET, %s) error(%v)", key, err)
		}
	}

	if origin <= 0 {
		return 0, nil
	}

	return origin, nil
}

// GetChangeSrcFromCache 查询上一次cdn
func (d *Dao) GetChangeSrcFromCache(c context.Context, rid int64) (int64, error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := d.getChangeSrcKey(rid)
	origin, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		if err != redis.ErrNil {
			return 0, fmt.Errorf("redis: conn.Do(GET, %s) error(%v)", key, err)
		}
	}

	if origin <= 0 {
		return 0, nil
	}

	return origin, nil
}

// DeleteStreamByRIDFromCache 删除一个房间的的缓存信息
func (d *Dao) DeleteStreamByRIDFromCache(c context.Context, rid int64) (err error) {
	// todo 删除内存

	// 删除redis
	// redis删除失败，进行重试，确保缓存的信息是最新的
	conn := d.redis.Get(c)
	defer conn.Close()

	roomKey := d.getRoomIDKey(rid)
	for i := 0; i < 3; i++ {
		_, err = conn.Do("DEL", roomKey)
		if err != nil {
			log.Error("conn.Do(DEL, %s) error(%v)", roomKey, err)
			continue
		} else {
			return nil
		}
	}

	return err
}

// DeleteLastCDNFromCache 删除上一次到cdn
func (d *Dao) DeleteLastCDNFromCache(c context.Context, rid int64) error {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := d.getLastCDNKey(rid)
	_, err := conn.Do("DEL", key)

	return err
}

// UpdateRoomOptionsCache 更新Options状态
func (d *Dao) UpdateRoomOptionsCache(c context.Context, rid int64, streamname string, options int64) error {
	conn := d.redis.Get(c)
	defer conn.Close()

	roomKey := d.getRoomIDKey(rid)
	optionsKey := d.getRoomFieldOption(streamname)
	_, err := conn.Do("HSET", roomKey, optionsKey, options)
	return err
}

// UpdateRoomHotStatusCache 更新房间冷热流状态
func (d *Dao) UpdateRoomHotStatusCache(c context.Context, rid int64, hot int64) error {
	conn := d.redis.Get(c)
	defer conn.Close()

	roomKey := d.getRoomIDKey(rid)
	hotKey := d.getRoomFieldHotKey(rid)
	_, err := conn.Do("HSET", roomKey, hotKey, hot)
	return err
}
