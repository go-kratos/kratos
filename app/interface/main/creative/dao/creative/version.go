package creative

import (
	"context"
	"database/sql"

	"fmt"
	"strconv"
	"strings"

	"go-common/library/log"

	"go-common/app/interface/main/creative/model/version"
)

const (
	_getAllByTypesSQL   = "SELECT id,type,title,content,link,ctime,dateline from version_log where type IN (%s) ORDER BY dateline desc"
	_getLatestByTypeSQL = "SELECT id,type,title,content,link,ctime,dateline from version_log where type = ? ORDER BY dateline desc limit 1"
)

// AllByTypes fn
func (d *Dao) AllByTypes(c context.Context, tyInts []int) (vss []*version.Version, err error) {
	tpStrs := []string{}
	for _, i := range tyInts {
		tpStrs = append(tpStrs, strconv.Itoa(i))
	}
	rows, err := d.creativeDb.Query(c, fmt.Sprintf(_getAllByTypesSQL, strings.Join(tpStrs, ",")))
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		vs := &version.Version{}
		if err = rows.Scan(&vs.ID, &vs.Ty, &vs.Title, &vs.Content, &vs.Link, &vs.Ctime, &vs.Dateline); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		vss = append(vss, vs)
	}
	return
}

// LatestByType fn
func (d *Dao) LatestByType(c context.Context, tid int) (vs *version.Version, err error) {
	vs = &version.Version{}
	row := d.creativeDb.QueryRow(c, _getLatestByTypeSQL, tid)
	if err = row.Scan(&vs.ID, &vs.Ty, &vs.Title, &vs.Content, &vs.Link, &vs.Ctime, &vs.Dateline); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			vs = nil
		}
		log.Error("row.Scan error(%v)", err)
		return
	}
	return
}
