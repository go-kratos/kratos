package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

const (
	_chunkedSize       = 200
	_getReplyStatsByID = "SELECT id, `like`, hate, rcount, ctime FROM reply_%d WHERE id IN (%s)"
	_getReplyStats     = "SELECT id, `like`, hate, rcount, ctime FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) ORDER BY `like` DESC LIMIT 2000"
	_getReplyReport    = "SELECT rpid, count FROM reply_report_%d WHERE rpid IN (%s)"
	_getSubjectStat    = "SELECT ctime from reply_subject_%d where oid=? and type=?"
	_getRpID           = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) ORDER BY `like` DESC LIMIT 2000"

	_getSlotStats = "SELECT slot, name, algorithm, weight FROM reply_abtest_strategy"

	_getSlotsMapping = "SELECT name, slot FROM reply_abtest_strategy"

	_upsertStatisticsT = "INSERT INTO reply_abtest_statistics (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s"
)

var (
	_upsertStatistics = genSQL()
)

func genSQL() string {
	var (
		slot1 []string
		slot2 []string
		slot3 []string
	)
	slot1 = append(slot1, model.StatisticsDatabaseI...)
	slot1 = append(slot1, model.StatisticsDatabaseU...)
	slot1 = append(slot1, model.StatisticsDatabaseS...)
	for range model.StatisticsDatabaseI {
		slot2 = append(slot2, "?")
	}
	for _, c := range model.StatisticsDatabaseU {
		slot2 = append(slot2, "?")
		slot3 = append(slot3, c+"="+c+"+?")
	}
	for _, c := range model.StatisticsDatabaseS {
		slot2 = append(slot2, "?")
		slot3 = append(slot3, c+"="+"?")
	}
	return fmt.Sprintf(_upsertStatisticsT, strings.Join(slot1, ","), strings.Join(slot2, ","), strings.Join(slot3, ","))
}

func reportHit(oid int64) int64 {
	return oid % 200
}

func replyHit(oid int64) int64 {
	return oid % 200
}

func subjectHit(oid int64) int64 {
	return oid % 50
}

func splitReplyScore(buf []*model.ReplyScore, limit int) [][]*model.ReplyScore {
	var chunk []*model.ReplyScore
	chunks := make([][]*model.ReplyScore, 0, len(buf)/limit+1)
	for len(buf) >= limit {
		chunk, buf = buf[:limit], buf[limit:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

func splitString(buf []string, limit int) [][]string {
	var chunk []string
	chunks := make([][]string, 0, len(buf)/limit+1)
	for len(buf) >= limit {
		chunk, buf = buf[:limit], buf[limit:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

func split(buf []int64, limit int) [][]int64 {
	var chunk []int64
	chunks := make([][]int64, 0, len(buf)/limit+1)
	for len(buf) >= limit {
		chunk, buf = buf[:limit], buf[limit:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

// SlotStats get slot stat
func (d *Dao) SlotStats(ctx context.Context) (ss []*model.SlotStat, err error) {
	rows, err := d.db.Query(ctx, _getSlotStats)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", _getSlotStats, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		s := new(model.SlotStat)
		if err = rows.Scan(&s.Slot, &s.Name, &s.Algorithm, &s.Weight); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		ss = append(ss, s)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// RpIDs return rpIDs should in hot reply list.
func (d *Dao) RpIDs(ctx context.Context, oid int64, tp int) (rpIDs []int64, err error) {
	query := fmt.Sprintf(_getRpID, replyHit(oid))
	rows, err := d.dbSlave.Query(ctx, query, oid, tp)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", query, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var ID int64
		if err = rows.Scan(&ID); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		rpIDs = append(rpIDs, ID)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// ReportStatsByID get report stats from database by ID
func (d *Dao) ReportStatsByID(ctx context.Context, oid int64, rpIDs []int64) (reportMap map[int64]*model.ReplyStat, err error) {
	reportMap = make(map[int64]*model.ReplyStat)
	chunkedRpIDs := split(rpIDs, _chunkedSize)
	for _, ids := range chunkedRpIDs {
		var (
			query = fmt.Sprintf(_getReplyReport, reportHit(oid), xstr.JoinInts(ids))
			rows  *sql.Rows
		)
		rows, err = d.dbSlave.Query(ctx, query)
		if err != nil {
			log.Error("db.Query(%s) error(%v)", query, err)
			return
		}
		for rows.Next() {
			var stat = new(model.ReplyStat)
			if err = rows.Scan(&stat.RpID, &stat.Report); err != nil {
				log.Error("rows.Scan() error(%v)", err)
				return
			}
			reportMap[stat.RpID] = stat
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			log.Error("rows.Err() error(%v)", err)
			return
		}
		rows.Close()
	}
	return
}

// ReplyLHRCStatsByID return a reply like hate reply ctime stat by rpid.
func (d *Dao) ReplyLHRCStatsByID(ctx context.Context, oid int64, rpIDs []int64) (replyMap map[int64]*model.ReplyStat, err error) {
	replyMap = make(map[int64]*model.ReplyStat)
	chunkedRpIDs := split(rpIDs, _chunkedSize)
	for _, ids := range chunkedRpIDs {
		var (
			query = fmt.Sprintf(_getReplyStatsByID, replyHit(oid), xstr.JoinInts(ids))
			rows  *sql.Rows
		)
		rows, err = d.dbSlave.Query(ctx, query)
		if err != nil {
			log.Error("db.Query(%s) error(%v)", query, err)
			return
		}
		for rows.Next() {
			var (
				ctime xtime.Time
				stat  = new(model.ReplyStat)
			)
			if err = rows.Scan(&stat.RpID, &stat.Like, &stat.Hate, &stat.Reply, &ctime); err != nil {
				log.Error("rows.Scan() error(%v)", err)
				return
			}
			stat.ReplyTime = ctime
			replyMap[stat.RpID] = stat
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			log.Error("rows.Err() error(%v)", err)
			return
		}
		rows.Close()
	}
	return
}

// SubjectStats get subject ctime from database
func (d *Dao) SubjectStats(ctx context.Context, oid int64, tp int) (ctime xtime.Time, err error) {
	query := fmt.Sprintf(_getSubjectStat, subjectHit(oid))
	if err = d.dbSlave.QueryRow(ctx, query, oid, tp).Scan(&ctime); err != nil {
		log.Error("db.QueryRow(%s) args(%d, %d) error(%v)", query, oid, tp, err)
		return
	}
	return
}

// ReplyLHRCStats get reply like, hate, reply, ctime stat from database, only get root reply which like>3, call it when back to source.
func (d *Dao) ReplyLHRCStats(ctx context.Context, oid int64, tp int) (replyMap map[int64]*model.ReplyStat, err error) {
	replyMap = make(map[int64]*model.ReplyStat)
	query := fmt.Sprintf(_getReplyStats, replyHit(oid))
	rows, err := d.dbSlave.Query(ctx, query, oid, tp)
	if err != nil {
		log.Error("db.Query(%s) args(%d, %d) error(%v)", query, oid, tp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			ctime xtime.Time
			stat  = new(model.ReplyStat)
		)
		if err = rows.Scan(&stat.RpID, &stat.Like, &stat.Hate, &stat.Reply, &ctime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		stat.ReplyTime = ctime
		replyMap[stat.RpID] = stat
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// SlotsMapping get slots and name mapping.
func (d *Dao) SlotsMapping(ctx context.Context) (slotsMap map[string]*model.SlotsMapping, err error) {
	slotsMap = make(map[string]*model.SlotsMapping)
	rows, err := d.db.Query(ctx, _getSlotsMapping)
	if err != nil {
		log.Error("db.Query(%s) args(%s) error(%v)", _getSlotsMapping, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name string
			slot int
		)
		if err = rows.Scan(&name, &slot); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		slotsMapping, ok := slotsMap[name]
		if ok {
			slotsMapping.Slots = append(slotsMapping.Slots, slot)
		} else {
			slotsMapping = &model.SlotsMapping{
				Name:  name,
				Slots: []int{slot},
			}
		}
		slotsMap[name] = slotsMapping
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UpsertStatistics insert or update statistics into database
func (d *Dao) UpsertStatistics(ctx context.Context, name string, date, hour int, s *model.StatisticsStat) (err error) {
	if _, err = d.db.Exec(ctx, _upsertStatistics,
		name, date, hour,
		s.HotLike, s.HotHate, s.HotReport, s.HotChildReply, s.TotalLike, s.TotalHate, s.TotalReport, s.TotalRootReply, s.TotalChildReply,
		s.HotLikeUV, s.HotHateUV, s.HotReportUV, s.HotChildUV, s.TotalLikeUV, s.TotalHateUV, s.TotalReportUV, s.TotalChildUV, s.TotalRootUV,
		s.HotLike, s.HotHate, s.HotReport, s.HotChildReply, s.TotalLike, s.TotalHate, s.TotalReport, s.TotalRootReply, s.TotalChildReply,
		s.HotLikeUV, s.HotHateUV, s.HotReportUV, s.HotChildUV, s.TotalLikeUV, s.TotalHateUV, s.TotalReportUV, s.TotalChildUV, s.TotalRootUV,
	); err != nil {
		log.Error("upsert statistics failed. error(%v)", err)
		return
	}
	return
}
