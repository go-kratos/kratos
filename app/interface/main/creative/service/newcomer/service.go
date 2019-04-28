package newcomer

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/academy"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/article"
	"go-common/app/interface/main/creative/dao/data"
	"go-common/app/interface/main/creative/dao/medal"
	"go-common/app/interface/main/creative/dao/newcomer"
	"go-common/app/interface/main/creative/dao/order"
	"go-common/app/interface/main/creative/dao/watermark"
	ncMDL "go-common/app/interface/main/creative/model/newcomer"
	"go-common/app/interface/main/creative/service"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
	"sync"
)

//Service struct
type Service struct {
	c     *conf.Config
	newc  *newcomer.Dao
	arc   *archive.Dao
	art   *article.Dao
	acc   *account.Dao
	data  *data.Dao
	aca   *academy.Dao
	wm    *watermark.Dao
	medal *medal.Dao
	order *order.Dao
	//reward
	RewardCache         []*ncMDL.Reward
	RewardMapCache      map[int64]*ncMDL.Reward
	RewardTyPIDMapCache map[int8]int64 //存储奖励 类型-分类id对应关系
	RewardPIDTyMapCache map[int64]int8 //存储奖励 分类id-类型对应关系
	//task
	TaskCache          []*ncMDL.Task
	TaskMapCache       map[int64]*ncMDL.Task
	TaskTypeMapCache   map[int8][]*ncMDL.Task
	TaskRewardMapCache map[int64][]*ncMDL.TaskRewardEntity
	//taskgroup-reward
	TaskGroupRewardMapCache map[int64][]*ncMDL.TaskGroupReward
	//gift-reward
	GiftRewardMapCache map[int8][]*ncMDL.GiftReward
	//initiative check task
	checkTaskChan  chan int64
	checkTaskQueue []int64
	checkTaskDone  chan struct{}
	// task-group
	TaskGroupCache    []*ncMDL.TaskGroupEntity
	TaskGroupMapCache map[int64]*ncMDL.TaskGroupEntity
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:              c,
		newc:           newcomer.New(c),
		arc:            rpcdaos.Arc,
		art:            rpcdaos.Art,
		acc:            rpcdaos.Acc,
		data:           data.New(c),
		aca:            academy.New(c),
		wm:             watermark.New(c),
		medal:          medal.New(c),
		order:          order.New(c),
		checkTaskChan:  make(chan int64, 1000),
		checkTaskQueue: make([]int64, 0, 1000),
		checkTaskDone:  make(chan struct{}),
	}

	s.loadRewards()
	s.loadTasks()
	s.loadTaskGroupRewards()
	s.loadGiftRewards()
	s.loadTaskRewards()
	s.loadTaskGroups()
	go s.loadProc()
	go s.checkTaskStateByMid()

	return s
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.newc.Ping(c); err != nil {
		log.Error("s.newc.Ping err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.newc.Close()
	s.checkTaskDone <- struct{}{}
}

// loadproc
func (s *Service) loadProc() {
	for {
		time.Sleep(3 * time.Minute)
		s.loadRewards()
		s.loadTasks()
		s.loadTaskGroupRewards()
		s.loadGiftRewards()
		s.loadTaskRewards()
		s.loadTaskGroups()
	}
}

//load tags
func (s *Service) loadRewards() {
	res, err := s.newc.Rewards(context.Background())
	if err != nil {
		log.Error("s.newc.Rewards error(%v)", err)
		return
	}
	if len(res) == 0 {
		return
	}
	s.RewardCache = res
	tempRewardMapCache := make(map[int64]*ncMDL.Reward)
	tempRewardTyPIDMapCache := make(map[int8]int64)
	tempRewardPIDTyMapCache := make(map[int64]int8)
	for _, v := range res {
		tempRewardMapCache[v.ID] = v
		if v.ParentID == 0 {
			tempRewardTyPIDMapCache[v.Type] = v.ID
			tempRewardPIDTyMapCache[v.ID] = v.Type
		}
	}
	s.RewardMapCache = tempRewardMapCache
	s.RewardTyPIDMapCache = tempRewardTyPIDMapCache
	s.RewardPIDTyMapCache = tempRewardPIDTyMapCache

}

//load tags
func (s *Service) loadTasks() {
	res, err := s.newc.Tasks(context.Background())
	if err != nil {
		log.Error("s.newc.Tasks error(%v)", err)
		return
	}
	if len(res) == 0 {
		return
	}
	s.TaskCache = res

	temp := make(map[int64]*ncMDL.Task)
	tempMap := make(map[int8][]*ncMDL.Task)
	for _, v := range res {
		temp[v.ID] = v
		switch v.Type {
		case ncMDL.NewcomerTaskType:
			tempMap[ncMDL.NewcomerTaskType] = append(tempMap[ncMDL.NewcomerTaskType], v)
		case ncMDL.AdvancedTaskType:
			tempMap[ncMDL.AdvancedTaskType] = append(tempMap[ncMDL.AdvancedTaskType], v)
		case ncMDL.MonthTaskType:
			tempMap[ncMDL.MonthTaskType] = append(tempMap[ncMDL.MonthTaskType], v)
		}
	}
	s.TaskMapCache = temp
	// 默认分类包含：新手与进阶任务
	tempMap[ncMDL.DefualtTaskType] = append(tempMap[ncMDL.DefualtTaskType], tempMap[ncMDL.NewcomerTaskType]...)
	tempMap[ncMDL.DefualtTaskType] = append(tempMap[ncMDL.DefualtTaskType], tempMap[ncMDL.AdvancedTaskType]...)
	s.TaskTypeMapCache = tempMap
}

// SendRewardReceiveLog for reward receive.
func (s *Service) SendRewardReceiveLog(c context.Context, mid int64) (err error) {
	if env.DeployEnv == env.DeployEnvDev {
		return
	}
	uInfo := &report.UserInfo{
		Business: 141, //创作中心奖励领取行为日志业务ID
		Mid:      mid,
		Type:     0, //0-激励计划奖品领取
		Oid:      mid,
		Platform: "激励计划",
		Ctime:    time.Now(),
		Action:   "incentive_plan_reward_receive", //激励计划奖品领取
		Index:    []interface{}{mid},
		IP:       metadata.String(c, metadata.RemoteIP),
	}
	report.User(uInfo)
	log.Info("s.SendRewardReceiveLog mid(%d)|uInfo(%+v)", mid, uInfo)
	return
}

// SendPendantReceiveLog for reward receive.
func (s *Service) SendPendantReceiveLog(c context.Context, mid int64) (err error) {
	if env.DeployEnv == env.DeployEnvDev {
		return
	}
	uInfo := &report.UserInfo{
		Business: 141, //创作中心奖励领取行为日志业务ID
		Mid:      mid,
		Type:     1, //1-头像挂件奖品领取
		Oid:      mid,
		Platform: "头像挂件",
		Ctime:    time.Now(),
		Action:   "pendant_reward_receive", //头像挂件奖品领取
		Index:    []interface{}{mid},
		IP:       metadata.String(c, metadata.RemoteIP),
	}
	report.User(uInfo)
	log.Info("s.SendPendantReceiveLog mid(%d)|uInfo(%+v)", mid, uInfo)
	return
}

//load gift-reward
func (s *Service) loadGiftRewards() {
	res, err := s.newc.AllGiftRewards(context.Background())
	if err != nil {
		log.Error("s.newc.AllGiftRewards error(%v)", err)
		return
	}
	if len(res) == 0 {
		return
	}
	s.GiftRewardMapCache = res
}

//load taskgroup-reward
func (s *Service) loadTaskGroupRewards() {
	res, err := s.newc.AllTaskGroupRewards(context.Background())
	if err != nil {
		log.Error("s.newc.AllTaskGroupRewards error(%v)", err)
		return
	}
	if len(res) == 0 {
		return
	}
	s.TaskGroupRewardMapCache = res
}

// checkTaskStateByMid read bindMid to check task state
func (s *Service) checkTaskStateByMid() {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case id := <-s.checkTaskChan:
			s.checkTaskQueue = append(s.checkTaskQueue, id)
		case <-ticker.C:
			if len(s.checkTaskQueue) > 0 {
				mid := s.checkTaskQueue[0]
				s.checkTaskQueue = s.checkTaskQueue[1:]
				s.DriveStateByUser(context.Background(), mid)
			}
		case <-s.checkTaskDone:
			log.Info("checkTaskStateByMid close")
			return
		}
	}
}

// putCheckTask put mid to checkTaskQueue
func (s *Service) putCheckTaskState(mid int64) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s.checkTaskChan <- mid
		wg.Done()
	}()
	wg.Wait()
}

//load task-group
func (s *Service) loadTaskGroups() {
	res, err := s.newc.TaskGroups(context.Background())
	if err != nil {
		log.Error("s.newc.TaskGroups error(%v)", err)
		return
	}
	if len(res) == 0 {
		return
	}
	s.TaskGroupCache = res

	temp := make(map[int64]*ncMDL.TaskGroupEntity)
	for _, v := range res {
		temp[v.ID] = v
	}
	s.TaskGroupMapCache = temp
}

//load task-group
func (s *Service) loadTaskRewards() {
	res, err := s.newc.TaskRewards(context.Background())
	if err != nil {
		log.Error("s.newc.TaskRewards error(%v)", err)
		return
	}
	if len(res) == 0 {
		return
	}
	s.TaskRewardMapCache = res
}

// getRewardsByTaskType get rewards by taskType
func (s *Service) getRewardsByTaskType(taskType int8) (res []*ncMDL.Reward) {
	res = make([]*ncMDL.Reward, 0)
	gifts, ok := s.GiftRewardMapCache[taskType]
	if !ok || len(gifts) == 0 {
		return
	}
	for _, v := range gifts {
		if v == nil {
			continue
		}
		r, ok := s.RewardMapCache[v.RewardID]
		if !ok {
			continue
		}

		res = append(res, r)
	}
	return
}

// getRewardsByGroupID get reward by groupID
func (s *Service) getRewardsByGroupID(groupID int64) []*ncMDL.Reward {
	rewards := make([]*ncMDL.Reward, 0)
	rs, ok := s.TaskGroupRewardMapCache[groupID]
	if !ok || len(rs) == 0 {
		log.Error("getRewardsByGroupID reward not exist groupID(%d)", groupID)
		return rewards
	}

	for _, r := range rs {
		if r == nil {
			continue
		}
		if rr, ok := s.RewardMapCache[r.RewardID]; ok {
			rewards = append(rewards, rr)
		}
	}
	return rewards
}

// getTaskByGroupID get tasks by groupID
func (s *Service) getTasksByGroupID(groupID int64) (res []*ncMDL.Task) {
	res = make([]*ncMDL.Task, 0)
	for _, task := range s.TaskCache {
		if task == nil {
			continue
		}
		if _, ok := s.TaskGroupMapCache[task.GroupID]; !ok {
			continue
		}
		if task.GroupID == groupID && task.State == ncMDL.NormalState {
			res = append(res, task)
		}
	}
	return
}

// getTaskByType get tasks by type
func (s *Service) getTasksByType(ty int8) (res []*ncMDL.Task) {
	res = make([]*ncMDL.Task, 0)
	for _, task := range s.TaskCache {
		if task == nil {
			continue
		}
		if _, ok := s.TaskGroupMapCache[task.GroupID]; !ok {
			continue
		}
		if task.Type == ty && task.State == ncMDL.NormalState {
			res = append(res, task)
		}
	}
	return
}

// getTasksInfo get tasks info by type
func (s *Service) getTasksInfoByType(ts []*ncMDL.Task, ty int8) (res []*ncMDL.Task) {
	tsMap := make(map[int64]*ncMDL.Task)
	for _, t := range ts {
		if t == nil {
			continue
		}
		tsMap[t.ID] = t
	}

	tasks, ok := s.TaskTypeMapCache[ty]
	if !ok || len(tasks) == 0 {
		return
	}

	res = make([]*ncMDL.Task, 0, len(tasks))
	for _, v := range tasks {
		if v == nil || v.State != 0 { //-1-删除；0-正常；1-隐藏
			continue
		}
		if _, ok := s.TaskGroupMapCache[v.GroupID]; !ok { // check taskGroup
			continue
		}

		task := *v
		t, ok := tsMap[v.ID]
		if !ok {
			task.CompleteSate = ncMDL.TaskIncomplete
		} else {
			task.CompleteSate = t.CompleteSate
		}
		res = append(res, &task)
	}

	return
}

// isHiddenTaskType check taskType is hidden
func (s *Service) isHiddenTaskType(ty int8) bool {
	tasks, ok := s.TaskTypeMapCache[ty]
	if !ok || len(tasks) == 0 {
		return true
	}
	for _, t := range tasks {
		if t == nil {
			continue
		}
		_, ok := s.TaskGroupMapCache[t.GroupID]
		if !ok {
			continue
		}
		if t.State == 0 {
			return false
		}
	}
	return true
}
