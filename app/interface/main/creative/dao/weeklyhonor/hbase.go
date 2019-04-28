package weeklyhonor

import (
	"context"
	"strconv"

	"go-common/app/admin/main/up/util/hbaseutil"
	model "go-common/app/interface/main/creative/model/weeklyhonor"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

// reverse for string.
func reverseString(s string) string {
	rs := []rune(s)
	l := len(rs)
	for f, t := 0, l-1; f < t; f, t = f+1, t-1 {
		rs[f], rs[t] = rs[t], rs[f]
	}
	ns := string(rs)
	if l < 10 {
		for i := 0; i < 10-l; i++ {
			ns = ns + "0"
		}
	}
	return ns
}

func honorRowKey(id int64, date string) string {
	idStr := strconv.FormatInt(id, 10)
	s := reverseString(idStr) + date
	return s
}

// HonorStat get up honor.
func (d *Dao) HonorStat(c context.Context, mid int64, date string) (stat *model.HonorStat, err error) {
	var (
		result      *hrpc.Result
		ctx, cancel = context.WithTimeout(c, d.hbaseTimeOut)
		tableName   = "up_honorary_weekly"
		key         = honorRowKey(mid, date)
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, tableName, key); err != nil {
		log.Error("HonorStat d.hbase.GetStr tableName(%s) mid(%d) key(%v) error(%v)", tableName, mid, key, err)
		err = ecode.CreativeDataErr
		return
	}
	if result == nil {
		return
	}
	stat = new(model.HonorStat)
	parser := hbaseutil.Parser{}
	err = parser.Parse(result.Cells, stat)
	if err != nil {
		log.Error("failed to parser hbase, tableName(%s) mid(%d) key(%v) stat(%+v) err (%v)", tableName, mid, key, stat, err)
		return
	}
	log.Info("HonorStat d.hbase.GetStr tableName(%s) mid(%d) key(%v) stat(%+v)", tableName, mid, key, stat)
	return
}
