package timemachine

//import (
//	"bytes"
//	"context"
//	"go-common/app/interface/main/activity/model/timemachine"
//	"go-common/library/log"
//	"io"
//
//	"github.com/tsuna/gohbase/hrpc"
//)
//
//func (d *Dao) timemachineScan(c context.Context, startRow, endRow string) (err error) {
//	var scanner hrpc.Scanner
//	if scanner, err = d.hbase.ScanRangeStr(c, _hBaseUpItemTableName, startRow, endRow); err != nil {
//		log.Error("TimemachineScan d.hbase.Scan(%s) error(%v)", _hBaseUpItemTableName, err)
//		return
//	}
//	for {
//		if d.tmProcStop != 0 {
//			err = scanner.Close()
//			return
//		}
//		if result, e := scanner.Next(); e != nil {
//			if e == io.EOF {
//				return
//			}
//			log.Error("TimemachineScan scanner.Next error(%v)", e)
//			continue
//		} else {
//			if result == nil {
//				continue
//			}
//			item := &timemachine.Item{}
//			for _, c := range result.Cells {
//				if c == nil {
//					continue
//				}
//				if !bytes.Equal(c.Family, []byte("m")) {
//					continue
//				}
//				tmFillFields(item, c)
//			}
//			if item.Mid > 0 {
//				for {
//					if e := d.cache.Do(c, func(ctx context.Context) {
//						if e := d.AddCacheTimemachine(ctx, item.Mid, item); e != nil {
//							log.Error("timemachineScand.AddCacheTimemachine(%d) error(%v)", item.Mid, e)
//						}
//					}); e != nil {
//						log.Warn("timemachineScan d.AddCacheTimemachine channel full(%v)", e)
//						d.limiter.Wait(context.Background())
//					} else {
//						break
//					}
//				}
//			}
//		}
//	}
//}
