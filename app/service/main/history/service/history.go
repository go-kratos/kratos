package service

import (
	"context"
	"sort"
	"strconv"
	"sync"

	pb "go-common/app/service/main/history/api/grpc"
	"go-common/app/service/main/history/model"
	"go-common/library/sync/errgroup"
	"go-common/library/sync/pipeline"
)

const _delBatch = 20

// AddHistory add history
func (s *Service) AddHistory(c context.Context, arg *pb.AddHistoryReq) (reply *pb.AddHistoryReply, err error) {
	if err = s.checkBusiness(arg.Business); err != nil {
		return
	}
	reply = &pb.AddHistoryReply{}
	// 用户忽略播放历史
	userReply, _ := s.UserHide(c, &pb.UserHideReq{Mid: arg.Mid})
	if userReply != nil && userReply.Hide {
		return
	}
	if err = s.dao.AddHistoryCache(c, arg); err != nil {
		return
	}
	s.addMerge(c, arg.Business, arg.Mid, arg.Kid, arg.ViewAt)
	return
}

// AddHistories 增加多条播放历史记录
func (s *Service) AddHistories(c context.Context, arg *pb.AddHistoriesReq) (reply *pb.AddHistoriesReply, err error) {
	var his []*pb.AddHistoryReq
	reply = &pb.AddHistoriesReply{}
	g := &errgroup.Group{}
	merges := make([]*pb.AddHistoryReq, 0, 100)
	mutex := &sync.Mutex{}
	for _, a := range arg.Histories {
		if err = s.checkBusiness(a.Business); err != nil {
			return
		}
		a := a
		g.Go(func() error {
			// 用户忽略播放历史
			userReply, _ := s.UserHide(c, &pb.UserHideReq{Mid: a.Mid})
			if userReply != nil && userReply.Hide {
				return nil
			}
			mutex.Lock()
			his = append(his, a)
			merges = append(merges, a)
			mutex.Unlock()
			return nil
		})
	}
	g.Wait()
	if err = s.dao.AddHistoriesCache(c, his); err != nil {
		return
	}
	for _, a := range merges {
		s.addMerge(c, a.Business, a.Mid, a.Kid, a.ViewAt)
	}
	return
}

// DelHistories delete histories
func (s *Service) DelHistories(c context.Context, arg *pb.DelHistoriesReq) (reply *pb.DelHistoriesReply, err error) {
	reply = &pb.DelHistoriesReply{}
	for _, r := range arg.Records {
		if err = s.checkBusiness(r.Business); err != nil {
			return
		}
	}
	if len(arg.Records) > _delBatch {
		g := errgroup.Group{}
		for i := 0; i < len(arg.Records); i += _delBatch {
			a := &pb.DelHistoriesReq{Mid: arg.Mid}
			if i+_delBatch > len(arg.Records) {
				a.Records = arg.Records[i:len(arg.Records)]
			} else {
				a.Records = arg.Records[i : i+_delBatch]
			}
			g.Go(func() (err error) {
				err = s.dao.DeleteHistories(c, a)
				return
			})
		}
		if err = g.Wait(); err != nil {
			return
		}
	} else {
		if err = s.dao.DeleteHistories(c, arg); err != nil {
			return
		}
	}
	err = s.dao.DelHistoryCache(c, arg)
	return
}

// ClearHistory clear histories
func (s *Service) ClearHistory(c context.Context, arg *pb.ClearHistoryReq) (reply *pb.ClearHistoryReply, err error) {
	reply = &pb.ClearHistoryReply{}
	for _, business := range arg.Businesses {
		if err = s.checkBusiness(business); err != nil {
			return
		}
	}
	// mid下数据量很大 异步处理 防止超时错误
	s.asyncFunc(func() {
		if len(arg.Businesses) > 0 {
			s.dao.ClearHistory(context.Background(), arg.Mid, arg.Businesses)
			return
		}
		s.dao.ClearAllHistory(context.Background(), arg.Mid)
	})
	var businesses []string
	if len(arg.Businesses) == 0 {
		for b := range s.businessNames {
			businesses = append(businesses, b)
		}
	} else {
		businesses = arg.Businesses
	}
	err = s.dao.ClearHistoryCache(c, arg.Mid, businesses)
	return
}

// UserHistories 查询用户的播放历史列表
func (s *Service) UserHistories(c context.Context, arg *pb.UserHistoriesReq) (reply *pb.UserHistoriesReply, err error) {
	g := &errgroup.Group{}
	var hisIds map[string][]int64
	var his map[string][]*model.History
	var cacheHis map[string]map[int64]*model.History
	var err1, err2 error
	g.Go(func() error {
		var names = arg.Businesses
		if len(names) == 0 {
			names = make([]string, 0)
			for name := range s.businessNames {
				names = append(names, name)
			}
		}
		if hisIds, err1 = s.dao.ListsCacheByTime(c, names, arg.Mid, arg.ViewAt, arg.Ps); err1 != nil {
			return nil
		}
		cacheHis, err1 = s.dao.HistoriesCache(c, arg.Mid, hisIds)
		return nil
	})
	g.Go(func() error {
		his, err2 = s.dao.UserHistories(c, arg.Businesses, arg.Mid, arg.ViewAt, arg.Ps)
		return nil
	})
	g.Wait()
	if err1 != nil && err2 != nil {
		err = err2
		return
	}
	if cacheHis == nil {
		cacheHis = make(map[string]map[int64]*model.History)
	}
	// 去重 优先用缓存数据
	for business, hs := range his {
		if cacheHis[business] == nil {
			cacheHis[business] = make(map[int64]*model.History)
		}
		for _, h := range hs {
			if _, ok := cacheHis[business][h.Kid]; !ok {
				cacheHis[business][h.Kid] = h
			}
		}
	}
	histories := make([]*model.History, 0, len(cacheHis))
	for _, hs := range cacheHis {
		for _, h := range hs {
			// 过滤上一条
			if h.Kid != arg.Kid || arg.Business != h.Business {
				histories = append(histories, h)
			}
		}
	}
	sort.Slice(histories, func(i, j int) bool { return histories[i].ViewAt > histories[j].ViewAt })
	if int64(len(histories)) > arg.Ps {
		histories = histories[0:arg.Ps]
	}
	reply = &pb.UserHistoriesReply{Histories: histories}
	return
}

// Histories 根据id查询播放历史
func (s *Service) Histories(c context.Context, arg *pb.HistoriesReq) (reply *pb.HistoriesReply, err error) {
	var (
		cacheHis map[string]map[int64]*model.History
		his      map[int64]*model.History
	)
	cacheHis, _ = s.dao.HistoriesCache(c, arg.Mid, map[string][]int64{arg.Business: arg.Kids})
	var miss []int64
	for _, id := range arg.Kids {
		if (cacheHis[arg.Business] == nil) || (cacheHis[arg.Business][id] == nil) {
			miss = append(miss, id)
		}
	}
	if cacheHis[arg.Business] != nil {
		his = cacheHis[arg.Business]
	}
	if his == nil {
		his = make(map[int64]*model.History)
	}
	if len(miss) > 0 {
		hs, _ := s.dao.Histories(c, arg.Business, arg.Mid, miss)
		for k, v := range hs {
			his[k] = v
		}
	}
	reply = &pb.HistoriesReply{Histories: his}
	return
}

func (s *Service) addMerge(c context.Context, business string, mid, kid, time int64) {
	s.merge.Add(c, strconv.FormatInt(mid, 10), &model.Merge{
		Mid:  mid,
		Kid:  kid,
		Bid:  s.businessNames[business].ID,
		Time: time,
	})
}

func (s *Service) initMerge() {
	s.merge = pipeline.NewPipeline(s.c.Merge)
	s.merge.Split = func(a string) int {
		mid, _ := strconv.ParseInt(a, 10, 64)
		return int(mid) % s.c.Merge.Worker
	}
	s.merge.Do = func(c context.Context, ch int, values map[string][]interface{}) {
		var merges []*model.Merge
		for _, vs := range values {
			for _, v := range vs {
				merges = append(merges, v.(*model.Merge))
			}
		}
		s.dao.AddHistoryMessage(c, ch, merges)
	}
	s.merge.Start()
}
