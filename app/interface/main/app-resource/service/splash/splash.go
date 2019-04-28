package splash

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/app-resource/conf"
	addao "go-common/app/interface/main/app-resource/dao/ad"
	locdao "go-common/app/interface/main/app-resource/dao/location"
	spdao "go-common/app/interface/main/app-resource/dao/splash"
	"go-common/app/interface/main/app-resource/model"
	"go-common/app/interface/main/app-resource/model/splash"
	locmdl "go-common/app/service/main/location/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	farm "github.com/dgryski/go-farm"
)

const (
	_birthType     = 2
	_vipType       = 4
	_defaultType   = 0
	_initVersion   = "splash_version"
	_initSplashKey = "splash_key_%d_%d_%d"
)

var (
	_emptySplashs = []*splash.Splash{}
)

// Service is splash service.
type Service struct {
	dao *spdao.Dao
	ad  *addao.Dao
	loc *locdao.Dao
	// tick
	tick time.Duration
	// splash duration
	splashTick string
	// screen
	andScreen map[int8]map[float64][2]int
	iosScreen map[int8]map[float64][2]int
	// cache
	cache        map[string][]*splash.Splash
	defaultCache map[string][]*splash.Splash
	birthCache   map[string][]*splash.Splash
	vipCache     map[string][]*splash.Splash
	// splash random
	splashRandomIds map[int8]map[int64]struct{}
}

// New new a splash service.
func New(c *conf.Config) *Service {
	s := &Service{
		dao: spdao.New(c),
		ad:  addao.New(c),
		loc: locdao.New(c),
		// tick
		tick: time.Duration(c.Tick),
		// splash duration
		splashTick: c.Duration.Splash,
		// screen
		andScreen: map[int8]map[float64][2]int{},
		iosScreen: map[int8]map[float64][2]int{},
		// splash cache
		cache:        map[string][]*splash.Splash{},
		defaultCache: map[string][]*splash.Splash{},
		birthCache:   map[string][]*splash.Splash{},
		vipCache:     map[string][]*splash.Splash{},
		// splash random
		splashRandomIds: map[int8]map[int64]struct{}{},
	}
	s.load()
	s.loadBirth()
	s.loadSplashRandomIds(c)
	go s.loadproc()
	return s
}

// Display dispaly data.
func (s *Service) Display(c context.Context, plat int8, w, h, build int, channel, ver string, now time.Time) (res []*splash.Splash, version string, err error) {
	// get from cache
	res, version, err = s.getCache(c, plat, w, h, build, channel, ver, now)
	return
}

func (s *Service) Birthday(c context.Context, plat int8, w, h int, birth string) (res *splash.Splash, err error) {
	// get from cache
	res, err = s.getBirthCache(c, plat, w, h, birth)
	return
}

// AdList ad splash list
func (s *Service) AdList(c context.Context, plat int8, mobiApp, device, buvid, birth, adExtra string, height, width, build int, mid int64) (res *splash.CmSplash, err error) {
	var (
		list   []*splash.List
		show   []*splash.Show
		config *splash.CmConfig
	)
	if ok := model.IsOverseas(plat); ok {
		err = ecode.NotModified
		return
	}
	g, ctx := errgroup.WithContext(c)
	g.Go(func() error {
		var e error
		if list, config, e = s.ad.SplashList(ctx, mobiApp, device, buvid, birth, adExtra, height, width, build, mid); e != nil {
			log.Error("cm s.ad.SplashList error(%v)", e)
			return e
		}
		return nil
	})
	g.Go(func() error {
		var e error
		if show, e = s.ad.SplashShow(ctx, mobiApp, device, buvid, birth, adExtra, height, width, build, mid); e != nil {
			log.Error("cm s.ad.SplashShow error(%v)", e)
			return e
		}
		return nil
	})
	if err = g.Wait(); err != nil {
		log.Error("cm splash errgroup.WithContext error(%v)", err)
		return
	}
	res = &splash.CmSplash{
		CmConfig: config,
		List:     list,
		Show:     show,
	}
	return
}

// getCache cache display data.
func (s *Service) getCache(c context.Context, plat int8, w, h, build int, channel, ver string, now time.Time) (res []*splash.Splash, version string, err error) {
	var (
		ip     = metadata.String(c, metadata.RemoteIP)
		screen map[int8]map[float64][2]int
	)
	if model.IsIOS(plat) {
		screen = s.iosScreen
	} else if model.IsAndroid(plat) {
		screen = s.andScreen
	}
	// TODO fate go start
	var (
		fgId     int64
		oldIdStr string
		fgids    = s.splashRandomIds[plat]
		pids     []string
		auths    map[string]*locmdl.Auth
	)
	vers := strings.Split(ver, strconv.Itoa(now.Year()))
	if len(vers) > 1 {
		ver = vers[0]
		oldIdStr = vers[1]
	}
	for id, _ := range fgids {
		fgId = id
		idStr := strconv.FormatInt(fgId, 10)
		if oldIdStr != idStr {
			break
		}
	}
	for tSplash, tScreen := range screen {
		var ss []*splash.Splash
		if ss = s.cache[fmt.Sprintf(_initSplashKey, plat, w, h)]; len(ss) == 0 {
			wh := s.similarScreen(plat, w, h, tScreen)
			width := wh[0]
			height := wh[1]
			ss = s.cache[fmt.Sprintf(_initSplashKey, plat, width, height)]
		}
		for _, splash := range ss {
			if splash.Type != tSplash {
				continue
			}
			if splash.Area != "" {
				pids = append(pids, splash.Area)
			}
		}
	}
	if len(pids) > 0 {
		auths, _ = s.loc.AuthPIDs(c, strings.Join(pids, ","), ip)
	}
	// TODO fate go end
	for tSplash, tScreen := range screen {
		var (
			ss []*splash.Splash
			// advance time
			advance, _ = time.ParseDuration(s.splashTick)
		)
		if ss = s.cache[fmt.Sprintf(_initSplashKey, plat, w, h)]; len(ss) == 0 {
			wh := s.similarScreen(plat, w, h, tScreen)
			width := wh[0]
			height := wh[1]
			ss = s.cache[fmt.Sprintf(_initSplashKey, plat, width, height)]
		}
		for _, splash := range ss {
			if splash.Type != tSplash {
				continue
			}
			// gt splash start time
			if splash.NoPreview == 1 {
				if h1 := now.Add(advance); int64(splash.Start) > h1.Unix() {
					continue
				}
			}
			// TODO fate go start
			if fgids != nil && splash.ID != fgId {
				if _, ok := fgids[splash.ID]; ok {
					continue
				}
			}
			// TODO fate go end
			if model.InvalidBuild(build, splash.Build, splash.Condition) {
				continue
			}
			if splash.Area != "" {
				if auth, ok := auths[splash.Area]; ok && auth.Play == locmdl.Forbidden {
					log.Warn("s.invalid area(%v) ip(%v) error(%v)", splash.Area, ip, err)
					continue
				}
			}
			res = append(res, splash)
		}
	}
	if vSplash := s.getVipCache(plat, w, h, screen, now); vSplash != nil {
		res = append(res, vSplash)
	}
	if dSplash := s.getDefaultCache(plat, w, h, screen, now); dSplash != nil {
		res = append(res, dSplash)
	}
	if len(res) == 0 {
		res = _emptySplashs
	}
	if version = s.hash(res); version == ver {
		err = ecode.NotModified
		res = nil
	}
	version = version + strconv.Itoa(now.Year()) + strconv.FormatInt(fgId, 10)
	return
}

// getBirthCache get birthday splash.
func (s *Service) getBirthCache(c context.Context, plat int8, w, h int, birth string) (res *splash.Splash, err error) {
	var (
		screen map[int8]map[float64][2]int
		wh     [2]int
	)
	if model.IsIOS(plat) {
		screen = s.iosScreen
	} else if model.IsAndroid(plat) {
		screen = s.andScreen
	}
	if v, ok := screen[_birthType]; !ok {
		return
	} else {
		wh = s.similarScreen(plat, w, h, v)
		w = wh[0]
		h = wh[1]
	}
	sps := s.birthCache[fmt.Sprintf(_initSplashKey, plat, w, h)]
	for _, sp := range sps {
		if sp.BirthStartMonth == "12" && sp.BirthEndMonth == "01" {
			if (sp.BirthStart <= birth && "1231" >= birth) || ("0101" <= birth && sp.BirthEnd >= birth) {
				res = sp
				return
			}
		}
		if sp.BirthStart <= birth && sp.BirthEnd >= birth {
			res = sp
			return
		}
	}
	err = ecode.NothingFound
	return
}

// getVipCache
func (s *Service) getVipCache(plat int8, w, h int, screen map[int8]map[float64][2]int, now time.Time) (res *splash.Splash) {
	var (
		ss []*splash.Splash
	)
	if v, ok := screen[_vipType]; !ok {
		return
	} else if ss = s.vipCache[fmt.Sprintf(_initSplashKey, plat, w, h)]; len(ss) == 0 {
		wh := s.similarScreen(plat, w, h, v)
		width := wh[0]
		height := wh[1]
		ss = s.vipCache[fmt.Sprintf(_initSplashKey, plat, width, height)]
	}
	if len(ss) == 0 {
		return
	}
	res = ss[(now.Day() % len(ss))]
	return
}

// getDefaultCache
func (s *Service) getDefaultCache(plat int8, w, h int, screen map[int8]map[float64][2]int, now time.Time) (res *splash.Splash) {
	var (
		ss []*splash.Splash
	)
	if v, ok := screen[_defaultType]; !ok {
		return
	} else if ss = s.defaultCache[fmt.Sprintf(_initSplashKey, plat, w, h)]; len(ss) == 0 {
		wh := s.similarScreen(plat, w, h, v)
		width := wh[0]
		height := wh[1]
		ss = s.defaultCache[fmt.Sprintf(_initSplashKey, plat, width, height)]
	}
	if len(ss) == 0 {
		return
	}
	res = ss[(now.Day() % len(ss))]
	return
}

func (s *Service) hash(v []*splash.Splash) string {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Error("json.Marshal error(%v)", err)
		return _initVersion
	}
	return strconv.FormatUint(farm.Hash64(bs), 10)
}

// cacheproc load splash into cache.
func (s *Service) load() {
	res, err := s.dao.ActiveAll(context.TODO())
	if err != nil {
		log.Error("s.dao.GetActiveAll() error(%v)", err)
		return
	}
	var (
		tmp        []*splash.Splash
		tmpdefault []*splash.Splash
	)
	for _, r := range res {
		if r.Type == _defaultType {
			tmpdefault = append(tmpdefault, r)
		} else {
			tmp = append(tmp, r)
		}
	}
	s.cache = s.dealCache(tmp)
	log.Info("splash cacheproc success")
	s.defaultCache = s.dealCache(tmpdefault)
	log.Info("splash default cacheproc tmpdefault")
	resVip, err := s.dao.ActiveVip(context.TODO())
	if err != nil {
		log.Error("s.dao.ActiveVip() error(%v)", err)
		return
	}
	s.vipCache = s.dealCache(resVip)
	log.Info("splash Vip cacheproc success")
}

// loadBirth load birthday splash.
func (s *Service) loadBirth() {
	res, err := s.dao.ActiveBirth(context.TODO())
	if err != nil {
		log.Error("s.dao.ActiveBirthday() error(%v)", err)
		return
	}
	s.birthCache = s.dealCache(res)
	log.Info("splash Birthday cacheproc success")
}

// dealCache
func (s *Service) dealCache(sps []*splash.Splash) (res map[string][]*splash.Splash) {
	res = map[string][]*splash.Splash{}
	tmpand := map[int8]map[float64][2]int{}
	tmpios := map[int8]map[float64][2]int{}
	for plat, v := range s.andScreen {
		for r, value := range v {
			if _, ok := tmpand[plat]; ok {
				tmpand[plat][r] = value
			} else {
				tmpand[plat] = map[float64][2]int{
					r: value,
				}
			}
		}
	}
	for plat, v := range s.iosScreen {
		for r, value := range v {
			if _, ok := tmpios[plat]; ok {
				tmpios[plat][r] = value
			} else {
				tmpios[plat] = map[float64][2]int{
					r: value,
				}
			}
		}
	}
	for _, v := range sps {
		v.URI = model.FillURI(v.Goto, v.Param, nil)
		key := fmt.Sprintf(_initSplashKey, v.Plat, v.Width, v.Height)
		res[key] = append(res[key], v)
		// generate screen
		if model.IsAndroid(v.Plat) {
			if _, ok := tmpand[v.Type]; ok {
				tmpand[v.Type][splash.Ratio(v.Width, v.Height)] = [2]int{v.Width, v.Height}
			} else {
				tmpand[v.Type] = map[float64][2]int{
					splash.Ratio(v.Width, v.Height): [2]int{v.Width, v.Height},
				}
			}
		} else if model.IsIOS(v.Plat) {
			if _, ok := tmpios[v.Type]; ok {
				tmpios[v.Type][splash.Ratio(v.Width, v.Height)] = [2]int{v.Width, v.Height}
			} else {
				tmpios[v.Type] = map[float64][2]int{
					splash.Ratio(v.Width, v.Height): [2]int{v.Width, v.Height},
				}
			}
		}
	}
	s.andScreen = tmpand
	s.iosScreen = tmpios
	return
}

func (s *Service) loadSplashRandomIds(c *conf.Config) {
	splashIds := map[int8]map[int64]struct{}{}
	for k, v := range c.Splash.Random {
		key := model.Plat(k, "")
		splashIds[key] = map[int64]struct{}{}
		for _, idStr := range v {
			idInt, _ := strconv.ParseInt(idStr, 10, 64)
			splashIds[key][idInt] = struct{}{}
		}
	}
	s.splashRandomIds = splashIds
	log.Info("splash Random cache success")
}

// loadproc load process.
func (s *Service) loadproc() {
	for {
		time.Sleep(s.tick)
		s.load()
		s.loadBirth()
	}
}

// similarScreen android screnn size
func (s *Service) similarScreen(plat int8, w, h int, screen map[float64][2]int) (wh [2]int) {
	if model.IsIOS(plat) {
		switch {
		case w == 750:
			h = 1334
		case w == 640 && h > 960:
			h = 1136
		case w == 640:
			h = 960
		case w == 2732:
			h = 2048
		case w == 2048:
			h = 1536
		case w == 1024:
			h = 768
		case w == 1242:
			h = 2208
		case w == 1496 || w == 1536:
			w = 2048
			h = 1536
		case w == 748 || w == 768:
			w = 1024
			h = 768
		}
	}
	min := float64(1<<64 - 1)
	for r, s := range screen {
		if s[0] == w && s[1] == h {
			wh = s
			return
		}
		abs := math.Abs(splash.Ratio(w, h) - r)
		if abs < min {
			min = abs
			wh = s
		}
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
