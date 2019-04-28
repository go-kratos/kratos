package data

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/admin/main/up/model/datamodel"
	"go-common/app/admin/main/up/util/hbaseutil"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"github.com/tsuna/gohbase/hrpc"
)

// 这里是集群2的数据请求
const (
	// HbaseUpArchiveInfoPrefix archive info表
	HbaseUpArchiveInfoPrefix = "upcrm:up_influence"
	// HbaseUpArchiveTagInfoPrefix tag表
	HbaseUpArchiveTagInfoPrefix = "upcrm:up_archive_tag_info"
	// HbaseUpArchiveTypeInfoPrefix 分区表
	HbaseUpArchiveTypeInfoPrefix = "upcrm:up_archive_type_info"
)

var (
	//ErrInvalidDataType invalid data type
	ErrInvalidDataType = fmt.Errorf("invalid data type")
)

//UpArchiveDataType data type
type UpArchiveDataType int

const (
	//DataTypeDay7 7 day
	DataTypeDay7 UpArchiveDataType = 1
	//DataTypeDay30 30 day
	DataTypeDay30 UpArchiveDataType = 2
	//DataTypeDay90 90 day
	DataTypeDay90 UpArchiveDataType = 3
	//DataTypeDay180 180 day
	DataTypeDay180 UpArchiveDataType = 4
	//DataTypeDayAll accumulated
	DataTypeDayAll UpArchiveDataType = 5
)

var (
	dataType2FamilyMap = map[UpArchiveDataType]string{
		DataTypeDay7:   "d7",
		DataTypeDay30:  "d30",
		DataTypeDay90:  "d90",
		DataTypeDay180: "d180",
		DataTypeDayAll: "dall",
	}
)

//UpArchiveInfo get up archive info
func (d *Dao) UpArchiveInfo(c context.Context, mids []int64, dataType UpArchiveDataType) (result map[int64]*datamodel.UpArchiveData, err error) {
	var now = time.Now()
	var date = now.AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	var tableName = generateTableName(HbaseUpArchiveInfoPrefix, date)

	var family, ok = dataType2FamilyMap[dataType]
	if !ok {
		log.Error("invalid data type, type=%d", dataType)
		err = ErrInvalidDataType
		return
	}
	var group, ctx = errgroup.WithContext(c)
	var lock sync.Mutex
	result = make(map[int64]*datamodel.UpArchiveData)
	for _, mid := range mids {
		var copymid = mid
		group.Go(
			func() error {
				var data = &datamodel.UpArchiveData{}
				var key = hbaseMd5Key(copymid)
				if e := d.getHbaseRowResult(ctx, tableName, key, data, hrpc.Families(map[string][]string{family: nil})); e != nil {
					log.Error("get up archive info fail, mid=%d, err=%v", err)
					return nil
				}
				lock.Lock()
				result[copymid] = data
				lock.Unlock()
				return nil
			})
	}

	if err = group.Wait(); err != nil {
		log.Error("batch get fail, err=%v", err)
		return
	}
	log.Info("get archive info ok, find result count=%d", len(result))
	return
}

//UpArchiveTagInfo get up archive tag info
func (d *Dao) UpArchiveTagInfo(c context.Context, mid int64) (result datamodel.UpArchiveTagData, err error) {
	var now = time.Now()
	var date = now.AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	var tableName = generateTableName(HbaseUpArchiveTagInfoPrefix, date)
	var key = hbaseMd5Key(mid)

	if err = d.getHbaseRowResult(c, tableName, key, &result); err != nil {
		log.Error("get up archive tag info fail, err=%v", err)
	}
	return
}

//UpArchiveTypeInfo get up archive type info
func (d *Dao) UpArchiveTypeInfo(c context.Context, mid int64) (result datamodel.UpArchiveTypeData, err error) {
	var now = time.Now()
	var date = now.AddDate(0, 0, -1).Add(-12 * time.Hour).Format("20060102")
	var tableName = generateTableName(HbaseUpArchiveTypeInfoPrefix, date)
	var key = hbaseMd5Key(mid)
	if err = d.getHbaseRowResult(c, tableName, key, &result); err != nil {
		log.Error("get up archive type info fail, err=%v", err)
	}
	return
}

func (d *Dao) getHbaseRowResult(c context.Context, table, key string, result interface{}, options ...func(hrpc.Call) error) (err error) {
	var ctx, cancel = context.WithTimeout(c, time.Duration(d.c.HBase.ReadTimeout))
	defer cancel()
	res, err := d.hbase.GetStr(ctx, table, key, options...)
	if err != nil {
		log.Error("fail to get data from hbase, table=%s, key=%s, err=%v", table, key, err)
		return
	}

	if len(res.Cells) == 0 {
		log.Warn("no cells get, table=%s, key=%s, err=%v", table, key, err)
		return
	}

	var parser = hbaseutil.Parser{}
	if err = parser.Parse(res.Cells, result); err != nil {
		log.Error("parse data fail, table=%s, key=%s, err=%v", table, key, err)
		return
	}
	return
}

// dont use date string
func generateTableName(prefix, date string) string {
	return prefix
}
