package dao

import (
	"context"
	"fmt"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

var (
	_selectUserShareToken = "select token from user_share_token where mid = ?;"
	_insertUserShareToken = "insert into user_share_token (mid, token) values (?, ?);"
)

// GetUserShareToken .
func (d *Dao) GetUserShareToken(ctx context.Context, mid int64) string {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	raw, err := redis.Bytes(conn.Do("GET", fmt.Sprintf(model.CacheKeyUserShareToken, mid)))
	if err == redis.ErrNil || raw == nil {
		rows, err := d.db.Query(ctx, _selectUserShareToken, mid)
		if err != nil {
			log.Errorv(ctx, log.KV("GetUserShareToken", err))
			return ""
		}

		var token string
		for rows.Next() {
			rows.Scan(&token)
		}
		return token
	}

	return string(raw)
}

// SetUserShareToken .
func (d *Dao) SetUserShareToken(ctx context.Context, mid int64, token string) (int64, error) {
	result, err := d.db.Exec(ctx, _insertUserShareToken, mid, token)
	fmt.Println(result, err)
	if err != nil {
		return 0, err
	}

	if n, _ := result.RowsAffected(); n > 0 {
		conn := d.redis.Get(ctx)
		defer conn.Close()
		conn.Do("SET", fmt.Sprintf(model.CacheKeyUserShareToken, mid), token)
	}

	return result.LastInsertId()
}
