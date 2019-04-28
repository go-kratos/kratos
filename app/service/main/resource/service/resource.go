package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/service/main/resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_archiveURL  = "https://www.bilibili.com/video/av"
	_archiveURL2 = "http://www.bilibili.com/video/av"
)

// LoadRes load Res info to cache
func (s *Service) loadRes() (err error) {
	// load default banner.
	var posTmp map[int][]int
	resTmp, err := s.res.Resources(context.TODO())
	if err != nil {
		log.Error("s.res.Resources error(%v)", err)
		return
	}
	s.resCache = resTmp
	resCacheMap := make(map[int]*model.Resource)
	posTmp = make(map[int][]int)
	for _, res := range resTmp {
		resCacheMap[res.ID] = res
		if res.Counter == 0 && res.Parent != 0 {
			posTmp[res.Parent] = append(posTmp[res.Parent], res.ID)
		}
	}
	s.posCache = posTmp
	s.resCacheMap = resCacheMap
	// load default banner.
	asgTmp, err := s.res.Assignment(context.TODO())
	if err != nil {
		log.Error("s.res.Assignment error(%v)", err)
		return
	}
	asgNewTmp, err := s.res.AssignmentNew(context.TODO())
	if err != nil {
		log.Error("s.res.AssignmentNew error(%v)", err)
		return
	}
	asgNewTmp = append(asgNewTmp, asgTmp...)
	categoryTmp, err := s.res.CategoryAssignment(context.TODO())
	if err != nil {
		log.Error("s.res.CategoryAssignment error(%v)", err)
		return
	}
	asgNewTmp = append(asgNewTmp, categoryTmp...)
	s.asgCache = asgNewTmp
	resArchiveWarn, resURLWarn := s.formWarnInfo(asgNewTmp)
	s.resArchiveWarnCache = resArchiveWarn
	s.resURLWarnCache = resURLWarn
	asgCacheMap := make(map[int][]*model.Assignment)
	for _, asg := range asgNewTmp {
		asgCacheMap[asg.ResID] = append(asgCacheMap[asg.ResID], asg)
	}
	s.asgCacheMap = asgCacheMap
	// load default banner.
	bannerTmp, err := s.res.DefaultBanner(context.TODO())
	if err != nil {
		log.Error("s.res.DefaultBanner error(%v)", err)
		return
	}
	s.defBannerCache = bannerTmp
	// index icon
	tmpIndexIcon, err := s.res.IndexIcon(context.TODO())
	if err != nil {
		log.Error("s.res.IndexIcon() error(%v)", err)
		return
	}
	s.indexIcon = tmpIndexIcon
	return
}

func (s *Service) formWarnInfo(asgNewTmp []*model.Assignment) (resArchive map[int64][]*model.ResWarnInfo, resURL map[string][]*model.ResWarnInfo) {
	resArchive = make(map[int64][]*model.ResWarnInfo)
	resURL = make(map[string][]*model.ResWarnInfo)
	for _, asg := range asgNewTmp {
		var (
			aid int64
			url string
			err error
			rw  *model.ResWarnInfo
		)
		if (asg.Atype == model.AsgTypeVideo) || (asg.Atype == model.AsgTypeAv) {
			if aid, err = strconv.ParseInt(asg.URL, 10, 64); err != nil {
				log.Error("formWarnInfo url(%v) error(%v)", asg.URL, err)
				err = nil
				continue
			}
		} else if (asg.Atype == model.AsgTypePic) || (asg.Atype == model.AsgTypeURL) {
			if strings.HasPrefix(asg.URL, _archiveURL) {
				urls := strings.Split(asg.URL, "?")
				aidURL := strings.TrimPrefix(urls[0], _archiveURL)
				aidURL = strings.TrimSuffix(aidURL, "/")
				if aid, err = strconv.ParseInt(aidURL, 10, 64); err != nil {
					log.Error("formWarnInfo url(%v) error(%v)", asg.URL, aidURL, err)
					err = nil
					continue
				}
			} else if strings.HasPrefix(asg.URL, _archiveURL2) {
				urls := strings.Split(asg.URL, "?")
				aidURL := strings.TrimPrefix(urls[0], _archiveURL2)
				aidURL = strings.TrimSuffix(aidURL, "/")
				if aid, err = strconv.ParseInt(aidURL, 10, 64); err != nil {
					log.Error("formWarnInfo url(%v) error(%v)", asg.URL, err)
					err = nil
					continue
				}
			} else {
				url = asg.URL
			}
		}
		if aid == 0 && url == "" {
			continue
		}
		rw = &model.ResWarnInfo{
			AssignmentID:   asg.AsgID,
			AssignmentName: asg.Name,
			STime:          asg.STime,
			ETime:          asg.ETime,
			UserName:       asg.Username,
			ApplyGroupID:   asg.ApplyGroupID,
			MaterialID:     asg.ID,
		}
		if re, ok := s.resCacheMap[asg.ResID]; ok {
			if re.Counter > 0 {
				rw.ResourceID = re.ID
				if rep, ok := s.resCacheMap[re.Parent]; ok {
					rw.ResourceName = fmt.Sprintf("%v_%v", rep.Name, re.Name)
					continue
				}
				rw.ResourceName = re.Name
			} else {
				rw.ResourceID = re.Parent
				rw.ResourceName = re.Name
			}
		}
		if aid != 0 {
			rw.AID = aid
			resArchive[aid] = append(resArchive[aid], rw)
		} else {
			rw.URL = url
			resURL[url] = append(resURL[url], rw)
		}
	}
	return
}

// ResourceAll get all resource
func (s *Service) ResourceAll(c context.Context) (res []*model.Resource) {
	res = s.resCache
	return
}

// AssignmentAll get all assignment
func (s *Service) AssignmentAll(c context.Context) (ass []*model.Assignment) {
	// TODO delete
	for _, asc := range s.asgCache {
		as := &model.Assignment{}
		*as = *asc
		as.Weight = 0
		ass = append(ass, as)
	}
	return
}

// Resource get resource by resource_id or positon_id
func (s *Service) Resource(c context.Context, resID int) (res *model.Resource) {
	var (
		ok  bool
		pos []int
	)
	if res, ok = s.resCacheMap[resID]; !ok {
		return
	}
	// Safe first!! Prevent res nil panic.
	if res.Counter == 0 {
		if len(s.asgCacheMap[resID]) > 0 {
			res.Assignments = s.asgCacheMap[resID]
			return
		}
		res.Assignments = s.asgCacheMap[res.Parent]
	} else {
		if pos, ok = s.posCache[resID]; !ok {
			return
		}
		var (
			tmpNormalRes   []*model.Assignment
			tmpCategoryRes = s.asgCacheMap[resID]
		)
		for _, pid := range pos {
			tmpNormalRes = append(tmpNormalRes, s.asgCacheMap[pid]...)
		}
		for _, nr := range tmpNormalRes {
			if nr.Weight > len(tmpCategoryRes) {
				tmpCategoryRes = append(tmpCategoryRes, nr)
			} else {
				tmpCategoryRes = append(tmpCategoryRes[:nr.Weight-1], append([]*model.Assignment{nr}, tmpCategoryRes[nr.Weight-1:]...)...)
			}
		}
		if len(tmpCategoryRes) > res.Counter {
			res.Assignments = tmpCategoryRes[:res.Counter]
		} else {
			res.Assignments = tmpCategoryRes
		}
	}
	return
}

// Resources get resources by resource_ids or position_ids
func (s *Service) Resources(c context.Context, resIDs []int) (res map[int]*model.Resource) {
	if len(resIDs) == 0 {
		res = _emptyResources
		return
	}
	res = make(map[int]*model.Resource)
	for _, rid := range resIDs {
		if resTmp := s.Resource(c, rid); resTmp != nil {
			res[rid] = resTmp
		}
	}
	return
}

// DefBanner get defbanner config
func (s *Service) DefBanner(c context.Context) (defbanner *model.Assignment) {
	defbanner = s.defBannerCache
	return
}

// IndexIcon get index icon
func (s *Service) IndexIcon(c context.Context) (icons map[string][]*model.IndexIcon) {
	icons = map[string][]*model.IndexIcon{
		model.IconTypes[model.IconTypeFix]:    s.indexIcon[model.IconTypeFix],
		model.IconTypes[model.IconTypeRandom]: s.indexIcon[model.IconTypeRandom],
	}
	return
}

// PlayerIcon get player icon
func (s *Service) PlayerIcon(c context.Context) (re *model.PlayerIcon, err error) {
	if re = s.playIcon; re == nil {
		err = ecode.NothingFound
	}
	return
}

// Cmtbox get live danmaku box
func (s *Service) Cmtbox(c context.Context, id int64) (re *model.Cmtbox, err error) {
	var ok bool
	if re, ok = s.cmtbox[id]; !ok {
		err = ecode.NothingFound
	}
	return
}
