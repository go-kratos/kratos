package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"time"

	"go-common/app/service/main/secure/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	tableLog  = "ugc:login_log"
	familyTS  = "ts"
	familyTSB = []byte(familyTS)

	tableSecure = "ugc:secure"
	familyEcpt  = "ecpt"
	familyEcptB = []byte(familyEcpt)
	columnid    = "id"
	columnloc   = "loc"
	columntp    = "tp"
	columnip    = "ip"
	columnts    = "ts"
	// familyFB     = "feedback"
	// familyFBB    = []byte(familyFB)
)

// rowKey return key string.
func rowKey(mid int64) string {
	sum := md5.Sum([]byte(strconv.FormatInt(mid, 10)))
	return fmt.Sprintf("%x", sum)
}
func exKey(mid, ts int64, ip uint32) string {
	return fmt.Sprintf("%d%d_%d_%d", mid%10, mid, ts, ip)
}
func exStartkey(mid int64) string {
	return fmt.Sprintf("%d%d", mid%10, mid)
}
func exStopKey(mid int64) string {
	return fmt.Sprintf("%d%d", mid%10, mid+1)
}

// AddLocs add login log.
func (d *Dao) AddLocs(c context.Context, mid, locid, ts int64) (err error) {
	var (
		locB        = make([]byte, 8)
		key         = rowKey(mid)
		column      = strconv.FormatInt(ts, 10)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint64(locB, uint64(locid))
	values := map[string]map[string][]byte{familyTS: {column: locB}}
	if _, err = d.hbase.PutStr(ctx, tableLog, key, values); err != nil {
		log.Error("hbase.Put error(%v)", err)
	}
	return
}

// Locs get all login location.
func (d *Dao) Locs(c context.Context, mid int64) (locs map[int64]int64, err error) {
	var (
		result      *hrpc.Result
		key         = rowKey(mid)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableLog, key); err != nil {
		log.Error("d.hbase.Get error(%v)", err)
		return
	}
	if result == nil {
		return
	}
	locs = make(map[int64]int64)
	for _, c := range result.Cells {

		if c != nil && len(c.Value) == 8 && bytes.Equal(c.Family, familyTSB) {
			locid := int64(binary.BigEndian.Uint64(c.Value))
			locs[locid]++
		}
	}
	return
}

// AddException add feedback.
func (d *Dao) AddException(c context.Context, l *model.Log) (err error) {
	var (
		idB         = make([]byte, 8)
		ipB         = make([]byte, 8)
		tsB         = make([]byte, 8)
		key         = exKey(l.Mid, int64(l.Time), l.IP)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint64(idB, uint64(l.LocationID))
	binary.BigEndian.PutUint64(ipB, uint64(l.IP))
	binary.BigEndian.PutUint64(tsB, uint64(l.Time))
	values := map[string]map[string][]byte{familyEcpt: {
		columnid:  idB,
		columnloc: []byte(l.Location),
		columnts:  tsB,
		columnip:  ipB,
	}}
	if _, err = d.hbase.PutStr(ctx, tableSecure, key, values); err != nil {
		log.Error("hbase.Put error(%v)", err)
	}
	return
}

// AddFeedBack add feedback
func (d *Dao) AddFeedBack(c context.Context, l *model.Log) (err error) {
	var (
		tpB         = make([]byte, 8)
		key         = exKey(l.Mid, int64(l.Time), l.IP)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint64(tpB, uint64(l.Type))
	values := map[string]map[string][]byte{familyEcpt: {columntp: tpB}}
	if _, err = d.hbase.PutStr(ctx, tableSecure, key, values); err != nil {
		log.Error("hbase.Put error(%v)", err)
	}
	return
}

// ExceptionLoc get exception loc.
func (d *Dao) ExceptionLoc(c context.Context, mid int64) (eps []*model.Expection, err error) {
	var (
		scanner     hrpc.Scanner
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
	)
	defer cancel()
	scanner, err = d.hbase.ScanRangeStr(ctx, tableSecure, exStartkey(mid), exStopKey(mid))
	if err != nil {
		log.Error("d.hbase.ScanRangeStr error(%v)", err)
		return
	}
	for {
		result, err = scanner.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		ep := new(model.Expection)
		for _, c := range result.Cells {
			if c != nil && bytes.Equal(c.Family, familyEcptB) {
				switch string(c.Qualifier) {
				case columnts:
					ep.Time = xtime.Time(binary.BigEndian.Uint64(c.Value))
				case columntp:
					ep.FeedBack = int8(binary.BigEndian.Uint64(c.Value))
				case columnip:
					ep.IP = binary.BigEndian.Uint64(c.Value)
				}
			}
		}
		eps = append(eps, ep)
	}
}
