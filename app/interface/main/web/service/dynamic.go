package service

import (
	"context"

	"go-common/app/interface/main/web/conf"
	"go-common/app/service/main/archive/api"
	dymdl "go-common/app/service/main/dynamic/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// DynamicRegion get dynamic region.
func (s *Service) DynamicRegion(c context.Context, rid int32, pn, ps int) (rs *dymdl.DynamicArcs3, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if rs, err = s.dy.RegionArcs3(c, &dymdl.ArgRegion3{RegionID: rid, Pn: pn, Ps: ps, RealIP: ip}); err != nil {
		log.Error("s.dy.RegionArcs3(%d,%d,%d) error(%v)", rid, pn, ps, err)
		err = nil
	} else if rs != nil && len(rs.Archives) > 0 {
		fmtArcs3(rs.Archives)
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetRegionBakCache(c, rid, pn, ps, rs)
		})
		return
	}
	if rs, err = s.dao.RegionBakCache(c, rid, pn, ps); err != nil {
		return
	}
	if rs == nil {
		err = ecode.NothingFound
	}
	return
}

func fmtArcs3(arcs []*api.Arc) {
	for _, v := range arcs {
		if v.Access >= 10000 {
			v.Stat.View = -1
		}
	}
}

// DynamicRegionTag get dynamic region tag.
func (s *Service) DynamicRegionTag(c context.Context, tagID int64, rid int32, pn, ps int) (rs *dymdl.DynamicArcs3, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if rs, err = s.dy.RegionTagArcs3(c, &dymdl.ArgRegionTag3{TagID: tagID, RegionID: rid, Pn: pn, Ps: ps, RealIP: ip}); err != nil {
		log.Error("s.dy.RegionTagArcs3(%d,%d,%d,%d) error(%v)", tagID, rid, pn, ps, err)
		err = nil
	} else if rs != nil && len(rs.Archives) > 0 {
		fmtArcs3(rs.Archives)
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetRegionTagBakCache(c, tagID, rid, pn, ps, rs)
		})
		return
	}
	if rs, err = s.dao.RegionTagBakCache(c, tagID, rid, pn, ps); err != nil {
		return
	}
	if rs == nil {
		err = ecode.NothingFound
	}
	return
}

// DynamicRegionTotal get dynamic region total.
func (s *Service) DynamicRegionTotal(c context.Context) (map[string]int, error) {
	rs, err := s.dy.RegionTotal(c, &dymdl.ArgRegionTotal{RealIP: metadata.String(c, metadata.RemoteIP)})
	if err != nil {
		log.Error("s.dy.RegionTotal error(%v)", err)
		return nil, err
	}
	return rs, nil
}

// DynamicRegions get dynamic regions.
func (s *Service) DynamicRegions(c context.Context) (rs map[int32][]*api.Arc, err error) {
	var (
		rids       []int32
		common, bg map[int32][]*api.Arc
		bgid       = int32(13)
		ip         = metadata.String(c, metadata.RemoteIP)
	)
	// get first type id
	for rid := range s.rids {
		if rid == bgid { //bangumi ignore.
			continue
		} else if rid == 167 { //guochuang use second rid 168.
			rid = 168
		}
		rids = append(rids, rid)
	}
	rs = make(map[int32][]*api.Arc, len(rids)+1)
	if common, err = s.dy.RegionsArcs3(c, &dymdl.ArgRegions3{RegionIDs: rids, Count: 10, RealIP: ip}); err != nil {
		log.Error("s.dy.RegionsArcs3(%v) error(%v)", rids, err)
		err = nil
	}
	for _, rid := range rids {
		rs[rid] = common[rid]
	}
	// bangumi type id 13 find 200,condition mid == 928123.
	if bg, err = s.dy.RegionsArcs3(c, &dymdl.ArgRegions3{RegionIDs: []int32{bgid}, Count: conf.Conf.Rule.BangumiCount, RealIP: ip}); err != nil {
		log.Error("s.dy.RegionsArcs3 error(%v)", err)
		err = nil
	} else {
		n := 1
		count := 1
		for _, arc := range bg[bgid] {
			count++
			if arc.Author.Mid == 928123 {
				rs[bgid] = append(rs[bgid], arc)
			} else {
				continue
			}
			n++
			if n > conf.Conf.Rule.RegionsCount {
				log.Info("s.dy.RegionsArcs bangumi count(%d)", count)
				break
			}
		}
		// not enough add other.
		if n <= conf.Conf.Rule.RegionsCount {
			for _, arc := range bg[bgid] {
				count++
				if arc.Author.Mid == 928123 {
					continue
				} else {
					rs[bgid] = append(rs[bgid], arc)
				}
				n++
				if n > conf.Conf.Rule.RegionsCount {
					log.Info("s.dy.RegionsArcs bangumi count(%d)", count)
					break
				}
			}
		}
	}
	if len(rs) > 0 {
		countCheck := true
		for rid, region := range rs {
			if len(region) < conf.Conf.Rule.MinDyCount {
				countCheck = false
				log.Info("countCheck rid(%d) len(%d) false", rid, len(region))
				break
			}
		}
		if countCheck {
			s.cache.Do(c, func(c context.Context) {
				s.dao.SetRegionsBakCache(c, rs)
			})
			return
		}
	}
	rs, err = s.dao.RegionsBakCache(c)
	return
}
