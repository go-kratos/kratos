package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/log"
)

func (s *Service) loadAlgorithm() (err error) {
	ss, err := s.dao.SlotStats(context.Background())
	if err != nil {
		return
	}
	// 按name聚合slot
	ssMap := make(map[string]*model.SlotsStat)
	for _, s := range ss {
		if v, exists := ssMap[s.Name]; exists {
			v.Slots = append(v.Slots, s.Slot)
		} else {
			ssMap[s.Name] = &model.SlotsStat{
				Name:      s.Name,
				Slots:     []int{s.Slot},
				Algorithm: s.Algorithm,
				Weight:    s.Weight,
			}
		}
	}
	var algorithms []model.Algorithm
	for name, ss := range ssMap {
		if ss.Weight == "" {
			continue
		}
		var (
			algorithm model.Algorithm
			w         interface{}
		)
		if ss.Algorithm == model.WilsonLHRRAlgorithm || ss.Algorithm == model.WilsonLHRRFluidAlgorithm {
			if err = json.Unmarshal([]byte(ss.Weight), &w); err != nil {
				log.Error("json.Unmarshal() error(%v), name (%s), algorightm (%s), weight (%s)", err, name, ss.Algorithm, ss.Weight)
				return
			}
		}
		switch ss.Algorithm {
		case model.WilsonLHRRAlgorithm:
			weight := w.(map[string]interface{})
			algorithm = model.NewWilsonLHRR(name, ss.Slots, &model.WilsonLHRRWeight{
				Like:   weight["like"].(float64),
				Hate:   weight["hate"].(float64),
				Reply:  weight["reply"].(float64),
				Report: weight["report"].(float64),
			})
		case model.WilsonLHRRFluidAlgorithm:
			weight := w.(map[string]interface{})
			algorithm = model.NewWilsonLHRRFluid(name, ss.Slots, &model.WilsonLHRRFluidWeight{
				Like:   weight["like"].(float64),
				Hate:   weight["hate"].(float64),
				Reply:  weight["reply"].(float64),
				Report: weight["report"].(float64),
				Slope:  weight["slope"].(float64),
			})
		case model.OriginAlgorithm:
			algorithm = model.NewOrigin(name, ss.Slots)
		case model.LikeDescAlgorithm:
			algorithm = model.NewLikeDesc(name, ss.Slots)
		case model.DefaultAlgorithm:
			continue
		default:
			log.Warn("invalid algorithm")
			continue
		}
		if algorithm != nil {
			algorithms = append(algorithms, algorithm)
		}
	}
	s.algorithmsLock.Lock()
	s.algorithms = algorithms
	s.algorithmsLock.Unlock()
	return
}

func (s *Service) loadSlots() (err error) {
	ctx := context.Background()
	slotsMap, err := s.dao.SlotsMapping(ctx)
	if err != nil {
		return
	}
	s.statisticsLock.Lock()
	for name, mapping := range slotsMap {
		for _, slot := range mapping.Slots {
			s.statisticsStats[slot].Name = name
			s.statisticsStats[slot].Slot = slot
		}
		log.Warn("name stat (name: %s, slots: %v)", name, mapping.Slots)
	}
	log.Warn("statistics stat (%v)", s.statisticsStats)
	s.statisticsLock.Unlock()
	return
}
