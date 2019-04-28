package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go-common/app/service/main/reply-feed/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_getSlotsStat        = "SELECT name,slot,state FROM reply_abtest_strategy"
	_getSlotsStatManager = "SELECT name,slot,algorithm,weight,state FROM reply_abtest_strategy"
	_setSlot             = "UPDATE reply_abtest_strategy SET name=?,algorithm=?,weight=?,state=? WHERE slot IN (%s)"
	_setWeight           = "UPDATE reply_abtest_strategy SET weight=? WHERE name=?"
	_modifyState         = "UPDATE reply_abtest_strategy SET state=? WHERE name=?"
	_getSlotsStatByName  = "SELECT slot,algorithm,weight FROM reply_abtest_strategy WHERE name=?"

	_getStatisticsDate = "SELECT name,date,hour,view,total_view,hot_view,hot_click,hot_like,hot_hate,hot_child,hot_report,total_like,total_hate,total_report,total_root,total_child,hot_like_uv,hot_hate_uv,hot_report_uv,hot_child_uv,total_like_uv,total_hate_uv,total_report_uv,total_child_uv,total_root_uv" +
		" FROM reply_abtest_statistics WHERE date>=? AND date<=? AND name!='default'"

	_upsertLog = "INSERT INTO reply_abtest_statistics (name,date,hour,view,hot_click,hot_view,total_view) VALUES(?,?,?,?,?,?,?)" +
		" ON DUPLICATE KEY UPDATE view=view+?,hot_click=hot_click+?,hot_view=hot_view+?,total_view=total_view+?"
)

var (
	_countIdleSlot = fmt.Sprintf("SELECT COUNT(*) FROM reply_abtest_strategy WHERE state=1 AND name='%s'", model.DefaultSlotName)
	_getIdelSlots  = fmt.Sprintf("SELECT slot FROM reply_abtest_strategy WHERE state=1 AND name='%s' ORDER BY slot ASC LIMIT ?", model.DefaultSlotName)
)

/*
SlotsStat
*/

// SlotsMapping get slot name stat from database.
func (d *Dao) SlotsMapping(ctx context.Context) (slotsMap map[string]*model.SlotsMapping, err error) {
	slotsMap = make(map[string]*model.SlotsMapping)
	rows, err := d.db.Query(ctx, _getSlotsStat)
	if err != nil {
		log.Error("db.Query(%s) args(%s) error(%v)", _getSlotsStat, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name  string
			slot  int
			state int
		)
		if err = rows.Scan(&name, &slot, &state); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		mapping, ok := slotsMap[name]
		if ok {
			mapping.Slots = append(mapping.Slots, slot)
		} else {
			mapping = &model.SlotsMapping{
				Name:  name,
				Slots: []int{slot},
				State: state,
			}
		}
		slotsMap[name] = mapping
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// SlotsStatManager get slots stat from database, used by manager.
func (d *Dao) SlotsStatManager(ctx context.Context) (s []*model.SlotsStat, err error) {
	slotsMap := make(map[string]*model.SlotsStat)
	rows, err := d.db.Query(ctx, _getSlotsStatManager)
	if err != nil {
		log.Error("db.Query(%s)  error(%v)", _getSlotsStatManager, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, algorithm, weight string
			slot, state             int
		)
		if err = rows.Scan(&name, &slot, &algorithm, &weight, &state); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		if mapping, ok := slotsMap[name]; ok {
			mapping.Slots = append(mapping.Slots, slot)
		} else {
			slotsMap[name] = &model.SlotsStat{
				Name:      name,
				Slots:     []int{slot},
				Algorithm: algorithm,
				Weight:    weight,
				State:     state,
			}
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
		return
	}
	for _, stat := range slotsMap {
		s = append(s, stat)
	}
	return
}

// CountIdleSlot count idle slot which name="default" and state=1
func (d *Dao) CountIdleSlot(ctx context.Context) (count int, err error) {
	if err = d.db.QueryRow(ctx, _countIdleSlot).Scan(&count); err != nil {
		log.Error("db.QueryRow() error(%v)", err)
	}
	return
}

// IdleSlots get idle slots
func (d *Dao) IdleSlots(ctx context.Context, count int) (slots []int64, err error) {
	rows, err := d.db.Query(ctx, _getIdelSlots, count)
	if err != nil {
		log.Error("db.Query(%s) args(%d) error(%v)", _getIdelSlots, count, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var slot int64
		if err = rows.Scan(&slot); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		slots = append(slots, slot)
	}
	// 槽位不够新创建实验组
	if len(slots) < count {
		slots = nil
		err = errors.New("out of slot")
		return
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// ModifyState ModifyState
func (d *Dao) ModifyState(ctx context.Context, name string, state int) (err error) {
	if _, err = d.db.Exec(ctx, _modifyState, state, name); err != nil {
		log.Error("db.Exec(%s) args(%d, %s) error(%v)", _modifyState, state, name, err)
	}
	return
}

// UpdateSlotsStat UpdateSlotStat and set state inactive.
func (d *Dao) UpdateSlotsStat(ctx context.Context, name, algorithm, weight string, slots []int64, state int) (err error) {
	if _, err = d.db.Exec(ctx, fmt.Sprintf(_setSlot, xstr.JoinInts(slots)), name, algorithm, weight, state); err != nil {
		log.Error("db.Exec() error(%v)", err)
	}
	return
}

// SlotsStatByName get slots stat by name.
func (d *Dao) SlotsStatByName(ctx context.Context, name string) (slots []int64, algorithm, weight string, err error) {
	rows, err := d.db.Query(ctx, _getSlotsStatByName, name)
	if err != nil {
		log.Error("db.Query(%s) args(%s) error(%v)", _getSlotsStatByName, name, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			slot int64
		)
		if err = rows.Scan(&slot, &algorithm, &weight); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		slots = append(slots, slot)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UpdateWeight update a test set weight by name and algorithm.
func (d *Dao) UpdateWeight(ctx context.Context, name string, weight interface{}) (err error) {
	b, err := json.Marshal(weight)
	if err != nil {
		return
	}
	if _, err = d.db.Exec(ctx, _setWeight, string(b), name); err != nil {
		log.Error("db.Exec(%s), error(%v)", _setWeight, err)
	}
	return
}

/*
Statistics
*/

// UpsertStatistics insert or update statistics from database, if err != nil, retry
func (d *Dao) UpsertStatistics(ctx context.Context, name string, date int, hour int, s *model.StatisticsStat) (err error) {
	if _, err = d.db.Exec(ctx, _upsertLog,
		name, date, hour,
		s.View, s.HotClick, s.HotView, s.TotalView,
		s.View, s.HotClick, s.HotView, s.TotalView); err != nil {
		return
	}
	return
}

// StatisticsByDate StatisticsByDate
func (d *Dao) StatisticsByDate(ctx context.Context, begin, end int64) (stats model.StatisticsStats, err error) {
	rows, err := d.db.Query(ctx, _getStatisticsDate, begin, end)
	if err != nil {
		log.Error("db.Query(%s) args(%d, %d) error(%v)", _getStatisticsDate, begin, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s = new(model.StatisticsStat)
		err = rows.Scan(&s.Name, &s.Date, &s.Hour, &s.View, &s.TotalView, &s.HotView, &s.HotClick, &s.HotLike, &s.HotHate, &s.HotChildReply,
			&s.HotReport, &s.TotalLike, &s.TotalHate, &s.TotalReport, &s.TotalRootReply, &s.TotalChildReply,
			&s.HotLikeUV, &s.HotHateUV, &s.HotReportUV, &s.HotChildUV, &s.TotalLikeUV, &s.TotalHateUV, &s.TotalReportUV, &s.TotalChildUV, &s.TotalRootUV)
		if err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		stats = append(stats, s)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
		return
	}
	return
}
