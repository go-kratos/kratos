package dao

import (
	"context"

	"go-common/library/cache/redis"
	"go-common/library/log"
)

//IsWhietListUID 通过UID判读是否是白名单用户
func (d *Dao) IsWhietListUID(ctx context.Context, key string) (isWhite bool) {
	conn := d.whitelistredis.Get(ctx)
	defer conn.Close()

	var err error
	isWhite, err = redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Error("[DM]  GetWhietListByUID redis err:%+v", err)
		return true
	}
	return
}
