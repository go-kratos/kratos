package dao

import (
	"bytes"
	"context"
	"fmt"

	"go-common/app/service/main/filter/model"
	xsql "go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_areaKeyRule = "SELECT a.id,b.id,a.area,a.`key`,b.filter,b.mode,b.level,b.comment,b.etime,b.stime FROM filter_key" +
		" AS a INNER JOIN filter_content AS b ON a.filterid=b.id WHERE a.`key`=? AND a.area in(%s) AND a.state=0" +
		" AND b.stime<now() AND b.etime>now()"
)

// KeyAreas .
func (d *Dao) KeyAreas(c context.Context, key string, areas []string) (rs []*model.KeyAreaInfo, err error) {
	var (
		querySQL string
		rows     *xsql.Rows
	)
	if rows, err = d.mysql.Query(c, fmt.Sprintf(_areaKeyRule, d.CoverStr(areas)), key); err != nil {
		err = errors.Wrapf(err, "d.mysql.Query(%s,%s) error(%+v)", querySQL, key, err)
		return
	}
	for rows.Next() {
		var (
			r = &model.KeyAreaInfo{}
		)
		if err = rows.Scan(&r.FKID, &r.ID, &r.Area, &r.Key, &r.Filter, &r.Mode, &r.Level, &r.Comment, &r.ETime, &r.STime); err != nil {
			err = errors.WithStack(err)
			return
		}
		rs = append(rs, r)
	}
	if err = rows.Err(); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// CoverStr return 'strs[0]','strs[1]',...
func (d *Dao) CoverStr(strs []string) string {
	var buf = bytes.NewBuffer(nil)
	for _, str := range strs {
		buf.WriteString("'")
		buf.WriteString(str)
		buf.WriteString("'")
		buf.WriteString(",")
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}
