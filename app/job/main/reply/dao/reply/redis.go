package reply

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go-common/app/job/main/reply/conf"
	"go-common/app/job/main/reply/model/reply"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_prefixIdx       = "i_"
	_prefixNewRtIdx  = "ri_"
	_prefixRpt       = "rt_"
	_prefixLike      = "l_"
	_prefixUAct      = "uha_"
	_prefixAuditIdx  = "ai_%d_%d"
	_prefixDialogIdx = "d_%d"

	_prefixSpamRec    = "sr_"
	_prefixSpamDaily  = "sd_"
	_prefixSpamAct    = "sa_"
	_prefixTopOid     = "tro_"
	_prefixNotifyCnt  = "rn_%d_%d"
	_prefixMaxLikeCnt = "mlc_%d"

	_oidOverflow = 1 << 48
	_maxCount    = 20000
	_maxHotReply = 2000

	// f_{折叠类型，根评论还是评论区}_{评论区ID或者根评论ID}
	_foldedReplyFmt = "f_%s_%d"
)

type rItem struct {
	ID    int64
	Score int64
}

// RedisDao define redis dao
type RedisDao struct {
	redis         *redis.Pool
	expireRdsIdx  int
	expireRdsRpt  int
	expireRdsUC   int
	expireUserAct int
	expireNotify  int
}

// NewRedisDao new redis dao
func NewRedisDao(c *conf.Redis) *RedisDao {
	return &RedisDao{
		redis:         redis.NewPool(c.Config),
		expireRdsIdx:  int(time.Duration(c.IndexExpire) / time.Second),
		expireRdsRpt:  int(time.Duration(c.ReportExpire) / time.Second),
		expireRdsUC:   int(time.Duration(c.UserCntExpire) / time.Second),
		expireUserAct: int(time.Duration(c.UserActExpire) / time.Second),
		expireNotify:  int(time.Duration(c.NotifyExpire) / time.Second),
	}
}

func keyFolderIdx(kind string, ID int64) string {
	return fmt.Sprintf(_foldedReplyFmt, kind, ID)
}

func keyDialogIdx(dialogID int64) string {
	return fmt.Sprintf(_prefixDialogIdx, dialogID)
}

func keyIdx(oid int64, tp, sort int8) string {
	if oid > _oidOverflow {
		return fmt.Sprintf("%s_%d_%d_%d", _prefixIdx, oid, tp, sort)
	}
	return _prefixIdx + strconv.FormatInt((oid<<16)|(int64(tp)<<8)|int64(sort), 10)
}

func keyNewRtIdx(rpID int64) string {
	return _prefixNewRtIdx + strconv.FormatInt(rpID, 10)
}

func keyAuditIdx(oid int64, tp int8) string {
	return fmt.Sprintf(_prefixAuditIdx, oid, tp)
}

func keyRpt(mid int64, now time.Time) string {
	return _prefixRpt + strconv.FormatInt(mid, 10) + "_" + strconv.Itoa(now.Day())
}

func keyLike(rpID int64) string {
	return _prefixLike + strconv.FormatInt(rpID, 10)
}

func keyUAct(mid int64) string {
	return _prefixUAct + strconv.FormatInt(mid, 10)
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

func keyNotifyCnt(oid int64, typ int8) string {
	return fmt.Sprintf(_prefixNotifyCnt, oid, typ)
}

func keyMaxLikeCnt(rpid int64) string {
	return fmt.Sprintf(_prefixMaxLikeCnt, rpid)
}

// Ping check connection success.
func (dao *RedisDao) Ping(c context.Context) (err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("GET", "PING")
	return
}

// DelAuditIndexs delete aduit reply cache.
func (dao *RedisDao) DelAuditIndexs(c context.Context, rs ...*reply.Reply) (err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, rp := range rs {
		key := keyAuditIdx(rp.Oid, rp.Type)
		if err = conn.Send("ZREM", key, rp.RpID); err != nil {
			log.Error("conn.Send(ZREM %s) error(%v)", key, err)
			return
		}
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(rs); i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddAuditIndex add audit reply index by user.
func (dao *RedisDao) AddAuditIndex(c context.Context, rp *reply.Reply) (err error) {
	key := keyAuditIdx(rp.Oid, rp.Type)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, rp.Mid, rp.RpID); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddFloorIndexEnd AddFloorIndexEnd
func (dao *RedisDao) AddFloorIndexEnd(c context.Context, oid int64, tp int8) (err error) {
	key := keyIdx(oid, tp, reply.SortByFloor)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, -1, -1); err != nil {
		log.Error("conn.Send ZADD error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
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

// AddCountIndexBatch add index by count.
func (dao *RedisDao) AddCountIndexBatch(c context.Context, oid int64, tp int8, rs ...*reply.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	if len(rs) > _maxCount {
		sort.Slice(rs, func(i, j int) bool {
			return rs[i].RCount > rs[j].RCount
		})
		rs = rs[:_maxCount]
	}
	key := keyIdx(oid, tp, reply.SortByCount)
	conn := dao.redis.Get(c)
	defer conn.Close()
	var count int
	for _, r := range rs {
		if !r.IsTop() {
			if err = conn.Send("ZADD", key, int64(r.RCount)<<32|(int64(r.Floor)&0xFFFFFFFF), r.RpID); err != nil {
				log.Error("conn.Send error(%v)", err)
				return
			}
			count++
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
	for i := 0; i < count+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddCountIndex add index by count.
func (dao *RedisDao) AddCountIndex(c context.Context, oid int64, tp int8, rp *reply.Reply) (err error) {
	var count int
	if count, err = dao.CountReplies(c, oid, tp, reply.SortByCount); err != nil {
		return
	} else if count >= _maxCount {
		var min int
		if min, err = dao.MinScore(c, oid, tp, reply.SortByCount); err != nil {
			return
		}
		if rp.RCount <= min {
			return
		}
	}

	key := keyIdx(oid, tp, reply.SortByCount)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, int64(rp.RCount)<<32|(int64(rp.Floor)&0xFFFFFFFF), rp.RpID); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddLikeIndexBatch add index by like.
func (dao *RedisDao) AddLikeIndexBatch(c context.Context, oid int64, tp int8, rpts map[int64]*reply.Report, rs ...*reply.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	items := make([]rItem, 0, len(rs))
	for _, r := range rs {
		if !r.IsTop() && r.IsRoot() {
			var rptCn int
			if rpt, ok := rpts[r.RpID]; ok {
				rptCn = rpt.Count
			}
			score := int64((float32(r.Like+conf.Conf.Weight.Like) / float32(r.Hate+conf.Conf.Weight.Hate+rptCn)) * 100)
			score = score<<32 | (int64(r.RCount) & 0xFFFFFFFF)
			items = append(items, rItem{ID: r.RpID, Score: score})
		}
	}
	if len(items) > _maxCount {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Score > items[j].Score
		})
		items = items[:_maxCount]
	}

	key := keyIdx(oid, tp, reply.SortByLike)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, item := range items {
		if err = conn.Send("ZADD", key, item.Score, item.ID); err != nil {
			log.Error("conn.Send(ZADD %s,%d) error(%v)", key, item.ID, err)
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
	for i := 0; i < len(items)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddLikeIndex add index by like.
func (dao *RedisDao) AddLikeIndex(c context.Context, oid int64, tp int8, rpts map[int64]*reply.Report, r *reply.Reply) (err error) {
	if r.IsTop() && r.IsRoot() {
		return
	}
	var rptCn int
	if rpt, ok := rpts[r.RpID]; ok {
		rptCn = rpt.Count
	}
	score := int64((float32(r.Like+conf.Conf.Weight.Like) / float32(r.Hate+conf.Conf.Weight.Hate+rptCn)) * 100)
	score = score<<32 | (int64(r.RCount) & 0xFFFFFFFF)
	key := keyIdx(oid, tp, reply.SortByLike)
	var count int
	if count, err = dao.CountReplies(c, oid, tp, reply.SortByLike); err != nil {
		return
	} else if count >= _maxCount {
		var min int
		if min, err = dao.MinScore(c, oid, tp, reply.SortByLike); err != nil {
			return
		}
		if score <= int64(min) {
			return
		}
	}

	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", key, score, r.RpID); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// AddIndex add reply index to redis.
func (dao *RedisDao) AddIndex(c context.Context, oid int64, tp int8, rpt *reply.Report, rp *reply.Reply, isRecover bool) (err error) {
	var (
		ok   bool
		isRt = rp.Root == 0 && rp.Parent == 0
	)
	if isRt {
		if ok, err = dao.ExpireIndex(c, oid, tp, reply.SortByFloor); err == nil && ok {
			if isRecover {
				var min int
				min, err = dao.MinScore(c, oid, tp, reply.SortByFloor)
				if err != nil {
					log.Error("s.dao.Redis.AddFloorIndex failed , oid(%d) type(%d) err(%v)", oid, tp, err)
				}
				if err == nil && rp.Floor > min {
					if err = dao.AddFloorIndex(c, oid, tp, rp); err != nil {
						log.Error("s.dao.Redis.AddFloorIndex failed , oid(%d) type(%d) err(%v)", oid, tp, err)
					}
				}
			} else {
				if err = dao.AddFloorIndex(c, oid, tp, rp); err != nil {
					log.Error("s.dao.Redis.AddFloorIndex failed , oid(%d) type(%d) err(%v)", oid, tp, err)
				}
			}
		}
		if ok, err = dao.ExpireIndex(c, oid, tp, reply.SortByCount); err == nil && ok {
			if err = dao.AddCountIndex(c, oid, tp, rp); err != nil {
				log.Error("s.dao.Redis.AddCountIndex failed , oid(%d) type(%d) err(%v)", oid, tp, err)
			}
		}
		if ok, err = dao.ExpireIndex(c, oid, tp, reply.SortByLike); err == nil && ok {
			rpts := make(map[int64]*reply.Report, 1)
			if rpt != nil {
				rpts[rpt.RpID] = rpt
			}
			if err = dao.AddLikeIndex(c, oid, tp, rpts, rp); err != nil {
				log.Error("s.dao.Redis.AddLikeIndex failed , oid(%d) type(%d) err(%v)", oid, tp, err)
			}
		}
	} else {
		if ok, err := dao.ExpireNewChildIndex(c, rp.Root); err == nil && ok {
			if err = dao.AddNewChildIndex(c, rp.Root, rp); err != nil {
				log.Error("s.dao.Redis.AddFloorIndexByRoot failed , rproot(%d),  err(%v)", rp.Root, err)
			}
		}
	}
	return
}

// DelIndexBySortType del index by sort.
func (dao *RedisDao) DelIndexBySortType(c context.Context, rp *reply.Reply, sortType int8) (err error) {
	key := keyIdx(rp.Oid, rp.Type, sortType)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", key, rp.RpID); err != nil {
		log.Error("redisDao.ZREM error(%v)", err)
	}
	return
}

// DelIndex delete reply index.
func (dao *RedisDao) DelIndex(c context.Context, rp *reply.Reply) (err error) {
	var (
		key string
		n   int
	)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if rp.Root == 0 {
		key = keyIdx(rp.Oid, rp.Type, reply.SortByFloor)
		err = conn.Send("ZREM", key, rp.RpID)
		key = keyIdx(rp.Oid, rp.Type, reply.SortByCount)
		err = conn.Send("ZREM", key, rp.RpID)
		key = keyIdx(rp.Oid, rp.Type, reply.SortByLike)
		err = conn.Send("ZREM", key, rp.RpID)
		n += 3
	} else {
		key = keyNewRtIdx(rp.Root)
		err = conn.Send("ZREM", key, rp.RpID)
		n++
		if rp.Dialog != 0 {
			key = keyDialogIdx(rp.Dialog)
			err = conn.Send("ZREM", key, rp.RpID)
			n++
		}
	}
	if err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < n; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// AddNewChildIndex add root reply index by floor.
func (dao *RedisDao) AddNewChildIndex(c context.Context, root int64, rs ...*reply.Reply) (err error) {
	if len(rs) == 0 {
		return
	}
	key := keyNewRtIdx(root)
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

//AddTopOid add oid to set  twhich has top reply
func (dao *RedisDao) AddTopOid(c context.Context, oid int64, tp int8) (err error) {
	key := keyTopOid(tp)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SADD", key, oid); err != nil {
		log.Error("redisDao.SADD error(%v)", err)
	}
	return
}

// DelTopOid delete oid from set
func (dao *RedisDao) DelTopOid(c context.Context, oid int64, tp int8) (err error) {
	key := keyTopOid(tp)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("SREM", key, oid); err != nil {
		log.Error("redisDao.SREM error(%v)", err)
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

// DelLike del user like from redis
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

// Range range replyies.
func (dao *RedisDao) Range(c context.Context, oid int64, tp, sort int8, start, end int) (rpIds []int64, err error) {
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
	return
}

// FloorEnd FloorEnd
func (dao *RedisDao) FloorEnd(c context.Context, oid int64, tp int8) (score int, found bool, err error) {
	key := keyIdx(oid, tp, reply.SortByFloor)
	conn := dao.redis.Get(c)
	defer conn.Close()
	err = conn.Send("ZSCORE", key, -1)
	if err != nil {
		log.Error("conn.Send ZSCORE error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	if v, e := redis.Int(conn.Receive()); e == nil {
		score = v
		found = true
	} else if err != redis.ErrNil {
		log.Error("redis.Int64(conn.Receive()) error(%v)", err)
		err = e
	}
	return
}

// MinScore get the lowest score from sorted set
func (dao *RedisDao) MinScore(c context.Context, oid int64, tp int8, sort int8) (score int, err error) {
	key := keyIdx(oid, tp, sort)
	conn := dao.redis.Get(c)
	defer conn.Close()
	values, err := redis.Values(conn.Do("ZRANGE", key, 0, 0, "WITHSCORES"))
	if err != nil {
		log.Error("conn.Do(ZREVRANGE, %s) error(%v)", key, err)
		return
	}
	if len(values) != 2 {
		err = fmt.Errorf("redis zrange items(%v) length not 2", values)
		return
	}
	var id int64
	redis.Scan(values, &id, &score)
	return
}

// CountReplies get count of reply.
func (dao *RedisDao) CountReplies(c context.Context, oid int64, tp, sort int8) (count int, err error) {
	key := keyIdx(oid, tp, sort)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if count, err = redis.Int(conn.Do("ZCARD", key)); err != nil {
		log.Error("CountReplies error(%v)", err)
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

// ExpireIndex set expire time for index.
func (dao *RedisDao) ExpireIndex(c context.Context, oid int64, tp, sort int8) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyIdx(oid, tp, sort), dao.expireRdsIdx)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// ExpireNewChildIndex set expire time for root's index.
func (dao *RedisDao) ExpireNewChildIndex(c context.Context, root int64) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyNewRtIdx(root), dao.expireRdsIdx)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// AddDialogIndex add reply to a dialog
func (dao *RedisDao) AddDialogIndex(c context.Context, dialogID int64, rps []*reply.Reply) (err error) {
	key := keyDialogIdx(dialogID)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for _, rp := range rps {
		if err = conn.Send("ZADD", key, rp.Floor, rp.RpID); err != nil {
			log.Error("conn.Send ZADD error(%v)", err)
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
	for i := 0; i < len(rps)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
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

// ExpireUserAct set expire time for user actions.
func (dao *RedisDao) ExpireUserAct(c context.Context, mid int64) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyUAct(mid), dao.expireRdsIdx)); err != nil {
		log.Error("conn.DO(EXPIRE) error(%v)", err)
	}
	return
}

// AddUserActs add user actions into redis.
func (dao *RedisDao) AddUserActs(c context.Context, mid int64, actions map[int64]int8) (err error) {
	key := keyUAct(mid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	for rpID, act := range actions {
		if err = conn.Send("HSET", key, rpID, act); err != nil {
			log.Error("conn.Send(HSET) error(%v)", err)
			return
		}
	}
	if err = conn.Send("EXPIRE", key, dao.expireUserAct); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < len(actions)+1; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// DelUserAct add user actions into redis.
func (dao *RedisDao) DelUserAct(c context.Context, mid int64, rpID int64) (err error) {
	key := keyUAct(mid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("HDEL", key, rpID); err != nil {
		log.Error("conn.SREM error(%v)", err)
	}
	return
}

// UserAct get user action from redis.
func (dao *RedisDao) UserAct(c context.Context, mid int64, rpID int64) (act int, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if act, err = redis.Int(conn.Do("HGET", keyUAct(mid), rpID)); err != nil {
		log.Error("conn.HGET(mid:%d) error(%v)", mid, err)
	}
	return
}

// UserActs get user actions from redis.
// NOTE: HGETALL quicker than HMEGT BUT transfer more data
func (dao *RedisDao) UserActs(c context.Context, mid int64, rpids []int64) (acts map[int64]int8, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	var tmpAct map[string]int
	if tmpAct, err = redis.IntMap(conn.Do("HGETALL", keyUAct(mid))); err != nil {
		log.Error("redis.IntMap(mid:%d)err(%v)", mid, err)
	}
	acts = make(map[int64]int8, len(rpids))
	for _, rpID := range rpids {
		acts[rpID] = int8(tmpAct[strconv.FormatInt(rpID, 10)])
	}
	return
}

// SpamReply return spam of add reply
func (dao *RedisDao) SpamReply(c context.Context, mid int64) (rec, daily int, err error) {
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
		return
	}
	rec = ii[0]
	daily = ii[1]
	return
}

// SpamAction return spam of add action
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

// NotifyCnt return notify max count.
func (dao *RedisDao) NotifyCnt(c context.Context, oid int64, typ int8) (cnt int, err error) {
	key := keyNotifyCnt(oid, typ)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if cnt, err = redis.Int(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET %s) error(%v)", key, err)
		}
	}
	return
}

// SetNotifyCnt set notify max count.
func (dao *RedisDao) SetNotifyCnt(c context.Context, oid int64, typ int8, cnt int) (err error) {
	key := keyNotifyCnt(oid, typ)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("SET", key, cnt); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, dao.expireNotify); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// MaxLikeCnt return reply max like count.
func (dao *RedisDao) MaxLikeCnt(c context.Context, rpid int64) (cnt int, err error) {
	key := keyMaxLikeCnt(rpid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if cnt, err = redis.Int(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("conn.Do(GET %s) error(%v)", key, err)
		}
	}
	return
}

// SetMaxLikeCnt set reply max like count.
func (dao *RedisDao) SetMaxLikeCnt(c context.Context, rpid, cnt int64) (err error) {
	key := keyMaxLikeCnt(rpid)
	conn := dao.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("SET", key, cnt); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Send("EXPIRE", key, dao.expireNotify); err != nil {
		log.Error("conn.Send error(%v)", err)
		return
	}
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < 2; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive error(%v)", err)
			return
		}
	}
	return
}

// ExpireFolder ...
func (dao *RedisDao) ExpireFolder(c context.Context, kind string, ID int64) (ok bool, err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	if ok, err = redis.Bool(conn.Do("EXPIRE", keyFolderIdx(kind, ID), dao.expireRdsIdx)); err != nil {
		log.Error("conn.Do(EXPIRE) error(%v)", err)
	}
	return
}

// AddFolder ...
func (dao *RedisDao) AddFolderBatch(c context.Context, kind string, ID int64, rps []*reply.Reply) (err error) {
	var (
		conn  = dao.redis.Get(c)
		count = 0
		key   = keyFolderIdx(kind, ID)
		args  []interface{}
	)
	defer conn.Close()
	args = append(args, key)
	for _, rp := range rps {
		args = append(args, rp.Floor)
		args = append(args, rp.RpID)
	}
	if err = conn.Send("ZADD", args...); err != nil {
		log.Error("conn.Send(ZADD %s) error(%v)", key, err)
		return
	}
	count++
	if err = conn.Send("EXPIRE", key, dao.expireRdsIdx); err != nil {
		return
	}
	count++
	if err = conn.Flush(); err != nil {
		log.Error("conn.Flush error(%v)", err)
		return
	}
	for i := 0; i < count; i++ {
		if _, err = conn.Receive(); err != nil {
			log.Error("conn.Receive() error(%v)", err)
			return
		}
	}
	return
}

// RemFolder ...
func (dao *RedisDao) RemFolder(c context.Context, kind string, ID, rpID int64) (err error) {
	conn := dao.redis.Get(c)
	defer conn.Close()
	key := keyFolderIdx(kind, ID)
	if _, err = conn.Do("ZREM", key, rpID); err != nil {
		log.Error("conn.Do(ZREM) error(%v)", err)
	}
	return
}
