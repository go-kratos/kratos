package weeklyhonor

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/up"
	"go-common/app/interface/main/creative/dao/weeklyhonor"
	model "go-common/app/interface/main/creative/model/weeklyhonor"
	"go-common/app/interface/main/creative/service"
	accmdl "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	upmdl "go-common/app/service/main/up/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/sync/pipeline/fanout"
	xtime "go-common/library/time"
)

const (
	layout = "20060102"
)

// Service struct.
type Service struct {
	c      *conf.Config
	honDao *weeklyhonor.Dao
	arc    *archive.Dao
	acc    *account.Dao
	up     *up.Dao
	// cache chan
	cache    *fanout.Fanout
	honorMap map[int][]*model.HonorWord
}

// New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:      c,
		honDao: weeklyhonor.New(c),
		arc:    rpcdaos.Arc,
		acc:    rpcdaos.Acc,
		up:     rpcdaos.Up,
		// cache
		cache: fanout.New("cache"),
	}
	s.honorMap = model.HMap()
	return s
}

// Ping service.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.honDao.Ping(c); err != nil {
		log.Error("s.honor.Ping err(%v)", err)
	}
	return
}

// Close dao.
func (s *Service) Close() {
	s.honDao.Close()
}

// ChangeSubState change subscribe state
func (s *Service) ChangeSubState(c context.Context, mid int64, state uint8) (err error) {
	err = s.honDao.ChangeUpSwitch(c, mid, state)
	if err != nil {
		log.Error("s.honDao.ChangeUpSwitch mid(%d) state(%d) err(%v)", mid, state, err)
	}
	return
}

// WeeklyHonor .
func (s *Service) WeeklyHonor(c context.Context, mid, uid int64, token string) (h *model.Honor, err error) {
	var upMid = mid
	h = &model.Honor{}
	if uid != 0 && uid != upMid {
		if token != s.genToken(uid) {
			return nil, ecode.RequestErr
		}
		upMid = uid
	}
	h.MID = upMid
	var upInfo *upmdl.UpInfo
	upInfo, err = s.up.UpInfo(c, upMid, 0, "")
	if err != nil {
		log.Error("s.up.UpInfo(%d) error(%v)", upMid, err)
		return
	}
	if upInfo.IsAuthor != 1 {
		return nil, ecode.CreativeNotUper
	}
	if upMid == mid {
		h.SubState, err = s.honDao.GetUpSwitch(c, upMid)
		if err != nil {
			log.Error("s.honDao.GetUpSwitch upMid(%d),err(%v)", upMid, err)
		}
		go s.addHonorClickCount(context.Background(), mid)
	}
	var timeUp bool
	h.DateBegin, h.DateEnd, timeUp = honTimeFrame()
	endStr := time.Unix(int64(h.DateEnd), 0).Format(layout)
	h.ShareToken = s.genToken(upMid)
	hs, err := s.honorStat(c, upMid, endStr)
	if hs == nil || err != nil || s.c.HonorDegradeSwitch {
		return s.degrade(c, h)
	}
	hl, err := s.honDao.HonorMC(c, h.MID, endStr)
	if err != nil {
		log.Error("s.honDao.HonorMC() error(%v)", err)
		return s.degrade(c, h)
	}
	if hl != nil {
		h.HID = hl.HID
		h.HonorCount = hl.Count
	} else {
		h, err = s.genHonor(c, h, endStr, timeUp, hs)
		if err != nil {
			return
		}
	}
	if h.HID == 0 {
		return s.degrade(c, h)
	}
	if err = s.rpcFill(c, hs, h); err != nil {
		return
	}
	s.wordFill(hs, h)
	h.RiseStage = s.stars(hs)
	return
}

func (s *Service) addHonorClickCount(c context.Context, mid int64) {
	err := s.honDao.ClickMC(c, mid)
	if err != memcache.ErrNotFound {
		if err != nil {
			log.Error("s.honDao.ClickMC mid(%d) err(%+v)", mid, err)
		}
		return
	}
	err = s.honDao.SetClickMC(c, mid)
	if err != nil {
		log.Error("s.honDao.SetClickMC mid(%d) err(%+v)", mid, err)
		return
	}
	err = s.honDao.UpsertClickCount(c, mid)
	if err != nil {
		log.Error("failed to add honor click count,mid(%d) err(%+v)", mid, err)
	}
}

func (s *Service) honorStat(c context.Context, mid int64, date string) (hs *model.HonorStat, err error) {
	hs, err = s.honDao.StatMC(c, mid, date)
	if hs != nil && err == nil {
		return
	}
	for i := 0; i < 3; i++ {
		hs, err = s.honDao.HonorStat(c, mid, date)
		if err != nil {
			log.Error("s.honDao.HonorStat(%d,%v) error(%v)", mid, date, err)
			continue
		}
		if hs != nil {
			break
		}
	}
	_ = s.cache.Do(c, func(c context.Context) {
		_ = s.honDao.SetStatMC(c, mid, date, hs)
	})
	return
}

func (s *Service) degrade(c context.Context, h *model.Honor) (*model.Honor, error) {
	h.Word = s.honorMap[58][0].Word
	h.Text = s.honorMap[58][0].Text
	h.Priority = s.honorMap[58][0].Priority
	us, err := s.acc.Infos(c, []int64{h.MID}, "")
	if err != nil {
		log.Error("s.acc.Infos(%v) error(%v)", h.MID, err)
		return nil, err
	}
	if u, ok := us[h.MID]; ok {
		h.Uname = u.Name
		h.Face = u.Face
	}
	return h, nil
}

func (s *Service) wordFill(hs *model.HonorStat, h *model.Honor) {
	hw := s.wordFormat(h, hs)
	h.Word = hw.Word
	h.Text = hw.Text
	h.Desc = hw.Desc
	h.Priority = hw.Priority
}

func (s *Service) wordFormat(h *model.Honor, stat *model.HonorStat) *model.HonorWord {
	hid := h.HID
	uname := h.Uname
	hm := s.honorMap
	hw := new(model.HonorWord)
	if _, ok := hm[hid]; !ok {
		return hw
	}
	hwfmt := s.getHWFmt(h)
	hw.Word = hwfmt.Word
	hw.Text = hwfmt.Text
	hw.Desc = hwfmt.Desc
	hw.Priority = hwfmt.Priority
	switch hid {
	case 2, 3, 4, 59:
		hw.Desc = fmt.Sprintf(hwfmt.Desc, stat.Rank0)
	case 5:
		var str string
		if stat.Play/100000000 > 0 {
			str = fmt.Sprintf("%d个亿", stat.Play/100000000)
		} else {
			str = fmt.Sprintf("%d千万", stat.Play/10000000)
		}
		hw.Desc = fmt.Sprintf(hwfmt.Desc, str)
	case 6:
		hw.Text = fmt.Sprintf(hwfmt.Text, uname)
		hw.Desc = fmt.Sprintf(hwfmt.Desc, stat.Fans/1000000)
	case 8:
		_, partion, rank := stat.PartionRank()
		hw.Desc = fmt.Sprintf(hwfmt.Desc, partion, rank)
	case 9:
		num := stat.Play / 100000
		if stat.Play/1000000 > 0 {
			num = stat.Play / 1000000 * 10
		}
		hw.Text = fmt.Sprintf(hwfmt.Text, num*10*2)
		hw.Desc = fmt.Sprintf(hwfmt.Desc, num*10)
	case 10:
		hw.Desc = fmt.Sprintf(hwfmt.Desc, stat.Fans/10000)
	case 12, 30:
		hw.Text = fmt.Sprintf(hwfmt.Text, uname)
	case 17:
		hw.Text = fmt.Sprintf(hwfmt.Text, stat.Play/10000*2)
		hw.Desc = fmt.Sprintf(hwfmt.Desc, stat.Play/10000)
	case 18:
		hw.Desc = fmt.Sprintf(hwfmt.Desc, stat.Fans/10000)
	case 26:
		hw.Desc = fmt.Sprintf(hwfmt.Desc, stat.Play/1000)
	case 27:
		hw.Desc = fmt.Sprintf(hwfmt.Desc, stat.Fans/1000)
	}
	return hw
}

func (s *Service) stars(hs *model.HonorStat) *model.RiseStage {
	_, stars := hs.PrioritySR()
	_, starsR := hs.PriorityR()
	stars = s.outputBig(stars, starsR)
	_, starsA := hs.PriorityA()
	stars = s.outputBig(stars, starsA)
	_, starsB := hs.PriorityB()
	stars = s.outputBig(stars, starsB)
	_, starsC := hs.PriorityC()
	stars = s.outputBig(stars, starsC)
	return stars
}

func (s *Service) outputBig(a, b *model.RiseStage) *model.RiseStage {
	if b.Coin > a.Coin {
		a.Coin = b.Coin
	}
	if b.Fans > a.Fans {
		a.Fans = b.Fans
	}
	if b.Like > a.Like {
		a.Like = b.Like
	}
	if b.Play > a.Play {
		a.Play = b.Play
	}
	if b.Share > a.Share {
		a.Share = b.Share
	}
	return a
}

// honTimeFrame  get honor time frame return start=last week Sun. end= this week Sat.
func honTimeFrame() (start, end xtime.Time, timeUp bool) {
	lastSun := model.LatestSunday()
	// saturday
	endTime := lastSun.AddDate(0, 0, -1)
	timeUp = time.Now().Unix() > lastSun.Unix()+18*3600
	if !timeUp {
		endTime = lastSun.AddDate(0, 0, -8)
	}
	// the week's (end's week-1) sunday
	start = xtime.Time(endTime.AddDate(0, 0, -6).Unix())
	end = xtime.Time(endTime.Unix())
	return
}

// corHW get correspond HonorWord fmt
func (s *Service) getHWFmt(h *model.Honor) (hw *model.HonorWord) {
	hm := s.honorMap
	hws := hm[h.HID]
	if len(hws) == 1 {
		hw = hws[0]
		return
	}
	var (
		l = 0
		r = len(hws) - 1
		m = (l + r) / 2
	)
	if h.DateEnd > hws[r].Start {
		hw = hws[r]
		return
	}
	if h.DateEnd <= hws[l].End {
		hw = hws[l]
		return
	}
	for l < r {
		if h.DateEnd > hws[m].Start && h.DateEnd <= hws[m].End {
			break
		} else if h.DateEnd <= hws[m].Start {
			r = m - 1
		} else {
			l = m + 1
		}
		m = (l + r) / 2
	}
	hw = hws[m]
	return
}

func (s *Service) genHonor(c context.Context, h *model.Honor, endStr string, timeUp bool, hs *model.HonorStat) (*model.Honor, error) {
	hls, err := s.honDao.HonorLogs(c, h.MID)
	end := h.DateEnd
	if err != nil {
		log.Error("s.honDao.HonorLogs(%d) error(%v)", h.MID, err)
		return s.degrade(c, h)
	}
	var lastHid int
	for _, v := range hls {
		if int64(v.MTime) < int64(end) && int64(v.MTime) > int64(end)-7*24*3600 {
			lastHid = v.HID
		}
	}
	h.HID = hs.GenHonor(h.MID, lastHid)
	if h.HID == 0 {
		return s.degrade(c, h)
	}
	var needUpdate bool
	if v, ok := hls[h.HID]; !ok {
		needUpdate = true
		h.HonorCount = 1
	} else {
		h.HonorCount = v.Count
		if int64(v.MTime) < int64(end) {
			needUpdate = true
			h.HonorCount = v.Count + 1
		}
	}
	if !timeUp {
		return h, nil
	}
	if needUpdate {
		_ = s.cache.Do(c, func(c context.Context) {
			if res, err1 := s.honDao.HonorMC(c, h.MID, endStr); err1 != nil || res != nil {
				return
			}
			if err1 := s.honDao.UpsertCount(c, h.MID, h.HID); err1 != nil {
				log.Error("s.honDao.UpsertCount() error(%v)", err1)
				return
			}
			hl := &model.HonorLog{
				MID:   h.MID,
				HID:   h.HID,
				Count: h.HonorCount,
			}
			if err1 := s.honDao.SetHonorMC(c, h.MID, endStr, hl); err1 != nil {
				log.Error("s.honDao.SetHonorMC(%d,%s,%v) error(%v)", h.MID, endStr, hl, err1)
			}
		})
	} else {
		hl := &model.HonorLog{
			MID:   h.MID,
			HID:   h.HID,
			Count: h.HonorCount,
		}
		if err1 := s.honDao.SetHonorMC(c, h.MID, endStr, hl); err1 != nil {
			log.Error("s.honDao.SetHonorMC(%d,%s,%v) error(%v)", h.MID, endStr, hl, err1)
		}
	}
	return h, nil
}

func (s *Service) rpcFill(c context.Context, hs *model.HonorStat, h *model.Honor) error {
	mids := []int64{h.MID}
	h.LoveFans = make([]*accmdl.Info, 0)
	h.PlayFans = make([]*accmdl.Info, 0)
	loveFans := make([]int64, 0)
	playFans := make([]int64, 0)
	aids := make([]int64, 0)
	// users
	if hs.Act1 != 0 {
		loveFans = append(loveFans, int64(hs.Act1))
	}
	if hs.Act2 != 0 {
		loveFans = append(loveFans, int64(hs.Act2))
	}
	if hs.Act3 != 0 {
		loveFans = append(loveFans, int64(hs.Act3))
	}
	mids = append(mids, loveFans...)
	if hs.Dr1 != 0 {
		playFans = append(playFans, int64(hs.Dr1))
	}
	if hs.Dr2 != 0 {
		playFans = append(playFans, int64(hs.Dr2))
	}
	if hs.Dr3 != 0 {
		playFans = append(playFans, int64(hs.Dr3))
	}
	mids = append(mids, playFans...)
	// arcs
	if hs.HottestAvNew != 0 {
		aids = append(aids, int64(hs.HottestAvNew))
	}
	if hs.HottestAvInc != 0 {
		aids = append(aids, int64(hs.HottestAvInc))
	}
	if hs.HottestAvAll != 0 {
		aids = append(aids, int64(hs.HottestAvAll))
	}
	g := new(errgroup.Group)
	var (
		us   map[int64]*accmdl.Info
		arcs map[int64]*api.Arc
	)
	if len(aids) > 0 {
		g.Go(func() (err error) {
			arcs, err = s.arc.Archives(c, aids, "")
			if err != nil {
				log.Error("s.arc.Archives(%v) error(%v)", aids, err)
				return err
			}
			return nil
		})
	}
	g.Go(func() (err error) {
		us, err = s.acc.Infos(c, mids, "")
		if err != nil {
			log.Error("s.acc.Infos(%v) error(%v)", mids, err)
			return err
		}
		return nil
	})
	err := g.Wait()
	if u, ok := us[h.MID]; ok {
		h.Uname = u.Name
		h.Face = u.Face
	}
	for _, uid := range loveFans {
		if u, ok := us[uid]; ok {
			h.LoveFans = append(h.LoveFans, u)
		}
	}
	for _, uid := range playFans {
		if u, ok := us[uid]; ok {
			h.PlayFans = append(h.PlayFans, u)
		}
	}
	if hs.HottestAvInc != 0 {
		if arc, ok := arcs[int64(hs.HottestAvInc)]; ok {
			h.NewArchive = arc
		}
	}
	if hs.HottestAvNew != 0 {
		if arc, ok := arcs[int64(hs.HottestAvNew)]; ok {
			h.NewArchive = arc
		}
	}
	if hs.HottestAvAll != 0 {
		if arc, ok := arcs[int64(hs.HottestAvAll)]; ok {
			h.HotArchive = arc
		}
	}
	return err
}

func (s *Service) genToken(mid int64) string {
	bs := []byte(fmt.Sprintf("bili%d", mid*3+223333))
	hs := md5.Sum(bs)
	return hex.EncodeToString(hs[:])
}
