package show

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/show"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	resource "go-common/app/service/main/resource/model"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_cnt              = 4
	_initShowKey      = "show_key_%d_%v"
	_initCardKey      = "card_key_%d"
	_initlanguage     = "hans"
	_bangumiSeasonID  = 1
	_bangumiEpisodeID = 2
)

var (
	_emptyShow      = []*show.Show{}
	_emptyItem      = &show.Item{}
	_emptyShowItems = []*show.Item{}
	// ad
	_recommend = map[int8]string{
		model.PlatIPhone:   "1508",
		model.PlatAndroid:  "1515",
		model.PlatIPad:     "1522",
		model.PlatIPhoneI:  "1529",
		model.PlatAndroidG: "1543",
		model.PlatAndroidI: "1777",
		model.PlatIPadI:    "1536",
	}
	_bangumiReids = map[int]struct{}{
		167: struct{}{},
	}
)

// Display display show data.
func (s *Service) Display(c context.Context, mid int64, plat int8, build int, buvid, channel, ip, ak, network, mobiApp,
	device, language, adExtra string, isTmp bool, now time.Time) (res []*show.Show) {
	res = s.showDisplay(c, mid, plat, build, buvid, channel, ip, ak, network, mobiApp, device, language, adExtra, isTmp, false, false, now)
	return
}

// RegionDisplay display region show data.
func (s *Service) RegionDisplay(c context.Context, mid int64, plat int8, build int, buvid, channel, ip, ak, network, mobiApp,
	device, language, adExtra string, isTmp bool, now time.Time) (res []*show.Show) {
	res = s.showDisplay(c, mid, plat, build, buvid, channel, ip, ak, network, mobiApp, device, language, adExtra, isTmp, true, false, now)
	return
}

func (s *Service) Index(c context.Context, mid int64, plat int8, build int, buvid, channel, ip, ak, network, mobiApp,
	device, language, adExtra string, isTmp bool, now time.Time) (res []*show.Show) {
	res = s.showDisplay(c, mid, plat, build, buvid, channel, ip, ak, network, mobiApp, device, language, adExtra, isTmp, true, true, now)
	if cards := s.showCardDisplay(plat, build); len(cards) > 0 {
		cards = append(cards, res...)
		res = cards
	}
	return
}

// Display display show data.
func (s *Service) showDisplay(c context.Context, mid int64, plat int8, build int, buvid, channel, ip, ak, network, mobiApp,
	device, language, adExtra string, isTmp, isRegion, isIndex bool, now time.Time) (res []*show.Show) {
	var (
		bnr            string
		banners        map[int][]*resource.Banner
		showRec        []*show.Item
		showLive       []*show.Item
		isBangumi      = false
		isRegionBanner = false
		ss             []*show.Show
		resIDStr       = _bannersPlat[plat]
	)
	if language == "" {
		language = _initlanguage
	}
	key := fmt.Sprintf(_initShowKey, plat, language)
	if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
		ss = s.cacheBgEp[key]
	} else if ((mobiApp == "iphone" && build > 5600) || (mobiApp == "android" && build > 507000)) && isIndex {
		ss = s.cacheBg[key]
	} else {
		ss = s.cache[key]
	}
	if isTmp {
		ss = s.tempCache[key]
	}
	if len(ss) == 0 {
		res = _emptyShow
		return
	}
	res = make([]*show.Show, 0, len(ss))
	if (mobiApp == "iphone" && build > 4310) || (mobiApp == "android" && build > 502000) || isIndex {
		isBangumi = true
	}
	if (mobiApp == "iphone" && build > 4350) || (mobiApp == "android" && build > 503000) {
		isRegionBanner = true
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() error {
		banners = s.resBanners(ctx, plat, build, mid, resIDStr, channel, ip, buvid, network, mobiApp, device, adExtra)
		return nil
	})
	if !isRegion {
		g.Go(func() error {
			showRec = s.getRecommend(ctx, mid, build, plat, buvid, network, mobiApp, device, ip)
			return nil
		})
		g.Go(func() error {
			showLive = s.getLive(ctx, mid, ak, ip, 0, now)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Error("showDisplay errgroup.WithContext error(%v)", err)
	}
	for i, sw := range ss {
		if mobiApp == "white" && 101220 >= build && sw.Param == "165" { // 165 ad region
			continue
		} else if sw.Param != "165" || ((mobiApp != "iphone" || device != "pad") || build <= 3590) {
			if model.InvalidBuild(build, sw.Build, sw.Condition) {
				continue
			}
		}
		if sw.Type == "recommend" {
			if isRegion {
				continue
			}
			sw = s.dealRecommend(c, sw, plat, mid, build, buvid, network, mobiApp, device, ip, showRec)
			bnr = "0"
		} else if sw.Type == "live" {
			if isRegion {
				continue
			}
			sw = s.dealLive(c, sw, showLive)
			bnr = "65537"
		} else if sw.Type == "bangumi" {
			if ok := s.auditRegion(mobiApp, plat, build, "13"); ok {
				continue
			}
			if isRegion && isBangumi && !isRegionBanner {
				bnr = "-1"
			} else if isRegion && !isBangumi && !isRegionBanner {
				continue
			} else {
				bnr = "13"
			}
		} else {
			bnr = sw.Param
			if isRegion {
				if ok := s.auditRegion(mobiApp, plat, build, sw.Param); ok {
					continue
				}
				if !isRegionBanner {
					if sw.Param == "1" && !isBangumi {
						bnr = "-1"
					}
					if sw.Type == "topic" && i > 0 && (ss[i-1].Type == "bangumi" || ss[i-1].Type == "1") {
						continue
					}
				}
			}
		}
		sw.Banner = s.getBanners(c, plat, build, bnr, channel, ip, banners, isIndex)
		res = append(res, sw)
	}
	return
}

// showCardDisplay
func (s *Service) showCardDisplay(plat int8, build int) (res []*show.Show) {
	var ss []*show.Show
	key := fmt.Sprintf(_initCardKey, plat)
	ss = s.cardCache[key]
	if len(ss) == 0 {
		res = _emptyShow
		return
	}
	res = []*show.Show{}
	for _, sw := range ss {
		if model.InvalidBuild(build, sw.Build, sw.Condition) {
			continue
		}
		tmp := &show.Show{}
		*tmp = *sw
		tmp.FillBuildURI(plat, build)
		res = append(res, tmp)
	}
	return
}

// Change change display show data.
func (s *Service) Change(c context.Context, mid int64, build int, plat int8, rand int, buvid, ip, network, mobiApp, device string) (sis []*show.Item) {
	cnt := s.itemNum(plat)
	// first get recommend data.
	tmp := s.userRecommend(c, mid, build, plat, buvid, network, mobiApp, device, ip, cnt)
	if len(tmp) == cnt {
		sis = append(sis, tmp...)
	}
	if len(sis) < cnt {
		start := cnt * rand
		end := start + cnt
		rcLen := len(s.rcmmndCache)
		if rcLen < end {
			rand = 0
			start = cnt * rand
			end = start + cnt
		}
		if rcLen > end {
			sis = s.rcmmndCache[start:end]
		}
	}
	return
}

// RegionChange change show region data.
func (s *Service) RegionChange(c context.Context, rid, rand int, plat int8, build int, mobiApp string) (sis []*show.Item) {
	if rand < 0 {
		rand = 0
	}
	var (
		cnt         = 4
		pn          = rand + 1
		isOsea      = model.IsOverseas(plat)
		bangumiType = 0
		tmp         []*show.Item
	)
	if (mobiApp == "iphone" && build > 5600) || (mobiApp == "android" && build > 507000) {
		if _, isBangumi := _bangumiReids[rid]; isBangumi {
			if (plat == model.PlatIPhone && build > 6050) || (plat == model.PlatAndroid && build > 512007) {
				bangumiType = _bangumiEpisodeID
			} else {
				bangumiType = _bangumiSeasonID
			}
		}
	}
	if model.IsIPad(plat) {
		cnt = 8
	}
	as, aids, err := s.dyn.RegionDynamic(c, rid, pn, cnt)
	if err != nil {
		log.Error("s.rcmmnd.RegionDynamic(%d, %d, %d) error(%v)", rid, pn, cnt, err)
		sis = []*show.Item{}
		return
	}
	if bangumiType != 0 {
		tmp = s.fromArchivesBangumiOsea(c, as, aids, isOsea, bangumiType)
	} else {
		tmp = s.fromArchivesOsea(as, isOsea)
	}
	sis = append(sis, tmp...)
	return
}

// BangumiChange change show bangumi data.
func (s *Service) BangumiChange(c context.Context, rand int, plat int8) (sis []*show.Item) {
	if rand < 0 {
		rand = 0
	}
	rand = rand + 1
	var (
		cnt = 4
	)
	if model.IsIPad(plat) {
		cnt = 8
	}
	start := cnt * rand
	end := start + cnt
	if bgms, ok := s.bgmCache[plat]; ok {
		bcLen := len(bgms)
		if bcLen < end {
			rand = 0
			start = cnt * rand
			end = start + cnt
		}
		if bcLen > end {
			sis = bgms[start:end]
		}
	}
	return
}

// Dislike dislike show data
func (s *Service) Dislike(c context.Context, mid int64, plat int8, id int64, buvid, mobiApp, device, gt, ip string) (si *show.Item) {
	var (
		cnt       = 1
		changeAid string
		port      string
	)
	// first get recommend data.
	tmp := s.userRecommend(c, mid, 0, plat, buvid, "", mobiApp, device, ip, cnt)
	if len(tmp) > 0 {
		si = tmp[0]
		port = "userRecommend"
	} else {
		si = s.rcmmndCache[0]
		port = "loadRcmmndCache"
	}
	if si != nil {
		changeAid = si.Param
	}
	if err := s.dbus.Pub(c, buvid, gt, id, mid); err != nil {
		log.Error("s.dbus.Pub(%s,%s,%d,%d) error(%v)", buvid, gt, id, mid, err)
		log.Error("dbus_Pub_dislike error mid:%v , dislike_aid:%v , change_aid:%v , interface_name:%v", mid, id, changeAid, port)
		return
	}
	log.Info("dbus_Pub_dislike success mid:%v , dislike_aid:%v , change_aid:%v , interface_name:%v", mid, id, changeAid, port)
	return
}

// Widget
func (s *Service) Widget(c context.Context, plat int8) (res []*show.Item) {
	var (
		isOsea   = model.IsOverseas(plat) //is overseas
		resCache []*show.Item
		randID   int
	)
	if isOsea {
		resCache = s.rcmmndOseaCache
	} else {
		resCache = s.rcmmndCache
	}
	resCacheLen := len(resCache)
	if resCacheLen >= 3 {
		for {
			if len(res) >= 3 || len(resCache) == 0 {
				log.Info("Widget len 3")
				break
			}
			if randInt := rand.Intn(resCacheLen); randInt != randID && resCache[randInt] != nil {
				randID = randInt
				res = append(res, resCache[randInt])
			}
		}
	} else if resCacheLen > 0 {
		log.Info("Widget resCache")
		res = resCache
	} else {
		log.Info("Widget is null")
		res = _emptyShowItems
	}
	return
}

// LiveChange live change.
func (s *Service) LiveChange(c context.Context, mid int64, ak, ip string, rand int, now time.Time) (sis []*show.Item) {
	return s.getLive(c, mid, ak, ip, rand, now)
}

// dealRecommend deal recommend.
func (s *Service) dealRecommend(c context.Context, sw *show.Show, plat int8, mid int64, build int, buvid, network, mobiApp, device, ipaddr string, showRec []*show.Item) (rs *show.Show) {
	cnt := s.itemNum(plat)
	sis := make([]*show.Item, 0, cnt)
	// first get recommend data.
	if len(showRec) == cnt {
		sis = append(sis, showRec...)
	}
	// if recommend data not enough, get from @hetongzi.
	if len(sis) < cnt {
		rcLen := len(s.rcmmndCache)
		if rcLen < cnt {
			sis = s.rcmmndCache[0:rcLen]
		} else {
			sis = s.rcmmndCache[0:cnt]
		}
		if rcLen > 0 {
			sis = s.adVideo(c, mid, build, plat, buvid, network, mobiApp, device, ipaddr, sis)
		}
	}
	if len(sis) == 0 {
		sis = []*show.Item{}
	}
	rs = &show.Show{
		Head: sw.Head,
		Body: sis,
	}
	return
}

// getRecommend user recommend data
func (s *Service) getRecommend(c context.Context, mid int64, build int, plat int8, buvid, network, mobiApp, device, ipaddr string) (sis []*show.Item) {
	cnt := s.itemNum(plat)
	// first get recommend data.
	sis = s.userRecommend(c, mid, build, plat, buvid, network, mobiApp, device, ipaddr, cnt)
	return
}

// userRecommend user recommend data.
func (s *Service) userRecommend(ctx context.Context, mid int64, build int, plat int8, buvid, network, mobiApp, device, ipaddr string, cnt int) (sis []*show.Item) {
	// get redis seed whether or not hit
	if !s.rcmmndOn {
		return
	}
	var (
		key  = buvid
		i    int
		aids []int64
		rcs  []*rcmmndCfg
		err  error
	)
	if mid > 0 {
		key = strconv.FormatInt(mid, 10)
	}
	if key == "" {
		return
	}
Retry:
	for i = 0; i < 2; i++ {
		if aids, err = s.dao.PopRcmmndCache(ctx, key, cnt); err != nil {
			log.Error("s.dao.PopRcmmndCache(%d) error(%v)", key, err)
			return
		}
		if len(aids) < cnt {
			break
		}
		for _, aid := range aids {
			if _, ok := s.blackCache[aid]; ok {
				continue Retry
			}
		}
		var isOsea = model.IsOverseas(plat)
		if sis = s.fromAidsOsea(ctx, aids, isOsea); len(sis) < cnt {
			log.Warn("recommend aids(%v) get from archive have not normal(%v)", aids, sis)
			continue Retry
		}
		return
	}
	// if i==2, mean retry two counts, else if i<2, means break and recommend not enough.
	if i == 2 {
		return
	}
	if host := s.rcmmndHost(mid); host != "" {
		rcs, aids = s.apiRecommend(ctx, plat, key, host, mid)
	}
	var (
		clen  = len(rcs)
		caids = make([]int64, 0, cnt)
		fill  = cnt - clen
	)
	if clen+len(aids) < cnt {
		return
	}
	if cnt < clen {
		fill = 0
	}
	for _, rc := range rcs {
		if rc.Goto == "" || rc.Goto == model.GotoAv {
			caids = append(caids, rc.Aid)
			if len(caids) == cnt {
				break
			}
		}
	}
	if fill > 0 {
		caids = append(caids, aids[:fill]...)
	}
	if aids = aids[fill:]; len(aids) >= cnt {
		select {
		case s.rcmmndCh <- recommend{key: key, aids: aids[fill:]}:
		default:
			log.Warn("recommendProc chan full")
		}
	}
	var isOsea = model.IsOverseas(plat)                            //is overseas
	if sis = s.fromAidsOsea(ctx, caids, isOsea); len(sis) < clen { // NOTE: if cnt=1 means dislike change one
		for {
			var (
				over  = cnt - len(sis)
				start = 0
			)
			if over == 0 || start+over > len(aids) {
				break
			}
			if tmp := s.fromAidsOsea(ctx, aids[start:over], isOsea); len(tmp) > 0 {
				sis = append(sis, tmp...)
			}
		}
		return
	}
	for i, rc := range rcs {
		if rc.Goto != "" && rc.Goto != model.GotoAv {
			sis[i].Param = strconv.FormatInt(rc.Aid, 10)
			sis[i].Goto = rc.Goto
			sis[i].URI = model.FillURI(rc.Goto, sis[i].Param, nil)
		}
		if rc.Title != "" {
			sis[i].Title = rc.Title
		}
		if rc.Cover != "" {
			sis[i].Cover = rc.Cover
		}
	}
	sis = s.adVideo(ctx, mid, build, plat, buvid, network, mobiApp, device, ipaddr, sis)
	return
}

// rcmmndHost get recommend host
func (s *Service) rcmmndHost(mid int64) (host string) {
	// if mid=0, let host is 1: base recommend
	yu := mid % 20
	g := s.rcmmndGroup[yu]
	if hosts, ok := s.rcmmndHosts[g]; ok {
		if len(hosts) == 1 {
			host = hosts[0]
		} else {
			host = hosts[rand.Intn(len(hosts))]
		}
	}
	return
}

// apiRecommend get recommend fron big data.
func (s *Service) apiRecommend(ctx context.Context, plat int8, key, host string, mid int64) (rcs []*rcmmndCfg, aids []int64) {
	var (
		uri    string
		recURL = conf.Conf.Host.Data + "/mobile/home/%s"
	)
	uri = fmt.Sprintf(recURL, key)
	params := url.Values{}
	params.Set("plat", strconv.Itoa(int(plat)))
	params.Set("v2", "1")
	var res struct {
		Code    int          `json:"code"`
		Data    []int64      `json:"data"`
		Configs []*rcmmndCfg `json:"config"`
	}
	if err := s.client.Post(ctx, uri, "", params, &res); err != nil {
		log.Error("recommend url(%s) error(%v)", uri+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) res code(%d) or res.result(%v,%v)", uri, res.Code, res.Data, res.Configs)
		return
	}
	aids = res.Data
	rcs = res.Configs
	return
}

// itemNum get item number by plat.
func (s *Service) itemNum(plat int8) int {
	// cnt is items number
	cnt := 6
	if plat == model.PlatAndroid || plat == model.PlatAndroidI || plat == model.PlatAndroidG {
		cnt = 4
	} else if plat == model.PlatIPad || plat == model.PlatIPadI {
		cnt = 8
	} else if plat == model.PlatAndroidTV {
		cnt = 16
	}
	return cnt
}

// dealLive dela live data
func (s *Service) dealLive(c context.Context, sw *show.Show, sis []*show.Item) (rs *show.Show) {
	rs = &show.Show{
		Head: sw.Head,
		Body: sis,
		Ext:  sw.Ext,
	}
	return
}

// getLive get lives: feed, moe, hot.
func (s *Service) getLive(c context.Context, mid int64, ak, ip string, rand int, now time.Time) (sis []*show.Item) {
	const (
		_halfCnt = 2
	)
	sis = make([]*show.Item, _cnt) // _cnt=4 [0,1,2,3]: 0 1 feed and hot, 2 3 moe and hot
	feed, err := s.lv.Feed(c, mid, ak, ip, now)
	if err != nil {
		log.Error("s.live.Feed(%d) error(%d)", mid, err)
	}
	var have int
	// get two feed
	if feed != nil {
		for i := 0; i < _halfCnt && i < len(feed.Lives); i++ {
			si := &show.Item{}
			si.FromLive(feed.Lives[i])
			sis[i] = si
			have++
		}
	}
	// get two moe
	fdCnt := have
	start := _halfCnt * rand
	if len(s.liveMoeCache) < start+_halfCnt {
		start = 0
	}
	index := _halfCnt
MOENEXT:
	for _, l := range s.liveMoeCache[start:] {
		for i := 0; i < fdCnt; i++ {
			if sis[i].Param == l.Param {
				continue MOENEXT
			}
		}
		sis[index] = l
		index++
		have++
		if index >= _cnt {
			break
		}
	}
	// if feed and moe not enough, get hot
	yu := _cnt - have
	if yu > 0 {
		start := yu * rand
		if len(s.liveHotCache) < start+yu {
			start = 0
		}
		var nilI int
	HOTNEXT:
		for _, l := range s.liveHotCache[start:] {
			nilI = -1
			for i := len(sis) - 1; i >= 0; i-- {
				if sis[i] == nil {
					nilI = i
				} else if sis[i].Param == l.Param {
					continue HOTNEXT
				}
			}
			if nilI != -1 {
				sis[nilI] = l
				have++
			} else {
				return
			}
		}
	}
	if have < _cnt {
		for k, v := range sis {
			if v == nil {
				sis[k] = _emptyItem
			}
		}
	}
	return
}

// fromArchives return region show items from archive archives.
func (s *Service) fromArchivesPB(as []*api.Arc) (sis, sisOsea []*show.Item) {
	var asLen = len(as)
	if asLen == 0 {
		sis = []*show.Item{}
		return
	}
	sis = make([]*show.Item, 0, asLen)
	for _, a := range as {
		i := &show.Item{}
		i.FromArchivePB(a)
		if a.AttrVal(archive.AttrBitOverseaLock) == 0 {
			sisOsea = append(sisOsea, i)
		}
		sis = append(sis, i)
	}
	return
}

// fromArchivesBangumi aid to sid
func (s *Service) fromArchivesBangumi(c context.Context, as []*api.Arc, aids []int64, sids map[int32]*seasongrpc.CardInfoProto, bangumiType int) (sis, sisOsea []*show.Item) {
	var (
		asLen = len(as)
		err   error
		// bangumi
	)
	if asLen == 0 {
		sis = []*show.Item{}
		return
	}
	if sids == nil {
		if sids, err = s.fromSeasonID(c, aids); err != nil {
			log.Error("s.fromSeasonID error(%v)", err)
			return
		}
	}
	sis = make([]*show.Item, 0, asLen)
	for _, a := range as {
		i := &show.Item{}
		if sid, ok := sids[int32(a.Aid)]; ok && sid.SeasonId != 0 {
			i.FromArchivePBBangumi(a, sid, bangumiType)
		} else {
			i.FromArchivePB(a)
		}
		sis = append(sis, i)
		if a.AttrVal(archive.AttrBitOverseaLock) == 0 {
			sisOsea = append(sisOsea, i)
		}
	}
	return
}

// fromArchivesOsea  isOverseas
func (s *Service) fromArchivesOsea(as []*api.Arc, isOsea bool) (sis []*show.Item) {
	tmp, tmpOsea := s.fromArchivesPB(as)
	if isOsea {
		sis = tmpOsea
	} else {
		sis = tmp
	}
	return
}

// fromArchivesOsea  isOverseas
func (s *Service) fromArchivesBangumiOsea(c context.Context, as []*api.Arc, aids []int64, isOsea bool, bangumiType int) (sis []*show.Item) {
	tmp, tmpOsea := s.fromArchivesBangumi(c, as, aids, nil, bangumiType)
	if isOsea {
		sis = tmpOsea
	} else {
		sis = tmp
	}
	return
}

// fromAids get Aids.
func (s *Service) fromAids(ctx context.Context, aids []int64) (sis, sisOsea []*show.Item) {
	as, err := s.arc.ArchivesPB(ctx, aids)
	if err != nil {
		log.Error("s.arc.ArchivesPB aids(%v) error(%v)", aids, err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB(%v) length is 0", aids)
		return
	}
	sis = make([]*show.Item, 0, len(aids))
	for _, aid := range aids {
		var isOverseas int32
		si := &show.Item{}
		si.Goto = model.GotoAv
		si.Param = strconv.FormatInt(aid, 10)
		si.URI = model.FillURI(si.Goto, si.Param, nil)
		if v, ok := as[aid]; ok {
			isOverseas = v.AttrVal(archive.AttrBitOverseaLock)
			si.Danmaku = int(v.Stat.Danmaku)
			si.Play = int(v.Stat.View)
			si.Title = v.Title
			si.Duration = v.Duration
			si.Rname = v.TypeName
			si.Name = v.Author.Name
			si.Like = int(v.Stat.Like)
			si.Cover = model.CoverURL(v.Pic)
		}
		if isOverseas == 0 {
			sisOsea = append(sisOsea, si)
		}
		sis = append(sis, si)
	}
	return
}

// fromCardAids get Aids.
func (s *Service) fromCardAids(ctx context.Context, aids []int64) (sis map[int64]*show.Item) {
	as, err := s.arc.ArchivesPB(ctx, aids)
	if err != nil {
		log.Error("s.arc.ArchivesPB aids(%v) error(%v)", aids, err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB(%v) length is 0", aids)
		return
	}
	sis = map[int64]*show.Item{}
	for _, aid := range aids {
		si := &show.Item{}
		si.Goto = model.GotoAv
		si.Param = strconv.FormatInt(aid, 10)
		si.URI = model.FillURI(si.Goto, si.Param, nil)
		if v, ok := as[aid]; ok {
			if !v.IsNormal() {
				continue
			}
			si.Danmaku = int(v.Stat.Danmaku)
			si.Play = int(v.Stat.View)
			si.Title = v.Title
			si.Duration = v.Duration
			if region, ok := s.reRegionCache[int(v.TypeID)]; ok {
				si.Desc = region.Name
				si.Reid = region.Rid
			}
			si.Rid = int(v.TypeID)
			si.Rname = v.TypeName
			si.Name = v.Author.Name
			si.Like = int(v.Stat.Like)
			si.Cover = model.CoverURL(v.Pic)
		}
		sis[aid] = si
	}
	return
}

// fromRankAids
func (s *Service) fromRankAids(ctx context.Context, aids []int64, scores map[int64]int64, as map[int64]*api.Arc) (sis, sisOsea []*show.Item) {
	var (
		aid int64
		arc *api.Arc
		ok  bool
	)
	for _, aid = range aids {
		if arc, ok = as[aid]; ok {
			i := &show.Item{}
			if region, ok := s.reRegionCache[int(arc.TypeID)]; ok {
				i.Desc = region.Name
			}
			i.FromArchiveRank(arc, scores)
			if arc.AttrVal(archive.AttrBitOverseaLock) == 0 {
				sisOsea = append(sisOsea, i)
			}
			sis = append(sis, i)
		}
	}
	return
}

// fromAids get Aids.
func (s *Service) fromBgAids(ctx context.Context, aids []int64, sids map[int32]*seasongrpc.CardInfoProto, bangumiType int) (sis, sisOsea, sisbg, sisbgOsea, sisbgep, sisbgepOsea []*show.Item) {
	var (
		err error
	)
	as, err := s.arc.ArchivesPB(ctx, aids)
	if err != nil {
		log.Error("s.arc.ArchivesPB aids(%v) error(%v)", aids, err)
		return
	}
	if len(as) == 0 {
		log.Warn("s.arc.ArchivesPB(%v) length is 0", aids)
		return
	}
	sis = make([]*show.Item, 0, len(aids))
	if sids == nil {
		if sids, err = s.fromSeasonID(ctx, aids); err != nil {
			log.Error("s.fromSeasonID error(%v)", err)
			return
		}
	}
	for _, aid := range aids {
		var isOverseas int32
		si := &show.Item{}
		sibg := &show.Item{}
		sibgep := &show.Item{}
		if v, ok := as[aid]; ok {
			isOverseas = v.AttrVal(archive.AttrBitOverseaLock)
			if sid, ok := sids[int32(aid)]; ok && sid.SeasonId != 0 {
				sibg.FromArchivePBBangumi(v, sid, _bangumiSeasonID)
				sibgep.FromArchivePBBangumi(v, sid, _bangumiEpisodeID)
			} else {
				sibg.FromArchivePB(v)
				sibgep.FromArchivePB(v)
			}
			si.FromArchivePB(v)
			if isOverseas == 0 {
				sisOsea = append(sisOsea, si)
				sisbgOsea = append(sisbgOsea, sibg)
				sisbgepOsea = append(sisbgepOsea, sibg)
			}
			sis = append(sis, si)
			sisbg = append(sisbg, sibg)
			sisbgep = append(sisbgep, sibgep)
		}
	}
	return
}

// fromSeasonID
func (s *Service) fromSeasonID(c context.Context, arcAids []int64) (seasonID map[int32]*seasongrpc.CardInfoProto, err error) {
	if seasonID, err = s.bgm.CardsByAids(c, arcAids); err != nil {
		log.Error("s.bgm.Seasonid CardsByAids %v", err)
	}
	return
}

// isOverseas
func (s *Service) fromAidsOsea(ctx context.Context, aids []int64, isOsea bool) (sis []*show.Item) {
	tmp, tmpOsea := s.fromAids(ctx, aids)
	if isOsea {
		sis = tmpOsea
	} else {
		sis = tmp
	}
	return
}

// adVideo
func (s *Service) adVideo(ctx context.Context, mid int64, build int, plat int8, buvid, network, mobiApp, device, ipaddr string, sis []*show.Item) (res []*show.Item) {
	var cpmsis map[int]*show.Item
	if resID, ok := _recommend[plat]; ok {
		cpmsis = s.cpmRecommend(ctx, mid, build, buvid, resID, network, mobiApp, device, ipaddr)
	}
	for rank, ad := range cpmsis {
		if len(sis) >= rank {
			if ad.IsAdReplace {
				sis[rank-1] = ad
			} else {
				sis[rank-1].IsAdLoc = true
				sis[rank-1].IsAd = ad.IsAd
				sis[rank-1].CmMark = ad.CmMark
				sis[rank-1].SrcId = ad.SrcId
				sis[rank-1].RequestId = ad.RequestId
				sis[rank-1].ClientIp = ad.ClientIp
			}
		}
	}
	res = sis
	return
}
