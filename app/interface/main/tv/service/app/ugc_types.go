package service

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
)

// prepareUGCData loads ugc data
func (s *Service) prepareUGCData(c context.Context) (err error) {
	var (
		tids       []int32
		originData = make(map[int][]*model.Card)
	)
	if tids, err = s.arcDao.TargetTypes(); err != nil {
		log.Error("[PrepareUGCData] TargetTypes Err %v", err)
		return
	}
	for _, v := range tids {
		if tidData, errUGC := s.ugcData(c, v); errUGC != nil {
			log.Error("[PrepareUGCData] Tid %d, Err %v", v, errUGC)
			continue
		} else {
			originData[int(v)] = tidData
		}
	}
	if len(originData) > 0 {
		s.UGCOrigins = originData
	}
	return
}

// ugcData gets the origin ugc data and intervene with TV CMS data
func (s *Service) ugcData(c context.Context, tid int32) (data []*model.Card, err error) {
	var (
		arcMetas []*model.ArcCMS
		origin   []*model.AIData
		aids     []int64
	)
	if origin, err = s.dao.UgcAIData(c, int16(tid)); err != nil {
		log.Error("[ugcData] Can't Pick AI Data, Tid: %d, Err: %v", tid, err)
		return
	}
	for _, v := range origin {
		aids = append(aids, int64(v.AID))
	}
	if arcMetas, err = s.cmsDao.LoadArcsMedia(c, aids); err != nil {
		log.Error("[ugcData] Can't Pick MediaCache Data, Tid: %d, Err: %v", tid, err)
		return
	}
	for _, v := range arcMetas {
		data = append(data, v.ToCard())
	}
	return
}

// SearchTypes return the ugc types
func (s *Service) SearchTypes() (res []*model.ArcType, err error) {
	var typeMap map[int32]*model.ArcType
	if typeMap, err = s.arcDao.FirstTypes(); err != nil {
		log.Error("SearchTypes ArcDao FirstTypes Err %v", err)
		return
	}
	for _, v := range s.conf.Cfg.ZonesInfo.UgcTypes {
		if tp, ok := typeMap[v]; ok {
			res = append(res, tp)
		}
	}
	return
}
