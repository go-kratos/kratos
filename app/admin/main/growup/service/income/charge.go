package income

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/main/growup/dao/resource"
	model "go-common/app/admin/main/growup/model/income"

	"go-common/library/log"
)

// ArchiveChargeStatis av_charge statis
func (s *Service) ArchiveChargeStatis(c context.Context, categoryID []int64, typ, groupType int, fromTime, toTime int64) (data interface{}, err error) {
	table := setChargeTableByGroup(typ, groupType)
	from := getDateByGroup(groupType, time.Unix(fromTime, 0))
	to := getDateByGroup(groupType, time.Unix(toTime, 0))
	query := formatArchiveQuery(categoryID, from, to)

	archives, err := s.GetArchiveChargeStatis(c, table, query)
	if err != nil {
		log.Error("s.GetArchiveChargeStatis error(%v)", err)
		return
	}

	data = archiveChargeStatis(archives, from, to, groupType)
	return
}

func archiveChargeStatis(archs []*model.ArchiveChargeStatis, from, to time.Time, groupType int) interface{} {
	avsMap := make(map[string]*model.ArchiveChargeStatis)
	ctgyMap := make(map[string]bool)
	for _, arch := range archs {
		date := formatDateByGroup(arch.CDate.Time(), groupType)
		ctgykey := date + strconv.FormatInt(arch.CategroyID, 10)
		if val, ok := avsMap[date]; ok {
			val.Avs += arch.Avs
			if !ctgyMap[ctgykey] {
				val.Charge += arch.Charge
				ctgyMap[ctgykey] = true
			}
		} else {
			avsMap[date] = &model.ArchiveChargeStatis{
				Avs:    arch.Avs,
				Charge: arch.Charge,
			}
			ctgyMap[ctgykey] = true
		}
	}

	charge, counts, xAxis := []string{}, []int64{}, []string{}
	// get result by date
	to = to.AddDate(0, 0, 1)
	for from.Before(to) {
		dateStr := formatDateByGroup(from, groupType)
		xAxis = append(xAxis, dateStr)
		if val, ok := avsMap[dateStr]; ok {
			charge = append(charge, fmt.Sprintf("%.2f", float64(val.Charge)/float64(100)))
			counts = append(counts, val.Avs)
		} else {
			charge = append(charge, "0")
			counts = append(counts, int64(0))
		}
		from = addDayByGroup(groupType, from)
	}

	return map[string]interface{}{
		"counts":  counts,
		"charges": charge,
		"xaxis":   xAxis,
	}
}

// ArchiveChargeSection get av_charge section
func (s *Service) ArchiveChargeSection(c context.Context, categoryID []int64, typ, groupType int, fromTime, toTime int64) (data interface{}, err error) {
	table := setChargeTableByGroup(typ, groupType)
	from := getDateByGroup(groupType, time.Unix(fromTime, 0))
	to := getDateByGroup(groupType, time.Unix(toTime, 0))
	query := formatArchiveQuery(categoryID, from, to)

	archives, err := s.GetArchiveChargeStatis(c, table, query)
	if err != nil {
		log.Error("s.GetArchiveChargeStatis error(%v)", err)
		return
	}

	data = archiveChargeSection(archives, from, to, groupType)
	return
}

func archiveChargeSection(archs []*model.ArchiveChargeStatis, from, to time.Time, groupType int) interface{} {
	ret := make([]map[string]interface{}, 0)
	avsMap := make(map[string][]int64)
	for _, arch := range archs {
		date := formatDateByGroup(arch.CDate.Time(), groupType)
		if val, ok := avsMap[date]; ok {
			val[arch.MoneySection] += arch.Avs
		} else {
			avsMap[date] = make([]int64, 12)
			avsMap[date][arch.MoneySection] = arch.Avs
			ret = append(ret, map[string]interface{}{
				"date_format": date,
				"sections":    avsMap[date],
			})
		}
	}
	return ret
}

// ArchiveChargeDetail archive charge details
func (s *Service) ArchiveChargeDetail(c context.Context, aid int64, typ int) (archives []*model.ArchiveCharge, err error) {
	switch typ {
	case _video:
		archives, err = s.GetAvCharges(c, aid)
	case _column:
		archives, err = s.dao.GetColumnCharges(c, aid)
	case _bgm:
		archives, err = s.dao.GetBgmCharges(c, aid)
	default:
		err = fmt.Errorf("type error")
	}
	if err != nil {
		log.Error("s.GetArchives(%d) error(%v)", typ, err)
		return
	}
	err = s.archiveChargeDetail(c, archives, aid, typ)
	return
}

func (s *Service) archiveChargeDetail(c context.Context, archs []*model.ArchiveCharge, aid int64, typ int) (err error) {
	if len(archs) == 0 {
		return
	}
	// get up nickname
	nickname, err := resource.NameByMID(c, archs[0].MID)
	if err != nil {
		return
	}

	var table, query string
	switch typ {
	case _video:
		table, query = "av_charge_statis", fmt.Sprintf("av_id = %d", aid)
	case _column:
		table, query = "column_charge_statis", fmt.Sprintf("aid = %d", aid)
	case _bgm:
		table, query = "bgm_charge_statis", fmt.Sprintf("sid = %d", aid)
	}
	totalCharge, err := s.dao.GetTotalCharge(c, table, query)
	if err != nil {
		log.Error("s.GetTotalCharge error(%v)", err)
		return
	}
	sort.Slice(archs, func(i, j int) bool {
		return archs[i].Date > archs[j].Date
	})
	for _, arch := range archs {
		arch.TotalCharge = totalCharge
		arch.Nickname = nickname
		totalCharge -= arch.Charge
	}
	return
}

// BgmChargeDetail bgm charge detail
func (s *Service) BgmChargeDetail(c context.Context, sid int64) (archives []*model.ArchiveCharge, err error) {
	archives = make([]*model.ArchiveCharge, 0)
	bgms, err := s.dao.GetBgmCharges(c, sid)
	if err != nil {
		log.Error("s.dao.GetBgmCharges error(%v)", err)
		return
	}
	avIDs := make(map[int64]struct{})
	for _, bgm := range bgms {
		avIDs[bgm.AvID] = struct{}{}
	}
	for avID := range avIDs {
		var avs []*model.ArchiveCharge
		avs, err = s.ArchiveChargeDetail(c, avID, _video)
		if err != nil {
			log.Error("s.ArchiveChargeDetail error(%v)", err)
			return
		}
		archives = append(archives, avs...)
	}
	return
}

// UpRatio up charge ratio
func (s *Service) UpRatio(c context.Context, from, limit int64) (map[int64]int64, error) {
	return s.dao.UpRatio(c, from, limit)
}

// GetAvCharges get av charge by av id
func (s *Service) GetAvCharges(c context.Context, avID int64) (avs []*model.ArchiveCharge, err error) {
	avs = make([]*model.ArchiveCharge, 0)
	for i := 1; i <= 12; i++ {
		var av []*model.ArchiveCharge
		av, err = s.dao.GetAvDailyCharge(c, i, avID)
		if err != nil {
			return
		}
		avs = append(avs, av...)
	}
	return
}

// GetArchiveChargeStatis get archive charge date statis
func (s *Service) GetArchiveChargeStatis(c context.Context, table, query string) (archs []*model.ArchiveChargeStatis, err error) {
	offset, size := 0, 2000
	for {
		var arch []*model.ArchiveChargeStatis
		arch, err = s.dao.GetArchiveChargeStatis(c, table, query, offset, size)
		if err != nil {
			return
		}
		archs = append(archs, arch...)
		if len(arch) < size {
			break
		}
		offset += len(arch)
	}
	return
}
