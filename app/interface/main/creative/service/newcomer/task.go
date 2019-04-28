package newcomer

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"time"

	"fmt"
	"go-common/app/interface/main/creative/model/newcomer"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	creatorMID        int64 = 37090048 //创作中心哔哩哔哩号
	msgPendantCode          = "1_17_4"
	msgPendantTitle         = "恭喜你通关UP主任务-新手任务！"
	msgPendantContent       = "新手任务通关奖励“小试身手”头像挂件已发送至你的“个人中心-我的头像-我的挂件”中，快来查看和佩戴吧！更多任务奖励正在准备中，敬请期待哦~"
)

// UserTaskInfo def user task info
type UserTaskInfo struct {
	Mid           int64
	UserTaskLevel int8
	Follower      int64
	IsChange      bool
}

// AppIndexNewcomer index newcomer
func (s *Service) AppIndexNewcomer(c context.Context, mid int64, plat string) (res *newcomer.AppIndexNewcomer, err error) {
	if !s.c.TaskCondition.AppIndexSwitch {
		return
	}
	var (
		count     int64
		todoTasks []*newcomer.UserTask
	)
	res = &newcomer.AppIndexNewcomer{
		H5URL: s.c.H5Page.Mission,
	}
	// check task bind
	count, err = s.CheckUserTaskBind(c, mid)
	if err != nil {
		return
	}
	showTaskIDs := make([]int64, 0)
	if count > 0 {
		// already get tasks
		res.TaskReceived = newcomer.BindTask
		todoTasks, err = s.newc.UserTasksByMID(c, mid)
		if err != nil {
			log.Error("AppIndexNewcomer s.newc.UserTasks mid(%d)|error(%v)", mid, err)
			return
		}
		// get at most 2 tasks
		maxCnt := 0
		for _, v := range todoTasks {
			if maxCnt >= 2 {
				break
			}
			showTaskIDs = append(showTaskIDs, v.TaskID)
			maxCnt++
		}
	} else {
		// not bind
		res.TaskReceived = newcomer.NoBindTask
		if len(s.TaskCache) < 2 {
			return
		}
		for i := 0; i < 2; i++ {
			t := s.TaskCache[i]
			showTaskIDs = append(showTaskIDs, t.ID)
		}
	}

	// fill app tasks struct
	for _, v := range showTaskIDs {
		if t, ok := s.TaskMapCache[v]; ok {
			appt := &newcomer.AppTasks{
				ID:    t.ID,
				Title: t.Title,
				Type:  t.Type,
			}
			// get redirect url
			if _, ok := newcomer.TaskRedirectMap[plat]; ok {
				if r, ook := newcomer.TaskRedirectMap[plat][t.TargetType]; ook {
					appt.Label = r[0]
					appt.Redirect = r[1]
				}
			}
			res.AppTasks = append(res.AppTasks, appt)
		}
	}
	if len(res.AppTasks) == 0 {
		res = nil
	}
	return
}

// IndexNewcomer index newcomer
func (s *Service) IndexNewcomer(c context.Context, mid int64) (res *newcomer.IndexNewcomer, err error) {
	res = &newcomer.IndexNewcomer{}

	// check task bind
	count, err := s.CheckUserTaskBind(c, mid)
	if err != nil {
		return
	}
	if count > 0 {
		res.TaskReceived = newcomer.BindTask
	} else {
		res.TaskReceived = newcomer.NoBindTask
	}

	var (
		upCount   int   // Number of submissions
		viewCount int64 // Number of plays
	)
	// if archives == 0 && plays == 0 , means zero data
	dt := time.Now().AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	stat, err := s.data.UpStat(c, mid, dt)
	if err != nil {
		log.Error("IndexNewcomer s.data.UpStat mid(%d)|error(%v)", mid, err)
	} else {
		viewCount = stat.View
	}
	up, err := s.arc.UpCount(c, mid)
	if err != nil {
		log.Error("IndexNewcomer s.arc.UpCount mid(%d)|error(%v)", mid, err)
	} else {
		upCount = up
	}
	log.Info("IndexNewcomer mid(%d)|stat.view(%d)|upcount(%d)", mid, viewCount, upCount)

	if upCount > 0 || viewCount > 0 {
		log.Info("IndexNewcomer upCount>0||viewCount>0 mid(%d)|stat.view(%d)|upcount(%d)", mid, viewCount, upCount)
		// check no receive reward
		noReceiveCount := 0
		noReceiveCount, CountErr := s.getNoReceiveRewardCount(c, mid)
		if CountErr != nil {
			log.Error("IndexNewcomer s.newc.noReceiveRewardCount mid(%d)|error(%v)", mid, err)
		}
		res.NoReceive = noReceiveCount
		res.SubZero = false
	} else {
		res.SubZero = true
	}
	// add three task
	tasks := make([]*newcomer.Task, 0)
	if len(s.TaskCache) > 3 {
		for i := 0; i < 3; i++ {
			tasks = append(tasks, s.TaskCache[i])
		}
		res.Tasks = tasks
	} else {
		log.Error("IndexNewcomer s.newc len(s.TaskCache)<3 mid(%d)|error(%v)", mid, err)
	}

	return
}

// getNoReceiveRewardCount get no receive reward count
func (s *Service) getNoReceiveRewardCount(c context.Context, mid int64) (count int, err error) {
	userTasks, err := s.newc.UserTasks(c, mid)
	if err != nil {
		log.Error("getNoReceiveRewardCount s.newc.UserTasks mid(%d)|error(%v)", mid, err)
		return
	}
	tasks := s.getTasksInfoByType(userTasks, newcomer.DefualtTaskType)
	if len(tasks) == 0 {
		return
	}
	var (
		taskGroupMap   map[int64][]*newcomer.Task
		taskTypeMap    map[int8][]*newcomer.Task
		availableCount int
		receivedCount  int
		groupIDs       []int64
		giftIDs        []int64
	)
	// group by groupID & taskType
	taskGroupMap, taskTypeMap = s.groupByTasks(tasks)
	groupIDs = make([]int64, 0, len(taskGroupMap))
	for k, v := range taskGroupMap {
		if _, ok := s.TaskGroupRewardMapCache[k]; ok {
			if s.getTaskCompleteCount(v) == len(v) {
				availableCount++
			}
			groupIDs = append(groupIDs, k)
		}
	}
	giftIDs = make([]int64, 0, len(taskTypeMap))
	for k, v := range taskTypeMap {
		if _, ok := s.GiftRewardMapCache[k]; ok {
			if s.getTaskCompleteCount(v) == len(v) {
				availableCount++
			}
			giftIDs = append(giftIDs, int64(k))
		}
	}
	r1, err := s.newc.BaseRewardCount(c, mid, groupIDs) // 基础奖励 已领取个数
	if err != nil {
		log.Error("getNoReceiveRewardCount s.newc.BaseRewardCount mid(%v)|error(%v)", mid, err)
		return
	}
	r2, err := s.newc.GiftRewardCount(c, mid, giftIDs) // 礼包奖励 已领取个数
	if err != nil {
		log.Error("getNoReceiveRewardCount s.newc.GiftRewardCount mid(%v)|error(%v)", mid, err)
		return
	}
	receivedCount = int(r1 + r2)
	// 可领取未领取的奖励 = 可领取 - 已领取
	count = availableCount - receivedCount
	if count < 0 {
		count = 0
	}
	log.Info("getNoReceiveRewardCount mid(%d)|availableCount(%d)|receivedCount(%d)|count(%d)", mid, availableCount, receivedCount, count)
	return
}

// getUserTaskInfo determine the userTaskLevel
func (s *Service) getUserTaskInfo(c context.Context, mid int64, tasks []*newcomer.Task) (u *UserTaskInfo, err error) {
	u = &UserTaskInfo{
		Mid:           mid,
		UserTaskLevel: newcomer.UserTaskLevel01,
		IsChange:      false,
	}
	var (
		count       int // Number of unfinished novice tasks
		follower    int64
		profileStat *accapi.ProfileStatReply
		taskTypeMap map[int8][]*newcomer.Task
		taskMap     map[int64]*newcomer.Task
	)

	// whether the advanced task is hidden
	if s.isHiddenTaskType(newcomer.AdvancedTaskType) {
		u.UserTaskLevel = newcomer.UserTaskLevel01
		return
	}

	taskMap = make(map[int64]*newcomer.Task)
	taskTypeMap = make(map[int8][]*newcomer.Task)
	for _, task := range tasks {
		if task == nil {
			continue
		}
		t, ok := s.TaskMapCache[task.ID]
		if !ok {
			continue
		}
		tp := *t
		tp.CompleteSate = task.CompleteSate
		taskTypeMap[tp.Type] = append(taskTypeMap[tp.Type], &tp)
		taskMap[tp.ID] = &tp
	}

	// If the user already has an advanced task, return UserTaskLevel02
	if _, ok := taskTypeMap[newcomer.AdvancedTaskType]; ok {
		u.UserTaskLevel = newcomer.UserTaskLevel02
		return
	}

	// Calculate whether the novice task is not fully completed (compare with cache)
	newcomerTasks := s.getTasksByType(newcomer.NewcomerTaskType)
	if len(newcomerTasks) == 0 {
		return
	}
	count = len(newcomerTasks)
	for _, t := range newcomerTasks {
		if task, ok := taskMap[t.ID]; ok {
			if task.CompleteSate == newcomer.TaskCompleted {
				count--
			}
		}
	}
	// judge fans count
	profileStat, profileErr := s.acc.ProfileWithStat(c, mid)
	if profileStat == nil || profileErr != nil {
		log.Error("genUserTaskInfo s.acc.ProfileWithStat mid(%d)|error(%v)", mid, err)
		follower = 0
	} else {
		follower = profileStat.Follower
	}
	//Number of unfinished tasks || fans >=100
	if count == 0 || follower >= s.c.TaskCondition.Fans {
		// insert advancedTask
		tasks = s.TaskTypeMapCache[newcomer.AdvancedTaskType]
		if len(tasks) == 0 {
			log.Warn("genUserTaskInfo no taskType==newcomer.AdvancedTaskType mid(%d)", mid)
			return
		}
		args, placeStr := genBatchParamsBindTasks(tasks, mid)
		_, err = s.newc.BindTasks(c, mid, placeStr, args)
		if err != nil {
			log.Error("genUserTaskInfo s.newc.BindTasks mid(%d)|error(%v)", mid, err)
			return
		}
		u.UserTaskLevel = newcomer.UserTaskLevel02
		u.IsChange = true        // user level changed
		s.putCheckTaskState(mid) // add to checkTaskQueue
	}

	return
}

//RewardReceive insert reward receive records.
func (s *Service) RewardReceive(c context.Context, mid int64, rid int64, ty int8, ip string) (res string, err error) {
	var (
		rewardIDs   []int64
		rewards     []*newcomer.Reward
		lockSuccess bool
	)

	// check completed
	err = s.isRewardComplete(c, mid, rid, ty)
	if err != nil {
		return
	}

	// prevent concurrent collection
	key := s.getReceiveKey(rid, ty, mid)
	if lockSuccess, err = s.newc.Lock(c, key, 1000); !lockSuccess || err != nil {
		if err == nil {
			log.Info("RewardReceive s.newc.Lock mid(%d)|rid(%d)|ty(%d)", mid, rid, ty)
			res = s.c.TaskCondition.ReceiveMsg
			return
		}
		log.Error("RewardReceive s.newc.Lock mid(%d)|rid(%d)|ty(%d)|error(%v)", mid, rid, ty, err)
	}

	// check received
	err = s.isRewardReceived(c, mid, rid, ty)
	if err != nil {
		return
	}

	// get rewards
	if ty == newcomer.RewardGiftType {
		rewardIDs, err = s.newc.GiftRewards(c, int8(rid))
		if err != nil {
			log.Error("RewardReceive s.newc.GiftRewards mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
			return
		}
	} else {
		rewardIDs, err = s.newc.TaskGroupRewards(c, rid)
		if err != nil {
			log.Error("RewardReceive s.newc.TaskGroupRewards mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
			return
		}
	}
	if len(rewardIDs) == 0 {
		err = ecode.RequestErr
		return
	}
	rewards = make([]*newcomer.Reward, 0)
	for _, v := range rewardIDs {
		if r, ok := s.RewardMapCache[v]; ok {
			rewards = append(rewards, r)
		}
	}
	// get user info
	profile, err := s.acc.Profile(c, mid, ip)
	if err != nil || profile == nil {
		log.Error("RewardReceive s.acc.Profile mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
		return
	}

	// receive rewards
	rrs := make([]*newcomer.RewardReceive2, 0)
	for _, v := range rewards {
		r := &newcomer.RewardReceive2{}
		r.MID = mid
		r.OID = rid
		r.Type = ty
		r.RewardID = v.ID
		r.RewardType = v.Type
		if v.IsActive == newcomer.RewardNeedActivate {
			r.State = newcomer.RewardCanActivate
		} else {
			err = s.callBusiness(c, profile, v)
			if err != nil {
				log.Error("RewardReceive s.callBusiness mid(%d)|reward(%+v)|error(%v)", mid, v, err)
				r.State = newcomer.RewardCanActivate
			} else {
				r.State = newcomer.RewardActivatedNotClick
			}
		}
		rrs = append(rrs, r)
	}
	_, err = s.newc.RewardReceive2(c, mid, rrs)
	if err != nil {
		log.Error("RewardReceive s.newc.RewardReceive mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
		s.newc.UnLock(c, key)
		return
	}
	res = s.c.TaskCondition.ReceiveMsg
	log.Info("RewardReceive success mid(%d)|rewardID(%d)|ty(%d)", mid, rid, ty)
	return
}

// Determine if the reward can be collected
func (s *Service) isRewardComplete(c context.Context, mid int64, rid int64, ty int8) error {
	tids := make([]int64, 0)
	tasks := make([]*newcomer.Task, 0)
	if ty == newcomer.RewardBaseType {
		tasks = s.getTasksByGroupID(rid)
	} else if ty == newcomer.RewardGiftType {
		tasks = s.getTasksByType(int8(rid))
	}
	for _, t := range tasks {
		if t == nil {
			continue
		}
		tids = append(tids, t.ID)
	}
	if len(tids) == 0 {
		log.Error("isRewardComplete len(tids) == 0 | mid(%d)", mid)
		return ecode.CreativeNewcomerNotCompleteErr
	}

	count, err := s.newc.RewardCompleteState(c, mid, tids)
	if err != nil {
		log.Error("isRewardComplete s.newc.RewardCompleteState mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
		return err
	}
	if count != 0 {
		err = ecode.CreativeNewcomerNotCompleteErr
		log.Error("isRewardComplete mission not completed mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
		return err
	}
	return nil
}

// Determine whether to receive it repeatedly
func (s *Service) isRewardReceived(c context.Context, mid int64, rid int64, ty int8) error {
	isReceived, err := s.newc.IsRewardReceived(c, mid, rid, ty)
	if err != nil {
		log.Error("isRewardReceived s.newc.IsGiftReceived mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
		return err
	}
	if isReceived {
		err = ecode.CreativeNewcomerRepeatRewardErr
		log.Error("isRewardReceived receive multiple times mid(%d)|rewardID(%d)|ty(%d)|err(%v)", mid, rid, ty, err)
		return err
	}
	return nil
}

// Calling the business side interface
func (s *Service) callBusiness(c context.Context, profile *accapi.Profile, v *newcomer.Reward) (err error) {
	switch v.Type {
	case newcomer.Bcoin: //B币券
		err = s.newc.BCoin(c, profile.Mid, v.PriceID, int64(v.PrizeUnit))
	case newcomer.MemberBuy: //会员购
		err = s.newc.Mall(c, profile.Mid, v.PriceID, profile.Name)
	case newcomer.IncentivePlan: //激励计划
		err = s.SendRewardReceiveLog(c, profile.Mid)
	case newcomer.BigMember: //大会员服务
		err = s.newc.BigMemberCoupon(c, profile.Mid, v.PriceID)
	case newcomer.PersonalCenter: //个人中心
		err = s.newc.Pendant(c, profile.Mid, v.PriceID, int64(v.Expire))
		if err != nil {
			return
		}
		//发送消息通知
		if e := s.newc.SendNotify(c, []int64{profile.Mid}, msgPendantCode, msgPendantTitle, msgPendantContent); e != nil {
			log.Error("callBusiness s.newc.SendNotify mid(%d)|err(%v)", profile.Mid, err)
		}
	}
	return
}

//RewardActivate insert reward receive records.
func (s *Service) RewardActivate(c context.Context, mid, id int64, ip string) (res int64, err error) {
	r, err := s.newc.RewardReceiveByID(c, mid, id)
	if err != nil {
		log.Error("RewardActivate s.newc.RewardIDByReceiveID mid(%d)|receiveID(%d)|err(%v)", mid, id, err)
		return
	}
	// check repeat collection
	if r.State != newcomer.RewardCanActivate {
		err = ecode.CreativeNewcomerRepeatRewardErr
		log.Error("RewardActivate check repeat collection mid(%d)|receiveID(%d)|state(%d)|err(%v)", mid, id, r.State, err)
		return
	}

	// check expire
	reward, ok := s.RewardMapCache[r.RewardID]
	if !ok || reward == nil {
		err = ecode.RequestErr
		log.Error("RewardActivate check expire mid(%d)|receiveID(%d)|rewardID(%d)|err(%v)", mid, id, r.RewardID, err)
		return
	}
	var expireTime int64
	if r.RewardType == newcomer.Bcoin { //奖品到期时间
		expireTime = reward.CTime.Time().Unix() + int64(reward.Expire*24*3600) //B币券显示截止时间
	} else {
		expireTime = r.ReceiveTime.Time().Unix() + int64(reward.Expire*24*3600)
	}
	if (time.Now().Unix() - expireTime) > 0 {
		err = ecode.CreativeNewcomerReceiveExpireErr
		log.Error("RewardActivate check expire mid(%d)|receiveID(%d)|rewardID(%d)|err(%v)", mid, id, r.RewardID, err)
		return
	}

	// get new receive record
	nr, err := s.newc.RewardReceiveByOldInfo(c, r)
	if err != nil || nr == nil {
		log.Error("RewardActivate s.newc.RewardReceiveByOldInfo mid(%d)|receive(%+v)|error(%v)", mid, r, err)
		return
	}

	// prevent concurrent collection
	var lockSuccess bool
	key := s.getActivateKey(mid, id)
	if lockSuccess, err = s.newc.Lock(c, key, 1000); !lockSuccess || err != nil {
		if err == nil {
			log.Info("RewardActivate s.newc.Lock mid(%d)|id(%d)", mid, id)
			res = 0
			return
		}
		log.Error("RewardActivate s.newc.Lock mid(%d)|id(%d)|error(%v)", mid, id, err)
	}

	profile, err := s.acc.Profile(c, mid, ip)
	if err != nil || profile == nil {
		log.Error("RewardActivate s.acc.Profile mid(%d)|receiveID(%d)|rewardID(%d)|err(%v)", mid, id, r.RewardID, err)
		return
	}
	if err = s.callBusiness(c, profile, reward); err != nil {
		log.Error("RewardActivate callBusiness mid(%d)|receiveID(%d)|rewardID(%d)|err(%v)", mid, id, r.RewardID, err)
		return
	}

	res, err = s.newc.RewardActivate2(c, mid, id, nr.ID)
	if err != nil {
		log.Error("s.newc.RewardActivate mid(%d)|id(%d)|err(%v)", mid, id, err)
		s.newc.UnLock(c, key)
		return
	}
	log.Info("RewardActivate success mid(%d)|receiveID(%d)", mid, id)
	return
}

//RewardReceives get reward receive records.
func (s *Service) RewardReceives(c context.Context, mid int64) (res []*newcomer.RewardReceiveGroup, err error) {
	var items map[int8][]*newcomer.RewardReceive
	items, err = s.newc.RewardReceives(c, mid)
	if err != nil {
		log.Error("s.newc.RewardReceives mid(%d)|err(%v)", mid, err)
		return
	}
	if len(items) == 0 {
		log.Error("s.newc.RewardReceives len(items) == 0")
		return
	}

	keys := reflect.ValueOf(items).MapKeys()
	pids := make([]int64, 0, len(keys)) //存储当前奖励 类型对应的 分类id
	for _, v := range keys {
		if k, ok := v.Interface().(int8); ok {
			if pid, ook := s.RewardTyPIDMapCache[k]; ook {
				pids = append(pids, pid)
			}
		}
	}
	sort.Slice(pids, func(i, j int) bool { //按奖励 类型对应id升序
		return pids[i] < pids[j]
	})
	rkTypes := make([]int8, 0, len(keys)) //存储按照 奖励分类添加顺序排列的 奖励类型
	for _, pid := range pids {
		if ty, ok := s.RewardPIDTyMapCache[pid]; ok {
			rkTypes = append(rkTypes, ty)
		}
	}

	log.Info("RewardReceives mid(%d)|pids(%+v)|rkTypes(%+v)", mid, pids, rkTypes)

	res = make([]*newcomer.RewardReceiveGroup, 0)
	for _, k := range rkTypes {

		item, ok := items[k]
		if !ok || len(item) == 0 {
			log.Error("RewardReceives items[k] k(%v) item(%+v) ok(%v)", k, item, ok)
			return
		}

		pid, ok := s.RewardTyPIDMapCache[k]
		if !ok || pid == 0 {
			log.Error("RewardReceives s.RewardTyPIDMapCache[k] pid(%v) ok(%v)", pid, ok)
			return
		}

		pr, ok := s.RewardMapCache[pid]
		if !ok || pr == nil { //获取奖励类别名称和logo
			log.Error("RewardReceives s.RewardMapCache[pid] pid(%d) pr(%+v) ok3(%v)", pid, pr, ok)
			return
		}

		s0 := make([]*newcomer.RewardReceive, 0, len(item))
		s1 := make([]*newcomer.RewardReceive, 0, len(item))
		s2 := make([]*newcomer.RewardReceive, 0, len(item))
		for _, v := range item {
			r, ok := s.RewardMapCache[v.RewardID]
			if !ok || r == nil {
				log.Error("RewardReceives s.RewardMapCache[v.RewardID] v.RewardID(%v) r(%+v) ok(%v)", v.RewardID, r, ok)
				return
			}

			v.RewardName = r.Name //获取奖品名称
			var expireTime int64
			if v.RewardType == newcomer.Bcoin { //奖品到期时间
				expireTime = r.CTime.Time().Unix() + int64(r.Expire*24*3600)
				v.ExpireTime = xtime.Time(expireTime) //B币券显示截止时间
			} else {
				expireTime = v.ReceiveTime.Time().Unix() + int64(r.Expire*24*3600)
				v.ExpireTime = xtime.Time(expireTime)
			}
			// set RewardExpireNotClick state
			if (time.Now().Unix() - expireTime) > 0 {
				v.State = newcomer.RewardExpireNotClick
			}

			switch v.State { //按照0-可激活 >1-已激活不可点击>2-已过期不可点击 优先级展示
			case newcomer.RewardCanActivate:
				s0 = append(s0, v)
			case newcomer.RewardActivatedNotClick:
				s1 = append(s1, v)
			case newcomer.RewardExpireNotClick:
				s2 = append(s2, v)
			}
		}

		r := &newcomer.RewardReceiveGroup{
			Count:          len(item),
			RewardType:     pr.Type,
			RewardTypeName: pr.Name,
			RewardTypeLogo: pr.Logo,
			Comment:        pr.Comment,
			Items:          append(append(s0, s1...), s2...),
		}
		res = append(res, r)
	}
	return
}

// TaskBind user bind tasks
func (s *Service) TaskBind(c context.Context, mid int64) (res int64, err error) {
	count, err := s.CheckUserTaskBind(c, mid)
	if err != nil {
		return
	}
	if count > 0 {
		err = ecode.CreativeNewcomerReBindTaskErr
		return
	}

	// Determining the number of fans owned by users
	profileStat, err := s.acc.ProfileWithStat(c, mid)
	if err != nil {
		log.Error("TaskBind s.acc.ProfileWithStat mid(%d)|error(%v)", mid, err)
		return
	}
	log.Info("TaskBind s.acc.ProfileWithStat mid(%d)|follower(%d)", mid, profileStat.Follower)
	var tasks []*newcomer.Task
	if profileStat.Follower >= s.c.TaskCondition.Fans {
		tasks = s.TaskTypeMapCache[newcomer.DefualtTaskType]
	} else {
		tasks = s.TaskTypeMapCache[newcomer.NewcomerTaskType]
	}

	args, placeStr := genBatchParamsBindTasks(tasks, mid)
	res, err = s.newc.BindTasks(c, mid, placeStr, args)
	if err != nil {
		log.Error("TaskBind s.newc.BindTasks mid(%v)|error(%v)", mid, err)
		return
	}

	// sync check task status
	s.syncCheckTaskStatus(c, mid, tasks)
	return
}

// genBatchParamsBindTasks generate batch insert parameters
func genBatchParamsBindTasks(tasks []*newcomer.Task, mid int64) ([]interface{}, string) {
	place := make([]string, 0)
	args := make([]interface{}, 0)
	for _, v := range tasks {
		place = append(place, "(?, ?, ?, ?, ?)")
		args = append(args, mid, v.ID, v.GroupID, v.Type, -1)
	}
	placeStr := strings.Join(place, ",")
	return args, placeStr
}

//CheckUserTaskBind determine if the user has bound the task
func (s *Service) CheckUserTaskBind(c context.Context, mid int64) (count int64, err error) {
	count, err = s.newc.UserTaskBind(c, mid)
	if err != nil {
		log.Error("CheckUserTaskBind s.newc.UserTaskBind mid(%v)|error(%v)", mid, err)
		return
	}
	return
}

//TaskMakeup fix unfinish task state
func (s *Service) TaskMakeup(c context.Context, mid int64) (err error) {
	tasks, err := s.newc.UserTasks(c, mid)
	if err != nil {
		log.Error("TaskMakeup s.newc.UserTasks mid(%d)|error(%v)", mid, err)
		err = ecode.RequestErr
		return
	}
	infoTasks := s.getTasksInfoByType(tasks, newcomer.DefualtTaskType)
	ts := make([]*newcomer.Task, 0)
	for _, v := range infoTasks {
		if v == nil {
			continue
		}
		if v.CompleteSate == newcomer.TaskIncomplete {
			ts = append(ts, v)
		}
	}
	s.syncCheckTaskStatus(c, mid, ts)
	return
}

//TaskPubList to apply task list
func (s *Service) TaskPubList(c context.Context, mid int64) (res *newcomer.PubTaskList, err error) {
	res = &newcomer.PubTaskList{}
	tasks, err := s.newc.UserTasks(c, mid)
	if err != nil {
		log.Error("PubTaskList s.newc.UserTasks mid(%d)|error(%v)", mid, err)
		return
	}
	if len(tasks) == 0 {
		log.Warn("PubTaskList No binding task | mid(%d)", mid)
		res.TaskReceived = newcomer.NoBindTask
		return
	}
	res.TaskReceived = newcomer.BindTask
	ts := s.getTasksInfoByType(tasks, newcomer.DefualtTaskType)
	if len(ts) == 0 {
		log.Error("PubTaskList no task | mid(%d)", mid)
		return
	}
	for _, task := range ts {
		if task == nil {
			continue
		}
		tt := &newcomer.PubTask{
			ID:    task.ID,
			Title: task.Title,
			Desc:  task.Desc,
			Type:  task.Type,
			State: task.CompleteSate,
		}
		res.Tasks = append(res.Tasks, tt)
	}

	return
}

// TaskList for task detail
func (s *Service) TaskList(c context.Context, mid int64, ty int8) (res *newcomer.TaskRewardList, err error) {
	var (
		u            *UserTaskInfo
		tasks        []*newcomer.Task
		taskTypeMap  = make(map[int8][]*newcomer.Task)
		taskGroupMap = make(map[int64][]*newcomer.Task)
	)

	// get user tasks
	userTasks, err := s.newc.UserTasks(c, mid)
	if err != nil {
		log.Error("TaskList s.newc.UserTasks mid(%d)|error(%v)", mid, err)
		return
	}
	if len(userTasks) == 0 {
		// return ：User did not receive the task
		res = &newcomer.TaskRewardList{
			TaskReceived: newcomer.NoBindTask,
			TaskType:     ty,
		}
		log.Warn("TaskList user did not receive the task mid(%d)", mid)
		return
	}

	// get user info
	u, err = s.getUserTaskInfo(c, mid, userTasks)
	if err != nil {
		return
	}

	// Judging according to the "user logic level", showing the task type
	if ty == newcomer.DefualtTaskType {
		ty = u.UserTaskLevel
	} else if ty == newcomer.AdvancedTaskType {
		if ty > u.UserTaskLevel {
			log.Warn("TaskList user unlocked this task mid(%d)", mid)
			// return ：User has not unlocked this type of task
			res = &newcomer.TaskRewardList{
				TaskReceived: newcomer.BindTask,
				TaskType:     ty,
			}
			return
		}
	}

	// get tasks
	tasks = s.getTasksInfoByType(userTasks, newcomer.DefualtTaskType)
	if len(tasks) == 0 {
		err = ecode.CreativeNewcomerNoTask
		log.Error("TaskList s.GetTaskByType len(tasks) == 0")
		return
	}

	// group by groupID & taskType
	for _, t := range tasks {
		if t == nil {
			continue
		}
		if t.Type == ty {
			taskGroupMap[t.GroupID] = append(taskGroupMap[t.GroupID], t)
		}
		taskTypeMap[t.Type] = append(taskTypeMap[t.Type], t)
	}

	// task_kind
	taskKinds := s.getTaskKindsData(u, taskTypeMap)

	// task_gift
	taskGift, err := s.getTaskGiftData(c, mid, taskTypeMap, newcomer.FromWeb)
	if err != nil {
		return
	}

	// task_groups
	tgs, err := s.getTaskGroupData(c, mid, taskGroupMap)
	if err != nil {
		return
	}

	res = &newcomer.TaskRewardList{
		TaskReceived: newcomer.BindTask,
		TaskType:     ty,
		TaskKinds:    taskKinds,
		TaskGroups:   tgs,
		TaskGift:     taskGift,
	}
	return
}

// getTaskGroupData for get task_group data
func (s *Service) getTaskGroupData(c context.Context, mid int64, taskMap map[int64][]*newcomer.Task) (tgs []*newcomer.TaskRewardGroup, err error) {
	if len(taskMap) == 0 {
		return
	}
	tgs = make([]*newcomer.TaskRewardGroup, 0, len(taskMap))
	ranks := make([]int64, 0, len(taskMap))
	groups := make([]int64, 0, len(taskMap))
	rankGroupMap := make(map[int64]int64)
	// sort by groupID
	for key := range taskMap {
		taskGroup, ok := s.TaskGroupMapCache[key]
		if !ok {
			continue
		}
		ranks = append(ranks, taskGroup.Rank)
		groups = append(groups, key)
		rankGroupMap[taskGroup.Rank] = key
	}
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] < ranks[j]
	})

	// get TaskRewardGroup
	for _, k := range ranks {
		taskGroup, ok := taskMap[rankGroupMap[k]]
		if !ok || len(taskGroup) == 0 {
			log.Error("genTaskGroupData taskGroup not exist ID(%d)", rankGroupMap[k])
			continue
		}
		rewardState := newcomer.RewardNotAvailable    //奖励领取状态
		total := len(taskGroup)                       // 任务组总数
		complete := s.getTaskCompleteCount(taskGroup) // 完成的任务数量
		if complete == total {
			rewardState = newcomer.RewardAvailable
		}
		rewards := s.getRewardsByGroupID(rankGroupMap[k]) // 任务组对应的奖励

		tg := &newcomer.TaskRewardGroup{
			TaskType:    taskGroup[0].Type,
			Tasks:       taskGroup,
			Rewards:     rewards,
			GroupID:     rankGroupMap[k],
			RewardState: int8(rewardState),
			Completed:   int64(complete),
			Total:       int64(total),
		}
		tgs = append(tgs, tg)
	}

	// determine if the reward is received
	rewardRecvGroup, err := s.newc.RewardReceivedGroup(c, mid, groups)
	if err != nil {
		log.Error("genTaskGroupData s.newc.RewardReceivedGroup mid(%d)|error(%v)", mid, err)
		return
	}
	for _, v := range tgs {
		for _, j := range rewardRecvGroup {
			if v.GroupID == int64(j) {
				v.RewardState = newcomer.RewardReceived
			}
		}
	}
	return
}

// getTaskGiftData for get task_gift data
func (s *Service) getTaskGiftData(c context.Context, mid int64, taskTypeMap map[int8][]*newcomer.Task, from int8) (taskGift []*newcomer.TaskGift, err error) {
	if len(taskTypeMap) == 0 {
		return
	}
	taskGift = make([]*newcomer.TaskGift, 0, len(taskTypeMap))
	// get gifts by taskType
	gifts := make(map[int8][]*newcomer.Reward)
	for k := range taskTypeMap {
		gifts[k] = s.getRewardsByTaskType(k)
	}

	for tType, giftRewards := range gifts {
		giftState := newcomer.RewardAvailable
		complete := s.getTaskCompleteCount(taskTypeMap[tType])
		total := len(taskTypeMap[tType])
		if complete != total {
			if from == newcomer.FromWeb {
				giftState = newcomer.RewardNotAvailable
			} else if from == newcomer.FromH5 {
				giftState = newcomer.RewardUnlock
			}
		}

		// judge the gift is received
		isGiftReceived, err := s.newc.IsRewardReceived(c, mid, int64(tType), newcomer.RewardGiftType)
		if err != nil {
			log.Error("genTaskGiftData s.newc.IsGiftReceived mid(%d)|error(%v)", mid, err)
			continue
		}
		if isGiftReceived {
			giftState = newcomer.RewardReceived
		}

		tg := &newcomer.TaskGift{
			State:   int8(giftState),
			Type:    tType,
			Rewards: giftRewards,
		}
		taskGift = append(taskGift, tg)
	}
	return
}

// getTaskKindsData for get task_kinds data
func (s *Service) getTaskKindsData(u *UserTaskInfo, taskTypeMap map[int8][]*newcomer.Task) (taskKinds []*newcomer.TaskKind) {
	taskKinds = make([]*newcomer.TaskKind, 0, len(taskTypeMap))
	for tType, tasks := range taskTypeMap {
		complete := s.getTaskCompleteCount(tasks) // 完成数量
		total := len(tasks)                       // 总数
		state := newcomer.RewardNotAvailable      // 完成状态
		if complete == len(tasks) {
			state = newcomer.RewardAvailable
		}
		kind := &newcomer.TaskKind{
			Type:      tType,
			Completed: int64(complete),
			Total:     int64(total),
			State:     int8(state),
		}
		// determine the unlock type by user_task level
		switch tType {
		case newcomer.NewcomerTaskType:
			taskKinds = append(taskKinds, kind)
		case newcomer.AdvancedTaskType:
			if u.UserTaskLevel == newcomer.UserTaskLevel01 {
				kind.State = newcomer.RewardUnlock
			}
			taskKinds = append(taskKinds, kind)
		}
	}
	return
}

// getTaskCompleteCount get tasks complete count
func (s *Service) getTaskCompleteCount(task []*newcomer.Task) int {
	complete := 0
	for _, t := range task {
		if t == nil {
			continue
		}
		if t.CompleteSate == newcomer.TaskCompleted {
			complete++
		}
	}
	return complete
}

// groupByTasks group by groupID & taskType
func (s *Service) groupByTasks(tasks []*newcomer.Task) (taskGroupMap map[int64][]*newcomer.Task, taskTypeMap map[int8][]*newcomer.Task) {
	taskGroupMap = make(map[int64][]*newcomer.Task)
	taskTypeMap = make(map[int8][]*newcomer.Task)
	for _, t := range tasks {
		if t == nil {
			continue
		}
		taskGroupMap[t.GroupID] = append(taskGroupMap[t.GroupID], t)
		taskTypeMap[t.Type] = append(taskTypeMap[t.Type], t)
	}
	return
}

// getReceiveKey get receive redis key
func (s *Service) getReceiveKey(oid int64, ty int8, mid int64) string {
	return fmt.Sprintf("%d_%d_%d", oid, ty, mid)
}

// getActivateKey get active key
func (s *Service) getActivateKey(mid int64, id int64) string {
	return fmt.Sprintf("%d_%d", mid, id)
}
