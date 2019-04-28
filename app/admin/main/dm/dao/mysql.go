package dao

import (
	"context"
	"fmt"
	"sync"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_pageSize        = 1000
	_subjectSharding = 100
	_indexSharding   = 1000
	_contentSharding = 1000

	// subject
	_updateSubAttr        = "UPDATE dm_subject_%02d SET attr=? WHERE type=? AND oid=?"
	_updateSubCount       = "UPDATE dm_subject_%02d SET acount=?,count=? WHERE type=? AND oid=?"
	_updateSubMCountSQL   = "UPDATE dm_subject_%02d SET mcount=? WHERE type=? AND oid=?"
	_incrSubMoveCntSQL    = "UPDATE dm_subject_%02d SET move_count=move_count+? WHERE type=? AND oid=?"
	_updateSubPoolSQL     = "UPDATE dm_subject_%02d SET childpool=? WHERE type=? AND oid=?"
	_updateSubStateSQL    = "UPDATE dm_subject_%02d SET state=? WHERE type=? AND oid=?"
	_updateSubMaxlimitSQL = "UPDATE dm_subject_%02d SET maxlimit=? WHERE type=? AND oid=?"
	_incrSubCountSQL      = "UPDATE dm_subject_%02d SET count=count+? WHERE type=? AND oid=?"
	_selectDMCount        = "SELECT count(*) from dm_index_%03d WHERE type=? AND oid=? AND state IN(%s)"
	_subSQL               = "SELECT id,type,oid,pid,mid,state,attr,acount,count,mcount,move_count,maxlimit,childpool,ctime,mtime FROM dm_subject_%02d WHERE type=? AND oid=?"
	_getSubjectsSQL       = "SELECT id,type,oid,pid,mid,state,attr,acount,count,mcount,move_count,maxlimit,childpool,ctime,mtime FROM dm_subject_%02d WHERE type=? AND oid IN(%s)"
	_contentsSQL          = "SELECT dmid,fontsize,color,mode,ip,plat,msg,ctime,mtime FROM dm_content_%03d WHERE dmid IN(%s)"
	_contentsSpeSQL       = "SELECT dmid,msg,ctime,mtime FROM dm_special_content WHERE dmid IN(%s)"
	_idxsByidSQL          = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND oid=? AND id IN(%s)"
	_updatePoolID         = "UPDATE dm_index_%03d SET pool=? WHERE type=? AND id IN (%s)"
	_updateDMState        = "UPDATE dm_index_%03d SET state=? WHERE type=? AND id IN (%s)"
	_updateDMAttr         = "UPDATE dm_index_%03d SET attr=? WHERE type=? AND id IN (%s)"
)

func (d *Dao) hitSubject(oid int64) int64 {
	return oid % _subjectSharding
}

func (d *Dao) hitIndex(oid int64) int64 {
	return oid % _indexSharding
}

func (d *Dao) hitContent(oid int64) int64 {
	return oid % _contentSharding
}

// UpSubjectAttr update subject attr.
func (d *Dao) UpSubjectAttr(c context.Context, tp int32, oid int64, attr int32) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_updateSubAttr, d.hitSubject(oid)), attr, tp, oid)
	if err != nil {
		log.Error("dmMetaWriter.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpSubjectCount update acount,count of subject.
func (d *Dao) UpSubjectCount(c context.Context, tp int32, oid, acount, count int64) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_updateSubCount, d.hitSubject(oid)), acount, count, tp, oid)
	if err != nil {
		log.Error("dmMetaWriter.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// IncrSubjectCount update count.
func (d *Dao) IncrSubjectCount(c context.Context, tp int32, oid, count int64) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_incrSubCountSQL, d.hitSubject(oid)), count, tp, oid)
	if err != nil {
		log.Error("dmMetaWriter.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// UpSubjectMCount update monitor dm count.
func (d *Dao) UpSubjectMCount(c context.Context, tp int32, oid, cnt int64) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_updateSubMCountSQL, d.hitSubject(oid)), cnt, tp, oid)
	if err != nil {
		log.Error("d.dmMetaWriter.Exect(%s,%d,%d,%d) error(%v)", _updateSubMCountSQL, cnt, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// IncrSubMoveCount update move_count in dm_subject.
func (d *Dao) IncrSubMoveCount(c context.Context, tp int32, oid, count int64) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_incrSubMoveCntSQL, d.hitSubject(oid)), count, tp, oid)
	if err != nil {
		log.Error("d.dmMetaWriter.Exec(%s,%d,%d,%d) error(%v)", _incrSubMoveCntSQL, tp, oid, count, err)
		return
	}
	return res.RowsAffected()
}

// UpSubjectPool update childpool in dm subject.
func (d *Dao) UpSubjectPool(c context.Context, tp int32, oid int64, childpool int32) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_updateSubPoolSQL, d.hitSubject(oid)), childpool, tp, oid)
	if err != nil {
		log.Error("d.dmWriter.Exect(%s,%d,%d,%d) error(%v)", _updateSubPoolSQL, childpool, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpSubjectState update state in dm subject.
func (d *Dao) UpSubjectState(c context.Context, tp int32, oid int64, state int32) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_updateSubStateSQL, d.hitSubject(oid)), state, tp, oid)
	if err != nil {
		log.Error("d.dmWriter.Exect(%s,%d,%d,%d) error(%v)", _updateSubStateSQL, state, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpSubjectMaxlimit update maxlimit in dm subject.
func (d *Dao) UpSubjectMaxlimit(c context.Context, tp int32, oid, maxlimit int64) (affect int64, err error) {
	res, err := d.dmMetaWriter.Exec(c, fmt.Sprintf(_updateSubMaxlimitSQL, d.hitSubject(oid)), maxlimit, tp, oid)
	if err != nil {
		log.Error("d.dmWriter.Exect(%s,%d,%d,%d) error(%v)", _updateSubMaxlimitSQL, maxlimit, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// DMCount statistics dm count by dm state.
func (d *Dao) DMCount(c context.Context, tp int32, oid int64, states []int64) (count int64, err error) {
	row := d.dmMetaReader.QueryRow(c, fmt.Sprintf(_selectDMCount, d.hitIndex(oid), xstr.JoinInts(states)), tp, oid)
	if err = row.Scan(&count); err != nil {
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// Subject get subject info from db.
func (d *Dao) Subject(c context.Context, tp int32, oid int64) (s *model.Subject, err error) {
	s = &model.Subject{}
	row := d.dmMetaReader.QueryRow(c, fmt.Sprintf(_subSQL, d.hitSubject(oid)), tp, oid)
	if err = row.Scan(&s.ID, &s.Type, &s.Oid, &s.Pid, &s.Mid, &s.State, &s.Attr, &s.ACount, &s.Count, &s.MCount, &s.MoveCnt, &s.Maxlimit, &s.Childpool, &s.Ctime, &s.Mtime); err != nil {
		if err == sql.ErrNoRows {
			s = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// Subjects multi get subjects.
func (d *Dao) Subjects(c context.Context, tp int32, oids []int64) (res map[int64]*model.Subject, err error) {
	var (
		oidMap     = make(map[int64][]int64)
		mutext     = &sync.Mutex{}
		wg, errCtx = errgroup.WithContext(c)
	)
	if len(oids) == 0 {
		return
	}
	res = make(map[int64]*model.Subject, len(oids))
	for _, oid := range oids {
		if _, ok := oidMap[d.hitSubject(oid)]; !ok {
			oidMap[d.hitSubject(oid)] = make([]int64, 0)
		}
		oidMap[d.hitSubject(oid)] = append(oidMap[d.hitSubject(oid)], oid)
	}
	for key, value := range oidMap {
		k := key
		v := value
		wg.Go(func() (err error) {
			rows, err := d.dmMetaReader.Query(errCtx, fmt.Sprintf(_getSubjectsSQL, k, xstr.JoinInts(v)), tp)
			if err != nil {
				log.Error("dmMetaReader.Query() error(%v)", err)
				return
			}
			defer rows.Close()
			for rows.Next() {
				s := &model.Subject{}
				err = rows.Scan(&s.ID, &s.Type, &s.Oid, &s.Pid, &s.Mid, &s.State, &s.Attr, &s.ACount, &s.Count, &s.MCount, &s.MoveCnt, &s.Maxlimit, &s.Childpool, &s.Ctime, &s.Mtime)
				if err != nil {
					log.Error("rows.Scan() error(%v)", err)
					return
				}
				mutext.Lock()
				res[s.Oid] = s
				mutext.Unlock()
			}
			return
		})
	}
	if err = wg.Wait(); err != nil {
		log.Error("d.Subjects() error(%v)", err)
		return
	}
	if len(res) == 0 {
		res = nil
	}
	return
}

// SetStateByIDs change danmu state in new db
func (d *Dao) SetStateByIDs(c context.Context, tp int32, oid int64, ids []int64, state int32) (affect int64, err error) {
	sqlStr := fmt.Sprintf(_updateDMState, d.hitIndex(oid), xstr.JoinInts(ids))
	res, err := d.dmMetaWriter.Exec(c, sqlStr, state, tp)
	if err != nil {
		log.Error("d.dmMetaWriter.Exec(%s %d) error(%v)", sqlStr, state, err)
		return
	}
	return res.RowsAffected()
}

// SetAttrByIDs set attr by mult dmids in new db
func (d *Dao) SetAttrByIDs(c context.Context, tp int32, oid int64, ids []int64, val int32) (affect int64, err error) {
	sqlStr := fmt.Sprintf(_updateDMAttr, d.hitIndex(oid), xstr.JoinInts(ids))
	res, err := d.dmMetaWriter.Exec(c, sqlStr, val, tp)
	if err != nil {
		log.Error("d.dmMetaWriter.Exec(%s %d) error(%v)", sqlStr, val, err)
		return
	}
	return res.RowsAffected()
}

// SetPoolIDByIDs change danmu poolid
func (d *Dao) SetPoolIDByIDs(c context.Context, tp int32, oid int64, pool int32, dmids []int64) (affect int64, err error) {
	sqlStr := fmt.Sprintf(_updatePoolID, d.hitIndex(oid), xstr.JoinInts(dmids))
	res, err := d.dmMetaWriter.Exec(c, sqlStr, pool, tp)
	if err != nil {
		log.Error("d.dbComment.Exec(%v) error(%v)", sqlStr, err)
		return
	}
	return res.RowsAffected()
}

// IndexsByID get dm index by dmids.
func (d *Dao) IndexsByID(c context.Context, tp int32, oid int64, dmids []int64) (idxMap map[int64]*model.DM, special []int64, err error) {
	query := fmt.Sprintf(_idxsByidSQL, d.hitIndex(oid), xstr.JoinInts(dmids))
	rows, err := d.dmMetaReader.Query(c, query, tp, oid)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	idxMap = make(map[int64]*model.DM, len(dmids))
	for rows.Next() {
		idx := new(model.DM)
		if err = rows.Scan(&idx.ID, &idx.Type, &idx.Oid, &idx.Mid, &idx.Progress, &idx.State, &idx.Pool, &idx.Attr, &idx.Ctime, &idx.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		idxMap[idx.ID] = idx
		if idx.Pool == model.PoolSpecial {
			special = append(special, idx.ID)
		}
	}
	return
}

// Contents multi get dm content by dmids.
func (d *Dao) Contents(c context.Context, oid int64, dmids []int64) (res []*model.Content, err error) {
	var (
		wg   errgroup.Group
		lock sync.Mutex
	)
	pageNum := len(dmids) / _pageSize
	if len(dmids)%_pageSize > 0 {
		pageNum = pageNum + 1
	}
	for i := 0; i < pageNum; i++ {
		start := i * _pageSize
		end := (i + 1) * _pageSize
		if end > len(dmids) {
			end = len(dmids)
		}
		wg.Go(func() (err error) {
			rows, err := d.dmMetaReader.Query(c, fmt.Sprintf(_contentsSQL, d.hitContent(oid), xstr.JoinInts(dmids[start:end])))
			if err != nil {
				log.Error("db.Query(%s) error(%v)", fmt.Sprintf(_contentsSQL, d.hitContent(oid), xstr.JoinInts(dmids)), err)
				return
			}
			defer rows.Close()
			for rows.Next() {
				ct := &model.Content{}
				if err = rows.Scan(&ct.ID, &ct.FontSize, &ct.Color, &ct.Mode, &ct.IP, &ct.Plat, &ct.Msg, &ct.Ctime, &ct.Mtime); err != nil {
					log.Error("rows.Scan() error(%v)", err)
					return
				}
				lock.Lock()
				res = append(res, ct)
				lock.Unlock()
			}
			return
		})
	}
	err = wg.Wait()
	return
}

// SpecialContents multi get special dm content by dmids.
func (d *Dao) SpecialContents(c context.Context, dmids []int64) (res map[int64]*model.ContentSpecial, err error) {
	res = make(map[int64]*model.ContentSpecial, len(dmids))
	rows, err := d.dmMetaReader.Query(c, fmt.Sprintf(_contentsSpeSQL, xstr.JoinInts(dmids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		content := &model.ContentSpecial{}
		if err = rows.Scan(&content.ID, &content.Msg, &content.Ctime, &content.Mtime); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		res[content.ID] = content
	}
	return
}
