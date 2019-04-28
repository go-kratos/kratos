package dao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/reply/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_inSQL                  = "INSERT IGNORE INTO reply_%d (id,oid,type,mid,root,parent,floor,state,attr,ctime,mtime) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	_selReplySQL            = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,attr,ctime,mtime FROM reply_%d WHERE id=?"
	_txSelReplySQLForUpdate = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,attr,ctime,mtime FROM reply_%d WHERE id=? for update"
	_selRepliesSQL          = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,floor,state,attr,ctime,mtime FROM reply_%d WHERE id IN (%s)"
	_incrRCountSQL          = "UPDATE reply_%d SET rcount=rcount+1,mtime=? WHERE id=?"
	_upReplyStateSQL        = "UPDATE reply_%d SET state=?,mtime=? WHERE id=?"
	_decrCntSQL             = "UPDATE reply_%d SET rcount=rcount-1,mtime=? WHERE id=? AND rcount > 0"
	_upAttrSQL              = "UPDATE reply_%d SET attr=?,mtime=? WHERE id=?"
	_selExportRepliesSQL    = "SELECT id,oid,type,mid,root,parent,count,rcount,`like`,hate,floor,state,attr,T1.ctime,message from reply_%d as T1 inner join reply_content_%d as T2 on id=rpid where oid=? and type=? and T1.ctime>=? and T2.ctime<=?"
)

// InsertReply insert reply by transaction.
func (d *Dao) InsertReply(c context.Context, r *model.Reply) (id int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_inSQL, hit(r.Oid)), r.ID, r.Oid, r.Type, r.Mid, r.Root, r.Parent, r.Floor, r.State, r.Attr, r.CTime, r.MTime)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// Reply get a reply from database.
func (d *Dao) Reply(c context.Context, oid, rpID int64) (r *model.Reply, err error) {
	r = new(model.Reply)
	row := d.db.QueryRow(c, fmt.Sprintf(_selReplySQL, hit(oid)), rpID)
	if err = row.Scan(&r.ID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
		if err == xsql.ErrNoRows {
			r = nil
			err = nil
		}
	}
	return
}

// Replies get replies by reply ids.
func (d *Dao) Replies(c context.Context, oids []int64, rpIds []int64) (rpMap map[int64]*model.Reply, err error) {
	hitMap := make(map[int64][]int64)
	for i, oid := range oids {
		hitMap[hit(oid)] = append(hitMap[hit(oid)], rpIds[i])
	}
	rpMap = make(map[int64]*model.Reply, len(rpIds))
	for hit, ids := range hitMap {
		var rows *xsql.Rows
		rows, err = d.db.Query(c, fmt.Sprintf(_selRepliesSQL, hit, xstr.JoinInts(ids)))
		if err != nil {
			return
		}
		for rows.Next() {
			r := &model.Reply{}
			if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
				rows.Close()
				return
			}
			rpMap[r.ID] = r
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return
		}
		rows.Close()
	}
	return
}

// UpdateReplyState update reply state.
func (d *Dao) UpdateReplyState(c context.Context, oid, rpID int64, state int32) (rows int64, err error) {
	now := time.Now()
	res, err := d.db.Exec(c, fmt.Sprintf(_upReplyStateSQL, hit(oid)), state, now, rpID)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxUpdateReplyState tx update reply state.
func (d *Dao) TxUpdateReplyState(tx *xsql.Tx, oid, rpID int64, state int32, mtime time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upReplyStateSQL, hit(oid)), state, mtime, rpID)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxIncrReplyRCount incr rcount of reply by transaction.
func (d *Dao) TxIncrReplyRCount(tx *xsql.Tx, oid, rpID int64, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrRCountSQL, hit(oid)), now, rpID)
	if err != nil {
		return
	}
	return res.RowsAffected()
}

// TxDecrReplyRCount decr rcount of reply by transaction.
func (d *Dao) TxDecrReplyRCount(tx *xsql.Tx, oid, rpID int64, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrCntSQL, hit(oid)), now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxUpReplyAttr update subject attr.
func (d *Dao) TxUpReplyAttr(tx *xsql.Tx, oid, rpID int64, attr uint32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upAttrSQL, hit(oid)), attr, now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxReplyForUpdate get a reply from database.
func (d *Dao) TxReplyForUpdate(tx *xsql.Tx, oid, rpID int64) (r *model.Reply, err error) {
	r = new(model.Reply)
	row := tx.QueryRow(fmt.Sprintf(_txSelReplySQLForUpdate, hit(oid)), rpID)
	if err = row.Scan(&r.ID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
		if err == xsql.ErrNoRows {
			r = nil
			err = nil
		}
	}
	return
}

// TxReply get a reply from database.
func (d *Dao) TxReply(tx *xsql.Tx, oid, rpID int64) (r *model.Reply, err error) {
	r = new(model.Reply)
	row := tx.QueryRow(fmt.Sprintf(_selReplySQL, hit(oid)), rpID)
	if err = row.Scan(&r.ID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
		if err == xsql.ErrNoRows {
			r = nil
			err = nil
		}
	}
	return
}

// ExportReplies export replies
func (d *Dao) ExportReplies(c context.Context, oid, mid int64, tp int8, state string, startTime, endTime time.Time) (data [][]string, err error) {
	var rows *xsql.Rows
	title := []string{"id", "oid", "type", "mid", "root", "parent", "count", "rcount", "like", "hate", "floor", "state", "attr", "ctime", "message"}
	data = append(data, title)
	query := fmt.Sprintf(_selExportRepliesSQL, hit(oid), hit(oid))
	if state != "" {
		query += fmt.Sprintf(" and state in (%s)", state)
	}
	if mid != 0 {
		query += fmt.Sprintf(" and mid=%d", mid)
	}
	rows, err = d.dbSlave.Query(c, query, oid, tp, startTime, endTime)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &model.ExportedReply{}
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.Message); err != nil {
			rows.Close()
			return
		}
		data = append(data, r.String())
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
