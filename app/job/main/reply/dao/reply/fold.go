package reply

import (
	"context"
	"fmt"
	"go-common/app/job/main/reply/model/reply"
	"go-common/library/database/sql"
)

const (
	_foldedReplies      = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,floor,state,attr,ctime,mtime FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=12"
	_countFoldedReplies = "SELECT COUNT(*) FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=12"
)

// TxCountFoldedReplies ...
func (dao *RpDao) TxCountFoldedReplies(tx *sql.Tx, oid int64, tp int8, root int64) (count int, err error) {
	if err = tx.QueryRow(fmt.Sprintf(_countFoldedReplies, dao.hit(oid)), oid, tp, root).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return
}

// FoldedReplies ...
func (dao *RpDao) FoldedReplies(ctx context.Context, oid int64, tp int8, root int64) (rps []*reply.Reply, err error) {
	rows, err := dao.mysql.Query(ctx, fmt.Sprintf(_foldedReplies, dao.hit(oid)), oid, tp, root)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(reply.Reply)
		if err = rows.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
			return
		}
		rps = append(rps, r)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
