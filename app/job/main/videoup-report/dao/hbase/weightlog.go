package hbase

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/videoup-report/model/task"
	"go-common/library/log"
)

var (
	tableInfo = "ugc:ArchiveTaskWeight"
	family    = "weightlog"
)

// hashRowKey create rowkey(md5(tid)[:2]+tid) for track by tid.
func hashRowKey(tid int64) string {
	var bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(tid))
	rk := md5.Sum(bs)
	return fmt.Sprintf("%x%d", rk[:2], tid)
}

// AddLog task weight log.
func (d *Dao) AddLog(c context.Context, alog *task.WeightLog) (err error) {
	var (
		value   []byte
		fvalues = make(map[string][]byte)
		column  string
		key     = hashRowKey(alog.TaskID)
	)
	column = strconv.FormatInt(int64(alog.Uptime.TimeValue().Unix()), 10)
	if value, err = json.Marshal(alog); err != nil {
		log.Error("json.Marshal(%v) error(%v)", value, err)
		return
	}
	fvalues[column] = value
	values := map[string]map[string][]byte{family: fvalues}
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.Hbase.WriteTimeout))
	defer cancel()
	// hbase info
	if _, err = d.hbase.PutStr(ctx, tableInfo, key, values); err != nil {
		log.Error("d.hbase.PutStr(%s,%s,%+v) error(%v)", tableInfo, key, values, err)
	}
	return
}
