package dao

import (
	"context"
	"encoding/binary"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/app/job/main/figure/model"

	"github.com/pkg/errors"
)

// BlockedRage .
func (d *Dao) BlockedRage(c context.Context, mid int64, vs int16) (err error) {
	var (
		rageByte    = make([]byte, 8)
		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.HBase.WriteTimeout))
	)
	defer cancel()
	binary.BigEndian.PutUint16(rageByte, uint16(vs))
	values := map[string]map[string][]byte{model.USFamilyUser: map[string][]byte{model.USColumnBlockedRage: rageByte}}
	if _, err = d.hbase.PutStr(ctx, model.UserInfoTable, d.rowKey(mid), values); err != nil {
		err = errors.Wrapf(err, "mid(%v), hbase.Put(key: %s values: %v)", mid, mid, values)
	}
	return
}
