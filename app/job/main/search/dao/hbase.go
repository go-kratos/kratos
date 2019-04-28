package dao

// import (
// 	"context"
// 	"time"

// 	"go-common/app/job/main/search/model"
// 	"go-common/log"
// 	"golang/gohbase/hrpc"
// )

// // stat archive stat
// func (d *Dao) stat(c context.Context, tableName, startRow, endRow string, from, to uint64, limit int) (res []*model.HbaseArchiveStat, err error) {
// 	var (
// 		scan        *hrpc.Scan
// 		results     []*hrpc.Result
// 		ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadsTimeout))
// 	)
// 	defer cancel()
// 	if scan, err = hrpc.NewScanRangeStr(ctx, tableName, startRow, endRow, from, to); err != nil {
// 		log.Error("d.hbase.stat hrpc.NewScanRangeStr table(%s) startRow(%s) endRow(%s) from(%d) to(%d) error(%v)", tableName, startRow, endRow, from, to, err)
// 		return
// 	}
// 	scan.SetLimit(limit)
// 	if results, err = d.hbase.Scan(ctx, scan); err != nil {
// 		log.Error("d.hbase.Scan error(%v)", err)
// 		return
// 	}
// 	for _, r := range results {
// 		for _, c := range r.Cells {
// 			oneRes := &model.HbaseArchiveStat{
// 				Row:       string(c.Row),
// 				TimeStamp: uint64(*c.Timestamp),
// 				Value:     string(c.Value),
// 			}
// 			res = append(res, oneRes)
// 		}
// 	}
// 	return
// }
