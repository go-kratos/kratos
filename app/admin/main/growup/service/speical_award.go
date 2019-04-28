package service

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/growup/dao"
	"go-common/app/admin/main/growup/dao/resource"
	"go-common/app/admin/main/growup/model"
	"go-common/app/admin/main/growup/util"

	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	awardEditable = iota + 1
	awardInfoEditable
	awardDivisionEditable
	awardPrizeEditable
	awardResourceEditable
)

const (
	_awardResourceTypeRule   = 1
	_awardResourceTypeDetail = 3
	_awardResourceTypeQ      = 5
	_awardResourceTypeA      = 6
)

func awardEditPerms(award *model.Award) (res map[int]bool) {
	var (
		now           = time.Now()
		notDisplay    = award.DisplayStatus == 1
		notFinished   = award.OpenStatus == 1
		preCycleStart = now.Before(award.CycleStart)
	)
	res = map[int]bool{
		awardEditable:         notDisplay || notFinished,
		awardInfoEditable:     notDisplay || preCycleStart,
		awardDivisionEditable: notDisplay || preCycleStart,
		awardPrizeEditable:    notDisplay || notFinished,
		awardResourceEditable: notDisplay || notFinished,
	}
	return
}

func (s *Service) generateAwardID() (awardID int64, err error) {
	awardID, err = util.NewSnowFlake().Generate()
	return
}

func (s *Service) validateAwardCycle(c context.Context, awardID int64, start, end time.Time) (ok bool, err error) {
	// simplified version: check all
	awardBases, err := s.dao.ListAward(c)
	if err != nil {
		return
	}
	between := func(t, a, b time.Time) bool {
		return !(t.Before(a) || t.After(b))
	}
	for _, v := range awardBases {
		if awardID == v.AwardID || v.DisplayStatus != 2 {
			continue
		}
		if between(start, v.CycleStart, v.CycleEnd) || between(end, v.CycleStart, v.CycleEnd) {
			ok = false
			return
		}
	}
	ok = true
	return
}

// AddAward .
func (s *Service) AddAward(c context.Context, arg *model.AddAwardArg, username string) (awardID int64, err error) {
	// 1. validation
	if !(arg.DisplayStatus == 1 || arg.DisplayStatus == 2) {
		err = ecode.RequestErr
		return
	}
	// 2. args
	// generate awardID
	awardID, err = s.generateAwardID()
	if err != nil {
		return
	}
	// cycle
	start := util.ToDayStart(time.Unix(arg.CycleStart, 0))
	end := util.ToDayEnd(time.Unix(arg.CycleEnd, 0))
	if arg.CycleStart > 0 && arg.CycleEnd > 0 && arg.DisplayStatus == 2 {
		var validCycle bool
		validCycle, err = s.validateAwardCycle(c, awardID, start, end)
		if err != nil {
			return
		}
		if !validCycle {
			err = ecode.Error(ecode.RequestErr, "评选周期重叠")
			return
		}
	}
	// total
	totalWinner, totalBonus := 0, 0
	if len(arg.Prizes) > 0 {
		for _, v := range arg.Prizes {
			totalWinner += v.Quota
			totalBonus += v.Bonus * v.Quota
		}
	}
	// 3. tx insert
	var (
		cycleStart   = start.Format("2006-01-02 15:04:05")
		cycleEnd     = end.Format("2006-01-02 15:04:05")
		announce     = util.ToDayNoon(time.Unix(arg.AnnounceDate, 0))
		announceDate = announce.Format("2006-01-02")
		openTime     = announce.Format("2006-01-02 15:04:05")
	)
	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		_, err = dao.AddAward(tx, awardID, arg.AwardName, cycleStart, cycleEnd, announceDate, openTime,
			arg.DisplayStatus, totalWinner, totalBonus, username)
		if err != nil {
			return
		}
		if len(arg.Divisions) > 0 {
			fields, values := saveDivisionStr(awardID, arg.Divisions)
			if values != "" {
				_, err = dao.SaveDivisions(tx, fields, values)
				if err != nil {
					return
				}
			}
		}
		if len(arg.Prizes) > 0 {
			fields, values := savePrizesStr(awardID, arg.Prizes)
			if values != "" {
				_, err = dao.SavePrizes(tx, fields, values)
				if err != nil {
					return
				}
			}
		}
		if arg.Resources != nil {
			fields, values := saveResourcesStr(awardID, arg.Resources)
			if values != "" {
				_, err = dao.SaveResource(tx, fields, values)
				if err != nil {
					return
				}
			}
		}
		return nil
	})

	return
}

func saveWinnerStr(winnersM map[int64]*model.AwardWinner) (fields, values string) {
	if len(winnersM) == 0 {
		return
	}
	fields = "award_id,mid,division_id,prize_id,tag_id"
	var buf bytes.Buffer
	for _, v := range winnersM {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(v.AwardID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(v.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(v.DivisionID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(v.PrizeID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(v.TagID, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	values += " ON DUPLICATE KEY UPDATE division_id=VALUES(division_id), prize_id=VALUES(prize_id), tag_id=VALUES(tag_id), is_deleted=0"
	buf.Reset()
	return
}

func saveDivisionStr(awardID int64, divisions []*model.AwardDivision) (fields, values string) {
	if len(divisions) == 0 {
		return
	}
	vs := make([]string, 0)
	for i, v := range divisions {
		vs = append(vs, fmt.Sprintf("(%d,%d,'%s',%d)", awardID, i+1, v.DivisionName, v.TagID))
	}
	if len(vs) > 0 {
		fields = "award_id,division_id,division_name,tag_id"
		values = strings.Join(vs, ",") + " ON DUPLICATE KEY UPDATE division_name=VALUES(division_name), tag_id=VALUES(tag_id), is_deleted=0"
	}
	return
}

func savePrizesStr(awardID int64, prizes []*model.AwardPrize) (fields, values string) {
	if len(prizes) == 0 {
		return
	}
	vs := make([]string, 0)
	for i, v := range prizes {
		vs = append(vs, fmt.Sprintf("(%d,%d,%d,%d)", awardID, i+1, v.Bonus, v.Quota))
	}
	if len(vs) > 0 {
		fields = "award_id,prize_id,bonus,quota"
		values = strings.Join(vs, ",") + " ON DUPLICATE KEY UPDATE bonus=VALUES(bonus), quota=VALUES(quota), is_deleted=0"
	}
	return
}

func delResourcesStr(awardID int64, arg *model.AwardResource) (delWhere string) {
	if len(arg.QA) == 0 {
		return fmt.Sprintf("award_id = %d AND resource_type IN (%d,%d)", awardID, _awardResourceTypeQ, _awardResourceTypeA)
	}
	idx := make([]int64, 0)
	for i := range arg.QA {
		idx = append(idx, int64(i+1))
	}
	return fmt.Sprintf("award_id = %d AND resource_type IN (%d,%d) AND resource_index NOT IN (%s)",
		awardID, _awardResourceTypeQ, _awardResourceTypeA, xstr.JoinInts(idx))
}

func saveResourcesStr(awardID int64, arg *model.AwardResource) (fields, values string) {
	if arg == nil {
		return
	}
	vs := make([]string, 0)
	vs = append(vs, fmt.Sprintf("(%d,%d,1,'%s')", awardID, _awardResourceTypeRule, arg.Rule))
	vs = append(vs, fmt.Sprintf("(%d,%d,1,'%s')", awardID, _awardResourceTypeDetail, arg.Detail))
	if len(arg.QA) > 0 {
		for i, qa := range arg.QA {
			vs = append(vs, fmt.Sprintf("(%d,%d,%d,'%s')", awardID, _awardResourceTypeQ, i+1, qa.Q))
			vs = append(vs, fmt.Sprintf("(%d,%d,%d,'%s')", awardID, _awardResourceTypeA, i+1, qa.A))
		}
	}
	if len(vs) > 0 {
		fields = "award_id,resource_type,resource_index,content"
		values = strings.Join(vs, ",") + " ON DUPLICATE KEY UPDATE content=VALUES(content), is_deleted=0"
	}
	return
}

// UpdateAward .
func (s *Service) UpdateAward(c context.Context, arg *model.SaveAwardArg) (err error) {
	// 1.
	if !(arg.DisplayStatus == 1 || arg.DisplayStatus == 2) || arg.AwardID == 0 {
		return ecode.Error(ecode.RequestErr, "illegal param")
	}
	awardID := arg.AwardID
	award, err := s.dao.Award(c, awardID)
	if err != nil {
		return
	}
	if award == nil {
		err = ecode.Error(ecode.RequestErr, "illegal award_id")
		return
	}
	//
	permission := awardEditPerms(award)
	if !permission[awardEditable] {
		err = ecode.Error(ecode.RequestErr, "award no longer editable")
		return
	}
	// sum
	totalQuota, totalBonus := 0, 0
	if len(arg.Prizes) > 0 {
		for _, v := range arg.Prizes {
			totalQuota += v.Quota
			totalBonus += v.Bonus * v.Quota
		}
	}
	newAward := &model.Award{
		AwardName:     arg.AwardName,
		CycleStart:    util.ToDayStart(time.Unix(arg.CycleStart, 0)),
		CycleEnd:      util.ToDayEnd(time.Unix(arg.CycleEnd, 0)),
		AnnounceDate:  time.Unix(arg.AnnounceDate, 0),
		OpenTime:      util.ToDayNoon(time.Unix(arg.AnnounceDate, 0)),
		DisplayStatus: arg.DisplayStatus,
		TotalQuota:    totalQuota,
		TotalBonus:    totalBonus,
	}

	if permission[awardInfoEditable] && arg.DisplayStatus == 2 {
		if !newAward.CycleStart.Equal(award.CycleStart) || !newAward.CycleEnd.Equal(award.CycleEnd) {
			var ok bool
			ok, err = s.validateAwardCycle(c, awardID, newAward.CycleStart, newAward.CycleEnd)
			if err != nil {
				return
			}
			if !ok {
				err = ecode.Error(ecode.RequestErr, "评选周期重叠")
				return
			}
		}
	}

	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		// lock
		award, err = dao.SelectAwardForUpdate(tx, awardID)
		if err != nil {
			return
		}
		if award == nil {
			return ecode.Error(ecode.RequestErr, "award not found")
		}
		//// award info
		if permission[awardInfoEditable] {
			updateAwardStr := fmt.Sprintf(`award_name='%s', cycle_start='%s', cycle_end='%s', 
		announce_date='%s', open_time='%s', display_status=%d, total_quota=%d, total_bonus=%d`,
				newAward.AwardName,
				newAward.CycleStart.Format("2006-01-02 15:04:05"),
				newAward.CycleEnd.Format("2006-01-02 15:04:05"),
				newAward.AnnounceDate.Format("2006-01-02"),
				newAward.OpenTime.Format("2006-01-02 15:04:05"),
				newAward.DisplayStatus, newAward.TotalQuota, newAward.TotalBonus)
			_, err = dao.UpdateAward(tx, awardID, updateAwardStr)
			if err != nil {
				return
			}
		}
		//// divisions
		if permission[awardDivisionEditable] {
			switch {
			case len(arg.Divisions) == 0:
				_, err = dao.DelDivisionAll(tx, awardID)
				if err != nil {
					return
				}
			default:
				// del
				argDivisionIDs := make([]int64, 0)
				for i, v := range arg.Divisions {
					v.DivisionID = int64(i + 1)
					argDivisionIDs = append(argDivisionIDs, v.DivisionID)
				}
				_, err = dao.DelDivisionsExclude(tx, awardID, argDivisionIDs)
				if err != nil {
					return
				}
				// save
				fields, values := saveDivisionStr(awardID, arg.Divisions)
				_, err = dao.SaveDivisions(tx, fields, values)
				if err != nil {
					return
				}
			}
		}
		//// prizes
		if permission[awardPrizeEditable] {
			switch {
			case len(arg.Prizes) == 0:
				_, err = dao.DelPrizeAll(tx, awardID)
				if err != nil {
					return
				}
			default:
				// del
				argPrizeIDs := make([]int64, 0)
				for i, v := range arg.Prizes {
					v.PrizeID = int64(i + 1)
					argPrizeIDs = append(argPrizeIDs, v.PrizeID)
				}
				_, err = dao.DelPrizesExclude(tx, awardID, argPrizeIDs)
				if err != nil {
					return
				}
				fields, values := savePrizesStr(awardID, arg.Prizes)
				_, err = dao.SavePrizes(tx, fields, values)
				if err != nil {
					return
				}
			}
		}
		//// resources
		if permission[awardResourceEditable] {
			switch {
			case arg.Resources == nil:
				_, err = dao.DelResources(tx, fmt.Sprintf("award_id = %d", awardID))
				if err != nil {
					return
				}
			default:
				delWhere := delResourcesStr(awardID, arg.Resources)
				_, err = dao.DelResources(tx, delWhere)
				if err != nil {
					return
				}
				fields, values := saveResourcesStr(awardID, arg.Resources)
				_, err = dao.SaveResource(tx, fields, values)
				if err != nil {
					return
				}
			}
		}
		return nil
	})
	return err
}

// ListAward .
func (s *Service) ListAward(c context.Context, from, limit int) (total int64, data []*model.AwardListModel, err error) {
	data = make([]*model.AwardListModel, 0)
	var (
		awardToWinnerC map[int64]int
		awards         []*model.Award
		divisions      []*model.AwardDivision
	)
	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		total, err = dao.CountAward(tx)
		if err != nil {
			return
		}
		if total == 0 {
			return
		}
		awards, err = dao.ListAward(tx, from, limit)
		if err != nil {
			return
		}
		if len(awards) == 0 {
			return
		}
		awardIDs := make([]int64, 0)
		for _, v := range awards {
			awardIDs = append(awardIDs, v.AwardID)
		}
		divisions, err = dao.ListAwardsDivision(tx, fmt.Sprintf("award_id IN (%s)", xstr.JoinInts(awardIDs)))
		if err != nil {
			return
		}
		awardToWinnerC, err = dao.GroupCountAwardWinner(tx, fmt.Sprintf("award_id IN (%s)", xstr.JoinInts(awardIDs)))
		if err != nil {
			return
		}
		return nil
	})
	if err != nil || total == 0 || len(awards) == 0 {
		return
	}
	// divisions group by award_id
	var (
		tagIDToName      map[int64]string
		awardToDivisions = make(map[int64][]*model.AwardDivision)
	)
	tagIDToName, err = resource.VideoCategoryIDToName(c)
	if err != nil {
		return
	}
	for _, division := range divisions {
		if _, ok := awardToDivisions[division.AwardID]; !ok {
			awardToDivisions[division.AwardID] = make([]*model.AwardDivision, 0)
		}
		awardToDivisions[division.AwardID] = append(awardToDivisions[division.AwardID], division)
	}
	for _, award := range awards {
		v := &model.AwardListModel{
			ID:              award.ID,
			AwardID:         award.AwardID,
			AwardName:       award.AwardName,
			CycleStart:      award.CycleStart.Unix(),
			CycleEnd:        award.CycleEnd.Unix(),
			TotalQuota:      award.TotalQuota,
			TotalBonus:      award.TotalBonus,
			AnnounceDate:    award.AnnounceDate.Unix(),
			OpenStatus:      award.OpenStatus,
			OpenTime:        award.OpenTime.Unix(),
			CTime:           award.CTime.Unix(),
			CreatedBy:       award.CreatedBy,
			SelectionStatus: 1,
			Tags:            make([]string, 0),
			DivisionNames:   make([]string, 0),
		}
		// status
		if count, ok := awardToWinnerC[v.AwardID]; ok && count == v.TotalQuota {
			v.SelectionStatus = 2
		}
		tagIDs := make([]int64, 0)
		for _, division := range awardToDivisions[v.AwardID] {
			tagIDs = append(tagIDs, division.TagID)
			v.DivisionNames = append(v.DivisionNames, division.DivisionName)
		}
		for _, tagID := range tagIDs {
			v.Tags = append(v.Tags, tagIDToName[tagID])
		}
		data = append(data, v)
	}
	return
}

// DetailAward .
func (s *Service) DetailAward(c context.Context, awardID int64) (data *model.AwardDetail, err error) {
	data = &model.AwardDetail{}
	var (
		winners        = 0
		typeIdxContent map[int]map[int]string
	)
	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		data.Award, err = dao.Award(tx, awardID)
		if err != nil {
			return
		}
		if data.Award == nil {
			err = ecode.Errorf(ecode.RequestErr, "award not found, awardID{%d}", awardID)
			return
		}
		data.Divisions, err = dao.ListDivision(tx, awardID)
		if err != nil {
			return
		}
		data.Prizes, err = dao.ListPrize(tx, awardID)
		if err != nil {
			return
		}
		typeIdxContent, err = dao.ListResource(tx, awardID)
		if err != nil {
			return
		}
		total, err := dao.CountAwardWinner(tx, awardID, "")
		if err != nil {
			return
		}
		winners = int(total)
		return nil
	})
	if err != nil {
		return
	}

	data.Award.SelectionStatus = 1
	if data.Award.TotalQuota == winners {
		data.Award.SelectionStatus = 2
	}
	// tags
	if len(data.Divisions) > 0 {
		var tagNames map[int64]string
		tagNames, err = resource.VideoCategoryIDToName(c)
		if err != nil {
			return
		}
		for _, division := range data.Divisions {
			if tagName, ok := tagNames[division.TagID]; ok {
				division.Tag = tagName
			}
		}
	}

	data.Resources = &model.AwardResource{
		Rule:   typeIdxContent[1][1],
		Detail: typeIdxContent[3][1],
		QA:     make([]*model.AwardQA, 0),
	}
	questions, answers := typeIdxContent[5], typeIdxContent[6]
	if len(questions) > 0 && len(answers) > 0 {
		for idx, question := range questions {
			data.Resources.QA = append(data.Resources.QA, &model.AwardQA{Index: idx, Q: question, A: answers[idx]})
		}
	}
	sort.Slice(data.Resources.QA, func(i, j int) bool {
		return data.Resources.QA[i].Index < data.Resources.QA[j].Index
	})

	data.Award.GenStr()

	return
}

// ListAwardWinner .
func (s *Service) ListAwardWinner(c context.Context, arg *model.QueryAwardWinnerArg) (total int64, res []*model.AwardWinner, err error) {
	res = make([]*model.AwardWinner, 0)
	if pass := s.validateQueryAwardWinner(c, arg); !pass {
		return
	}
	where := s.queryAwardWinnerWhere(arg)
	var (
		winnerInfo   map[int64]*model.AwardWinner   // <mid, *winner>
		prizeInfo    map[int64]*model.AwardPrize    // <prize_id, *prize>
		divisionInfo map[int64]*model.AwardDivision // <division_id, *division>
		mids         = make([]int64, 0)
	)
	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		total, err = dao.CountAwardWinner(tx, arg.AwardID, where)
		if err != nil {
			return
		}
		if total == 0 {
			return
		}
		winnerInfo, err = dao.QueryAwardWinner(tx, arg.AwardID, fmt.Sprintf("%s ORDER BY id DESC LIMIT %d,%d", where, arg.From, arg.Limit))
		if err != nil {
			return
		}
		if len(winnerInfo) > 0 {
			// divisions
			divisionInfo, err = dao.DivisionInfo(tx, arg.AwardID)
			if err != nil {
				return
			}
			// prizes
			prizeInfo, err = dao.PrizeInfo(tx, arg.AwardID)
			if err != nil {
				return
			}
		}
		return nil
	})
	if err != nil || total == 0 || len(winnerInfo) == 0 {
		return
	}
	categories, err := resource.VideoCategoryIDToName(c)
	if err != nil {
		return
	}
	for mid := range winnerInfo {
		mids = append(mids, mid)
	}
	upNames, err := resource.NamesByMIDs(c, mids)
	if err != nil {
		return
	}
	for mid, info := range winnerInfo {
		if division, ok := divisionInfo[info.DivisionID]; ok {
			info.DivisionName = division.DivisionName
		}
		if prize, ok := prizeInfo[info.PrizeID]; ok {
			info.Bonus = prize.Bonus
		}
		// tag
		if data, ok := categories[info.TagID]; ok {
			info.Tag = data
		}
		// nickname
		if nickname, ok := upNames[mid]; ok {
			info.Nickname = nickname
		}
	}
	res = make([]*model.AwardWinner, 0)
	for _, v := range winnerInfo {
		res = append(res, v)
	}
	return
}

// ExportAwardWinner .
func (s *Service) ExportAwardWinner(c context.Context, arg *model.QueryAwardWinnerArg) (res []byte, err error) {
	records, err := s.listAwardWinnerAll(c, arg)
	if err != nil {
		return
	}
	data := make([][]string, 0, len(records)+1)
	data = append(data, model.AwardWinnerExportFields())
	for _, v := range records {
		data = append(data, v.ExportStrings())
	}
	if res, err = FormatCSV(data); err != nil {
		log.Error("FormatCSV error(%v)", err)
	}
	return
}

func (s *Service) listAwardWinnerAll(c context.Context, arg *model.QueryAwardWinnerArg) (records []*model.AwardWinner, err error) {
	records = make([]*model.AwardWinner, 0)
	if pass := s.validateQueryAwardWinner(c, arg); !pass {
		return
	}
	where := s.queryAwardWinnerWhere(arg)
	var (
		winnerInfo   map[int64]*model.AwardWinner   // <mid, *winner>
		prizeInfo    map[int64]*model.AwardPrize    // <prize_id, *prize>
		divisionInfo map[int64]*model.AwardDivision // <division_id, *division>
		mids         = make([]int64, 0)
	)
	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		winnerInfo, err = dao.QueryAwardWinner(tx, arg.AwardID, fmt.Sprintf("%s ORDER BY id DESC", where))
		if err != nil {
			return
		}
		if len(winnerInfo) > 0 {
			// divisions
			divisionInfo, err = dao.DivisionInfo(tx, arg.AwardID)
			if err != nil {
				return
			}
			// prizes
			prizeInfo, err = dao.PrizeInfo(tx, arg.AwardID)
			if err != nil {
				return
			}
		}
		return nil
	})
	if err != nil {
		return
	}
	if len(winnerInfo) == 0 {
		return
	}
	categories, err := resource.VideoCategoryIDToName(c)
	if err != nil {
		return
	}
	for mid := range winnerInfo {
		mids = append(mids, mid)
	}
	upNames, err := resource.NamesByMIDs(c, mids)
	if err != nil {
		return
	}
	for mid, info := range winnerInfo {
		if division, ok := divisionInfo[info.DivisionID]; ok {
			info.DivisionName = division.DivisionName
		}
		if prize, ok := prizeInfo[info.PrizeID]; ok {
			info.Bonus = prize.Bonus
		}
		// tag
		if data, ok := categories[info.TagID]; ok {
			info.Tag = data
		}
		// nickname
		if nickname, ok := upNames[mid]; ok {
			info.Nickname = nickname
		}
	}
	for _, v := range winnerInfo {
		records = append(records, v)
	}
	return
}

func (s *Service) validateQueryAwardWinner(c context.Context, arg *model.QueryAwardWinnerArg) (pass bool) {
	if arg.Nickname != "" {
		mid, err := resource.MidByNickname(c, arg.Nickname)
		if err != nil || mid == 0 {
			return
		}
		if arg.MID == 0 {
			arg.MID = mid
		}
		if arg.MID != mid {
			log.Error("illegal mid(%d) and nickname(%s) pair", arg.MID, arg.Nickname)
			return
		}
	}
	return true
}

func (s *Service) queryAwardWinnerWhere(arg *model.QueryAwardWinnerArg) string {
	str := ""
	if arg.MID > 0 {
		str += fmt.Sprintf(" AND mid=%d", arg.MID)
	}
	if arg.TagID > 0 {
		str += fmt.Sprintf(" AND tag_id=%d", arg.TagID)
	}
	return str
}

// ReplaceAwardWinner .
func (s *Service) ReplaceAwardWinner(c context.Context, awardID, prevMID, mid int64) (err error) {
	records, err := s.dao.ListAwardRecord(c, awardID, fmt.Sprintf("AND (mid=%d OR mid=%d)", prevMID, mid))
	if err != nil {
		return
	}
	midsM := make(map[int64]bool)
	for _, v := range records {
		midsM[v.MID] = true
	}
	if !midsM[prevMID] {
		return ecode.Errorf(ecode.RequestErr, "mid(%d) not in award(%d)", prevMID, awardID)
	}
	if !midsM[mid] {
		return ecode.Errorf(ecode.RequestErr, "mid(%d) not in award(%d)", mid, awardID)
	}

	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		// lock award
		award, err := dao.SelectAwardForUpdate(tx, awardID)
		if err != nil {
			return
		}
		if award == nil {
			return ecode.Error(ecode.RequestErr, "award not found")
		}
		if award.OpenStatus != 1 {
			return ecode.Error(ecode.RequestErr, "illegal operation")
		}
		// winner info
		res, err := dao.QueryAwardWinner(tx, awardID, fmt.Sprintf("AND (mid=%d OR mid=%d)", prevMID, mid))
		if err != nil {
			return
		}
		prevWinner, ok := res[prevMID]
		if !ok {
			return ecode.Error(ecode.RequestErr, "invalid old winner")
		}
		if _, ok = res[mid]; ok {
			return ecode.Error(ecode.RequestErr, "invalid new winner")
		}
		// replace
		rows, err := dao.DelWinner(tx, awardID, fmt.Sprintf(" AND mid=%d", prevMID))
		if err != nil {
			return
		}
		if rows != 1 {
			return ecode.Error(ecode.ServerErr, "failed to del old winner")
		}
		fields, values := saveWinnerStr(map[int64]*model.AwardWinner{
			awardID: {
				AwardID:    awardID,
				MID:        mid,
				DivisionID: prevWinner.DivisionID,
				PrizeID:    prevWinner.PrizeID,
				TagID:      prevWinner.TagID,
			},
		})
		rows, err = dao.SaveWinners(tx, fields, values)
		if err != nil {
			return
		}
		if rows != 1 {
			return ecode.Error(ecode.ServerErr, "failed to add new winner")
		}
		return
	})
	return
}

// SaveAwardResult .
func (s *Service) SaveAwardResult(c context.Context, arg *model.AwardResult) (err error) {
	var (
		awardID      = arg.AwardID
		winnersM     = make(map[int64]*model.AwardWinner)
		prizeWinnerC = make(map[int64]int)
		divisionInfo map[int64]*model.AwardDivision
	)
	divisionInfo, err = s.dao.AwardDivisionInfo(c, awardID)
	if err != nil {
		return
	}
	for divisionI, division := range arg.Divisions {
		divisionID := int64(divisionI + 1)
		for prizeI, prize := range division.Prizes {
			prizeID := int64(prizeI + 1)
			prizeWinnerC[prizeID] += len(prize.MIDs)
			for _, mid := range prize.MIDs {
				if _, ok := winnersM[mid]; ok {
					return ecode.Error(ecode.RequestErr, "UID重复，请重新录入")
				}
				winnersM[mid] = &model.AwardWinner{
					MID:        mid,
					AwardID:    awardID,
					DivisionID: divisionID,
					PrizeID:    prizeID,
					TagID:      divisionInfo[divisionID].TagID,
				}
			}
		}
	}
	mids := make([]int64, 0, len(winnersM))
	for k := range winnersM {
		mids = append(mids, k)
	}
	if len(winnersM) == 0 {
		return ecode.Error(ecode.RequestErr, "奖项人数不符，保存失败")
	}
	// 是否已报名专项奖
	{
		var winnerRecords []*model.AwardRecord
		winnerRecords, err = s.dao.ListAwardRecord(c, awardID, fmt.Sprintf("AND mid IN (%s)", xstr.JoinInts(mids)))
		if err != nil {
			return
		}
		if len(winnersM) != len(winnerRecords) {
			var (
				awardRecordM = make(map[int64]bool)
				unSignedMIDs []int64
			)
			for _, rcr := range winnerRecords {
				awardRecordM[rcr.MID] = true
			}
			for _, mid := range mids {
				if !awardRecordM[mid] {
					unSignedMIDs = append(unSignedMIDs, mid)
				}
			}
			return ecode.Errorf(ecode.RequestErr, "%s 未报名专项奖，保存失败", xstr.JoinInts(unSignedMIDs))
		}
	}
	// 是否已签约激励
	{
		var signedUPs map[int64]struct{}
		signedUPs, err = s.dao.GetUpInfoByState(c, "up_info_video", mids, 3)
		if err != nil {
			return
		}
		if len(mids) != len(signedUPs) {
			var unSignedMIDs []int64
			for _, mid := range mids {
				if _, ok := signedUPs[mid]; !ok {
					unSignedMIDs = append(unSignedMIDs, mid)
				}
			}
			return ecode.Errorf(ecode.RequestErr, "%s 非签约状态，保存失败", xstr.JoinInts(unSignedMIDs))
		}
	}
	// save winners
	updateAwardStr := fmt.Sprintf("open_time = '%s'", time.Unix(arg.OpenTime, 0).Format("2006-01-02 15:04:05"))
	fields, values := saveWinnerStr(winnersM)
	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		award, err := dao.SelectAwardForUpdate(tx, awardID)
		if err != nil {
			return
		}
		if award == nil {
			return ecode.Errorf(ecode.RequestErr, "award not found, awardID{%d}", awardID)
		}
		if award.OpenStatus != 1 {
			return ecode.Error(ecode.RequestErr, "已发奖，不能再修改名单")
		}
		if award.TotalQuota != len(mids) {
			return ecode.Error(ecode.RequestErr, "奖项人数不符，保存失败")
		}
		prizeInfo, err := dao.PrizeInfo(tx, awardID)
		if err != nil {
			return
		}
		for prizeID, info := range prizeInfo {
			if prizeWinnerC[prizeID] != info.Quota {
				return ecode.Error(ecode.RequestErr, "奖项人数不符，保存失败")
			}
		}
		_, err = dao.UpdateAward(tx, awardID, updateAwardStr)
		if err != nil {
			return
		}
		_, err = dao.DelWinnerAll(tx, awardID)
		if err != nil {
			return
		}
		_, err = dao.SaveWinners(tx, fields, values)
		if err != nil {
			return
		}
		return nil
	})
	return
}

// AwardResult .
func (s *Service) AwardResult(c context.Context, awardID int64) (res *model.AwardResult, err error) {
	var (
		award        *model.Award
		winners      []*model.AwardWinner
		prizeInfo    map[int64]*model.AwardPrize
		divisionInfo map[int64]*model.AwardDivision
		data         = make(map[int64]map[int64][]int64)
	)
	err = s.dao.DoInTx(c, func(tx *sql.Tx) (err error) {
		award, err = dao.Award(tx, awardID)
		if err != nil {
			return
		}
		if award == nil {
			err = ecode.Errorf(ecode.RequestErr, "award not found, awardID{%d}", awardID)
			return
		}
		winners, err = dao.AwardWinnerAll(tx, awardID)
		if err != nil {
			return
		}
		// divisions
		divisionInfo, err = dao.DivisionInfo(tx, awardID)
		if err != nil {
			return
		}
		// prizes
		prizeInfo, err = dao.PrizeInfo(tx, awardID)
		if err != nil {
			return
		}
		return nil
	})
	if err != nil {
		return
	}
	// init
	for divisionID := range divisionInfo {
		data[divisionID] = make(map[int64][]int64)
		for prizeID := range prizeInfo {
			data[divisionID][prizeID] = make([]int64, 0)
		}
	}
	// mids
	for _, v := range winners {
		data[v.DivisionID][v.PrizeID] = append(data[v.DivisionID][v.PrizeID], v.MID)
	}
	// res
	res = &model.AwardResult{
		AwardID:      awardID,
		OpenTime:     award.OpenTime.Unix(),
		AnnounceDate: award.AnnounceDate.Unix(),
		CycleEnd:     award.CycleEnd.Unix(),
		Divisions:    make([]*model.AwardDivisionResult, 0),
	}
	for divisionID, division := range divisionInfo {
		dv := &model.AwardDivisionResult{
			DivisionID:   divisionID,
			DivisionName: division.DivisionName,
			Prizes:       make([]*model.AwardPrizeResult, 0),
		}
		for prizeID := range prizeInfo {
			mids := data[divisionID][prizeID]
			dv.Prizes = append(dv.Prizes, &model.AwardPrizeResult{MIDs: mids, PrizeID: prizeID})
		}
		sort.Slice(dv.Prizes, func(i, j int) bool {
			return dv.Prizes[i].PrizeID < dv.Prizes[j].PrizeID
		})
		res.Divisions = append(res.Divisions, dv)
	}
	sort.Slice(res.Divisions, func(i, j int) bool {
		return res.Divisions[i].DivisionID < res.Divisions[j].DivisionID
	})

	return
}
