package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// sql cachekey
const (
	_addUserFavSQL    = "INSERT INTO user_favorite (mid, item_id, type, status) VALUES(?, ?, ?, ?)"
	_selectUserFavSQL = "SELECT status, mtime FROM user_favorite WHERE mid=? AND type=? AND item_id=? LIMIT 1"
	_updateUserFavSQL = "UPDATE user_favorite SET status=? where mid=? AND type=? AND item_id=? AND mtime=? LIMIT 1"

	_userFavList  = "%d:USERFAVLIST"
	_userFavState = "%d:USERFAVSTATE:%d:%d"
)

// FavUpdate 收藏状态更新
func (d *Dao) FavUpdate(c context.Context, itemID int64, mid int64, typ int32, status int32) (err error) {
	res := d.db.Raw(_selectUserFavSQL, mid, typ, itemID)
	if err = res.Error; err != nil {
		log.Error("d.FavUpdate() error(%v)", err)
		return
	}

	var (
		oldStatus int32
		lastMtime xtime.Time
		notFound  bool
	)

	if err = res.Row().Scan(&oldStatus, &lastMtime); err != nil {
		if err == sql.ErrNoRows {
			notFound = true
		} else {
			log.Error("d.FavUpdate() res.Row().Scan() error(%v)", err)
			return
		}
	}

	if notFound {
		// 插入
		if err = d.db.Exec(_addUserFavSQL, mid, itemID, typ, status).Error; err != nil {
			log.Error("d.FavUpdate(%d, %d, %d, %d) Exec(%s) error(%v)", mid, itemID, typ, status, _addUserFavSQL, err)
			return
		}
	} else {
		// 更新
		if oldStatus == status {
			log.Info("d.FavUpdate(%d, %d, %d, %d) 前后状态相同 oldStatus %d ", mid, itemID, typ, status, oldStatus)
			return
		}

		if err = d.db.Exec(_updateUserFavSQL, status, mid, typ, itemID, lastMtime.Time().Format("2006-01-02 15:04:05")).Error; err != nil {
			log.Error("d.FavUpdate(%d, %d, %d, %d) Exec(%s) error(%v)", mid, itemID, typ, status, _updateUserFavSQL, err)
			return
		}
	}

	return
}

// UserFavStateCache 设置用户收藏缓存 支持多种类型
func (d *Dao) UserFavStateCache(c context.Context, itemID int64, mid int64, typ int32, status int32) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	_, err = conn.Do("SETEX", fmt.Sprintf(_userFavState, mid, itemID, typ), _expireHalfhour, status)
	if err != nil {
		log.Error("d.UserFavState(%d, %d, %d, %d) SETEX error(%v)", itemID, mid, typ, status, err)
		return
	}
	return
}

// UpdateFavListCache 更新收藏列表缓存
func (d *Dao) UpdateFavListCache(c context.Context, itemID int64, mid int64, status int32) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	if status == 1 {
		_, err = conn.Do("ZADD", fmt.Sprintf(_userFavList, mid), time.Now().Unix(), itemID)
	} else {
		_, err = conn.Do("ZREM", fmt.Sprintf(_userFavList, mid), itemID)
	}

	if err != nil {
		log.Error("d.UpdateFavListCache(%d, %d, %d) error(%v)", itemID, mid, status, err)
		return
	}
	return
}
