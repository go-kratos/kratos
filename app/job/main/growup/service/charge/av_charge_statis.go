package charge

import (
	"bytes"
	"context"
	"strconv"

	model "go-common/app/job/main/growup/model/charge"

	"go-common/library/log"
)

// GetAvChargeStatisMap get av charge statistics map
func (s *Service) GetAvChargeStatisMap(c context.Context) (chargeStatisMap map[int64]*model.AvChargeStatis, err error) {
	avChargeStatis, err := s.GetAvChargeStatis(c)
	if err != nil {
		log.Error("s.GetAvChargeStatis error(%v)", err)
		return
	}

	chargeStatisMap = make(map[int64]*model.AvChargeStatis)
	for _, chargeStatis := range avChargeStatis {
		chargeStatisMap[chargeStatis.AvID] = chargeStatis
	}
	return
}

// GetAvChargeStatis get av charge statistics
func (s *Service) GetAvChargeStatis(c context.Context) (avChargeStatis []*model.AvChargeStatis, err error) {
	var id int64
	for {
		statis, err1 := s.dao.AvChargeStatis(c, id, _limitSize)
		if err1 != nil {
			err = err1
			return
		}
		avChargeStatis = append(avChargeStatis, statis...)
		if len(statis) < _limitSize {
			break
		}
		id = statis[len(statis)-1].ID
	}
	return
}

// CalAvChargeStatis cal av charge statis
func (s *Service) CalAvChargeStatis(dailyCharge *model.AvCharge, chargeStatisMap map[int64]*model.AvChargeStatis) {
	if statisCharge, ok := chargeStatisMap[dailyCharge.AvID]; ok {
		updateAvChargeStatis(statisCharge, dailyCharge)
	} else {
		chargeStatisMap[dailyCharge.AvID] = addAvChargeStatis(dailyCharge)
	}
}

func addAvChargeStatis(daily *model.AvCharge) *model.AvChargeStatis {
	return &model.AvChargeStatis{
		AvID:        daily.AvID,
		MID:         daily.MID,
		TagID:       daily.TagID,
		IsOriginal:  daily.IsOriginal,
		UploadTime:  daily.UploadTime,
		TotalCharge: daily.IncCharge,
		DBState:     _dbInsert,
	}
}

func updateAvChargeStatis(avChargeStatis *model.AvChargeStatis, daily *model.AvCharge) {
	avChargeStatis.TotalCharge += daily.IncCharge
	avChargeStatis.DBState = _dbUpdate
}

// AvChargeStatisDBStore store data
func (s *Service) AvChargeStatisDBStore(c context.Context, chargeStatisMap map[int64]*model.AvChargeStatis) error {
	insert, update := make([]*model.AvChargeStatis, _batchSize), make([]*model.AvChargeStatis, _batchSize)
	insertIndex, updateIndex := 0, 0
	for _, charge := range chargeStatisMap {
		if charge.DBState == _dbInsert {
			insert[insertIndex] = charge
			insertIndex++
		} else if charge.DBState == _dbUpdate {
			update[updateIndex] = charge
			updateIndex++
		}

		if insertIndex >= _batchSize {
			_, err := s.avChargeStatisBatchInsert(c, insert[:insertIndex])
			if err != nil {
				log.Error("s.avChargeStatisBatchInsert error(%v)", err)
				return err
			}
			insertIndex = 0
		}

		if updateIndex >= _batchSize {
			_, err := s.avChargeStatisBatchInsert(c, update[:updateIndex])
			if err != nil {
				log.Error("s.avChargeStatisBatchInsert error(%v)", err)
				return err
			}
			updateIndex = 0
		}
	}

	if insertIndex > 0 {
		_, err := s.avChargeStatisBatchInsert(c, insert[:insertIndex])
		if err != nil {
			log.Error("s.avChargeStatisBatchInsert error(%v)", err)
			return err
		}
	}

	if updateIndex > 0 {
		_, err := s.avChargeStatisBatchInsert(c, update[:updateIndex])
		if err != nil {
			log.Error("s.avChargeStatisBatchInsert error(%v)", err)
			return err
		}
	}

	return nil
}

func assembleAvChargeStatis(avChargeStatis []*model.AvChargeStatis) (vals string) {
	var buf bytes.Buffer
	for _, row := range avChargeStatis {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.IsOriginal))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalCharge, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + row.UploadTime.Time().Format(_layoutSec) + "'")
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

func (s *Service) avChargeStatisBatchInsert(c context.Context, avChargeStatis []*model.AvChargeStatis) (rows int64, err error) {
	vals := assembleAvChargeStatis(avChargeStatis)
	rows, err = s.dao.InsertAvChargeStatisBatch(c, vals)
	return
}
