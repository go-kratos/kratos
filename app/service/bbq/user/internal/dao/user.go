package dao

import (
	"context"
	"fmt"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/model"
	accountv1 "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	_userLikeNum   = 128
	_userFollowNum = 128
	_userFanNum    = 128
)

//calcTableID .
func (d *Dao) calcTableID(num, mid int64) string {
	id := mid % 100
	return fmt.Sprintf("%02d", id)
}

func (d *Dao) getTableIndex(mid int64) int64 {
	return mid % 100
}

//userLikeSQL .
func (d *Dao) userLikeSQL(mid int64, sql string) string {
	tableName := "user_like_" + d.calcTableID(_userLikeNum, mid)
	return fmt.Sprintf(sql, tableName)
}

//userFollowSQL .
func (d *Dao) userFollowSQL(mid int64, sql string) string {
	tableName := "user_follow_" + d.calcTableID(_userFollowNum, mid)
	return fmt.Sprintf(sql, tableName)
}

//userFanSQL .
func (d *Dao) userFanSQL(mid int64, sql string) string {
	tableName := "user_fan_" + d.calcTableID(_userFanNum, mid)
	return fmt.Sprintf(sql, tableName)
}

// isMidIn 获取mid的关注up主
// 如果key在map中，那么value值肯定为1
func (d *Dao) isMidIn(c context.Context, mid int64, candidateMIDs []int64, sql string) (MIDMap map[int64]bool) {
	if len(candidateMIDs) == 0 {
		return
	}
	MIDMap = make(map[int64]bool)
	tableName := d.getTableIndex(mid)
	midstr := xstr.JoinInts(candidateMIDs)
	querySQL := fmt.Sprintf(sql, tableName, midstr)
	log.V(1).Infov(c, log.KV("event", "fetch_list"), log.KV("sql", querySQL))
	rows, err := d.db.Query(c, querySQL, mid)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_select"), log.KV("sql", querySQL))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var followedMid int64
		if err = rows.Scan(&followedMid); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("sql", querySQL))
			continue
		}
		MIDMap[followedMid] = true
	}
	log.Infov(c, log.KV("event", "mysql_select"), log.KV("sql", querySQL), log.KV("mid", mid),
		log.KV("req_size", len(candidateMIDs)), log.KV("rsp_size", len(MIDMap)))
	return MIDMap
}

// fetchPartRelationUserList 获取相关的用户列表，可以是关注列表也可以是粉丝列表，根据sql区别
func (d *Dao) fetchPartRelationUserList(c context.Context, mid int64, cursor model.CursorValue, sql string) (
	MID2IDMap map[int64]time.Time, relationMIDs []int64, err error) {

	MID2IDMap = make(map[int64]time.Time)

	querySQL := fmt.Sprintf(sql, d.getTableIndex(mid), model.UserListLen)
	log.V(1).Infov(c, log.KV("event", "fetch_follow_list"), log.KV("sql", querySQL))
	rows, err := d.db.Query(c, querySQL, mid, cursor.CursorTime)
	if err != nil {
		log.Errorv(c, log.KV("event", "mysql_select"), log.KV("sql", querySQL))
		return
	}
	defer rows.Close()
	conflict := bool(true)
	for rows.Next() {
		var relationMID int64
		var mtime time.Time
		if err = rows.Scan(&relationMID, &mtime); err != nil {
			log.Errorv(c, log.KV("event", "mysql_scan"), log.KV("sql", querySQL))
			return
		}
		// 为了解决同一个mtime的冲突问题
		if mtime == cursor.CursorTime && conflict {
			if relationMID == cursor.CursorID {
				conflict = false
			}
			continue
		}
		relationMIDs = append(relationMIDs, relationMID)
		MID2IDMap[relationMID] = mtime
	}
	log.Infov(c, log.KV("event", "mysql_select"), log.KV("sql", querySQL),
		log.KV("mid", mid), log.KV("relation_num", len(relationMIDs)))
	return
}

//GetUserBProfile 获取用户全量b站信息
func (d *Dao) GetUserBProfile(c context.Context, in *api.PhoneCheckReq) (res *accountv1.ProfileReply, err error) {
	req := &accountv1.MidReq{
		Mid:    in.Mid,
		RealIp: "",
	}
	res, err = d.accountClient.Profile3(c, req)
	return
}
