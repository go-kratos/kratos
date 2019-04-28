package space

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go-common/app/interface/main/app-interface/model/space"
	accmdl "go-common/app/service/main/account/model"
	relmdl "go-common/app/service/main/relation/model"
	"go-common/library/log"

	"go-common/library/sync/errgroup.v2"
)

const (
	_initSidebarKey = "sidebar_%d_%d"
	_selfCenter     = 6
	_mySerive       = 7
	_creative       = 11
	_ipadSelfCenter = 12
	_ipadCreative   = 13
)

var (
	iphoneMenu = map[int]string{
		_creative:   "创作中心",
		_selfCenter: "个人中心",
		_mySerive:   "我的服务",
	}
	iphoneNormalMenu = []int{_creative, _selfCenter, _mySerive}
	iphoneFilterMenu = []int{_selfCenter, _mySerive}
	ipadNormalMenu   = []int{_ipadCreative, _ipadSelfCenter}
	ipadFilterMenu   = []int{_ipadSelfCenter}
)

// Mine mine center for iphone/android
func (s *Service) Mine(c context.Context, mid int64, platform, filtered string, build int, plat int8) (mine *space.Mine, err error) {
	var whiteMap, rdMap map[int64]bool
	mine = new(space.Mine)
	mine.Official.Type = -1
	if mid > 0 {
		if mine, whiteMap, rdMap, err = s.userInfo(c, mid, platform, plat); err != nil {
			return
		}
	}
	if platform == "ios" {
		mine.Sections = s.sections(c, whiteMap, rdMap, mid, build, filtered == "1", plat)
	}
	return
}

// MineIpad mine center for ipad
func (s *Service) MineIpad(c context.Context, mid int64, platform, filtered string, build int, plat int8) (mine *space.Mine, err error) {
	var whiteMap, rdMap map[int64]bool
	mine = new(space.Mine)
	mine.Official.Type = -1
	if mid > 0 {
		if mine, whiteMap, rdMap, err = s.userInfo(c, mid, platform, plat); err != nil {
			return
		}
	}
	mine.IpadSections, mine.IpadUpperSections = s.ipadSections(c, whiteMap, rdMap, mid, build, filtered == "1", plat)
	return
}

func (s *Service) userInfo(c context.Context, mid int64, platform string, plat int8) (mine *space.Mine, whiteMap, rdMap map[int64]bool, err error) {
	mine = new(space.Mine)
	whiteMap = make(map[int64]bool)
	rdMap = make(map[int64]bool)
	eg := errgroup.WithContext(c)
	// account userinfo
	eg.Go(func(ctx context.Context) (err error) {
		var ps *accmdl.ProfileStat
		if ps, err = s.accDao.Profile3(ctx, mid); err != nil {
			log.Error("s.accDao.UserInfo(%d) error(%v)", mid, err)
			return
		}
		if ps.Silence == 1 {
			if mine.EndTime, err = s.accDao.BlockTime(ctx, mid); err != nil {
				log.Error("%+v", err)
				err = nil
			}
		}
		mine.Silence = ps.Silence
		mine.Mid = ps.Mid
		mine.Name = ps.Name
		mine.Face = ps.Face
		mine.Coin = ps.Coins
		if ps.Pendant.Image != "" {
			mine.Pendant = &space.Pendant{Image: ps.Pendant.Image}
		}
		switch ps.Sex {
		case "男":
			mine.Sex = 1
		case "女":
			mine.Sex = 2
		default:
			mine.Sex = 0
		}
		mine.Rank = ps.Rank
		mine.Level = ps.Level
		if ps.Vip.Status == 1 {
			mine.VipType = ps.Vip.Type
		}
		if ps.Official.Role == 0 {
			mine.Official.Type = -1
		} else {
			if ps.Official.Role <= 2 {
				mine.Official.Type = 0
			} else {
				mine.Official.Type = 1
			}
			mine.Official.Desc = ps.Official.Title
		}
		return
	})
	// music card
	eg.Go(func(ctx context.Context) (err error) {
		cardm, err := s.audioDao.Card(ctx, mid)
		if err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if card, ok := cardm[mid]; ok && card.Type == 1 && card.Status == 1 {
			mine.AudioType = card.Type
		}
		return
	})
	// following and follower
	eg.Go(func(ctx context.Context) (err error) {
		var stat *relmdl.Stat
		if stat, err = s.relDao.Stat(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		mine.Following = stat.Following
		mine.Follower = stat.Follower
		return
	})
	// new followers
	eg.Go(func(ctx context.Context) (err error) {
		if mine.NewFollowers, err = s.relDao.FollowersUnreadCount(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		return
	})
	// dynamic count
	eg.Go(func(ctx context.Context) (err error) {
		var count int64
		if count, err = s.bplusDao.DynamicCount(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		mine.Dynamic = count
		return
	})
	// bcoin
	eg.Go(func(ctx context.Context) (err error) {
		var bp float64
		if bp, err = s.payDao.UserWalletInfo(ctx, mid, platform); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		mine.BCoin = bp
		return
	})
	// creative
	eg.Go(func(ctx context.Context) (err error) {
		var (
			isUp int
			show int
		)
		if isUp, show, err = s.memberDao.Creative(ctx, mid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		mine.ShowVideoup = show
		mine.ShowCreative = isUp
		return
	})
	if platform == "ios" {
		var mutex sync.Mutex
		//white
		for _, v := range s.white[plat] {
			tmpID := v.ID
			tmpURL := v.URL
			eg.Go(func(ctx context.Context) (err error) {
				ok, err := s.accDao.UserCheck(ctx, mid, tmpURL)
				if err != nil {
					log.Error("s.accDao.UserCheck error(%+v)", err)
					err = nil
					return
				}
				if ok {
					mutex.Lock()
					whiteMap[tmpID] = true
					mutex.Unlock()
				}
				return
			})
		}
		//redDot
		for _, v := range s.redDot[plat] {
			tmpID := v.ID
			tmpURL := v.URL
			eg.Go(func(ctx context.Context) (err error) {
				ok, err := s.accDao.RedDot(ctx, mid, tmpURL)
				if err != nil {
					log.Error("s.accDao.RedDot error(%+v)", err)
					err = nil
					return
				}
				if ok {
					mutex.Lock()
					rdMap[tmpID] = true
					mutex.Unlock()
				}
				return
			})
		}
	}
	// when account info error,return,ingore else error
	err = eg.Wait()
	return
}

func (s *Service) sections(c context.Context, whiteMap, rdMap map[int64]bool, mid int64, build int, filtered bool, plat int8) (sections []*space.Section) {
	menus := iphoneNormalMenu
	if filtered {
		menus = iphoneFilterMenu
	}
	for _, module := range menus {
		key := fmt.Sprintf(_initSidebarKey, plat, module)
		ss, ok := s.sectionCache[key]
		if !ok {
			continue
		}
		var items []*space.SectionItem
		for _, si := range ss {
			if !si.CheckLimit(build) {
				continue
			}
			if si.Item.Name == "离线缓存" && filtered {
				continue
			}
			if si.Item.WhiteURL != "" && !whiteMap[si.Item.ID] {
				continue
			}
			tmpItem := &space.SectionItem{
				Title:     si.Item.Name,
				Icon:      si.Item.Logo,
				NeedLogin: si.Item.NeedLogin,
				URI:       si.Item.Param,
			}
			if si.Item.Red != "" && rdMap[si.Item.ID] {
				tmpItem.RedDot = 1
			}
			items = append(items, tmpItem)
		}
		if len(items) == 0 {
			continue
		}
		sections = append(sections, &space.Section{
			Title: iphoneMenu[module],
			Items: items,
		})
	}
	return
}

func (s *Service) ipadSections(c context.Context, whiteMap, rdMap map[int64]bool, mid int64, build int, filtered bool, plat int8) (ipadSections, ipadUpperSections []*space.SectionItem) {
	menus := ipadNormalMenu
	if filtered {
		menus = ipadFilterMenu
	}
	for _, module := range menus {
		key := fmt.Sprintf(_initSidebarKey, plat, module)
		ss, ok := s.sectionCache[key]
		if !ok {
			continue
		}
		for _, si := range ss {
			if !si.CheckLimit(build) {
				continue
			}
			if si.Item.Name == "离线缓存" && filtered {
				continue
			}
			if si.Item.WhiteURL != "" && !whiteMap[si.Item.ID] {
				continue
			}
			tmpItem := &space.SectionItem{
				Title:     si.Item.Name,
				Icon:      si.Item.Logo,
				NeedLogin: si.Item.NeedLogin,
				URI:       si.Item.Param,
			}
			if si.Item.Red != "" && rdMap[si.Item.ID] {
				tmpItem.RedDot = 1
			}
			if module == _ipadCreative {
				ipadUpperSections = append(ipadUpperSections, tmpItem)
			} else {
				ipadSections = append(ipadSections, tmpItem)
			}
		}
	}
	return
}

// Myinfo simple myinfo
func (s *Service) Myinfo(c context.Context, mid int64) (myinfo *space.Myinfo, err error) {
	var pf *accmdl.ProfileStat
	if pf, err = s.accDao.Profile3(c, mid); err != nil {
		log.Error("%+v", err)
		return
	}
	p, _ := json.Marshal(pf)
	log.Warn("myinfo mid(%d) pf(%s)", mid, p)
	myinfo = new(space.Myinfo)
	if pf.Silence == 1 {
		if myinfo.EndTime, err = s.accDao.BlockTime(c, mid); err != nil {
			log.Error("%+v", err)
			err = nil
		}
	}
	myinfo.Coins = pf.Coins
	myinfo.Sign = pf.Sign
	switch pf.Sex {
	case "男":
		myinfo.Sex = 1
	case "女":
		myinfo.Sex = 2
	default:
		myinfo.Sex = 0
	}
	myinfo.Mid = mid
	myinfo.Birthday = pf.Birthday.Time().Format("2006-01-02")
	myinfo.Name = pf.Name
	myinfo.Face = pf.Face
	myinfo.Rank = pf.Rank
	myinfo.Level = pf.Level
	myinfo.Vip = pf.Vip
	myinfo.Silence = pf.Silence
	myinfo.EmailStatus = pf.EmailStatus
	myinfo.TelStatus = pf.TelStatus
	myinfo.Official = pf.Official
	myinfo.Identification = pf.Identification
	if pf.Pendant.Image != "" {
		myinfo.Pendant = &space.Pendant{Image: pf.Pendant.Image}
	}
	return
}
