package dao

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/tsuna/gohbase/hrpc"
	"go-common/library/log"
)

var (
	hbaseTable    = "ugc:filtercontent"
	hbaseFamilyCt = []byte("ct")
	hbaseColumnCt = []byte("ct")
)

func rowKeyContent(id int64, typ string) string {
	return fmt.Sprintf("%d_%s", id, typ)
}

// Content return the content by id and type.
func (d *Dao) Content(c context.Context, id int64, typ string) (content string, err error) {
	var (
		result      *hrpc.Result
		key         = rowKeyContent(id, typ)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.conf.HBase.ReadTimeout))
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, hbaseTable, key); err != nil {
		log.Error("d.hbase.Get error(%v)", err)
		return
	}
	for _, c := range result.Cells {
		if c != nil && bytes.Equal(c.Family, hbaseFamilyCt) && bytes.Equal(c.Qualifier, hbaseColumnCt) {
			content = string(c.Value)
			break
		}
	}
	return
}
