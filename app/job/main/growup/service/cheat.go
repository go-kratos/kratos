package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	_spyArchiveCoin     = 26 //稿件硬币
	_spyArchiveFavorite = 27 //稿件收藏
	_spyArchivePlay     = 28 //稿件播放
	_spyUpFans          = 29 //异常粉丝数
	query               = "{\"select\": [{\"name\": \"id\", \"as\": \"id\"}," +
		"{\"name\": \"log_date\",\"as\": \"log_date\"},{\"name\": \"target_mid\",\"as\": \"target_mid\"},{\"name\": \"target_id\",\"as\": \"target_id\"}," +
		"{\"name\": \"event_id\",\"as\": \"event_id\"},{\"name\": \"state\",\"as\": \"state\"},{\"name\": \"type\",\"as\": \"type\"},{\"name\": \"quantity\",\"as\": \"quantity\"}," +
		"{\"name\": \"isdel\",\"as\": \"isdel\"}],\"where\": {\"log_date\": {\"in\": [\"%s\"]}},\"sort\": {\"id\": -1},\"page\": {\"skip\": %d,\"limit\": %d}}"
)

// UpdateCheatHTTP update cheat by http
func (s *Service) UpdateCheatHTTP(c context.Context, date time.Time) (err error) {
	err = s.CheatStatistics(c, date)
	if err != nil {
		log.Error("s.UpdateSpy UpdateSpyData error(%v)", err)
	}
	return
}

// CheatStatistics task update cheat
func (s *Service) CheatStatistics(c context.Context, date time.Time) (err error) {
	spies, err := s.getSpy(c, date)
	if err != nil {
		return
	}
	// first filter
	cs := cheats(spies)

	// aggregation by mid
	cm := aggreByMID(cs)
	log.Info("agreegation by mid spies:", len(cm))

	// aggregation by av_id
	am := aggreByAvID(cs)
	log.Info("agreegation by av_id spies:", len(cm))

	// get up base info
	upm, err := s.getUps(c, cm)
	if err != nil {
		return
	}
	// get av base info
	avm, err := s.getAvs(c, am, time.Now().Add(-30*24*time.Hour))
	if err != nil {
		return
	}

	pcs, err := s.playCount(c, cm)
	if err != nil {
		return
	}

	deducted, err := s.breachRecord(c)
	if err != nil {
		return
	}

	var ups []*model.Cheating
	for mid, cheat := range cm {
		if up, ok := upm[mid]; ok {
			cheat.Nickname = up.Nickname
			cheat.Fans = up.Fans
			cheat.SignedAt = up.SignedAt
			cheat.AccountState = 3
			cheat.PlayCount = pcs[mid]
			ups = append(ups, cheat)
		}
	}

	var avs []*model.Cheating
	for avID, cheat := range am {
		if av, ok := avm[avID]; ok {
			cheat.UploadTime = av.UploadTime
			cheat.TotalIncome = av.TotalIncome
			cheat.Nickname = cm[cheat.MID].Nickname
			if deducted[avID] {
				cheat.Deducted = 1
			} else {
				cheat.Deducted = 0
			}
			avs = append(avs, cheat)
		}
	}

	log.Info("signed cheat up count:", len(ups))
	err = s.batchInsertCheats(c, ups, s.batchInsertCheatUps)
	if err != nil {
		log.Error("batchInsertCheatUps error(%v)", err)
		return
	}

	log.Info("signed cheat av count:", len(avs))
	err = s.batchInsertCheats(c, avs, s.batchInsertCheatArchives)
	if err != nil {
		log.Error("batchInsertCheatArchives error(%v)", err)
		return
	}
	return
}

func (s *Service) breachRecord(c context.Context) (deducted map[int64]bool, err error) {
	deducted = make(map[int64]bool)
	var id int64
	for {
		var ds map[int64]bool
		id, ds, err = s.dao.AvBreachRecord(c, id, 2000)
		if err != nil {
			return
		}
		if len(ds) == 0 {
			break
		}
		for k, v := range ds {
			deducted[k] = v
		}
	}
	return
}

func (s *Service) getSpy(c context.Context, date time.Time) (spies []*model.Spy, err error) {
	from, limit := 0, 500
	var info []*model.Spy
	for {
		dateStr := date.Format("20060102")
		info, err = s.dp.SendSpyRequest(c, fmt.Sprintf(query, dateStr, from, limit))
		if err != nil {
			log.Error("s.getSpyData error(%v)", err)
			return
		}
		if len(info) == 0 {
			break
		}
		spies = append(spies, info...)
		from += len(info)
	}
	log.Info("get spy data total (%d) rows", from)
	return
}

func cheats(spies []*model.Spy) (cs []*model.Cheating) {
	for _, spy := range spies {
		c := &model.Cheating{}
		c.MID = spy.TargetMID
		switch spy.EventID {
		case _spyArchiveCoin:
			c.CheatCoin = spy.Quantity
			c.AvID = spy.TargetID
		case _spyArchiveFavorite:
			c.CheatFavorite = spy.Quantity
			c.AvID = spy.TargetID
		case _spyArchivePlay:
			c.CheatPlayCount = spy.Quantity
			c.AvID = spy.TargetID
		case _spyUpFans:
			c.CheatFans = spy.Quantity
		}
		cs = append(cs, c)
	}
	return
}

func aggreByMID(source []*model.Cheating) (cheats map[int64]*model.Cheating) {
	cheats = make(map[int64]*model.Cheating)
	for _, c := range source {
		if cheat, ok := cheats[c.MID]; !ok {
			cheats[c.MID] = &model.Cheating{
				MID:            c.MID,
				CheatFans:      c.CheatFans,
				CheatPlayCount: c.CheatPlayCount,
			}
		} else {
			cheat.CheatFans += c.CheatFans
			cheat.CheatPlayCount += c.CheatPlayCount
		}
	}
	return
}

func aggreByAvID(source []*model.Cheating) (cheats map[int64]*model.Cheating) {
	cheats = make(map[int64]*model.Cheating)
	for _, c := range source {
		if c.AvID == 0 {
			continue
		}
		if cheat, ok := cheats[c.AvID]; !ok {
			cheats[c.AvID] = &model.Cheating{
				MID:            c.MID,
				AvID:           c.AvID,
				CheatPlayCount: c.CheatPlayCount,
				CheatCoin:      c.CheatCoin,
				CheatFavorite:  c.CheatFavorite,
			}
		} else {
			cheat.CheatPlayCount += c.CheatPlayCount
			cheat.CheatCoin += c.CheatCoin
			cheat.CheatFavorite += c.CheatFavorite
		}
	}
	return
}

// Ups get ups in up_info_video
func (s *Service) getUps(c context.Context, cheats map[int64]*model.Cheating) (ups map[int64]*model.Cheating, err error) {
	ups = make(map[int64]*model.Cheating)
	var mids []int64
	for mid := range cheats {
		mids = append(mids, mid)
		if len(mids) == 200 {
			var nc map[int64]*model.Cheating
			nc, err = s.dao.Ups(c, mids)
			if err != nil {
				return
			}
			for k, v := range nc {
				if v.IsDeleted == 0 {
					ups[k] = v
				}
			}
			mids = make([]int64, 0)
			time.Sleep(200 * time.Millisecond)
		}
	}

	if len(mids) > 0 {
		var nc map[int64]*model.Cheating
		nc, err = s.dao.Ups(c, mids)
		if err != nil {
			return
		}
		for k, v := range nc {
			ups[k] = v
		}
	}
	return
}

// avs result key: av_id, value: cheating with total_income, upload_time
func (s *Service) getAvs(c context.Context, cheats map[int64]*model.Cheating, mtime time.Time) (avs map[int64]*model.Cheating, err error) {
	avs = make(map[int64]*model.Cheating)
	var avIds []int64
	for avID := range cheats {
		avIds = append(avIds, avID)
		if len(avIds) == 200 {
			var nc map[int64]*model.Cheating
			nc, err = s.dao.Avs(c, mtime, avIds)
			if err != nil {
				return
			}
			for k, v := range nc {
				avs[k] = v
			}
			avIds = make([]int64, 0)
			time.Sleep(200 * time.Millisecond)
		}
	}
	if len(avIds) > 0 {
		var na map[int64]*model.Cheating
		na, err = s.dao.Avs(c, mtime, avIds)
		if err != nil {
			return
		}
		for k, v := range na {
			avs[k] = v
		}
	}
	return
}

// key: mid, value: play_count
func (s *Service) playCount(c context.Context, cheats map[int64]*model.Cheating) (pcs map[int64]int64, err error) {
	var mids []int64
	pcs = make(map[int64]int64)
	for mid := range cheats {
		mids = append(mids, mid)
		if len(mids) == 200 {
			var pc map[int64]int64
			pc, err = s.dao.PlayCount(c, mids)
			if err != nil {
				return
			}
			for k, v := range pc {
				pcs[k] = v
			}
			mids = make([]int64, 200)
		}
	}
	if len(mids) > 0 {
		var pc map[int64]int64
		pc, err = s.dao.PlayCount(c, mids)
		if err != nil {
			return
		}
		for k, v := range pc {
			pcs[k] = v
		}
	}
	return
}

type insertCheats func(c context.Context, cheats []*model.Cheating) (err error)

func (s *Service) batchInsertCheats(c context.Context, cheats []*model.Cheating, insert insertCheats) (err error) {
	var end int
	for range cheats {
		end++
		if end%2000 == 0 {
			err = insert(c, cheats[:end])
			if err != nil {
				return
			}
			cheats = cheats[end:]
			end = 0
		}
	}
	if end > 0 {
		err = insert(c, cheats)
	}
	return
}

func (s *Service) batchInsertCheatUps(c context.Context, cheats []*model.Cheating) (err error) {
	var buf bytes.Buffer
	for _, cheat := range cheats {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(cheat.MID, 10))
		buf.WriteByte(',')
		buf.WriteString("\"" + cheat.SignedAt.Time().Format("2006-01-02 15:04:05") + "\"")
		buf.WriteByte(',')
		buf.WriteString("\"" + cheat.Nickname + "\"")
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.Fans))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.CheatFans))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(cheat.PlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.CheatPlayCount))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.AccountState))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values := buf.String()
	buf.Reset()
	_, err = s.dao.InsertCheatUps(c, values)
	return
}

func (s *Service) batchInsertCheatArchives(c context.Context, cheats []*model.Cheating) (err error) {
	var buf bytes.Buffer
	for _, cheat := range cheats {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(cheat.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(cheat.MID, 10))
		buf.WriteByte(',')
		buf.WriteString("\"" + cheat.Nickname + "\"")
		buf.WriteByte(',')
		buf.WriteString("\"" + cheat.UploadTime.Time().Format("2006-01-02 15:04:05") + "\"")
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.TotalIncome))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.CheatPlayCount))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.CheatFavorite))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.CheatCoin))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(cheat.Deducted))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values := buf.String()
	buf.Reset()
	_, err = s.dao.InsertCheatArchives(c, values)
	return
}
