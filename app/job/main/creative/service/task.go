package service

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"time"

	"go-common/app/job/main/creative/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

//commitTask for commit task
func (s *Service) commitTask() {
	for i := 0; i < s.c.Task.TableJobNum; i++ {
		go func(i int) {
			s.dispatchData(fmt.Sprintf("%02d", i))
		}(i)
	}
}

//shardingQueueIndex sharding queue index
func (s *Service) shardingQueueIndex(name string, ql int) (i int) { ////注：使用校验和取模的原因是：使得获取的消息均匀散入不同的worker队列中
	ch := crc32.ChecksumIEEE([]byte(name))
	i = int(ch) % ql
	return
}

func (s *Service) dispatchData(index string) { //index 分表后缀名
	var id int64
	limit := s.c.Task.RowLimit
	for {
		res, err := s.newc.UserTasks(context.Background(), index, id, limit)
		if err != nil {
			log.Error("s.newc.UserTasks table index(%s)|id(%d)|limit(%d)|err(%v)", index, id, limit, err)
			return
		}
		if len(res) == 0 {
			time.Sleep(600 * time.Second)
			id = 0 //reset id
			continue
		}
		id = res[len(res)-1].ID
		tks := make([]*model.UserTask, 0, len(res))
		for _, v := range res {
			target, ok := s.TaskMapCache[v.TaskID]
			if !ok || target == nil {
				continue
			}
			//过滤 非T+1的任务
			if target.TargetType != model.TargetType004 &&
				target.TargetType != model.TargetType005 &&
				target.TargetType != model.TargetType006 &&
				target.TargetType != model.TargetType007 &&
				target.TargetType != model.TargetType008 &&
				target.TargetType != model.TargetType009 {
				continue
			}
			tks = append(tks, v)
		}
		if len(tks) > 0 {
			s.taskQueue[s.shardingQueueIndex(index, s.c.Task.TableConsumeNum)] <- tks
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Service) initTaskQueue() {
	for i := 0; i < s.c.Task.TableConsumeNum; i++ {
		ut := make(chan []*model.UserTask, s.chanSize)
		s.taskQueue[i] = ut
		go func(ch chan []*model.UserTask) {
			s.updateTaskStateByGRPC(ch)
		}(ut)
	}
}

//updateTaskStateByGRPC for check task by call grpc.
func (s *Service) updateTaskStateByGRPC(c chan []*model.UserTask) {
	for msg := range c {
		for _, v := range msg {
			mid, tid := v.MID, v.TaskID
			reply, err := s.newc.CheckTaskState(context.Background(), mid, tid)
			if err != nil {
				if ec := ecode.Cause(err); ec.Code() == ecode.ServiceUnavailable.Code() {
					log.Error("s.newc.CheckTaskState mid(%d)|task id(%d)|err(%v)", mid, tid, err)
					return
				}
				log.Warn("s.newc.CheckTaskState mid(%d)|task id(%d)|err(%v)", mid, tid, err)
				continue
			}

			log.Info("updateTaskStateByGRPC mid(%d)|task id(%d)", mid, tid)
			if reply != nil && reply.FinishState {
				_, err := s.newc.UpUserTask(context.Background(), mid, tid)
				if err != nil {
					log.Error("DriveStateByUser s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, tid, err)
					return
				}
			}

			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (s *Service) initStatViewQueue() {
	for i := 0; i < s.statViewQueueLen; i++ {
		view := make(chan *model.StatView, s.chanSize)
		s.statViewSubQueue[i] = view
		go func(m chan *model.StatView) { //播放
			for v := range m {
				log.Info("StatView v(%+v)|指标：该UID下任意avid的获得-点击量(%d)", v, v.Count)
				s.completeUserTask(s.getMIDByAID(v.ID), v.ID, model.TargetType015, v.Count)
				time.Sleep(time.Millisecond * 10)
			}
		}(view)
	}
}

func (s *Service) initStatLikeQueue() {
	for i := 0; i < s.statLikeQueueLen; i++ {
		li := make(chan *model.StatLike, s.chanSize)
		s.statLikeSubQueue[i] = li
		go func(m chan *model.StatLike) { //点赞
			for v := range m {
				log.Info("StatLike v(%+v)|指标：该UID下任意avid的获得-点赞量(%d)", v, v.Count)
				s.completeUserTask(s.getMIDByAID(v.ID), v.ID, model.TargetType020, v.Count)
				time.Sleep(time.Millisecond * 10)
			}
		}(li)
	}
}

func (s *Service) initDatabusQueue() {
	for i := 0; i < s.databusQueueLen; i++ {
		tk := make(chan *databus.Message, s.chanSize)
		sh := make(chan *model.ShareMsg, s.chanSize)
		fl := make(chan *model.Stat, s.chanSize)       //粉丝数
		rela := make(chan *model.Relation, s.chanSize) //关注
		np := make(chan *model.Up, s.chanSize)         //最新投稿
		op := make(chan *model.Up, s.chanSize)         //最新投稿
		mp := make(chan *model.Up, s.chanSize)         //手机投稿

		//单个稿件计数
		stsh := make(chan *model.StatShare, s.chanSize)
		coin := make(chan *model.StatCoin, s.chanSize)
		fav := make(chan *model.StatFav, s.chanSize)
		rep := make(chan *model.StatReply, s.chanSize)
		dm := make(chan *model.StatDM, s.chanSize)

		s.taskSubQueue[i] = tk
		s.shareSubQueue[i] = sh
		s.followerQueue[i] = fl   //粉丝
		s.relationQueue[i] = rela //关注
		s.newUpQueue[i] = np      //新投稿
		s.oldUpQueue[i] = op      //投下5个稿
		s.mobileUpQueue[i] = mp   //手机投稿

		//单个稿件计数
		s.statShareSubQueue[i] = stsh
		s.statCoinSubQueue[i] = coin
		s.statFavSubQueue[i] = fav
		s.statReplySubQueue[i] = rep
		s.statDMSubQueue[i] = dm

		//水印设置、观看创作学院视频、参加激励计划
		go func(m chan *databus.Message) {
			s.startByTask(m)
		}(tk)
		//分享自己的稿件
		go func(m chan *model.ShareMsg) {
			for v := range m {
				log.Info("startByShare mid(%d)|v(%+v)|指标：该UID分享自己视频的次数≥1", v.MID, v)
				s.completeUserTask(v.MID, v.OID, model.TargetType002, 1)
				time.Sleep(time.Millisecond * 10)
			}
		}(sh)

		//关注哔哩哔哩创作中心和粉丝数判断
		go func(m chan *model.Stat) {
			for v := range m {
				log.Info("followerStat mid(%d)|v(%+v)|指标：该UID的粉丝数≥10 或者 1000", v.MID, v)
				s.completeUserTask(v.MID, 0, model.TargetType010, v.Follower) //新手任务粉丝数
				s.completeUserTask(v.MID, 0, model.TargetType022, v.Follower) //进阶任务粉丝数
			}
		}(fl)
		go func(m chan *model.Relation) {
			for v := range m {
				log.Info("relationMID mid(%d)|v(%+v)|指标：该UID的关注列表含有“哔哩哔哩创作中心", v.MID, v)
				s.completeUserTask(v.MID, 0, model.TargetType012, 1)
			}
		}(rela)

		// 该UID下开放浏览的稿件≥1
		go func(m chan *model.Up) {
			for v := range m {
				log.Info("newUP mid(%d)|v(%+v)|指标：该UID下开放浏览的稿件≥1", v.MID, v)
				s.completeUserTask(v.MID, v.AID, model.TargetType001, 1)
			}
		}(np)

		// 该UID下开放浏览的稿件≥5
		go func(m chan *model.Up) {
			for v := range m {
				log.Info("oldUP mid(%d)|v(%+v)|指标：该UID下开放浏览的稿件≥5", v.MID, v)
				s.completeUserTask(v.MID, v.AID, model.TargetType014, 1)
			}
		}(op)

		//单个稿件计数
		go func(m chan *model.StatReply) {
			for v := range m {
				log.Info("StatReply v(%+v)|指标：该UID下任意avid的获得-评论量(%d)", v, v.Count)
				s.completeUserTask(s.getMIDByAID(v.ID), v.ID, model.TargetType016, v.Count)
				time.Sleep(time.Millisecond * 10)
			}
		}(rep)

		go func(m chan *model.StatShare) {
			for v := range m {
				log.Info("StatShare v(%+v)|指标：该UID下任意avid的获得-分享量(%d)", v, v.Count)
				s.completeUserTask(s.getMIDByAID(v.ID), v.ID, model.TargetType017, v.Count)
				time.Sleep(time.Millisecond * 10)
			}
		}(stsh)

		go func(m chan *model.StatFav) {
			for v := range m {
				log.Info("StatFav v(%+v)|指标：该UID下任意avid的获得-收藏量(%d)", v, v.Count)
				s.completeUserTask(s.getMIDByAID(v.ID), v.ID, model.TargetType018, v.Count)
				time.Sleep(time.Millisecond * 10)
			}

		}(fav)

		go func(m chan *model.StatCoin) {
			for v := range m {
				log.Info("StatCoin v(%+v)|指标：该UID下任意avid的获得-硬币量(%d)", v, v.Count)
				s.completeUserTask(s.getMIDByAID(v.ID), v.ID, model.TargetType019, v.Count)
				time.Sleep(time.Millisecond * 10)
			}
		}(coin)

		go func(m chan *model.StatDM) {
			for v := range m {
				log.Info("StatDM v(%+v)|指标：该UID下任意avid的获得-弹幕量(%d)", v, v.Count)
				s.completeUserTask(s.getMIDByAID(v.ID), v.ID, model.TargetType021, v.Count)
				time.Sleep(time.Millisecond * 10)
			}
		}(dm)

		go func(m chan *model.Up) {
			for v := range m {
				log.Info("Mobile mid(%d)|v(%+v)|指标：该UID通过手机投稿的稿件≥1", v.MID, v)
				s.completeUserTask(v.MID, v.AID, model.TargetType013, 1)
				time.Sleep(time.Millisecond * 10)
			}
		}(mp)
	}
}

func (s *Service) startByTask(c chan *databus.Message) {
	for msg := range c {
		v := &model.TaskMsg{}
		if err := json.Unmarshal(msg.Value, v); err != nil {
			log.Error("startByTask json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		switch v.From {
		case model.MsgForWaterMark:
			log.Info("startByTask WaterMark mid(%d)|v(%+v)|指标：任务完成期间该UID的水印开关为打开状态", v.MID, v)
			s.completeUserTask(v.MID, 0, model.TargetType011, v.Count)

		case model.MsgForAcademyFavVideo:
			log.Info("startByTask AcademyFavVideo mid(%d)|v(%+v)|指标：该UID在创作学院的观看记录≥1", v.MID, v)
			s.completeUserTask(v.MID, 0, model.TargetType003, v.Count)

		case model.MsgForGrowAccount:
			log.Info("startByTask GrowAccount mid(%d)|v(%+v)|指标：该UID的激励计划状态为已开通", v.MID, v)
			s.completeUserTask(v.MID, 0, model.TargetType023, v.Count)

		case model.MsgForOpenFansMedal:
			log.Info("startByTask OpenFansMedal mid(%d)|v(%+v)|指标：该UID粉丝勋章为开启状态", v.MID, v)
			s.completeUserTask(v.MID, 0, model.TargetType024, v.Count)
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func (s *Service) getMIDByAID(aid int64) (mid int64) {
	arc, err := s.arc.Archive(context.Background(), aid)
	if err != nil || arc == nil {
		log.Error("getMIDByAID s.arc.Archive aid(%d)|err(%v)", aid, err)
		return
	}
	mid = arc.Author.Mid
	return
}

func (s *Service) getTaskIDByTargetType(mid int64, ty int8) (tid int64) {
	userTasks, err := s.newc.UserTasksByMIDAndState(context.Background(), mid, model.TaskIncomplete)
	if err != nil {
		log.Error("s.newc.UserTasksByMIDAndState mid(%d)|err(%v)", mid, err)
		return
	}
	if len(userTasks) == 0 {
		return
	}

	for _, v := range userTasks {
		t, ok := s.TaskMapCache[v.TaskID]
		if ok && t.TargetType == ty {
			tid = v.TaskID
			break
		}
	}
	return
}

// completeUserTask update user task to complete
func (s *Service) completeUserTask(mid, aid int64, ty int8, count int64) {
	tid := s.getTaskIDByTargetType(mid, ty)
	if tid == 0 {
		return
	}
	target, ok := s.TaskMapCache[tid]
	if !ok || target == nil {
		return
	}

	if count >= target.TargetValue {
		if _, err := s.newc.UpUserTask(context.Background(), mid, tid); err != nil {
			log.Error("s.newc.UpUserTask mid(%d)|tid(%d)|err(%v)", mid, tid, err)
			return
		}
		log.Info("completeUserTask mid(%d)|aid(%d)|count(%d)|taskID(%d)|targetType(%d)|targetValue(%d)", mid, aid, count, tid, ty, target.TargetValue)
	}
}

// 对第30天未完成新手任务的UP主，发送消息通知；记录时间点为用户加入任务成就的时间;该消息有且仅发送一次。
func (s *Service) expireTaskNotify() {
	for i := 0; i < s.c.Task.TaskTableJobNum; i++ {
		go func(i int) {
			s.dispatchTasksNotify(fmt.Sprintf("%02d", i))
		}(i)
	}
}

func (s *Service) dispatchTasksNotify(index string) {
	var id int64
	ext := s.c.Task.TaskExpireTime
	th := s.c.Task.TaskSendHour
	tm := s.c.Task.TaskSendMiniute
	ts := s.c.Task.TaskSendSecond

	limit := s.c.Task.TaskRowLimitNum
	batchSize := s.c.Task.TaskBatchMidNum //每次发送mid数量
	for {
		now := time.Now()
		if now.Hour() != th || now.Minute() != tm || now.Second() != ts {
			// log.Info("dispatchTasksNotify minuts(%d) second(%d)", now.Minute(), now.Second())
			time.Sleep(1 * time.Second)
			continue
		}
		ctime := now.Unix() - ext //检查任务是否超过30天未完成
		year, month, day := time.Unix(ctime, 0).Date()
		start := time.Date(year, month, day, 0, 0, 0, 0, time.Local).Format("2006-01-02 15:04:05")
		end := time.Date(year, month, day, 23, 59, 59, 999, time.Local).Format("2006-01-02 15:04:05")
		log.Info("dispatchTasksNotify now(%s)|start(%s)|end(%s)", now.Format("2006-01-02 15:04:05"), start, end)

		midMap := make(map[int64]*model.UserTask)
		for {
			res, err := s.newc.UserTasksNotify(context.Background(), index, id, start, end, limit)
			if err != nil {
				log.Error("s.newc.UserTasksNotify table index(%s)|start(%s)|end(%s)|limit(%d)|err(%v)", index, start, end, limit, err)
				return
			}

			if len(res) == 0 {
				id = 0
				break
			}

			for _, v := range res {
				midMap[v.MID] = v
			}

			id = res[len(res)-1].ID //next limit
			time.Sleep(1 * time.Second)
		}

		if len(midMap) == 0 {
			continue
		}
		mids := make([]int64, 0, len(midMap))
		for mid := range midMap {
			mids = append(mids, mid)
		}

		var tmids []int64
		count := len(mids)/batchSize + 1
		for i := 0; i < count; i++ {
			if i == count-1 {
				tmids = mids[i*batchSize:]
			} else {
				tmids = mids[i*batchSize : (i+1)*batchSize]
			}
			if len(tmids) > 0 {
				s.taskNotifyQueue[s.shardingQueueIndex(index, s.c.Task.TaskTableConsumeNum)] <- tmids
			}
		}
	}
}

func (s *Service) initTaskNotifyQueue() {
	for i := 0; i < s.c.Task.TaskTableConsumeNum; i++ {
		ut := make(chan []int64, s.chanSize)
		s.taskNotifyQueue[i] = ut
		go func(ch chan []int64) {
			s.sendTaskNotify(ch)
		}(ut)
	}
}

func (s *Service) sendTaskNotify(c chan []int64) {
	for mids := range c {
		if len(mids) == 0 {
			time.Sleep(time.Second * 60)
			continue
		}

		for i := 1; i < 3; i++ {
			if err := s.newc.SendNotify(context.Background(), mids, s.c.Task.TaskMsgCode, s.c.Task.TaskTitle, s.c.Task.TaskContent); err != nil {
				log.Error("sendTaskNotify s.newc.SendNotify(%v) error(%v)", mids, err)
				time.Sleep(time.Millisecond * 10)
				continue
			} else {
				log.Info("sendTaskNotify s.newc.SendNotify mids(%+v)", mids)
				break
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// loadproc
func (s *Service) loadProc() {
	for {
		time.Sleep(3 * time.Minute)
		s.loadTasks()
		s.loadGiftRewards()
	}
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
	temp := make(map[int64]*model.Task)
	for _, v := range s.TaskCache {
		temp[v.ID] = v
	}
	s.TaskMapCache = temp
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
	s.GiftRewardCache = res
}

// 检查用户是否有奖励可领取
func (s *Service) checkRewardReceive(c context.Context, mid int64) (is bool) {
	tasks, err := s.newc.UserTasksByMID(c, mid)
	if err != nil {
		log.Error("s.newc.UserTasksByMID mid(%v)|error(%v)", mid, err)
		return
	}
	groupMap := make(map[int64][]*model.UserTask)
	giftMap := make(map[int8][]*model.UserTask)
	for _, v := range tasks {
		groupMap[v.TaskGroupID] = append(groupMap[v.TaskGroupID], v)
		if _, ok := s.GiftRewardCache[v.TaskType]; ok {
			giftMap[v.TaskType] = append(giftMap[v.TaskType], v)
		}
	}

	groupNum := 0                  // 组奖励 已完成个数
	giftNum := make(map[int8]bool) // 礼包奖励 已完成个数
	for _, ts := range groupMap {
		for _, t := range ts {
			if t.State == model.TaskIncomplete {
				groupNum++
				if _, ok := s.GiftRewardCache[t.TaskType]; ok {
					giftNum[t.TaskType] = true
				}
				break
			}
		}
	}

	r1, err := s.newc.BaseRewardCount(c, mid) // 组奖励 已领取个数
	if err != nil {
		log.Error("s.newc.BaseRewardCount mid(%v)|error(%v)", mid, err)
		return
	}
	r2, err := s.newc.GiftRewardCount(c, mid) // 礼包奖励 已领取个数
	if err != nil {
		log.Error("s.newc.GiftRewardCount mid(%v)|error(%v)", mid, err)
		return
	}

	total := len(groupMap) + len(giftMap) //奖励总数
	untotal := groupNum + len(giftNum)    //未完成的奖励
	receive := r1 + r2                    //已领取奖励

	// 可领取的奖励 = 奖励总数 -未完成的奖励 - 已领取奖励
	count := total - untotal - receive

	log.Info("checkRewardReceive mid(%d)|奖励总数(%d)|未完成奖励总数(%d)|已领取奖励总数(%d)|可领取奖励总数(%d)", mid, total, untotal, receive, count)
	if count > 0 {
		is = true
	}
	return
}

// 该消息每周最多发送 1 条，发送时间为每周六的20:00，用户为上周周六18:00 - 本周周六17:59所有达到领取奖励且 未领取 的用户。
// 通知仅限用户有未领取的奖励时发送：若在该时间段，用户已领取全部可领取的奖励，
// 则不发送通知，如果用户已领取部分可领取的奖励，仍有部分奖励未领取，则仍然发送通知
func (s *Service) rewardReceiveNotify() {
	for i := 0; i < s.c.Task.RewardTableJobNum; i++ {
		go func(i int) {
			s.dispatchRewardNotify(fmt.Sprintf("%02d", i))
		}(i)
	}
}

func (s *Service) dispatchRewardNotify(index string) {
	var id int64
	week := s.c.Task.RewardWeek      //星期几
	ld := s.c.Task.RewardLastDay     //从过去多少天开始查询
	lh := s.c.Task.RewardLastHour    //几点开始查询
	lm := s.c.Task.RewardLastMiniute //几分开始查询
	ls := s.c.Task.RewardLastSecond  //几秒开始查询

	nh := s.c.Task.RewardNowHour    //从当前时间几点开始
	nm := s.c.Task.RewardNowMiniute //从当前时间几分开始
	ns := s.c.Task.RewardNowSecond  //从当前时间几秒开始

	limit := s.c.Task.RewardRowLimitNum
	batchSize := s.c.Task.RewardBatchMidNum //每次发送mid数量
	for {
		now := time.Now()
		if int(now.Weekday()) != week || now.Hour() != nh || now.Minute() != nm || now.Second() != ns {
			// log.Info("dispatchRewardNotify Weekday(%d) Hour(%d) Minute(%d) Second(%d)", now.Weekday(), now.Hour(), now.Minute(), now.Second())
			time.Sleep(1 * time.Second)
			continue
		}

		last := now.AddDate(0, 0, ld).Add(time.Hour * time.Duration(lh)).Add(time.Minute * time.Duration(lm)).Add(time.Second * time.Duration(ls))
		log.Info("dispatchRewardNotify last(%s) now(%s)\n", last.Format("2006-01-02 15:04:05"), now.Format("2006-01-02 15:04:05"))
		midMap := make(map[int64]*model.UserTask)
		for {
			res, err := s.newc.CheckTasksForRewardNotify(context.Background(), index, id, last, now, limit)
			if err != nil {
				log.Error("s.newc.CheckTasksForRewardNotify table index(%s)|id(%d)|limit(%d)|err(%v)", index, id, limit, err)
				return
			}

			if len(res) == 0 {
				id = 0
				break
			}
			for _, v := range res {
				midMap[v.MID] = v
			}

			id = res[len(res)-1].ID //next limit
			time.Sleep(1 * time.Second)
		}

		if len(midMap) == 0 {
			continue
		}
		mids := make([]int64, 0, len(midMap))
		for mid := range midMap {
			if s.checkRewardReceive(context.Background(), mid) {
				mids = append(mids, mid)
			}
		}

		var tmids []int64
		count := len(mids)/batchSize + 1
		for i := 0; i < count; i++ {
			if i == count-1 {
				tmids = mids[i*batchSize:]
			} else {
				tmids = mids[i*batchSize : (i+1)*batchSize]
			}
			if len(tmids) > 0 {
				s.rewardNotifyQueue[s.shardingQueueIndex(index, s.c.Task.RewardTableConsumeNum)] <- tmids
			}
		}
	}
}

func (s *Service) initRewardNotifyQueue() {
	for i := 0; i < s.c.Task.RewardTableConsumeNum; i++ {
		ut := make(chan []int64, s.chanSize)
		s.rewardNotifyQueue[i] = ut
		go func(ch chan []int64) {
			s.sendRewardNotify(ch)
		}(ut)
	}
}

func (s *Service) sendRewardNotify(c chan []int64) {
	for mids := range c {
		if len(mids) == 0 {
			time.Sleep(time.Second * 1)
			continue
		}

		for i := 1; i < 3; i++ {
			if err := s.newc.SendNotify(context.Background(), mids, s.c.Task.RewardMsgCode, s.c.Task.RewardTitle, s.c.Task.RewardContent); err != nil {
				log.Error("sendRewardNotify s.newc.SendNotify mids(%+v) error(%v)", mids, err)
				time.Sleep(time.Millisecond * 100)
				continue
			} else {
				log.Info("sendRewardNotify s.newc.SendNotify mids(%+v)", mids)
				break
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}
