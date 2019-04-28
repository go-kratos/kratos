package data

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"

	"go-common/library/database/hbase.v2"

	"github.com/tsuna/gohbase/hrpc"
)

func (d *Dao) getDataWithBackup(c context.Context, client *hbase.Client, tableNameFunc func(retryCount int) string, maxRetry int, key string, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	if client == nil {
		err = fmt.Errorf("hbase client is nil")
		return
	}

	for i := 0; i < maxRetry; i++ {
		var tableName = tableNameFunc(i)
		if result, err = d.hbase.GetStr(c, tableName, key, options...); err != nil {
			log.Error("hbase GetStr tableName(%s)|key(%v)|error(%v)", tableName, key, err)
			continue
		}
		break
	}
	return
}

func getTableName(tablePrefix string, date time.Time) string {
	return tablePrefix + date.Format("20060102")
}

func generateTableNameFunc(tablePrefix string, date time.Time, dayDiff int) func(retryCount int) string {
	return func(retryCount int) string {
		var backdate = date.AddDate(0, 0, dayDiff*retryCount)
		return getTableName(tablePrefix, backdate)
	}
}
