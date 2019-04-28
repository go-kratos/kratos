package service

import (
	"context"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/task"
	"go-common/app/job/main/videoup-report/model/utils"
	"go-common/library/log"
)

const (
	_pm = int64(3) //period minute
	//普通任务参数
	_nc1 = int64(3)  // 等待时长9分钟
	_nc2 = int64(5)  // 等待时长15分钟
	_nc3 = int64(9)  // 等待时长27分钟
	_nc4 = int64(15) // 等待时长45分钟

	//定时任务参数
	_tc1 = int64(80) // 距离发布4小时
	_tc2 = int64(40) // 距离发布2小时
	_tc3 = int64(20) // 距离发布1小时
)

func (s *Service) upTaskWeightCache() (err error) {
	var (
		ids      []int64
		firstid  int64
		lastid   int64
		tasksMap map[int64]*task.WeightParams
	)
	c := context.TODO()
	s.jumplist.Reset()

	defer func() {
		s.lastjumpMap = make(map[int64]struct{})
		for item := s.jumplist.POP(); item != nil; item = s.jumplist.POP() {
			select {
			case s.jumpchan <- item:
				s.lastjumpMap[item.TaskID] = struct{}{}
			default:
				log.Warn("jumpchan full")
			}
		}
	}()

	for {
		cacheMap := make(map[int64]*task.WeightParams)
		// 先从数据库批量读取出id来
		if ids, lastid, err = s.arc.TaskIDforWeight(c, firstid); len(ids) == 0 {
			if err != nil {
				log.Error("s.task.TaskIDforWeight(%d) error(%v)", firstid, err)
			}
			return
		}
		// 先从redis取权重配置,redis挂了从数据库读取权重配置
		if tasksMap, err = s.redis.GetWeight(c, ids); len(tasksMap) != len(ids) {
			log.Warn("GetTaskWeight from redis need len(%d) while len(%d)", len(ids), len(tasksMap))
			if tasksMap, err = s.arc.GetTaskWeight(c, firstid); err != nil {
				log.Error("s.arc.GetTaskWeight(%d) error(%v)", firstid, err)
				return
			}
		}
		for id, tp := range tasksMap {
			if tp.State == 0 {
				s.setTask(c, id, tp, cacheMap)
			}
		}
		if len(cacheMap) > 0 {
			s.redis.SetWeight(c, cacheMap)
		}
		firstid = lastid
		time.Sleep(time.Second)
	}
}

func (s *Service) upTaskWeightDB(c context.Context, item *task.WeightLog) (err error) {
	// 超时还没入库的直接丢掉，免得阻塞
	if time.Since(item.Uptime.TimeValue()).Minutes() > 3.0 {
		log.Warn("task weight item(%v) expired! ", item)
		return
	}
	_, err = s.arc.UpTaskWeight(c, item.TaskID, item.Weight)
	if err != nil {
		log.Error("UpTaskWeight(%d,%d) error(%v)", item.Weight, item.TaskID, err)
		return
	}
	return
}

func (s *Service) setTask(c context.Context, taskid int64, tp *task.WeightParams, cm map[int64]*task.WeightParams) {
	tpc := s.newPriority(c, tp)
	// 上轮次权重前2000的这轮也必须更新
	if _, ok := s.lastjumpMap[tpc.TaskID]; ok {
		select {
		case s.jumpchan <- tpc:
		default:
			log.Warn("jumpchan full")
		}
	} else {
		s.jumplist.PUSH(tpc) //插队序列
	}

	tp.Weight = tpc.Weight
	if cm != nil {
		cm[taskid] = tp
	}
	select {
	case s.tasklogchan <- tpc:
	default:
		log.Info("s.tasklogchan full(%d)", len(s.tasklogchan))
	}
}

// 任务权重计算
func (s *Service) newPriority(c context.Context, tp *task.WeightParams) (tpc *task.WeightLog) {
	var round int64
	var wcf = task.WLVConf
	tpc = &task.WeightLog{
		TaskID: tp.TaskID,
		Mid:    tp.Mid,
	}

	if tp.AccFailed {
		if tp.Fans, tp.AccFailed = s.getUpperFans(c, tp.Mid); !tp.AccFailed {
			tp.Special = s.getUpSpecial(tp.Mid, tp.Fans)
			s.arc.SetUpSpecial(c, tp.TaskID, tp.Special)
		}
	}

	round = int64(time.Since(tp.Ctime.TimeValue()).Minutes()) / _pm
	//定时任务
	if !tp.Ptime.TimeValue().IsZero() {
		pround := int64((time.Until(tp.Ptime.TimeValue()).Minutes())) / 3
		sumtesmp := wcf.Tlv4 * (int64(tp.Ptime.TimeValue().Sub(tp.Ctime.TimeValue()).Minutes()) - 60*4) / 3
		switch {
		case pround >= _tc1:
			tpc.TWeight = round * wcf.Tlv4
		case _tc2 <= pround && pround < _tc1:
			tpc.TWeight = (_tc1-pround)*wcf.Tlv1 + sumtesmp
		case _tc3 <= pround && pround < _tc2:
			tpc.TWeight = (_tc2-pround)*wcf.Tlv2 + wcf.Tsum2h + sumtesmp
		case pround < _tc3:
			tpc.TWeight = (_tc3-pround)*wcf.Tlv3 + wcf.Tsum1h + sumtesmp
		}
		tpc.NWeight = 0
	} else { // 普通任务加权
		switch {
		case round < _nc1:
			tpc.NWeight = wcf.Nlv5 * round
		case _nc1 <= round && round < _nc2:
			tpc.NWeight = wcf.Nlv1*(round-_nc1) + wcf.Nsum9
		case _nc2 <= round && round < _nc3:
			tpc.NWeight = wcf.Nlv2*(round-_nc2) + wcf.Nsum15
		case _nc3 <= round && round < _nc4:
			tpc.NWeight = wcf.Nlv3*(round-_nc3) + wcf.Nsum27
		default:
			tpc.NWeight = wcf.Nlv4*(round-_nc4) + wcf.Nsum45
		}
		tpc.TWeight = 0
	}

	// 特殊任务加权
	switch tp.Special {
	case task.UpperBigNormal:
		tpc.SWeight = wcf.Slv1 * round
	case task.UpperSuperNormal:
		tpc.SWeight = wcf.Slv2 * round
	case task.UpperWhite:
		tpc.SWeight = wcf.Slv3 * round
	case task.UpperBigWhite:
		tpc.SWeight = wcf.Slv4 * round
	case task.UpperSuperWhite:
		tpc.SWeight = wcf.Slv5 * round
	case task.UpperSuperBlack:
		tpc.SWeight = wcf.Slv6 * round
	case task.UpperBlack:
		tpc.SWeight = wcf.Slv7 * round
	default:
		tpc.SWeight = 0
	}

	//配置任务加权
	tpc.CWeight = 0
	if len(tp.CfItems) > 0 {
		tpc.CfItems = tp.CfItems
		for _, item := range tp.CfItems {
			if item.Rule == 0 {
				round2 := int64(time.Since(item.Mtime.TimeValue()).Minutes()) / _pm
				if round2 < round {
					round = round2
				}
				tpc.CWeight += round * item.Weight
			} else {
				tpc.CWeight += item.Weight
			}
		}
		if tpc.CWeight <= wcf.MinWeight {
			tpc.CWeight = wcf.MinWeight
		}
	}

	tpc.Weight = tpc.NWeight + tpc.SWeight + tpc.TWeight + tpc.CWeight
	if tpc.Weight >= wcf.MaxWeight {
		tpc.Weight = wcf.MaxWeight
	}
	tpc.Uptime = utils.NewFormatTime(time.Now())
	return
}

// 设置特殊用户的审核任务
func (s *Service) setTaskUPSpecial(c context.Context, t *task.Task, mid int64) (fans int64, failed bool) {
	fans, failed = s.getUpperFans(c, mid)
	t.UPSpecial = s.getUpSpecial(mid, fans)
	return
}

func (s *Service) getUpSpecial(mid int64, fans int64) (upspecial int8) {
	switch {
	case s.isWhite(mid): //优质
		if fans >= task.SuperUpperTH {
			upspecial = task.UpperSuperWhite
		} else if fans >= task.BigUpperTH {
			upspecial = task.UpperSuperWhite
		} else {
			upspecial = task.UpperWhite
		}
	case s.isBlack(mid):
		if fans >= task.SuperUpperTH {
			upspecial = task.UpperSuperBlack
		} else {
			upspecial = task.UpperBlack
		}

	default:
		if fans >= task.SuperUpperTH {
			upspecial = task.UpperSuperNormal
		} else if fans >= task.BigUpperTH {
			upspecial = task.UpperBigNormal
		}
	}
	return
}

// 设置定时发布的审核任务
func (s *Service) setTaskTimed(c context.Context, t *task.Task) {
	adelay, err := s.arc.Delay(c, t.Aid)
	if err != nil {
		log.Error("s.arc.Delay(%d) error(%v)", t.Aid, err)
		return
	}
	if adelay != nil && adelay.State == 0 && !adelay.DTime.IsZero() {
		t.Ptime = utils.NewFormatTime(adelay.DTime)
	}
}

func (s *Service) setTaskUpFrom(c context.Context, aid int64, tp *task.WeightParams) {
	if addit, _ := s.arc.Addit(c, aid); addit != nil {
		tp.UpFrom = addit.UpFrom
	}
}

func (s *Service) setTaskUpGroup(c context.Context, mid int64, tp *task.WeightParams) {
	ugs := []int8{}
	for gid, cache := range s.upperCache {
		if _, ok := cache[mid]; ok {
			ugs = append(ugs, gid)
		}
	}
	tp.UpGroups = ugs
}

// 设置配置用户权重的审核任务
func (s *Service) getConfWeight(c context.Context, t *task.Task, a *archive.Archive) (cfitems []*task.ConfigItem) {
	// 1. 按照用户
	if confs, ok := s.weightCache[task.WConfMid]; ok {
		if conf, ok := confs[a.Mid]; ok {
			cfitems = append(cfitems, conf)
		}
	}
	// 2. 按照分区
	if confs, ok := s.weightCache[task.WConfType]; ok {
		if conf, ok := confs[int64(a.TypeID)]; ok {
			cfitems = append(cfitems, conf)
		}
	}
	// 3. 按照投稿来源
	if confs, ok := s.weightCache[task.WConfUpFrom]; ok {
		addit, err := s.arc.Addit(c, a.ID)
		if addit == nil {
			if err != nil {
				log.Error(" s.arc.Addit(%d) error(%v)", t.Aid, err)
			}
			return
		}

		if conf, ok := confs[int64(addit.UpFrom)]; ok {
			cfitems = append(cfitems, conf)
		}
	}
	return
}

// 权重配置
func (s *Service) weightConf(c context.Context) (cfc map[int8]map[int64]*task.ConfigItem, err error) {
	var (
		ids    []int64
		arrcfs []*task.ConfigItem
	)
	if arrcfs, err = s.arc.WeightConf(context.TODO()); err != nil {
		log.Error("s.arc.WeightConf error(%v)", err)
		return
	}
	cfc = map[int8]map[int64]*task.ConfigItem{}
	for _, item := range arrcfs {
		if !item.Et.TimeValue().IsZero() && item.Et.TimeValue().Before(time.Now()) {
			ids = append(ids, item.ID)
			continue
		}
		if _, ok := cfc[item.Radio]; !ok {
			cfc[item.Radio] = make(map[int64]*task.ConfigItem)
		}
		cfc[item.Radio][item.CID] = item
	}
	if len(ids) > 0 {
		log.Info("task config(%v) 权重配置已过期,自动失效", ids)
		s.arc.DelWeightConfs(c, ids)
	}
	return
}

func (s *Service) taskWeightConsumer() {
	var c = context.TODO()
	defer s.waiter.Done()

	for {
		if s.closed {
			return
		}
		select {
		case item, ok := <-s.jumpchan: //插队序列
			if !ok {
				log.Error("s.jumpchan closed")
				return
			}
			s.upTaskWeightDB(c, item)
		case item, ok := <-s.tasklogchan:
			if !ok {
				log.Error("s.tasklogchan closed")
				return
			}
			s.hbase.AddLog(c, item)
		default:
			time.Sleep(time.Second)
		}
	}
}

func (s *Service) taskweightproc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}
		s.upTaskWeightCache()
		log.Info("taskweightproc 插队序列(%d) 日志序列(%d)", len(s.jumpchan), len(s.tasklogchan))
		time.Sleep(3 * time.Minute)
	}
}
