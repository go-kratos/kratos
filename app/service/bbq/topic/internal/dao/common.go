package dao

import (
	"context"
	"encoding/json"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

/*各种情况下，是否需要查询db，前三列表示情景条件
NEXT	offset		rank		是否需要db操作
1									1
0		 1			  0				0
0		n>1			  0				1
0		 0			 n>0			0
*/
func parseCursor(ctx context.Context, cursorPrev, cursorNext string) (cursor model.CursorValue, directionNext bool, err error) {
	// 判断是向前还是向后查询
	directionNext = true
	cursorStr := cursorNext
	if len(cursorNext) == 0 && len(cursorPrev) > 0 {
		directionNext = false
		cursorStr = cursorPrev
	}
	// 解析cursor中的cursor_id
	if len(cursorStr) != 0 {
		var cursorData = []byte(cursorStr)
		err = json.Unmarshal(cursorData, &cursor)
		if err != nil {
			err = ecode.ReqParamErr
			return
		}
	}

	// 最后做一次校验，保证cursor的值是对的
	if (cursor.StickRank > 0 && cursor.Offset > 0) || (!directionNext && cursor.Offset == 0 && cursor.StickRank == 0) {
		err = ecode.TopicReqParamErr
		log.Errorw(ctx, "log", "cursor value error", "prev", cursorPrev, "next", cursorNext)
		return
	}

	return
}

func (d *Dao) getRedisList(ctx context.Context, key string) (list []int64, err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	str, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		err = nil
		return
	}
	if err != nil {
		log.Errorw(ctx, "log", "get redis list fail", "key", key)
		return
	}
	list, err = xstr.SplitInts(str)
	if err != nil {
		log.Errorw(ctx, "log", "split list_str fail", "key", key, "str", str)
		return
	}
	return
}
func (d *Dao) setRedisList(ctx context.Context, key string, list []int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()
	if _, err = conn.Do("SET", key, xstr.JoinInts(list)); err != nil {
		log.Errorw(ctx, "log", "set redis list fail", "key", key, "list", list)
		return
	}
	return
}
