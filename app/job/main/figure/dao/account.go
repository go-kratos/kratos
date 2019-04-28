package dao

import (
	"context"
	"encoding/binary"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/app/job/main/figure/model"

	"github.com/pkg/errors"
)

// UpdateAccountExp update user exp
func (d *Dao) UpdateAccountExp(c context.Context, mid, exp int64) (err error) {
	var (
		expByte     = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint64(expByte, uint64(exp))
	values := map[string]map[string][]byte{model.USFamilyUser: map[string][]byte{model.USColumnExp: expByte}}
	if _, err = d.hbase.PutStr(ctx, model.UserInfoTable, d.rowKey(mid), values); err != nil {
		err = errors.Wrapf(err, "mid(%d), hbase.Put(key: %s, values: %v)", d.rowKey(mid), values)
	}
	return
}

// IncArchiveViews .
func (d *Dao) IncArchiveViews(c context.Context, mid int64) (err error) {
	var (
		incrByte    = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint64(incrByte, uint64(1))
	values := map[string]map[string][]byte{model.USFamilyUser: map[string][]byte{model.USColumnArchiveViews: incrByte}}
	if _, err = d.hbase.Increment(ctx, model.UserInfoTable, d.rowKey(mid), values); err != nil {
		err = errors.Wrapf(err, "msg(%v), hbase.Increment(key: %s values: %v)", mid, d.rowKey(mid), values)
	}
	return
}
