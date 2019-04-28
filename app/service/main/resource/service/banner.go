package service

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	locmdl "go-common/app/service/main/location/model"
	"go-common/app/service/main/resource/model"
	"go-common/library/log"

	farm "github.com/dgryski/go-farm"
)

// loadBannerCahce load banner cache.
func (s *Service) loadBannerCahce() (err error) {
	// get all banners
	nbs, err := s.res.Banner(context.TODO())
	if err != nil {
		log.Error("s.res.Banner error(%v)", err)
		return
	}
	s.bannerCache = nbs
	log.Info("load bannerCache success")
	var bannerHashTpm = make(map[int8]string, len(nbs))
	for plat, bnnr := range nbs {
		bannerHashTpm[plat] = hash(bnnr)
	}
	s.bannerHashCache = bannerHashTpm
	log.Info("load BannerHashCache success")
	cbs, err := s.res.Category(context.TODO())
	if err != nil {
		log.Error("s.res.Category error(%v)", err)
		return
	}
	s.categoryBannerCache = cbs
	log.Info("load categoryBannerCache success")
	// banner limit
	limit, err := s.res.Limit(context.TODO())
	if err != nil {
		log.Error("s.dao.Limit error(%v)", err)
		return
	}
	s.bannerLimitCache = limit
	log.Info("load BannerLimitCache success")
	return
}

// hash get banner hash.
func hash(v map[int][]*model.Banner) (value string) {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}
	value = strconv.FormatUint(farm.Hash64(bs), 10)
	return
}

// Banners get banners by plat, build channel, ip for app-feed.
func (s *Service) Banners(c context.Context, plat int8, build int, aid, mid int64, resIdsStr, channel, ip, buvid, network, mobiApp, device, openEvent, adExtra, version string, isAd bool) (res *model.Banners) {
	res = &model.Banners{}
	if version != "" && (version == s.bannerHashCache[plat]) {
		log.Warn("Banners() plat(%v) version(%v) same as hash cache, return nil", plat, version)
		return
	}
	var (
		cpmResBus map[int]map[int]*model.Banner
		resIds    []string
		banner    map[int][]*model.Banner
	)
	if isAd {
		cpmResBus = s.cpmBanners(c, aid, mid, build, resIdsStr, mobiApp, device, buvid, network, ip, openEvent, adExtra)
	}
	resIds = strings.Split(resIdsStr, ",")
	banner = map[int][]*model.Banner{}
	for _, resIDStr := range resIds {
		if resIDStr == "" {
			continue
		}
		resID, err := strconv.Atoi(resIDStr)
		if err != nil {
			log.Warn("strconv.Atoi(%s) error(%v)", resIDStr, err)
			err = nil
			continue
		}
		var (
			resBs, cbcs, resAll []*model.Banner
			cbc                 = s.categoryBannerCache[plat] // operater category banner
			bArea               []string
			ok                  bool
			maxBannerIndex      int
		)
		if len(cbc) > 0 {
			if cbcs, ok = cbc[resID]; ok {
				btime := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
				for _, b := range cbcs {
					if s.filterBs(c, plat, build, channel, b) {
						continue
					}
					if b.Area != "" {
						bArea = append(bArea, b.Area)
					}
					tmp := &model.Banner{}
					*tmp = *b
					tmp.ServerType = 0
					tmp.RequestId = btime
					if tmp.Rank > maxBannerIndex {
						maxBannerIndex = tmp.Rank
					}
					resBs = append(resBs, tmp)
				}
			}
		}
		if len(resBs) > maxBannerIndex {
			maxBannerIndex = len(resBs)
		}
		var (
			cpmBus  map[int]*model.Banner // cpm ad
			allRank []int
			tmpCmps []*model.Banner
		)
		// append cpm banner
		if cpmBus, ok = cpmResBus[resID]; ok && len(cpmBus) > 0 {
			var cpmMs = map[int]*model.Banner{}
			for _, cpm := range cpmBus {
				if cpm.IsAdReplace {
					cpmMs[cpm.Rank] = cpm
					allRank = append(allRank, cpm.Rank)
					delete(cpmBus, cpm.Rank)
				}
			}
			if len(allRank) > 0 {
				sort.Ints(allRank)
				for _, key := range allRank {
					if cpmMs[key].Rank > maxBannerIndex {
						maxBannerIndex = cpmMs[key].Rank
					}
					tmpCmps = append(tmpCmps, cpmMs[key])
				}
			}
		}
		if (len(resBs) + len(tmpCmps)) > maxBannerIndex {
			maxBannerIndex = len(resBs) + len(tmpCmps)
		}
		var (
			plm         = s.bannerCache[plat] // operater normal banner
			tmpBs, plbs []*model.Banner
		)
		// append normal banner
		if len(plm) > 0 {
			if plbs, ok = plm[resID]; ok {
				btime := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
				for _, b := range plbs {
					if s.filterBs(c, plat, build, channel, b) {
						continue
					}
					if b.Area != "" {
						bArea = append(bArea, b.Area)
					}
					tmp := &model.Banner{}
					*tmp = *b
					tmp.ServerType = 0
					tmp.RequestId = btime
					if tmp.Rank > maxBannerIndex {
						maxBannerIndex = tmp.Rank
					}
					tmpBs = append(tmpBs, tmp)
				}
			}
		}
		if (len(resBs) + len(tmpCmps) + len(tmpBs)) > maxBannerIndex {
			maxBannerIndex = len(resBs) + len(tmpCmps) + len(tmpBs)
		}
		var tcIndex, tbIndex, cIndex int
		for i := 1; i <= maxBannerIndex; i++ {
			if tcIndex < len(tmpCmps) {
				tc := tmpCmps[tcIndex]
				if tc.Rank == i {
					resAll = append(resAll, tc)
					tcIndex++
					continue
				}
			}
			if tbIndex < len(tmpBs) {
				tb := tmpBs[tbIndex]
				if tb.Rank <= i {
					resAll = append(resAll, tb)
					tbIndex++
					continue
				}
			}
			if cIndex < len(resBs) {
				cb := resBs[cIndex]
				resAll = append(resAll, cb)
				cIndex++
			}
		}
		for i, b := range resAll {
			if cpm, ok := cpmBus[i+1]; ok && !b.IsAdReplace { // NOTE: surplus cpm is ad loc
				b.IsAdLoc = true
				b.IsAd = cpm.IsAd
				b.CmMark = cpm.CmMark
				b.SrcId = cpm.SrcId
				b.RequestId = cpm.RequestId
				b.ClientIp = cpm.ClientIp
			}
		}
		if max, ok := s.bannerLimitCache[resID]; ok && len(resAll) > max {
			resAll = resAll[:max]
		}
		for i := 0; i < len(resAll); i++ {
			resAll[i].Index = i + 1
			resAll[i].ResourceID = resID
		}
		if len(resAll) > 0 {
			var (
				auths  map[int64]*locmdl.Auth
				resBs2 []*model.Banner
			)
			if len(bArea) > 0 {
				if auths, err = s.locationRPC.AuthPIDs(c, &locmdl.ArgPids{Pids: strings.Join(bArea, ","), IP: ip}); err != nil {
					log.Error("%v", err)
					err = nil
				}
			}
			for _, resB := range resAll {
				if resB.Area != "" {
					var pid int64
					if pid, err = strconv.ParseInt(resB.Area, 10, 64); err != nil {
						log.Warn("banner strconv.ParseInt(%v) error(%v)", resB.Area, err)
						err = nil
					} else {
						if auth, ok := auths[pid]; ok && auth.Play == locmdl.Forbidden {
							log.Warn("resID(%v) pid(%v) ip(%v) in zone limit", resID, resB.Area, ip)
							continue
						}
					}
				}
				resBs2 = append(resBs2, resB)
			}
			banner[resID] = resBs2
		}
	}
	res.Banner = banner
	res.Version = s.bannerHashCache[plat]
	return
}

// cpmBanners
func (s *Service) cpmBanners(c context.Context, aid, mid int64, build int, resource, mobiApp, device, buvid, network, ipaddr, openEvent, adExtra string) (banners map[int]map[int]*model.Banner) {
	ipInfo, err := s.locationRPC.Info(c, &locmdl.ArgIP{IP: ipaddr})
	if err != nil || ipInfo == nil {
		log.Error("CpmsBanners s.locationRPC.Zone(%s) error(%v) or ipinfo is nil", ipaddr, err)
		ipInfo = &locmdl.Info{Addr: ipaddr}
	}
	adr, err := s.cpm.CpmsAPP(c, aid, mid, build, resource, mobiApp, device, buvid, network, openEvent, adExtra, ipInfo)
	if err != nil || adr == nil {
		log.Error("s.ad.ADRequest error(%v)", err)
		return
	}
	banners = adr.ConvertBanner(ipInfo.Addr, mobiApp, build)
	return
}

// filterBs filter banner.
func (s *Service) filterBs(c context.Context, plat int8, build int, channel string, b *model.Banner) bool {
	if model.InvalidBuild(build, b.Build, b.Condition) {
		return true
	}
	if model.InvalidChannel(plat, channel, b.Channel) && b.Channel != "" {
		return true
	}
	return false
}
