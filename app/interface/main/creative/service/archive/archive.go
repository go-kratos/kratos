package archive

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/interface/main/creative/model/activity"
	"go-common/app/interface/main/creative/model/appeal"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/search"
	"go-common/app/interface/main/creative/model/tag"
	pubSvc "go-common/app/interface/main/creative/service"
	"go-common/app/service/main/archive/api"
	mdlarc "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
	"strings"
	"time"
)

// SimpleArchiveVideos fn
func (s *Service) SimpleArchiveVideos(c context.Context, mid, aid int64, ak, ck, ip string) (ap *archive.SimpleArchiveVideos, err error) {
	var (
		sa  *archive.SpArchive
		svs []*archive.SpVideo
	)
	// User identity && permission check, logic from project 'member'
	card, err := s.acc.Card(c, mid, ip)
	if err != nil {
		log.Error("s.acc.Profile(%d,%s) error(%v)", mid, ck, err)
		return
	}
	if sa, err = s.arc.SimpleArchive(c, aid, ip); err != nil {
		log.Error("s.arc.SimpleArchive(%d) error(%v)", aid, err)
		return
	}
	if sa == nil {
		log.Error("s.arc.SimpleArchive(%d) not found", aid)
		err = ecode.NothingFound
		return
	}
	if sa.Mid != mid {
		err = ecode.NothingFound
		return
	}
	if svs, err = s.arc.SimpleVideos(c, aid, ip); err != nil {
		log.Error("s.arc.SimpleVideos(%d) error(%v)", aid, err)
		return
	}
	for index, sv := range svs {
		if sv.Status == -100 {
			sv.DmActive = 0
		} else {
			sv.DmActive = 1
		}
		svs[index] = sv
	}
	acceptAss := card.Rank > 15000 || card.Rank == 20000
	ap = &archive.SimpleArchiveVideos{Archive: sa, SpVideos: svs, AcceptAss: acceptAss}
	return
}

// View get archive.
func (s *Service) View(c context.Context, mid, aid int64, ip, platform string) (av *archive.ArcVideo, err error) {
	if av, err = s.arc.View(c, mid, aid, ip, archive.NeedPoi(platform), archive.NeedVote(platform)); err != nil {
		log.Error("s.arc.View(%d,%d) error(%v)", mid, aid, err)
		return
	}
	if av == nil || av.Archive == nil {
		log.Error("s.arc.View(%d) not found", mid)
		err = ecode.NothingFound
		return
	}
	av.Archive.NilPoiObj(platform)
	av.Archive.NilVote()
	// check white
	isWhite := false
	for _, m := range s.c.Whitelist.ArcMids {
		if m == mid {
			isWhite = true
			break
		}
	}
	if !isWhite {
		if av.Archive.Mid != mid {
			err = ecode.AccessDenied
			return
		}
	}
	var (
		a  = av.Archive
		vs = av.Videos
	)
	// add reject reason
	if (a.State == mdlarc.StateForbidRecicle || a.State == mdlarc.StateForbidLock) && a.RejectReason == "" {
		var c = 0
		for _, v := range vs {
			if v.Status == mdlarc.StateForbidRecicle ||
				v.Status == mdlarc.StateForbidLock ||
				v.Status == mdlarc.StateForbidXcodeFail {
				c++
			}
		}
		a.RejectReason = fmt.Sprintf("稿件中发现%d个问题。", c)
	}
	a.Cover = pubSvc.CoverURL(a.Cover)
	if a.OrderID != 0 {
		if a.OrderID, a.OrderName, _, err = s.order.OrderByAid(c, aid); err != nil {
			log.Error("s.order.OrderByAid(%d,%d) error(%v)", mid, aid, err)
			err = nil
		}
	}
	if a.MissionID > 0 {
		var act *activity.Activity
		if act, err = s.act.Subject(c, a.MissionID); err != nil {
			log.Error("s.act.Subject a.MissionID(%d) error(%v)", a.MissionID, err)
			err = nil
			return
		}
		a.MissionName = act.Name
	}
	if a.Porder != nil && a.Porder.Official == 1 && a.Porder.IndustryID == 1 {
		if v, ok := s.gameMap[a.Porder.BrandID]; ok {
			a.Porder.BrandName = v.GameName
		}
	}
	if a.UgcPay == 1 {
		a.UgcPayInfo = s.p.FillPayInfo(c, a, s.c.UgcPay, ip)
	}
	return
}

// Del fn
// 1,unbind act;2,decrease coin;3,unbind order;4,clean memcache
func (s *Service) Del(c context.Context, mid, aid int64, ip string) (err error) {
	av, err := s.View(c, mid, aid, ip, archive.PlatformWeb)
	if err != nil {
		return
	}
	if av == nil || av.Archive == nil {
		log.Error("s.arc.Del NothingFound (%d,%d,%d) error(%v)", mid, av.Archive.Mid, aid, err)
		err = ecode.NothingFound
		return
	}
	if av.Archive.Mid != mid {
		log.Error("s.arc.Del AccessDenied (%d,%d,%d) error(%v)", mid, av.Archive.Mid, aid, err)
		err = ecode.AccessDenied
		return
	}
	if av.Archive.UgcPay == 1 {
		canDelTime := xtime.Time(av.Archive.PTime.Time().AddDate(0, 0, s.c.UgcPay.AllowDeleteDays).Unix())
		if av.Archive.CTime != av.Archive.PTime && xtime.Time(time.Now().Unix()) < canDelTime {
			log.Error("checkEditPay CreativePayForbidDeleteAfterOpen aid(%d) ctime(%v) ptime(%v)", av.Archive.Aid, av.Archive.CTime, av.Archive.PTime)
			err = ecode.CreativePayForbidDeleteAfterOpen
			return
		}
	}
	if err = s.arc.Del(c, mid, aid, ip); err != nil {
		log.Error("s.arc.Del(%d,%d) error(%v)", mid, aid, err)
		return
	}
	var (
		a   = av.Archive
		g   = &errgroup.Group{}
		ctx = context.TODO()
	)
	g.Go(func() error {
		if e := s.act.Unbind(ctx, aid, a.MissionID, ip); e != nil {
			log.Error("s.act.UpdateByAid(%d,%d) error(%v)", aid, a.MissionID, e)
		}
		return nil
	})
	g.Go(func() error {
		var coins float64
		if a.State >= 0 {
			coins = -2
		} else {
			coins = -1
		}
		if e := s.coin.AddCoin(ctx, mid, aid, coins, ip); e != nil {
			log.Error("s.coin.AddCoin(%d,%d,%f,%s) error(%v)", mid, aid, coins, ip, e)
		}
		return nil
	})
	g.Go(func() error {
		if e := s.order.Unbind(ctx, mid, aid, ip); e != nil {
			log.Error("s.order.Unbind(%d,%d,%s) error(%v)", mid, aid, ip, e)
		}
		return nil
	})
	g.Go(func() error {
		if e := s.arc.DelSubmitCache(ctx, mid, a.Title); e != nil {
			log.Error("s.arc.DelSubmitCache (%d,%s,%s) error(%v)", mid, a.Title, ip, e)
		}
		return nil
	})
	g.Wait()
	s.prom.Incr("archive_del")
	return
}

// Archives for app all achives.
func (s *Service) Archives(c context.Context, mid int64, tid int16, kw, order, class, ip string, pn, ps, coop int) (res *search.Result, err error) {
	if res, err = s.Search(c, mid, tid, kw, order, class, ip, pn, ps, coop); err != nil {
		log.Error("s.Search err(%v)", err)
		return
	}
	if res == nil || res.Archives == nil || len(res.Archives) == 0 {
		return
	}
	for _, av := range res.Archives {
		a := &api.Arc{
			Aid:       av.Archive.Aid,
			TypeID:    int32(av.Archive.TypeID),
			TypeName:  av.TypeName,
			Copyright: int32(av.Archive.Copyright),
			Title:     av.Archive.Title,
			Desc:      av.Archive.Desc,
			Attribute: av.Archive.Attribute,
			Videos:    int64(len(av.Videos)),
			Pic:       pubSvc.CoverURL(av.Archive.Cover),
			State:     int32(av.Archive.State),
			Access:    int32(av.Archive.Access),
			Tags:      strings.Split(av.Archive.Tag, ","),
			Duration:  av.Archive.Duration,
			MissionID: av.Archive.MissionID,
			OrderID:   av.Archive.OrderID,
			PubDate:   av.Archive.PTime,
			Ctime:     av.Archive.CTime,
		}
		a.Author = api.Author{
			Mid:  av.Archive.Mid,
			Name: av.Archive.Author,
		}
		a.Stat = api.Stat{
			Aid:     av.Stat.Aid,
			View:    int32(av.Stat.View),
			Danmaku: int32(av.Stat.Danmaku),
			Reply:   int32(av.Stat.Reply),
			Fav:     int32(av.Stat.Fav),
			Coin:    int32(av.Stat.Coin),
			Share:   int32(av.Stat.Share),
			NowRank: int32(av.Stat.NowRank),
			HisRank: int32(av.Stat.HisRank),
		}
		ava := &archive.OldArchiveVideoAudit{
			Arc:         a,
			StatePanel:  av.StatePanel,
			ParentTName: av.ParentTName,
			Dtime:       av.Archive.DTime,
			StateDesc:   av.Archive.StateDesc,
			RejectReson: av.Archive.RejectReason,
			UgcPay:      av.Archive.UgcPay,
			Attrs:       av.Archive.Attrs,
		}
		if len(av.ArcVideo.Videos) > 0 {
			var vas = make([]*archive.OldVideoAudit, 0, len(av.ArcVideo.Videos))
			for _, vd := range av.ArcVideo.Videos {
				va := &archive.OldVideoAudit{
					IndexOrder: vd.Index,
					Eptitle:    vd.Title,
					Reason:     vd.RejectReason,
				}
				ava.VideoAudits = append(vas, va)
			}
		}
		res.OldArchives = append(res.OldArchives, ava)
	}
	return
}

// Search all achive.
func (s *Service) Search(c context.Context, mid int64, tid int16, keyword, order, class, ip string, pn, ps, coop int) (res *search.Result, err error) {
	defer func() {
		if res != nil {
			// get pending apply num
			if coop > 0 {
				count, _ := s.arc.CountByMID(c, mid)
				res.Applies = &search.ApplyStateCount{
					Pending: count,
				}
			}
		}
	}()
	if res, err = s.sear.ArchivesES(c, mid, tid, keyword, order, class, ip, pn, ps, coop); err != nil {
		// search err, use archive-service
		log.Error("s.arc.Search(%d) error(%v)", mid, err)
		res, err = s.ArchivesFromService(c, mid, class, ip, pn, ps)
		if err != nil {
			return
		}
	}
	// add arctype  *must return
	res.ArrType = []*search.TypeCount{}
	for _, v := range s.p.TopTypesCache {
		t := &search.TypeCount{
			Tid:  v.ID,
			Name: v.Name,
		}
		if _, ok := res.Type[v.ID]; ok {
			t.Count = res.Type[v.ID].Count
		}
		res.ArrType = append(res.ArrType, t)
	}
	if res == nil || len(res.Aids) == 0 {
		return
	}
	avm, err := s.arc.Views(c, mid, res.Aids, ip)
	if err != nil {
		log.Error("s.arc.Views res.Aids(%v), ip(%s) err(%v)", res.Aids, ip, err)
		return
	}
	// get arc stats
	as, _ := s.arc.Stats(c, res.Aids, ip)
	// archives
	for _, aid := range res.Aids {
		av := avm[aid]
		if av == nil {
			continue
		}
		a := av.Archive
		a.Cover = pubSvc.CoverURL(a.Cover)
		ava := &archive.ArcVideoAudit{ArcVideo: av}
		// arc stat info
		if _, ok := as[aid]; ok {
			ava.Stat = as[aid]
		} else {
			ava.Stat = &api.Stat{}
		}
		// typename
		if _, ok := s.p.TypeMapCache[a.TypeID]; ok {
			ava.TypeName = s.p.TypeMapCache[a.TypeID].Name
		}
		// parent typename
		if _, ok := s.p.TypeMapCache[a.TypeID]; ok {
			if _, ok := s.p.TypeMapCache[s.p.TypeMapCache[a.TypeID].Parent]; ok {
				ava.ParentTName = s.p.TypeMapCache[s.p.TypeMapCache[a.TypeID].Parent].Name
			}
		}
		// state panel
		ava.StatePanel = archive.StatePanel(a.State)
		// state desc
		ava.Archive.StateDesc = s.c.StatDesc(int(a.State))
		// not pubbed videos for reason
		unpubedVideos := make([]*archive.Video, 0, len(ava.Videos))
		for _, v := range ava.Videos {
			if v.Status == -2 || v.Status == -4 || v.Status == -16 {
				unpubedVideos = append(unpubedVideos, v)
			}
		}
		c := len(unpubedVideos)
		if c > 0 {
			a.RejectReason = fmt.Sprintf("稿件中发现%d个问题。", c)
		}
		// set attrs
		attrs := &archive.Attrs{
			IsCoop:  int8(ava.AttrVal(archive.AttrBitIsCoop)),
			IsOwner: ava.IsOwner(mid),
		}
		ava.Archive.Attrs = attrs
		ava.Videos = unpubedVideos
		res.Archives = append(res.Archives, ava)
	}
	return
}

// ApplySearch all achive.
func (s *Service) ApplySearch(c context.Context, mid int64, tid int16, keyword, state string, pn, ps int) (res *search.StaffApplyResult, err error) {
	var mids []int64
	if res, err = s.sear.ArchivesStaffES(c, mid, tid, keyword, state, pn, ps); err != nil {
		log.Error("s.arc.ArchivesStaffES(%d) error(%v)", mid, err)
		return
	}
	// add arctype  *must return
	res.ArrType = []*search.TypeCount{}
	for _, v := range s.p.TopTypesCache {
		t := &search.TypeCount{
			Tid:  v.ID,
			Name: v.Name,
		}
		if _, ok := res.Type[v.ID]; ok {
			t.Count = res.Type[v.ID].Count
		}
		res.ArrType = append(res.ArrType, t)
	}
	if res == nil || len(res.Aids) == 0 {
		return
	}
	// get archives
	avm, err := s.arc.Views(c, mid, res.Aids, "")
	if err != nil {
		log.Error("s.arc.Views res.Aids(%v), ip(%s) err(%v)", res.Aids, err)
		return
	}
	// get applies
	apm := make(map[int64]*archive.StaffApply)
	applies, err := s.arc.StaffApplies(c, mid, res.Aids)
	for _, ap := range applies {
		apm[ap.ApplyAID] = ap
	}
	// combine
	for _, aid := range res.Aids {
		av := avm[aid]
		if av == nil {
			continue
		}
		a := av.Archive
		a.Cover = pubSvc.CoverURL(a.Cover)
		ava := &archive.ArcVideoAudit{ArcVideo: av}
		// typename
		if _, ok := s.p.TypeMapCache[a.TypeID]; ok {
			ava.TypeName = s.p.TypeMapCache[a.TypeID].Name
		}
		// parent typename
		if _, ok := s.p.TypeMapCache[a.TypeID]; ok {
			if _, ok := s.p.TypeMapCache[s.p.TypeMapCache[a.TypeID].Parent]; ok {
				ava.ParentTName = s.p.TypeMapCache[s.p.TypeMapCache[a.TypeID].Parent].Name
			}
		}
		// state panel
		ava.StatePanel = archive.StatePanel(a.State)
		// state desc
		ava.Archive.StateDesc = s.c.StatDesc(int(a.State))
		// not pubbed videos for reason
		unpubedVideos := make([]*archive.Video, 0, len(ava.Videos))
		for _, v := range ava.Videos {
			if v.Status == -2 || v.Status == -4 || v.Status == -16 {
				unpubedVideos = append(unpubedVideos, v)
			}
		}
		c := len(unpubedVideos)
		if c > 0 {
			a.RejectReason = fmt.Sprintf("稿件中发现%d个问题。", c)
		}
		ava.Videos = unpubedVideos
		// get apply
		apply := &search.StaffApply{}
		if v, ok := apm[aid]; ok {
			apply.ID = v.ID
			apply.Type = v.Type
			apply.Mid = v.ApplyUpMID
			apply.State = v.StaffState
			apply.ApplyState = v.State
			apply.ApplyTitle = v.ApplyTitle
		}
		// set archive
		apply.Archive = ava
		mids = append(mids, apply.Mid)
		res.Applies = append(res.Applies, apply)
	}
	// get name
	users, _ := s.acc.Infos(c, mids, "")
	for _, v := range res.Applies {
		if u, ok := users[v.Mid]; ok {
			v.Uname = u.Name
		}
	}
	return
}

// Types get typelist.
func (s *Service) Types(c context.Context, lang string) (tps []*archive.Type) {
	if _, ok := s.p.TypesCache[lang]; !ok {
		lang = "ch"
	}
	tps = s.p.TypesCache[lang]
	return
}

// StaffTitles get staff titles.
func (s *Service) StaffTitles(c context.Context) (titles []*tag.StaffTitle) {
	return s.p.StaffTitlesCache
}

// AppTypes fn
func (s *Service) AppTypes(c context.Context, lang string) (tps []*archive.Type) {
	if _, ok := s.p.CTypesCache[lang]; !ok {
		lang = "ch"
	}
	tps = s.p.CTypesCache[lang]
	for _, val := range tps {
		for _, child := range val.Children {
			child.Notice = child.AppNotice
		}
	}
	return
}

// Activities get activity list.
func (s *Service) Activities(c context.Context) (acts []*activity.Activity) {
	acts = s.p.ActVideoAllCache
	return
}

// WebArchives achive list with appeal.
func (s *Service) WebArchives(c context.Context, mid int64, tid int16, keyword, order, class, ip string, pn, ps, coop int) (res *search.Result, err error) {
	if res, err = s.Search(c, mid, tid, keyword, order, class, ip, pn, ps, coop); err != nil {
		log.Error("s.Search err(%v)", err)
		return
	}
	if res == nil || len(res.Aids) == 0 {
		return
	}
	//做降级
	aps, err := s.ap.AppealList(c, mid, appeal.Business, ip)
	if err != nil {
		log.Error("s.ap.AppealList error(%v)", err)
		err = nil
	}
	if len(aps) == 0 {
		return
	}
	aaMap := make(map[int64]int64, len(aps))
	for _, v := range aps {
		if appeal.IsOpen(v.BusinessState) {
			aaMap[v.Oid] = v.ID
		}
	}
	for _, v := range res.Archives {
		v.OpenAppeal = aaMap[v.Archive.Aid]
	}
	return
}

// Videos get Simple Archive and Videos Info.
func (s *Service) Videos(c context.Context, mid, aid int64, ip string) (sa archive.SimpleArchive, svs []*archive.SimpleVideo, err error) {
	var av *archive.ArcVideo
	if av, err = s.arc.View(c, mid, aid, ip, 0, 0); err != nil {
		log.Error("s.arc.View(%d,%d) error(%v)", mid, aid, err)
		return
	}
	if av == nil {
		log.Error("s.arc.View(%d) not found", mid)
		err = ecode.RequestErr
		return
	}
	// white list check
	isWhite := false
	for _, m := range s.c.Whitelist.DataMids {
		if m == mid {
			isWhite = true
			break
		}
	}
	var (
		a  = av.Archive
		vs = av.Videos
	)
	if !isWhite {
		if a.Mid != mid {
			err = ecode.AccessDenied
			return
		}
	}
	sa.Aid = a.Aid
	sa.Title = a.Title
	for _, v := range vs {
		svs = append(svs, &archive.SimpleVideo{
			Cid:   v.Cid,
			Title: v.Title,
			Index: v.Index,
		})
	}
	return
}

// ArchivesFromService get archives from service
func (s *Service) ArchivesFromService(c context.Context, mid int64, class, ip string, pn, ps int) (sres *search.Result, err error) {
	aids, count, err := s.arc.UpArchives(c, mid, int64(pn), int64(ps), 0, ip)
	if err != nil {
		return
	}
	sres = &search.Result{}
	sres.Aids = aids
	sres.Page.Pn = pn
	sres.Page.Ps = ps
	sres.Page.Count = int(count)
	return
}

// DescFormat get desc format
func (s *Service) DescFormat(c context.Context, typeid, copyright int64, langStr, ip string) (desc *archive.DescFormat, err error) {
	lang := archive.ToLang(langStr)
	desc, ok := s.p.DescFmtsCache[typeid][int8(copyright)][lang]
	if !ok {
		err = nil
	}
	if desc != nil {
		var Components []*struct {
			Name interface{}
		}
		if err = json.Unmarshal([]byte(desc.Components), &Components); err != nil || len(Components) == 0 {
			desc = nil
			err = nil
		}
	}
	return
}

// AppFormats for app portal list.
func (s *Service) AppFormats(c context.Context) (af []*archive.AppFormat, err error) {
	for _, f := range s.p.DescFmtsArrCache {
		format := &archive.AppFormat{ID: f.ID, Copyright: f.Copyright, TypeID: f.TypeID}
		af = append(af, format)
	}
	return
}

// Video get video by aid and cid
func (s *Service) Video(c context.Context, mid, aid, cid int64, ip string) (video *api.Page, err error) {
	if video, err = s.arc.Video(c, aid, cid, ip); err != nil {
		log.Error("s.arc.Video %d,%d,%s | err(%v)", aid, cid, ip, err)
	}
	return
}

// VideoJam get video traffic jam level from service
// level为0的时候，可以忽略错误处理，映射表里已经做了前端容错
func (s *Service) VideoJam(c context.Context, ip string) (j *archive.VideoJam, err error) {
	level, _ := s.arc.VideoJam(c, ip)
	if jam, ok := archive.VjInfo[level]; ok {
		j = jam
	} else {
		j = archive.VjInfo[0]
	}
	return
}

// DescFormatForApp get desc format length
func (s *Service) DescFormatForApp(c context.Context, typeid, copyright int64, langStr, ip string) (desc *archive.DescFormat, length int, err error) {
	var (
		descLengthMax = 2000
		descLengthMin = 250
		ok            bool
	)
	lang := archive.ToLang("")
	if desc, ok = s.p.DescFmtsCache[typeid][int8(copyright)][lang]; !ok {
		err = nil
	}
	if typeid == 0 {
		length = descLengthMin
	} else if desc != nil {
		length = descLengthMax
	} else {
		length = descLengthMin
	}
	return
}

// Dpub for app view.
func (s *Service) Dpub() (dpub *archive.Dpub) {
	now := time.Now()
	dpub = &archive.Dpub{
		Deftime:    xtime.Time(now.Add(time.Duration(14400) * time.Second).Unix()),
		DeftimeEnd: xtime.Time(now.Add(time.Duration(15*24) * time.Hour).Unix()),
		DeftimeMsg: ecode.String(ecode.VideoupDelayTimeErr.Error()).Message(),
	}
	return
}

// SimpleArcVideos fn
func (s *Service) SimpleArcVideos(c context.Context, mid int64, tid int16, kw, order, class, ip string, pn, ps, coop int) (res *search.SimpleResult, err error) {
	var sres *search.Result
	res = &search.SimpleResult{}
	if sres, err = s.sear.ArchivesES(c, mid, tid, kw, order, class, ip, pn, ps, coop); err != nil {
		log.Error("s.arc.Search mid(%d)|error(%v)", mid, err)
		sres, err = s.ArchivesFromService(c, mid, class, ip, pn, ps) // search err, use archive-service
		if err != nil {
			log.Error("s.ArchivesFromService mid(%d)|error(%v)", mid, err)
			return
		}
	}
	if sres == nil || len(sres.Aids) == 0 {
		return
	}
	res.Class = sres.Class
	res.Page = sres.Page
	avm, err := s.arc.Views(c, mid, sres.Aids, ip)
	if err != nil {
		log.Error("s.arc.Views mid(%d)|aids(%v)|ip(%s)|err(%v)", mid, sres.Aids, ip, err)
		return
	}
	savs := make([]*search.SimpleArcVideos, 0, len(avm))
	for _, aid := range sres.Aids {
		av, ok := avm[aid]
		if !ok || av == nil || av.Archive == nil {
			continue
		}
		vds := make([]*archive.SimpleVideo, 0, len(av.Videos))
		for _, v := range av.Videos {
			if v == nil {
				continue
			}
			vd := &archive.SimpleVideo{
				Cid:   v.Cid,
				Title: v.Title,
				Index: v.Index,
			}
			vds = append(vds, vd)
		}
		sav := &search.SimpleArcVideos{}
		sav.Archive = &archive.SimpleArchive{
			Aid:   av.Archive.Aid,
			Title: av.Archive.Title,
		}
		sav.Videos = vds
		savs = append(savs, sav)
	}
	res.ArchivesVideos = savs
	return
}
