package dao

import (
	"context"
	"encoding/binary"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/app/job/main/figure/model"

	"github.com/pkg/errors"
)

// UpdateVipStatus .
func (d *Dao) UpdateVipStatus(c context.Context, mid int64, vs int32) (err error) {
	var (
		vipByte     = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint64(vipByte, uint64(vs))
	values := map[string]map[string][]byte{model.USFamilyUser: map[string][]byte{model.USColumnVipStatus: vipByte}}
	if _, err = d.hbase.PutStr(ctx, model.UserInfoTable, d.rowKey(mid), values); err != nil {
		err = errors.Wrapf(err, "mid(%v), hbase.Put(key: %s values: %v)", mid, d.rowKey(mid), values)
	}
	return
}
