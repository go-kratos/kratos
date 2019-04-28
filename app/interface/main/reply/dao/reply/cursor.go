package reply

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	selectRootIDsByLatestFloorSQL   = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) ORDER BY floor DESC limit 0,%d"
	selectRootIDsByCursorOnFloorSQL = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) AND floor %s %d ORDER BY floor limit 0,%d"
	selectRootIDsByRootStateSQL     = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state in (0,1,2,5,6) ORDER BY floor limit ?,?"

	selectChildrenIDsByLatestFloorSQL   = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state in (0,1,2,5,6) ORDER BY floor ASC limit 0,%d"
	selectChildrenIDsByCursorOnFloorSQL = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state in (0,1,2,5,6) AND floor %s %d ORDER BY floor limit 0,%d"
)

// ErrCursorDirection ErrCursorDirection
var ErrCursorDirection = errors.New("error cursor direction")

// ChildrenIDsOfRootReply ChildrenIDsOfRootReply
func (dao *RpDao) ChildrenIDsOfRootReply(ctx context.Context,
	oid, rootID int64, tp int8, offset, limit int) ([]int64, error) {
	rows, err := dao.dbSlave.Query(ctx,
		fmt.Sprintf(selectRootIDsByRootStateSQL, dao.hit(oid)),
		oid, tp, rootID, offset, limit)

	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("%v", err)
			return nil, err
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return nil, err
	}
	return ids, nil
}

// CacheKeyRootReplyIDs CacheKeyRootReplyIDs
func (dao *RedisDao) CacheKeyRootReplyIDs(oid int64, tp, sort int8) string {
	// subject:root_reply_ids
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d_%d", _prefixIdx, oid, tp, sort)
	}
	return _prefixIdx + strconv.FormatInt((oid<<16)|(int64(tp)<<8)|int64(sort), 10)
}

// ParentChildrenReplyIDMap ParentChildrenReplyIDMap
func (dao *RedisDao) ParentChildrenReplyIDMap(ctx context.Context,
	parentIDs []int64, start, end int) (parentChildrenMap map[int64][]int64, missedIDs []int64, err error) {

	parentChildrenMap = make(map[int64][]int64)
	arrayOfChildrenIDs, missedKeys, err := dao.RangeChildrenReplyIDs(ctx,
		genChildrenKeyByRootReplyIDs(parentIDs), start, end)
	if err != nil {
		return nil, nil, err
	}
	m := genChildrenKeyParentIDMap(parentIDs)
	for _, k := range missedKeys {
		if pid, ok := m[k]; ok {
			missedIDs = append(missedIDs, pid)
		}
	}
	for i, pid := range parentIDs {
		parentChildrenMap[pid] = arrayOfChildrenIDs[i]
	}
	return parentChildrenMap, missedIDs, nil
}

// RangeChildrenReplyIDs RangeChildrenReplyIDs
func (dao *RedisDao) RangeChildrenReplyIDs(ctx context.Context,
	keys []string, start, end int) (arrOfChildrenReplyIDs [][]int64, missedKeys []string, err error) {

	if len(keys) == 0 {
		return
	}

	conn := dao.redis.Get(ctx)
	defer conn.Close()

	for _, key := range keys {
		if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
			log.Error("%v", err)
			return nil, nil, err
		}
		if err = conn.Send("ZRANGE", key, start, end); err != nil {
			log.Error("%v", err)
			return nil, nil, err
		}
	}

	if err = conn.Flush(); err != nil {
		log.Error("%v", err)
		return nil, nil, err
	}

	arrOfChildrenReplyIDs = make([][]int64, 0)
	missedKeys = make([]string, 0)
	for _, key := range keys {
		if exists, err := redis.Bool(conn.Receive()); err != nil {
			log.Error("%v", err)
			return nil, nil, err
		} else if !exists {
			missedKeys = append(missedKeys, key)
		}
		values, err := redis.Values(conn.Receive())
		if err != nil {
			log.Error("%v", err)
			return nil, nil, err
		}
		if len(values) == 0 {
			arrOfChildrenReplyIDs = append(arrOfChildrenReplyIDs, []int64{})
			continue
		}
		var ids []int64
		if err = redis.ScanSlice(values, &ids); err != nil {
			log.Error("%v ", err)
			return nil, nil, err
		}
		arrOfChildrenReplyIDs = append(arrOfChildrenReplyIDs, ids)
	}
	return arrOfChildrenReplyIDs, missedKeys, nil
}

// RangeChildrenIDByCursorScore RangeChildrenIDByCursorScore
func (dao *RedisDao) RangeChildrenIDByCursorScore(ctx context.Context, key string, cursor *model.Cursor) ([]int64, error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	var (
		vals []interface{}
		err  error
	)
	if cursor.Latest() {
		vals, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, 0, "+inf", "LIMIT", 0, cursor.Len()))
	} else if cursor.Descrease() {
		vals, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, cursor.Current(), "-inf", "LIMIT", 0, cursor.Len()))
	} else if cursor.Increase() {
		vals, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, cursor.Current(), "+inf", "LIMIT", 0, cursor.Len()))
	} else {
		err = ErrCursorDirection
	}
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	replyIDs := make([]int64, 0)
	if err = redis.ScanSlice(vals, &replyIDs); err != nil {
		return nil, err
	}
	if cursor.Descrease() {
		// ZREVRANGEBYSCORE to ASC
		for i, j := 0, len(replyIDs)-1; i < j; i, j = i+1, j-1 {
			replyIDs[i], replyIDs[j] = replyIDs[j], replyIDs[i]
		}
	}
	return replyIDs, nil
}

// RangeRootIDByCursorScore RangeRootIDByCursorScore
func (dao *RedisDao) RangeRootIDByCursorScore(ctx context.Context, key string, cursor *model.Cursor) ([]int64, bool, error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()

	var (
		vals []interface{}
		err  error
	)
	if cursor.Latest() {
		vals, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, "+inf", 0, "LIMIT", 0, cursor.Len()))
	} else if cursor.Increase() {
		vals, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, cursor.Current(), "+inf", "LIMIT", 0, cursor.Len()))
	} else if cursor.Descrease() {
		vals, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, cursor.Current(), "-inf", "LIMIT", 0, cursor.Len()))
	} else {
		err = ErrCursorDirection
	}
	if err != nil {
		log.Error("%v", err)
		return nil, false, err
	}
	replyIDs := make([]int64, 0)
	if err = redis.ScanSlice(vals, &replyIDs); err != nil {
		return nil, false, err
	}
	if len(replyIDs) > 0 && replyIDs[len(replyIDs)-1] == -1 {
		replyIDs = replyIDs[:len(replyIDs)-1]
		return replyIDs, true, nil
	}
	return replyIDs, false, nil
}

// RangeRootReplyIDs RangeRootReplyIDs
func (dao *RedisDao) RangeRootReplyIDs(ctx context.Context, key string, start, end int) ([]int64, error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()

	vals, err := redis.Values(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	replyIDs := make([]int64, 0)
	if err = redis.ScanSlice(vals, &replyIDs); err != nil {
		log.Error("%v", err)
		return nil, err
	}
	return replyIDs, nil
}

// ExpireCache ExpireCache
func (dao *RedisDao) ExpireCache(ctx context.Context, key string) (bool, error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXPIRE", key, dao.expireRdsIdx))
	if err != nil {
		log.Error("%v", err)
		return false, err
	}
	return ok, nil
}

// genChildrenKeyParentIDMap genChildrenKeyParentIDMap
func genChildrenKeyParentIDMap(ids []int64) map[string]int64 {
	m := make(map[string]int64)
	for _, id := range ids {
		m[GenNewChildrenKeyByRootReplyID(id)] = id
	}
	return m
}

// genChildrenKeyByRootReplyIDs genChildrenKeyByRootReplyIDs
func genChildrenKeyByRootReplyIDs(ids []int64) []string {
	ks := make([]string, len(ids))
	for i, id := range ids {
		ks[i] = GenNewChildrenKeyByRootReplyID(id)
	}
	return ks
}

// GenNewChildrenKeyByRootReplyID GenNewChildrenKeyByRootReplyID
func GenNewChildrenKeyByRootReplyID(id int64) string {
	return _prefixRtIdx + strconv.FormatInt(id, 10)
}

// genChildrenKeyByRootReplyID genChildrenKeyByRootReplyID
func genChildrenKeyByRootReplyID(id int64) string {
	// score: timestamp
	// reply:id:ids
	return _prefixRtIdx + strconv.FormatInt(id, 10)
}

// genReplyKeyByID genReplyKeyByID
func genReplyKeyByID(id int64) string {
	// "subject:reply:ids"
	return _prefixRp + strconv.FormatInt(id, 10)
}

func contains(arr []string, b string) bool {
	for _, a := range arr {
		if a == b {
			return true
		}
	}
	return false
}

// GetReplyByIDs GetReplyByIDs
func (dao *MemcacheDao) GetReplyByIDs(ctx context.Context, ids []int64) ([]*model.Reply, []int64, error) {
	if len(ids) == 0 {
		return []*model.Reply{}, []int64{}, nil
	}
	keys := make([]string, len(ids))
	keyIDMap := make(map[string]int64)
	for i, id := range ids {
		key := genReplyKeyByID(id)
		keys[i] = key

		keyIDMap[key] = id
	}

	conn := dao.mc.Get(ctx)
	defer conn.Close()

	items, err := conn.GetMulti(keys)
	if err != nil {
		log.Error("%v", err)
		return nil, nil, err
	}

	foundKeys := make([]string, 0)
	for _, item := range items {
		foundKeys = append(foundKeys, item.Key)
	}

	missedKeys := make([]string, 0)
	if len(foundKeys) < len(keys) {
		for _, key := range keys {
			if !contains(foundKeys, key) {
				missedKeys = append(missedKeys, key)
			}
		}
	}

	missedIDs := make([]int64, 0)
	for _, mk := range missedKeys {
		missedIDs = append(missedIDs, keyIDMap[mk])
	}

	var rs = make([]*model.Reply, 0)
	for _, item := range items {
		rp := new(model.Reply)
		if err := conn.Scan(item, rp); err != nil {
			log.Error("%v", err)
			missedIDs = append(missedIDs, keyIDMap[item.Key])
			continue
		}
		rs = append(rs, rp)
	}
	return rs, missedIDs, nil
}

// ChildrenIDSortByFloorCursor ChildrenIDSortByFloorCursor
func (dao *RpDao) ChildrenIDSortByFloorCursor(ctx context.Context, oid int64, tp int8, rootID int64, cursor *model.Cursor) ([]int64, error) {
	var rawSQL string
	if cursor.Latest() {
		rawSQL = fmt.Sprintf(selectChildrenIDsByLatestFloorSQL, dao.hit(oid), cursor.Len())
	} else if cursor.Descrease() {
		rawSQL = fmt.Sprintf(selectChildrenIDsByCursorOnFloorSQL, dao.hit(oid), "<=", cursor.Current(), cursor.Len())
	} else if cursor.Increase() {
		rawSQL = fmt.Sprintf(selectChildrenIDsByCursorOnFloorSQL, dao.hit(oid), ">=", cursor.Current(), cursor.Len())
	} else {
		log.Error("%v", ErrCursorDirection)
		return nil, ErrCursorDirection
	}
	rows, err := dao.dbSlave.Query(ctx, rawSQL, oid, tp, rootID)
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	defer rows.Close()

	var id int64
	res := make([]int64, 0, cursor.Len())
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("%v", err)
			return nil, err
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return nil, err
	}
	return res, nil
}

// RootIDSortByFloorCursor RootIDSortByFloorCursor
func (dao *RpDao) RootIDSortByFloorCursor(ctx context.Context, oid int64, tp int8, cursor *model.Cursor) ([]int64, error) {
	var rawSQL string
	if cursor.Latest() {
		rawSQL = fmt.Sprintf(selectRootIDsByLatestFloorSQL, dao.hit(oid), cursor.Len())
	} else if cursor.Increase() {
		rawSQL = fmt.Sprintf(selectRootIDsByCursorOnFloorSQL, dao.hit(oid), ">=", cursor.Current(), cursor.Len())
	} else if cursor.Descrease() {
		rawSQL = fmt.Sprintf(selectRootIDsByCursorOnFloorSQL, dao.hit(oid), "<=", cursor.Current(), cursor.Len())
	} else {
		log.Error("%v", ErrCursorDirection)
		return nil, ErrCursorDirection
	}
	rows, err := dao.dbSlave.Query(ctx, rawSQL, oid, tp)
	if err != nil {
		log.Error("%v", err)
		return nil, err
	}
	defer rows.Close()

	var id int64
	ids := make([]int64, 0, cursor.Len())
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("%v", err)
			return nil, err
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return nil, err
	}
	return ids, nil
}
