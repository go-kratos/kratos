package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/member/model"
	"go-common/library/cache/redis"
)

const (
	_expShard       = 10000
	_expAddedPrefix = "ea_%s_%d_%d"
	_expCoinPrefix  = "ecoin_%d_%d"
)
const (
	_share = "shareClick"
	_view  = "watch"
	_login = "login"
)

func expCoinKey(mid, day int64) string {
	return fmt.Sprintf(_expCoinPrefix, day, mid)
}

func expAddedKey(tp string, mid, day int64) string {
	return fmt.Sprintf(_expAddedPrefix, tp, day, mid/_expShard)
}

// StatCache get exp stat cache.
func (d *Dao) StatCache(c context.Context, mid, day int64) (st *model.ExpStat, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	conn.Send("GETBIT", expAddedKey(_login, mid, day), mid%_expShard)
	conn.Send("GETBIT", expAddedKey(_view, mid, day), mid%_expShard)
	conn.Send("GETBIT", expAddedKey(_share, mid, day), mid%_expShard)
	conn.Send("GET", expCoinKey(mid, day))
	err = conn.Flush()
	if err != nil {
		return
	}
	st = new(model.ExpStat)
	st.Login, err = redis.Bool(conn.Receive())
	if err != nil && err != redis.ErrNil {
		return
	}
	st.Watch, err = redis.Bool(conn.Receive())
	if err != nil && err != redis.ErrNil {
		return

	}
	st.Share, err = redis.Bool(conn.Receive())
	if err != nil && err != redis.ErrNil {
		return
	}
	st.Coin, err = redis.Int64(conn.Receive())
	if err != nil && err != redis.ErrNil {
		return
	}
	if err == redis.ErrNil {
		err = nil
	}
	return
}
