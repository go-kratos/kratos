package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/workflow/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
)

const (
	_listKeyFormat        = "wf_%d_%d"
	_fbRound              = 11
	_auditFlowDealType    = 0
	_feedbackFlowDealType = 1

	_adminOnlineHash        = "wkf_online_hash"
	_prefixAdminOnlineField = "wkf_online_uid_%d"
)

func (d *Dao) pingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// redis list key wf_business_dealtype like wf_2_1
// dealtype=0 audit  dealtype=1 feedback

// RedisRPOPCids returns cids from a list
func (d *Dao) RedisRPOPCids(c context.Context, business int8, round int64, num int8) (cids []int64, err error) {
	var (
		key   string
		conn  = d.redis.Get(c)
		flow  int
		cid   int64
		chall *model.Chall
	)
	defer conn.Close()
	cids = make([]int64, 0)
	if round == _fbRound {
		flow = _feedbackFlowDealType
	} else {
		flow = _auditFlowDealType
	}

	key = fmt.Sprintf(_listKeyFormat, business, flow)
	for {
		exist := false
		if exist, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
			log.Error("redis.Bool key(%s) err(%v)", key, err)
			return
		}
		if !exist {
			log.Warn("key(%s) not exist", key)
			return
		}
		if cid, err = redis.Int64(conn.Do("RPOP", key)); err != nil {
			log.Error("conn.Do(RPOP,%s) error(%v)", key, err)
			return
		}
		// judge if business_state is queue state
		if chall, err = d.Chall(c, cid); err != nil {
			return
		}
		chall.FromState()
		if chall.BusinessState == model.QueueState {
			cids = append(cids, cid)
		}

		if len(cids) >= int(num) {
			break
		}
	}
	return
}

// IsOnline judge if admin is online
func (d *Dao) IsOnline(c context.Context, assigneeAdminID int64) (online bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	key := _adminOnlineHash
	field := d.fieldOnlineList(assigneeAdminID)
	return redis.Bool(conn.Do("HEXISTS", key, field))
}

// AddOnline checkin if start subscribe mission in platform, set key
func (d *Dao) AddOnline(c context.Context, assigneeAdminID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := _adminOnlineHash
	field := d.fieldOnlineList(assigneeAdminID)
	_, err = conn.Do("HSET", key, field, assigneeAdminID)
	return
}

// DelOnline checkout if exit subscribe mission in platform, delete key
func (d *Dao) DelOnline(c context.Context, assigneeAdminID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := _adminOnlineHash
	field := d.fieldOnlineList(assigneeAdminID)
	_, err = conn.Do("HDEL", key, field)
	return
}

// ListOnline list online admin
func (d *Dao) ListOnline(c context.Context) (ids []int64, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := _adminOnlineHash
	return redis.Int64s(conn.Do("HVALS", key))
}

// LogInOutTime show last online or offline time
func (d *Dao) LogInOutTime(c context.Context, uids []int64) {

}

// uid field in key wkf_online_hash
func (d *Dao) fieldOnlineList(assigneeAdminID int64) string {
	return fmt.Sprintf(_prefixAdminOnlineField, assigneeAdminID)
}

func (d *Dao) keyChallCount(assigneeAdminID int64) string {
	return fmt.Sprintf(_prefixChallPendingCount, assigneeAdminID)
}
