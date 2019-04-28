package newcomer

import (
	"context"
	"fmt"
	"time"

	"go-common/app/job/main/creative/model"
	"go-common/library/log"
)

//UserTasks get user unfinish  task.
func (d *Dao) UserTasks(c context.Context, index string, id int64, limit int) (res []*model.UserTask, err error) {
	_getUserTaskSQL := "SELECT id, mid, task_id, task_group_id, task_type, state, task_bind_time, ctime, mtime FROM newcomers_task_user_%s WHERE state=-1 AND id > ? order by id asc limit ?"

	rows, err := d.db.Query(c, fmt.Sprintf(_getUserTaskSQL, index), id, limit)
	if err != nil {
		log.Error("UserTasks d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.UserTask, 0)
	for rows.Next() {
		r := &model.UserTask{}
		if err = rows.Scan(&r.ID, &r.MID, &r.TaskID, &r.TaskGroupID, &r.TaskType, &r.State, &r.TaskBindTime, &r.CTime, &r.MTime); err != nil {
			log.Error("UserTasks rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// getTableName by mid%100
func getTableName(mid int64) string {
	return fmt.Sprintf("%02d", mid%100)
}

// UpUserTask update user task finish state
func (d *Dao) UpUserTask(c context.Context, mid, tid int64) (int64, error) {
	_upUserTaskSQL := "UPDATE newcomers_task_user_%s SET state=? WHERE mid=? AND task_id=?"

	res, err := d.db.Exec(c, fmt.Sprintf(_upUserTaskSQL, getTableName(mid)), 0, mid, tid)
	if err != nil {
		log.Error("RewardActivate d.db.Exec mid(%d) id(%d) error(%v)", mid, tid, err)
		return 0, err
	}
	return res.RowsAffected()
}

//UserTasksByMIDAndState get user unfinish task by mid & state.
func (d *Dao) UserTasksByMIDAndState(c context.Context, mid int64, state int) (res []*model.UserTask, err error) {
	_getUserTaskByMIDSQL := "SELECT id, mid, task_id, task_group_id, task_type, state, ctime, mtime FROM newcomers_task_user_%s WHERE mid=? AND state=?"

	rows, err := d.db.Query(c, fmt.Sprintf(_getUserTaskByMIDSQL, getTableName(mid)), mid, state)
	if err != nil {
		log.Error("UserTasksByMIDAndState d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.UserTask, 0)
	for rows.Next() {
		r := &model.UserTask{}
		if err = rows.Scan(&r.ID, &r.MID, &r.TaskID, &r.TaskGroupID, &r.TaskType, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("UserTasksByMIDAndState rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

//Tasks get all task.
func (d *Dao) Tasks(c context.Context) (res []*model.Task, err error) {
	_getTaskSQL := "SELECT id, group_id, type, state, target_type, target_value, title, `desc`, comment, ctime, mtime FROM newcomers_task WHERE state=0"

	rows, err := d.db.Query(c, _getTaskSQL)
	if err != nil {
		log.Error("Tasks d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make([]*model.Task, 0)
	for rows.Next() {
		r := &model.Task{}
		if err = rows.Scan(&r.ID, &r.GroupID, &r.Type, &r.State, &r.TargetType, &r.TargetValue, &r.Title, &r.Desc, &r.Comment, &r.CTime, &r.MTime); err != nil {
			log.Error("Tasks rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// TaskByTID get task by task_id
func (d *Dao) TaskByTID(c context.Context, mid int64, tid int64) (res *model.Task, err error) {
	_getOneTaskSQL := fmt.Sprintf("SELECT task_id,task_group_id,task_type,state FROM newcomers_task_user_%s WHERE mid=? AND task_id=?", getTableName(mid))

	row := d.db.QueryRow(c, _getOneTaskSQL, mid, tid)
	res = &model.Task{}
	if err = row.Scan(&res.ID, &res.GroupID, &res.Type, &res.State); err != nil {
		log.Error("TaskByTID ow.Scan error(%v)", err)
		return
	}

	return
}

// CheckTaskComplete check task complete state
func (d *Dao) CheckTaskComplete(c context.Context, mid int64, tid int64) bool {
	task, err := d.TaskByTID(c, mid, tid)
	if err != nil || task == nil {
		return false
	}
	if task.State == 0 {
		return true
	}
	return false
}

//UserTasksNotify get unfinish task send notify to user.
func (d *Dao) UserTasksNotify(c context.Context, index string, id int64, start, end string, limit int) (res []*model.UserTask, err error) {
	_getUserTaskSQL := "SELECT id, mid, task_id, task_group_id, task_type, state, ctime, mtime FROM newcomers_task_user_%s WHERE state=-1 AND task_type=1 AND id>? AND ctime >= ? AND ctime <= ? order by id asc limit ?"

	rows, err := d.db.Query(c, fmt.Sprintf(_getUserTaskSQL, index), id, start, end, limit)
	if err != nil {
		log.Error("UserTasks d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.UserTask, 0)
	for rows.Next() {
		r := &model.UserTask{}
		if err = rows.Scan(&r.ID, &r.MID, &r.TaskID, &r.TaskGroupID, &r.TaskType, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("UserTasks rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

//UserTasksByMID get user unfinish task by mid.
func (d *Dao) UserTasksByMID(c context.Context, mid int64) (res []*model.UserTask, err error) {
	_getUserTaskByMIDSQL := "SELECT id, mid, task_id, task_group_id, task_type, state, ctime, mtime FROM newcomers_task_user_%s WHERE mid=?"

	rows, err := d.db.Query(c, fmt.Sprintf(_getUserTaskByMIDSQL, getTableName(mid)), mid)
	if err != nil {
		log.Error("UserTasksByMID d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.UserTask, 0)
	for rows.Next() {
		r := &model.UserTask{}
		if err = rows.Scan(&r.ID, &r.MID, &r.TaskID, &r.TaskGroupID, &r.TaskType, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("UserTasksByMID rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// AllGiftRewards get all GiftRewards for cache
func (d *Dao) AllGiftRewards(c context.Context) (res map[int8][]*model.GiftReward, err error) {
	_getAllGiftRewardSQL := "SELECT task_type,reward_id,state,comment,ctime,mtime FROM newcomers_gift_reward WHERE state=0"

	rows, err := d.db.Query(c, _getAllGiftRewardSQL)
	if err != nil {
		log.Error("AllGiftRewards d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int8][]*model.GiftReward)
	for rows.Next() {
		t := &model.GiftReward{}
		if err = rows.Scan(&t.TaskType, &t.RewardID, &t.State, &t.Comment, &t.CTime, &t.MTime); err != nil {
			log.Error("AllGiftRewards rows.Scan error(%v)", err)
			return
		}
		res[t.TaskType] = append(res[t.TaskType], t)
	}
	return
}

// GiftRewardCount get received gift reward count
func (d *Dao) GiftRewardCount(c context.Context, mid int64) (res int, err error) {
	sqlStr := "SELECT count(DISTINCT task_gift_id) FROM newcomers_reward_receive WHERE NOT task_gift_id=0 AND mid=?"
	row := d.db.QueryRow(c, sqlStr, mid)
	err = row.Scan(&res)
	if err != nil {
		log.Error("GiftRewardCount d.db.QueryRow error(%v)", err)
	}
	return
}

// BaseRewardCount get received base reward count
func (d *Dao) BaseRewardCount(c context.Context, mid int64) (res int, err error) {
	sqlStr := "SELECT count(DISTINCT task_group_id) FROM newcomers_reward_receive WHERE NOT task_group_id=0 AND mid=?"
	row := d.db.QueryRow(c, sqlStr, mid)
	err = row.Scan(&res)
	if err != nil {
		log.Error("BaseRewardCount d.db.QueryRow error(%v)", err)
	}
	return
}

//CheckTasksForRewardNotify get finish task send notify to user.
func (d *Dao) CheckTasksForRewardNotify(c context.Context, index string, id int64, startMtime, endMtime time.Time, limit int) (res []*model.UserTask, err error) {
	_getUserTaskSQL := "SELECT id, mid, task_id, task_group_id, task_type, state, ctime, mtime FROM newcomers_task_user_%s WHERE state=0 AND id>? AND mtime>=? AND mtime<=? order by id asc limit ?"

	rows, err := d.db.Query(c, fmt.Sprintf(_getUserTaskSQL, index), id, startMtime.Format("2006-01-02 15:04:05"), endMtime.Format("2006-01-02 15:04:05"), limit)
	if err != nil {
		log.Error("UserTasks d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.UserTask, 0)
	for rows.Next() {
		r := &model.UserTask{}
		if err = rows.Scan(&r.ID, &r.MID, &r.TaskID, &r.TaskGroupID, &r.TaskType, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("UserTasks rows.Scan error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}
