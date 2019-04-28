package reply

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/reply/conf"
	"go-common/app/interface/main/reply/model/reply"
	"go-common/app/interface/main/reply/model/xreply"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixIdx       = "i_"
	_prefixRtIdx     = "ri_"
	_prefixRpt       = "rt_"
	_prefixLike      = "l_"
	_prefixAuditIdx  = "ai_%d_%d"
	_prefixDialogIdx = "d_%d"

	// f_{折叠类型，根评论还是评论区}_{评论区ID或者根评论ID}
	_foldedReplyFmt = "f_%s_%d"

	_prefixSpamRec   = "sr_"
	_prefixSpamDaily = "sd_"
	_prefixSpamAct   = "sa_"
	_prefixTopOid    = "tro_"
)

const (
	_oidOverflow = 1 << 48
)

// RedisDao RedisDao
type RedisDao struct {
	redis         *redis.Pool
	expireRdsIdx  int
	expireRdsRpt  int
	expireRdsUC   int
	expireUserAct int
}

// NewRedisDao NewRedisDao
func NewRedisDao(c *conf.Redis) *RedisDao {
	r := &RedisDao{
		redis:         redis.NewPool(c.Config),
		expireRdsIdx:  int(time.Duration(c.IndexExpire) / time.Second),
		expireRdsRpt:  int(time.Duration(c.ReportExpire) / time.Second),
		expireRdsUC:   int(time.Duration(c.UserCntExpire) / time.Second),
		expireUserAct: int(time.Duration(c.UserActExpire) / time.Second),
	}
	return r
}

func keyRcntCnt(mid int64) string {
	return "rc_" + strconv.FormatInt(mid, 10)
}

func keyUpRcntCnt(mid int64) string {
	return "urc_" + strconv.FormatInt(mid, 10)
}

func keyDialogIdx(dialogID int64) string {
	return fmt.Sprintf(_prefixDialogIdx, dialogID)
}

func keyFolderIdx(kind string, ID int64) string {
	return fmt.Sprintf(_foldedReplyFmt, kind, ID)
}

func keyIdx(oid int64, tp, sort int8) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d_%d", _prefixIdx, oid, tp, sort)
	}
	return _prefixIdx + strconv.FormatInt((oid<<16)|(int64(tp)<<8)|int64(sort), 10)
}

func keyAuditIdx(oid int64, tp int8) string {
	return fmt.Sprintf(_prefixAuditIdx, oid, tp)
}

func keyRtIdx(rpID int64) string {
	return _prefixRtIdx + strconv.FormatInt(rpID, 10)
}

func keyRpt(mid int64, now time.Time) string {
	return _prefixRpt + strconv.FormatInt(mid, 10) + "_" + strconv.Itoa(now.Day())
}

func keyLike(rpID int64) string {
	return _prefixLike + strconv.FormatInt(rpID, 10)
}

func keySpamRpRec(mid int64) string {
	return _prefixSpamRec + strconv.FormatInt(mid, 10)
}

func keySpamRpDaily(mid int64) string {
	return _prefixSpamDaily + strconv.FormatInt(mid, 10)
}

func keySpamActRec(mid int64) string {
	return _prefixSpamAct + strconv.FormatInt(mid, 10)
}

func keyTopOid(tp int8) string {
	return _prefixTopOid + strconv.FormatInt(int64(tp), 10)
}

// Ping check connection success.
func (dao *RedisDao) Ping(c context.Context) (err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("GET", "PING")
	return
}

// AddFloorIndex add index by floor.
func (dao *RedisDao) AddFloorIndex(c context.Context, oid int64, tp int8, rs ...*reply.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	key := keyIdx(oid, tp, reply.SortByFloor)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, r := range rs {
		if err = conn.Send("ZADD", key, r.Floor, r.RpID); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(rs)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddCountIndex add index by count.
func (dao *RedisDao) AddCountIndex(c context.Context, oid int64, tp int8, rs ...*reply.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	key := keyIdx(oid, tp, reply.SortByCount)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, r := range rs {
		if err = conn.Send("ZADD", key, int64(r.RCount)<<32|(int64(r.Floor)&0xFFFFFFFF), r.RpID); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(rs)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddLikeIndex add index by like.
func (dao *RedisDao) AddLikeIndex(c context.Context, oid int64, tp int8, r *reply.Reply, rpt *reply.Report) (err error) {
	var (
		count  int
		rptCnt int
	)
	key := keyIdx(oid, tp, reply.SortByLike)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if r.Like >= 3 && (r.Attr&0x3 == 0) {
		if rpt != nil {
			rptCnt = rpt.Count
		}
		score := int64((float32(r.Like+2) / float32(r.Hate+rptCnt+4)) * 100)
		if err = conn.Send("ZADD", key, score<<32|(int64(r.RCount)&0xFFFFFFFF), r.RpID); err != nil {
			log.Error("conn.Send(ZADD %s,%d) error(%v)", key, r.RpID, err)
			return
		}
		count++
	} else if r.Like < 3 {
		if err = conn.Send("ZREM", key, r.RpID); err != nil {
			log.Error("conn.Send(ZREM %s,%d) error(%v)", key, r.RpID, err)
			return
		}
		count++
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelIndex delete reply index.
func (dao *RedisDao) DelIndex(c context.Context, rp *reply.Reply) (err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if rp.Root == 0 {
		key := keyIdx(rp.Oid, rp.Type, reply.SortByFloor)
		err = conn.Send("ZREM", key, rp.RpID)
		key = keyIdx(rp.Oid, rp.Type, reply.SortByCount)
		err = conn.Send("ZREM", key, rp.RpID)
		key = keyIdx(rp.Oid, rp.Type, reply.SortByLike)
		err = conn.Send("ZREM", key, rp.RpID)
	} else {
		key := keyRtIdx(rp.Root)
		err = conn.Send("ZREM", key, rp.RpID)
	}
	if err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if rp.Root == 0 {
		for i := 0; i < 3; i++ {
			if _, err = conn.Receive(); err != nil {
				log.Error("conn.Receive error(%v)", err)
				return
			}
		}
	} else {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
		}
	}
	return
}

// AddFloorIndexByRoot add root reply index by floor.
func (dao *RedisDao) AddFloorIndexByRoot(c context.Context, root int64, rs ...*reply.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	key := keyRtIdx(root)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, r := range rs {
		if err = conn.Send("ZADD", key, r.CTime, r.RpID); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(rs)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddLike add actions into redis
func (dao *RedisDao) AddLike(c context.Context, rpID int64, ras ...*reply.Action) (err error) {
	if len(ras) == 0 {
		return
	}
	key := keyLike(rpID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, r := range ras {
		if err = conn.Send("ZADD", key, r.CTime, r.Mid); err != nil {
			log.Error("conn.Send error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(ras)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// DelLike add actions into redis
func (dao *RedisDao) DelLike(c context.Context, rpID int64, ra *reply.Action) (err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", keyLike(rpID), ra.Mid); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

// ExpireLike set expire time for action.
func (dao *RedisDao) ExpireLike(c context.Context, rpID int64) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyLike(rpID), dao.expireRdsIdx)); err != nil {
		log.Error("conn.Receive() error(%v)", err)
	}
	return
}

func (dao *RedisDao) Range(c context.Context, oid int64, tp, sort int8, start, end int) (rpIds []int64, isEnd bool, err error) {
	key := keyIdx(oid, tp, sort)
	conn := dao.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZREVRANGE", key, start, end))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	err = redis.ScanSlice(values, &rpIds)
	if len(rpIds) > 0 && rpIds[len(rpIds)-1] == -1 {
		rpIds = rpIds[:len(rpIds)-1]
		isEnd = true
	}
	return
}

// CountReplies CountReplies
func (dao *RedisDao) CountReplies(c context.Context, oid int64, tp, sort int8) (count int, err error) {
	key := keyIdx(oid, tp, sort)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		log.Error("conn.Do(ZCARD, %s) error(%v)", key, err)
	}
	return
}

// UserAuditReplies return user audit replies.
func (dao *RedisDao) UserAuditReplies(c context.Context, mid, oid int64, tp int8) (rpIds []int64, err error) {
	key := keyAuditIdx(oid, tp)
	conn := dao.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGEBYSCORE", key, mid, mid))
	if err != nil {
		log.Error("conn.Do(RANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	err = redis.ScanSlice(values, &rpIds)
	return
}

// RangeByRoot range root's replyies.
func (dao *RedisDao) RangeByRoot(c context.Context, root int64, start, end int) (rpIds []int64, err error) {
	key := keyRtIdx(root)
	conn := dao.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, start, end))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	err = redis.ScanSlice(values, &rpIds)
	return
}

// RangeByOids range oids
func (dao *RedisDao) RangeByOids(c context.Context, oids []int64, tp, sort, start, end int8) (oidMap map[int64][]int64, miss []int64, err error) {
	oidMap = make(map[int64][]int64)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, oid := range oids {
		if err = conn.Send("EXPIRE", keyIdx(oid, tp, sort), dao.expireRdsIdx); err != nil {
			log.Error("conn.Send(EXPIRE) err(%v)", err)
			return
		}
		if err = conn.Send("ZREVRANGE", keyIdx(oid, tp, sort), start, end-1); err != nil {
			log.Error("conn.Send(ZREVRANGE) err(%v)", err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.SEND(FLUSH) err(%v)", err)
		return
	}
	for _, oid := range oids {
		var (
			rpids  []int64
			values []interface{}
		)
		if _, err = conn.Receive(); err != nil {
			log.Error("redis.Bool() err(%v)", err)
			return
		}
		if values, err = redis.Values(conn.Receive()); err != nil {
			log.Error("redis.Values() err(%v)", err)
			return
		}
		if len(values) == 0 {
			miss = append(miss, oid)
			continue
		}
		if err = redis.ScanSlice(values, &rpids); err != nil {
			log.Error("redis.ScanSlice() err(%v) ", err)
			return
		}
		oidMap[oid] = rpids
	}
	return
}

// RangeByRoots range roots's replyies.
func (dao *RedisDao) RangeByRoots(c context.Context, roots []int64, start, end int) (mrpids map[int64][]int64, idx, miss []int64, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, root := range roots {
		// if exist delay expire time
		if err = conn.Send("EXPIRE", keyRtIdx(root), dao.expireRdsIdx); err != nil {
			log.Error("conn.Send(EXPIRE) err(%v)", err)
			return
		}
		if err = conn.Send("ZRANGE", keyRtIdx(root), start, end); err != nil {
			log.Error("conn.Send(ZRANGE) err(%v)", err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.SEND(FLUSH) err(%v)", err)
		return
	}
	mrpids = make(map[int64][]int64, len(roots))
	for _, root := range roots {
		var (
			rpids  []int64
			values []interface{}
		)
		if _, err = conn.Receive(); err != nil {
			log.Error("redis.Bool() err(%v)", err)
			return
		}
		if values, err = redis.Values(conn.Receive()); err != nil {
			log.Error("redis.Values() err(%v)", err)
			return
		}
		if len(values) == 0 {
			miss = append(miss, root)
			continue
		}
		if err = redis.ScanSlice(values, &rpids); err != nil {
			log.Error("redis.ScanSlice() err(%v) ", err)
			return
		}
		idx = append(idx, rpids...)
		mrpids[root] = rpids
	}
	return
}

// ExpireIndex set expire time for index.
func (dao *RedisDao) ExpireIndex(c context.Context, oid int64, tp, sort int8) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyIdx(oid, tp, sort), dao.expireRdsIdx)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// ExpireIndexByRoot set expire time for root's index.
func (dao *RedisDao) ExpireIndexByRoot(c context.Context, root int64) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyRtIdx(root), dao.expireRdsIdx)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// SetUserReportCnt set user report count.
func (dao *RedisDao) SetUserReportCnt(c context.Context, mid int64, count int, now time.Time) (err error) {
	key := keyRpt(mid, now)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SETEX", key, dao.expireRdsRpt, count); err != nil {
		log.Error("conn.Do(SETEX) error(%v)", err)
	}
	return
}

// GetUserReportCnt get user report count.
func (dao *RedisDao) GetUserReportCnt(c context.Context, mid int64, now time.Time) (count int, err error) {
	key := keyRpt(mid, now)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET) error(%v)", err)
		}
	}
	return
}

// GetUserReportTTL get TTL of user report count redis.
func (dao *RedisDao) GetUserReportTTL(c context.Context, mid int64, now time.Time) (ttl int, err error) {
	key := keyRpt(mid, now)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ttl, err = redis.Int(conn.Do("TTL", key)); err != nil {
		log.Error("conn.Do(TTl) error(%v)", err)
	}
	return
}

// RankIndex get rank from reply index.
func (dao *RedisDao) RankIndex(c context.Context, oid int64, tp int8, rpID int64, sort int8) (rank int, err error) {
	key := keyIdx(oid, tp, sort)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if rank, err = redis.Int(conn.Do("ZREVRANK", key, rpID)); err != nil {
		if err == redis.ErrNil {
			rank = -1
			err = nil
		} else {
			log.Error("conn.Do(ZREVRANK) error(%v)", err)
		}

	}
	return
}

// RankIndexByRoot get rank from root reply index.
func (dao *RedisDao) RankIndexByRoot(c context.Context, root int64, rpID int64) (rank int, err error) {
	key := keyRtIdx(root)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if rank, err = redis.Int(conn.Do("ZRANK", key, rpID)); err != nil {
		if err == redis.ErrNil {
			rank = -1
			err = nil
		} else {
			log.Error("conn.Do(ZRANK) error(%v)", err)
		}
	}
	return
}

// OidHaveTop OidHaveTop
func (dao *RedisDao) OidHaveTop(c context.Context, oid int64, tp int8) (ok bool, err error) {
	key := keyTopOid(tp)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("SISMEMBER", key, oid)); err != nil {
		log.Error("OidHavaTop.Do error(%v)", err)
	}
	return
}

// DelReplyIncr reply cd key
func (dao *RedisDao) DelReplyIncr(c context.Context, mid int64, isUp bool) (err error) {
	key := keyRcntCnt(mid)
	if isUp {
		key = keyUpRcntCnt(mid)
	}
	conn := dao.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err != nil {
		log.Error("DelReplyIncr redis failed!err:=%v key:%s", err, key)
	}
	return
}

// DelReplyIncr reply cd key
func (dao *RedisDao) DelReplySpam(c context.Context, mid int64) (err error) {
	key := keySpamRpRec(mid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("DEL", key)
	if err != nil {
		log.Error("DelReplySpam redis failed!err:=%v key:%s", err, key)
	}
	return
}

// SpamReply SpamReply
func (dao *RedisDao) SpamReply(c context.Context, mid int64) (recent, daily int, err error) {
	rkey, dkey := keySpamRpRec(mid), keySpamRpDaily(mid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	ii, err := redis.Ints(conn.Do("MGET", rkey, dkey))
	if err != nil {
		log.Error("conn.Do(MGET, %s, %s) error(%v)", rkey, dkey, err)
		// no need for redis.ErrNil check
		return
	}
	if len(ii) != 2 {
		err = fmt.Errorf("ReplySpam redis result: %v, len not 2", ii)
		log.Error("%v", err)
		return
	}
	recent = ii[0]
	daily = ii[1]
	return
}

// SpamAction SpamAction
func (dao *RedisDao) SpamAction(c context.Context, mid int64) (code int, err error) {
	key := keySpamActRec(mid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if code, err = redis.Int(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET, %s, ), err (%v)", key, err)
		}
	}
	return
}

// ExpireDialogIndex expire time for dialog index
func (dao *RedisDao) ExpireDialogIndex(c context.Context, dialogID int64) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyDialogIdx(dialogID), dao.expireRdsIdx)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// RangeRpsByDialog return replies by dialog ID
func (dao *RedisDao) RangeRpsByDialog(c context.Context, dialog int64, start, end int) (rpIDs []int64, err error) {
	key := keyDialogIdx(dialog)
	conn := dao.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, start, end))
	if err != nil {
		log.Error("conn.Do(ZRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) == 0 {
		return
	}
	err = redis.ScanSlice(values, &rpIDs)
	if err != nil {
		log.Error("redis.ScanSlice Error (%v)", err)
	}
	return
}

// DialogBySide ...
func (dao *RedisDao) DialogDesc(c context.Context, dialog int64, floor, size int) (rpIDs []int64, err error) {
	var vals []interface{}
	key := keyDialogIdx(dialog)
	conn := dao.redis.Get(c)
	defer conn.Close()
	vals, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, floor, "-inf", "LIMIT", 0, size))
	if err = redis.ScanSlice(vals, &rpIDs); err != nil {
		log.Error("redis.ScanSlice Error (%v)", err)
		return nil, err
	}
	return
}

// DialogMaxMinFloor return min and max floor
func (dao *RedisDao) DialogMinMaxFloor(c context.Context, dialog int64) (minFloor, maxFloor int, err error) {
	var RpID int64
	key := keyDialogIdx(dialog)
	conn := dao.redis.Get(c)
	defer conn.Close()
	err = conn.Send("ZRANGE", key, 0, 0, "WITHSCORES")
	if err != nil {
		log.Error("redis.Send key(%s) error(%v)", key, err)
		return
	}
	err = conn.Send("ZRANGE", key, -1, -1, "WITHSCORES")
	if err != nil {
		log.Error("redis.Send key(%s) error(%v)", key, err)
		return
	}
	err = conn.Flush()
	if err != nil {
		log.Error("redis.Flush (%s) error(%v)", key, err)
		return
	}

	minValue, err := redis.Values(conn.Receive())
	if err != nil {
		log.Error("redis.Values key(%s) error(%v)", key, err)
		return
	}
	if _, err = redis.Scan(minValue, &RpID, &minFloor); err != nil {
		log.Error("redis.Scan() error(%v)", err)
		return
	}

	maxValue, err := redis.Values(conn.Receive())
	if err != nil {
		log.Error("redis.Values key(%s) error(%v)", key, err)
		return
	}
	if _, err = redis.Scan(maxValue, &RpID, &maxFloor); err != nil {
		log.Error("redis.Scan() error(%v)", err)
		return
	}
	return
}

// DialogByCursor return replies by dialog
func (dao *RedisDao) DialogByCursor(c context.Context, dialog int64, cursor *reply.Cursor) (rpIDs []int64, err error) {
	var vals []interface{}
	key := keyDialogIdx(dialog)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if cursor.Latest() {
		vals, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, 0, "+inf", "LIMIT", 0, cursor.Len()))
	} else if cursor.Descrease() {
		vals, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, cursor.Current(), "-inf", "LIMIT", 0, cursor.Len()))
	} else if cursor.Increase() {
		vals, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, cursor.Current(), "+inf", "LIMIT", 0, cursor.Len()))
	} else {
		err = ErrCursorDirection
	}
	if err = redis.ScanSlice(vals, &rpIDs); err != nil {
		log.Error("redis.ScanSlice() error(%v)", err)
		return nil, err
	}
	return
}

// ExpireFolder ...
func (dao *RedisDao) ExpireFolder(c context.Context, kind string, ID int64) (ok bool, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = keyFolderIdx(kind, ID)
	)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", key, dao.expireRdsIdx)); err != nil {
		log.Error("redis EXPIRE(%s) error(%v)", key, err)
	}
	return
}

// FolderByCursor ...
func (dao *RedisDao) FolderByCursor(c context.Context, kind string, ID int64, cursor *xreply.Cursor) (rpIDs []int64, err error) {
	var (
		conn = dao.redis.Get(c)
		key  = keyFolderIdx(kind, ID)
		vals []interface{}
	)
	defer conn.Close()
	if cursor.Latest() {
		vals, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "LIMIT", 0, cursor.Ps))
	} else if cursor.Forward() {
		vals, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, fmt.Sprintf("(%d", cursor.Next), "-inf", "LIMIT", 0, cursor.Ps))
	} else {
		vals, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, fmt.Sprintf("(%d", cursor.Prev), "+inf", "LIMIT", 0, cursor.Ps))
		// 这里保持一致都是降序的输出
		for left, right := 0, len(vals)-1; left < right; left, right = left+1, right-1 {
			vals[left], vals[right] = vals[right], vals[left]
		}
	}
	if err = redis.ScanSlice(vals, &rpIDs); err != nil {
		log.Error("redis.ScanSlice() error(%v)", err)
	}
	return
}
