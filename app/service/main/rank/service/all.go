package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/rank/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) all(c context.Context, minAid, maxAid int64) error {
	aid := minAid
	limit := s.c.Rank.RowsLimit
	for {
		log.Info("do all aid:%d", aid)
		if maxAid != 0 && aid > maxAid {
			break
		}
		arcs, err := s.dao.ArchiveMetas(c, aid, limit)
		if err != nil {
			log.Error("s.dao.ArchiveMetas(%d,%d) error(%v)", aid, limit, err)
			continue
		}
		if len(arcs) == 0 {
			break
		}
		aid = arcs[len(arcs)-1].ID
		var (
			typeids []int64
			aids    []int64
		)
		typeidsMap := make(map[int64]struct{})
		for _, v := range arcs {
			if _, ok := typeidsMap[v.Typeid]; !ok {
				typeids = append(typeids, v.Typeid)
				typeidsMap[v.Typeid] = struct{}{}
			}
			aids = append(aids, v.ID)
		}
		// ptypeid
		typesMap, err := s.dao.ArchiveTypes(c, typeids)
		if err != nil {
			log.Error("s.dao.ArchiveTypes(%+v) error(%v)", typeids, err)
			continue
		}
		// view
		statsMap, err := s.dao.ArchiveStats(c, aids)
		if err != nil {
			log.Error("s.dao.ArchiveStats(%+v) error(%v)", aids, err)
			continue
		}
		// tv
		tvsMap, err := s.dao.ArchiveTVs(c, aids)
		if err != nil {
			log.Error("s.dao.ArchiveTVs(%+v) error(%v)", aids, err)
			continue
		}
		// data append
		for _, a := range arcs {
			f := new(model.Field)
			f.Flag = model.FlagExist
			f.Oid = a.ID
			f.Pubtime = a.SetPubtime()
			if v, ok := typesMap[a.Typeid]; ok {
				f.Pid = v.SetPid()
			}
			if v, ok := statsMap[a.ID]; ok {
				f.Click = v.SetClick()
			}
			if v, ok := tvsMap[a.ID]; ok {
				f.Result = v.Result
				f.Deleted = v.Deleted
				f.Valid = v.Valid
			}
			s.setField(a.ID, f) // write map
		}
		time.Sleep(time.Duration(s.c.Rank.BatchSleep))
	}
	fmt.Println("all map len:", minAid, maxAid, len(s.rmap))
	log.Info("do all(%d,%d) successful,map len(%d)", minAid, maxAid, len(s.rmap))
	return nil
}

func (s *Service) patch(c context.Context, begin, end time.Time) error {
	step := int64(time.Duration(xtime.Time(s.c.Rank.BatchStep)) / time.Second)
	limit := s.c.Rank.RowsLimit
	for i := begin.Unix(); i <= end.Unix(); i += step {
		// archive meta and type
		var aid int64
		for {
			arcs, err := s.dao.ArchiveMetasIncrs(c, aid, xtime.Time(i), xtime.Time(i+step), limit)
			if err != nil {
				log.Error("s.dao.ArchiveMetas(%d,%d) error(%v)", aid, limit, err)
				continue
			}
			if len(arcs) == 0 {
				break
			}
			aid = arcs[len(arcs)-1].ID
			var typeids []int64
			typeidsMap := make(map[int64]struct{})
			for _, v := range arcs {
				if _, ok := typeidsMap[v.Typeid]; !ok {
					typeids = append(typeids, v.Typeid)
					typeidsMap[v.Typeid] = struct{}{}
				}
			}
			// ptypeid
			typesMap, err := s.dao.ArchiveTypes(c, typeids)
			if err != nil {
				log.Error("s.dao.ArchiveTypes(%+v) error(%v)", typeids, err)
				continue
			}
			// data append
			for _, a := range arcs {
				f := new(model.Field)
				f.Flag = model.FlagExist
				f.Oid = a.ID
				f.Pubtime = a.SetPubtime()
				if v, ok := typesMap[a.Typeid]; ok {
					f.Pid = v.SetPid()
				}
				s.setField(a.ID, f) // write map
			}
			time.Sleep(time.Duration(s.c.Rank.BatchSleep))
		}
		// archive tv
		var id int64
		for {
			tvs, err := s.dao.ArchiveTVsIncrs(c, id, xtime.Time(i), xtime.Time(i+step), limit)
			if err != nil {
				log.Error("s.dao.ArchiveTVsIncrs(%d,%s,%s,%d) error(%v)", id, xtime.Time(i), xtime.Time(i+step), limit, err)
				continue
			}
			if len(tvs) == 0 {
				break
			}
			id = tvs[len(tvs)-1].ID
			// data append
			for _, a := range tvs {
				s.field(a.Aid).Result = a.Result
				s.field(a.Aid).Deleted = a.Deleted
				s.field(a.Aid).Valid = a.Valid
			}
			time.Sleep(time.Duration(s.c.Rank.BatchSleep))
		}
		// archive stats
		for tbl := 0; tbl < 100; tbl++ {
			var id int64
			for {
				stats, err := s.dao.ArchiveStatsIncrs(c, tbl, id, xtime.Time(i), xtime.Time(i+step), limit)
				if err != nil {
					log.Error("s.dao.ArchiveTVsIncrs(%d,%s,%s,%d) error(%v)", id, xtime.Time(i), xtime.Time(i+step), limit, err)
					continue
				}
				if len(stats) == 0 {
					break
				}
				id = stats[len(stats)-1].ID
				// data append
				for _, a := range stats {
					s.field(a.Aid).Click = a.SetClick()
				}
				time.Sleep(time.Duration(s.c.Rank.BatchSleep))
			}
		}
	}
	log.Info("patch(%s,%s) successful,len(%d)", begin, end, len(s.rmap))
	return nil
}
