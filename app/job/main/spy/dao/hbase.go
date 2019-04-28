package dao

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/app/job/main/spy/model"
	"go-common/library/log"

	"github.com/pkg/errors"
	"github.com/tsuna/gohbase/hrpc"
)

var (
	// table
	_hbaseTableActive = "active_data"
	// family
	_familyDays  = "activeDays"
	_familyDaysB = []byte(_familyDays)
)

func strRowKey(mid int64) string {
	return fmt.Sprintf("%d", mid)
}

// GetActiveData get user active days and watched bangumi video time
func (dao *Dao) GetActiveData(c context.Context, mid int64) (active *model.Active, err error) {
	var (
		result      *hrpc.Result
		key         = strRowKey(mid)
		ctx, cancel = context.WithTimeout(c, time.Duration(dao.c.HBase.ReadTimeout))
	)
	defer cancel()
	if result, err = dao.hbase.GetStr(ctx, _hbaseTableActive, key); err != nil {
		err = errors.Wrapf(err, "hbase.GetStr(%s,%s)", _hbaseTableActive, key)
		return
	}
	active = &model.Active{}
	for _, c := range result.Cells {
		h := &model.Active{}
		if c != nil && bytes.Equal(c.Family, _familyDaysB) {
			days, err := strconv.ParseInt(string(c.Qualifier), 10, 64)
			if err != nil {
				log.Error("strconv.ParseInt err(%v)", err)
				continue
			}
			h.Active = days
			active = h
		}
	}
	return
}
