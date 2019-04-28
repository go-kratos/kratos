package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/resource/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// loadVideoAds load Video_ads to cache.
func (s *Service) loadVideoAds() (err error) {
	var (
		ok        bool
		vdoAds    []*model.VideoAD
		vdoAdsAPP = make(map[int8]map[int8]map[int8]map[string]*model.VideoAD)
		tmpAIDm   []map[int64]int64
		tmpAIDs   = make(map[int64]int64)
	)
	if vdoAds, err = s.ads.VideoAds(context.TODO()); err != nil {
		log.Error("s.ads.VideoAds error(%v)", err)
		return
	}
	for _, ad := range vdoAds {
		var (
			tmpAPP  map[int8]map[int8]map[string]*model.VideoAD
			tmpAPP2 map[int8]map[string]*model.VideoAD
			tmpAPP3 map[string]*model.VideoAD
		)
		if tmpAPP, ok = vdoAdsAPP[ad.Platform]; !ok {
			tmpAPP = make(map[int8]map[int8]map[string]*model.VideoAD)
			vdoAdsAPP[ad.Platform] = tmpAPP
		}
		if tmpAPP2, ok = tmpAPP[ad.Type]; !ok {
			tmpAPP2 = make(map[int8]map[string]*model.VideoAD)
			tmpAPP[ad.Type] = tmpAPP2
		}
		if tmpAPP3, ok = tmpAPP2[ad.Target]; !ok {
			tmpAPP3 = make(map[string]*model.VideoAD)
			tmpAPP2[ad.Target] = tmpAPP3
		}
		switch ad.Target {
		case model.VdoAdsTargetArchive:
			if ad.Aids == "" {
				continue
			}
			xids := strings.Split(ad.Aids, ",")
			for _, xid := range xids {
				tmpAPP3[xid] = ad
				if ad.FrontAid > 0 {
					tmpAIDs[ad.FrontAid] = ad.FrontAid
				}
			}
		case model.VdoAdsTargetBangumi:
			if ad.SeasonID <= 0 {
				continue
			}
			sid := strconv.Itoa(ad.SeasonID)
			tmpAPP3[sid] = ad
			if ad.Platform == model.VdoAdsPC {
				if ad.AdCid > 0 {
					tmpAIDs[ad.AdCid] = ad.AdCid
				}
			} else {
				if ad.FrontAid > 0 {
					tmpAIDs[ad.FrontAid] = ad.FrontAid
				}
			}
		case model.VdoAdsTargetType:
			if ad.TypeID <= 0 {
				continue
			}
			tid := strconv.Itoa(int(ad.TypeID))
			tmpAPP3[tid] = ad
			if ad.FrontAid > 0 {
				tmpAIDs[ad.FrontAid] = ad.FrontAid
			}
		}
		if len(tmpAIDs) == 50 {
			tmpAIDm = append(tmpAIDm, tmpAIDs)
			tmpAIDs = make(map[int64]int64)
		}
	}
	if len(tmpAIDs) > 0 {
		tmpAIDm = append(tmpAIDm, tmpAIDs)
	}
	s.videoAdsAPPCache = vdoAdsAPP
	s.PasterAIDCache = tmpAIDm
	return
}

// PasterAPP get paster for app nologin
func (s *Service) PasterAPP(c context.Context, plat, adType int8, aid, typeID, buvid string) (res *model.Paster, err error) {
	var (
		vdoApp  map[int8]map[int8]map[string]*model.VideoAD
		vdoApp2 map[int8]map[string]*model.VideoAD
		vdoApp3 map[string]*model.VideoAD
		vdoApp4 *model.VideoAD
		ok      bool
	)
	platform := model.PasterPlat(int8(plat))
	res = new(model.Paster)
	if vdoApp, ok = s.videoAdsAPPCache[platform]; !ok {
		return
	}
	if vdoApp2, ok = vdoApp[adType]; !ok {
		return
	}
	var (
		pages                          []*api.Page
		faid, playCount, confPlayCount int64
	)
	// aid first.
	if vdoApp3, ok = vdoApp2[model.VdoAdsTargetArchive]; ok && vdoApp3[aid] != nil {
		faid = vdoApp3[aid].FrontAid
		confPlayCount = vdoApp3[aid].PlayCount
	}
	if faid <= 0 {
		if vdoApp3, ok = vdoApp2[model.VdoAdsTargetType]; !ok {
			return
		}
		if vdoApp4, ok = vdoApp3[typeID]; !ok && len(s.typeList) > 0 {
			pid := s.typeList[typeID]
			if pid != "" && pid != "0" {
				vdoApp4 = vdoApp3[pid]
			}
		}
		if vdoApp4 == nil {
			return
		}
		if vdoApp4.FrontAid <= 0 {
			return
		}
		faid = vdoApp4.FrontAid
		confPlayCount = vdoApp4.PlayCount
	}
	// check buvid count.
	if playCount, err = s.ads.BuvidCount(c, faid, buvid); err != nil || playCount == confPlayCount {
		// skip ad when redis error or playCount reached
		return
	}
	if pages, err = s.arcRPC.Page3(c, &arcmdl.ArgAid2{Aid: faid}); err != nil || len(pages) == 0 {
		// skip ad when av error
		log.Error("s.arcRPC.Page3(%d) error(%v)", faid, err)
		return
	}
	res.AID = faid
	res.CID = pages[0].Cid
	res.Duration = pages[0].Duration
	res.Type = adType
	res.AllowJump = vdoApp4.Skipable
	// update buvid count.
	bcCache := map[string]map[int64]int64{
		buvid: map[int64]int64{
			faid: playCount + 1,
		},
	}
	s.addCache(bcCache)
	return
}

// PasterPGC get paster for pgc
func (s *Service) PasterPGC(c context.Context, plat, adType int8, sid string) (res *model.Paster, err error) {
	var (
		vdoApp  map[int8]map[int8]map[string]*model.VideoAD
		vdoApp2 map[int8]map[string]*model.VideoAD
		vdoApp3 map[string]*model.VideoAD
		ok      bool
	)
	platform := model.PasterPlat(int8(plat))
	res = new(model.Paster)
	if vdoApp, ok = s.videoAdsAPPCache[platform]; !ok {
		return
	}
	if vdoApp2, ok = vdoApp[adType]; !ok {
		return
	}
	var (
		pages []*api.Page
		faid  int64
	)
	// aid first.
	if vdoApp3, ok = vdoApp2[model.VdoAdsTargetBangumi]; !ok || vdoApp3[sid] == nil {
		return
	}
	if platform == model.VdoAdsPC {
		faid = vdoApp3[sid].AdCid
	} else {
		faid = vdoApp3[sid].FrontAid
	}
	if pages, err = s.arcRPC.Page3(c, &arcmdl.ArgAid2{Aid: faid}); err != nil || len(pages) == 0 {
		// skip ad when av error
		log.Error("PasterPGC s.arcRPC.Page3(%d) error(%v)", faid, err)
		return
	}
	res.AID = faid
	res.CID = pages[0].Cid
	res.Duration = pages[0].Duration
	res.Type = adType
	res.AllowJump = vdoApp3[sid].Skipable
	if platform == model.VdoAdsPC {
		url := vdoApp3[sid].AdURL
		if url != "" {
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = fmt.Sprintf("http://%v", url)
			}
		}
		res.URL = url
	}
	return
}

// PasterCID get all Paster's cid.
func (s *Service) PasterCID(c context.Context) (res map[int64]int64, err error) {
	var mutex = sync.Mutex{}
	res = make(map[int64]int64)
	g, ctx := errgroup.WithContext(c)
	for _, aidm := range s.PasterAIDCache {
		var aids []int64
		for aid, _ := range aidm {
			aids = append(aids, aid)
		}
		if len(aids) > 0 {
			g.Go(func() (err error) {
				arcs, err := s.arcRPC.Archives3(ctx, &arcmdl.ArgAids2{Aids: aids})
				if err != nil {
					log.Error("%v", err)
					err = nil
					return
				}
				for _, arc := range arcs {
					mutex.Lock()
					res[arc.FirstCid] = arc.Aid
					mutex.Unlock()
				}
				return
			})
		}
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	return
}
