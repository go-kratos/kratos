package reply

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_replySharding int64 = 200
)

const (
	_inSQL       = "INSERT IGNORE INTO reply_%d (id,oid,type,mid,root,parent,dialog,floor,state,attr,ctime,mtime) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)"
	_incrCntSQL  = "UPDATE reply_%d SET count=count+1,rcount=rcount+1,mtime=? WHERE id=?"
	_incrFCntSQL = "UPDATE reply_%d SET count=count+1,mtime=? WHERE id=?"
	_incrRCntSQL = "UPDATE reply_%d SET rcount=rcount+1,mtime=? WHERE id=?"
	_decrCntSQL  = "UPDATE reply_%d SET rcount=rcount-1,mtime=? WHERE id=? AND rcount > 0"

	_upStateSQL         = "UPDATE reply_%d SET state=?,mtime=? WHERE id=?"
	_upAttrSQL          = "UPDATE reply_%d SET attr=?,mtime=? WHERE id=?"
	_upLikeSQL          = "UPDATE reply_%d SET `like`=?,hate=?,mtime=? WHERE id=?"
	_selSQLForUpdate    = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,attr,ctime,mtime FROM reply_%d WHERE id=? for update"
	_selSQL             = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,attr,ctime,mtime FROM reply_%d WHERE id=?"
	_selAllSQL          = "SELECT id,rcount,`like`,hate,floor,attr FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6)"
	_selAllByRtSQL      = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,ctime,mtime FROM reply_%d WHERE oid=? AND type=? AND root=? AND state in (0,1,2,5,6)"
	_selIncrByDialogSQL = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,attr,ctime,mtime FROM reply_%d WHERE oid=? AND type=? AND root=? AND state IN (0,1,2,5,6) AND dialog=? and id>? limit 10000"
	_selByRootSQL       = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,ctime,mtime FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=?"
	_selTopSQL          = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,floor,state,attr,ctime,mtime FROM reply_%d WHERE oid=? AND type=? AND root=0 AND attr&(1<<?) limit 1"
	_selAllByFloorSQL   = "SELECT id,rcount,`like`,hate,floor,attr FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) and floor>=? and floor<?"
	_selByFloorLimitSQL = "SELECT id,rcount,`like`,hate,floor,attr FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) and floor<? order by floor desc limit ?"
	_selByLikeLimitSQL  = "SELECT id,rcount,`like`,hate,floor,attr FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) order by `like` desc limit ?"
	_selByCountLimitSQL = "SELECT id,rcount,`like`,hate,floor,attr FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) order by `rcount` desc limit ?"

	_fixDialogSelSQL = "select id, parent, floor from reply_%d where oid=? and type=? and root=? and id>? limit ?"
	_fixDialogSetSQL = "update reply_%d set dialog=?, mtime=? where id in (%s)"
)

// RpDao define reply db sqlStmt
type RpDao struct {
	// stmt
	upStateStmt     []*sql.Stmt
	upLikeStmt      []*sql.Stmt
	selStmt         []*sql.Stmt
	selAllStmt      []*sql.Stmt
	selAllByRtStmt  []*sql.Stmt
	selByDialogStmt []*sql.Stmt
	selByRootStmt   []*sql.Stmt
	selTopStmt      []*sql.Stmt
	fixDialogStmt   []*sql.Stmt
	mysql           *sql.DB
}

// NewReplyDao new replyDao and return.
func NewReplyDao(db *sql.DB) (dao *RpDao) {
	dao = &RpDao{
		mysql:           db,
		upStateStmt:     make([]*sql.Stmt, _replySharding),
		upLikeStmt:      make([]*sql.Stmt, _replySharding),
		selTopStmt:      make([]*sql.Stmt, _repSharding),
		selStmt:         make([]*sql.Stmt, _replySharding),
		selAllStmt:      make([]*sql.Stmt, _replySharding),
		selByDialogStmt: make([]*sql.Stmt, _replySharding),
		selAllByRtStmt:  make([]*sql.Stmt, _replySharding),
		selByRootStmt:   make([]*sql.Stmt, _replySharding),
		fixDialogStmt:   make([]*sql.Stmt, _replySharding),
	}

	for i := int64(0); i < _replySharding; i++ {
		dao.upStateStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_upStateSQL, i))
		dao.upLikeStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_upLikeSQL, i))
		dao.selStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_selSQL, i))
		dao.selAllStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_selAllSQL, i))
		dao.selTopStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_selTopSQL, i))
		dao.selByDialogStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_selIncrByDialogSQL, i))
		dao.selAllByRtStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_selAllByRtSQL, i))
		dao.selByRootStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_selByRootSQL, i))
		dao.fixDialogStmt[i] = dao.mysql.Prepared(fmt.Sprintf(_fixDialogSelSQL, i))
	}
	return
}

func (dao *RpDao) hit(oid int64) int64 {
	return oid % _replySharding
}

// TxInsert insert reply by transaction.
func (dao *RpDao) TxInsert(tx *sql.Tx, r *reply.Reply) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_inSQL, dao.hit(r.Oid)), r.RpID, r.Oid, r.Type, r.Mid, r.Root, r.Parent, r.Dialog, r.Floor, r.State, r.Attr, r.CTime, r.MTime)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrCount incr count and rcount of reply by transaction.
func (dao *RpDao) TxIncrCount(tx *sql.Tx, oid, rpID int64, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrCntSQL, dao.hit(oid)), now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrFCount incr rcount of reply by transaction
func (dao *RpDao) TxIncrFCount(tx *sql.Tx, oid, rpID int64, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrFCntSQL, dao.hit(oid)), now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxIncrRCount incr rcount of reply by transaction
func (dao *RpDao) TxIncrRCount(tx *sql.Tx, oid, rpID int64, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrRCntSQL, dao.hit(oid)), now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxDecrCount decr rcount of reply by transaction.
func (dao *RpDao) TxDecrCount(tx *sql.Tx, oid, rpID int64, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_decrCntSQL, dao.hit(oid)), now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// GetForUpdate decr rcount of reply by transaction.
func (dao *RpDao) GetForUpdate(tx *sql.Tx, oid, rpID int64) (r *reply.Reply, err error) {
	r = new(reply.Reply)
	row := tx.QueryRow(fmt.Sprintf(_selSQLForUpdate, dao.hit(oid)), rpID)
	if err = row.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// TxUpState update reply state by transaction.
func (dao *RpDao) TxUpState(tx *sql.Tx, oid, rpID int64, state int8, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upStateSQL, dao.hit(oid)), state, now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpState update reply state.
func (dao *RpDao) UpState(c context.Context, oid, rpID int64, state int8, now time.Time) (rows int64, err error) {
	res, err := dao.upStateStmt[dao.hit(oid)].Exec(c, state, now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxUpAttr update reply state.
func (dao *RpDao) TxUpAttr(tx *sql.Tx, oid, rpID int64, attr uint32, now time.Time) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upAttrSQL, dao.hit(oid)), attr, now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpLike incr or decr reply like.
func (dao *RpDao) UpLike(c context.Context, oid, rpID int64, like, hate int, now time.Time) (rows int64, err error) {
	res, err := dao.upLikeStmt[dao.hit(oid)].Exec(c, like, hate, now, rpID)
	if err != nil {
		log.Error("mysqlDB.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// Get get reply.
func (dao *RpDao) Get(c context.Context, oid, rpID int64) (r *reply.Reply, err error) {
	r = &reply.Reply{}
	row := dao.selStmt[dao.hit(oid)].QueryRow(c, rpID)
	if err = row.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// GetTop get top reply
func (dao *RpDao) GetTop(c context.Context, oid int64, tp int8, bit uint32) (r *reply.Reply, err error) {
	r = &reply.Reply{}
	row := dao.selTopStmt[dao.hit(oid)].QueryRow(c, oid, tp, bit)
	if err = row.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
		if err == sql.ErrNoRows {
			r = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// GetByDialog get all reply by dialog
func (dao *RpDao) GetByDialog(c context.Context, oid int64, typ int8, root int64, dialog int64) (rps []*reply.Reply, err error) {
	var minID int64
	for {
		count := 0
		rows, err := dao.selByDialogStmt[dao.hit(oid)].Query(c, oid, typ, root, dialog, minID)
		if err != nil {
			log.Error("mysql.QueryGetByDialog error(%v)", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			r := &reply.Reply{}
			if err = rows.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
				log.Error("row.Scan() error(%v)", err)
				return nil, err
			}
			rps = append(rps, r)
			count++
			if r.RpID > minID {
				minID = r.RpID
			}
		}
		if err = rows.Err(); err != nil {
			log.Error("mysql rows.Err() error(%v)", err)
			return nil, err
		}
		// 每次查10000个，直到某一次查出来count小于10000
		if count < 10000 {
			break
		}
	}
	return
}

// GetAllInSlice GetAllInSlice
func (dao *RpDao) GetAllInSlice(c context.Context, oid int64, typ int8, maxFloor int, shard int) (rs []*reply.Reply, err error) {
	if shard < 1 {
		log.Error("shard(%d) is too small", shard)
		return nil, fmt.Errorf("shard(%d) is too small", shard)
	}

	start := 1
	startTs := time.Now()
	for {
		if start > maxFloor {
			nowTs := time.Now()
			if nowTs.Sub(startTs) > time.Second*3 {
				log.Warn("GetAllInSlice (%d,%d,%d,%d) running over 3 secs,total time is:%v", oid, typ, maxFloor, shard, nowTs.Sub(startTs))
			}
			return
		}
		end := start + shard
		if end > maxFloor {
			end = maxFloor + 1
		}
		var result []*reply.Reply
		result, err = dao.GetAllByFloor(c, oid, typ, start, end)
		if err != nil {
			//try again
			result, err = dao.GetAllByFloor(c, oid, typ, start, end)
			if err != nil {
				return
			}
		}
		rs = append(rs, result...)
		start += shard
	}
}

// GetByFloorLimit GetByFloorLimit
func (dao *RpDao) GetByFloorLimit(ctx context.Context, oid int64, typ int8, floor int, limit int) (rs []*reply.Reply, err error) {
	lastFloor := floor

	for lastFloor > 1 && limit > 0 {
		count := _maxCount
		if limit <= _maxCount {
			count = limit
		}
		limit -= _maxCount
		var temp []*reply.Reply
		temp, err = dao.getByFloorLimit(ctx, oid, typ, lastFloor, count)
		if err != nil {
			return
		}
		if len(temp) > 0 {
			lastFloor = temp[len(temp)-1].Floor
			rs = append(rs, temp...)
		}
		if len(temp) < count {
			break
		}
	}
	return
}

func (dao *RpDao) getByFloorLimit(ctx context.Context, oid int64, typ int8, floor int, limit int) (rs []*reply.Reply, err error) {
	rows, err := dao.mysql.Query(ctx, fmt.Sprintf(_selByFloorLimitSQL, dao.hit(oid)), oid, typ, floor, limit)
	if err != nil {
		log.Error("mysql.Query %s error(%v)", _selByFloorLimitSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.Attr); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetByLikeLimit GetByLikeLimit
func (dao *RpDao) GetByLikeLimit(ctx context.Context, oid int64, typ int8, limit int) (rs []*reply.Reply, err error) {
	rows, err := dao.mysql.Query(ctx, fmt.Sprintf(_selByLikeLimitSQL, dao.hit(oid)), oid, typ, limit)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.Attr); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetByCountLimit GetByCountLimit
func (dao *RpDao) GetByCountLimit(ctx context.Context, oid int64, typ int8, limit int) (rs []*reply.Reply, err error) {
	rows, err := dao.mysql.Query(ctx, fmt.Sprintf(_selByCountLimitSQL, dao.hit(oid)), oid, typ, limit)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.Attr); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetAllByFloor GetAllByFloor
func (dao *RpDao) GetAllByFloor(ctx context.Context, oid int64, typ int8, start int, end int) (rs []*reply.Reply, err error) {
	rows, err := dao.mysql.Query(ctx, fmt.Sprintf(_selAllByFloorSQL, dao.hit(oid)), oid, typ, start, end)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.Attr); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}

	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetAll get all replies.
func (dao *RpDao) GetAll(c context.Context, oid int64, tp int8) (rs []*reply.Reply, err error) {
	rows, err := dao.selAllStmt[dao.hit(oid)].Query(c, oid, tp)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.Attr); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetAllByRoot get all replies by root reply.
func (dao *RpDao) GetAllByRoot(c context.Context, oid, rpID int64, tp int8) (rs []*reply.Reply, err error) {
	rows, err := dao.selAllByRtStmt[dao.hit(oid)].Query(c, oid, tp, rpID)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetsByRoot get all replies by root reply.
func (dao *RpDao) GetsByRoot(c context.Context, oid, rpID int64, tp, state int8) (rs []*reply.Reply, err error) {
	rows, err := dao.selByRootStmt[dao.hit(oid)].Query(c, oid, tp, rpID, state)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.CTime, &r.MTime); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rs = append(rs, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// FixDialogGetRepliesByRoot ...
func (dao *RpDao) FixDialogGetRepliesByRoot(c context.Context, oid int64, tp int8, rootID int64) (rps []*reply.RpItem, err error) {
	var (
		minID int64
	)
	for {
		count := 0
		rows, err := dao.fixDialogStmt[dao.hit(oid)].Query(c, oid, tp, rootID, minID, 10000)
		if err == sql.ErrNoRows {
			err = nil
			break
		} else if err != nil {
			log.Error("stmt.Query() error(%v)", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			rp := &reply.RpItem{}
			if err = rows.Scan(&rp.ID, &rp.Parent, &rp.Floor); err != nil {
				log.Error("rows.Scan() error(%v)", err)
				return nil, err
			}
			rps = append(rps, rp)
			count++
			if rp.ID > minID {
				minID = rp.ID
			}
		}
		if err = rows.Err(); err != nil {
			log.Error("rows.Err() error(%v)", err)
			return nil, err
		}
		if count < 10000 {
			break
		}
	}
	return
}

// FixDialogSetDialogBatch ...
func (dao *RpDao) FixDialogSetDialogBatch(c context.Context, oid, dialog int64, rpIDs []int64) {
	length := len(rpIDs)
	if length == 0 {
		return
	}
	ids := strings.Trim(strings.Join(strings.Split(fmt.Sprint(rpIDs), " "), ","), "[]")
	setQuery := fmt.Sprintf(_fixDialogSetSQL, dao.hit(oid), ids)
	res, err := dao.mysql.Exec(c, setQuery, dialog, time.Now())
	if err != nil {
		log.Error("db.exec(query: %s, dialog: %d, oid: %d, rpids:%v, error(%v))", setQuery, dialog, oid, rpIDs, err)
		return
	}
	rows, err := res.RowsAffected()
	if rows != int64(length) || err != nil {
		log.Error("s.dao.SetDialogBatch RowsAffected(%d) actual length(%d) error(%v)", rows, length, err)
		return
	}
}
