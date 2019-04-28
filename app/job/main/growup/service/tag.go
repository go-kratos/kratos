package service

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

// DeleteAvRatio clean the av_charge_ratio
func (s *Service) DeleteAvRatio(c context.Context, limit int64) (rows int64, err error) {
	rows, err = s.dao.DelAvRatio(c, limit)
	if err != nil {
		log.Error("growup-job s.DeleteAvRatio error(%v)", err)
	}
	return
}

// DeleteUpIncome clean the up_tag_income
func (s *Service) DeleteUpIncome(c context.Context, limit int64) (rows int64, err error) {
	rows, err = s.dao.DelIncome(c, limit)
	if err != nil {
		log.Error("growup-job s.DeleteUpIncome error(%v)", err)
	}
	return
}

// DelActivity clean activity_info table.
func (s *Service) DelActivity(c context.Context, limit int64) (rows int64, err error) {
	rows, err = s.dao.DelActivity(c, limit)
	if err != nil {
		log.Error("growup-job s.DelActivity error(%v)", err)
	}
	return
}

func (s *Service) aids(c context.Context, categoryID int, date time.Time) (aids []int64, err error) {
	var offset int64
	m := int(date.Month())
	var sm string
	if m < 10 {
		sm += "0"
		sm += strconv.Itoa(m)
	}
	for {
		var as []*model.AID
		as, err = s.dao.AIDs(c, categoryID, offset, conf.Conf.Ratio.Limit, sm, date)
		if err != nil {
			return
		}
		if len(as) == 0 {
			break
		}
		offset += int64(len(as))
		for _, a := range as {
			if a.IsDeleted == 0 && a.IncCharge > 0 {
				aids = append(aids, a.AvID)
			}
		}
	}
	return
}

func (s *Service) aidsByMID(c context.Context, mid int64, categoryID int, date time.Time) (aids []int64, err error) {
	var offset int64
	m := int(date.Month())
	var sm string
	if m < 10 {
		sm += "0"
		sm += strconv.Itoa(m)
	}

	for {
		var as []*model.AID
		as, err = s.dao.AIDsByMID(c, mid, categoryID, offset, conf.Conf.Ratio.Limit, sm, date)
		if err != nil {
			return
		}
		if len(as) == 0 {
			break
		}
		offset += int64(len(as))
		for _, a := range as {
			if a.IsDeleted == 0 && a.IncCharge > 0 {
				aids = append(aids, a.AvID)
			}
		}
	}
	return
}

// ExecRatioForHTTP exec http
func (s *Service) ExecRatioForHTTP(c context.Context, year int, month int, day int) (err error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	err = s.execAvRatio(c, t)
	if err != nil {
		log.Error("growup-job s.execAvRatio error!(%v)", err)
	}
	return
}

func (s *Service) execAvRatio(c context.Context, date time.Time) (err error) {
	commonTags, err := s.dao.CommonTagInfo(c, date)
	if err != nil {
		return
	}
	m := make(map[int64]*model.TagUpInfo)
	//m := make(map[int64]int)
	for _, commonTag := range commonTags {
		var as []int64
		as, err = s.aids(c, commonTag.Category, date)
		if err != nil {
			return
		}
		getAvRatio(as, commonTag.TagID, commonTag.Ratio, m)
	}
	noCommonTags, err := s.dao.NoCommonTagInfo(c, date)
	if err != nil {
		return
	}
	for _, noCommonTag := range noCommonTags {
		var as []int64
		as, err = s.aidsByMID(c, noCommonTag.MID, noCommonTag.Category, date)
		if err != nil {
			return
		}
		getAvRatio(as, noCommonTag.TagID, noCommonTag.Ratio, m)
	}

	activityTags, err := s.dao.ActivityTagInfo(c, date)
	if err != nil {
		log.Error("s.dao.ActivityTagInfo error(%v)", err)
		return
	}
	var res []*model.ActivityAVInfo
	res, err = s.getActivityInfo(c, activityTags)
	if err != nil {
		log.Error("s.getActivityaInfo error(%v)", err)
		return
	}
	tm, err := s.handleActivityData(c, activityTags, res)
	if err != nil {
		log.Error("s.handleActivityData error(%v)", err)
		return
	}

	for key, value := range tm {
		var avIDs []int64
		var ratio int
		for _, v := range value {
			avIDs = append(avIDs, v.AVID)
			ratio = v.Ratio
		}
		getAvRatio(avIDs, key, ratio, m)
	}

	err = s.insertAvRatio(c, m)
	if err != nil {
		log.Error("s.insertAvRatio error(%v)", err)
		return
	}
	log.Info("s.insertAvRatio av_map len:%d", len(m))
	return
}

func (s *Service) insertAvRatio(c context.Context, m map[int64]*model.TagUpInfo) (err error) {
	var buf bytes.Buffer
	var count int64
	for k, v := range m {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(v.TagID, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.FormatInt(k, 10))
		buf.WriteString(",")
		buf.WriteString(strconv.Itoa(v.Ratio))
		buf.WriteString("),")
		count++
		if count%2000 == 0 {
			buf.Truncate(buf.Len() - 1)
			_, err = s.dao.InsertRatio(c, buf.String())
			if err != nil {
				log.Error("growup-job s.insertAvRatio error(%v)", err)
				return
			}
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
		_, err = s.dao.InsertRatio(c, buf.String())
		if err != nil {
			log.Error("growup-job s.insertAvRatio error(%v)", err)
			return
		}
	}
	return
}

func (s *Service) handleActivityData(c context.Context, tag []*model.TagUpInfo, res []*model.ActivityAVInfo) (tm map[int64][]*model.ActivityAVInfo, err error) {
	tm = make(map[int64][]*model.ActivityAVInfo)
	types, err := s.getAllTypes(c)
	if err != nil {
		log.Error("s.handleActivityData getAllTypes error(%v)", err)
		return
	}
	for _, act := range tag {
		for _, r := range res {
			if act.IsCommon == 0 {
				_, err = s.dao.ActivityMIDExist(c, act.TagID, r.MID)
				if err != nil {
					err = nil
					continue
				}
			}
			if (getPIDByID(types, int16(r.Category)) == int16(act.Category)) && (r.ActivityID == act.ActivityID) && (act.Business == 1) {
				_, ok := tm[act.TagID]
				if !ok {
					arr := make([]*model.ActivityAVInfo, 0)
					v := &model.ActivityAVInfo{ActivityID: r.ActivityID, TagID: act.TagID, AVID: r.AVID, MID: r.MID, Category: r.Category, Ratio: act.Ratio}
					arr = append(arr, v)
					tm[act.TagID] = arr
				} else {
					v := &model.ActivityAVInfo{ActivityID: r.ActivityID, TagID: act.TagID, AVID: r.AVID, MID: r.MID, Category: r.Category, Ratio: act.Ratio}
					tm[act.TagID] = append(tm[act.TagID], v)
				}
			}
		}
	}
	return
}

func (s *Service) getActivityInfo(c context.Context, activityTags []*model.TagUpInfo) (res []*model.ActivityAVInfo, err error) {
	res = make([]*model.ActivityAVInfo, 0)
	am := make(map[int64]bool)
	var ai []int64
	for _, act := range activityTags {
		if _, ok := am[act.ActivityID]; !ok {
			am[act.ActivityID] = true
			ai = append(ai, act.ActivityID)
		}
	}
	var ts []*model.ActivityAVInfo
	var tv []int64
	var total int
	for i, a := range ai {
		tv = append(tv, a)
		if (i+1)%20 == 0 {
			pn := 1
			ps := 30
			ts, total, err = s.dao.GetActivityAVInfo(c, pn, ps, tv)
			if err != nil {
				log.Error("s.dao.ActivityTagInfo error(%v)", err)
				return
			}
			if len(ts) == 0 {
				return
			}

			for ps*pn <= (total + ps) {
				ts, total, err = s.dao.GetActivityAVInfo(c, pn, ps, tv)
				if err != nil {
					log.Error("s.dao.ActivityTagInfo error(%v)", err)
					return
				}
				pn++
				res = append(res, ts...)
			}
			res = append(res, ts...)
			ts = ts[0:0:0]
			tv = tv[0:0:0]
		}
	}
	if len(tv) > 0 {
		pn := 1
		ps := 30
		ts, total, err = s.dao.GetActivityAVInfo(c, pn, ps, tv)
		if err != nil {
			log.Error("s.dao.ActivityTagInfo error(%v)", err)
			return
		}
		if len(ts) == 0 {
			return
		}
		for pn*ps <= (total + ps) {
			ts, total, err = s.dao.GetActivityAVInfo(c, pn, ps, tv)
			if err != nil {
				log.Error("s.dao.ActivityTagInfo error(%v)", err)
				return
			}
			pn++
			res = append(res, ts...)
		}
	}
	log.Info("res len(%d), total(%d)", len(res), total)
	return
}

// UpdateAvRatio update av ratio everyday
func (s *Service) UpdateAvRatio() (err error) {
	c := context.TODO()
	yesterday := time.Now().Add(-time.Hour * 24)
	t := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	err = s.execAvRatio(c, t)
	if err != nil {
		log.Error("growup-job s.execAvRatio error!")
	}
	return
}

func getAvRatio(aids []int64, tagID int64, ratio int, m map[int64]*model.TagUpInfo) {
	for _, avid := range aids {
		v, ok := m[avid]
		if ok {
			if v.Ratio < ratio {
				m[avid].Ratio = ratio
				m[avid].TagID = tagID
			}
		} else {
			a := &model.TagUpInfo{TagID: tagID, Ratio: ratio}
			m[avid] = a
		}
	}
}

// ExecIncomeForHTTP income
func (s *Service) ExecIncomeForHTTP(c context.Context, year int, month int, day int) (err error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	for {
		rows, _ := s.DeleteAvRatio(c, conf.Conf.Ratio.Num)
		time.Sleep(time.Duration(conf.Conf.Ratio.Sleep) * time.Millisecond)
		if rows == 0 {
			break
		}
	}
	for {
		rows, _ := s.DelActivity(c, conf.Conf.Ratio.Num)
		time.Sleep(time.Duration(conf.Conf.Ratio.Sleep) * time.Millisecond)
		if rows == 0 {
			break
		}
	}
	err = s.ExecRatioForHTTP(c, year, month, day)
	if err != nil {
		log.Error("Exec avRatio from http error!(%v)", err)
		return
	}
	err = s.InsertTagIncome(c, t)
	if err != nil {
		log.Error("s.InsertTagIncome error!")
	}
	return
}

func (s *Service) getAllTypes(c context.Context) (tm map[int16]int16, err error) {
	tm, err = s.dao.GetAllTypes(c)
	if err != nil {
		log.Error("s.getAllTypes error(%v)", err)
	}
	return
}

func getPIDByID(tm map[int16]int16, id int16) (pid int16) {
	for k, v := range tm {
		if k == id {
			return v
		}
	}
	return
}
