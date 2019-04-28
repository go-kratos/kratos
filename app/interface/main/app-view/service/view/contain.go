package view

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"sort"
	"time"

	"go-common/app/interface/main/app-view/model"
	"go-common/app/interface/main/app-view/model/ad"
	"go-common/app/interface/main/app-view/model/bangumi"
	"go-common/app/interface/main/app-view/model/game"
	"go-common/app/interface/main/app-view/model/manager"
	"go-common/app/interface/main/app-view/model/tag"
	"go-common/app/interface/main/app-view/model/view"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	location "go-common/app/service/main/location/model"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	_rate = map[int]int64{15: 464, 16: 464, 32: 1028, 48: 1328, 64: 2192, 74: 3192, 80: 3192, 112: 6192, 116: 6192, 66: 1820}
)

const (
	_dmformat         = "http://comment.bilibili.com/%d.xml"
	_qn480            = 32
	_qnAndroidBuildGt = 5325000
	_qnIosBuildGt     = 8170
	_qnAndroidBuildLt = 5335000
	_qnIosBuildLt     = 8190
)

// initReqUser init Req User
func (s *Service) initReqUser(c context.Context, v *view.View, mid int64, plat int8, build int) {
	const (
		_androidOld = 427000
		_iosOld     = 4000
		_ipadOld    = 4300
		_ipadHD     = 10410
	)
	// owner ext
	var (
		owners []int64
		cards  map[int64]*account.Card
		fls    map[int64]int8
	)
	g, ctx := errgroup.WithContext(c)
	if v.Author.Mid > 0 {
		owners = append(owners, v.Author.Mid)
		for _, staffInfo := range v.StaffInfo {
			owners = append(owners, staffInfo.Mid)
		}
		g.Go(func() (err error) {
			v.OwnerExt.OfficialVerify.Type = -1
			cards, err = s.accDao.Cards3(ctx, owners)
			if err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if card, ok := cards[v.Author.Mid]; ok && card != nil {
				otp := -1
				odesc := ""
				if card.Official.Role != 0 {
					if card.Official.Role <= 2 {
						otp = 0
					} else {
						otp = 1
					}
					odesc = card.Official.Title
				}
				v.OwnerExt.OfficialVerify.Type = otp
				v.OwnerExt.OfficialVerify.Desc = odesc
				v.OwnerExt.Vip.Type = int(card.Vip.Type)
				v.OwnerExt.Vip.VipStatus = int(card.Vip.Status)
				v.OwnerExt.Vip.DueDate = card.Vip.DueDate
				v.Author.Name = card.Name
				v.Author.Face = card.Face
			}
			if l, ok := s.liveCache[v.Author.Mid]; ok {
				v.OwnerExt.Live = l
			}
			return
		})
		g.Go(func() (err error) {
			stat, err := s.relDao.Stat(c, v.Author.Mid)
			if err != nil {
				log.Error("%+v", err)
				err = nil
				return
			}
			if stat != nil {
				v.OwnerExt.Fans = int(stat.Follower)
			}
			return
		})
		g.Go(func() error {
			if ass, err := s.assDao.Assist(ctx, v.Author.Mid); err != nil {
				log.Error("%+v", err)
			} else {
				v.OwnerExt.Assists = ass
			}
			return nil
		})
	}
	// req user
	v.ReqUser = &view.ReqUser{Favorite: 0, Attention: -999, Like: 0, Dislike: 0}
	// check req user
	if mid > 0 {
		g.Go(func() error {
			var is bool
			if (model.IsAndroid(plat) && build < _androidOld) || (model.IsIPhone(plat) && build < _iosOld) || ((plat == model.PlatIPad && build < _ipadOld) || (plat == model.PlatIpadHD && build < _ipadHD) || (plat == model.PlatIPadI)) {
				is, _ = s.favDao.IsFavDefault(ctx, mid, v.Aid)
			} else {
				is, _ = s.favDao.IsFav(ctx, mid, v.Aid)
			}
			if is {
				v.ReqUser.Favorite = 1
			}
			return nil
		})
		g.Go(func() error {
			res, err := s.thumbupDao.HasLike(ctx, mid, _businessLike, []int64{v.Aid})
			if err != nil {
				log.Error("s.thumbupDao.HasLike err(%+v)", err)
				return nil
			}
			if res.States == nil {
				return nil
			}
			if typ, ok := res.States[v.Aid]; ok {
				if typ.State == thumbup.StateLike {
					v.ReqUser.Like = 1
				} else if typ.State == thumbup.StateDislike {
					v.ReqUser.Dislike = 1
				}
			}
			return nil
		})
		g.Go(func() (err error) {
			res, err := s.coinDao.ArchiveUserCoins(ctx, v.Aid, mid, _avTypeAv)
			if err != nil {
				log.Error("%+v", err)
				err = nil
			}
			if res != nil && res.Multiply > 0 {
				v.ReqUser.Coin = 1
			}
			return
		})
		if v.Author.Mid > 0 {
			g.Go(func() error {
				fls = s.accDao.IsAttention(ctx, owners, mid)
				if _, ok := fls[v.Author.Mid]; ok {
					v.ReqUser.Attention = 1
				}
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		log.Error("%+v", err)
	}
	// fill staff
	for _, owner := range owners {
		if card, ok := cards[owner]; ok && card != nil {
			staff := &view.Staff{Mid: owner}
			if owner == v.Author.Mid {
				staff.Title = "UP主"
			} else {
				for _, s := range v.StaffInfo {
					if s.Mid == owner {
						staff.Title = s.Title
					}
				}
			}
			staff.Name = card.Name
			staff.Face = card.Face
			staff.OfficialVerify.Type = -1
			otp := -1
			odesc := ""
			if card.Official.Role != 0 {
				if card.Official.Role <= 2 {
					otp = 0
				} else {
					otp = 1
				}
				odesc = card.Official.Title
			}
			staff.OfficialVerify.Type = otp
			staff.OfficialVerify.Desc = odesc
			staff.Vip.Type = int(card.Vip.Type)
			staff.Vip.VipStatus = int(card.Vip.Status)
			staff.Vip.DueDate = card.Vip.DueDate
			if _, ok := fls[owner]; ok {
				staff.Attention = 1
			}
			v.Staff = append(v.Staff, staff)
		}
	}
}

func (s *Service) initRelateCMTag(c context.Context, v *view.View, plat int8, build, qn, fnver, fnval, forceHost, parentMode int, mid int64, buvid, mobiApp, device, network, adExtra, from string, now time.Time) {
	const (
		_iPhoneRelateGame  = 6500
		_androidRelateGame = 5210000
	)
	var (
		rls              []*view.Relate
		aidm             map[int64]struct{}
		mr               *manager.Relate
		rGameID, cGameID int64
		advert           *ad.Ad
		adm              map[int]*ad.AdInfo
		relateRsc        int64
		cmRsc            int64
		err              error
		adminfo          json.RawMessage
		hasDalao         int
		dalaoExp         int
	)
	tids := s.initTag(c, v, mid, plat)
	g, ctx := errgroup.WithContext(c)
	g.Go(func() (err error) {
		if mid > 0 || buvid != "" {
			if rls, v.UserFeature, v.ReturnCode, hasDalao, dalaoExp, err = s.newRcmdRelate(ctx, plat, v.Aid, mid, buvid, mobiApp, from, build, qn, fnver, fnval, forceHost, parentMode); err != nil {
				log.Error("s.newRcmdRelate(%d) error(%+v)", v.Aid, err)
			}
		}
		if len(rls) == 0 {
			rls, _, err = s.dealRcmdRelate(ctx, plat, v.Aid, build)
			log.Warn("s.dealRcmdRelate aid(%d) mid(%d) buvid(%s)", v.Aid, mid, buvid)
		} else {
			v.IsRec = 1
			log.Warn("s.newRcmdRelate returncode(%s) aid(%d) mid(%d) buvid(%s) hasDalao(%d) dalaoExp(%d)", v.ReturnCode, v.Aid, mid, buvid, hasDalao, dalaoExp)
		}
		err = nil
		return
	})
	if !model.IsIPad(plat) {
		g.Go(func() (err error) {
			const (
				_iphoneRelateRsc  = 2029
				_androidRelateRsc = 2028
				_iphoneCMRsc      = 2335
				_androidCMRsc     = 2337
			)
			if model.IsIPhone(plat) {
				relateRsc = _iphoneRelateRsc
				cmRsc = _iphoneCMRsc
			} else {
				relateRsc = _androidRelateRsc
				cmRsc = _androidCMRsc
			}
			if advert, err = s.adDao.Ad(ctx, mobiApp, device, buvid, build, mid, v.Author.Mid, v.Aid, v.TypeID, tids, []int64{relateRsc, cmRsc}, network, adExtra); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
		if v.OrderID > 0 {
			g.Go(func() (err error) {
				if adminfo, err = s.adDao.MonitorInfo(ctx, v.Aid); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
	}
	if dalaoExp != 1 {
		g.Go(func() (err error) {
			mr = s.relateCache(ctx, plat, build, now, v.Aid, tids, v.TypeID)
			return
		})
	}
	if (plat == model.PlatAndroid && build >= _androidRelateGame) || (plat == model.PlatIPhone && build >= _iPhoneRelateGame) || plat == model.PlatIPhoneB {
		if buvid != "" && crc32.ChecksumIEEE([]byte(buvid))%10 == 1 {
			g.Go(func() (err error) {
				rGameID = s.relateGame(ctx, v.Aid)
				return
			})
		}
		if v.AttrVal(archive.AttrBitIsPorder) == archive.AttrYes || v.OrderID > 0 {
			g.Go(func() (err error) {
				if cGameID, err = s.arcDao.Commercial(ctx, v.Aid); err != nil {
					log.Error("%+v", err)
					err = nil
				}
				return
			})
		}
	}
	if err = g.Wait(); err != nil {
		log.Error("%+v", err)
		return
	}
	//ad config
	if advert != nil {
		if advert.AdsControl != nil {
			v.CMConfig = &view.CMConfig{
				AdsControl: advert.AdsControl,
			}
		}
	}
	if adminfo != nil {
		if v.CMConfig == nil {
			v.CMConfig = &view.CMConfig{
				MonitorInfo: adminfo,
			}
		} else {
			v.CMConfig.MonitorInfo = adminfo
		}
	}
	//ad
	if len(rls) == 0 {
		s.prom.Incr("zero_relates")
		return
	}
	var (
		r  *view.Relate
		rm map[int]*view.Relate
	)
	if advert != nil {
		if adm, err = s.dealCM(c, advert, relateRsc); err != nil {
			log.Error("%+v", err)
		}
		initCM(c, v, advert, cmRsc)
	}
	//ai已经有dalao卡则直接返回，没有则看要不要第一位拼接其他卡
	if hasDalao == 1 {
		v.Relates = rls
		return
	}
	if dalaoExp != 1 {
		if r, err = s.dealManagerRelate(c, plat, mr, build); err != nil {
			log.Error("%+v", err)
		}
	}
	if r == nil {
		if r, err = s.dealGame(c, plat, cGameID, model.FromOrder); err != nil {
			log.Error("%+v", err)
		}
	}
	if r == nil {
		if r, err = s.dealGame(c, plat, rGameID, model.FromRcmd); err != nil {
			log.Error("%+v", err)
		}
	}
	if r != nil {
		rm = map[int]*view.Relate{0: r}
		aidm = map[int64]struct{}{r.Aid: struct{}{}}
	} else if len(adm) != 0 {
		rm = make(map[int]*view.Relate, len(adm))
		for idx, ad := range adm {
			r = &view.Relate{}
			r.FromCM(ad)
			rm[idx] = r
		}
	}
	if len(rm) != 0 {
		var tmp []*view.Relate
		for _, rl := range rls {
			if _, ok := aidm[rl.Aid]; ok {
				continue
			}
			tmp = append(tmp, rl)
		}
		v.Relates = make([]*view.Relate, 0, len(tmp)+len(rm))
		for _, rl := range tmp {
		LABEL:
			if r, ok := rm[len(v.Relates)]; ok {
				if r.IsAdLoc && r.AdCb == "" {
					rel := &view.Relate{}
					*rel = *rl
					rel.IsAdLoc = r.IsAdLoc
					rel.RequestID = r.RequestID
					rel.SrcID = r.SrcID
					rel.ClientIP = r.ClientIP
					rel.AdIndex = r.AdIndex
					rel.Extra = r.Extra
					rel.CardIndex = r.CardIndex
					v.Relates = append(v.Relates, rel)
				} else if r.Aid != v.Aid {
					v.Relates = append(v.Relates, r)
					goto LABEL
				} else {
					v.Relates = append(v.Relates, rl)
				}
			} else {
				v.Relates = append(v.Relates, rl)
			}
		}
	} else {
		v.Relates = rls
	}
}

func initCM(c context.Context, v *view.View, advert *ad.Ad, resource int64) {
	ads, _ := advert.Convert(resource)
	sort.Sort(ad.AdInfos(ads))
	if len(ads) == 0 {
		return
	}
	v.CMs = make([]*view.CM, 0, len(ads))
	for _, ad := range ads {
		cm := &view.CM{}
		cm.FromCM(ad)
		v.CMs = append(v.CMs, cm)
	}
}

func (s *Service) initMovie(c context.Context, v *view.View, mid int64, build int, mobiApp, device string, nMovie bool) (err error) {
	s.pHit.Incr("is_movie")
	var m *bangumi.Movie
	if m, err = s.banDao.Movie(c, v.Aid, mid, build, mobiApp, device); err != nil || m == nil {
		log.Error("%+v", err)
		err = ecode.NothingFound
		s.pMiss.Incr("err_is_PGC")
		return
	}
	if v.Rights.HD5 == 1 && m.PayUser.Status == 0 && !s.checkVIP(c, mid) {
		v.Rights.HD5 = 0
	}
	if len(m.List) == 0 {
		err = ecode.NothingFound
		return
	}
	vps := make([]*view.Page, 0, len(m.List))
	for _, l := range m.List {
		vp := &view.Page{
			Page3: &archive.Page3{Cid: l.Cid, Page: int32(l.Page), From: l.Type, Part: l.Part, Vid: l.Vid},
		}
		vps = append(vps, vp)
	}
	m.List = nil
	// view
	v.Pages = vps
	v.Rights.Download = int32(m.AllowDownload)
	m.AllowDownload = 0
	v.Rights.Bp = 0
	if nMovie {
		v.Movie = m
		v.Desc = m.Season.Evaluate
	}
	return
}

func (s *Service) initPGC(c context.Context, v *view.View, mid int64, build int, mobiApp, device string) (err error) {
	s.pHit.Incr("is_PGC")
	var season *bangumi.Season
	if season, err = s.banDao.PGC(c, v.Aid, mid, build, mobiApp, device); err != nil {
		log.Error("%+v", err)
		err = ecode.NothingFound
		s.pMiss.Incr("err_is_PGC")
		return
	}
	if season != nil {
		if season.Player != nil {
			if len(v.Pages) != 0 {
				if season.Player.Cid != 0 {
					v.Pages[0].Cid = season.Player.Cid
				}
				if season.Player.From != "" {
					v.Pages[0].From = season.Player.From
				}
				if season.Player.Vid != "" {
					v.Pages[0].Vid = season.Player.Vid
				}
			}
			season.Player = nil
		}
		if season.AllowDownload == "1" {
			v.Rights.Download = 1
		} else {
			v.Rights.Download = 0
		}
		if season.SeasonID != "" {
			season.AllowDownload = ""
			v.Season = season
		}
	}
	if v.Rights.HD5 == 1 && !s.checkVIP(c, mid) {
		v.Rights.HD5 = 0
	}
	v.Rights.Bp = 0
	return
}

func (s *Service) initPages(c context.Context, vs *view.ViewStatic, ap []*archive.Page3) {
	pages := make([]*view.Page, 0, len(ap))
	for _, v := range ap {
		page := &view.Page{}
		metas := make([]*view.Meta, 0, 4)
		for q, r := range _rate {
			meta := &view.Meta{
				Quality: q,
				Size:    int64(float64(r*v.Duration) * 1.1 / 8.0),
			}
			metas = append(metas, meta)
		}
		if vs.AttrVal(archive.AttrBitIsBangumi) == archive.AttrYes {
			v.From = "bangumi"
		}
		page.Page3 = v
		page.Metas = metas
		page.DMLink = fmt.Sprintf(_dmformat, v.Cid)
		pages = append(pages, page)
	}
	vs.Pages = pages
}

func (s *Service) initDownload(c context.Context, v *view.View, mid int64, cdnIP string) (err error) {
	var download int64
	if v.AttrVal(archive.AttrBitLimitArea) == archive.AttrYes {
		if download, err = s.ipLimit(c, mid, v.Aid, cdnIP); err != nil {
			return
		}
	} else {
		download = location.AllowDown
	}
	if download == location.ForbiddenDown {
		v.Rights.Download = int32(download)
		return
	}
	for _, p := range v.Pages {
		if p.From == "qq" {
			download = location.ForbiddenDown
			break
		}
	}
	v.Rights.Download = int32(download)
	return
}

func (s *Service) initUGCPay(c context.Context, v *view.View, plat int8, mid int64) (err error) {
	var (
		asset    *view.Asset
		platform = model.Platform(plat)
	)
	if asset, err = s.ugcpayDao.AssetRelationDetail(c, mid, v.Aid, platform); err != nil {
		log.Error("%+v", err)
		return
	}
	if asset != nil {
		v.Asset = asset
	}
	return
}

func (s *Service) initContributions(c context.Context, v *view.View) {
	const _count = 5
	if v.ReqUser != nil && v.ReqUser.Attention == 1 {
		return
	}
	var hasAudio bool
	for _, page := range v.Pages {
		if page.Audio != nil {
			hasAudio = true
			break
		}
	}
	if hasAudio || v.OwnerExt.Archives < 20 {
		return
	}
	aids, err := s.arcDao.ViewContributeCache(c, v.Author.Mid)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(aids) < _count+1 {
		return
	}
	as, err := s.arcDao.Archives(c, aids)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(as) == 0 {
		return
	}
	ctbt := make([]*view.Contribution, 0, len(as))
	for _, aid := range aids {
		if a, ok := as[aid]; ok && a.IsNormal() {
			if a.Aid != v.Aid {
				vc := &view.Contribution{Aid: a.Aid, Title: a.Title, Pic: a.Pic, Author: a.Author, Stat: a.Stat, CTime: a.PubDate}
				ctbt = append(ctbt, vc)
			}
		}
	}
	if len(ctbt) > _count {
		ctbt = ctbt[:_count]
	}
	if len(ctbt) == _count {
		v.Contributions = ctbt
	}
}

func (s *Service) initAudios(c context.Context, v *view.View) {
	pLen := len(v.Pages)
	if pLen == 0 || pLen > 100 {
		return
	}
	if pLen > 50 {
		pLen = 50
	}
	cids := make([]int64, 0, len(v.Pages[:pLen]))
	for _, p := range v.Pages[:pLen] {
		cids = append(cids, p.Cid)
	}
	vam, err := s.audioDao.AudioByCids(c, cids)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(vam) != 0 {
		for _, p := range v.Pages[:pLen] {
			if va, ok := vam[p.Cid]; ok {
				p.Audio = va
			}
		}
		if len(v.Pages) == 1 {
			if va, ok := vam[v.Pages[0].Cid]; ok {
				v.Audio = va
			}
		}
	}
}

func (s *Service) initElec(c context.Context, v *view.View, mid int64) {
	if _, ok := s.allowTypeIds[int16(v.TypeID)]; !ok || int8(v.Copyright) != archive.CopyrightOriginal {
		return
	}
	info, err := s.elcDao.Info(c, v.Author.Mid, mid)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if info != nil {
		v.Rights.Elec = 1
		info.Show = true
		v.Elec = info
	}
}

func (s *Service) initTag(c context.Context, v *view.View, mid int64, plat int8) (tids []int64) {
	var (
		actTag     []*tag.Tag
		arcTag     []*tag.Tag
		actTagName string
	)
	if v.MissionID > 0 {
		protocol, err := s.actDao.ActProtocol(c, v.MissionID)
		if err != nil {
			log.Error("s.actDao.ActProtocol err(%+v)", err)
			err = nil
		} else {
			if protocol.SubjectItem != nil {
				v.ActivityURL = protocol.SubjectItem.AndroidURL
				if model.IsIOS(plat) {
					v.ActivityURL = protocol.SubjectItem.IosURL
				}
			}
			if protocol.ActSubjectProtocol != nil {
				actTagName = protocol.ActSubjectProtocol.Tags
			}
		}
	}
	tags, err := s.tagDao.ArcTags(c, v.Aid, mid)
	if err != nil {
		log.Error("s.tagDao.ArcTags err(%+v)", err)
		return
	}
	tids = make([]int64, 0, len(tags))
	for _, t := range tags {
		if actTagName == t.Name {
			t.IsActivity = 1
			actTag = append(actTag, t)
		} else {
			arcTag = append(arcTag, t)
		}
		tids = append(tids, t.TagID)
	}
	//活动稿件tag放在第一位
	v.Tag = append(actTag, arcTag...)
	return
}

func (s *Service) initDM(c context.Context, v *view.View) {
	const (
		_dmTypeAv    = 1
		_dmPlatMobie = 1
	)
	pLen := len(v.Pages)
	if pLen == 0 || pLen > 100 {
		return
	}
	if pLen > 50 {
		pLen = 50
	}
	cids := make([]int64, 0, len(v.Pages[:pLen]))
	for _, p := range v.Pages[:pLen] {
		cids = append(cids, p.Cid)
	}
	res, err := s.dmDao.SubjectInfos(c, _dmTypeAv, _dmPlatMobie, cids...)
	if err != nil {
		log.Error("%+v", err)
		return
	}
	if len(res) == 0 {
		return
	}
	for _, p := range v.Pages[:pLen] {
		if r, ok := res[p.Cid]; ok {
			p.DM = r
		}
	}
}

func (s *Service) dealCM(c context.Context, advert *ad.Ad, resource int64) (adm map[int]*ad.AdInfo, err error) {
	ads, aids := advert.Convert(resource)
	if len(ads) == 0 {
		return
	}
	adm = make(map[int]*ad.AdInfo, len(ads))
	for _, ad := range ads {
		adm[ad.CardIndex-1] = ad
	}
	if len(aids) == 0 {
		return
	}
	as, err := s.arcDao.Archives(c, aids)
	if err != nil {
		log.Error("%+v", err)
		err = nil
		return
	}
	for _, ad := range adm {
		if ad.Goto == model.GotoAv && ad.CreativeContent != nil {
			if a, ok := as[ad.CreativeContent.VideoID]; ok {
				ad.View = int(a.Stat.View)
				ad.Danmaku = int(a.Stat.Danmaku)
				if ad.CreativeContent.Desc == "" {
					ad.CreativeContent.Desc = a.Desc
				}
				ad.URI = model.FillURI(ad.Goto, ad.Param, model.AvHandler(a, "", nil))
			}
		}
	}
	return
}

func (s *Service) dealRcmdRelate(c context.Context, plat int8, aid int64, build int) (rls []*view.Relate, aidm map[int64]struct{}, err error) {
	if rls, err = s.arcDao.RelatesCache(c, aid); err != nil {
		return
	}
	if len(rls) != 0 {
		aidm = make(map[int64]struct{}, len(rls))
		for _, rl := range rls {
			if rl.Aid != 0 {
				aidm[rl.Aid] = struct{}{}
			}
		}
		return
	}
	s.prom.Incr("need_relates")
	var aids []int64
	if aids, err = s.arcDao.RelateAids(c, aid); err != nil {
		return
	}
	if len(aids) == 0 {
		return
	}
	var as map[int64]*api.Arc
	if as, err = s.arcDao.Archives(c, aids); err != nil {
		return
	}
	aidm = make(map[int64]struct{}, len(as))
	for _, aid := range aids {
		if a, ok := as[aid]; ok {
			if s.overseaCheck(archive.BuildArchive3(a), plat) || !a.IsNormal() {
				continue
			}
			r := &view.Relate{}
			var cooperation bool
			if (model.IsAndroid(plat) && build > s.c.BuildLimit.CooperationAndroid) || (model.IsIPhone(plat) && build > s.c.BuildLimit.CooperationIOS) {
				cooperation = true
			}
			r.FromAv(a, "", "", nil, cooperation)
			rls = append(rls, r)
			aidm[aid] = struct{}{}
		}
	}
	if len(rls) != 0 {
		s.arcDao.AddRelatesCache(aid, rls)
	}
	return
}

func (s *Service) dealManagerRelate(c context.Context, plat int8, mr *manager.Relate, build int) (r *view.Relate, err error) {
	if mr == nil || mr.Param < 1 {
		return
	}
	var cooperation bool
	if (model.IsAndroid(plat) && build > s.c.BuildLimit.CooperationAndroid) || (model.IsIPhone(plat) && build > s.c.BuildLimit.CooperationIOS) {
		cooperation = true
	}
	switch mr.Goto {
	case model.GotoAv:
		var a *api.Arc
		if a, err = s.arcDao.Archive3(c, mr.Param); err != nil {
			return
		}
		if a != nil {
			r = &view.Relate{}
			r.FromOperateOld(mr, a, nil, nil, model.FromOperation, cooperation)
		}
	case model.GotoGame:
		var info *game.Info
		if info, err = s.gameDao.Info(c, mr.Param, plat); err != nil {
			return
		}
		if info != nil && info.IsOnline {
			r = &view.Relate{}
			r.FromOperateOld(mr, nil, info, nil, model.FromOperation, cooperation)
		}
	case model.GotoSpecial:
		if sp, ok := s.specialCache[mr.Param]; ok {
			r = &view.Relate{}
			r.FromOperateOld(mr, nil, nil, sp, model.FromOperation, cooperation)
		}
	}
	return
}

func (s *Service) dealGame(c context.Context, plat int8, id int64, from string) (r *view.Relate, err error) {
	if id < 1 {
		return
	}
	var info *game.Info
	if info, err = s.gameDao.Info(c, id, plat); err != nil {
		return
	}
	if info != nil && info.IsOnline {
		r = &view.Relate{}
		r.FromGame(info, from)
	}
	return
}

func (s *Service) newRcmdRelate(c context.Context, plat int8, aid, mid int64, buvid, mobiApp, from string, build, qn, fnver, fnval, forceHost, parentMode int) (rls []*view.Relate, userFeature, returnCode string, hasDalao, dalaoExp int, err error) {
	recData, userFeature, returnCode, dalaoExp, err := s.arcDao.NewRelateAids(c, aid, mid, build, parentMode, buvid, from, plat)
	if err != nil || len(recData) == 0 {
		return
	}
	var (
		aids      []int64
		ssIDs     []int32
		gameID    int64
		specialID int64
		arcm      map[int64]*api.Arc
		banm      map[int32]*v1.CardInfoProto
		gameInfo  *game.Info
	)
	for _, rec := range recData {
		switch rec.Goto {
		case model.GotoAv:
			aids = append(aids, rec.Oid)
		case model.GotoBangumi:
			ssIDs = append(ssIDs, int32(rec.Oid))
		case model.GotoGame:
			gameID = rec.Oid
		case model.GotoSpecial:
			specialID = rec.Oid
		}
	}
	eg := errgroup.Group{}
	if len(aids) > 0 {
		eg.Go(func() (err error) {
			if arcm, err = s.arcDao.Archives(context.Background(), aids); err != nil {
				log.Error("s.arcDao.Archives err(%+v)", err)
			}
			return
		})
	}
	if len(ssIDs) > 0 {
		eg.Go(func() (err error) {
			if banm, err = s.banDao.CardsInfoReply(context.Background(), ssIDs); err != nil {
				log.Error("s.banDao.CardsInfoReply err(%+v)", err)
			}
			return
		})
	}
	if gameID > 0 {
		eg.Go(func() (err error) {
			if gameInfo, err = s.gameDao.Info(c, gameID, plat); err != nil {
				log.Error("s.gameDao.Info err(%+v)", err)
			}
			return
		})
	}
	eg.Wait()
	players := make(map[int64]*archive.PlayerInfo)
	if (model.IsAndroid(plat) && build > _qnAndroidBuildGt) || (model.IsIOSNormal(plat) && build > _qnIosBuildGt) || model.IsIPhoneB(plat) {
		var cids []int64
		if aid%100 < s.c.RelateGray {
			for k, v := range aids {
				if k == s.c.RelateCnt {
					break
				}
				if arcm[v].Rights.Autoplay == 1 {
					cids = append(cids, arcm[v].FirstCid)
				}
			}
		}
		if len(cids) > 0 {
			playerInfo := make(map[uint32]*archive.BvcVideoItem)
			if (model.IsAndroid(plat) && build < _qnAndroidBuildLt) || (model.IsIOSNormal(plat) && build <= _qnIosBuildLt) || qn <= 0 {
				qn = _qn480
			}
			if playerInfo, err = s.arcDao.PlayerInfos(c, cids, qn, fnver, fnval, forceHost, mobiApp); err != nil {
				log.Error("%+v", err)
				err = nil
			} else if len(playerInfo) > 0 {
				for k, pi := range playerInfo {
					cid := int64(k)
					players[cid] = new(archive.PlayerInfo)
					players[cid].Cid = pi.Cid
					players[cid].ExpireTime = pi.ExpireTime
					players[cid].FileInfo = make(map[int][]*archive.PlayerFileInfo)
					for qn, files := range pi.FileInfo {
						for _, f := range files.Infos {
							players[cid].FileInfo[int(qn)] = append(players[cid].FileInfo[int(qn)], &archive.PlayerFileInfo{
								FileSize:   f.Filesize,
								TimeLength: f.Timelength,
							})
						}
					}
					players[cid].SupportQuality = pi.SupportQuality
					players[cid].SupportFormats = pi.SupportFormats
					players[cid].SupportDescription = pi.SupportDescription
					players[cid].Quality = pi.Quality
					players[cid].URL = pi.Url
					players[cid].VideoCodecid = pi.VideoCodecid
					players[cid].VideoProject = pi.VideoProject
					players[cid].Fnver = pi.Fnver
					players[cid].Fnval = pi.Fnval
					players[cid].Dash = pi.Dash
				}
			}
		}
	}
	var cooperation bool
	if (model.IsAndroid(plat) && build > s.c.BuildLimit.CooperationAndroid) || (model.IsIPhone(plat) && build > s.c.BuildLimit.CooperationIOS) {
		cooperation = true
	}
	for _, rec := range recData {
		r := &view.Relate{AvFeature: rec.AvFeature, Source: rec.Source, TrackID: rec.TrackID}
		switch rec.Goto {
		case model.GotoAv:
			arc, ok := arcm[rec.Oid]
			if !ok || s.overseaCheck(archive.BuildArchive3(arc), plat) || !arc.IsNormal() {
				continue
			}
			if rec.IsDalao == 1 {
				r.FromOperate(rec, arc, nil, nil, model.FromOperation, cooperation)
			} else {
				r.FromAv(arc, "", rec.TrackID, players[arc.FirstCid], cooperation)
			}
		case model.GotoBangumi:
			ban, ok := banm[int32(rec.Oid)]
			if !ok {
				continue
			}
			r.FromBangumi(ban)
		case model.GotoGame:
			if gameInfo == nil || !gameInfo.IsOnline {
				continue
			}
			r.FromOperate(rec, nil, gameInfo, nil, model.FromOperation, cooperation)
		case model.GotoSpecial:
			sp, ok := s.specialCache[specialID]
			if !ok {
				continue
			}
			r.FromOperate(rec, nil, nil, sp, model.FromOperation, cooperation)
		}
		if rec.IsDalao == 1 {
			hasDalao = 1
		}
		rls = append(rls, r)
	}
	return
}
