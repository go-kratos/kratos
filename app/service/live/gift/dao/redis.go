package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1pb "go-common/app/service/live/gift/api/grpc/v1"
	"go-common/app/service/live/gift/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"math/rand"
	"time"
)

func dailyBagKey(uid int64) string {
	return fmt.Sprintf("gift:daily_bag:%s:%d", time.Now().Format("20060102"), uid)
}

// GetDailyBagCache GetDailyBagCache
func (d *Dao) GetDailyBagCache(ctx context.Context, uid int64) (res []*v1pb.DailyBagResp_BagList, err error) {
	key := dailyBagKey(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	item, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(item, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(item), err)
	}
	return
}

// SetDailyBagCache SetDailyBagCache
func (d *Dao) SetDailyBagCache(ctx context.Context, uid int64, data []*v1pb.DailyBagResp_BagList, expire int64) (err error) {
	key := dailyBagKey(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	bs, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal(%v) err(%v)", data, err)
		return
	}
	_, err = conn.Do("SETEX", key, expire, bs)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

func dailyMedalBagKey(uid int64) string {
	return fmt.Sprintf("gift:medal:daily_gift_bag:%s:%d", time.Now().Format("20060102"), uid)
}

// GetMedalDailyBagCache GetMedalDailyBagCache
func (d *Dao) GetMedalDailyBagCache(ctx context.Context, uid int64) (res *model.BagGiftStatus, err error) {
	key := dailyMedalBagKey(uid)
	fmt.Println(key)
	res = &model.BagGiftStatus{}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	item, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(item, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(item), err)
	}
	return
}

// SetMedalDailyBagCache SetMedalDailyBagCache
func (d *Dao) SetMedalDailyBagCache(ctx context.Context, uid int64, data *model.BagGiftStatus, expire int64) (err error) {
	key := dailyMedalBagKey(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	bs, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal(%v) err(%v)", data, err)
		return
	}
	_, err = conn.Do("SETEX", key, expire, bs)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

func weekLevelBagKey(uid, level int64) string {
	_, week := time.Now().ISOWeek()
	return fmt.Sprintf("gift:level:week_gift_bag:%d:%d:%d", week, uid, level)
}

// GetWeekLevelBagCache GetWeekLevelBagCache
func (d *Dao) GetWeekLevelBagCache(ctx context.Context, uid, level int64) (res *model.BagGiftStatus, err error) {
	key := weekLevelBagKey(uid, level)
	res = &model.BagGiftStatus{}
	conn := d.redis.Get(ctx)
	defer conn.Close()
	item, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(item, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(item), err)
	}
	return
}

// SetWeekLevelBagCache SetWeekLevelBagCache
func (d *Dao) SetWeekLevelBagCache(ctx context.Context, uid, level int64, data *model.BagGiftStatus, expire int64) (err error) {
	key := weekLevelBagKey(uid, level)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	bs, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal(%v) err(%v)", data, err)
		return
	}
	_, err = conn.Do("SETEX", key, expire, bs)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

//Lock Lock
func (d *Dao) Lock(ctx context.Context, key string, ttl int, retry int, retryDelay int) (gotLock bool, lockValue string, err error) {

	if retry <= 0 {
		retry = 1
	}
	lockValue = "locked:" + randomString(5)
	retryTimes := 0
	conn := d.redis.Get(ctx)
	defer conn.Close()

	realKey := lockKey(key)

	for ; retryTimes < retry; retryTimes++ {
		var res interface{}
		res, err = conn.Do("SET", realKey, lockValue, "PX", ttl, "NX")
		if err != nil {
			log.Error("redis_lock failed:%s:%v", realKey, err)
			break
		}

		if res != nil {
			gotLock = true
			break
		}
		time.Sleep(time.Duration(retryDelay) * time.Millisecond)
	}
	return
}

// UnLock UnLock
func (d *Dao) UnLock(ctx context.Context, key string, lockValue string) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	realKey := lockKey(key)
	res, err := redis.String(conn.Do("GET", realKey))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	if res != lockValue {
		err = errors.New("unlock value error")
		return
	}

	_, err = conn.Do("DEL", realKey)

	return
}

//ForceUnLock UnLock without lockValue
func (d *Dao) ForceUnLock(ctx context.Context, key string) (err error) {
	realKey := lockKey(key)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("DEL", realKey)
	return
}

func lockKey(key string) string {
	return fmt.Sprintf("gift_lock:%s", key)
}

func randomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func bagIDCache(uid, giftID, expireAt int64) string {
	return fmt.Sprintf("bag_id:%d:%d:%d", uid, giftID, expireAt)
}

// GetBagIDCache GetBagIDCache
func (d *Dao) GetBagIDCache(ctx context.Context, uid, giftID, expireAt int64) (bagID int64, err error) {
	key := bagIDCache(uid, giftID, expireAt)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	bagID, err = redis.Int64(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	return
}

// SetBagIDCache SetBagIDCache
func (d *Dao) SetBagIDCache(ctx context.Context, uid, giftID, expireAt, bagID, expire int64) (err error) {
	key := bagIDCache(uid, giftID, expireAt)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", key, expire, bagID)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

func bagListKey(uid int64) string {
	return fmt.Sprintf("bag_list:%d", uid)
}

// GetBagListCache GetBagListCache
func (d *Dao) GetBagListCache(ctx context.Context, uid int64) (res []*model.BagGiftList, err error) {
	key := bagListKey(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	item, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			res = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	if err = json.Unmarshal(item, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(item), err)
	}
	return
}

// SetBagListCache SetBagListCache
func (d *Dao) SetBagListCache(ctx context.Context, uid int64, data []*model.BagGiftList, expire int64) (err error) {
	key := bagListKey(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	bs, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal(%v) err(%v)", data, err)
		return
	}
	_, err = conn.Do("SETEX", key, expire, bs)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

// ClearBagListCache ClearBagListCache
func (d *Dao) ClearBagListCache(ctx context.Context, uid int64) (err error) {
	key := bagListKey(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", key, err)
	}
	return
}

func bagNumKey(uid, giftID, expireAt int64) string {
	return fmt.Sprintf("bag_num:%d:%d:%d", uid, giftID, expireAt)
}

// SetBagNumCache SetBagNumCache
func (d *Dao) SetBagNumCache(ctx context.Context, uid, giftID, expireAt, giftNum, expire int64) (err error) {
	key := bagNumKey(uid, giftID, expireAt)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", key, expire, giftNum)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}

func vipMonthBag(uid int64) string {
	return fmt.Sprintf("gift:vip_month:%s:%d", time.Now().Format("200601"), uid)
}

// GetVipStatusCache GetVipStatusCache
func (d *Dao) GetVipStatusCache(ctx context.Context, uid int64) (status int64, err error) {
	key := vipMonthBag(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	status, err = redis.Int64(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	return
}

// ClearVipStatusCache ClearVipStatusCache
func (d *Dao) ClearVipStatusCache(ctx context.Context, uid int64) (err error) {
	key := vipMonthBag(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err != nil {
		log.Error("conn.Do(DEL, %s) error(%v)", key, err)
	}
	return
}

func giftBagStatus(uid int64) string {
	return fmt.Sprintf("gift:bag:status:%d", uid)
}

// GetBagStatusCache GetBagStatusCache
func (d *Dao) GetBagStatusCache(ctx context.Context, uid int64) (status int64, err error) {
	key := giftBagStatus(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	status, err = redis.Int64(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			status = -100 // means cache miss
		} else {
			log.Error("conn.Do(GET, %s) error(%v)", key, err)
		}
		return
	}
	return
}

// SetBagStatusCache SetBagStatusCache
func (d *Dao) SetBagStatusCache(ctx context.Context, uid, status int64, expire int64) (err error) {
	key := giftBagStatus(uid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", key, expire, status)
	if err != nil {
		log.Error("conn.Do(SETEX, %s) error(%v)", key, err)
	}
	return
}
