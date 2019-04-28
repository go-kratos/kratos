package reply

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/reply/model/reply"
	"go-common/app/interface/main/reply/model/xreply"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_replySharding int64 = 200
)

const (
	_selByIdsSQL             = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,attr,ctime,mtime FROM reply_%d WHERE id IN (%s)"
	_selSQL                  = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,hate,floor,state,attr,ctime,mtime FROM reply_%d WHERE id=?"
	_selByDialogSQL          = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state IN (0,1,2,5,6) AND dialog=? ORDER BY floor ASC limit ?,?"
	_selByDialogDescSQL      = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state IN (0,1,2,5,6) AND dialog=? AND floor<=? ORDER BY floor DESC limit ?"
	_selByDialogAscSQL       = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state IN (0,1,2,5,6) AND dialog=? AND floor>=? ORDER BY floor ASC limit ?"
	_selMinMaxFloorDialogSQL = "SELECT IFNULL(MIN(floor),0), IFNULL(MAX(floor),0) FROM reply_%d WHERE dialog=? AND oid=? AND type=? AND root=? AND state IN (0,1,2,5,6)"
	_selIdsByFloorSQL        = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) ORDER BY floor DESC limit ?,?"
	_selIdsByCountSQL        = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) ORDER BY rcount DESC limit ?,?"
	_selIdsByLikeSQL         = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) ORDER BY `like` DESC limit ?,?"       // like >= 3
	_selCountByLikeSQL       = "SELECT COUNT(*) FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,5,6) AND attr&1!=1 AND attr>>1&1!=1" // like >= 3

	_selIdsByRootStateSQL   = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state in (0,1,2,5,6) ORDER BY floor limit ?,?"
	_selIdsByRootSQL        = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? ORDER BY floor limit ?,?"
	_selIdsByFloorOffsetSQL = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=0 AND state in (0,1,2,3,4,5,6,7,8,9,10,11) AND floor >= ? AND floor <= ? ORDER BY floor"
	_setHateSQL             = "UPDATE reply_%d SET `hate`=?,mtime=? WHERE id=?"
	_setLikeSQL             = "UPDATE reply_%d SET `like`=?,mtime=? WHERE id=?"

	_insFilteredReply = "INSERT INTO reply_filtered (rpid, oid, mid, type, message, level, ctime, mtime) VALUES (?,?,?,?,?,?,?,?)"

	_selFilteredReply = "select rpid,message FROM reply_filtered WHERE rpid in (%s)"

	// 针对折叠评论
	_foldedReplyLatest = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=12 ORDER BY floor DESC LIMIT ?"
	_foldedReplyDesc   = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=12 AND floor<? ORDER BY floor DESC LIMIT ?"
	_foldedReplyAsc    = "SELECT id FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=12 AND floor>? ORDER BY floor DESC LIMIT ?"
)

// RpDao reply dao.
type RpDao struct {
	db      *sql.DB
	dbSlave *sql.DB
}

// NewReplyDao new replyDao and return.
func NewReplyDao(db, dbSlave *sql.DB) (dao *RpDao) {
	dao = &RpDao{
		db:      db,
		dbSlave: dbSlave,
	}
	return
}

func (dao *RpDao) hit(oid int64) int64 {
	return oid % _replySharding
}

// FilterContents get filtered contents from db
func (d *Dao) FilterContents(ctx context.Context, rpMaps map[int64]string) error {
	if len(rpMaps) == 0 {
		return nil
	}
	var rpids []int64
	for k := range rpMaps {
		rpids = append(rpids, k)
	}
	rows, err := d.dbSlave.Query(ctx, fmt.Sprintf(_selFilteredReply, xstr.JoinInts(rpids)))
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var message string
		if err = rows.Scan(&id, &message); err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				log.Error("row.Scan error(%v)", err)
				return err
			}
		}
		rpMaps[id] = message
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return err
	}
	return nil
}

// CountLike get count of reply like
func (dao *RpDao) CountLike(c context.Context, oid int64, tp int8) (count int, err error) {
	row := dao.dbSlave.QueryRow(c, fmt.Sprintf(_selCountByLikeSQL, dao.hit(oid)), oid, tp)
	if err = row.Scan(&count); err != nil {
		log.Error("row.scan err(%v)", err)
		return
	}
	return
}

// Get get reply.
func (dao *RpDao) Get(c context.Context, oid, rpID int64) (r *reply.Reply, err error) {
	r = &reply.Reply{}
	row := dao.db.QueryRow(c, fmt.Sprintf(_selSQL, dao.hit(oid)), rpID)
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

// GetIDsByDialogAsc ...
func (dao *RpDao) GetIDsByDialogAsc(c context.Context, oid int64, tp int8, root, dialog, maxFloor int64, count int) (rpIDs []int64, err error) {
	query := fmt.Sprintf(_selByDialogAscSQL, dao.hit(oid))
	rows, err := dao.dbSlave.Query(c, query, oid, tp, root, dialog, maxFloor, count)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rpIDs = append(rpIDs, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetDialogMinMaxFloor ...
func (dao *RpDao) GetDialogMinMaxFloor(c context.Context, oid int64, tp int8, root, dialog int64) (minFloor, maxFloor int, err error) {
	query := fmt.Sprintf(_selMinMaxFloorDialogSQL, dao.hit(oid))
	err = dao.dbSlave.QueryRow(c, query, dialog, oid, tp, root).Scan(&minFloor, &maxFloor)
	if err != nil {
		log.Error("row.scan err(%v)", err)
		return
	}
	return
}

// GetIDsByDialogDesc ...
func (dao *RpDao) GetIDsByDialogDesc(c context.Context, oid int64, tp int8, root, dialog, minFloor int64, count int) (rpIDs []int64, err error) {
	query := fmt.Sprintf(_selByDialogDescSQL, dao.hit(oid))
	rows, err := dao.dbSlave.Query(c, query, oid, tp, root, dialog, minFloor, count)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rpIDs = append(rpIDs, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetIDsByDialog get replies by dialog
func (dao *RpDao) GetIDsByDialog(c context.Context, oid int64, tp int8, root, dialog int64, offset, count int) (rpIDs []int64, err error) {
	query := fmt.Sprintf(_selByDialogSQL, dao.hit(oid))
	rows, err := dao.dbSlave.Query(c, query, oid, tp, root, dialog, offset, count)
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		rpIDs = append(rpIDs, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetByIds get replies by reply ids.
func (dao *RpDao) GetByIds(c context.Context, oid int64, tp int8, rpIds []int64) (rpMap map[int64]*reply.Reply, err error) {
	if len(rpIds) == 0 {
		return
	}
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selByIdsSQL, dao.hit(oid), xstr.JoinInts(rpIds)))
	if err != nil {
		log.Error("mysql.Query error(%v)", err)
		return
	}
	defer rows.Close()
	rpMap = make(map[int64]*reply.Reply, len(rpIds))
	for rows.Next() {
		r := &reply.Reply{}
		if err = rows.Scan(&r.RpID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Hate, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
			if err == sql.ErrNoRows {
				r = nil
			} else {
				log.Error("row.Scan error(%v)", err)
				return
			}
		}
		rpMap[r.RpID] = r
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetIdsSortFloor limit get reply ids and order by floor desc.
func (dao *RpDao) GetIdsSortFloor(c context.Context, oid int64, tp int8, offset, count int) (res []int64, err error) {
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selIdsByFloorSQL, dao.hit(oid)), oid, tp, offset, count)
	if err != nil {
		log.Error("dao.selIdsByFloorStmt query err(%v)", err)
		return
	}
	defer rows.Close()
	var id int64
	res = make([]int64, 0, count)
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.scan err is (%v)", err)
			return
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetIdsSortCount limit get reply ids and order by rcount desc.
func (dao *RpDao) GetIdsSortCount(c context.Context, oid int64, tp int8, offset, count int) (res []int64, err error) {
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selIdsByCountSQL, dao.hit(oid)), oid, tp, offset, count)
	if err != nil {
		log.Error("dao.selIdsByCountStmt query err(%v)", err)
		return
	}
	defer rows.Close()
	var id int64
	res = make([]int64, 0, count)
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.scan err is (%v)", err)
			return
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetIdsSortLike limit get reply ids and order by like desc.
func (dao *RpDao) GetIdsSortLike(c context.Context, oid int64, tp int8, offset, count int) (res []int64, err error) {
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selIdsByLikeSQL, dao.hit(oid)), oid, tp, offset, count)
	if err != nil {
		log.Error(" dao.selIdsByLikeStmt query err(%v)", err)
		return
	}
	defer rows.Close()
	var id int64
	res = make([]int64, 0, count)
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.scan err is (%v)", err)
			return
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetIdsByRoot limit get reply ids of root reply and order by floor.
func (dao *RpDao) GetIdsByRoot(c context.Context, oid, root int64, tp int8, offset, count int) (res []int64, err error) {
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selIdsByRootStateSQL, dao.hit(oid)), oid, tp, root, offset, count)
	if err != nil {
		log.Error("dao.selIdsByRtStmt,oid(%d),root(%d),tp(%d),offset(%d),query err(%v)", oid, root, tp, offset, err)
		return
	}
	defer rows.Close()
	var id int64
	res = make([]int64, 0, count)
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.scan err is (%v,%v,%v,%v,%v)", oid, root, offset, count, err)
			return
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetIDsByRootWithoutState limit get reply ids of root reply and order by floor.
func (dao *RpDao) GetIDsByRootWithoutState(c context.Context, oid, root int64, tp int8, offset, count int) (res []int64, err error) {
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selIdsByRootSQL, dao.hit(oid)), oid, tp, root, offset, count)
	if err != nil {
		log.Error("dao.selIdsByRtStmt,oid(%d),root(%d),tp(%d),offset(%d),query err(%v)", oid, root, tp, offset, err)
		return
	}
	defer rows.Close()
	var id int64
	res = make([]int64, 0, count)
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.scan err is (%v,%v,%v,%v,%v)", oid, root, offset, count, err)
			return
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// GetIDsByFloorOffset  get reply ids in floor range
func (dao *RpDao) GetIDsByFloorOffset(c context.Context, oid int64, tp int8, start, end int) (res []int64, err error) {
	rows, err := dao.dbSlave.Query(c, fmt.Sprintf(_selIdsByFloorOffsetSQL, dao.hit(oid)), oid, tp, start, end)
	if err != nil {
		log.Error(" dao.db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	var id int64
	res = make([]int64, 0, end-start)
	for rows.Next() {
		if err = rows.Scan(&id); err != nil {
			log.Error("rows.scan err is (%v)", err)
			return
		}
		res = append(res, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.err error(%v)", err)
		return
	}
	return
}

// SetHateCount set a reply hate count.
func (dao *RpDao) SetHateCount(c context.Context, oid, rpID int64, count int32, now time.Time) (rows int64, err error) {
	res, err := dao.db.Exec(c, fmt.Sprintf(_setHateSQL, dao.hit(oid)), count, now, rpID)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// SetLikeCount set a reply hate count.
func (dao *RpDao) SetLikeCount(c context.Context, oid, rpID int64, count int32, now time.Time) (rows int64, err error) {
	res, err := dao.db.Exec(c, fmt.Sprintf(_setLikeSQL, dao.hit(oid)), count, now, rpID)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// AddFilteredReply AddFilteredReply
func (dao *RpDao) AddFilteredReply(c context.Context, rpID, oid, mid int64, tp, level int8, message string, now time.Time) (err error) {
	_, err = dao.db.Exec(c, _insFilteredReply, rpID, oid, mid, tp, message, level, now, now)
	if err != nil {
		log.Error("db.Exec error(%v)", err)
		return
	}
	return
}

// FoldedRepliesCursor ...
func (dao *RpDao) FoldedRepliesCursor(c context.Context, oid int64, tp int8, root int64, cursor *xreply.Cursor) (ids []int64, err error) {
	var (
		rows *sql.Rows
	)
	if cursor.Latest() {
		rows, err = dao.dbSlave.Query(c, fmt.Sprintf(_foldedReplyLatest, dao.hit(oid)), oid, tp, root, cursor.Ps)
	} else if cursor.Forward() {
		rows, err = dao.dbSlave.Query(c, fmt.Sprintf(_foldedReplyDesc, dao.hit(oid)), oid, tp, root, cursor.Next, cursor.Ps)
	} else {
		rows, err = dao.dbSlave.Query(c, fmt.Sprintf(_foldedReplyAsc, dao.hit(oid)), oid, tp, root, cursor.Prev, cursor.Ps)
		return
	}
	if err != nil {
		log.Error("mysql Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Error("mysql Scan() error(%v)", err)
			return
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		log.Error("mysql rows.Err() error(%v)", err)
	}
	return
}
