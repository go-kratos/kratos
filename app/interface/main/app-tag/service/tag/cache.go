package tag

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-tag/model"
	"go-common/app/interface/main/app-tag/model/region"
	"go-common/app/interface/main/app-tag/model/tag"
	"go-common/library/log"
)

// loadSimilar
func (s *Service) loadShowChildTags() {
	// default use android regions TODO
	regionkey := fmt.Sprintf(_initRegionKey, model.PlatAndroid, _initlanguage)
	res := s.regionCache[regionkey]
	var (
		// tag tmp
		similarTagTmp     = map[int64][]*tag.SimilarTag{}  // tmp tags
		tagsDetailTmp     = map[int64][]*region.ShowItem{} //tmp tags Detail
		tagsDetailOseaTmp = map[int64][]*region.ShowItem{} //tmp tags osea Detail
		tagsDetailAidsTmp = map[int64][]int64{}            //tmp tags aids Det
		// tags change ranking detail
		tagsDetailRankingTmp     = map[string][]*region.ShowItem{} //tmp tags Detail
		tagsDetailOseaRankingTmp = map[string][]*region.ShowItem{} //tmp tags osea Detail
		tagsDetailAidsRankingTmp = map[string][]int64{}            //tmp tags aids Detail
	)
	for _, v := range res {
		if v.Reid != 0 {
			var (
				tmp, tmpOsea []*region.ShowItem
				aidsTmp      []int64
			)
			//tag
			if tinfos, ok := s.regionTagCache[v.Rid]; ok && len(tinfos) > 0 {
				for _, tinfo := range tinfos {
					tid := tinfo.TagId
					if _, ok := similarTagTmp[tid]; !ok {
						tmpStag := s.loadSimilarTag(tid)
						similarTagTmp[tid] = s.upTagCache(tmpStag, s.similarTagCache[tid])
						tmp, tmpOsea, aidsTmp = s.loadTagDetail(tid)
						tagsDetailTmp[tid], tagsDetailOseaTmp[tid] = s.upCache(tmp, tmpOsea, s.tagsDetailCache[tid], s.tagsDetailOseaCache[tid])
						tagsDetailAidsTmp[tid] = s.upAidsCache(aidsTmp, s.tagsDetailAidsCache[tid])
					}
					tagRankKey := fmt.Sprintf(_initRegionTagKey, v.Reid, tid)
					if _, ok := tagsDetailRankingTmp[tagRankKey]; !ok {
						tmp, tmpOsea, aidsTmp = s.loadTagDetailRanking(v.Reid, tid)
						tagsDetailRankingTmp[tagRankKey], tagsDetailOseaRankingTmp[tagRankKey] = s.upCache(tmp, tmpOsea, s.tagsDetailRankingCache[tagRankKey], s.tagsDetailRankingOseaCache[tagRankKey])
						tagsDetailAidsRankingTmp[tagRankKey] = s.upAidsCache(aidsTmp, s.tagsDetailRankingAidsCache[tagRankKey])
					}
				}
			}
		}
	}
	// tag
	s.similarTagCache = similarTagTmp
	s.tagsDetailCache = tagsDetailTmp
	s.tagsDetailOseaCache = tagsDetailOseaTmp
	s.tagsDetailAidsCache = tagsDetailAidsTmp
	// tags change ranking detail
	s.tagsDetailRankingCache = tagsDetailRankingTmp
	s.tagsDetailRankingOseaCache = tagsDetailOseaRankingTmp
	s.tagsDetailRankingAidsCache = tagsDetailAidsRankingTmp
}

func (s *Service) loadShowChildTagsInfo() {
	// default use android regions TODO
	regionkey := fmt.Sprintf(_initRegionKey, model.PlatAndroid, _initlanguage)
	res := s.regionCache[regionkey]
	reslist := s.regionListCache[regionkey]
	var (
		// tag tmp
		tagsRegionTmp = map[int][]*tag.SimilarTag{} // region tags
		tagsTmp       = map[string]string{}         // tagid cache
		tagsNameTmp   = map[string]int64{}
	)
	for _, v := range res {
		if v.Reid != 0 {
			//tag
			var rTmp *region.Region
			if r, ok := reslist[v.Reid]; ok {
				rTmp = r
			}
			if tids := s.loadTagIDs(v.Rid); len(tids) > 0 {
				for _, t := range tids {
					tagInfo := &tag.SimilarTag{
						TagId:   t.ID,
						TagName: t.Name,
						Rid:     v.Rid,
						Rname:   v.Name,
					}
					if rTmp != nil {
						tagInfo.Reid = rTmp.Rid
						tagInfo.Rename = rTmp.Name
					}
					//tags info
					tagsRegionTmp[v.Rid] = append(tagsRegionTmp[v.Rid], tagInfo)
					key := fmt.Sprintf(_initRegionTagKey, v.Rid, t.ID)
					tagsTmp[key] = t.Name
					tgkey := fmt.Sprintf(_initTagNameKey, t.Name)
					tagsNameTmp[tgkey] = t.ID
				}
			}
		}
	}
	// region child aids
	s.regionTagCache = tagsRegionTmp
	s.tagsCache = tagsTmp
	s.tagsNameCache = tagsNameTmp
}

// loadRegion regions cache.
func (s *Service) loadRegion() {
	res, err := s.regiondao.All(context.TODO())
	if err != nil {
		log.Error("s.regiondao.All error(%v)", err)
		return
	}
	tmp := map[string][]*region.Region{}
	tmpRegion := map[string]map[int]*region.Region{}
	for _, v := range res {
		key := fmt.Sprintf(_initRegionKey, v.Plat, v.Language)
		tmp[key] = append(tmp[key], v)
	}
	if tmp != nil && tmpRegion != nil {
		s.regionCache = tmp
	}
	log.Info("region cacheproc success")
}

// loadShowChildTagIDs
func (s *Service) loadTagIDs(rid int) (tags []*tag.Tag) {
	tags, err := s.tg.TagHotsId(context.TODO(), rid, time.Now())
	if err != nil || len(tags) == 0 {
		log.Error("s.tag.loadShowChildTagIDs(%d) error(%v)", rid, err)
		return
	}
	return
}

// loadSimilarTag
func (s *Service) loadSimilarTag(tid int64) (res []*tag.SimilarTag) {
	res, err := s.tg.SimilarTagChange(context.TODO(), tid, time.Now())
	if err != nil || len(res) == 0 {
		log.Error("s.tag.loadSimilarChangeTag(%d) error(%v)", tid, err)
		return
	}
	log.Info("loadSimilarChangeTag(%d) success", tid)
	return
}

// loadShowChileNewTagHot
func (s *Service) loadTagDetail(tid int64) (resData, resOseaData []*region.ShowItem, arcAids []int64) {
	const (
		strtNum = 1
		maxNum  = 50
		newNum  = 20
	)
	arcAids, err := s.tg.Detail(context.TODO(), tid, strtNum, maxNum, time.Now())
	if err != nil || len(arcAids) < 20 {
		log.Error("s.tg.Detail(%d) error(%v)", tid, err)
		return
	}
	if len(arcAids) > newNum {
		arcAids = arcAids[:newNum]
	}
	resData, resOseaData = s.fromAids(context.TODO(), arcAids, false, 0)
	log.Info("loadTagDetail(%d) success", tid)
	return
}

// loadShowChileNewTagHot
func (s *Service) loadTagDetailRanking(rid int, tid int64) (resData, resOseaData []*region.ShowItem, arcAids []int64) {
	const (
		strtNum = 1
		maxNum  = 50
		newNum  = 20
	)
	arcAids, err := s.tg.DetailRanking(context.TODO(), rid, tid, strtNum, maxNum, time.Now())
	if err != nil || len(arcAids) < 20 {
		log.Error("s.tg.DetailRanking(%d) error(%v)", tid, err)
		return
	}
	if len(arcAids) > newNum {
		arcAids = arcAids[:newNum]
	}
	resData, resOseaData = s.fromAids(context.TODO(), arcAids, false, 0)
	log.Info("loadTagDetailRanking(%d) success", tid)
	return
}

// upTagCache
func (s *Service) upTagCache(new, old []*tag.SimilarTag) (res []*tag.SimilarTag) {
	if len(new) > 0 {
		res = new
	} else {
		res = old
	}
	return
}

// upCahce update cache
func (s *Service) upCache(new, newOsea, old, oldOsea []*region.ShowItem) (res, resOsea []*region.ShowItem) {
	if len(new) > 0 {
		res = new
	} else {
		res = old
	}
	if len(newOsea) > 0 {
		resOsea = newOsea
	} else {
		resOsea = oldOsea
	}
	return
}

// upAidsCache update aids  cache
func (s *Service) upAidsCache(new, old []int64) (aids []int64) {
	if len(new) > 0 {
		aids = new
	} else {
		aids = old
	}
	return
}
