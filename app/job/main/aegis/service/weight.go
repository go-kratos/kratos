package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/aegis/model"
	"go-common/library/log"
)

// WeightManager weight manager
type WeightManager struct {
	s *Service

	businessID, flowID int64
	toplen, batchlen   int64
	minute             int64

	// cache
	topweightList []*model.WeightItem

	// channel
	redisWeightList chan *model.WeightItem
	dbWeightList    chan *model.WeightItem
	asignList       chan *model.Task
	//dbstartSig, dbstopSig chan struct{}
	redisFinish chan struct{}

	//closeChan chan struct{}
	close bool
}

var _defaultopt = &model.WeightOPT{
	TopListLen:   1000,
	BatchListLen: 1000,
	RedisListLen: 10000,
	DbListLen:    2000,
	AssignLen:    100,
	Minute:       3,
}

// NewWeightManager new
func NewWeightManager(s *Service, opt *model.WeightOPT, key string) (wm *WeightManager) {
	if opt == nil {
		opt = _defaultopt
	} else {
		if opt.TopListLen <= 0 {
			opt.TopListLen = _defaultopt.TopListLen
		}
		if opt.BatchListLen <= 0 {
			opt.BatchListLen = _defaultopt.BatchListLen
		}
		if opt.RedisListLen <= 0 {
			opt.RedisListLen = _defaultopt.RedisListLen
		}
		if opt.DbListLen <= 0 {
			opt.DbListLen = _defaultopt.DbListLen
		}
		if opt.AssignLen <= 0 {
			opt.AssignLen = _defaultopt.AssignLen
		}
		if opt.Minute <= 0 {
			opt.Minute = _defaultopt.Minute
		}
	}

	if len(key) > 0 {
		bizid, flowid := parseKey(key)
		opt.BusinessID = int64(bizid)
		opt.FlowID = int64(flowid)
	}

	wm = &WeightManager{
		s:               s,
		businessID:      opt.BusinessID,
		flowID:          opt.FlowID,
		toplen:          opt.TopListLen,
		batchlen:        opt.BatchListLen,
		minute:          opt.Minute,
		redisWeightList: make(chan *model.WeightItem, opt.RedisListLen),
		dbWeightList:    make(chan *model.WeightItem, opt.DbListLen),
		asignList:       make(chan *model.Task, opt.AssignLen),
		redisFinish:     make(chan struct{}),
	}

	go wm.weightProc()
	go wm.weightWatcher()
	log.Info("启动权重计算器 bizid(%d) flowid(%d) opt(%+v)", wm.businessID, wm.flowID, opt)
	return
}

func parseKey(key string) (bizid, flowid int) {
	pos := strings.Index(key, "-")
	bizids := key[:pos]
	flowids := key[pos+1:]
	bizid, _ = strconv.Atoi(bizids)
	flowid, _ = strconv.Atoi(flowids)
	return
}

func (s *Service) startWeightManager() {
	// 1.当前的所有业务线，需要计算权重的先枚举出来
	s.wmHash = make(map[string]*WeightManager)
	for key := range s.newactiveBizFlow {
		bizid, _ := parseKey(key)
		s.wmHash[key] = NewWeightManager(s, s.getWeightOpt(bizid), key)
	}
}

func (w *WeightManager) weightProc() {
	for !w.close {
		if err := w.weightRedisProcess(); err != nil {
			w.weightDBProcess()
		}
		time.Sleep(time.Duration(w.minute) * time.Minute)
	}
}

func (w *WeightManager) weightWatcher() {
	for !w.close {
		select {
		case <-w.redisFinish: //取出权重最大的一批，更新到数据库
			log.Info("redisFinish(%d-%d:%d)", w.businessID, w.flowID, w.toplen)
			w.handleRedisFinish(context.Background())
		case wi := <-w.redisWeightList:
			w.handleRedisWeightList(context.Background(), wi)
		case wi := <-w.dbWeightList:
			w.handleDBWeightList(context.Background(), wi)
		case task := <-w.asignList:
			w.handleAssign(context.Background(), task)
		}
	}
}

func (w *WeightManager) weightRedisProcess() (err error) {
	var c = context.Background()
	if err = w.s.dao.CreateUnionSet(c, w.businessID, w.flowID); err != nil {
		return
	}

	var (
		start = int64(0)
		stop  = w.batchlen
	)
	for {
		wis, err := w.s.dao.RangeUinonSet(c, w.businessID, w.flowID, start, stop)
		if err != nil {
			return err
		}
		log.Info("weightRedisProcess length(%d) start(%d) stop(%d)", len(wis), start, stop)
		start += w.batchlen
		stop += w.batchlen
		if len(wis) == 0 {
			break
		}
		for _, wi := range wis {
			if w.caculateWeight(c, wi) {
				log.Warn("weightRedisProcess 任务未找到 wi(%+v)", wi)
				continue
			}
			w.s.dao.SetWeight(c, w.businessID, w.flowID, wi.ID, wi.Weight)
		}
		time.Sleep(time.Second)
	}
	w.redisFinish <- struct{}{}
	w.s.dao.DeleteUinonSet(c, w.businessID, w.flowID)
	return nil
}

func (w *WeightManager) caculateWeight(c context.Context, wi *model.WeightItem) (skip bool) {
	task, err := w.s.dao.GetTask(c, wi.ID)
	if err != nil {
		return true
	}
	w.reAssign(c, task)

	wm := int64(time.Since(task.Ctime.Time()).Minutes())
	wl := &model.WeightLog{
		UPtime:   time.Now().Format("2006-01-02 15:04:05"),
		Mid:      task.MID,
		Fans:     task.Fans,
		Group:    task.Group,
		WaitTime: model.WaitTime(task.Ctime.Time()),
	}

	var wtRange, wtEqual int64

	wci, ewc := w.s.getWeightCache(c, task.BusinessID, task.FlowID)
	if wci != nil {
		wtRange = w.rangeCaculate(c, wci, task, wm, wl)
	}
	if ewc != nil {
		wtEqual = w.equalCaculate(c, ewc, task, wm, wl)
	}
	wi.Weight = wtRange + wtEqual
	wl.Weight = wi.Weight

	w.s.sendWeightLog(c, task, wl)
	return
}

func (w *WeightManager) rangeCaculate(c context.Context, wci map[string]*model.RangeWeightConfig, task *model.Task, wt int64, wl *model.WeightLog) (weight int64) {
	var wtWeight, fanWeight, groupWeight int64

	if cfg, ok := wci["waittime"]; ok {
		if wtlen := len(cfg.Range); wtlen > 0 { // 等待时长，要把之前等级的权重加上去
			for i := wtlen - 1; i >= 0; i-- {
				if wt >= cfg.Range[i].Threshold { // 命中配置
					wtWeight += cfg.Range[i].Weight * ((wt - cfg.Range[i].Threshold) / w.minute)

					// 计算0 到 (i-1) 累计权重
					for j := 0; j <= i-1; j++ {
						wtWeight += cfg.Range[j].Weight * ((cfg.Range[j+1].Threshold - cfg.Range[j].Threshold) / w.minute)
					}
					break
				}
			}
		}
	}

	if cfg, ok := wci["fans"]; ok {
		if fanLen := len(cfg.Range); fanLen > 0 {
			for i := fanLen - 1; i >= 0; i-- {
				if task.Fans >= cfg.Range[i].Threshold {
					fanWeight = cfg.Range[i].Weight * (wt / w.minute)
					break
				}
			}
		}
	}

	if cfg, ok := wci["group"]; ok {
		if len(cfg.Range) > 0 {
			for _, item := range cfg.Range {
				if strings.Contains(","+task.Group+",", fmt.Sprintf(",%d,", item.Threshold)) {
					groupWeight = item.Weight * (wt / w.minute)
				}
			}
		}
	}

	weight = wtWeight + fanWeight + groupWeight
	wl.WaitWeight = wtWeight
	wl.FansWeight = fanWeight
	wl.GroupWeight = groupWeight
	return
}

func (w *WeightManager) equalCaculate(c context.Context, ewc []*model.EqualWeightConfig, task *model.Task, wt int64, wl *model.WeightLog) (weight int64) {
	var midweight, taskweight int64
	for _, item := range ewc {
		if item.Name == "mid" {
			if strings.Contains(","+item.IDs+",", fmt.Sprintf(",%d,", task.MID)) {
				if item.Type == model.WeightTypeCycle {
					midweight += item.Weight * (wt / w.minute)
				} else {
					midweight += item.Weight
				}
				log.Info("equalCaculate task(%+v) hit (%+v)", task, item)
				wl.ConfigItems = append(wl.ConfigItems, &model.ConfigItem{
					Name:  item.Name,
					Desc:  item.Description,
					Uname: item.Uname,
				})
			}
		}
		if item.Name == "taskid" || item.Name == "task_id" {
			if strings.Contains(","+item.IDs+",", fmt.Sprintf(",%d,", task.ID)) {
				if item.Type == model.WeightTypeCycle {
					taskweight += item.Weight * (wt / w.minute)
				} else {
					taskweight += item.Weight
				}
				log.Info("equalCaculate task(%+v) hit (%+v)", task, item)
				wl.ConfigItems = append(wl.ConfigItems, &model.ConfigItem{
					Name:  item.Name,
					Desc:  item.Description,
					Uname: item.Uname,
				})
			}
		}
	}
	weight = midweight + taskweight

	wl.EqualWeight = weight
	return
}

func (w *WeightManager) reAssign(c context.Context, task *model.Task) {
	if task.UID == 0 {
		select {
		case w.asignList <- task:
			log.Info("指派判断 reAssign（%+v）", task)
		case <-time.NewTimer(10 * time.Millisecond).C:
			log.Warn("chan asignList full,len:%d", len(w.dbWeightList))
		}
	}
}

func (w *WeightManager) weightDBProcess() (err error) {
	// TODO 只用db更新权重的策略
	return nil
}

func (w *WeightManager) handleAssign(c context.Context, task *model.Task) (err error) {
	if w.s.setAssign(c, task) {
		if rows, err := w.s.dao.AssignTask(c, task); err == nil && rows == 1 {
			w.s.dao.SetTask(c, task)
		}
	}
	return
}

func (w *WeightManager) handleRedisWeightList(c context.Context, wi *model.WeightItem) (err error) {
	return w.s.dao.SetWeight(c, w.businessID, w.flowID, wi.ID, wi.Weight)
}

func (w *WeightManager) handleDBWeightList(c context.Context, wi *model.WeightItem) (rows int64, err error) {
	return w.s.dao.SetWeightDB(c, wi.ID, wi.Weight)
}

func (w *WeightManager) handleRedisFinish(c context.Context) (err error) {
	log.Info("handleRedisFinish")
	wis, err := w.s.dao.TopWeights(c, w.businessID, w.flowID, w.toplen)
	if err != nil {
		return
	}

	tempMap := make(map[int64]struct{})
	for _, wi := range wis {
		log.Info("handleRedisFinish:(%+v)", wi)
		w.addToDBList(wi)
		tempMap[wi.ID] = struct{}{}
	}

	for _, wi := range w.topweightList {
		if _, ok := tempMap[wi.ID]; !ok {
			weight, err := w.s.dao.GetWeight(c, w.businessID, w.flowID, wi.ID)
			if err != nil {
				continue
			}
			wi.Weight = weight
			w.addToDBList(wi)
		}
	}
	w.topweightList = wis
	log.Info("handleRedisFinish:topweightList(%d)", len(wis))

	return
}

func (w *WeightManager) addToDBList(wi *model.WeightItem) {
	select {
	case w.dbWeightList <- wi:
		log.Info("addToDBList (%+v)", wi)
	case <-time.NewTimer(10 * time.Millisecond).C:
		log.Warn("chan dbWeightList full,len:%d", len(w.dbWeightList))
	}
}
