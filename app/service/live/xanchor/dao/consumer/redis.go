package consumer

import (
	"context"

	"go-common/library/log"
)

//实时消费缓存设计，异步落地

const VALUE = "value"
const DATE = "date" //是否最新消息，是为1(需要刷新到DB) 否为0(不需要刷新到DB)
const DATE_1 = "1"

//Set 设置实时数据
func (d *Dao) Set(ctx context.Context, redisKey string, value string, timeOut int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("HMSET", redisKey, VALUE, value, DATE, DATE_1); err != nil {
		log.Error("consumer_set_err:key=%s;err=%v", redisKey, err)
		return
	}
	conn.Do("EXPIRE", redisKey, timeOut)
	return
}

//Incr  设置增加数据
func (d *Dao) Incr(ctx context.Context, redisKey string, num int64, timeOut int) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("HINCRBY", redisKey, VALUE, num, DATE, DATE_1); err != nil {
		log.Error("consumer_incr_err:key=%s;err=%v", redisKey, err)
		return
	}
	conn.Do("EXPIRE", redisKey, timeOut)
	return
}
