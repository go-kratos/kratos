package newcomer

import (
	"context"
	"fmt"

	"database/sql"
	"go-common/app/interface/main/creative/model/newcomer"
	"go-common/library/log"
	"go-common/library/xstr"
	"strings"
)

const (
	// select
	_getUserTaskBindSQL        = "SELECT count(id) FROM newcomers_task_user_%s WHERE mid=?"
	_getGiftRewardSQL          = "SELECT reward_id FROM newcomers_gift_reward WHERE state=0 AND task_type=?"
	_getTaskGroupRewardSQL     = "SELECT reward_id FROM newcomers_grouptask_reward WHERE state=0 AND task_group_id=?"
	_getIsReceivedSQL          = "SELECT count(id) FROM newcomers_reward_receive WHERE mid=?"
	_getRewardReceiveStateSQL  = "SELECT task_group_id FROM newcomers_reward_receive WHERE mid=?"
	_getRewardReceiveSQL       = "SELECT id, mid, task_gift_id, task_group_id, reward_id, reward_type, state, receive_time, ctime, mtime FROM newcomers_reward_receive WHERE mid=? AND state IN (0,1,2) ORDER BY receive_time ASC"
	_getRewardSQL              = "SELECT id, parent_id, type, state, is_active, prize_id, prize_unit, expire, name, logo, comment, unlock_logo, name_extra, ctime, mtime FROM newcomers_reward  WHERE state=0"
	_getUserTaskByMIDSQL       = "SELECT id, mid, task_id, task_group_id, task_type, state, task_bind_time, ctime, mtime FROM newcomers_task_user_%s WHERE state=-1 AND mid=?"
	_getTaskSQL                = "SELECT id, group_id, type, state, target_type, target_value, title, `desc`, comment, ctime, mtime, rank, extra, fan_range, up_time, down_time FROM newcomers_task WHERE state IN (0,1) ORDER BY rank ASC"
	_getAllTaskGroupRewardSQL  = "SELECT id,task_group_id,reward_id,state,comment,ctime,mtime FROM newcomers_grouptask_reward WHERE state=0"
	_getAllGiftRewardSQL       = "SELECT id,root_type,task_type,reward_id,state,comment,ctime,mtime FROM newcomers_gift_reward WHERE state=0"
	_getUserTaskSQL            = "SELECT task_id,state FROM newcomers_task_user_%s WHERE mid=?"
	_getRewardReceiveByIDSQL   = "SELECT id, mid, task_gift_id, task_group_id, reward_id, reward_type, state, receive_time, ctime, mtime FROM newcomers_reward_receive WHERE mid=? AND id=?"
	_getTaskGroupSQL           = "SELECT id, rank, state, root_type, type, ctime, mtime FROM newcomers_task_group WHERE state=0"
	_getTaskRewardSQL          = "SELECT id,task_id,reward_id,state,comment,ctime,mtime FROM newcomers_task_reward WHERE state=0"
	_getRewardReceiveByOldInfo = "SELECT id, mid, oid, type, reward_id, reward_type, state, ctime, mtime FROM newcomers_reward_receive_%s WHERE mid=? AND oid=? AND type=? AND reward_id=?"
	// insert
	// _inRewardReceiveSQL = "INSERT INTO newcomers_reward_receive (mid, task_gift_id, task_group_id, reward_id, reward_type, state, receive_time, ctime, mtime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	// update
	_upRewardReceiveSQL  = "UPDATE newcomers_reward_receive SET state=? WHERE mid=? AND id=?"
	_upRewardReceiveSQL2 = "UPDATE newcomers_reward_receive_%s SET state=? WHERE mid=? AND id=?"
	// update
	_upUserTaskSQL = "UPDATE newcomers_task_user_%s SET state=? WHERE mid=? AND task_id=?"
)

// getTableName by mid%100
func getTableName(mid int64) string {
	return fmt.Sprintf("%02d", mid%100)
}

// UserTaskBind determine if the user has bound the task
func (d *Dao) UserTaskBind(c context.Context, mid int64) (res int64, err error) {
	row := d.db.QueryRow(c, fmt.Sprintf(_getUserTaskBindSQL, getTableName(mid)), mid)
	if err = row.Scan(&res); err != nil {
		log.Error("UserTaskBind d.db.QueryRow error(%v)", err)
		return
	}
	return
}

// IsRewardReceived if the basic reward and the gift award have been received
func (d *Dao) IsRewardReceived(c context.Context, mid int64, rid int64, rewardType int8) (res bool, err error) {
	sqlStr := _getIsReceivedSQL
	if rewardType == newcomer.RewardGiftType {
		sqlStr += fmt.Sprintf(" AND task_gift_id=?")
	} else {
		sqlStr += fmt.Sprintf(" AND task_group_id=?")
	}
	row := d.db.QueryRow(c, sqlStr, mid, rid)
	count := 0
	if err = row.Scan(&count); err != nil {
		log.Error("isRewardReceived d.db.Query error(%v)", err)
		return
	}
	if count > 0 {
		return true, nil
	}
	return
}

// RewardReceivedGroup get []task_group_id by []group_id
func (d *Dao) RewardReceivedGroup(c context.Context, mid int64, ids []int64) (res []int, err error) {
	if len(ids) == 0 {
		return
	}
	sqlStr := _getRewardReceiveStateSQL
	sqlStr += fmt.Sprintf(" AND task_group_id IN (%s)", xstr.JoinInts(ids))
	rows, err := d.db.Query(c, sqlStr, mid)
	if err != nil {
		log.Error("RewardReceivedGroup d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make([]int, 0)
	for rows.Next() {
		gid := 0
		if err = rows.Scan(&gid); err != nil {
			log.Error("RewardReceivedGroup rows.Scan error(%v)", err)
			return
		}
		res = append(res, gid)
	}

	return
}

// GiftRewards get the rewards of gift by task_type
func (d *Dao) GiftRewards(c context.Context, taskType int8) (res []int64, err error) {
	rows, err := d.db.Query(c, _getGiftRewardSQL, taskType)
	if err != nil {
		log.Error("GiftRewards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make([]int64, 0)
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("GiftRewards rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}

	return
}

// TaskGroupRewards get the rewards of gift by group_id
func (d *Dao) TaskGroupRewards(c context.Context, groupID int64) (res []int64, err error) {
	rows, err := d.db.Query(c, _getTaskGroupRewardSQL, groupID)
	if err != nil {
		log.Error("TaskGroupRewards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make([]int64, 0)
	for rows.Next() {
		var r int64
		if err = rows.Scan(&r); err != nil {
			log.Error("TaskGroupRewards rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}

	return
}

// RewardCompleteState judge the completion of the task
func (d *Dao) RewardCompleteState(c context.Context, mid int64, tids []int64) (res int64, err error) {
	sqlStr := fmt.Sprintf(_getUserTaskSQL, getTableName(mid))
	sqlStr += fmt.Sprintf(" AND task_id IN (%s)", xstr.JoinInts(tids))
	rows, err := d.db.Query(c, sqlStr, mid)
	if err != nil {
		log.Error("RewardCompleteState d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = 0
	for rows.Next() {
		t := &newcomer.Task{}
		if err = rows.Scan(&t.ID, &t.CompleteSate); err != nil {
			log.Error("RewardCompleteState rows.Scan error(%v)", err)
			return
		}
		if t.CompleteSate == newcomer.TaskIncomplete {
			res++
		}
	}
	return
}

// RewardReceive  add reward receive records.
func (d *Dao) RewardReceive(c context.Context, places string, args []interface{}) (int64, error) {
	var (
		res sql.Result
		err error
	)

	sqlStr := fmt.Sprintf("INSERT INTO newcomers_reward_receive(mid, task_gift_id, task_group_id, reward_id, reward_type, state) VALUES %s", places)
	res, err = d.db.Exec(c, sqlStr, args...)
	if err != nil {
		log.Error("RewardReceive tx.Exec error(%v)", err)
		return 0, err
	}

	return res.LastInsertId()
}

// RewardActivate activate reward.
func (d *Dao) RewardActivate(c context.Context, mid, id int64) (int64, error) {
	res, err := d.db.Exec(c, _upRewardReceiveSQL, newcomer.RewardActivatedNotClick, mid, id)
	if err != nil {
		log.Error("RewardActivate d.db.Exec mid(%d) id(%d) error(%v)", mid, id, err)
		return 0, err
	}
	return res.RowsAffected()
}

//RewardReceives get reward receive by mid.
func (d *Dao) RewardReceives(c context.Context, mid int64) (res map[int8][]*newcomer.RewardReceive, err error) {
	rows, err := d.db.Query(c, _getRewardReceiveSQL, mid)
	if err != nil {
		log.Error("RewardReceives d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make(map[int8][]*newcomer.RewardReceive)
	for rows.Next() {
		r := &newcomer.RewardReceive{}
		if err = rows.Scan(&r.ID, &r.MID, &r.TaskGiftID, &r.TaskGroupID, &r.RewardID, &r.RewardType, &r.State, &r.ReceiveTime, &r.CTime, &r.MTime); err != nil {
			log.Error("RewardReceives rows.Scan error(%v)", err)
			return
		}
		res[r.RewardType] = append(res[r.RewardType], r)
	}
	return
}

//Rewards get all reward receive.
func (d *Dao) Rewards(c context.Context) (res []*newcomer.Reward, err error) {
	rows, err := d.db.Query(c, _getRewardSQL)
	if err != nil {
		log.Error("Rewards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*newcomer.Reward, 0)
	for rows.Next() {
		r := &newcomer.Reward{}
		if err = rows.Scan(&r.ID, &r.ParentID, &r.Type, &r.State, &r.IsActive, &r.PriceID, &r.PrizeUnit, &r.Expire, &r.Name, &r.Logo, &r.Comment, &r.UnlockLogo, &r.NameExtra, &r.CTime, &r.MTime); err != nil {
			log.Error("Rewards rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// BindTasks user bind tasks
func (d *Dao) BindTasks(c context.Context, mid int64, places string, args []interface{}) (res int64, err error) {
	sqlStr := fmt.Sprintf("INSERT INTO newcomers_task_user_%s(mid , task_id , task_group_id , task_type , state) VALUES %s", getTableName(mid), places)
	result, err := d.db.Exec(c, sqlStr, args...)
	if err != nil {
		log.Error("BindTasks d.db.Exec error(%v)", err)
		return
	}
	return result.LastInsertId()
}

// GiftRewardCount get received gift reward count
func (d *Dao) GiftRewardCount(c context.Context, mid int64, ids []int64) (res int64, err error) {
	if len(ids) == 0 {
		return
	}
	sqlStr := fmt.Sprintf("SELECT count(DISTINCT task_gift_id) FROM newcomers_reward_receive WHERE mid=? AND task_gift_id IN (%s)", xstr.JoinInts(ids))
	row := d.db.QueryRow(c, sqlStr, mid)
	err = row.Scan(&res)
	if err != nil {
		log.Error("GiftRewardCount d.db.QueryRow error(%v)", err)
		return
	}
	return
}

// BaseRewardCount get received base reward count
func (d *Dao) BaseRewardCount(c context.Context, mid int64, ids []int64) (res int64, err error) {
	if len(ids) == 0 {
		return
	}
	sqlStr := fmt.Sprintf("SELECT count(DISTINCT task_group_id) FROM newcomers_reward_receive WHERE mid=? AND task_group_id IN (%s)", xstr.JoinInts(ids))
	row := d.db.QueryRow(c, sqlStr, mid)
	err = row.Scan(&res)
	if err != nil {
		log.Error("BaseRewardCount d.db.QueryRow error(%v)", err)
		return
	}
	return
}

//Tasks get all task.
func (d *Dao) Tasks(c context.Context) (res []*newcomer.Task, err error) {
	rows, err := d.db.Query(c, _getTaskSQL)
	if err != nil {
		log.Error("Tasks d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*newcomer.Task, 0)

	for rows.Next() {
		r := &newcomer.Task{}
		if err = rows.Scan(&r.ID, &r.GroupID, &r.Type, &r.State, &r.TargetType, &r.TargetValue, &r.Title, &r.Desc, &r.Comment, &r.CTime, &r.MTime, &r.Rank, &r.Extra, &r.FanRange, &r.UpTime, &r.DownTime); err != nil {
			log.Error("Tasks rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

//UserTasksByMID get user unfinish task by mid.
func (d *Dao) UserTasksByMID(c context.Context, mid int64) (res []*newcomer.UserTask, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getUserTaskByMIDSQL, getTableName(mid)), mid)
	if err != nil {
		log.Error("UserTasksByMID d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*newcomer.UserTask, 0)
	for rows.Next() {
		r := &newcomer.UserTask{}
		if err = rows.Scan(&r.ID, &r.MID, &r.TaskID, &r.TaskGroupID, &r.TaskType, &r.State, &r.TaskBindTime, &r.CTime, &r.MTime); err != nil {
			log.Error("UserTasksByMID rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// UpUserTask update user task finish state
func (d *Dao) UpUserTask(c context.Context, mid, tid int64) (int64, error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_upUserTaskSQL, getTableName(mid)), 0, mid, tid)
	if err != nil {
		log.Error("RewardActivate d.db.Exec mid(%d) id(%d) error(%v)", mid, tid, err)
		return 0, err
	}
	return res.RowsAffected()
}

// UserTaskType Determining the type of task that users have
func (d *Dao) UserTaskType(c context.Context, mid int64) (res int8, err error) {
	sqlStr := fmt.Sprintf("SELECT count(DISTINCT task_type) FROM newcomers_task_user_%s WHERE mid=?", getTableName(mid))
	row := d.db.QueryRow(c, sqlStr, mid)
	if err = row.Scan(&res); err != nil {
		log.Error("UserTaskType row.Scan(&res) mid(%d)|error(%v)", mid, err)
		return
	}
	return
}

// AllTaskGroupRewards get all TaskGroupRewards for cache
func (d *Dao) AllTaskGroupRewards(c context.Context) (res map[int64][]*newcomer.TaskGroupReward, err error) {
	rows, err := d.db.Query(c, _getAllTaskGroupRewardSQL)
	if err != nil {
		log.Error("AllTaskGroupRewards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64][]*newcomer.TaskGroupReward)
	for rows.Next() {
		t := &newcomer.TaskGroupReward{}
		if err = rows.Scan(&t.ID, &t.TaskGroupID, &t.RewardID, &t.State, &t.Comment, &t.CTime, &t.MTime); err != nil {
			log.Error("AllTaskGroupRewards rows.Scan error(%v)", err)
			return
		}
		res[t.TaskGroupID] = append(res[t.TaskGroupID], t)
	}
	return
}

// AllGiftRewards get all GiftRewards for cache
func (d *Dao) AllGiftRewards(c context.Context) (res map[int8][]*newcomer.GiftReward, err error) {
	rows, err := d.db.Query(c, _getAllGiftRewardSQL)
	if err != nil {
		log.Error("AllTaskGroupRewards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int8][]*newcomer.GiftReward)
	for rows.Next() {
		t := &newcomer.GiftReward{}
		if err = rows.Scan(&t.ID, &t.RootType, &t.TaskType, &t.RewardID, &t.State, &t.Comment, &t.CTime, &t.MTime); err != nil {
			log.Error("AllTaskGroupRewards rows.Scan error(%v)", err)
			return
		}
		res[t.TaskType] = append(res[t.TaskType], t)
	}
	return
}

// UserTasks get user tasks
func (d *Dao) UserTasks(c context.Context, mid int64) (res []*newcomer.Task, err error) {
	sqlStr := fmt.Sprintf(_getUserTaskSQL, getTableName(mid))
	rows, err := d.db.Query(c, sqlStr, mid)
	if err != nil {
		log.Error("UserTasks d.db.Query mid(%d)|error(%v)", mid, err)
		return
	}
	defer rows.Close()

	res = make([]*newcomer.Task, 0)
	for rows.Next() {
		t := &newcomer.Task{}
		if err = rows.Scan(&t.ID, &t.CompleteSate); err != nil {
			log.Error("UserTasks rows.Scan mid(%d)|error(%v)", mid, err)
			return
		}
		res = append(res, t)
	}
	return
}

// RewardReceiveByID get rewardReceive by receiveID
func (d *Dao) RewardReceiveByID(c context.Context, mid, receiveID int64) (res *newcomer.RewardReceive, err error) {
	row := d.db.QueryRow(c, _getRewardReceiveByIDSQL, mid, receiveID)
	res = &newcomer.RewardReceive{}
	if err = row.Scan(&res.ID, &res.MID, &res.TaskGiftID, &res.TaskGroupID, &res.RewardID, &res.RewardType, &res.State, &res.ReceiveTime, &res.CTime, &res.MTime); err != nil {
		log.Error("RewardIDByReceiveID row.Scan error(%v)", err)
		return
	}
	return
}

// TaskGroups get task-group
func (d *Dao) TaskGroups(c context.Context) (res []*newcomer.TaskGroupEntity, err error) {
	rows, err := d.db.Query(c, _getTaskGroupSQL)
	if err != nil {
		log.Error("TaskGroups d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make([]*newcomer.TaskGroupEntity, 0)
	for rows.Next() {
		t := &newcomer.TaskGroupEntity{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.State, &t.RootType, &t.Type, &t.CTime, &t.MTime); err != nil {
			log.Error("TaskGroups rows.Scan error(%v)", err)
			return
		}
		res = append(res, t)
	}
	return
}

// TaskRewards get task-rewards
func (d *Dao) TaskRewards(c context.Context) (res map[int64][]*newcomer.TaskRewardEntity, err error) {
	rows, err := d.db.Query(c, _getTaskRewardSQL)
	if err != nil {
		log.Error("TaskRewards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make(map[int64][]*newcomer.TaskRewardEntity)
	for rows.Next() {
		t := &newcomer.TaskRewardEntity{}
		if err = rows.Scan(&t.ID, &t.TaskID, &t.RewardID, &t.State, &t.Comment, &t.CTime, &t.MTime); err != nil {
			log.Error("TaskRewards rows.Scan error(%v)", err)
			return
		}
		res[t.TaskID] = append(res[t.TaskID], t)
	}
	return
}

// RewardReceive2  add reward receive records.
func (d *Dao) RewardReceive2(c context.Context, mid int64, rrs []*newcomer.RewardReceive2) (res int64, err error) {
	if len(rrs) == 0 {
		return 0, nil
	}

	tx, err := d.db.Begin(c)
	if err != nil {
		log.Error("RewardReceive2 d.db.Begin mid(%d)|error(%v)", mid, err)
		return
	}

	place := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range rrs {
		place = append(place, "(?, ?, ?, ?, ? ,?)")
		if v.Type == newcomer.RewardBaseType {
			args = append(args, v.MID, 0, v.OID, v.RewardID, v.RewardType, v.State)
		} else if v.Type == newcomer.RewardGiftType {
			args = append(args, v.MID, v.OID, 0, v.RewardID, v.RewardType, v.State)
		}
	}
	placeStr := strings.Join(place, ",")
	sqlStr := fmt.Sprintf("INSERT INTO newcomers_reward_receive(mid, task_gift_id, task_group_id, reward_id, reward_type, state) VALUES %s", placeStr)
	_, err = tx.Exec(sqlStr, args...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error("RewardReceive2 tx.Rollback mid(%d)|error(%v)", mid, rollbackErr)
		}
		log.Error("RewardReceive2 tx.Exec mid(%d)|error(%v)", mid, err)
		return
	}

	// insert to new table
	placeNew := make([]string, 0)
	argsNew := make([]interface{}, 0)
	for _, v := range rrs {
		placeNew = append(placeNew, "(?, ?, ?, ?, ? ,?)")
		argsNew = append(argsNew, v.MID, v.OID, v.Type, v.RewardID, v.RewardType, v.State)
	}
	placeStrNew := strings.Join(placeNew, ",")
	sqlStrNew := fmt.Sprintf("INSERT INTO newcomers_reward_receive_%s(mid, oid, type, reward_id, reward_type, state) VALUES %s", getTableName(mid), placeStrNew)
	result, err := tx.Exec(sqlStrNew, argsNew...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error("RewardReceive2 tx.Rollback mid(%d)|error(%v)", mid, rollbackErr)
		}
		log.Error("RewardReceive2 tx.Exec mid(%d)|error(%v)", mid, err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("RewardReceive2 tx.Commit() mid(%d)|error(%v)", mid, err)
		return
	}

	return result.LastInsertId()
}

// RewardReceiveByOldInfo get rewardReceive2 by old info
func (d *Dao) RewardReceiveByOldInfo(c context.Context, r *newcomer.RewardReceive) (res *newcomer.RewardReceive2, err error) {
	sqlStr := fmt.Sprintf(_getRewardReceiveByOldInfo, getTableName(r.MID))
	var (
		oid int64
		ty  int8
	)
	if r.TaskGiftID == 0 && r.TaskGroupID != 0 {
		oid = r.TaskGroupID
		ty = newcomer.RewardBaseType
	} else if r.TaskGroupID == 0 && r.TaskGiftID != 0 {
		oid = r.TaskGiftID
		ty = newcomer.RewardGiftType
	}
	row := d.db.QueryRow(c, sqlStr, r.MID, oid, ty, r.RewardID)
	res = &newcomer.RewardReceive2{}
	if err = row.Scan(&res.ID, &res.MID, &res.OID, &res.Type, &res.RewardID, &res.RewardType, &res.State, &res.CTime, &res.MTime); err != nil {
		log.Error("RewardReceiveByOldInfo row.Scan error(%v)", err)
		return
	}
	return
}

// RewardActivate2 activate reward double write
func (d *Dao) RewardActivate2(c context.Context, mid, oid, nid int64) (int64, error) {
	tx, err := d.db.Begin(c)
	if err != nil {
		log.Error("RewardActivate2 d.db.Begin mid(%d)|error(%v)", mid, err)
		return 0, err
	}

	_, err = tx.Exec(_upRewardReceiveSQL, newcomer.RewardActivatedNotClick, mid, oid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error("RewardActivate2 tx.Rollback mid(%d)|error(%v)", mid, rollbackErr)
		}
		log.Error("RewardActivate2 tx.Exec mid(%d)|error(%v)", mid, err)
		return 0, err
	}

	// update new table
	sqlStr := fmt.Sprintf(_upRewardReceiveSQL2, getTableName(mid))
	res, err := tx.Exec(sqlStr, newcomer.RewardActivatedNotClick, mid, nid)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error("RewardActivate2 tx.Rollback mid(%d)|error(%v)", mid, rollbackErr)
		}
		log.Error("RewardActivate2 tx.Exec mid(%d)|error(%v)", mid, err)
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		log.Error("RewardActivate2 tx.Commit() mid(%d)|error(%v)", mid, err)
		return 0, err
	}

	return res.RowsAffected()
}
