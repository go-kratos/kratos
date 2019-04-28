package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"go-common/app/admin/main/creative/model/task"
	"go-common/library/ecode"
	"go-common/library/log"
)

//TaskGroupRewards for task group & rewards.
func (s *Service) TaskGroupRewards(ids []int64) (res map[int64][]*task.TaskGroupReward, err error) {
	tgrs := []*task.TaskGroupReward{}
	if err = s.DB.Model(&task.TaskGroupReward{}).Where("task_group_id IN (?)", ids).Where("state>=0").Find(&tgrs).Error; err != nil {
		log.Error("s.TaskGroupRewards ids(%+v) error(%v)", ids, err)
	}
	if len(tgrs) == 0 {
		return
	}

	res = make(map[int64][]*task.TaskGroupReward)
	for _, v := range tgrs {
		if v != nil {
			res[v.TaskGroupID] = append(res[v.TaskGroupID], v)
		}
	}
	return
}

func (s *Service) getRewards(rids []int64) (res []*task.RewardResult) {
	_, rewardMap := s.loadRewards()
	if len(rids) == 0 || len(rewardMap) == 0 {
		return
	}

	sort.Slice(rids, func(i, j int) bool {
		return rids[i] < rids[j]
	})
	res = make([]*task.RewardResult, 0, len(rids))
	for _, rid := range rids {
		if v, ok := rewardMap[rid]; ok {
			res = append(res, &task.RewardResult{
				RewardID:   v.ID,
				RewardName: v.Name,
			})
		}
	}
	return
}

//AddTaskGroup for add task group.
func (s *Service) AddTaskGroup(v *task.TaskGroup, rewardsIDs []int64) (id int64, err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	v.CTime = now
	tx := s.DB.Begin()
	v.State = task.StateHide
	if err = tx.Create(v).Error; err != nil {
		log.Error("addGroup error(%v)", err)
		tx.Rollback()
		return
	}

	if err = tx.Model(&task.TaskGroup{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"rank": v.ID,
	}).Error; err != nil {
		log.Error("addGroup error(%v)", err)
		tx.Rollback()
		return
	}

	if len(rewardsIDs) > 0 {
		valReward := make([]string, 0, len(rewardsIDs))
		valRewardArgs := make([]interface{}, 0)
		for _, rid := range rewardsIDs {
			valReward = append(valReward, "(?, ?, ?, ?, ?, ?)")
			valRewardArgs = append(valRewardArgs, v.ID, rid, task.StateNormal, v.Comment, now, now)
		}
		sqlRewardStr := fmt.Sprintf("INSERT INTO newcomers_grouptask_reward (task_group_id, reward_id, state, comment, ctime, mtime) VALUES %s", strings.Join(valReward, ","))
		if err = tx.Exec(sqlRewardStr, valRewardArgs...).Error; err != nil {
			log.Error("addGroup link reward error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return v.ID, nil
}

//EditTaskGroup for edit task group.
func (s *Service) EditTaskGroup(v *task.TaskGroup, rewardsIDs []int64) (id int64, err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	tg := &task.TaskGroup{}
	if err = s.DB.Model(&task.TaskGroup{}).Where("id=?", v.ID).Find(tg).Error; err != nil {
		log.Error("EditTaskGroup link reward error(%v)", err)
		return
	}
	if tg == nil {
		return
	}

	tx := s.DB.Begin()
	if err = tx.Model(&task.TaskGroup{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"root_type": v.RootType,
		"type":      v.Type,
	}).Error; err != nil {
		log.Error("editGroup error(%v)", err)
		tx.Rollback()
		return
	}

	if len(rewardsIDs) > 0 {
		var tgr task.TaskGroupReward
		if err = tx.Where("task_group_id =?", v.ID).Delete(&tgr).Error; err != nil {
			log.Error("editGroup delete old group id(%d)|error(%v)", v.ID, err)
			tx.Rollback()
			return
		}

		valReward := make([]string, 0, len(rewardsIDs))
		valRewardArgs := make([]interface{}, 0)
		for _, rid := range rewardsIDs {
			valReward = append(valReward, "(?, ?, ?, ?, ?, ?)")
			valRewardArgs = append(valRewardArgs, v.ID, rid, task.StateNormal, v.Comment, now, now)
		}
		sqlRewardStr := fmt.Sprintf("INSERT INTO newcomers_grouptask_reward (task_group_id, reward_id, state, comment, ctime, mtime) VALUES %s ON DUPLICATE KEY UPDATE task_group_id=VALUES(task_group_id), reward_id=VALUES(reward_id)", strings.Join(valReward, ","))
		if err = tx.Exec(sqlRewardStr, valRewardArgs...).Error; err != nil {
			log.Error("editGroup link reward error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return v.ID, nil
}

//OrderTaskGroup for order task group.
func (s *Service) OrderTaskGroup(v *task.OrderTask) (err error) {
	tg := &task.TaskGroup{}
	if err = s.DB.Find(tg, v.ID).Error; err != nil {
		log.Error("orderGroup error(%v)", err)
		return
	}
	stg := &task.TaskGroup{}
	if err = s.DB.Find(stg, v.SwitchID).Error; err != nil {
		log.Error("orderGroup error(%v)", err)
		return
	}

	tx := s.DB.Begin()
	if err = tx.Model(&task.TaskGroup{}).Where("id=?", v.ID).Updates(
		map[string]interface{}{
			"rank": v.SwitchRank,
		},
	).Error; err != nil {
		log.Error("orderGroup error(%v)", err)
		tx.Rollback()
		return
	}
	if err = tx.Model(&task.TaskGroup{}).Where("id=?", v.SwitchID).Updates(
		map[string]interface{}{
			"rank": v.Rank,
		},
	).Error; err != nil {
		log.Error("orderGroup error(%v)", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//UpStateGroup for update task group.
func (s *Service) UpStateGroup(id int64, state int8) (err error) {
	tg := &task.TaskGroup{}
	if err = s.DB.Find(tg, id).Error; err != nil {
		log.Error("UpStateGroup id(%d) error(%v)", id, err)
		return
	}
	if tg.ID == 0 {
		err = ecode.NothingFound
		return
	}

	if err = s.DB.Model(&task.TaskGroup{}).Where("id=?", id).Updates(map[string]interface{}{
		"state": state,
	}).Error; err != nil {
		log.Error("UpStateGroup id(%d) state(%d) error(%v)", id, state, err)
		return
	}
	return
}

//TaskGroup for task group.
func (s *Service) TaskGroup(id int64) (res *task.TaskGroup, err error) {
	var tg task.TaskGroup
	if err = s.DB.Model(&task.TaskGroup{}).Where("id=?", id).Find(&tg).Error; err != nil {
		log.Error("s.TaskGroup id (%d) error(%v)", id, err)
		return
	}
	if tg.ID == 0 {
		return
	}

	tgrsMap, _ := s.TaskGroupRewards([]int64{id})
	if rs, ok := tgrsMap[id]; ok {
		rids := make([]int64, 0, len(rs))
		for _, r := range rs {
			if r != nil {
				rids = append(rids, r.RewardID)
			}
		}
		tg.Reward = s.getRewards(rids)
	}

	res = &tg
	return
}

//AddSubtask for add sub task.
func (s *Service) AddSubtask(v *task.Task, rewardsIDs []int64) (id int64, err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	v.CTime = now
	v.MTime = now
	v.State = task.StateHide
	tx := s.DB.Begin()
	if err = tx.Create(v).Error; err != nil {
		log.Error("AddSubtask error(%v)", err)
		tx.Rollback()
		return
	}
	if v.ID == 0 {
		log.Error("AddSubtask v.ID(%d)", v.ID)
		tx.Rollback()
		return
	}

	if err = tx.Model(&task.Task{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"rank": v.ID,
	}).Error; err != nil {
		log.Error("AddSubtask error(%v)", err)
		tx.Rollback()
		return
	}

	if len(rewardsIDs) > 0 {
		valReward := make([]string, 0, len(rewardsIDs))
		valRewardArgs := make([]interface{}, 0)
		for _, rid := range rewardsIDs {
			valReward = append(valReward, "(?, ?, ?, ?, ?, ?)")
			valRewardArgs = append(valRewardArgs, v.ID, rid, task.StateNormal, v.Comment, now, now)
		}
		sqlRewardStr := fmt.Sprintf("INSERT INTO newcomers_task_reward (task_id, reward_id, state, comment, ctime, mtime) VALUES %s", strings.Join(valReward, ","))
		if err = tx.Exec(sqlRewardStr, valRewardArgs...).Error; err != nil {
			log.Error("AddSubtask link reward error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return v.ID, nil
}

//EditSubtask for edit sub task.
func (s *Service) EditSubtask(v *task.Task, rewardsIDs []int64) (id int64, err error) {
	tk := &task.Task{}
	if err = s.DB.Model(&task.TaskGroup{}).Where("id=?", v.ID).Find(tk).Error; err != nil {
		return
	}
	if tk == nil {
		err = ecode.NothingFound
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	v.CTime = tk.CTime
	v.MTime = now
	v.State = tk.State //编辑不更新状态
	tx := s.DB.Begin()
	if err = tx.Save(v).Error; err != nil { //Save将包括执行更新SQL时的所有字段，即使它没有更改
		log.Error("editSubtask error(%v)", err)
		tx.Rollback()
		return
	}

	if len(rewardsIDs) > 0 {
		var tr task.TaskReward
		if err = tx.Where("task_id =?", v.ID).Delete(&tr).Error; err != nil {
			log.Error("editSubtask delete old task id(%d)|error(%v)", v.ID, err)
			tx.Rollback()
			return
		}

		valReward := make([]string, 0, len(rewardsIDs))
		valRewardArgs := make([]interface{}, 0)
		for _, rid := range rewardsIDs {
			valReward = append(valReward, "(?, ?, ?, ?, ?, ?)")
			valRewardArgs = append(valRewardArgs, v.ID, rid, task.StateNormal, v.Comment, now, now)
		}
		sqlRewardStr := fmt.Sprintf("INSERT INTO newcomers_task_reward (task_id, reward_id, state, comment, ctime, mtime) VALUES %s ON DUPLICATE KEY UPDATE task_id=VALUES(task_id), reward_id=VALUES(reward_id)", strings.Join(valReward, ","))
		if err = tx.Exec(sqlRewardStr, valRewardArgs...).Error; err != nil {
			log.Error("editSubtask link reward error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return v.ID, nil
}

//OrderSubTask for order sub task.
func (s *Service) OrderSubTask(v *task.OrderTask) (err error) {
	tg := &task.Task{}
	if err = s.DB.Find(tg, v.ID).Error; err != nil {
		log.Error("OrderSubTask error(%v)", err)
		return
	}
	stg := &task.Task{}
	if err = s.DB.Find(stg, v.SwitchID).Error; err != nil {
		log.Error("OrderSubTask error(%v)", err)
		return
	}

	tx := s.DB.Begin()
	if err = tx.Model(&task.Task{}).Where("id=?", v.ID).Updates(
		map[string]interface{}{
			"rank": v.SwitchRank,
		},
	).Error; err != nil {
		log.Error("OrderSubTask error(%v)", err)
		tx.Rollback()
		return
	}
	if err = tx.Model(&task.Task{}).Where("id=?", v.SwitchID).Updates(
		map[string]interface{}{
			"rank": v.Rank,
		},
	).Error; err != nil {
		log.Error("OrderSubTask error(%v)", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//UpStateSubTask for update sub task state.
func (s *Service) UpStateSubTask(id int64, state int8) (err error) {
	tg := &task.Task{}
	if err = s.DB.Find(tg, id).Error; err != nil {
		return
	}
	if tg.ID == 0 {
		err = ecode.NothingFound
		return
	}

	if err = s.DB.Model(&task.Task{}).Where("id=?", id).Updates(map[string]interface{}{
		"state": state,
	}).Error; err != nil {
		log.Error("UpStateSubTask id(%d) state(%d) error(%v)", id, state, err)
		return
	}
	return
}

//TransferSubtask for transfer one sub task from some one group to another one.
func (s *Service) TransferSubtask(id, gid int64) (err error) {
	tg := &task.TaskGroup{}
	if err = s.DB.Find(tg, gid).Error; err != nil {
		return
	}
	if tg.ID == 0 {
		err = ecode.CreativeNewcomerGroupIDErr
		return
	}

	tk := &task.Task{}
	if err = s.DB.Find(tk, id).Error; err != nil {
		return
	}
	if tk.ID == 0 {
		err = ecode.NothingFound
		return
	}

	if err = s.DB.Model(&task.Task{}).Where("id=?", id).Updates(map[string]interface{}{
		"group_id": gid,
	}).Error; err != nil {
		log.Error("transferSubtask id(%+v) gid(%d) error(%v)", id, gid, err)
		return
	}
	return
}

//Task for task.
func (s *Service) Task(id int64) (res *task.Task, err error) {
	var t task.Task
	if err = s.DB.Model(&task.Task{}).Where("id=?", id).Find(&t).Error; err != nil {
		log.Error("s.Task id (%d) error(%v)", id, err)
		return
	}
	if t.ID == 0 {
		return
	}

	trsMap, _ := s.TaskRewards([]int64{id})
	if rs, ok := trsMap[id]; ok {
		rids := make([]int64, 0, len(rs))
		for _, r := range rs {
			if r != nil {
				rids = append(rids, r.RewardID)
			}
		}
		t.Reward = s.getRewards(rids)
	}
	res = &t
	return
}

//TaskRewards for task & rewards.
func (s *Service) TaskRewards(ids []int64) (res map[int64][]*task.TaskReward, err error) {
	trs := []*task.TaskReward{}
	if err = s.DB.Model(&task.TaskReward{}).Where("state>=0 AND task_id IN (?)", ids).Find(&trs).Error; err != nil {
		log.Error("s.TaskRewards ids(%+v) error(%v)", ids, err)
		return
	}
	if len(trs) == 0 {
		return
	}

	res = make(map[int64][]*task.TaskReward)
	for _, v := range trs {
		if v != nil {
			res[v.TaskID] = append(res[v.TaskID], v)
		}
	}
	return
}

//TasksByGroupIDsMap for task map
func (s *Service) TasksByGroupIDsMap(gids []int64) (res map[int64][]*task.Task, err error) {
	tks := []*task.Task{}
	if err = s.DB.Model(&task.Task{}).Where("state>=0 AND type>0 AND group_id IN (?)", gids).Find(&tks).Error; err != nil {
		log.Error("s.TasksByGroupIDsMap id (%+v) error(%v)", gids, err)
		return
	}
	if len(tks) == 0 {
		return
	}

	tkMap := make(map[int64]*task.Task)
	tgMap := make(map[int64][]*task.Task)
	ids := make([]int64, 0, len(tks))
	for _, v := range tks {
		if v == nil {
			continue
		}
		if v.Rank == 0 {
			v.Rank = v.ID
		}
		ids = append(ids, v.ID)
		tkMap[v.ID] = v
		tgMap[v.GroupID] = append(tgMap[v.GroupID], v)
	}

	trsMap, _ := s.TaskRewards(ids)
	for _, id := range ids {
		rs, ok := trsMap[id]
		if !ok || len(rs) == 0 {
			continue
		}
		rids := make([]int64, 0, len(rs))
		for _, r := range rs {
			if r != nil {
				rids = append(rids, r.RewardID)
			}
		}
		tkMap[id].Reward = s.getRewards(rids)
	}

	res = make(map[int64][]*task.Task)
	for _, gid := range gids {
		tks, ok := tgMap[gid]
		if !ok || len(tks) == 0 {
			continue
		}
		for _, tk := range tks {
			if v, ok := tkMap[tk.ID]; ok {
				res[gid] = append(res[gid], v)
			}
		}
	}
	for _, v := range res {
		sort.Slice(v, func(i, j int) bool {
			return v[i].Rank < v[j].Rank
		})
	}
	return
}

//TaskList for task list
func (s *Service) TaskList(ty int8) (res []*task.TaskGroup, err error) {
	tgs := []*task.TaskGroup{}
	db := s.DB.Model(&task.TaskGroup{}).Where("state>=0 AND root_type>0 AND type>0")
	if ty > 0 {
		db = db.Where("type=?", ty)
	}

	if err = db.Find(&tgs).Error; err != nil {
		log.Error("TaskList %v\n", err)
		return
	}
	if len(tgs) == 0 {
		return
	}

	gids := make([]int64, 0, len(tgs))
	tgMap := make(map[int64]*task.TaskGroup)
	for _, v := range tgs {
		if v == nil {
			continue
		}
		gids = append(gids, v.ID)
		tgMap[v.ID] = v
	}

	tgrsMap, _ := s.TaskGroupRewards(gids)
	for _, id := range gids {
		if rs, ok := tgrsMap[id]; ok {
			if len(rs) == 0 {
				continue
			}
			rids := make([]int64, 0, len(rs))
			for _, r := range rs {
				if r != nil {
					rids = append(rids, r.RewardID)
				}
			}
			tgMap[id].Reward = s.getRewards(rids)
		}
	}

	tkMap, _ := s.TasksByGroupIDsMap(gids)
	res = tgs
	for _, v := range res {
		if v == nil {
			continue
		}
		if g, ok := tgMap[v.ID]; ok {
			v.Reward = g.Reward
		}
		if tks, ok := tkMap[v.ID]; ok {
			v.Tasks = tks
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Rank < res[j].Rank
	})
	return
}

//loadRewards for task reward.
func (s *Service) loadRewards() (res []*task.Reward, rwMap map[int64]*task.Reward) {
	var (
		rds = []*task.Reward{}
		err error
	)

	if err = s.DB.Order("id ASC").Find(&rds).Error; err != nil {
		log.Error("loadRewards error(%v)", err)
		return
	}

	res = make([]*task.Reward, 0, len(rds))
	rwMap = make(map[int64]*task.Reward)
	top := make(map[int64]*task.Reward)
	for _, v := range rds {
		if v == nil {
			continue
		}
		rwMap[v.ID] = v

		if v.ParentID == 0 {
			top[v.ID] = v        //映射一级对象
			res = append(res, v) //追加一级对象
		}
	}

	for _, ch := range rds {
		if ch == nil {
			continue
		}
		if p, ok := top[ch.ParentID]; ok && p != nil && p.Type == ch.Type { //为一级对象增加子对象，注:要append满足条件的ch，若append p 则会造成递归导致stack overflow
			p.Children = append(p.Children, ch)
		}
	}
	return
}

//RewardTree for reward tree.
func (s *Service) RewardTree() (res []*task.Reward) {
	res, _ = s.loadRewards()
	return
}

//ViewReward for view one reward.
func (s *Service) ViewReward(id int64) (res *task.Reward, err error) {
	rd := task.Reward{}
	if err = s.DB.Model(&task.Reward{}).Where("id=?", id).Find(&rd).Error; err != nil {
		return
	}
	res = &rd
	return
}

//AddReward for add one reward.
func (s *Service) AddReward(v *task.Reward) (id int64, err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	v.CTime = now
	v.MTime = now
	if err = s.DB.Create(v).Error; err != nil {
		log.Error("AddReward v(%+v) error(%v)", v, err)
		return
	}
	return v.ID, nil
}

//EditReward for edit one reward.
func (s *Service) EditReward(v *task.Reward) (id int64, err error) {
	rd := &task.Reward{}
	if err = s.DB.Model(&task.Reward{}).Where("id=?", v.ID).Find(rd).Error; err != nil {
		return
	}
	if rd == nil {
		err = ecode.NothingFound
		return
	}

	v.CTime = rd.CTime
	v.MTime = time.Now().Format("2006-01-02 15:04:05")
	v.State = rd.State //编辑不更新状态
	if err = s.DB.Save(v).Error; err != nil {
		log.Error("EditReward v(%+v) error(%v)", v, err)
		return
	}
	return v.ID, nil
}

//UpStateReward for update reward state.
func (s *Service) UpStateReward(id int64, state int8) (err error) {
	rd := &task.Reward{}
	if err = s.DB.Find(rd, id).Error; err != nil {
		return
	}
	if rd.ID == 0 {
		err = ecode.NothingFound
		return
	}

	if rd.ParentID != 0 { //如果是子分类直接更新并返回
		if err = s.DB.Model(&task.Reward{}).Where("id=?", id).Updates(map[string]interface{}{
			"state": state,
		}).Error; err != nil {
			log.Error("UpStateReward parent id(%d) state(%d) error(%v)", id, state, err)
		}
		return
	}

	rds := []*task.Reward{}
	if err = s.DB.Model(&task.Reward{}).Where("parent_id=?", id).Find(&rds).Error; err != nil {
		log.Error("UpStateReward get childen by parent id(%d) error(%v)", rd.ParentID, err)
		return
	}
	ids := make([]int64, 0, len(rds)+1)
	for _, v := range rds {
		if v != nil {
			ids = append(ids, v.ID) //追加子分类id
		}
	}

	ids = append(ids, id) //追加父分类id
	if err = s.DB.Model(&task.Reward{}).Where("id IN (?)", ids).Updates(map[string]interface{}{
		"state": state,
	}).Error; err != nil {
		log.Error("UpStateReward childen ids(%+v) error(%v)", ids, err)
	}
	return
}

//ViewGiftReward for view gift.
func (s *Service) ViewGiftReward(v *task.GiftReward) (res *task.GiftReward, err error) {
	var gfs []*task.GiftReward
	if err = s.DB.Model(&task.GiftReward{}).Where("root_type=? AND task_type=?", v.RootType, v.TaskType).Find(&gfs).Error; err != nil {
		log.Error("ViewGiftReward v(%+v) error(%v)", v, err)
		return
	}
	if len(gfs) == 0 {
		return
	}

	res = &task.GiftReward{
		RootType: v.RootType,
		TaskType: v.TaskType,
	}
	var (
		state   int8
		comment string
	)
	rids := make([]int64, 0, len(gfs))
	for _, gf := range gfs {
		if gf != nil {
			state = gf.State
			comment = gf.Comment
			rids = append(rids, gf.RewardID)
		}
	}
	res.State = state
	res.Comment = comment
	res.Reward = s.getRewards(rids)
	return
}

//ListGiftReward for get gift list.
func (s *Service) ListGiftReward() (res []*task.GiftReward, err error) {
	gfs := []*task.GiftReward{}
	if err = s.DB.Model(&task.GiftReward{}).Where("root_type>0 AND task_type>0").Find(&gfs).Error; err != nil {
		log.Error("ListGiftReward error(%v)", err)
		return
	}
	if len(gfs) == 0 {
		return
	}

	gfMap := make(map[int64][]*task.GiftReward)
	for _, v := range gfs {
		if v != nil {
			gfMap[v.TaskType] = append(gfMap[v.TaskType], v)
		}
	}

	tys := make([]int64, 0, len(gfMap))
	for k := range gfMap {
		tys = append(tys, k)
	}
	sort.Slice(tys, func(i, j int) bool {
		return tys[i] < tys[j]
	})

	res = make([]*task.GiftReward, 0, len(tys))
	for _, ty := range tys {
		gfs, ok := gfMap[ty]
		if !ok && len(gfs) == 0 {
			continue
		}
		var rt uint8
		if task.CheckRootType(uint8(ty)) {
			rt = task.TaskManagement
		} else {
			rt = task.AchievementManagement
		}
		re := &task.GiftReward{
			RootType: rt,
			TaskType: ty,
		}

		var (
			state   int8
			comment string
		)
		rids := make([]int64, 0, len(gfs))
		for _, gf := range gfs {
			if gf != nil {
				state = gf.State
				comment = gf.Comment
				rids = append(rids, gf.RewardID)
			}
		}
		re.State = state
		re.Comment = comment
		re.Reward = s.getRewards(rids)
		res = append(res, re)
	}
	return
}

//AddGiftReward for add gift rewards.
func (s *Service) AddGiftReward(v *task.GiftReward, rewardsIDs []int64) (rows int64, err error) {
	var gfs []*task.GiftReward
	if err = s.DB.Model(&task.GiftReward{}).Where("root_type=? AND task_type=? AND reward_id IN (?)", v.RootType, v.TaskType, rewardsIDs).Find(&gfs).Error; err != nil {
		log.Error("UpGiftReward v(%+v) error(%v)", v, err)
		return
	}
	if len(gfs) != 0 {
		hitMap := make(map[int64]struct{})
		for _, gf := range gfs {
			hitMap[gf.RewardID] = struct{}{}
		}
		for _, rid := range rewardsIDs {
			if _, ok := hitMap[rid]; ok {
				err = ecode.CreativeNewcomerDuplicateGiftRewardIDErr
				log.Error("AddGiftReward rid(%d) error(%v)", rid, err)
				return
			}
		}
	}

	valReward := make([]string, 0, len(rewardsIDs))
	valRewardArgs := make([]interface{}, 0)
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, rid := range rewardsIDs {
		valReward = append(valReward, "(?, ?, ?, ?, ?, ?, ?)")
		valRewardArgs = append(valRewardArgs, v.RootType, v.TaskType, rid, task.StateNormal, v.Comment, now, now)
	}
	sqlRewardStr := fmt.Sprintf("INSERT INTO newcomers_gift_reward (root_type, task_type, reward_id, state, comment, ctime, mtime) VALUES %s", strings.Join(valReward, ","))
	if err = s.DB.Exec(sqlRewardStr, valRewardArgs...).Error; err != nil {
		log.Error("AddGiftReward error(%v)", err)
		return
	}
	return s.DB.RowsAffected, nil
}

//EditGiftReward for edit gift rewards.
func (s *Service) EditGiftReward(v *task.GiftReward, rewardsIDs []int64) (rows int64, err error) {
	var (
		gfs   []*task.GiftReward
		state int8
	)
	if err = s.DB.Model(&task.GiftReward{}).Where("root_type=? AND task_type=?", v.RootType, v.TaskType).Find(&gfs).Error; err != nil {
		log.Error("EditGiftReward v(%+v) error(%v)", v, err)
		return
	}
	if len(gfs) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, gf := range gfs { //获取原来的状态
		state = gf.State
		break
	}

	var gf task.GiftReward
	if err = s.DB.Model(&task.GiftReward{}).Where("root_type =?", v.RootType).Where("task_type =?", v.TaskType).Delete(&gf).Error; err != nil {
		log.Error("EditGiftReward delete old id(%d)|error(%v)", v.ID, err)
		return
	}

	valReward := make([]string, 0, len(rewardsIDs))
	valRewardArgs := make([]interface{}, 0)
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, rid := range rewardsIDs {
		valReward = append(valReward, "(?, ?, ?, ?, ?, ?, ?)")
		valRewardArgs = append(valRewardArgs, v.RootType, v.TaskType, rid, state, v.Comment, now, now)
	}
	sqlRewardStr := fmt.Sprintf("INSERT INTO newcomers_gift_reward (root_type, task_type, reward_id, state, comment, ctime, mtime) VALUES %s", strings.Join(valReward, ","))
	if err = s.DB.Exec(sqlRewardStr, valRewardArgs...).Error; err != nil {
		log.Error("EditGiftReward error(%v)", err)
		return
	}
	return s.DB.RowsAffected, nil
}

//UpGiftReward for update gift reward.
func (s *Service) UpGiftReward(v *task.GiftReward) (rows int64, err error) {
	var gfs []*task.GiftReward
	if err = s.DB.Model(&task.GiftReward{}).Where("root_type=? AND task_type=?", v.RootType, v.TaskType).Find(&gfs).Error; err != nil {
		log.Error("UpGiftReward v(%+v) error(%v)", v, err)
		return
	}
	if len(gfs) == 0 {
		return
	}

	if err = s.DB.Model(&task.GiftReward{}).Where("root_type=? AND task_type=?", v.RootType, v.TaskType).Updates(map[string]interface{}{
		"state": v.State,
	}).Error; err != nil {
		log.Error("UpGiftReward v(%+v) error(%v)", v, err)
		return
	}
	return s.DB.RowsAffected, nil
}

//BatchOnline for subtask & grouptask with state,rank and so on.
func (s *Service) BatchOnline(tgs []*task.TaskGroup) (err error) {
	tasks := make([]*task.Task, 0)
	groups := make([]*task.TaskGroup, 0, len(tgs))
	for _, v := range tgs {
		groups = append(groups, &task.TaskGroup{
			ID:       v.ID,
			Rank:     v.Rank,
			State:    v.State,
			RootType: v.RootType,
			Type:     v.Type,
		})
		tasks = append(tasks, v.Tasks...)
	}

	valGroups := make([]string, 0, len(groups))
	valGroupArgs := make([]interface{}, 0)
	for _, v := range groups {
		valGroups = append(valGroups, "(?, ?, ?)")
		valGroupArgs = append(valGroupArgs, v.ID, v.State, v.Rank)
	}
	sqlGroupStr := fmt.Sprintf("INSERT INTO newcomers_task_group (id, state, rank) VALUES %s ON DUPLICATE KEY UPDATE state=VALUES(state), rank=VALUES(rank)", strings.Join(valGroups, ","))
	if err = s.DB.Exec(sqlGroupStr, valGroupArgs...).Error; err != nil {
		log.Error("BatchOnline update groups error(%v)", err)
		return
	}

	valTasks := make([]string, 0, len(tasks))
	valTaskArgs := make([]interface{}, 0)
	for _, v := range tasks {
		valTasks = append(valTasks, "(?, ?, ?, ?)")
		valTaskArgs = append(valTaskArgs, v.ID, v.State, v.Rank, v.GroupID)
	}
	sqlTaskStr := fmt.Sprintf("INSERT INTO newcomers_task (id, state, rank, group_id) VALUES %s ON DUPLICATE KEY UPDATE state=VALUES(state), rank=VALUES(rank), group_id=VALUES(group_id)", strings.Join(valTasks, ","))
	if err = s.DB.Exec(sqlTaskStr, valTaskArgs...).Error; err != nil {
		log.Error("BatchOnline update tasks error(%v)", err)
	}
	return
}
