package dao

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/app/job/main/figure/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// DanmakuReport .
func (d *Dao) DanmakuReport(c context.Context, mid int64, column string, incr int64) (err error) {
	var (
		key         = d.rowVerKey(mid, time.Now())
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	log.Info("Put danmaku act key [%s] c [%s] incr [%d]", key, column, incr)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, incr)
	values := map[string]map[string][]byte{model.ACFamilyUser: map[string][]byte{column: bytesBuffer.Bytes()}}
	if _, err = d.hbase.Increment(ctx, model.ActionCounterTable, key, values); err != nil {
		err = errors.Wrapf(err, "msg(%d,%s,%d), hbase.Increment(key: %s values: %v)", mid, column, incr, key, values)
	}
	return
}
