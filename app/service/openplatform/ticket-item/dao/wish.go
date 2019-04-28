package dao

import (
	"context"

	"encoding/json"
	"fmt"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_addWishSQL = "INSERT INTO user_wish (mid, item_id, face) VALUES(?, ?, ?)"

	// 缓存 key
	_userWishCountKey  = "USER:WISH:COUNT:%d"
	_userWishActiveKey = "USER:WISH:ACTIVE:%d:%d"
	_userWishListKey   = "USER:WISH:LIST:%d"
)

// AddWish 添加想去
func (d *Dao) AddWish(c context.Context, wish *model.UserWish) (err error) {
	err = d.db.Exec(_addWishSQL, wish.MID, wish.ItemID, wish.Face).Error
	return
}

// WishCacheUpdate 更新想去缓存
func (d *Dao) WishCacheUpdate(c context.Context, wish *model.UserWish) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SETEX", fmt.Sprintf(_userWishActiveKey, wish.ItemID, wish.MID), _expireHalfhour, 1); err != nil {
		log.Error("d.WishCacheUpdate(%+v) SETEX error(%v)", wish, err)
		return
	}
	if _, err = conn.Do("INCR", fmt.Sprintf(_userWishCountKey, wish.ItemID)); err != nil {
		log.Error("d.WishCacheUpdate(%+v) INCR error(%v)", wish, err)
		return
	}

	wishCache, err := json.Marshal(map[string]interface{}{
		"mid":  wish.MID,
		"face": wish.Face,
	})
	if err != nil {
		log.Error("d.WishCacheUpdate(%+v) json.Marshal() error(%v)", wish, err)
		return
	}

	listKey := fmt.Sprintf(_userWishListKey, wish.ItemID)
	length, err := redis.Int64(conn.Do("RPUSH", listKey, wishCache))
	if err != nil {
		log.Error("d.WishCacheUpdate(%+v) RPUSH error(%v)", wish, err)
		return
	}
	if length > 5 {
		if _, err = conn.Do("LTRIM", listKey, -5, -1); err != nil {
			log.Error("d.WishCacheUpdate(%+v) LTRIM error(%v)", wish, err)
			return
		}
	}

	return
}
