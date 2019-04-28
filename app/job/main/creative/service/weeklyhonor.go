package service

import (
	"context"
	"strings"
	"time"

	wkhmdl "go-common/app/interface/main/creative/model/weeklyhonor"
	upgrpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	forbidMidMap = map[int64]struct{}{
		24:        {},
		2:         {},
		517999:    {},
		9099524:   {},
		208259:    {},
		202466803: {},
		40016273:  {},
		245482023: {},
		84089650:  {},
		31465698:  {},
		22160843:  {},
		3098848:   {},
		39592268:  {},
		223931175: {},
	}
	honorMap map[int][]*wkhmdl.HonorWord
)

// SendMsg .
func (s *Service) SendMsg() {
	log.Info("Start SendMsg")
	var (
		c      = context.TODO()
		lastid int64
		size   = 1000
	)
	for {
		var (
			upActives []*upgrpc.UpActivity
			newid     int64
			err       error
		)
		for i := 0; i < 5; i++ {
			upActives, newid, err = s.honDao.UpActivesList(c, lastid, size)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Error("s.honDao.UpActivesList(%d,%d) error(%v)", lastid, size, err)
			break
		}
		// filter mid
		mids, err := s.filterMids(c, upActives)
		if err != nil {
			continue
		}
		var errMids []int64
		if len(mids) > 0 {
			for i := 1; i < 3; i++ {
				if errMids, err = s.honDao.SendNotify(c, mids); err != nil {
					log.Error("s.honDao.SendNotify(%v) error(%v)", mids, err)
					continue
				}
				break
			}
		}
		go s.infocMsgStat(mids, errMids)
		if len(upActives) < size {
			break
		}
		lastid = newid
		time.Sleep(time.Second)
	}
	log.Info("Finish SendMsg")
}

// FlushHonor .
func (s *Service) FlushHonor() {
	log.Info("FlushHonor Start")
	var (
		c         = context.TODO()
		lastid    int64
		size      = 1000
		batchSize = s.c.HonorStep
	)
	for {
		upActives, newid, err := s.honDao.UpActivesList(c, lastid, size)
		if err != nil {
			log.Error("s.honDao.UpActivesList(%d,%d) error(%v)", lastid, size, err)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		var mids []int64
		for _, v := range upActives {
			mids = append(mids, v.Mid)
		}
		g := new(errgroup.Group)
		var pmids []int64
		routines := len(mids)/batchSize + 1
		for i := 0; i < routines; i++ {
			if i == routines-1 {
				pmids = mids[i*batchSize:]
			} else {
				pmids = mids[i*batchSize : (i+1)*batchSize]
			}
			t := pmids
			g.Go(func() (err error) {
				err = s.upsertHonor(c, t)
				if err != nil {
					log.Error("s.upsertHonor(%v) error(%v)", t, err)
					return err
				}
				return nil
			})
		}
		g.Wait()
		if len(mids) < size {
			break
		}
		lastid = newid
	}
	log.Info("FlushHonor done")
}

func (s *Service) upsertHonor(c context.Context, mids []int64) error {
	now := time.Now()
	day := int(now.Weekday())
	date := now.AddDate(0, 0, -day-1)
	saturday := date.Format("20060102")
LOOP:
	for _, mid := range mids {
		var hls map[int]*wkhmdl.HonorLog
		hls, err := s.honDao.HonorLogs(c, mid)
		if err != nil {
			log.Error("s.honDao.HonorLogs(%d) error(%v)", mid, err)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		var lastHid int
		for _, v := range hls {
			if int64(v.MTime) > date.Unix() && int64(v.MTime) < date.AddDate(0, 0, 7).Unix() {
				continue LOOP
			}
			if int64(v.MTime) < date.Unix() && int64(v.MTime) > date.AddDate(0, 0, -7).Unix() {
				lastHid = v.HID
			}
		}
		var hs *wkhmdl.HonorStat
		for i := 0; i < 3; i++ {
			hs, err = s.honDao.HonorStat(c, mid, saturday)
			if err != nil {
				log.Error("s.honDao.HonorStat(%d,%v) error(%v)", mid, saturday, err)
				continue
			}
			if hs != nil {
				break
			}
		}
		if hs == nil {
			log.Error("FlushHonor nil hs mid(%d)", mid)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		newHid := hs.GenHonor(mid, lastHid)
		affected, err := s.honDao.UpsertCount(c, mid, newHid)
		if err != nil || affected == 0 {
			log.Error("s.honDao.UpsertCount(%d,%d) affceted(%d) error(%v)", mid, newHid, affected, err)
		}
		log.Info("FlushHonor mid(%d)", mid)
	}
	return nil
}

// TestSendMsg .
func (s *Service) TestSendMsg(c context.Context, mids []int64) (err error) {
	for i := 1; i < 3; i++ {
		if _, err = s.honDao.SendNotify(c, mids); err != nil {
			log.Error("s.honDao.SendNotify(%v) error(%v)", mids, err)
			continue
		} else {
			break
		}
	}
	return
}

func (s *Service) filterMids(c context.Context, upActives []*upgrpc.UpActivity) (mids []int64, err error) {
	if len(upActives) == 0 {
		return
	}
	var rawMids []int64
	for _, v := range upActives {
		rawMids = append(rawMids, v.Mid)
	}
	hls, err := s.honDao.LatestHonorLogs(c, rawMids)
	if err != nil {
		log.Error("failed to get latest honor logs err(%v)", err)
		return
	}
	midClickMap, err := s.honDao.ClickCounts(c, rawMids)
	if err != nil {
		log.Error("failed to get honor click count err (%v)", err)
		return
	}
	highLevMap := highLevMidMap(hls)
	for _, v := range upActives {
		var subState uint8
		if subState, err = s.honDao.GetUpSwitch(c, v.Mid); subState == wkhmdl.HonorUnSub {
			continue
		}
		if err != nil {
			log.Error("s.honDao.GetUpSwitch mid(%d) err(%v)", v.Mid, err)
		}
		var cnt int
		if cnt, err = s.honDao.UpCount(c, v.Mid); err == nil && cnt == 0 {
			continue
		}
		if highLevMap[v.Mid] {
			mids = append(mids, v.Mid)
			continue
		}
		if _, ok := forbidMidMap[v.Mid]; ok {
			continue
		}
		if v.Activity > 3 {
			continue
		}
		if _, ok := midClickMap[v.Mid]; ok {
			mids = append(mids, v.Mid)
			continue
		}
		if sunday := wkhmdl.LatestSunday(); s.c.SendEveryWeek || isOddWeek(sunday) {
			mids = append(mids, v.Mid)
		}
	}
	return mids, err
}

func isOddWeek(date time.Time) bool {
	_, w := date.ISOWeek()
	return w%2 != 0
}

func highLevMidMap(hls []*wkhmdl.HonorLog) (res map[int64]bool) {
	if honorMap == nil {
		honorMap = wkhmdl.HMap()
	}
	res = make(map[int64]bool)
	for _, h := range hls {
		hws, ok := honorMap[h.HID]
		if !ok || len(hws) == 0 {
			continue
		}
		last := len(hws) - 1
		res[h.MID] = isHighLev(hws[last].Priority)
	}
	return
}

func isHighLev(p string) bool {
	return strings.HasPrefix(p, "A") || strings.HasPrefix(p, "R") || strings.HasPrefix(p, "S")
}

func (s *Service) infocMsgStat(mids, errMids []int64) {
	errMidsMap := make(map[int64]bool)
	for _, mid := range errMids {
		errMidsMap[mid] = true
	}
	for _, mid := range mids {
		var success int32
		if !errMidsMap[mid] {
			success = 1
		}
		err := s.honDao.HonorInfoc(context.Background(), mid, success)
		if err != nil {
			log.Error("failed to log honor infoc,mid(%d),success(%d),err(%v)", mid, success, err)
		}
	}
}
