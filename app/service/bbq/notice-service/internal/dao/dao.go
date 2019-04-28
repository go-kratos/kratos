package dao

import (
	"context"
	"fmt"

	"go-common/app/service/bbq/notice-service/api/v1"
	"go-common/app/service/bbq/notice-service/internal/conf"
	push "go-common/app/service/bbq/push/api/grpc/v1"
	"go-common/library/cache/redis"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
)

const (
	_listSQL         = "select id, mid, action_mid, svid, notice_type, title, text, jump_url, biz_type, biz_id, ctime from notice_%02d where mid = ? and notice_type = ? and id < ? order by id desc limit %d"
	_insertSQL       = "insert into notice_%02d (mid, action_mid, svid, notice_type, title, text, jump_url, biz_type, biz_id) values (?,?,?,?,?,?,?,?,?)"
	_noticeLen       = 10
	_redisUnreadKey  = "notice:unread:%d"
	_redisExpireTime = 7776000 // 90days
)

// Dao dao
type Dao struct {
	c          *conf.Config
	db         *xsql.DB
	redis      *redis.Pool
	pushClient push.PushClient
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:          c,
		db:         xsql.NewMySQL(c.MySQL),
		redis:      redis.NewPool(c.Redis),
		pushClient: newPushClient(c.GRPCClient["push"]),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return d.db.Ping(ctx)
}

func getTableIndex(id int64) int64 {
	return id % 100
}

// newPushClient .
func newPushClient(cfg *conf.GRPCClientConfig) push.PushClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return push.NewPushClient(cc)
}

// ListNotices 获取通知列表
func (d *Dao) ListNotices(ctx context.Context, mid, cursorID int64, noticeType int32) (list []*v1.NoticeBase, err error) {
	querySQL := fmt.Sprintf(_listSQL, getTableIndex(mid), _noticeLen)
	log.V(1).Infov(ctx, log.KV("mid", mid), log.KV("mid", mid), log.KV("notice_type", noticeType), log.KV("cursor_id", cursorID), log.KV("sql", querySQL))
	rows, err := d.db.Query(ctx, querySQL, mid, noticeType, cursorID)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "query mysql notice list fail"), log.KV("sql", querySQL), log.KV("mid", mid), log.KV("biz_type", noticeType), log.KV("cursor_id", cursorID))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var notice v1.NoticeBase
		if err = rows.Scan(&notice.Id, &notice.Mid, &notice.ActionMid, &notice.SvId, &notice.NoticeType, &notice.Title, &notice.Text, &notice.JumpUrl, &notice.BizType, &notice.BizId, &notice.NoticeTime); err != nil {
			log.Errorv(ctx, log.KV("log", "scan mysql notice list fail"), log.KV("sql", querySQL), log.KV("mid", mid), log.KV("biz_type", noticeType), log.KV("mid", mid), log.KV("cursor_id", cursorID))
			return
		}
		list = append(list, &notice)
	}

	// 只要用户读取数据，即清理未读数
	conn := d.redis.Get(ctx)
	defer conn.Close()
	redisKey := fmt.Sprintf(_redisUnreadKey, mid)
	if _, tmpErr := conn.Do("HSET", redisKey, noticeType, 0); tmpErr != nil {
		log.Warnv(ctx, log.KV("log", "clear unread info redis fail: key="+redisKey))
	}

	log.V(1).Infov(ctx, log.KV("req_size", _noticeLen), log.KV("rsp_size", len(list)))
	return
}

// CreateNotice 创建通知
func (d *Dao) CreateNotice(ctx context.Context, notice *v1.NoticeBase) (id int64, err error) {
	querySQL := fmt.Sprintf(_insertSQL, getTableIndex(notice.Mid))
	res, err := d.db.Exec(ctx, querySQL, notice.Mid, notice.ActionMid, notice.SvId, notice.NoticeType, notice.Title, notice.Text, notice.JumpUrl, notice.BizType, notice.BizId)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "exec mysql fail: create notice"), log.KV("sql", querySQL))
		return
	}
	id, _ = res.LastInsertId()

	return
}

// IncreaseUnread 增加未读
func (d *Dao) IncreaseUnread(ctx context.Context, mid int64, noticeType int32, num int64) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	redisKey := fmt.Sprintf(_redisUnreadKey, mid)
	expireResult, _ := redis.Int(conn.Do("EXPIRE", redisKey, _redisExpireTime))
	if expireResult == 0 {
		log.Infov(ctx, log.KV("log", "expire fail: key="+redisKey))
	}

	_, err = conn.Do("HINCRBY", redisKey, noticeType, num)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "HINCRBY notice unread fail: err="+err.Error()))
		return
	}
	log.V(1).Infov(ctx, log.KV("log", "hincrby notice unread : key="+redisKey), log.KV("notice_type", noticeType), log.KV("num", num))
	return
}

// ClearUnread 清理未读
func (d *Dao) ClearUnread(ctx context.Context, mid int64, noticeType int32) (err error) {
	conn := d.redis.Get(ctx)
	defer conn.Close()

	redisKey := fmt.Sprintf(_redisUnreadKey, mid)
	expireResult, _ := redis.Int(conn.Do("EXPIRE", redisKey, _redisExpireTime))
	if expireResult == 0 {
		log.Infov(ctx, log.KV("log", "expire fail and return: key="+redisKey))
		return
	}

	_, err = conn.Do("HSET", redisKey, noticeType, 0)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "HSET notice unread fail: err="+err.Error()))
		return
	}
	log.V(1).Infov(ctx, log.KV("log", "HSET clear notice unread : key="+redisKey), log.KV("notice_type", noticeType))

	// 清理推送用户
	err = d.ClearPushActionMid(ctx, mid, noticeType)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "ClearPushActionMid fail: err="+err.Error()))
		return
	}

	return
}

// GetUnreadInfo 获取未读情况
func (d *Dao) GetUnreadInfo(ctx context.Context, mid int64) (list []*v1.UnreadItem, err error) {
	redisKey := fmt.Sprintf(_redisUnreadKey, mid)
	conn := d.redis.Get(ctx)
	defer conn.Close()
	expireResult, _ := redis.Int(conn.Do("EXPIRE", redisKey, _redisExpireTime))
	if expireResult == 0 {
		log.V(1).Infov(ctx, log.KV("log", "expire fail: key="+redisKey))
		return
	}

	result, err := redis.Int64s(conn.Do("HMGET", redisKey, 1, 2, 3, 4))
	if err != nil {
		log.Errorv(ctx, log.KV("log", "hmget notice unread fail: err="+err.Error()))
		return
	}
	for i, val := range result {
		var item v1.UnreadItem
		item.NoticeType = int32(i + 1)
		item.UnreadNum = val
		list = append(list, &item)
	}
	return
}
