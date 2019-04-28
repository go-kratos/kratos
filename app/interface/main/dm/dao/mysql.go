package dao

import (
	"context"
	"fmt"
	"sync"

	"go-common/app/interface/main/dm/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_pagesize        = 1000
	_subjectSharding = 100
	_indexSharding   = 1000
	_contentSharding = 1000
	_subSQL          = "SELECT id,type,oid,pid,mid,state,attr,acount,count,mcount,move_count,maxlimit,childpool,ctime,mtime FROM dm_subject_%02d WHERE type=? AND oid=?"
	_idxSQL          = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE id=? AND oid=? AND type=?"
	_idxsByIDSQL     = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE id IN(%s) AND oid=? AND type=?"
	_contentSQL      = "SELECT dmid,fontsize,color,mode,ip,plat,msg,ctime,mtime FROM dm_content_%03d WHERE dmid=?"
	_contentsSQL     = "SELECT dmid,fontsize,color,mode,ip,plat,msg,ctime,mtime FROM dm_content_%03d WHERE dmid IN(%s)"
	_contentSpeSQL   = "SELECT dmid,msg,ctime,mtime FROM dm_special_content WHERE dmid=?"
	_contentsSpeSQL  = "SELECT dmid,msg,ctime,mtime FROM dm_special_content WHERE dmid IN(%s)"
	// dm transfer
	_addTransferJob  = "INSERT INTO dm_transfer_job SET from_cid=?,to_cid=?,mid=?,offset=?,state=?"
	_selTransferJob  = "SELECT id,from_cid,to_cid,mid,offset,state,ctime,mtime from dm_transfer_job WHERE from_cid=? AND to_cid=?"
	_transferListSQL = "SELECT id,from_cid,state,ctime FROM dm_transfer_job where to_cid=?"
	_uptTransferSQL  = "UPDATE dm_transfer_job SET state=? WHERE id=?"
	_selTransferID   = "SELECT id,from_cid,to_cid,mid,offset,state,ctime,mtime from dm_transfer_job WHERE id=?"
	//dm state update
	_updateIdxStatSQL = "UPDATE dm_index_%03d SET state=? WHERE type=? AND oid=? AND id IN(%s)"
)

func (d *Dao) hitSubject(oid int64) int64 {
	return oid % _subjectSharding
}

func (d *Dao) hitIndex(oid int64) int64 {
	return oid % _indexSharding
}

func (d *Dao) hitContent(dmid int64) int64 {
	return dmid % _contentSharding
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

// Index get index by dmid
func (d *Dao) Index(c context.Context, tp int32, oid, dmid int64) (dm *model.DM, err error) {
	dm = &model.DM{}
	row := d.dmMetaReader.QueryRow(c, fmt.Sprintf(_idxSQL, d.hitIndex(oid)), dmid, oid, tp)
	if err = row.Scan(&dm.ID, &dm.Type, &dm.Oid, &dm.Mid, &dm.Progress, &dm.State, &dm.Pool, &dm.Attr, &dm.Ctime, &dm.Mtime); err != nil {
		if err == sql.ErrNoRows {
			dm = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// IndexsByID get dm index by dmids.
func (d *Dao) IndexsByID(c context.Context, tp int32, oid int64, dmids []int64) (res map[int64]*model.DM, special []int64, err error) {
	rows, err := d.dmMetaReader.Query(c, fmt.Sprintf(_idxsByIDSQL, d.hitIndex(oid), xstr.JoinInts(dmids)), oid, tp)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.DM, len(dmids))
	for rows.Next() {
		dm := &model.DM{}
		if err = rows.Scan(&dm.ID, &dm.Type, &dm.Oid, &dm.Mid, &dm.Progress, &dm.State, &dm.Pool, &dm.Attr, &dm.Ctime, &dm.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res[dm.ID] = dm
		if dm.Pool == model.PoolSpecial {
			special = append(special, dm.ID)
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
	row := d.dmMetaReader.QueryRow(c, fmt.Sprintf(_contentSQL, d.hitContent(oid)), dmid)
	if err = row.Scan(&ct.ID, &ct.FontSize, &ct.Color, &ct.Mode, &ct.IP, &ct.Plat, &ct.Msg, &ct.Ctime, &ct.Mtime); err != nil {
		if err == sql.ErrNoRows {
			ct = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
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
	pageNum := len(dmids) / _pagesize
	if len(dmids)%_pagesize > 0 {
		pageNum = pageNum + 1
	}
	for i := 0; i < pageNum; i++ {
		start := i * _pagesize
		end := (i + 1) * _pagesize
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
			err = rows.Err()
			return
		})
	}
	if err = wg.Wait(); err != nil {
		log.Error("wg.Wait() error(%v)", err)
	}
	return
}

// ContentSpecial get special dm content by dmids.
func (d *Dao) ContentSpecial(c context.Context, dmid int64) (contentSpe *model.ContentSpecial, err error) {
	contentSpe = &model.ContentSpecial{}
	row := d.dmMetaReader.QueryRow(c, _contentSpeSQL, dmid)
	if err = row.Scan(&contentSpe.ID, &contentSpe.Msg, &contentSpe.Ctime, &contentSpe.Mtime); err != nil {
		if err == sql.ErrNoRows {
			contentSpe = nil
			err = nil
		} else {
			log.Error("rows.Scan() error(%v)", err)
		}
	}
	return
}

// ContentsSpecial multi get special dm content by dmids.
func (d *Dao) ContentsSpecial(c context.Context, dmids []int64) (res map[int64]*model.ContentSpecial, err error) {
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
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// AddTransferJob add transfer job.
func (d *Dao) AddTransferJob(c context.Context, fromCid, toCid, mid int64, offset float64, state int8) (affect int64, err error) {
	row, err := d.biliDM.Exec(c, _addTransferJob, fromCid, toCid, mid, offset, model.TransferJobStatInit)
	if err != nil {
		log.Error("d.biliDM.Exec(fromCid:%d,toCid:%d,mid:%d,offset:%v) error(%v)", fromCid, toCid, mid, offset, err)
	}
	return row.LastInsertId()
}

// CheckTransferJob check transfer job state.
func (d *Dao) CheckTransferJob(c context.Context, fromCid, toCid int64) (job *model.TransferJob, err error) {
	job = new(model.TransferJob)
	row := d.biliDM.QueryRow(c, _selTransferJob, fromCid, toCid)
	if err = row.Scan(&job.ID, &job.FromCID, &job.ToCID, &job.MID, &job.Offset, &job.State, &job.Ctime, &job.Mtime); err != nil {
		if err == sql.ErrNoRows {
			job = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// CheckTransferID check transfer job state by id
func (d *Dao) CheckTransferID(c context.Context, id int64) (job *model.TransferJob, err error) {
	job = new(model.TransferJob)
	row := d.biliDM.QueryRow(c, _selTransferID, id)
	if err = row.Scan(&job.ID, &job.FromCID, &job.ToCID, &job.MID, &job.Offset, &job.State, &job.Ctime, &job.Mtime); err != nil {
		if err == sql.ErrNoRows {
			job = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// TransferList cid 的转移历史
func (d *Dao) TransferList(c context.Context, cid int64) (l []*model.TransferHistory, err error) {
	l = make([]*model.TransferHistory, 0)
	rows, err := d.biliDM.Query(c, _transferListSQL, cid)
	if err != nil {
		log.Error("d.biliDM.Query(%s,%d) error(%v)", _transferListSQL, cid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		hs := new(model.TransferHistory)
		if err = rows.Scan(&hs.ID, &hs.CID, &hs.State, &hs.CTime); err != nil {
			log.Error("TransfrerList: rows.Scan() error(%v)", err)
			return
		}
		l = append(l, hs)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// SetTransferState change transfer state
func (d *Dao) SetTransferState(c context.Context, id int64, state int8) (affect int64, err error) {
	row, err := d.biliDM.Exec(c, _uptTransferSQL, state, id)
	if err != nil {
		log.Error("d.biliDM.Exec(%s,%d) error(%v)", _uptTransferSQL, id, err)
		return
	}
	return row.RowsAffected()
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
