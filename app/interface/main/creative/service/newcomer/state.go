package newcomer

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/model/newcomer"
	"go-common/library/log"

	"go-common/library/sync/errgroup.v2"
)

//DriveStateByUser drive new hand & advanced task state
func (s *Service) DriveStateByUser(c context.Context, mid int64) {
	res, err := s.newc.UserTasksByMID(c, mid)
	if err != nil {
		log.Error("DriveStateByUser s.newc.UserTasksByMID mid(%d)|err(%v)", mid, err)
		return
	}
	if len(res) == 0 {
		return
	}

	for _, v := range res {
		if s.CheckTaskState(context.Background(), &newcomer.CheckTaskStateReq{MID: v.MID, TaskID: v.TaskID}) {
			_, err := s.newc.UpUserTask(c, v.MID, v.TaskID)
			if err != nil {
				log.Error("DriveStateByUser s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", v.MID, v.TaskID, err)
			}
		}
	}
}

//CheckTaskState check task state by mid & task id
func (s *Service) CheckTaskState(c context.Context, req *newcomer.CheckTaskStateReq) (state bool) {
	if req == nil {
		return
	}

	mid, tid := req.MID, req.TaskID
	if _, ok := s.TaskMapCache[tid]; !ok {
		return
	}
	task := s.TaskMapCache[tid]
	//log.Info("CheckTaskState req(%+v)|task(%+v)", req, task)

	switch task.TargetType {
	case newcomer.TargetType001, newcomer.TargetType014: //新手任务-1 该UID下开放浏览的稿件≥1 / 进阶任务- 14 该UID下开放浏览的稿件≥5
		cnt, err := s.arc.UpCount(c, mid)
		if err != nil {
			log.Error("CheckTaskState s.arc.UpCount mid(%d)|error(%v)", mid, err)
			return
		}
		if cnt >= task.TargetValue {
			switch task.TargetType {
			case newcomer.TargetType001:
				log.Info("CheckTaskState TargetType001 mid(%d)|count(%d)|指标：新手任务-1 该UID下开放浏览的稿件≥1", mid, cnt)
			case newcomer.TargetType014:
				log.Info("CheckTaskState TargetType014 mid(%d)|count(%d)|指标：进阶任务- 14 该UID下开放浏览的稿件≥5", mid, cnt)
			}
			state = true
		}

	case newcomer.TargetType002: //该UID分享自己视频的次数≥1
	case newcomer.TargetType003: //该UID在创作学院的观看记录≥1
		cnt, err := s.aca.PlayCount(c, mid)
		if err != nil {
			log.Error("s.aca.PlayCount error(%v)", err)
			return
		}
		if cnt >= task.TargetValue {
			log.Info("CheckTaskState TargetType003 mid(%d)|count(%d)|指标：新手任务-1 该UID在创作学院的观看记录≥1", mid, cnt)
			state = true
		}

	case newcomer.TargetType004, newcomer.TargetType005, newcomer.TargetType006, newcomer.TargetType007, newcomer.TargetType008, newcomer.TargetType009: //该UID下所有avid的获得评论数/...≥3/5
		state = s.totalStat(c, mid, int64(task.TargetValue), task.TargetType)

	case newcomer.TargetType010, newcomer.TargetType022: //新手任务- 10 获得10个粉丝 / 进阶任务- 22 该UID的粉丝数≥1000
		pl, err := s.acc.ProfileWithStat(c, mid)
		if err != nil {
			log.Error("CheckTaskState s.acc.ProfileWithStat error(%v)", err)
			return
		}
		if pl == nil {
			return
		}
		if pl.Follower >= int64(task.TargetValue) {
			switch task.TargetType {
			case newcomer.TargetType010:
				log.Info("CheckTaskState TargetType010 mid(%d)|follower(%d)|指标：新手任务-10 该UID的粉丝数≥10", mid, pl.Follower)
			case newcomer.TargetType022:
				log.Info("CheckTaskState TargetType022 mid(%d)|follower(%d)|指标：进阶任务-22 该UID的粉丝数≥1000", mid, pl.Follower)
			}
			state = true
		}

	case newcomer.TargetType011: //任务完成期间该UID的水印开关为打开状态
		wm, err := s.wm.WaterMark(c, mid)
		if err != nil {
			log.Error("CheckTaskState s.wm.WaterMark mid(%d) error(%v)", mid, err)
			return
		}
		if wm != nil && wm.State == 1 && wm.URL != "" {
			log.Info("CheckTaskState TargetType011 mid(%d)已开启水印|指标：新手任务-11 任务完成期间该UID的水印开关为打开状态", mid)
			state = true
		}

	case newcomer.TargetType012: //该UID的关注列表含有“哔哩哔哩创作中心”
		fid := creatorMID //用户哔哩哔哩创作中心 mid
		fl, err := s.acc.Relations(c, mid, []int64{fid}, "")
		if err != nil {
			log.Error("CheckTaskState s.acc.Relations mid(%d)|ip(%s)|error(%v)", mid, "", err)
			return
		}
		if fl == nil {
			return
		}
		if st, ok := fl[fid]; ok && (st == 6 || st == 2) {
			log.Info("CheckTaskState TargetType012 mid(%d)已关注创作中心|指标：新手任务-12 该UID的关注列表含有“哔哩哔哩创作中心”", mid)
			state = true
		}

	case newcomer.TargetType013, newcomer.TargetType015, newcomer.TargetType016, newcomer.TargetType017, newcomer.TargetType018, newcomer.TargetType019, newcomer.TargetType020, newcomer.TargetType021: //13-该UID通过手机投稿的稿件≥1 / 15-21单稿件获得1000播放...
		state = s.singleStat(c, mid, int64(task.TargetValue), task.TargetType)
	case newcomer.TargetType023: //该UID的激励计划状态为已开通
		ac, err := s.order.GrowAccountState(c, mid, 0) //类型 0 视频 2 专栏 3 素材
		if err != nil {
			log.Error("CheckTaskState s.order.GrowUpPlan mid(%d) error(%v)", mid, err)
			return
		}
		if ac == nil {
			return
		}
		if ac.State == 3 { //账号状态; 1: 未申请; 2: 待审核; 3: 已签约; 4.已驳回; 5.主动退出; 6:被动退出; 7:封禁
			log.Info("CheckTaskState TargetType023 mid(%d)已签约激励计划|指标：进阶任务-23 该UID的激励计划状态为已开通", mid)
			state = true
		}

	case newcomer.TargetType024: //该UID粉丝勋章为开启状态
		medal, err := s.medal.Medal(c, mid)
		if err != nil {
			log.Error("CheckTaskState s.medal.Medal mid(%d) error(%v)", mid, err)
			return
		}
		if medal == nil {
			return
		}
		status, err := strconv.Atoi(medal.Status)
		if err != nil {
			log.Error("CheckTaskState strconv.Atoi medal.Status(%s) error(%v)", medal.Status, err)
			return
		}
		if status == 2 { //勋章审核状态 -1已拒绝 0未申请 1已申请 2已开通
			log.Info("CheckTaskState TargetType024 mid(%d)已开通粉丝勋章|指标：进阶任务-24 该UID粉丝勋章为开启状态", mid)
			state = true
		}
	}
	return
}

func (s *Service) totalStat(c context.Context, mid, val int64, ty int8) (state bool) {
	st, err := s.data.UpStat(c, mid, time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102"))
	if err != nil || st == nil {
		log.Error("totalStat s.data.NewStat mid(%d) st(%+v) error(%v) ", mid, st, err)
		return
	}

	switch ty {
	case newcomer.TargetType004:
		if st.Reply >= val {
			log.Info("CheckTaskState TargetType004 mid(%d)|v(%+v)|指标：新手任务-4 该UID下所有avid的获得评论数≥3", mid, st.Reply)
			state = true
		}
	case newcomer.TargetType005:
		if st.Share >= val {
			log.Info("CheckTaskState TargetType005 mid(%d)|v(%+v)|指标：新手任务-5 该UID下所有avid获得分享数≥3", mid, st.Share)
			state = true
		}
	case newcomer.TargetType006:
		if st.Fav >= val {
			log.Info("CheckTaskState TargetType006 mid(%d)|v(%+v)|指标：新手任务-6 该UID的所有avid的获得收藏数≥5", mid, st.Fav)
			state = true
		}
	case newcomer.TargetType007:
		if st.Coin >= val {
			log.Info("CheckTaskState TargetType007 mid(%d)|v(%+v)|指标：新手任务-7 该UID下所有avid的获得硬币数≥5", mid, st.Coin)
			state = true
		}
	case newcomer.TargetType008:
		if st.Like >= val {
			log.Info("CheckTaskState TargetType008 mid(%d)|v(%+v)|指标：新手任务-8 该UID下所有avid获得点赞数≥5", mid, st.Like)
			state = true
		}
	case newcomer.TargetType009:
		if st.Dm >= val {
			log.Info("CheckTaskState TargetType009 mid(%d)|v(%+v)|指标：新手任务-9 该UID下所有avid的获得弹幕数≥5", mid, st.Dm)
			state = true
		}
	}

	return
}

func (s *Service) singleStat(c context.Context, mid, val int64, ty int8) (state bool) {
	st, err := s.data.UpArchiveStatQuery(c, mid, time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102"))
	if err != nil || st == nil {
		log.Error("singleStat s.data.UpArchiveStatQuery mid(%d)|err(%v)", mid, err)
		return
	}

	switch ty {
	case newcomer.TargetType013:
		if st.FromPhoneNum >= val {
			log.Info("CheckTaskState TargetType013 mid(%d)|v(%+v)|指标：进阶任务-13 用手机投稿上传视频", mid, st.FromPhoneNum)
			state = true
		}
	case newcomer.TargetType015:
		if st.PlayV >= val {
			log.Info("CheckTaskState TargetType015 mid(%d)|v(%+v)|指标：进阶任务-15 该UID下任意avid的获得点击量≥1000", mid, st.PlayV)
			state = true
		}
	case newcomer.TargetType016:
		if st.ReplyV >= val {
			log.Info("CheckTaskState TargetType016 mid(%d)|v(%+v)|指标：进阶任务-16 该UID下任意avid的评论≥30", mid, st.ReplyV)
			state = true
		}
	case newcomer.TargetType017:
		if st.ShareV >= val {
			log.Info("CheckTaskState TargetType017 mid(%d)|v(%+v)|指标：进阶任务-17 该UID下任意avid的获得分享数≥10", mid, st.ShareV)
			state = true
		}
	case newcomer.TargetType018:
		if st.FavV >= val {
			log.Info("CheckTaskState TargetType018 mid(%d)|v(%+v)|指标：进阶任务-18 该UID下任意avid的获得收藏数≥30", mid, st.FavV)
			state = true
		}
	case newcomer.TargetType019:
		if st.CoinV >= val {
			log.Info("CheckTaskState TargetType019 mid(%d)|v(%+v)|指标：进阶任务-19 该UID下任意avid的获得硬币数≥50", mid, st.CoinV)
			state = true
		}
	case newcomer.TargetType020:
		if st.LikeV >= val {
			log.Info("CheckTaskState TargetType020 mid(%d)|v(%+v)|指标：进阶任务-20 该UID下任意avid的获得点赞数≥50", mid, st.LikeV)
			state = true
		}
	case newcomer.TargetType021:
		if st.DmV >= val {
			log.Info("CheckTaskState TargetType021 mid(%d)|v(%+v)|指标：进阶任务-21 该UID下任意avid的获得弹幕数≥50", mid, st.DmV)
			state = true
		}
	}
	return
}

//syncCheckTaskStatus check task status
func (s *Service) syncCheckTaskStatus(c context.Context, mid int64, tasks []*newcomer.Task) {
	log.Info("syncCheckTaskStatus mid(%d) | tasks count(%d)", mid, len(tasks))
	tsm := getTaskSortMap(tasks)
	if len(tsm) == 0 {
		return
	}
	g := &errgroup.Group{}
	for k := range tsm {
		switch k {
		case newcomer.ArcUpCount: //该UID下开放浏览的稿件数量
			g.Go(func(context.Context) error {
				s.arcUpCount(c, mid, tsm[newcomer.ArcUpCount])
				return nil
			})
		case newcomer.AcaPlayCount: //该UID在创作学院的观看记录
			g.Go(func(context.Context) error {
				s.acaPlayCount(c, mid, tsm[newcomer.AcaPlayCount])
				return nil
			})
		case newcomer.DataUpStat: //该UID下所有avid的最高计数
			g.Go(func(context.Context) error {
				s.dataUpStat(c, mid, tsm[newcomer.DataUpStat])
				return nil
			})
		case newcomer.AccProfileWithStat: //粉丝数量
			g.Go(func(context.Context) error {
				s.accProfileWithStat(c, mid, tsm[newcomer.AccProfileWithStat])
				return nil
			})
		case newcomer.WmWaterMark: //水印状态
			g.Go(func(context.Context) error {
				s.wmWaterMark(c, mid, tsm[newcomer.WmWaterMark])
				return nil
			})
		case newcomer.AccRelation: //该UID的关注列表含有“哔哩哔哩创作中心”
			g.Go(func(context.Context) error {
				s.accRelation(c, mid, tsm[newcomer.AccRelation])
				return nil
			})
		case newcomer.DataUpArchiveStat: //该UID下任意avid的计数
			g.Go(func(context.Context) error {
				s.dataUpArchiveStatQuery(c, mid, tsm[newcomer.DataUpArchiveStat])
				return nil
			})
		case newcomer.OrderGrowAccountState: //激励计划状态
			g.Go(func(context.Context) error {
				s.orderGrowAccountState(c, mid, tsm[newcomer.OrderGrowAccountState])
				return nil
			})
		case newcomer.MedalCheckMedal: //该UID粉丝勋章
			g.Go(func(context.Context) error {
				s.medalCheckMedal(c, mid, tsm[newcomer.MedalCheckMedal])
				return nil
			})
		}
	}
	g.Wait()
}

func (s *Service) arcUpCount(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	cnt, err := s.arc.UpCount(c, mid)
	if err != nil {
		log.Error("arcUpCount s.arc.UpCount mid(%d)|error(%v)", mid, err)
		return
	}
	for _, v := range tasks {
		if cnt >= v.TargetValue {
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("arcUpCount s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
			}
			log.Info("arcUpCount finish task mid(%d)|count(%d)|", mid, cnt)
		}
	}
}

func (s *Service) acaPlayCount(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	cnt, err := s.aca.PlayCount(c, mid)
	if err != nil {
		log.Error("acaPlayCount s.aca.PlayCount error(%v)", err)
		return
	}
	for _, v := range tasks {
		if cnt >= v.TargetValue {
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("acaPlayCount s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
				continue
			}
			log.Info("acaPlayCount finish task mid(%d)|count(%d)|", mid, cnt)
		}
	}
}

func (s *Service) dataUpStat(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	st, err := s.data.UpStat(c, mid, time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102"))
	if err != nil || st == nil {
		log.Error("dataUpStat s.data.NewStat mid(%d) st(%+v) error(%v) ", mid, st, err)
		return
	}

	for _, v := range tasks {
		val := int64(v.TargetValue)
		state := false
		switch v.TargetType {
		case newcomer.TargetType004:
			if st.Reply >= val {
				state = true
			}
			log.Info("dataUpStat TargetType004 mid(%d)|v(%+v)|", mid, st.Reply)
		case newcomer.TargetType005:
			if st.Share >= val {
				state = true
			}
			log.Info("dataUpStat TargetType005 mid(%d)|v(%+v)|", mid, st.Share)
		case newcomer.TargetType006:
			if st.Fav >= val {
				state = true
			}
			log.Info("dataUpStat TargetType006 mid(%d)|v(%+v)|", mid, st.Fav)
		case newcomer.TargetType007:
			if st.Coin >= val {
				state = true
			}
			log.Info("dataUpStat TargetType007 mid(%d)|v(%+v)|", mid, st.Coin)
		case newcomer.TargetType008:
			if st.Like >= val {
				state = true
			}
			log.Info("dataUpStat TargetType008 mid(%d)|v(%+v)|", mid, st.Like)
		case newcomer.TargetType009:
			if st.Dm >= val {
				state = true
			}
			log.Info("dataUpStat TargetType009 mid(%d)|v(%+v)|", mid, st.Dm)
		}
		if state {
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("dataUpStat s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
			}
		}
	}
}

func (s *Service) accProfileWithStat(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	pl, err := s.acc.ProfileWithStat(c, mid)
	if err != nil {
		log.Error("accProfileWithStat s.acc.ProfileWithStat error(%v)", err)
		return
	}
	if pl == nil {
		return
	}
	for _, v := range tasks {
		if pl.Follower >= int64(v.TargetValue) {
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("accProfileWithStat s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
				continue
			}
			log.Info("accProfileWithStat finish task mid(%d)|follower(%d)|", mid, pl.Follower)
		}
	}
}

func (s *Service) wmWaterMark(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	wm, err := s.wm.WaterMark(c, mid)
	if err != nil {
		log.Error("wmWaterMark s.wm.WaterMark mid(%d) error(%v)", mid, err)
		return
	}
	for _, v := range tasks {
		if wm != nil && wm.State == 1 && wm.URL != "" {
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("wmWaterMark s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
				continue
			}
			log.Info("wmWaterMark finish task mid(%d)|state(%d)|", mid, wm.State)
		}
	}
}

func (s *Service) accRelation(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	fid := creatorMID //用户哔哩哔哩创作中心 mid
	fl, err := s.acc.Relations(c, mid, []int64{fid}, "")
	if err != nil {
		log.Error("accRelation s.acc.Relations mid(%d)|ip(%s)|error(%v)", mid, "", err)
		return
	}
	if fl == nil {
		return
	}
	for _, v := range tasks {
		if st, ok := fl[fid]; ok && (st == 6 || st == 2) {
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("accRelation s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
				continue
			}
			log.Info("accRelation finish task mid(%d)|state(%d)|", mid, st)
		}
	}
}

func (s *Service) dataUpArchiveStatQuery(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	st, err := s.data.UpArchiveStatQuery(c, mid, time.Now().AddDate(0, 0, -1).Add(-12*time.Hour).Format("20060102"))
	if err != nil || st == nil {
		log.Error("dataUpArchiveStatQuery s.data.UpArchiveStatQuery mid(%d)|err(%v)", mid, err)
		return
	}

	for _, v := range tasks {
		val := int64(v.TargetValue)
		state := false
		switch v.TargetType {
		case newcomer.TargetType013:
			if st.FromPhoneNum >= val {
				log.Info("dataUpArchiveStatQuery TargetType013 mid(%d)|v(%+v)|", mid, st.FromPhoneNum)
				state = true
			}
		case newcomer.TargetType015:
			if st.PlayV >= val {
				log.Info("dataUpArchiveStatQuery TargetType015 mid(%d)|v(%+v)|", mid, st.PlayV)
				state = true
			}
		case newcomer.TargetType016:
			if st.ReplyV >= val {
				log.Info("dataUpArchiveStatQuery TargetType016 mid(%d)|v(%+v)|", mid, st.ReplyV)
				state = true
			}
		case newcomer.TargetType017:
			if st.ShareV >= val {
				log.Info("dataUpArchiveStatQuery TargetType017 mid(%d)|v(%+v)|", mid, st.ShareV)
				state = true
			}
		case newcomer.TargetType018:
			if st.FavV >= val {
				log.Info("dataUpArchiveStatQuery TargetType018 mid(%d)|v(%+v)|", mid, st.FavV)
				state = true
			}
		case newcomer.TargetType019:
			if st.CoinV >= val {
				log.Info("dataUpArchiveStatQuery TargetType019 mid(%d)|v(%+v)|", mid, st.CoinV)
				state = true
			}
		case newcomer.TargetType020:
			if st.LikeV >= val {
				log.Info("dataUpArchiveStatQuery TargetType020 mid(%d)|v(%+v)|", mid, st.LikeV)
				state = true
			}
		case newcomer.TargetType021:
			if st.DmV >= val {
				log.Info("dataUpArchiveStatQuery TargetType021 mid(%d)|v(%+v)|", mid, st.DmV)
				state = true
			}
		}
		if state {
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("dataUpArchiveStatQuery s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
			}
		}
	}
}

func (s *Service) orderGrowAccountState(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	ac, err := s.order.GrowAccountState(c, mid, 0)
	if err != nil {
		log.Error("orderGrowAccountState s.order.GrowUpPlan mid(%d) error(%v)", mid, err)
		return
	}
	if ac == nil {
		return
	}
	for _, v := range tasks {
		if ac.State == 3 { //账号状态; 1: 未申请; 2: 待审核; 3: 已签约; 4.已驳回; 5.主动退出; 6:被动退出; 7:封禁
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("orderGrowAccountState s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
				continue
			}
			log.Info("orderGrowAccountState finish task mid(%d)|state(%d)|", mid, ac.State)
		}
	}
}

func (s *Service) medalCheckMedal(c context.Context, mid int64, tasks []*newcomer.Task) {
	if len(tasks) == 0 {
		return
	}
	medal, err := s.medal.Medal(c, mid)
	if err != nil {
		log.Error("medalCheckMedal s.medal.Medal mid(%d) error(%v)", mid, err)
		return
	}
	if medal == nil {
		return
	}
	for _, v := range tasks {
		status, err := strconv.Atoi(medal.Status)
		if err != nil {
			log.Error("medalCheckMedal strconv.Atoi medal.Status(%s) error(%v)", medal.Status, err)
			continue
		}
		if status == 2 { //勋章审核状态 -1已拒绝 0未申请 1已申请 2已开通
			_, err := s.newc.UpUserTask(c, mid, v.ID)
			if err != nil {
				log.Error("medalCheckMedal s.newc.UpUserTask mid(%d)|task id(%d)|err(%v)", mid, v.ID, err)
				continue
			}
			log.Info("medalCheckMedal finish task mid(%d)|state(%d)|", mid, status)
		}
	}
}

func getTaskSortMap(tasks []*newcomer.Task) map[int8][]*newcomer.Task {
	taskSortMap := make(map[int8][]*newcomer.Task)
	for _, v := range tasks {
		switch v.TargetType {
		case newcomer.TargetType001, newcomer.TargetType014:
			taskSortMap[newcomer.ArcUpCount] = append(taskSortMap[newcomer.ArcUpCount], v)
		case newcomer.TargetType002: //分享
		case newcomer.TargetType003:
			taskSortMap[newcomer.AcaPlayCount] = append(taskSortMap[newcomer.AcaPlayCount], v)
		case newcomer.TargetType004, newcomer.TargetType005, newcomer.TargetType006, newcomer.TargetType007, newcomer.TargetType008, newcomer.TargetType009:
			taskSortMap[newcomer.DataUpStat] = append(taskSortMap[newcomer.DataUpStat], v)
		case newcomer.TargetType010, newcomer.TargetType022:
			taskSortMap[newcomer.AccProfileWithStat] = append(taskSortMap[newcomer.AccProfileWithStat], v)
		case newcomer.TargetType011:
			taskSortMap[newcomer.WmWaterMark] = append(taskSortMap[newcomer.WmWaterMark], v)
		case newcomer.TargetType012:
			taskSortMap[newcomer.AccRelation] = append(taskSortMap[newcomer.AccRelation], v)
		case newcomer.TargetType013, newcomer.TargetType015, newcomer.TargetType016, newcomer.TargetType017, newcomer.TargetType018, newcomer.TargetType019, newcomer.TargetType020, newcomer.TargetType021:
			taskSortMap[newcomer.DataUpArchiveStat] = append(taskSortMap[newcomer.DataUpArchiveStat], v)
		case newcomer.TargetType023:
			taskSortMap[newcomer.OrderGrowAccountState] = append(taskSortMap[newcomer.OrderGrowAccountState], v)
		case newcomer.TargetType024:
			taskSortMap[newcomer.MedalCheckMedal] = append(taskSortMap[newcomer.MedalCheckMedal], v)
		}
	}
	return taskSortMap
}
