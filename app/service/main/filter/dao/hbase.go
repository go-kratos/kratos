package dao

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/tsuna/gohbase/hrpc"
)

var (
	// _hbaseTable filter table.
	_hbaseTable = "ugc:filtercontent"
	// _hbaseFamilyCt content family.
	_hbaseFamilyCt = []byte("ct")
	// _hbaseColumnCt content column.
	_hbaseColumnCt = []byte("ct")
)

func rowKeyContent(id int64, typ string) string {
	return fmt.Sprintf("%d_%s", id, typ)
}

// SetContent set the content with id and type.
func (d *Dao) SetContent(c context.Context, id int64, typ, content string) (err error) {
	var (
		key         = rowKeyContent(id, typ)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.conf.HBase.WriteTimeout))
	)
	defer cancel()
	valueBytes := []byte(content)
	values := map[string]map[string][]byte{string(_hbaseFamilyCt): {string(_hbaseColumnCt): valueBytes}}
	if _, err = d.hbase.PutStr(ctx, _hbaseTable, key, values); err != nil {
		err = errors.Wrapf(err, "hbase.PutStr(%s,%s)", _hbaseTable, key)
	}
	return
}

// Content return the content by id and type.
func (d *Dao) Content(c context.Context, id int64, typ string) (content string, err error) {
	var (
		result      *hrpc.Result
		key         = rowKeyContent(id, typ)
		ctx, cancel = context.WithTimeout(c, time.Duration(d.conf.HBase.ReadTimeout))
	)
	defer cancel()
	if result, err = d.hbase.GetStr(ctx, _hbaseTable, key); err != nil {
		err = errors.Wrapf(err, "hbase.GetStr(%s,%s)", _hbaseTable, key)
		return
	}
	for _, c := range result.Cells {
		if c != nil && bytes.Equal(c.Family, _hbaseFamilyCt) && bytes.Equal(c.Qualifier, _hbaseColumnCt) {
			content = string(c.Value)
			break
		}
	}
	return
}
