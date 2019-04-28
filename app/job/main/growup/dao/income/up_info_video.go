package income

import (
	"context"
	"fmt"

	model "go-common/app/job/main/growup/model/income"
	"go-common/library/log"
)

const (
	_signedSQL = "SELECT id,mid,signed_at,quit_at,account_state,is_deleted FROM up_info_%s WHERE id > ? ORDER BY id LIMIT ?"
)

// Ups get ups by business type
func (d *Dao) Ups(c context.Context, business string, offset int64, limit int64) (last int64, ups map[int64]*model.Signed, err error) {
	ups = make(map[int64]*model.Signed)
	rows, err := d.db.Query(c, fmt.Sprintf(_signedSQL, business), offset, limit)
	if err != nil {
		log.Error("db Query Signed up Info error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		up := &model.Signed{}
		err = rows.Scan(&last, &up.MID, &up.SignedAt, &up.QuitAt, &up.AccountState, &up.IsDeleted)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		if up.IsDeleted == 0 {
			ups[up.MID] = up
		}
	}
	return
}
