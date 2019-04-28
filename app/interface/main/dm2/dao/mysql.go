package dao

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/time"
	"go-common/library/xstr"
)

const (
	_pageSize             = 1000
	_subjectSharding      = 100
	_indexSharding        = 1000
	_contentSharding      = 1000
	_subSQL               = "SELECT id,type,oid,pid,mid,state,attr,acount,count,mcount,move_count,maxlimit,childpool,ctime,mtime FROM dm_subject_%02d WHERE type=? AND oid=?"
	_getSubjectsSQL       = "SELECT id,type,oid,pid,mid,state,attr,acount,count,mcount,move_count,maxlimit,childpool,ctime,mtime FROM dm_subject_%02d WHERE type=? AND oid IN(%s)"
	_updateSubAttr        = "UPDATE dm_subject_%02d SET attr=? WHERE type=? AND oid=?"
	_idxByidSQL           = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND id=?"
	_idxsByidSQL          = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND id IN(%s)"
	_getIndexSQL          = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND oid=? AND state IN(0,6)"
	_getContentSQL        = "SELECT dmid,fontsize,color,mode,ip,plat,msg,ctime,mtime FROM dm_content_%03d WHERE dmid=?"
	_contentsSQL          = "SELECT dmid,fontsize,color,mode,ip,plat,msg,ctime,mtime FROM dm_content_%03d WHERE dmid IN(%s)"
	_contentsSpeSQL       = "SELECT dmid,msg,ctime,mtime FROM dm_special_content WHERE dmid IN(%s)"
	_idxSegIDSQL          = "SELECT id FROM dm_index_%03d WHERE type=? AND oid=? AND progress>=? AND progress<? AND state IN(0,6) AND pool = ? limit ?"
	_selectDMCount        = "SELECT count(*) from dm_index_%03d WHERE type=? AND oid=? AND state IN(%s)"
	_updateIdxStatSQL     = "UPDATE dm_index_%03d SET state=? WHERE type=? AND oid=? AND id IN(%s)"
	_updateUserIdxStatSQL = "UPDATE dm_index_%03d SET state=? WHERE type=? AND oid=? AND id IN(%s) and mid=?"
	_updateIdxPoolSQL     = "UPDATE dm_index_%03d SET pool=? WHERE type=? AND oid=? AND id IN(%s)"
	_updateIdxAttrSQL     = "UPDATE dm_index_%03d SET attr=? WHERE id=?"
	_judgeIdxPageSQL      = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND oid=? AND ctime >=? AND ctime <=? AND progress>=? AND progress<=?  AND state IN (0,1,2,6,8,9,10,11) limit 1000"
	_updateSubPoolSQL     = "UPDATE dm_subject_%02d SET childpool=? WHERE type=? AND oid=?"
	_incrSubMoveCntSQL    = "UPDATE dm_subject_%02d SET move_count=move_count+? WHERE type=? AND oid=?"
	_updateSubMCountSQL   = "UPDATE dm_subject_%02d SET mcount=? WHERE type=? AND oid=?"
	_incrSubCountSQL      = "UPDATE dm_subject_%02d SET count=count+? WHERE type=? AND oid=?"
	_getSpecialLocation   = "SELECT id,type,oid,locations FROM dm_special_content_location WHERE oid=? AND type=?"
	_getSpecialIdxSQL     = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE oid=? AND type=? AND state IN(0,6) AND pool=2"
	// upper config
	_addUpperCfgSQL = "REPLACE INTO dm_upper_config(mid,advance_permit) VALUES(?,?)"
	_selUpperCfgSQL = "SELECT advance_permit FROM dm_upper_config WHERE mid=?"
)

// Subject get subject info from db.
func (d *Dao) Subject(c context.Context, tp int32, oid int64) (s *model.Subject, err error) {
	s = &model.Subject{}
	row := d.dmReader.QueryRow(c, fmt.Sprintf(_subSQL, d.hitSubject(oid)), tp, oid)
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
			rows, err := d.dmReader.Query(errCtx, fmt.Sprintf(_getSubjectsSQL, k, xstr.JoinInts(v)), tp)
			if err != nil {
				log.Error("dmReader.Query() error(%v)", err)
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
			err = rows.Err()
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

// UptSubAttr update subject attr
func (d *Dao) UptSubAttr(c context.Context, tp int32, oid int64, attr int32) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateSubAttr, d.hitSubject(oid)), attr, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec(%s,%d,%d,%d) error(%v)", _updateSubAttr, attr, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// IncrSubMoveCount update move_count in dm_subject.
func (d *Dao) IncrSubMoveCount(c context.Context, tp int32, oid, count int64) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_incrSubMoveCntSQL, d.hitSubject(oid)), count, tp, oid)
	if err != nil {
		log.Error("d.dmWriter.Exec(%s,%d,%d,%d) error(%v)", _incrSubMoveCntSQL, tp, oid, count, err)
		return
	}
	return res.RowsAffected()
}

// UpSubjectMCount update monitor dm count.
func (d *Dao) UpSubjectMCount(c context.Context, tp int32, oid, cnt int64) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateSubMCountSQL, d.hitSubject(oid)), cnt, tp, oid)
	if err != nil {
		log.Error("d.dmWriter.Exect(%s,%d,%d,%d) error(%v)", _updateSubMCountSQL, cnt, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpSubjectPool update childpool in dm subject.
func (d *Dao) UpSubjectPool(c context.Context, tp int32, oid int64, childpool int32) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateSubPoolSQL, d.hitSubject(oid)), childpool, tp, oid)
	if err != nil {
		log.Error("d.dmWriter.Exect(%s,%d,%d,%d) error(%v)", _updateSubPoolSQL, childpool, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// IncrSubjectCount update count.
func (d *Dao) IncrSubjectCount(c context.Context, tp int32, oid, count int64) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_incrSubCountSQL, d.hitSubject(oid)), count, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec() error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DMIDs normal dm ids
func (d *Dao) DMIDs(c context.Context, tp int32, oid, ps, pe, limit int64, pool int32) (dmids []int64, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_idxSegIDSQL, d.hitIndex(oid)), tp, oid, ps, pe, pool, limit)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	var dmid int64
	for rows.Next() {
		if err = rows.Scan(&dmid); err != nil {
			log.Error("rows.Scan() error(%v)", err)
			return
		}
		dmids = append(dmids, dmid)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Indexs get dm index by oid.
func (d *Dao) Indexs(c context.Context, tp int32, oid int64) (idxMap map[int64]*model.DM, dmids, special []int64, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_getIndexSQL, d.hitIndex(oid)), tp, oid)
	if err != nil {
		log.Error("dmReader.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	idxMap = make(map[int64]*model.DM)
	for rows.Next() {
		idx := &model.DM{}
		if err = rows.Scan(&idx.ID, &idx.Type, &idx.Oid, &idx.Mid, &idx.Progress, &idx.State, &idx.Pool, &idx.Attr, &idx.Ctime, &idx.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		idxMap[idx.ID] = idx
		dmids = append(dmids, idx.ID)
		if idx.Pool == model.PoolSpecial {
			special = append(special, idx.ID)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// IndexByid get index by dmid.
func (d *Dao) IndexByid(c context.Context, tp int8, oid, dmid int64) (idx *model.DM, err error) {
	row := d.dmReader.QueryRow(c, fmt.Sprintf(_idxByidSQL, d.hitIndex(oid)), tp, dmid)
	idx = &model.DM{}
	err = row.Scan(&idx.ID, &idx.Type, &idx.Oid, &idx.Mid, &idx.Progress, &idx.State, &idx.Pool, &idx.Attr, &idx.Ctime, &idx.Mtime)
	if err != nil {
		if err == sql.ErrNoRows {
			idx = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// IndexsByid get dm index by dmids.
func (d *Dao) IndexsByid(c context.Context, tp int32, oid int64, dmids []int64) (idxMap map[int64]*model.DM, special []int64, err error) {
	query := fmt.Sprintf(_idxsByidSQL, d.hitIndex(oid), xstr.JoinInts(dmids))
	rows, err := d.dmReader.Query(c, query, tp)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	idxMap = make(map[int64]*model.DM, len(dmids))
	for rows.Next() {
		idx := &model.DM{}
		if err = rows.Scan(&idx.ID, &idx.Type, &idx.Oid, &idx.Mid, &idx.Progress, &idx.State, &idx.Pool, &idx.Attr, &idx.Ctime, &idx.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		idxMap[idx.ID] = idx
		if idx.Pool == model.PoolSpecial {
			special = append(special, idx.ID)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// JudgeIndex get judget index.
func (d *Dao) JudgeIndex(c context.Context, tp int8, oid int64, ctime1, ctime2 time.Time, prog1, prog2 int32) (idxs []*model.DM, special []int64, err error) {
	query := fmt.Sprintf(_judgeIdxPageSQL, d.hitIndex(oid))
	rows, err := d.dmReader.Query(c, query, tp, oid, ctime1, ctime2, prog1, prog2)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		idx := &model.DM{}
		if err = rows.Scan(&idx.ID, &idx.Type, &idx.Oid, &idx.Mid, &idx.Progress, &idx.State, &idx.Pool, &idx.Attr, &idx.Ctime, &idx.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		idxs = append(idxs, idx)
		if idx.Pool == model.PoolSpecial {
			special = append(special, idx.ID)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Content dm content by dmid
func (d *Dao) Content(c context.Context, oid, dmid int64) (ct *model.Content, err error) {
	ct = &model.Content{}
	row := d.dmReader.QueryRow(c, fmt.Sprintf(_getContentSQL, d.hitContent(oid)), dmid)
	if err = row.Scan(&ct.ID, &ct.FontSize, &ct.Color, &ct.Mode, &ct.IP, &ct.Plat, &ct.Msg, &ct.Ctime, &ct.Mtime); err != nil {
		ct = nil
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// Contents multi get dm content by dmids.
func (d *Dao) Contents(c context.Context, oid int64, dmids []int64) (ctsMap map[int64]*model.Content, err error) {
	var (
		wg   errgroup.Group
		lock sync.Mutex
	)
	ctsMap = make(map[int64]*model.Content)
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
			rows, err := d.dmReader.Query(c, fmt.Sprintf(_contentsSQL, d.hitContent(oid), xstr.JoinInts(dmids[start:end])))
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
				ctsMap[ct.ID] = ct
				lock.Unlock()
			}
			err = rows.Err()
			return
		})
	}
	if err = wg.Wait(); err != nil {
		log.Error("wg.Wait() error(%v)", err)
	}
	return
}

// ContentsSpecial multi get special dm content by dmids.
func (d *Dao) ContentsSpecial(c context.Context, dmids []int64) (res map[int64]*model.ContentSpecial, err error) {
	res = make(map[int64]*model.ContentSpecial, len(dmids))
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_contentsSpeSQL, xstr.JoinInts(dmids)))
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
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// UpdateDMStat edit dm state
func (d *Dao) UpdateDMStat(c context.Context, tp int32, oid int64, state int32, dmids []int64) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateIdxStatSQL, d.hitIndex(oid), xstr.JoinInts(dmids)), state, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec(%s %d) error(%v)", _updateIdxStatSQL, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpdateUserDMStat edit user dm state
func (d *Dao) UpdateUserDMStat(c context.Context, tp int32, oid, mid int64, state int32, dmids []int64) (affect int64, err error) {
	if mid <= 0 || tp <= 0 || oid <= 0 || len(dmids) <= 0 {
		affect = 0
		err = errors.New("d.UpdateUserDMStat: invalid arguments")
		return
	}
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateUserIdxStatSQL, d.hitIndex(oid), xstr.JoinInts(dmids)), state, tp, oid, mid)
	if err != nil {
		log.Error("dmWriter.Exec(%s %d) error(%v)", _updateIdxStatSQL, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpdateDMPool edit dm pool.
func (d *Dao) UpdateDMPool(c context.Context, tp int32, oid int64, pool int32, dmids []int64) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateIdxPoolSQL, d.hitIndex(oid), xstr.JoinInts(dmids)), pool, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec(%s %d) error(%v)", _updateIdxPoolSQL, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpdateDMAttr update dm attr
func (d *Dao) UpdateDMAttr(c context.Context, tp int32, oid, dmid int64, attr int32) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateIdxAttrSQL, d.hitIndex(oid)), attr, dmid)
	if err != nil {
		log.Error("dmWriter.Exec(%s oid:%d dmid:%d attr:%d) error(%v)", _updateIdxAttrSQL, oid, dmid, attr, err)
		return
	}
	return res.RowsAffected()
}

// DMCount statistics dm count by dm state.
func (d *Dao) DMCount(c context.Context, typ int32, oid int64, states []int64) (count int64, err error) {
	row := d.dmReader.QueryRow(c, fmt.Sprintf(_selectDMCount, d.hitIndex(oid), xstr.JoinInts(states)), typ, oid)
	if err = row.Scan(&count); err != nil {
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// SpecialDmLocation get special dm localtion url
func (d *Dao) SpecialDmLocation(c context.Context, tp int32, oid int64) (ds *model.DmSpecial, err error) {
	row := d.dmReader.QueryRow(c, _getSpecialLocation, oid, tp)
	ds = &model.DmSpecial{}
	if err = row.Scan(&ds.ID, &ds.Type, &ds.Oid, &ds.Locations); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			ds = nil
		} else {
			log.Error("SpecialDmLocation.Query(tp:%v,oid:%v) error(%v)", tp, oid, err)
		}
	}
	return
}

// SpecalDMs .
func (d *Dao) SpecalDMs(c context.Context, tp int32, oid int64) (dms map[int64]*model.DM, dmids []int64, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_getSpecialIdxSQL, d.hitIndex(oid)), oid, tp)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("SpecalDMs.Query(tp:%v,oid:%v) error(%v)", tp, oid, err)
		}
		return
	}
	defer rows.Close()
	dms = make(map[int64]*model.DM)
	for rows.Next() {
		dm := &model.DM{}
		if err = rows.Scan(&dm.ID, &dm.Type, &dm.Oid, &dm.Mid, &dm.Progress, &dm.State, &dm.Pool, &dm.Attr, &dm.Ctime, &dm.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		dms[dm.ID] = dm
		dmids = append(dmids, dm.ID)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// AddUpperConfig add upper config.
func (d *Dao) AddUpperConfig(c context.Context, mid int64, advPermit int8) (affect int64, err error) {
	res, err := d.dbDM.Exec(c, _addUpperCfgSQL, mid, advPermit)
	if err != nil {
		log.Error("dbDM.Exec(%s,%d,%d) error(%v)", _addUpperCfgSQL, mid, advPermit, err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// UpperConfig get upper config.
func (d *Dao) UpperConfig(c context.Context, mid int64) (advPermit int8, err error) {
	row := d.dbDM.QueryRow(c, _selUpperCfgSQL, mid)
	if err = row.Scan(&advPermit); err != nil {
		if err == sql.ErrNoRows {
			advPermit = model.AdvPermitAll
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}
