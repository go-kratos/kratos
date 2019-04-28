package dao

import (
	"context"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"

	"go-common/app/admin/main/member/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	tableFaceByMidHistory = "account:user_face"
	tableFaceByOPHistory  = "account:user_face_admin"
)

// reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func rpad(s string, c string, l int) string {
	dt := l - len(s)
	if dt <= 0 {
		return s
	}
	return s + strings.Repeat(c, dt)
}

// func lpad(s string, c string, l int) string {
// 	dt := l - len(s)
// 	if dt <= 0 {
// 		return s
// 	}
// 	return strings.Repeat(c, dt) + s
// }

func midKey(mid int64, time xtime.Time) string {
	suf := rpad(strconv.FormatInt(math.MaxInt64-int64(time), 10), "0", 10)
	return fmt.Sprintf("%s%s", rpad(reverse(strconv.FormatInt(mid, 10)), "0", 10), string(suf[len(suf)-10:]))
}

func operatorKey(operator string, time xtime.Time) string {
	suf := rpad(strconv.FormatInt(math.MaxInt64-int64(time), 10), "0", 10)
	return fmt.Sprintf("%s%s", rpad(operator, "-", 20), string(suf[len(suf)-10:]))
}

func scanTimes(tDuration xtime.Time) int {
	times := int(tDuration)/(60*60) + 5
	if times > 50 {
		times = 50
	}
	return times
}

// FaceHistoryByMid is.
func (d *Dao) FaceHistoryByMid(ctx context.Context, arg *model.ArgFaceHistory) (model.FaceRecordList, error) {
	size := arg.ETime - arg.STime
	stimes := scanTimes(size)
	chunk := xtime.Time(size / xtime.Time(stimes))
	if chunk == 0 {
		chunk = size
	}
	rendKey, rstartKey := midKey(arg.Mid, arg.STime), midKey(arg.Mid, arg.ETime)

	eg := errgroup.Group{}
	lock := sync.RWMutex{}
	records := make(model.FaceRecordList, 0)
	for i := 0; i < stimes; i++ {
		times := i + 1
		stime := arg.STime + (chunk * xtime.Time(i))
		etime := arg.STime + (chunk * xtime.Time(i+1))
		eg.Go(func() error {
			endKey, startKey := midKey(arg.Mid, stime), midKey(arg.Mid, etime)
			log.Info("FaceHistoryByMid: range: start: %s end: %s with times: %d scaning key: start: %s end: %s",
				rstartKey, rendKey, times, startKey, endKey)

			scanner, err := d.fhbymidhbase.ScanRangeStr(ctx, tableFaceByMidHistory, startKey, endKey)
			if err != nil {
				log.Error("hbase.ScanRangeStr(%s,%+v) error(%v)", tableFaceByMidHistory, arg, err)
				return err
			}
			for {
				r, err := scanner.Next()
				if err != nil {
					if err != io.EOF {
						return err
					}
					break
				}
				lock.Lock()
				records = append(records, toMidFaceRecord(r))
				lock.Unlock()
			}
			return nil
		})

		if etime >= arg.ETime {
			break
		}
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return records, nil
}

// FaceHistoryByOP is.
func (d *Dao) FaceHistoryByOP(ctx context.Context, arg *model.ArgFaceHistory) (model.FaceRecordList, error) {
	size := arg.ETime - arg.STime
	stimes := scanTimes(size)
	chunk := xtime.Time(size / xtime.Time(stimes))
	if chunk == 0 {
		chunk = size
	}
	rendKey, rstartKey := operatorKey(arg.Operator, arg.STime), operatorKey(arg.Operator, arg.ETime)

	eg := errgroup.Group{}
	lock := sync.RWMutex{}
	records := make(model.FaceRecordList, 0)
	for i := 0; i < stimes; i++ {
		times := i + 1
		stime := arg.STime + (chunk * xtime.Time(i))
		etime := arg.STime + (chunk * xtime.Time(i+1))
		eg.Go(func() error {
			endKey, startKey := operatorKey(arg.Operator, stime), operatorKey(arg.Operator, etime)
			log.Info("FaceHistoryByOP: range: start: %s end: %s with times: %d scaning key: start: %s end: %s",
				rstartKey, rendKey, times, startKey, endKey)

			scanner, err := d.fhbyophbase.ScanRangeStr(ctx, tableFaceByOPHistory, startKey, endKey)
			if err != nil {
				log.Error("hbase.ScanRangeStr(%s,%+v) error(%v)", tableFaceByOPHistory, arg, err)
				return err
			}

			for {
				r, err := scanner.Next()
				if err != nil {
					if err != io.EOF {
						return err
					}
					break
				}
				lock.Lock()
				records = append(records, toOPFaceRecord(r))
				lock.Unlock()
			}
			return nil
		})

		if etime >= arg.ETime {
			break
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return records, nil
}

func toOPFaceRecord(res *hrpc.Result) *model.FaceRecord {
	l := &model.FaceRecord{}
	for _, c := range res.Cells {
		key := string(c.Row)
		qf := string(c.Qualifier)
		v := string(c.Value)
		log.Info("Retrieving op face record history: key(%s) qualifier(%s) value(%s)", key, qf, v)
		// fill fields
		switch qf {
		case "id":
			l.ID, _ = strconv.ParseInt(v, 10, 64)
		case "mid":
			l.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "sta":
			l.Status = model.ParseStatus(v)
		case "op":
			l.Operator = v
		case "old":
			l.OldFace = v
		case "new":
			l.NewFace = v
		case "apply":
			l.ApplyTime = model.ParseApplyTime(v)
		case "mtm":
			l.ModifyTime, _ = model.ParseLogTime(v)
		}
	}
	return l
}

func toMidFaceRecord(res *hrpc.Result) *model.FaceRecord {
	l := &model.FaceRecord{}
	for _, c := range res.Cells {
		key := string(c.Row)
		qf := string(c.Qualifier)
		v := string(c.Value)
		log.Info("Retrieving mid face record history: key(%s) qualifier(%s) value(%s)", key, qf, v)
		// fill fields
		switch qf {
		// case "id":
		// 	l.ID, _ = strconv.ParseInt(v, 10, 64)
		case "mid":
			l.Mid, _ = strconv.ParseInt(v, 10, 64)
		case "s":
			l.Status = model.ParseStatus(v)
		case "op":
			l.Operator = v
		case "of":
			l.OldFace = v
		case "nf":
			l.NewFace = v
		case "at":
			l.ApplyTime = model.ParseApplyTime(v)
		case "mt":
			l.ModifyTime, _ = model.ParseLogTime(v)
		}
	}
	return l
}
