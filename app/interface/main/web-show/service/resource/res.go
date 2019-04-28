package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/web-show/dao/resource"
	rsmdl "go-common/app/interface/main/web-show/model/resource"
	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	locmdl "go-common/app/service/main/location/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_nullImage     = "https://static.hdslb.com/images/transparent.gif"
	_videoPrefix   = "http://www.bilibili.com/video/av"
	_bangumiPrefix = "bilibili://bangumi/season/"
	_GamePrefix    = "bilibili://game/"
	_LivePrefix    = "bilibili://live/"
	_AVprefix      = "bilibili://video/"
	_topicPrefix   = "//www.bilibili.com/tag/"
)

var (
	_emptyRelation = []*rsmdl.Relation{}
	_emptyAsgs     = []*rsmdl.Assignment{}
	_contractMap   = map[string]struct{}{
		"banner":     struct{}{},
		"focus":      struct{}{},
		"promote":    struct{}{},
		"app_banner": struct{}{},
		"text_link":  struct{}{},
		"frontpage":  struct{}{},
	}
	_bannerID = map[int64]struct{}{
		142:  struct{}{},
		1576: struct{}{},
		1580: struct{}{},
		1584: struct{}{},
		1588: struct{}{},
		1592: struct{}{},
		1596: struct{}{},
		1600: struct{}{},
		1604: struct{}{},
		1608: struct{}{},
		1612: struct{}{},
		1616: struct{}{},
		1620: struct{}{},
		1622: struct{}{},
		1634: struct{}{},
		1920: struct{}{},
		2260: struct{}{},
		2210: struct{}{},
	}
	_cpmGrayRate   = int64(0)
	_white         = map[int64]struct{}{}
	_cpmOn         = true
	_RelationResID = 162
)

// URLMonitor return all urls configured
func (s *Service) URLMonitor(c context.Context, pf int) (urls map[string]string) {
	return s.urlMonitor[pf]
}

// GrayRate return gray  percent
func (s *Service) GrayRate(c context.Context) (r int64, ws []int64, swt bool) {
	r = _cpmGrayRate
	for w := range _white {
		ws = append(ws, w)
	}
	swt = _cpmOn
	return
}

// SetGrayRate set gray percent
func (s *Service) SetGrayRate(c context.Context, swt bool, rate int64, white []int64) {
	_cpmGrayRate = rate
	tmp := map[int64]struct{}{}
	for _, w := range white {
		tmp[w] = struct{}{}
	}
	_cpmOn = swt
	_white = tmp
}

// Resources get resource info by pf,ids
func (s *Service) Resources(c context.Context, arg *rsmdl.ArgRess) (mres map[string][]*rsmdl.Assignment, count int, err error) {
	var (
		aids                    []int64
		arcs                    map[int64]*api.Arc
		country, province, city string
		info                    *locmdl.Info
	)
	arg.IP = metadata.String(c, metadata.RemoteIP)
	if info, err = s.locRPC.Info(c, &locmdl.ArgIP{IP: arg.IP}); err != nil {
		log.Error("Location RPC error %v", err)
		err = nil
	}
	if info != nil {
		country = info.Country
		province = info.Province
		city = info.City
	}
	area := checkAera(country)
	var cpmInfos map[int64]*rsmdl.Assignment
	if _cpmOn {
		cpmInfos = s.cpms(c, arg.Mid, arg.Ids, arg.Sid, arg.IP, country, province, city, arg.Buvid)
	} else if _, ok := _white[arg.Mid]; ok || (arg.Mid%100 < _cpmGrayRate && arg.Mid != 0) {
		cpmInfos = s.cpms(c, arg.Mid, arg.Ids, arg.Sid, arg.IP, country, province, city, arg.Buvid)
	}
	mres = make(map[string][]*rsmdl.Assignment)
	for _, id := range arg.Ids {
		pts := s.posCache[posKey(arg.Pf, int(id))]
		if pts == nil {
			continue
		}
		count = pts.Counter
		// add ads if exists
		res, as := s.res(c, cpmInfos, int(id), area, pts, arg.Mid)
		mres[strconv.FormatInt(id, 10)] = res
		aids = append(aids, as...)
	}
	// fill archive if has video ad
	if len(aids) != 0 {
		argAids := &arcmdl.ArgAids2{
			Aids: aids,
		}
		if arcs, err = s.arcRPC.Archives3(c, argAids); err != nil {
			resource.PromError("arcRPC.Archives3", "s.arcRPC.Archives3(arcAids:(%v), arcs), err(%v)", aids, err)
			return
		}
		for _, tres := range mres {
			for _, rs := range tres {
				if arc, ok := arcs[rs.Aid]; ok {
					rs.Archive = arc
					if rs.Name == "" {
						rs.Name = arc.Title
					}
					if rs.Pic == "" {
						rs.Pic = arc.Pic
					}
				}
			}
		}
	}
	// if id is banner and not content add defult
	for i, rs := range mres {
		if len(rs) == 0 {
			id, _ := strconv.ParseInt(i, 10, 64)
			if _, ok := _bannerID[id]; ok {
				mres[i] = append(mres[i], s.defBannerCache)
			}
		}
	}
	return
}

// cpmBanners
func (s *Service) cpms(c context.Context, mid int64, ids []int64, sid, ip, country, province, city, buvid string) (res map[int64]*rsmdl.Assignment) {
	cpmInfos, err := s.adDao.Cpms(c, mid, ids, sid, ip, country, province, city, buvid)
	if err != nil {
		log.Error("s.adDao.Cpms error(%v)", err)
		return
	}
	res = make(map[int64]*rsmdl.Assignment, len(cpmInfos.AdsInfo))
	for _, id := range ids {
		idStr := strconv.FormatInt(id, 10)
		if adsInfos := cpmInfos.AdsInfo[idStr]; len(adsInfos) > 0 {
			for srcStr, adsInfo := range adsInfos {
				// var url string
				srcIDInt, _ := strconv.ParseInt(srcStr, 10, 64)
				if adInfo := adsInfo.AdInfo; adInfo != nil {
					//switch adInfo.CreativeType {
					// case 0:
					// 	url = adInfo.CreativeContent.URL
					// case 1:
					// 	url = "www.bilibili.com/video/av" + adInfo.CreativeContent.VideoID
					// }
					ad := &rsmdl.Assignment{
						CreativeType: adInfo.CreativeType,
						Aid:          adInfo.CreativeContent.VideoID,
						RequestID:    cpmInfos.RequestID,
						SrcID:        srcIDInt,
						IsAdLoc:      true,
						IsAd:         adsInfo.IsAd,
						CmMark:       adsInfo.CmMark,
						CreativeID:   adInfo.CreativeID,
						AdCb:         adInfo.AdCb,
						ShowURL:      adInfo.CreativeContent.ShowURL,
						ClickURL:     adInfo.CreativeContent.ClickURL,
						Name:         adInfo.CreativeContent.Title,
						Pic:          adInfo.CreativeContent.ImageURL,
						LitPic:       adInfo.CreativeContent.ThumbnailURL,
						URL:          adInfo.CreativeContent.URL,
						PosNum:       int(adsInfo.Index),
						Title:        adInfo.CreativeContent.Title,
						ServerType:   rsmdl.FromCpm,
						IsCpm:        true,
					}
					res[srcIDInt] = ad
				} else {
					ad := &rsmdl.Assignment{
						IsAdLoc:   true,
						RequestID: cpmInfos.RequestID,
						IsAd:      false,
						SrcID:     srcIDInt,
						ResID:     int(id),
						CmMark:    adsInfo.CmMark,
					}
					res[srcIDInt] = ad
				}
			}
		}
	}
	return
}

func checkAera(country string) (area int8) {
	switch country {
	case "中国":
		area = 1
	case "香港", "台湾", "澳门":
		area = 2
	case "日本":
		area = 3
	case "美国":
		area = 4
	default:
		if _, ok := rsmdl.OverSeasCountry[country]; ok {
			area = 5
		} else {
			area = 0
		}
	}
	return
}

// Relation get relation archives by aid
func (s *Service) Relation(c context.Context, arg *rsmdl.ArgAid) (rls []*rsmdl.Relation, err error) {

	var aids []int64
	rls = _emptyRelation
	arg.IP = metadata.String(c, metadata.RemoteIP)
	if aids, err = s.dataDao.Related(c, arg.Aid, arg.IP); err != nil {
		log.Error("s.dataDao.Related aid(%v) error(%v)", arg.Aid, err)
		return
	}

	if len(aids) == 0 {
		log.Warn("zero_relates")
		return
	}
	argAdis := &arcmdl.ArgAids2{
		Aids:   aids,
		RealIP: arg.IP,
	}
	arcs, err := s.arcRPC.Archives3(c, argAdis)
	if err != nil {
		log.Info("s.arcDao.Archives3", err)
		return
	}
	var res []*rsmdl.Relation
	for _, arc := range arcs {
		res = append(res, &rsmdl.Relation{Arc: arc})
	}
	rls = res
	if len(res) < 3 {
		return
	}
	var (
		country, province, city string
		info                    *locmdl.Info
	)
	if info, err = s.locRPC.Info(c, &locmdl.ArgIP{IP: arg.IP}); err != nil {
		log.Error("Location RPC error %v", err)
		err = nil
	}
	if info != nil {
		country = info.Country
		province = info.Province
		city = info.City
	}
	area := checkAera(country)
	//pts := s.posCache[posKey(0, _RelationResID)]
	cpmInfos := s.cpms(c, arg.Mid, []int64{int64(_RelationResID)}, arg.Sid, arg.IP, country, province, city, arg.Buvid)
	for _, rs := range cpmInfos {
		// just fet one ad
		if rs.IsAd {
			arcAid := &arcmdl.ArgAid2{
				Aid: rs.Aid,
			}
			var arc *api.Arc
			arc, err = s.arcRPC.Archive3(c, arcAid)
			if err != nil {
				resource.PromError("arcRPC.Archive3", "s.arcRPC.Archive3(arcAid:(%v), arcs), err(%v)", rs.Aid, err)
				err = nil
				rls = res
				return
			}
			rl := &rsmdl.Relation{
				Arc:        arc,
				Area:       area,
				RequestID:  rs.RequestID,
				CreativeID: rs.CreativeID,
				AdCb:       rs.AdCb,
				SrcID:      rs.SrcID,
				ShowURL:    rs.ShowURL,
				ClickURL:   rs.ClickURL,
				IsAdLoc:    rs.IsAdLoc,
				ResID:      _RelationResID,
				IsAd:       true,
			}
			if rs.Pic != "" {
				rl.Pic = rs.Pic
			}
			if rs.Title != "" {
				rl.Title = rs.Title
			}
			rls = append(res[:2], append([]*rsmdl.Relation{rl}, res[2:]...)...)
			return
		}
		res[2].AdCb = rs.AdCb
		res[2].SrcID = rs.SrcID
		res[2].ShowURL = rs.ShowURL
		res[2].ClickURL = rs.ClickURL
		res[2].IsAdLoc = rs.IsAdLoc
		res[2].RequestID = rs.RequestID
		res[2].CreativeID = rs.CreativeID
		res[2].ResID = _RelationResID
		return
	}
	return
}

// Resource get resource info by pf,id
func (s *Service) Resource(c context.Context, arg *rsmdl.ArgRes) (res []*rsmdl.Assignment, count int, err error) {
	var (
		aids                    []int64
		arcs                    map[int64]*api.Arc
		country, province, city string
		info                    *locmdl.Info
	)
	arg.IP = metadata.String(c, metadata.RemoteIP)
	res = _emptyAsgs
	pts := s.posCache[posKey(arg.Pf, int(arg.ID))]
	if pts == nil {
		return
	}
	count = pts.Counter
	if info, err = s.locRPC.Info(c, &locmdl.ArgIP{IP: arg.IP}); err != nil {
		log.Error("Location RPC error %v", err)
		err = nil
	}
	if info != nil {
		country = info.Country
		province = info.Province
		city = info.City
	}
	area := checkAera(country)
	var cpmInfos map[int64]*rsmdl.Assignment
	if _cpmOn {
		cpmInfos = s.cpms(c, arg.Mid, []int64{arg.ID}, arg.Sid, arg.IP, country, province, city, arg.Buvid)
	} else if _, ok := _white[arg.Mid]; ok || (arg.Mid%100 < _cpmGrayRate && arg.Mid != 0) {
		cpmInfos = s.cpms(c, arg.Mid, []int64{arg.ID}, arg.Sid, arg.IP, country, province, city, arg.Buvid)
	}
	res, aids = s.res(c, cpmInfos, int(arg.ID), area, pts, arg.Mid)
	// fill archive if has video ad
	if len(aids) != 0 {
		argAids := &arcmdl.ArgAids2{
			Aids: aids,
		}
		if arcs, err = s.arcRPC.Archives3(c, argAids); err != nil {
			resource.PromError("arcRPC.Archive3", "s.arcRPC.Archive3(arcAid:(%v), arcs), err(%v)", aids, err)
			return
		}
		for _, rs := range res {
			if arc, ok := arcs[rs.Aid]; ok {
				rs.Archive = arc
				if rs.Name == "" {
					rs.Name = arc.Title
				}
				if rs.Pic == "" {
					rs.Pic = arc.Pic
				}
			}
		}
	}
	// add defBanner if contnent not exits
	if len(res) == 0 {
		if _, ok := _bannerID[int64(arg.ID)]; ok {
			// 142 is index banner
			if rs := s.resByID(142); rs != nil {
				res = append(res, rs)
			} else {
				res = append(res, s.defBannerCache)
			}
		}
	}
	return
}

func (s *Service) res(c context.Context, cpmInfos map[int64]*rsmdl.Assignment, id int, area int8, pts *rsmdl.Position, mid int64) (res []*rsmdl.Assignment, aids []int64) {
	// add ads if exists
	var (
		reqID    string
		index    int
		resIndex int
		ts       = strconv.FormatInt(time.Now().Unix(), 10)
		resBs    []*rsmdl.Assignment
	)
	for _, pt := range pts.Pos {
		var isAdLoc bool
		if rs, ok := cpmInfos[int64(pt.ID)]; ok && rs.IsCpm {
			rs.Area = area
			if mid != 0 {
				rs.Mid = strconv.FormatInt(mid, 10)
			}
			if rs.CreativeType == rsmdl.CreativeVideo {
				aids = append(aids, rs.Aid)
			}
			rs.ServerType = rsmdl.FromCpm
			res = append(res, rs)
			if rsb := s.resByID(pt.ID); rsb != nil {
				// url mean aid in Asgtypevideo
				resBs = append(resBs, rsb)
			}
			continue
		} else if ok && !rs.IsCpm {
			isAdLoc = true
			reqID = rs.RequestID
		}
		var rs *rsmdl.Assignment
		resBslen := len(resBs)
		if resIndex < resBslen {
			rs = resBs[resIndex]
			resIndex++
			if rsb := s.resByID(pt.ID); rsb != nil {
				resBs = append(resBs, rsb)
			}
		} else if rs = s.resByID(pt.ID); rs != nil {

		} else {
			rs = s.resByIndex(id, index)
			index++
		}
		if rs != nil {
			if rs.Atype == rsmdl.AsgTypeVideo || rs.Atype == rsmdl.AsgTypeAv {
				aids = append(aids, rs.Aid)
			}
			rs.PosNum = pt.PosNum
			rs.SrcID = int64(pt.ID)
			rs.Area = area
			rs.IsAdLoc = isAdLoc
			if isAdLoc {
				rs.RequestID = reqID
			} else {
				rs.RequestID = ts
			}
			if mid != 0 {
				rs.Mid = strconv.FormatInt(mid, 10)
			}
			res = append(res, rs)
		} else if isAdLoc {
			rs = &rsmdl.Assignment{
				PosNum:    pt.PosNum,
				SrcID:     int64(pt.ID),
				IsAdLoc:   isAdLoc,
				RequestID: reqID,
				Area:      area,
				Pic:       _nullImage,
			}
			if mid != 0 {
				rs.Mid = strconv.FormatInt(mid, 10)
			}
			res = append(res, rs)
		}
	}
	return
}

// resByIndex return res of index
func (s *Service) resByIndex(id, index int) (res *rsmdl.Assignment) {
	ss := s.asgCache[id]
	if index >= len(ss) {
		return
	}
	res = new(rsmdl.Assignment)
	*res = *(ss[index])
	return
}

// resByID return res of id
func (s *Service) resByID(id int) (res *rsmdl.Assignment) {
	ss := s.asgCache[id]
	l := len(ss)
	if l == 0 {
		return
	}
	res = ss[0]
	for _, s := range ss {
		// ContractId not in contractMap ,it is ad and ad first
		if _, ok := _contractMap[s.ContractID]; !ok {
			res = s
			return
		}
	}
	return
}

// rpc resourcesALL
func (s *Service) resourcesALL() (rscs []*rsmdl.Res, err error) {
	resourcesRPC, err := s.recrpc.ResourceAll(context.Background())
	if err != nil {
		resource.PromError("recRPC.ResourcesALL", "s.recrpc.resourcesRPC error(%v)", err)
		return
	}
	rscs = make([]*rsmdl.Res, 0)
	for _, res := range resourcesRPC {
		rsc := &rsmdl.Res{
			ID:       res.ID,
			Platform: res.Platform,
			Name:     res.Name,
			Parent:   res.Parent,
			Counter:  res.Counter,
			Position: res.Position,
		}
		rscs = append(rscs, rsc)
	}
	return
}

// rpc assignmentAll
func (s *Service) assignmentAll() (asgs []*rsmdl.Assignment, err error) {
	assignRPC, err := s.recrpc.AssignmentAll(context.Background())
	if err != nil {
		resource.PromError("recRPC.AssignmentAll", "s.recrpc.assignRPC error(%v)", err)
		return
	}
	asgs = make([]*rsmdl.Assignment, 0)
	for _, asgr := range assignRPC {
		asg := &rsmdl.Assignment{
			ID:         asgr.ID,
			Name:       asgr.Name,
			ContractID: asgr.ContractID,
			ResID:      asgr.ResID,
			Pic:        asgr.Pic,
			LitPic:     asgr.LitPic,
			URL:        asgr.URL,
			Atype:      asgr.Atype,
			Weight:     asgr.Weight,
			Rule:       asgr.Rule,
			Agency:     asgr.Agency,
			STime:      asgr.STime,
		}
		asgs = append(asgs, asg)
	}
	return
}

// default banner
func (s *Service) defBanner() (asg *rsmdl.Assignment, err error) {
	bannerRPC, err := s.recrpc.DefBanner(context.Background())
	if err != nil {
		resource.PromError("recRPC.defBanner", "s.recrpc.defBanner error(%v)", err)
		return
	}
	if bannerRPC != nil {
		asg = &rsmdl.Assignment{
			ID:         bannerRPC.ID,
			Name:       bannerRPC.Name,
			ContractID: bannerRPC.ContractID,
			ResID:      bannerRPC.ResID,
			Pic:        bannerRPC.Pic,
			LitPic:     bannerRPC.LitPic,
			URL:        bannerRPC.URL,
			Atype:      bannerRPC.Atype,
			Weight:     bannerRPC.Weight,
			Rule:       bannerRPC.Rule,
			Agency:     bannerRPC.Agency,
		}
	}
	return
}

// LoadRes load Res info to cache
func (s *Service) loadRes() (err error) {
	assign, err := s.assignmentAll()
	if err != nil {
		return
	}
	resources, err := s.resourcesALL()
	if err != nil {
		return
	}
	resMap := make(map[int]*rsmdl.Res)
	posMap := make(map[string]*rsmdl.Position)
	for _, res := range resources {
		resMap[res.ID] = res
		if res.Counter > 0 {
			key := posKey(res.Platform, res.ID)
			pos := &rsmdl.Position{
				Pos:     make([]*rsmdl.Loc, 0),
				Counter: res.Counter,
			}
			posMap[key] = pos
		} else {
			key := posKey(res.Platform, res.Parent)
			if pos, ok := posMap[key]; ok {
				loc := &rsmdl.Loc{
					ID:     res.ID,
					PosNum: res.Position,
				}
				pos.Pos = append(pos.Pos, loc)
			}
		}
	}
	for _, a := range assign {
		if res, ok := resMap[a.ResID]; ok {
			if err = s.convertURL(a); err != nil {
				return
			}
			var data struct {
				Cover        int32  `json:"is_cover"`
				Style        int32  `json:"style"`
				Label        string `json:"label"`
				Intro        string `json:"intro"`
				CreativeType int8   `json:"creative_type"`
			}
			//  unmarshal rule for frontpage style
			if a.Rule != "" {
				e := json.Unmarshal([]byte(a.Rule), &data)
				if e != nil {
					log.Error("json.Unmarshal (%s) error(%v)", a.Rule, e)
				} else {
					a.Style = data.Style
					a.CreativeType = data.CreativeType
					if a.ContractID == "rec_video" {
						a.Label = data.Label
						a.Intro = data.Intro
					}
				}
			}
			res.Assignments = append(res.Assignments, a)
		}

	}
	urlMonitor := make(map[int]map[string]string)
	tmp := make(map[int][]*rsmdl.Assignment, len(resMap))
	for _, res := range resMap {
		tmp[res.ID] = res.Assignments
		urlMap, ok := urlMonitor[res.Platform]
		if !ok {
			urlMap = make(map[string]string)
			urlMonitor[res.Platform] = urlMap
		}
		for _, a := range res.Assignments {
			urlMap[fmt.Sprintf("%d_%s", a.ResID, a.Name)] = a.URL
		}
	}
	s.asgCache = tmp
	s.posCache = posMap
	s.urlMonitor = urlMonitor
	// load default banner
	banner, err := s.defBanner()
	if err != nil {
		return
	} else if banner != nil {
		var data struct {
			Cover int32 `json:"is_cover"`
			Style int32 `json:"style"`
		}
		err := json.Unmarshal([]byte(banner.Rule), &data)
		if err != nil {
			log.Error("json.Unmarshal (%s) error(%v)", banner.Rule, err)
		} else {
			banner.Style = data.Style
		}
		s.defBannerCache = banner
	}
	return
}

func posKey(pf, id int) string {
	return fmt.Sprintf("%d_%d", pf, id)
}

func (s *Service) convertURL(a *rsmdl.Assignment) (err error) {
	switch a.Atype {
	case rsmdl.AsgTypeVideo:
		var aid int64
		if aid, err = strconv.ParseInt(a.URL, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) err(%v)", a.URL, err)
			return
		}
		a.Aid = aid
		a.URL = _videoPrefix + a.URL
	case rsmdl.AsgTypeURL:
	case rsmdl.AsgTypeBangumi:
		a.URL = _bangumiPrefix + a.URL
	case rsmdl.AsgTypeLive:
		a.URL = _LivePrefix + a.URL
	case rsmdl.AsgTypeGame:
		a.URL = _GamePrefix + a.URL
	case rsmdl.AsgTypeAv:
		var aid int64
		if aid, err = strconv.ParseInt(a.URL, 10, 64); err != nil {
			log.Error("strconv.ParseInt(%s) err(%v)", a.URL, err)
			return
		}
		a.Aid = aid
		a.URL = _AVprefix + a.URL
	case rsmdl.AsgTypeTopic:
		a.URL = _topicPrefix + a.URL
	}
	return
}
