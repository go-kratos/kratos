package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	tableInfo = "ugc:ArchiveTaskWeight"
	family    = "weightlog"
	familyB   = []byte(family)
)

// hashRowKey create rowkey(md5(tid)[:2]+tid) for track by tid.
func hashRowKey(tid int64) string {
	var bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(tid))
	rk := md5.Sum(bs)
	return fmt.Sprintf("%x%d", rk[:2], tid)
}

// WeightLog get weight log.
func (d *Dao) WeightLog(c context.Context, taskid int64) (ls []*model.TaskWeightLog, err error) {
	var (
		result      *hrpc.Result
		key         = hashRowKey(taskid)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableInfo, key); err != nil {
		log.Error("hrpc.NewGetStr(%s,%s) error(%v)", tableInfo, key, err)
		return
	}
	for _, c := range result.Cells {
		if c == nil || !bytes.Equal(c.Family, familyB) {
			return
		}
		aLog := &model.TaskWeightLog{}
		if err = json.Unmarshal(c.Value, aLog); err != nil {
			log.Warn("json.Unmarshal(%s) error(%v)", string(c.Value), err)
			err = nil
			continue
		}
		ls = append(ls, aLog)
	}
	return
}
