package service

import (
	"context"
	"time"

	"go-common/app/interface/main/history/model"
	hisapi "go-common/app/service/main/history/api/grpc"
	history "go-common/app/service/main/history/model"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

func (s *Service) serviceRun(f func()) {
	select {
	case s.serviceChan <- f:
	default:
		log.Error("serviceChan full")
	}
}

func (s *Service) serviceproc() {
	for {
		f := <-s.serviceChan
		f()
	}
}

func (s *Service) serviceAdd(arg *model.History) {
	s.serviceRun(func() {
		h := arg.ConvertServiceType()
		arg := &hisapi.AddHistoryReq{
			Mid:      h.Mid,
			Business: h.Business,
			Kid:      h.Kid,
			Aid:      h.Aid,
			Sid:      h.Sid,
			Epid:     h.Epid,
			Cid:      h.Cid,
			SubType:  h.SubType,
			Device:   h.Device,
			Progress: h.Progress,
			ViewAt:   h.ViewAt,
		}
		s.hisRPC.AddHistory(context.Background(), arg)
	})
}

func (s *Service) serviceAdds(mid int64, hs []*model.History) {
	s.serviceRun(func() {
		arg := &hisapi.AddHistoriesReq{}
		for _, a := range hs {
			h := a.ConvertServiceType()
			arg.Histories = append(arg.Histories, &hisapi.AddHistoryReq{
				Mid:      mid,
				Business: h.Business,
				Kid:      h.Kid,
				Aid:      h.Aid,
				Sid:      h.Sid,
				Epid:     h.Epid,
				Cid:      h.Cid,
				SubType:  h.SubType,
				Device:   h.Device,
				Progress: h.Progress,
				ViewAt:   h.ViewAt,
			})
		}
		s.hisRPC.AddHistories(context.Background(), arg)
	})
}

func (s *Service) serviceDel(ctx context.Context, mid int64, his []*model.History) error {
	arg := &hisapi.DelHistoriesReq{
		Mid:     mid,
		Records: []*hisapi.DelHistoriesReq_Record{},
	}
	var aids []int64
	for _, v := range his {
		if v.TP == model.TypePGC || v.TP == model.TypeUGC {
			aids = append(aids, v.Aid)
		}
	}
	seasonMap := make(map[int64]*model.BangumiSeason)
	if len(aids) > 0 {
		seasonMap, _ = s.season(ctx, mid, aids, metadata.String(ctx, metadata.RemoteIP))
	}
	for _, h := range his {
		if value, ok := seasonMap[h.Aid]; ok && value != nil {
			arg.Records = append(arg.Records, &hisapi.DelHistoriesReq_Record{
				Business: model.BusinessByTP(model.TypePGC),
				ID:       value.ID,
			})
			log.Warn("seasonMap(%d,%v)season:%d", mid, h.Aid, value.ID)
		}
		arg.Records = append(arg.Records, &hisapi.DelHistoriesReq_Record{
			Business: model.BusinessByTP(h.TP),
			ID:       h.Aid,
		})
	}
	if _, err := s.hisRPC.DelHistories(ctx, arg); err != nil {
		log.Error("s.hisRPC.DelHistories(%+v %+v) err:%+v", his, arg, errors.WithStack(err))
		return err
	}
	return nil
}

func (s *Service) serviceClear(mid int64, tps []int8) {
	s.serviceRun(func() {
		arg := &hisapi.ClearHistoryReq{
			Mid: mid,
		}
		for _, t := range tps {
			arg.Businesses = append(arg.Businesses, model.BusinessByTP(t))
		}
		s.hisRPC.ClearHistory(context.Background(), arg)
	})
}

func (s *Service) serviceDels(ctx context.Context, mid int64, aids []int64, typ int8) error {
	arg := &hisapi.DelHistoriesReq{
		Mid:     mid,
		Records: []*hisapi.DelHistoriesReq_Record{},
	}
	seasonMap := make(map[int64]*model.BangumiSeason)
	if typ == 0 {
		seasonMap, _ = s.season(ctx, mid, aids, metadata.String(ctx, metadata.RemoteIP))
	}
	b := model.BusinessByTP(typ)
	for _, aid := range aids {
		if value, ok := seasonMap[aid]; ok && value != nil {
			arg.Records = append(arg.Records, &hisapi.DelHistoriesReq_Record{
				Business: model.BusinessByTP(model.TypePGC),
				ID:       value.ID,
			})
			log.Warn("seasonMap(%d,%v)season:%d", mid, aid, value.ID)
		}
		arg.Records = append(arg.Records, &hisapi.DelHistoriesReq_Record{
			Business: b,
			ID:       aid,
		})
	}
	if _, err := s.hisRPC.DelHistories(ctx, arg); err != nil {
		log.Error("s.hisRPC.DelHistories(%v %+v) err:%+v", mid, arg, errors.WithStack(err))
		return err
	}
	return nil
}

func (s *Service) serviceHide(mid int64, hide bool) {
	s.serviceRun(func() {
		arg := &hisapi.UpdateUserHideReq{
			Mid:  mid,
			Hide: hide,
		}
		s.hisRPC.UpdateUserHide(context.Background(), arg)
	})
}

func (s *Service) serviceHistoryCursor(c context.Context, mid int64, kid int64, businesses []string, business string, viewAt int64, ps int) ([]*model.Resource, error) {
	if viewAt == 0 {
		viewAt = time.Now().Unix()
	}
	arg := &hisapi.UserHistoriesReq{
		Mid:        mid,
		Businesses: businesses,
		Business:   business,
		Kid:        kid,
		ViewAt:     viewAt,
		Ps:         int64(ps),
	}
	reply, err := s.hisRPC.UserHistories(c, arg)
	if err != nil {
		log.Error("s.hisRPC.UserHistories(%+v) err:%+v", arg, err)
		return nil, err
	}
	if reply == nil {
		return nil, err
	}
	his := make([]*model.Resource, 0)
	for _, v := range reply.Histories {
		tp, _ := model.CheckBusiness(v.Business)
		his = append(his, &model.Resource{
			Mid:      v.Mid,
			Oid:      v.Aid,
			Sid:      v.Sid,
			Epid:     v.Epid,
			TP:       tp,
			Business: v.Business,
			STP:      int8(v.SubType),
			Cid:      v.Cid,
			DT:       int8(v.Device),
			Pro:      int64(v.Progress),
			Unix:     v.ViewAt,
		})
	}
	return his, nil
}

func (s *Service) servicePnPsCursor(c context.Context, mid int64, businesses []string, pn, ps int) ([]*model.History, []int64, error) {
	if pn*ps > 1000 {
		return nil, nil, nil
	}
	arg := &hisapi.UserHistoriesReq{
		Mid:        mid,
		Businesses: businesses,
		Ps:         int64(pn * ps),
		ViewAt:     time.Now().Unix(),
	}
	reply, err := s.hisRPC.UserHistories(c, arg)
	if err != nil {
		log.Error("s.hisRPC.UserHistories(%+v) err:%+v", arg, err)
		return nil, nil, err
	}
	if reply == nil {
		return nil, nil, err
	}
	size := len(reply.Histories)
	start := (pn - 1) * ps
	end := start + ps - 1
	switch {
	case size > start && size > end:
		reply.Histories = reply.Histories[start : end+1]
	case size > start && size <= end:
		reply.Histories = reply.Histories[start:]
	default:
		reply.Histories = make([]*history.History, 0)
	}
	var epids []int64
	his := make([]*model.History, 0)
	for _, v := range reply.Histories {
		tp, _ := model.CheckBusiness(v.Business)
		if tp == model.TypePGC {
			epids = append(epids, v.Epid)
		}
		his = append(his, &model.History{
			Mid:      v.Mid,
			Aid:      v.Aid,
			Sid:      v.Sid,
			Epid:     v.Epid,
			TP:       tp,
			Business: v.Business,
			STP:      int8(v.SubType),
			Cid:      v.Cid,
			DT:       int8(v.Device),
			Pro:      int64(v.Progress),
			Unix:     v.ViewAt,
		})
	}
	return his, epids, nil
}

func (s *Service) servicePosition(c context.Context, mid int64, business string, kids []int64) (map[int64]*model.History, error) {
	arg := &hisapi.HistoriesReq{
		Mid:      mid,
		Business: business,
		Kids:     kids,
	}
	reply, err := s.hisRPC.Histories(c, arg)
	if err != nil {
		log.Error("s.hisRPC.Histories(%+v) err:%+v", arg, err)
		return nil, err
	}
	if reply == nil {
		return nil, err
	}
	now := time.Now().Unix() - 8*60*60
	his := make(map[int64]*model.History)
	for _, v := range reply.Histories {
		if business == model.BusinessByTP(model.TypeUGC) && v.ViewAt < now {
			continue
		}
		tp, _ := model.CheckBusiness(v.Business)
		his[v.Aid] = &model.History{
			Mid:      v.Mid,
			Aid:      v.Aid,
			Sid:      v.Sid,
			Epid:     v.Epid,
			TP:       tp,
			Business: v.Business,
			STP:      int8(v.SubType),
			Cid:      v.Cid,
			DT:       int8(v.Device),
			Pro:      int64(v.Progress),
			Unix:     v.ViewAt,
		}
	}
	return his, nil
}

func (s *Service) serviceHistoryType(c context.Context, mid int64, business string, kids []int64) ([]*model.History, error) {
	arg := &hisapi.HistoriesReq{
		Mid:      mid,
		Business: business,
		Kids:     kids,
	}
	reply, err := s.hisRPC.Histories(c, arg)
	if err != nil {
		log.Error("s.hisRPC.Histories(%+v) err:%+v", arg, err)
		return nil, err
	}
	if reply == nil {
		return nil, err
	}
	his := make([]*model.History, 0)
	for _, v := range reply.Histories {
		tp, _ := model.CheckBusiness(v.Business)
		his = append(his, &model.History{
			Mid:      v.Mid,
			Aid:      v.Aid,
			Sid:      v.Sid,
			Epid:     v.Epid,
			TP:       tp,
			Business: v.Business,
			STP:      int8(v.SubType),
			Cid:      v.Cid,
			DT:       int8(v.Device),
			Pro:      int64(v.Progress),
			Unix:     v.ViewAt,
		})
	}
	return his, nil
}

func (s *Service) serviceHideState(c context.Context, mid int64) (int64, error) {
	arg := &hisapi.UserHideReq{
		Mid: mid,
	}
	reply, err := s.hisRPC.UserHide(c, arg)
	if err != nil {
		log.Error("s.hisRPC.UserHide(%d) err:%+v", mid, err)
		return 0, err
	}
	if !reply.Hide {
		return 0, nil
	}
	return 1, nil
}
