package music

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/archive"
	material "go-common/app/interface/main/creative/dao/material"
	music "go-common/app/interface/main/creative/dao/music"
	"go-common/app/interface/main/creative/dao/up"
	appMdl "go-common/app/interface/main/creative/model/app"
	mMdl "go-common/app/interface/main/creative/model/music"
	sMdl "go-common/app/interface/main/creative/model/search"
	"go-common/app/interface/main/creative/service"
	accMdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/queue/databus/report"
	"go-common/library/sync/errgroup"
	"sort"
	"time"
)

const hotvalArcFactor = int(357)

//Service struct.
type Service struct {
	c                    *conf.Config
	music                *music.Dao
	material             *material.Dao
	acc                  *account.Dao
	archive              *archive.Dao
	up                   *up.Dao
	LatestBgm            *mMdl.Music
	MscWithTypes         map[int][]*mMdl.Music
	AllMsc               map[int64]*mMdl.Music
	Types                []*mMdl.Category
	Subtitles            []*mMdl.Subtitle
	Fonts                []*mMdl.Font
	Filters              []*mMdl.Filter
	FilterWithCategory   []*mMdl.FilterCategory
	VstickerWithCategory []*mMdl.VstickerCategory
	Hotwords             []*mMdl.Hotword
	Stickers             []*mMdl.Sticker
	Intros               []*mMdl.Intro
	Vstickers            []*mMdl.VSticker
	Transitions          []*mMdl.Transition
	// from app535 for sticker whitelist
	stickerUps map[int64]int64
	Themes     []*mMdl.Theme
	Cooperates []*mMdl.Cooperate
	p          *service.Public
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:                    c,
		music:                music.New(c),
		material:             material.New(c),
		acc:                  rpcdaos.Acc,
		archive:              rpcdaos.Arc,
		up:                   rpcdaos.Up,
		Types:                make([]*mMdl.Category, 0),
		MscWithTypes:         make(map[int][]*mMdl.Music),
		AllMsc:               make(map[int64]*mMdl.Music),
		Subtitles:            make([]*mMdl.Subtitle, 0),
		Fonts:                make([]*mMdl.Font, 0),
		Filters:              make([]*mMdl.Filter, 0),
		FilterWithCategory:   make([]*mMdl.FilterCategory, 0),
		VstickerWithCategory: make([]*mMdl.VstickerCategory, 0),
		Hotwords:             make([]*mMdl.Hotword, 0),
		Stickers:             make([]*mMdl.Sticker, 0),
		Intros:               make([]*mMdl.Intro, 0),
		Vstickers:            make([]*mMdl.VSticker, 0),
		Transitions:          make([]*mMdl.Transition, 0),
		Themes:               make([]*mMdl.Theme, 0),
		Cooperates:           make([]*mMdl.Cooperate, 0),
		p:                    p,
	}
	s.loadUpSpecialSticker()
	s.loadPreValues()
	s.loadMaterial()
	go s.loadproc()
	return s
}

// loadUpSpecialSticker 拍摄贴纸支持灰度分发，配合Sticker.white字段使用
func (s *Service) loadUpSpecialSticker() {
	ups, err := s.up.UpSpecial(context.TODO(), 18)
	if err != nil {
		return
	}
	s.stickerUps = ups
}
func (s *Service) loadproc() {
	for {
		time.Sleep(time.Duration(2 * time.Minute))
		s.loadPreValues()
		s.loadMaterial()
		s.loadUpSpecialSticker()
	}
}

func (s *Service) loadMaterial() {
	var (
		basic = make(map[string]interface{})
		err   error
		g     = &errgroup.Group{}
		ctx   = context.TODO()
	)
	c := context.TODO()
	if basic, err = s.material.Basic(c); err != nil {
		log.Error("s.material.basic err(%+v)", err)
		return
	}
	if subtitles, ok := basic["subs"].([]*mMdl.Subtitle); ok {
		log.Info("s.material.subs (%+d)", len(subtitles))
		s.Subtitles = subtitles
	}
	if fonts, ok := basic["fons"].([]*mMdl.Font); ok {
		log.Info("s.material.fons (%+d)", len(fonts))
		s.Fonts = fonts
	}
	if hotwords, ok := basic["hots"].([]*mMdl.Hotword); ok {
		log.Info("s.material.hots (%+d)", len(hotwords))
		s.Hotwords = hotwords
	}
	if stickers, ok := basic["stis"].([]*mMdl.Sticker); ok {
		log.Info("s.material.stis (%+d)", len(stickers))
		s.Stickers = stickers
	}
	if intros, ok := basic["ints"].([]*mMdl.Intro); ok {
		log.Info("s.material.ints (%+d)", len(intros))
		s.Intros = intros
	}
	if trans, ok := basic["trans"].([]*mMdl.Transition); ok {
		log.Info("s.material.trans (%+d)", len(trans))
		s.Transitions = trans
	}
	if themes, ok := basic["themes"].([]*mMdl.Theme); ok {
		log.Info("s.material.themes (%+d)", len(themes))
		s.Themes = themes
	}
	g.Go(func() error {
		s.getFilterAndItsCategory(ctx)
		return nil
	})
	g.Go(func() error {
		s.getVStickerAndItsCategory(ctx)
		return nil
	})
	g.Go(func() error {
		s.getCooperates(ctx)
		return nil
	})
	g.Wait()
}

func (s *Service) loadPreValues() {
	var (
		mcategory    []*mMdl.Mcategory
		resCategorys []*mMdl.Category
		tids         []int64
		sids         []int64
		err          error
		musicMap     map[int64]*mMdl.Music
		// latestBgm    *mMdl.Music
		jointimeMap       = make(map[int64]int64)
		jointimes         = make([]int64, 0)
		c                 = context.TODO()
		sidTidMapIdx      = make(map[int64][](map[int]int))
		sidTidMapJoinUnix = make(map[int64][](map[int]int64))
		categoryMaps      = make(map[int]*mMdl.Category)
	)
	if mcategory, err = s.music.MCategorys(c); err != nil {
		log.Error("s.music.MCategorys err(%+v)", err)
		return
	}
	for _, v := range mcategory {
		joinUnix := v.CTime.Time().Unix()
		jointimeMap[joinUnix] = v.SID
		jointimes = append(jointimes, joinUnix)
		tid := v.Tid
		tidx := v.Index
		sid := v.SID
		tids = append(tids, int64(tid))
		sids = append(sids, sid)
		if _, ok := sidTidMapIdx[sid]; !ok {
			sidTidMapIdx[sid] = make([](map[int]int), 0)
		}
		sidTidMapIdx[sid] = append(sidTidMapIdx[sid], map[int]int{
			tid: tidx,
		})
		if _, ok := sidTidMapJoinUnix[sid]; !ok {
			sidTidMapJoinUnix[sid] = make([](map[int]int64), 0)
		}
		sidTidMapJoinUnix[sid] = append(sidTidMapJoinUnix[sid], map[int]int64{
			tid: joinUnix,
		})
	}
	if len(tids) > 0 {
		if resCategorys, categoryMaps, err = s.music.Categorys(c, tids); err != nil {
			log.Error("s.music.Categorys tids(%+v)|err(%+v)", tids, err)
			return
		}
		if len(resCategorys) > 0 && len(sids) > 0 {
			if musicMap, err = s.music.Music(c, sids); err != nil {
				log.Error("s.music.Music tids(%+v)|sids(%+v)|err(%+v)", tids, sids, err)
				return
			}
			// get last jointime bgm
			if len(jointimes) > 0 {
				sort.Slice(jointimes, func(i, j int) bool {
					return jointimes[i] >= jointimes[j]
				})
				lastJoinUnix := jointimes[0]
				lastSid := jointimeMap[lastJoinUnix]
				if bgm, ok := musicMap[lastSid]; ok {
					s.LatestBgm = bgm
				}
			}
			s.AllMsc = musicMap
			bgms := make(map[int][]*mMdl.Music)
			upNamesMap, _ := s.getUpNames(c, musicMap)
			for sid, msc := range musicMap {
				if tpsSlice, okM := sidTidMapIdx[sid]; okM {
					for _, tpIdx := range tpsSlice {
						for tp, idx := range tpIdx {
							if _, ok := categoryMaps[tp]; !ok {
								continue
							}
							if _, ok := bgms[tp]; !ok {
								bgms[tp] = make([]*mMdl.Music, 0)
							}
							var (
								junix  int64
								upName string
							)
							if name, okU := upNamesMap[msc.UpMID]; okU {
								upName = name
							}
							for _, Jmap := range sidTidMapJoinUnix[sid] {
								if joinUnix, okJ := Jmap[tp]; okJ {
									junix = joinUnix
								}
							}
							bgm := &mMdl.Music{
								ID:             msc.ID,
								TID:            tp,
								Index:          idx,
								SID:            msc.SID,
								Name:           msc.Name,
								Musicians:      upName,
								UpMID:          msc.UpMID,
								Cover:          msc.Cover,
								Stat:           msc.Stat,
								Playurl:        msc.Playurl,
								Duration:       msc.Duration,
								FileSize:       msc.FileSize,
								CTime:          msc.CTime,
								MTime:          msc.MTime,
								Pubtime:        msc.Pubtime,
								Tl:             msc.Tl,
								RecommendPoint: msc.RecommendPoint,
								Cooperate:      msc.Cooperate,
								CooperateURL:   msc.CooperateURL,
							}
							// freash bgm.Tags by joinunix
							if junix+86400*7 >= time.Now().Unix() {
								bgm.New = 1
								bgm.Tags = []string{"NEW"}
							} else {
								bgm.Tags = make([]string, 0)
							}
							if len(msc.Tags) > 0 {
								bgm.Tags = append(bgm.Tags, msc.Tags...)
							}
							topLen := 3
							if len(bgm.Tags) > topLen {
								bgm.Tags = bgm.Tags[:topLen]
							}
							bgms[tp] = append(bgms[tp], bgm)
						}
					}
				}
			}
			if len(bgms) > 0 {
				s.MscWithTypes = bgms
			}
			var filterdCategorys []*mMdl.Category
			for _, t := range resCategorys {
				if len(bgms[t.ID]) > 0 {
					sort.Slice(bgms[t.ID], func(i, j int) bool {
						return bgms[t.ID][i].Index < bgms[t.ID][j].Index
					})
					t.Children = bgms[t.ID]
					filterdCategorys = append(filterdCategorys, t)
				}
			}
			s.Types = filterdCategorys
		}
	}
	log.Info("loadPreValues (%d)|(%d)", len(s.MscWithTypes), len(s.Types))
}

// BgmExt fn
func (s *Service) BgmExt(c context.Context, mid, sid int64) (ret *mMdl.BgmExt, err error) {
	var (
		bgm              *mMdl.Music
		upMid            int64
		shouldFollowMids []int64
		g                = &errgroup.Group{}
		ctx              = context.TODO()
		extSidMap        map[int]int64
	)
	log.Warn("BgmExt allMsc(%d)", len(s.AllMsc))
	if v, ok := s.AllMsc[sid]; ok {
		bgm = v
		ret = &mMdl.BgmExt{
			Msc: bgm,
		}
	} else {
		return
	}
	log.Warn("BgmExt finish find bgm allMsc(%d)|bgm(%+v)", len(s.AllMsc), bgm)
	if bgm == nil {
		return
	}
	log.Warn("BgmExt bgm info (%+v)", bgm)
	upMid = bgm.UpMID
	ret.ExtMscs, extSidMap = s.UperOtherBgmsFromRecom(ctx, upMid, sid)
	ip := metadata.String(c, metadata.RemoteIP)
	// step 1: get ext aids and rpc archives
	g.Go(func() error {
		ret.ExtArcs, ret.Msc.Hotval, err = s.ExtArcsWithSameBgm(ctx, sid)
		log.Warn("BgmExt step 1: extArcs(%+v)|hotVal(%+v)|sid(%+d)|err(%+v)", ret.ExtArcs, ret.Msc.Hotval, sid, err)
		return nil
	})
	// step 2: get ext mscs
	g.Go(func() error {
		ret.ExtMscs, err = s.ExtSidHotMapAndSort(ctx, ret.ExtMscs, extSidMap)
		log.Warn("BgmExt step 2: ExtMscs(%+v)|extSidMap(%+v)|err(%+v)", ret.ExtMscs, extSidMap, err)
		return nil
	})
	// step 3: get up info and if should follow
	g.Go(func() error {
		ret.UpProfile, err = s.acc.Profile(ctx, upMid, ip)
		log.Warn("BgmExt step 3: profile mid(%+v)", upMid)
		return nil
	})
	// step 4: get up info and if should follow
	g.Go(func() error {
		shouldFollowMids, err = s.acc.ShouldFollow(ctx, mid, []int64{upMid}, ip)
		if len(shouldFollowMids) == 1 {
			ret.ShouldFollow = true
		}
		log.Warn("BgmExt step 4: shouldFollow(%+v)", ret.ShouldFollow)
		return nil
	})
	g.Wait()
	return
}

// ExtSidHotMapAndSort fn,  sorry ...
func (s *Service) ExtSidHotMapAndSort(c context.Context, ExtMscs []*mMdl.Music, extSidMap map[int]int64) (res []*mMdl.Music, err error) {
	var (
		s1total, s2total, s3total int
		s1hot, s2hot, s3hot       int
	)
	if len(extSidMap) > 0 {
		var (
			g   = &errgroup.Group{}
			ctx = context.TODO()
		)
		if sid, ok := extSidMap[1]; ok {
			g.Go(func() error {
				if _, s1total, err = s.music.ExtAidsWithSameBgm(ctx, sid, 1); err != nil {
					log.Error("ExtAidsWithSameBgm S1 error(%v)", err)
				}
				s1hot = s1total * hotvalArcFactor
				return nil
			})
		}
		if sid, ok := extSidMap[2]; ok {
			g.Go(func() error {
				if _, s2total, err = s.music.ExtAidsWithSameBgm(ctx, sid, 1); err != nil {
					log.Error("ExtAidsWithSameBgm S2 error(%v)", err)
				}
				s2hot = s2total * hotvalArcFactor
				return nil
			})
		}
		if sid, ok := extSidMap[3]; ok {
			g.Go(func() error {
				if _, s3total, err = s.music.ExtAidsWithSameBgm(ctx, sid, 1); err != nil {
					log.Error("ExtAidsWithSameBgm S3 error(%v)", err)
				}
				s3hot = s3total * hotvalArcFactor
				return nil
			})
		}
		g.Wait()
	}
	for idx, v := range ExtMscs {
		if idx == 0 {
			v.Hotval = s1hot
		}
		if idx == 1 {
			v.Hotval = s2hot
		}
		if idx == 2 {
			v.Hotval = s3hot
		}
	}
	sort.Slice(ExtMscs, func(i, j int) bool {
		return ExtMscs[i].Hotval > ExtMscs[j].Hotval
	})
	res = ExtMscs
	return
}

// ExtArcsWithSameBgm fn
func (s *Service) ExtArcsWithSameBgm(c context.Context, sid int64) (res []*api.Arc, hot int, err error) {
	var (
		arcMap map[int64]*api.Arc
		aids   []int64
		total  int
	)
	aids, total, err = s.music.ExtAidsWithSameBgm(c, sid, 100)
	if len(aids) > 0 {
		ip := metadata.String(c, metadata.RemoteIP)
		if arcMap, err = s.archive.Archives(c, aids, ip); err != nil {
			log.Error("s.archive.Archives Stats (%v) error(%v)", aids, err)
			err = ecode.CreativeArcServiceErr
			return
		}
		for _, aid := range aids {
			if arc, ok := arcMap[aid]; ok && arc.State >= 0 {
				res = append(res, arc)
			}
		}
	}
	topLen := 20
	if len(res) > topLen {
		res = res[:topLen]
	}
	hot = total * hotvalArcFactor
	return
}

// UperOtherBgmsFromRecom fn, 最多三个同一个Up主的bgms
func (s *Service) UperOtherBgmsFromRecom(c context.Context, upmid, sid int64) (res []*mMdl.Music, extSidMap map[int]int64) {
	extSidMap = make(map[int]int64)
	idx := int(1)
	for _, mscs := range s.MscWithTypes {
		for _, msc := range mscs {
			if msc.SID != sid &&
				msc.UpMID == upmid &&
				len(extSidMap) < 3 {
				res = append(res, msc)
				extSidMap[idx] = msc.SID
				idx++
			}
		}
	}
	return
}

// BgmView fn
func (s *Service) BgmView(c context.Context, sid int64) (ret *mMdl.Music) {
	for _, msc := range s.AllMsc {
		if msc.ID == sid {
			ret = msc
			break
		}
	}
	return
}

// PreByFrom fn
func (s *Service) PreByFrom(c context.Context, from int) (types []*mMdl.Category) {
	if from == 1 {
		sort.Slice(s.Types, func(i, j int) bool {
			return s.Types[i].CameraIndex <= s.Types[j].CameraIndex
		})
	} else {
		sort.Slice(s.Types, func(i, j int) bool {
			return s.Types[i].Index <= s.Types[j].Index
		})
	}
	types = s.Types
	return
}

// BgmList fn
func (s *Service) BgmList(c context.Context, tid int) (ret []*mMdl.Music) {
	if len(s.MscWithTypes) > 0 {
		if mics, ok := s.MscWithTypes[tid]; ok {
			ret = mics
		}
	}
	return
}

// getUpNames fn
func (s *Service) getUpNames(c context.Context, mmap map[int64]*mMdl.Music) (ret map[int64]string, err error) {
	var (
		minfos map[int64]*accMdl.Info
		mids   []int64
	)
	ret = make(map[int64]string)
	for _, msc := range mmap {
		mids = append(mids, msc.UpMID)
	}
	if len(mids) > 0 {
		minfos, err = s.acc.Infos(c, mids, "localhost")
		if err != nil {
			log.Info("minfos err mids (%+v)|err(%+v)", mids, err)
			return
		}
		for _, info := range minfos {
			ret[info.Mid] = info.Name
		}
	}
	return
}

// Cooperate fn
func (s *Service) Cooperate(c context.Context, id, mid int64) (res *mMdl.Cooperate) {
	_, white := s.stickerUps[mid]
	for _, v := range s.Cooperates {
		if v.White == 1 && !white {
			return
		}
		if v.ID == id {
			return v
		}
	}
	return
}

// Material fn
func (s *Service) Material(c context.Context, id int64, tp int8, mid int64) (res interface{}) {
	if _, ok := mMdl.ViewTpMap[tp]; !ok {
		return
	}
	_, white := s.stickerUps[mid]
	switch tp {
	case appMdl.TypeSubtitle:
		for _, v := range s.Subtitles {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	case appMdl.TypeFont:
		for _, v := range s.Fonts {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	case appMdl.TypeFilter:
		for _, v := range s.Filters {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	case appMdl.TypeSticker:
		for _, v := range s.Stickers {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	case appMdl.TypeVideoupSticker:
		for _, v := range s.Vstickers {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	case appMdl.TypeTransition:
		for _, v := range s.Transitions {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	case appMdl.TypeCooperate:
		for _, v := range s.Cooperates {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	case appMdl.TypeTheme:
		for _, v := range s.Themes {
			if v.White == 1 && !white {
				return
			}
			if v.ID == id {
				return v
			}
		}
	}
	return
}

// MaterialPre fn
func (s *Service) MaterialPre(c context.Context, mid int64, platStr string, build int) (res map[string]interface{}) {
	var (
		hotwords             = []*mMdl.Hotword{}
		stickers             = []*mMdl.Sticker{}
		vstickers            = []*mMdl.VSticker{}
		trans                = []*mMdl.Transition{}
		cooperates           = []*mMdl.Cooperate{}
		themes               = []*mMdl.Theme{}
		intros               = []*mMdl.Intro{}
		intro                = &mMdl.Intro{}
		subs                 = []*mMdl.Subtitle{}
		fonts                = []*mMdl.Font{}
		filters              = []*mMdl.Filter{}
		filterWithCategory   = make([]*mMdl.FilterCategory, 0)
		vstickerWithCategory = make([]*mMdl.VstickerCategory, 0)
		white                bool
	)
	_, white = s.stickerUps[mid]
	for _, v := range s.Subtitles {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		subs = append(subs, v)
	}
	for _, v := range s.Fonts {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		fonts = append(fonts, v)
	}
	for _, v := range s.Filters {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		filters = append(filters, v)
	}
	for _, v := range s.Hotwords {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		hotwords = append(hotwords, v)
	}
	for _, v := range s.Stickers {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		stickers = append(stickers, v)
	}
	for _, v := range s.Intros {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		intros = append(intros, v)
	}
	for _, v := range s.Vstickers {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		vstickers = append(vstickers, v)
	}
	for _, v := range s.Transitions {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		trans = append(trans, v)
	}
	for _, v := range s.Themes {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		themes = append(themes, v)
	}
	for _, v := range s.Cooperates {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		cooperates = append(cooperates, v)
	}
	for _, fcategory := range s.FilterWithCategory {
		fcategoryChilds := make([]*mMdl.Filter, 0)
		for _, v := range fcategory.Children {
			if !v.AllowMaterial(v.Material, platStr, build, white) {
				continue
			}
			fcategoryChilds = append(fcategoryChilds, v)
		}
		if len(fcategoryChilds) > 0 {
			filterWithCategory = append(filterWithCategory, &mMdl.FilterCategory{
				ID:       fcategory.ID,
				Name:     fcategory.Name,
				Rank:     fcategory.Rank,
				Tp:       fcategory.Tp,
				Children: fcategoryChilds,
			})
		}
	}
	for _, vscategory := range s.VstickerWithCategory {
		vscategoryChilds := make([]*mMdl.VSticker, 0)
		for _, v := range vscategory.Children {
			if !v.AllowMaterial(v.Material, platStr, build, white) {
				continue
			}
			vscategoryChilds = append(vscategoryChilds, v)
		}
		if len(vscategoryChilds) > 0 {
			vstickerWithCategory = append(vstickerWithCategory, &mMdl.VstickerCategory{
				ID:       vscategory.ID,
				Name:     vscategory.Name,
				Rank:     vscategory.Rank,
				Tp:       vscategory.Tp,
				Children: vscategoryChilds,
			})
		}
	}
	if len(intros) > 0 {
		intro = intros[0]
	}
	res = map[string]interface{}{
		"filter":                 filters,
		"filter_with_category":   filterWithCategory,
		"hotword":                hotwords,
		"sticker":                stickers,
		"intro":                  intro,
		"subtitle":               subs,
		"font":                   fonts,
		"trans":                  trans,
		"themes":                 themes,
		"cooperates":             cooperates,
		"videoup_sticker":        vstickers,
		"vsticker_with_category": vstickerWithCategory,
		"latests":                s.genLatests(intros, filters, subs, fonts, hotwords, stickers, vstickers, trans, cooperates, themes),
	}
	return
}

// begin from app 536
func (s *Service) genLatests(intros []*mMdl.Intro, filters []*mMdl.Filter, subs []*mMdl.Subtitle, fonts []*mMdl.Font, hotwords []*mMdl.Hotword, stickers []*mMdl.Sticker, vstickers []*mMdl.VSticker, trans []*mMdl.Transition, coos []*mMdl.Cooperate, themes []*mMdl.Theme) (latests map[int8]interface{}) {
	latests = make(map[int8]interface{})
	latests[appMdl.TypeBGM] = s.LatestBgm
	if len(intros) > 0 {
		sort.Slice(intros, func(i, j int) bool {
			return intros[i].MTime >= intros[j].MTime
		})
		latests[appMdl.TypeIntro] = intros[0]
	}
	if len(themes) > 0 {
		cthemes := make([]*mMdl.Theme, len(themes))
		copy(cthemes, themes)
		sort.Slice(cthemes, func(i, j int) bool {
			return cthemes[i].MTime >= cthemes[j].MTime
		})
		latests[appMdl.TypeTheme] = cthemes[0]
	}
	if len(coos) > 0 {
		ccoos := make([]*mMdl.Cooperate, len(coos))
		copy(ccoos, coos)
		sort.Slice(ccoos, func(i, j int) bool {
			return ccoos[i].MTime >= ccoos[j].MTime
		})
		latests[appMdl.TypeCooperate] = ccoos[0]
	}
	if len(trans) > 0 {
		ctrans := make([]*mMdl.Transition, len(trans))
		copy(ctrans, trans)
		sort.Slice(ctrans, func(i, j int) bool {
			return ctrans[i].MTime >= ctrans[j].MTime
		})
		latests[appMdl.TypeTransition] = ctrans[0]
	}
	if len(vstickers) > 0 {
		cvstickers := make([]*mMdl.VSticker, len(vstickers))
		copy(cvstickers, vstickers)
		sort.Slice(cvstickers, func(i, j int) bool {
			return cvstickers[i].MTime >= cvstickers[j].MTime
		})
		latests[appMdl.TypeVideoupSticker] = cvstickers[0]
	}
	if len(stickers) > 0 {
		cstickers := make([]*mMdl.Sticker, len(stickers))
		copy(cstickers, stickers)
		sort.Slice(cstickers, func(i, j int) bool {
			return cstickers[i].MTime >= cstickers[j].MTime
		})
		latests[appMdl.TypeSticker] = cstickers[0]
	}
	if len(hotwords) > 0 {
		chotwords := make([]*mMdl.Hotword, len(hotwords))
		copy(chotwords, hotwords)
		sort.Slice(chotwords, func(i, j int) bool {
			return chotwords[i].MTime >= chotwords[j].MTime
		})
		latests[appMdl.TypeHotWord] = chotwords[0]
	}
	if len(fonts) > 0 {
		cfonts := make([]*mMdl.Font, len(fonts))
		copy(cfonts, fonts)
		sort.Slice(cfonts, func(i, j int) bool {
			return cfonts[i].MTime >= cfonts[j].MTime
		})
		latests[appMdl.TypeFont] = cfonts[0]
	}
	if len(filters) > 0 {
		cfilters := make([]*mMdl.Filter, len(filters))
		copy(cfilters, filters)
		sort.Slice(cfilters, func(i, j int) bool {
			return cfilters[i].MTime >= cfilters[j].MTime
		})
		latests[appMdl.TypeFilter] = cfilters[0]
	}
	if len(subs) > 0 {
		csubs := make([]*mMdl.Subtitle, len(subs))
		copy(csubs, subs)
		sort.Slice(csubs, func(i, j int) bool {
			return csubs[i].MTime >= csubs[j].MTime
		})
		latests[appMdl.TypeSubtitle] = csubs[0]
	}
	return
}

// AddBgmFeedBack send to log service
func (s *Service) AddBgmFeedBack(c context.Context, name, musicians, platform string, mid int64) (err error) {
	uInfo := &report.UserInfo{
		Mid:      mid,
		Platform: platform,
		Business: 82, // app投稿业务
		Type:     1,  //投稿bgm反馈
		Oid:      mid,
		Action:   "add_bgm_feedback",
		Ctime:    time.Now(),
		IP:       metadata.String(c, metadata.RemoteIP),
		Index:    []interface{}{name, musicians, mid},
	}
	uInfo.Content = map[string]interface{}{
		"name":      name,
		"musicians": musicians,
		"mid":       mid,
	}
	report.User(uInfo)
	log.Info("send AddBgmFeedBack Log data(%+v)", uInfo)
	return
}

// BgmSearch fn
func (s *Service) BgmSearch(c context.Context, kw string, mid int64, pn, ps int) (res *sMdl.BgmSearchRes) {
	res = &sMdl.BgmSearchRes{
		Bgms: make([]*mMdl.Music, 0),
		Pager: &sMdl.Pager{
			Num:   pn,
			Size:  ps,
			Total: 0,
		},
	}
	if len(kw) == 0 {
		return
	}
	var (
		err    error
		resIDS = make([]int64, 0)
		mids   = make([]int64, 0)
		pager  *sMdl.Pager
		minfos map[int64]*accMdl.Info
	)
	if resIDS, pager, err = s.music.SearchBgmSIDs(c, kw, pn, ps); err != nil {
		log.Error("s.music.SearchBgmSIDs kw(%s)|pn(%d)|ps(%d) error(%v)", kw, pn, ps, err)
		return
	}
	for _, sid := range resIDS {
		if msc, ok := s.AllMsc[sid]; ok {
			res.Bgms = append(res.Bgms, msc)
			mids = append(mids, msc.UpMID)
		}
	}
	if len(mids) > 0 {
		minfos, err = s.acc.Infos(c, mids, "localhost")
		if err != nil {
			log.Info("minfos err mids (%+v)|err(%+v)", mids, err)
			return
		}
		for _, v := range res.Bgms {
			if up, ok := minfos[v.UpMID]; ok {
				v.Musicians = up.Name
			}
		}
	}
	if pager != nil {
		res.Pager.Total = pager.Total
	}
	return
}

func (s *Service) getFilterAndItsCategory(ctx context.Context) {
	var (
		err            error
		filters        []*mMdl.Filter
		filterMap      map[int64]*mMdl.Filter
		filterCategory = make([]*mMdl.FilterCategory, 0)
		filterBinds    []*mMdl.MaterialBind
	)
	c := context.Background()
	if filters, filterMap, err = s.material.Filters(c); err != nil {
		log.Error("s.material.Filters err(%+v)", err)
		return
	}
	s.Filters = filters
	if filterBinds, err = s.material.CategoryBind(c, appMdl.TypeFilter); err != nil {
		log.Error("s.material.CategoryBind err(%+v)", err)
		return
	}
	mapsByCID := make(map[int64]*mMdl.FilterCategory)
	for _, bind := range filterBinds {
		oldF := filterMap[bind.MID]
		if oldF == nil {
			continue
		}
		newF := &mMdl.Filter{
			ID:          oldF.ID,
			Name:        oldF.Name,
			Cover:       oldF.Cover,
			DownloadURL: oldF.DownloadURL,
			Rank:        bind.BRank,
			Extra:       oldF.Extra,
			Material:    oldF.Material,
			MTime:       oldF.MTime,
			New:         oldF.New,
			Tags:        oldF.Tags,
			FilterType:  oldF.FilterType,
		}
		if _, ok := mapsByCID[bind.CID]; !ok {
			mapsByCID[bind.CID] = &mMdl.FilterCategory{
				ID:       bind.CID,
				Name:     bind.CName,
				Rank:     bind.CRank,
				Tp:       bind.Tp,
				New:      bind.New,
				Children: []*mMdl.Filter{newF},
			}
		} else {
			mapsByCID[bind.CID].Children = append(mapsByCID[bind.CID].Children, newF)
		}
	}
	for _, v := range mapsByCID {
		sort.Slice(v.Children, func(i, j int) bool {
			return v.Children[i].Rank <= v.Children[j].Rank
		})
		filterCategory = append(filterCategory, v)
	}
	sort.Slice(filterCategory, func(i, j int) bool {
		return filterCategory[i].Rank <= filterCategory[j].Rank
	})
	s.FilterWithCategory = filterCategory
}

func (s *Service) getCooperates(ctx context.Context) {
	var (
		err     error
		coos    []*mMdl.Cooperate
		daids   []int64
		darcMap map[int64]*api.Arc
	)
	c := context.Background()
	if coos, daids, err = s.material.Cooperates(c); err != nil {
		log.Error("s.material.Cooperates err(%+v)", err)
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	if len(daids) > 0 {
		if darcMap, err = s.archive.Archives(c, daids, ip); err != nil {
			log.Error("s.archive.Archives err(%+v)", err)
			return
		}
	}
	for _, v := range coos {
		mtime := v.MTime.Time()
		hot1 := 20000
		hot2 := mtime.Hour()*100 + mtime.Minute()
		hot3 := 567 * v.ArcCnt
		v.HotVal = hot1 + hot2 + hot3
		if arc, ok := darcMap[v.DemoAID]; ok {
			v.Cover = arc.Pic
		}
	}
	s.Cooperates = coos
}

func (s *Service) getVStickerAndItsCategory(ctx context.Context) {
	var (
		err              error
		vstickers        []*mMdl.VSticker
		vstickersMap     map[int64]*mMdl.VSticker
		vstickerCategory = make([]*mMdl.VstickerCategory, 0)
		vstickerBinds    []*mMdl.MaterialBind
	)
	c := context.Background()
	if vstickers, vstickersMap, err = s.material.Vstickers(c); err != nil {
		log.Error("s.material.vstickers err(%+v)", err)
		return
	}
	s.Vstickers = vstickers
	if vstickerBinds, err = s.material.CategoryBind(c, appMdl.TypeVideoupSticker); err != nil {
		log.Error("s.material.CategoryBind err(%+v)", err)
		return
	}
	mapsByCID := make(map[int64]*mMdl.VstickerCategory)
	for _, bind := range vstickerBinds {
		oldF := vstickersMap[bind.MID]
		if oldF == nil {
			continue
		}
		newF := &mMdl.VSticker{
			ID:          oldF.ID,
			Name:        oldF.Name,
			Cover:       oldF.Cover,
			DownloadURL: oldF.DownloadURL,
			Rank:        bind.BRank,
			Extra:       oldF.Extra,
			Material:    oldF.Material,
		}
		if _, ok := mapsByCID[bind.CID]; !ok {
			mapsByCID[bind.CID] = &mMdl.VstickerCategory{
				ID:       bind.CID,
				Name:     bind.CName,
				Rank:     bind.CRank,
				Tp:       bind.Tp,
				New:      bind.New,
				Children: []*mMdl.VSticker{newF},
			}
		} else {
			mapsByCID[bind.CID].Children = append(mapsByCID[bind.CID].Children, newF)
		}
	}
	for _, v := range mapsByCID {
		sort.Slice(v.Children, func(i, j int) bool {
			return v.Children[i].Rank <= v.Children[j].Rank
		})
		vstickerCategory = append(vstickerCategory, v)
	}
	sort.Slice(vstickerCategory, func(i, j int) bool {
		return vstickerCategory[i].Rank <= vstickerCategory[j].Rank
	})
	s.VstickerWithCategory = vstickerCategory
}

// CooperatePre fn
func (s *Service) CooperatePre(c context.Context, mid int64, platStr string, build int) (cooperates []*mMdl.Cooperate) {
	var (
		white bool
	)
	cooperates = make([]*mMdl.Cooperate, 0)
	_, white = s.stickerUps[mid]
	for _, v := range s.Cooperates {
		if !v.AllowMaterial(v.Material, platStr, build, white) {
			continue
		}
		cooperates = append(cooperates, v)
	}
	for _, x := range cooperates {
		if x.MissionID > 0 {
			if mis, ok := s.p.ActMapCache[x.MissionID]; ok {
				x.Mission = mis
			}
		}
	}
	return
}
