package charge

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
	xtime "go-common/library/time"
)

// SectionEntries section entries
type SectionEntries struct {
	daily   []*model.DateStatis
	weekly  []*model.DateStatis
	monthly []*model.DateStatis
}

func initChargeSections(charge, tagID int64, date xtime.Time) []*model.DateStatis {
	chargeSections := make([]*model.DateStatis, 12)
	chargeSections[0] = initChargeSection(0, 1, 0, charge, tagID, date)
	chargeSections[1] = initChargeSection(1, 5, 1, charge, tagID, date)
	chargeSections[2] = initChargeSection(5, 10, 2, charge, tagID, date)
	chargeSections[3] = initChargeSection(10, 20, 3, charge, tagID, date)
	chargeSections[4] = initChargeSection(20, 50, 4, charge, tagID, date)
	chargeSections[5] = initChargeSection(50, 100, 5, charge, tagID, date)
	chargeSections[6] = initChargeSection(100, 200, 6, charge, tagID, date)
	chargeSections[7] = initChargeSection(200, 500, 7, charge, tagID, date)
	chargeSections[8] = initChargeSection(500, 1000, 8, charge, tagID, date)
	chargeSections[9] = initChargeSection(1000, 3000, 9, charge, tagID, date)
	chargeSections[10] = initChargeSection(3000, 5000, 10, charge, tagID, date)
	chargeSections[11] = initChargeSection(5000, math.MaxInt32, 11, charge, tagID, date)
	return chargeSections
}

func initChargeSection(min, max, section, charge, tagID int64, date xtime.Time) *model.DateStatis {
	var tips string
	if max == math.MaxInt32 {
		tips = fmt.Sprintf("\"%d+\"", min)
	} else {
		tips = fmt.Sprintf("\"%d~%d\"", min, max)
	}
	return &model.DateStatis{
		MinCharge:    min,
		MaxCharge:    max,
		MoneySection: section,
		MoneyTips:    tips,
		Charge:       charge,
		CategoryID:   tagID,
		CDate:        date,
	}
}

func (s *Service) handleDateStatis(c context.Context, sourceCh chan []*model.Archive, date time.Time, table string) (sections []*model.DateStatis, err error) {
	// delete
	if table != "" {
		_, err = s.dao.DelStatisTable(c, table, date.Format(_layout))
		if err != nil {
			log.Error("s.dao.DelChargeStatisTable error(%v)", err)
			return
		}
	}
	// add
	sections = s.handleArchives(c, sourceCh, date)
	return
}

// HandleAv handle archive_charge_daily_statis, archive_charge_weekly_statis, archive_charge_monthly_statis
func (s *Service) handleArchives(c context.Context, archiveCh chan []*model.Archive, date time.Time) (sections []*model.DateStatis) {
	archiveTagMap := make(map[int64]map[int64]int64) // key TagID, value map[int64]int64 -> key avId, value charge
	tagChargeMap := make(map[int64]int64)            // key TagID, value TagID total Charge
	handleArchive(archiveCh, archiveTagMap, tagChargeMap, date)

	sections = make([]*model.DateStatis, 0)
	for tagID, archives := range archiveTagMap {
		section := countDateStatis(archives, tagChargeMap[tagID], tagID, date)
		sections = append(sections, section...)
	}
	return
}

func handleArchive(archiveCh chan []*model.Archive, archiveTagMap map[int64]map[int64]int64, tagChargeMap map[int64]int64, startDate time.Time) {
	for archives := range archiveCh {
		for _, ac := range archives {
			if !startDate.After(ac.Date.Time()) {
				tagChargeMap[ac.TagID] += ac.IncCharge
				if _, ok := archiveTagMap[ac.TagID]; !ok {
					archiveTagMap[ac.TagID] = make(map[int64]int64)
				}
				archiveTagMap[ac.TagID][ac.ID] += ac.IncCharge
			}
		}
	}
}

func countDateStatis(charges map[int64]int64, totalCharge, tagID int64, date time.Time) (sections []*model.DateStatis) {
	if len(charges) == 0 {
		return
	}
	sections = initChargeSections(totalCharge, tagID, xtime.Time(date.Unix()))
	for _, charge := range charges {
		for _, section := range sections {
			min, max := section.MinCharge*100, section.MaxCharge*100
			if charge >= min && charge < max {
				section.Count++
			}
		}
	}
	return
}

func (s *Service) dateStatisInsert(c context.Context, avChargeSection []*model.DateStatis, table string) (rows int64, err error) {
	var buf bytes.Buffer
	for _, row := range avChargeSection {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.Count, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MoneySection, 10))
		buf.WriteByte(',')
		buf.WriteString(row.MoneyTips)
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Charge, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CategoryID, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.CDate.Time().Format(_layout) + "'")
		buf.WriteString(")")
		buf.WriteByte(',')
	}

	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals := buf.String()
	buf.Reset()
	rows, err = s.dao.InsertStatisTable(c, table, vals)
	return
}
