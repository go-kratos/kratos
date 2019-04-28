package archive

import (
	"context"
	"fmt"

	"go-common/library/log"
)

const (
	_consumerOnlineSQL = "SELECT uid FROM task_consumer WHERE uid IN (%s) AND state=1"
)

//ConsumerOnline get online task_consumer
func (d *Dao) ConsumerOnline(c context.Context, uids string) (ids []int64, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_consumerOnlineSQL, uids))
	if err != nil {
		log.Error("d.db.Query(%s, %v) error(%v)", _consumerOnlineSQL, uids, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		ids = append(ids, id)
	}
	return
}
