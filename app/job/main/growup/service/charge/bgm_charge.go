package charge

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
	xtime "go-common/library/time"
	"golang.org/x/sync/errgroup"
)

func (s *Service) handleBgm(c context.Context, date time.Time,
	dailyChannel chan []*model.BgmCharge) (weekly, monthly map[string]*model.BgmCharge, statis map[string]*model.BgmStatis, err error) {
	var eg errgroup.Group
	weekly = make(map[string]*model.BgmCharge)
	monthly = make(map[string]*model.BgmCharge)
	statis = make(map[string]*model.BgmStatis)

	eg.Go(func() (err error) {
		weekly, err = s.getBgmCharge(c, getStartWeeklyDate(date), "bgm_weekly_charge")
		if err != nil {
			log.Error("s.GetBgmCharge(bgm_weekly_charge) error(%v)", err)
			return
		}
		return
	})

	eg.Go(func() (err error) {
		monthly, err = s.getBgmCharge(c, getStartMonthlyDate(date), "bgm_monthly_charge")
		if err != nil {
			log.Error("s.GetBgmCharge(bgm_monthly_charge) error(%v)", err)
			return
		}
		return
	})

	eg.Go(func() (err error) {
		statis, err = s.getBgmStatis(c)
		if err != nil {
			log.Error("s.getBgmStatisMap error(%v)", err)
		}
		return
	})

	if err = eg.Wait(); err != nil {
		log.Error("handleBgm eg.Wait error(%v)", err)
		return
	}

	s.calBgms(date, weekly, monthly, statis, dailyChannel)
	return
}

func (s *Service) calBgms(date time.Time, weeklyMap, monthlyMap map[string]*model.BgmCharge, statisMap map[string]*model.BgmStatis, dailyChannel chan []*model.BgmCharge) {
	for daily := range dailyChannel {
		s.calBgm(date, daily, weeklyMap, monthlyMap, statisMap)
	}
}

func (s *Service) calBgm(date time.Time, daily []*model.BgmCharge, weeklyMap, monthlyMap map[string]*model.BgmCharge, statisMap map[string]*model.BgmStatis) {
	for _, charge := range daily {
		key := fmt.Sprintf("%d+%d", charge.SID, charge.AID)
		if weekly, ok := weeklyMap[key]; ok {
			updateBgmCharge(weekly, charge)
		} else {
			weeklyMap[key] = addBgmCharge(charge, startWeeklyDate)
		}
		if monthly, ok := monthlyMap[key]; ok {
			updateBgmCharge(monthly, charge)
		} else {
			monthlyMap[key] = addBgmCharge(charge, startMonthlyDate)
		}
		if statis, ok := statisMap[key]; ok {
			updateBgmStatis(statis, charge)
		} else {
			statisMap[key] = addBgmStatis(charge)
		}
	}
}

func addBgmCharge(daily *model.BgmCharge, fixDate time.Time) *model.BgmCharge {
	return &model.BgmCharge{
		SID:       daily.SID,
		AID:       daily.AID,
		MID:       daily.MID,
		CID:       daily.CID,
		JoinAt:    daily.JoinAt,
		Title:     daily.Title,
		IncCharge: daily.IncCharge,
		Date:      xtime.Time(fixDate.Unix()),
		DBState:   _dbInsert,
	}
}

func updateBgmCharge(origin, daily *model.BgmCharge) {
	origin.IncCharge += daily.IncCharge
	origin.DBState = _dbUpdate
}

func addBgmStatis(daily *model.BgmCharge) *model.BgmStatis {
	return &model.BgmStatis{
		SID:         daily.SID,
		AID:         daily.AID,
		MID:         daily.MID,
		CID:         daily.CID,
		JoinAt:      daily.JoinAt,
		Title:       daily.Title,
		TotalCharge: daily.IncCharge,
		DBState:     _dbInsert,
	}
}

func updateBgmStatis(statis *model.BgmStatis, daily *model.BgmCharge) {
	statis.TotalCharge += daily.IncCharge
	statis.DBState = _dbUpdate
}

func (s *Service) getBgmCharge(c context.Context, date time.Time, table string) (bgms map[string]*model.BgmCharge, err error) {
	bgms = make(map[string]*model.BgmCharge)
	var id int64
	for {
		var bgm []*model.BgmCharge
		bgm, err = s.dao.BgmCharge(c, date, id, _limitSize, table)
		if err != nil {
			return
		}
		for _, b := range bgm {
			key := fmt.Sprintf("%d+%d", b.SID, b.AID)
			bgms[key] = b
		}
		if len(bgm) < _limitSize {
			break
		}
		id = bgm[len(bgm)-1].ID
	}
	return
}

func (s *Service) getBgmStatis(c context.Context) (bgms map[string]*model.BgmStatis, err error) {
	bgms = make(map[string]*model.BgmStatis)
	var id int64
	for {
		var bgm []*model.BgmStatis
		bgm, err = s.dao.BgmStatis(c, id, _limitSize)
		if err != nil {
			return
		}
		for _, b := range bgm {
			key := fmt.Sprintf("%d+%d", b.SID, b.AID)
			bgms[key] = b
		}
		if len(bgm) < _limitSize {
			break
		}
		id = bgm[len(bgm)-1].ID
	}
	return
}

func (s *Service) bgmCharges(c context.Context, date time.Time, ch chan []*model.BgmCharge, avBgmCharge chan []*model.AvCharge) (dailyMap map[string]*model.BgmCharge, err error) {
	dailyMap = make(map[string]*model.BgmCharge)
	bgms, err := s.GetBgms(c, int64(_limitSize))
	if err != nil {
		log.Error("s.GetBgms error(%v)", err)
		return
	}
	defer func() {
		close(ch)
	}()
	for avs := range avBgmCharge {
		for _, av := range avs {
			bgm, ok := bgms[av.AvID]
			if !ok {
				continue
			}
			incCharge := int64(round(div(mul(float64(av.IncCharge), float64(0.3)), float64(len(bgm))), 0))
			if incCharge == 0 {
				continue
			}
			charges := make([]*model.BgmCharge, 0, len(bgm))
			for _, b := range bgm {
				if b.MID == av.MID {
					continue
				}
				c := &model.BgmCharge{
					SID:       b.SID,
					MID:       b.MID,
					AID:       b.AID,
					CID:       b.CID,
					IncCharge: incCharge,
					JoinAt:    b.JoinAt,
					Date:      av.Date,
					Title:     b.Title,
					DBState:   _dbInsert,
				}
				dailyMap[fmt.Sprintf("%d+%d", c.SID, c.AID)] = c
				charges = append(charges, c)
			}
			ch <- charges
		}
	}
	return
}

// GetBgms map[av_id][]*model.Bgm
func (s *Service) GetBgms(c context.Context, limit int64) (bm map[int64][]*model.Bgm, err error) {
	var id int64
	bm = make(map[int64][]*model.Bgm)
	for {
		var bs []*model.Bgm
		bs, id, err = s.dao.GetBgm(c, id, limit)
		if err != nil {
			return
		}
		if len(bs) == 0 {
			break
		}
		for _, b := range bs {
			if _, ok := bm[b.AID]; ok {
				bm[b.AID] = append(bm[b.AID], b)
			} else {
				bm[b.AID] = []*model.Bgm{b}
			}
		}
	}
	return
}

func round(val float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}

func div(x, y float64) float64 {
	a := big.NewFloat(x)
	b := big.NewFloat(y)
	c := new(big.Float).Quo(a, b)
	d, _ := c.Float64()
	return d
}

func mul(x, y float64) float64 {
	a := big.NewFloat(x)
	b := big.NewFloat(y)
	c := new(big.Float).Mul(a, b)
	d, _ := c.Float64()
	return d
}

func (s *Service) bgmDBStore(c context.Context, table string, bgmMap map[string]*model.BgmCharge) error {
	insert, update := make([]*model.BgmCharge, _batchSize), make([]*model.BgmCharge, _batchSize)
	insertIndex, updateIndex := 0, 0
	for _, charge := range bgmMap {
		if charge.DBState == _dbInsert {
			insert[insertIndex] = charge
			insertIndex++
		} else if charge.DBState == _dbUpdate {
			update[updateIndex] = charge
			updateIndex++
		}

		if insertIndex >= _batchSize {
			_, err := s.bgmBatchInsert(c, insert[:insertIndex], table)
			if err != nil {
				log.Error("s.bgmBatchInsert error(%v)", err)
				return err
			}
			insertIndex = 0
		}

		if updateIndex >= _batchSize {
			_, err := s.bgmBatchInsert(c, update[:updateIndex], table)
			if err != nil {
				log.Error("s.bgmBatchInsert error(%v)", err)
				return err
			}
			updateIndex = 0
		}
	}

	if insertIndex > 0 {
		_, err := s.bgmBatchInsert(c, insert[:insertIndex], table)
		if err != nil {
			log.Error("s.bgmBatchInsert error(%v)", err)
			return err
		}
	}

	if updateIndex > 0 {
		_, err := s.bgmBatchInsert(c, update[:updateIndex], table)
		if err != nil {
			log.Error("s.bgmBatchInsert error(%v)", err)
			return err
		}
	}
	return nil
}

func assembleBgmCharge(charges []*model.BgmCharge) (vals string) {
	var buf bytes.Buffer
	for _, row := range charges {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.SID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CID, 10))
		buf.WriteByte(',')
		buf.WriteString("\"" + strings.Replace(row.Title, "\"", "\\\"", -1) + "\"")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.IncCharge, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.Date.Time().Format(_layout) + "'")
		buf.WriteByte(',')
		buf.WriteString("'" + row.JoinAt.Time().Format(_layoutSec) + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}

func (s *Service) bgmBatchInsert(c context.Context, charges []*model.BgmCharge, table string) (rows int64, err error) {
	vals := assembleBgmCharge(charges)
	rows, err = s.dao.InsertBgmChargeTable(c, vals, table)
	return
}

func (s *Service) bgmStatisDBStore(c context.Context, statisMap map[string]*model.BgmStatis) error {
	insert, update := make([]*model.BgmStatis, _batchSize), make([]*model.BgmStatis, _batchSize)
	insertIndex, updateIndex := 0, 0
	for _, charge := range statisMap {
		if charge.DBState == _dbInsert {
			insert[insertIndex] = charge
			insertIndex++
		} else if charge.DBState == _dbUpdate {
			update[updateIndex] = charge
			updateIndex++
		}
		if insertIndex >= _batchSize {
			_, err := s.bgmStatisBatchInsert(c, insert[:insertIndex])
			if err != nil {
				log.Error("s.bgmStatisBatchInsert error(%v)", err)
				return err
			}
			insertIndex = 0
		}

		if updateIndex >= _batchSize {
			_, err := s.bgmStatisBatchInsert(c, update[:updateIndex])
			if err != nil {
				log.Error("s.bgmStatisBatchInsert error(%v)", err)
				return err
			}
			updateIndex = 0
		}
	}

	if insertIndex > 0 {
		_, err := s.bgmStatisBatchInsert(c, insert[:insertIndex])
		if err != nil {
			log.Error("s.bgmStatisBatchInsert error(%v)", err)
			return err
		}
	}

	if updateIndex > 0 {
		_, err := s.bgmStatisBatchInsert(c, update[:updateIndex])
		if err != nil {
			log.Error("s.bgmStatisBatchInsert error(%v)", err)
			return err
		}
	}
	return nil
}

func assembleBgmStatis(statis []*model.BgmStatis) (vals string) {
	var buf bytes.Buffer
	for _, row := range statis {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.SID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CID, 10))
		buf.WriteByte(',')
		buf.WriteString("\"" + strings.Replace(row.Title, "\"", "\\\"", -1) + "\"")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalCharge, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.JoinAt.Time().Format(_layoutSec) + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals = buf.String()
	buf.Reset()
	return
}

func (s *Service) bgmStatisBatchInsert(c context.Context, statis []*model.BgmStatis) (rows int64, err error) {
	vals := assembleBgmStatis(statis)
	rows, err = s.dao.InsertBgmStatisBatch(c, vals)
	return
}
