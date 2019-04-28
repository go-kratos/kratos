package dao

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/relation/model"
	"go-common/library/log"
)

// reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func rpad(s string, l int) string {
	dt := l - len(s)
	if dt <= 0 {
		return s
	}
	return s + strings.Repeat("0", dt)
}

func logKey(mid, fid, ts int64) string {
	midStr := rpad(reverse(strconv.FormatInt(mid, 10)), 10)
	fidStr := rpad(reverse(strconv.FormatInt(fid, 10)), 10)
	tsStr := strconv.FormatInt(ts, 10)
	return midStr + fidStr + tsStr
}

// RelationLogs is used to retriev relation log.
func (d *Dao) RelationLogs(ctx context.Context, mid, fid int64, from time.Time, to time.Time) (model.RelationLogList, error) {
	scanner, err := d.hbase.ScanRangeStr(ctx, d.c.LogTable, logKey(mid, fid, from.Unix()), logKey(mid, fid, to.Unix()))
	if err != nil {
		log.Error("Failed to d.hbase.Scan(): %+v", err)
		return nil, err
	}

	logs := make(model.RelationLogList, 0)

	for {
		r, err := scanner.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		l := &model.RelationLog{
			Mid: mid,
			Fid: fid,
		}
		for _, c := range r.Cells {
			key := string(c.Row)
			qf := string(c.Qualifier)
			v := string(c.Value)
			log.Info("Retrieving relation log: mid(%d) fid(%d) key(%s) qualifier(%s) value(%s)", mid, fid, key, qf, v)

			// fill fields
			switch qf {
			case "att":
				l.Attention = model.ParseAction(v)
			case "bl":
				l.Black = model.ParseAction(v)
			case "wh":
				l.Whisper = model.ParseAction(v)
			case "src":
				l.Source = model.ParseSource(v)
			case "mt":
				l.MTime, _ = model.ParseLogTime(v)
			}
		}
		l.FillAttrField()
		logs = append(logs, l)
	}
	return logs, nil
}
