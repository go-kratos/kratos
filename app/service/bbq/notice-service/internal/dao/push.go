package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/bbq/notice-service/api/v1"
	push "go-common/app/service/bbq/push/api/grpc/v1"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_queryUserPushDev = "select `reg_id`, `sdk`, `platform` from `user_push_device` where `mid` = ? and `state` = 0 order by `last_login_time` desc limit 1;"
	_queryUserName    = "select `uname` from `user_base` where mid in %s;"
)

// PushNotice .
func (d *Dao) PushNotice(c context.Context, req *push.NotificationRequest) (result []*push.PushResult, err error) {
	resp, err := d.pushClient.AsyncNotification(c, req)
	if err != nil {
		log.Errorv(c, log.KV("log", "push notification error"), log.KV("error", err))
		return
	}
	result = resp.Result
	for _, item := range result {
		if item.Error != nil {
			log.Errorv(c, log.KV("log", "push device notifiaction error"), log.KV("error", item.Error))
		}
	}
	return
}

// PushMessage .
func (d *Dao) PushMessage(c context.Context, req *push.MessageRequest) (result []*push.PushResult, err error) {
	resp, err := d.pushClient.AsyncMessage(c, req)
	if err != nil {
		log.Errorv(c, log.KV("log", "push message error"), log.KV("error", err))
		return
	}
	result = resp.Result
	for _, item := range result {
		if item.Error != nil {
			log.Errorv(c, log.KV("log", "push device message error"), log.KV("error", item.Error))
		}
	}
	return nil, err
}

// FetchPushDev .
func (d *Dao) FetchPushDev(c context.Context, mid int64) (result *push.Device, err error) {
	result = &push.Device{}
	err = d.db.QueryRow(c, _queryUserPushDev, mid).Scan(&result.RegisterID, &result.SDK, &result.Platform)
	if err == sql.ErrNoRows {
		err = nil
		log.Warnw(c, "log", "no row in push", "sql", _queryUserPushDev, "mid", mid)
	}
	return
}

// IncrDailyPushCount .
func (d *Dao) IncrDailyPushCount(c context.Context, mid int64) (count int, err error) {
	dt := time.Now().Format("20060102")
	key := fmt.Sprintf("bbq:push:u:%d:%s:count", mid, dt)

	conn := d.redis.Get(c)
	defer conn.Close()

	count, err = redis.Int(conn.Do("INCR", key, 1))
	if err != nil {
		return
	}
	_, err = conn.Do("EXPIRE", 86400)

	return
}

// IncrHourPushAction .
func (d *Dao) IncrHourPushAction(c context.Context, mid int64, noticeType int32, t int) (count int, err error) {
	dt := time.Now().Format("2006010215")
	key := fmt.Sprintf("bbq:push:u:%d:%s:action:%d", mid, dt, noticeType)

	conn := d.redis.Get(c)
	defer conn.Close()

	count, err = redis.Int(conn.Do("INCR", key))
	if err != nil {
		return
	}

	_, err = conn.Do("EXPIRE", t)
	return
}

// ClearHourPushAction .
func (d *Dao) ClearHourPushAction(c context.Context, mid int64, noticeType int32) error {
	dt := time.Now().Format("2006010215")
	key := fmt.Sprintf("bbq:push:u:%d:%s:action:%d", mid, dt, noticeType)

	conn := d.redis.Get(c)
	defer conn.Close()

	_, err := conn.Do("DEL", key)

	return err
}

// SetPushActionMid .
func (d *Dao) SetPushActionMid(c context.Context, mid int64, actionMid int64, noticeType int32) error {
	dt := time.Now().Format("2006010215")
	key := fmt.Sprintf("bbq:v1:push:u:%d:%s:action:%d", mid, dt, noticeType)

	conn := d.redis.Get(c)
	defer conn.Close()

	values, _ := redis.Values(conn.Do("HGETALL", key))
	if len(values) >= 4 {
		return nil
	}

	_, err := conn.Do("HMSET", key, actionMid, actionMid)

	return err
}

// GetPushActionMid .
func (d *Dao) GetPushActionMid(c context.Context, mid int64, noticeType int32) ([]int64, error) {
	dt := time.Now().Format("2006010215")
	key := fmt.Sprintf("bbq:v1:push:u:%d:%s:action:%d", mid, dt, noticeType)

	conn := d.redis.Get(c)
	defer conn.Close()

	m, err := redis.Int64Map(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	result := make([]int64, 0)
	for _, v := range m {
		result = append(result, v)
	}

	return result, nil
}

// ClearPushActionMid .
func (d *Dao) ClearPushActionMid(c context.Context, mid int64, noticeType int32) error {
	conn := d.redis.Get(c)
	defer conn.Close()

	dt := time.Now().Format("2006010215")
	key := fmt.Sprintf("bbq:v1:push:u:%d:%s:action:%d", mid, dt, noticeType)

	_, err := conn.Do("DEL", key)
	return err
}

// GetUserName .
func (d *Dao) GetUserName(c context.Context, midList []int64, count int) (names []string, err error) {
	mids := []string{}
	list := midList
	if len(midList) > count {
		list = midList[:count]
	}
	for _, v := range list {
		mids = append(mids, strconv.Itoa(int(v)))
	}
	where := "(" + strings.Join(mids, ",") + ")"
	row, err := d.db.Query(c, fmt.Sprintf(_queryUserName, where))
	if err != nil {
		return
	}
	for row.Next() {
		n := ""
		err = row.Scan(&n)
		if err != nil {
			return
		}
		names = append(names, n)
	}
	return
}

// FetchUserPushDev .
func (d *Dao) FetchUserPushDev(c context.Context, mid int64, buvid string) (result *v1.UserPushDev, err error) {
	querySQL := "select `id` from `user_push_device` where `mid` = ? and `buvid` = ?"
	row := d.db.QueryRow(c, querySQL, mid, buvid)
	result = &v1.UserPushDev{}
	row.Scan(&result.Id)
	return
}

// InsertUserPushDev .
func (d *Dao) InsertUserPushDev(c context.Context, req *v1.UserPushDev) (int64, error) {
	_insertUserPushDev := "insert into `user_push_device` (`mid`, `reg_id`, `buvid`, `sdk`, `platform`) values (?, ?, ?, ?, ?);"
	result, err := d.db.Exec(c, _insertUserPushDev, req.Mid, req.RegisterId, req.Buvid, req.Sdk, req.Platform)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateUserPushDev .
func (d *Dao) UpdateUserPushDev(c context.Context, req *v1.UserPushDev) (int64, error) {
	_updateUserPushDev := "update `user_push_device` set `reg_id`=?,`sdk`=?,`platform`=?,`state`=?,`last_login_time`=now() where `mid`=? and `buvid`=?"
	result, err := d.db.Exec(c, _updateUserPushDev, req.RegisterId, req.Sdk, req.Platform, req.State, req.Mid, req.Buvid)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// DeleteUserPushDev .
func (d *Dao) DeleteUserPushDev(c context.Context, req *v1.UserPushDev) (int64, error) {
	_updateUserPushDev := "update `user_push_device` set `state`=1 where `mid`=? and `buvid`=?"
	result, err := d.db.Exec(c, _updateUserPushDev, req.Mid, req.Buvid)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
