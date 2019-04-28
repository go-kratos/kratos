package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/member/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_SelAsos = "SELECT mid,uname from aso_account where mid in(%s)"

	_SelAsoName = "SELECT uname from aso_account where mid=?"
)

// Accounts gets account info.
func (d *Dao) Accounts(c context.Context, mids []int64) (accs map[int64]*model.AccountInfo, errs map[int64]map[string]bool, err error) {
	//moralLogTpl := mids[0] % _shardMoralLog
	midsStr := xstr.JoinInts(mids)
	accs = make(map[int64]*model.AccountInfo)
	errs = make(map[int64]map[string]bool)

	// name
	rows, err := d.asodb.Query(c, fmt.Sprintf(_SelAsos, midsStr))
	if err != nil {
		log.Error("d.asodb.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.AccountInfo{}
		if err = rows.Scan(&r.Mid, &r.Name); err != nil {
			log.Error("row.Scan() error(%v) mid(%v)", err, r.Mid)
			err = nil
			continue
		}
		accs[r.Mid] = r
		errs[r.Mid] = make(map[string]bool)
		errs[r.Mid]["asoOK"] = true
	}
	return
}

//Name name.
func (d *Dao) Name(c context.Context, mid int64) (name string, err error) {
	arow := d.asodb.QueryRow(c, _SelAsoName, mid)
	if err = arow.Scan(&name); err != nil {
		log.Error("row.Scan() error(%v)", err)
		return
	}
	return
}
