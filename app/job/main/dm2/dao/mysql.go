package dao

import (
	"context"
	"fmt"
	"sync"

	"go-common/app/job/main/dm2/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_subtitleSharding = 100
	_subjectSharding  = 100
	_indexSharding    = 1000
	_contentSharding  = 1000

	_addSubjectSQL      = "INSERT INTO dm_subject_%02d(type,oid,pid,mid,maxlimit,attr) VALUES(?,?,?,?,?,?)"
	_updateChildpoolSQL = "UPDATE dm_subject_%02d SET childpool=? WHERE type=? AND oid=?"
	_updateSubMidSQL    = "UPDATE dm_subject_%02d SET mid=? WHERE type=? AND oid=?"
	_updateSubAttrSQL   = "UPDATE dm_subject_%02d SET attr=? WHERE type=? AND oid=?"
	_incrSubMCountSQL   = "UPDATE dm_subject_%02d SET mcount=mcount+1 WHERE type=? AND oid=?"
	_incrSubCountSQL    = "UPDATE dm_subject_%02d SET acount=acount+?,count=count+?,childpool=? WHERE type=? AND oid=?"
	_getSubjectSQL      = "SELECT id,type,oid,pid,mid,state,attr,acount,count,mcount,move_count,maxlimit,childpool,ctime,mtime FROM dm_subject_%02d WHERE type=? AND oid=?"
	_addIndexSQL        = "INSERT INTO dm_index_%03d(id,type,oid,mid,progress,state,pool,attr,ctime) VALUES(?,?,?,?,?,?,?,?,?)"
	_getIndexSQL        = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND oid=? AND state IN(0,6)"
	_idxSegIDSQL        = "SELECT id FROM dm_index_%03d WHERE type=? AND oid=? AND progress>=? AND progress<? AND state IN(0,6) AND pool = ? limit ?"
	_idxSegSQL          = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND oid=? AND state IN(0,6) AND progress>=? AND progress<? AND pool=? limit ?"
	_idxIDSQL           = "SELECT id FROM dm_index_%03d WHERE type=? AND oid=? AND state IN(0,6) AND pool=?"
	_idxsByidSQL        = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE id IN(%s)"
	_idxsByPoolSQL      = "SELECT id,type,oid,mid,progress,state,pool,attr,ctime,mtime FROM dm_index_%03d WHERE type=? AND oid=? AND state IN(0,6) AND pool=?"
	_addContentSQL      = "REPLACE INTO dm_content_%03d(dmid,fontsize,color,mode,ip,plat,msg,ctime) VALUES(?,?,?,?,?,?,?,?)"
	_getContentsSQL     = "SELECT dmid,fontsize,color,mode,ip,plat,msg,ctime,mtime FROM dm_content_%03d WHERE dmid IN(%s)"
	_getContentSQL      = "SELECT dmid,fontsize,color,mode,ip,plat,msg,ctime,mtime FROM dm_content_%03d WHERE dmid=?"
	_addContentSpeSQL   = "REPLACE INTO dm_special_content(dmid,msg,ctime) VALUES(?,?,?)"
	_getContentSpeSQL   = "SELECT dmid,msg,ctime,mtime FROM dm_special_content WHERE dmid=?"
	_getContentsSpeSQL  = "SELECT dmid,msg,ctime,mtime FROM dm_special_content WHERE dmid IN(%s)"
	//delete dm hide state
	_delDMHideState = "UPDATE dm_index_%03d SET state=? WHERE oid=? AND id=? AND state=?"
	// update subtitle upmid
	_getSubtitle    = "SELECT id,oid,type,lan,status,mid,up_mid,subtitle_url,pub_time,reject_comment from subtitle_%02d WHERE id=?"
	_getSubtitles   = "SELECT id,oid,type,lan,status,mid,up_mid,subtitle_url,pub_time,reject_comment from subtitle_%02d WHERE oid=? AND type=?"
	_updateSubtitle = "UPDATE subtitle_%02d SET up_mid=?,status=?,pub_time=?,reject_comment=? WHERE id=?"
	_addSubtitlePub = "INSERT INTO subtitle_pub(oid,type,lan,subtitle_id,is_delete) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE subtitle_id=?,is_delete=?"
	// get mask mid
	_getMaskMids = "SELECT mid from dm_mask_up where state=1"
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

func (d *Dao) hitSubtile(oid int64) int64 {
	return oid % _subtitleSharding
}

// AddSubject insert subject.
func (d *Dao) AddSubject(c context.Context, tp int32, oid, pid, mid, maxlimit int64, attr int32) (lastID int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_addSubjectSQL, d.hitSubject(oid)), tp, oid, pid, mid, maxlimit, attr)
	if err != nil {
		log.Error("dmWriter.Exec(%d,%d,%d,%d,%d,%d) error(%v)", tp, oid, pid, mid, maxlimit, attr, err)
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// UpdateSubAttr .
func (d *Dao) UpdateSubAttr(c context.Context, tp int32, oid int64, attr int32) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateSubAttrSQL, d.hitSubject(oid)), attr, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec(%s,%d,%d) error(%v)", _updateSubMidSQL, oid, attr, err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// UpdateSubMid update mid in dm_subject.
func (d *Dao) UpdateSubMid(c context.Context, tp int32, oid, mid int64) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateSubMidSQL, d.hitSubject(oid)), mid, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec(%s,%d,%d) error(%v)", _updateSubMidSQL, oid, mid, err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// Subject get subject info from db.
func (d *Dao) Subject(c context.Context, tp int32, oid int64) (s *model.Subject, err error) {
	s = &model.Subject{}
	row := d.dmReader.QueryRow(c, fmt.Sprintf(_getSubjectSQL, d.hitSubject(oid)), tp, oid)
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

// UpdateChildpool update childpool.
func (d *Dao) UpdateChildpool(c context.Context, tp int32, oid int64, childpool int32) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_updateChildpoolSQL, d.hitSubject(oid)), childpool, tp, oid)
	if err != nil {
		log.Error("dmWriter.Exec(%s %d) error(%v)", _updateChildpoolSQL, oid, err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// TxIncrSubjectCount update acount,count,childpool of dm by transcation.
func (d *Dao) TxIncrSubjectCount(tx *sql.Tx, tp int32, oid, acount, count int64, childpool int32) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubCountSQL, d.hitSubject(oid)), acount, count, childpool, tp, oid)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxAddIndex add index of dm by transcation.
func (d *Dao) TxAddIndex(tx *sql.Tx, m *model.DM) (id int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_addIndexSQL, d.hitIndex(m.Oid)), m.ID, m.Type, m.Oid, m.Mid, m.Progress, m.State, m.Pool, m.Attr, m.Ctime)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// Indexs get dm index by type and oid.
func (d *Dao) Indexs(c context.Context, tp int32, oid int64) (idxMap map[int64]*model.DM, dmids, special []int64, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_getIndexSQL, d.hitIndex(oid)), tp, oid)
	if err != nil {
		log.Error("dmReader.Query(%d,%d) error(%v)", tp, oid, err)
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
	err = rows.Err()
	return
}

// IndexsSeg get segment index info from db by ps and pe.
func (d *Dao) IndexsSeg(c context.Context, tp int32, oid, ps, pe, limit int64, pool int32) (res []*model.DM, dmids []int64, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_idxSegSQL, d.hitIndex(oid)), tp, oid, ps, pe, pool, limit)
	if err != nil {
		log.Error("db.Query(%d %d %d %d) error(%v)", tp, oid, ps, pe, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		dm := &model.DM{}
		if err = rows.Scan(&dm.ID, &dm.Type, &dm.Oid, &dm.Mid, &dm.Progress, &dm.State, &dm.Pool, &dm.Attr, &dm.Ctime, &dm.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		res = append(res, dm)
		dmids = append(dmids, dm.ID)
	}
	return
}

// IndexsSegID get segment dmids.
func (d *Dao) IndexsSegID(c context.Context, tp int32, oid, ps, pe, limit int64, pool int32) (dmids []int64, err error) {
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

// IndexsID get dmids.
func (d *Dao) IndexsID(c context.Context, tp int32, oid int64, pool int32) (dmids []int64, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_idxIDSQL, d.hitIndex(oid)), tp, oid, pool)
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

// IndexsByid get dm index by dmids.
func (d *Dao) IndexsByid(c context.Context, tp int32, oid int64, dmids []int64) (idxMap map[int64]*model.DM, special []int64, err error) {
	query := fmt.Sprintf(_idxsByidSQL, d.hitIndex(oid), xstr.JoinInts(dmids))
	rows, err := d.dmReader.Query(c, query)
	if err != nil {
		log.Error("db.Query(%s) error(%v)", query, err)
		return
	}
	defer rows.Close()
	idxMap = make(map[int64]*model.DM)
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
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// IndexsByPool get dm index by type,oid and pool.
func (d *Dao) IndexsByPool(c context.Context, tp int32, oid int64, pool int32) (dms []*model.DM, dmids []int64, err error) {
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_idxsByPoolSQL, d.hitIndex(oid)), tp, oid, pool)
	if err != nil {
		log.Error("dmReader.Query(tp:%v,oid:%v) error(%v)", tp, oid, err)
		return
	}
	defer rows.Close()
	dms = make([]*model.DM, 0, 100)
	for rows.Next() {
		dm := &model.DM{}
		if err = rows.Scan(&dm.ID, &dm.Type, &dm.Oid, &dm.Mid, &dm.Progress, &dm.State, &dm.Pool, &dm.Attr, &dm.Ctime, &dm.Mtime); err != nil {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		dms = append(dms, dm)
		dmids = append(dmids, dm.ID)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// TxAddContent add content of dm by transcation.
func (d *Dao) TxAddContent(tx *sql.Tx, oid int64, m *model.Content) (id int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_addContentSQL, d.hitContent(oid)), m.ID, m.FontSize, m.Color, m.Mode, m.IP, m.Plat, m.Msg, m.Ctime)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
}

// TxAddContentSpecial add special dm by transcation.
func (d *Dao) TxAddContentSpecial(tx *sql.Tx, m *model.ContentSpecial) (id int64, err error) {
	res, err := tx.Exec(_addContentSpeSQL, m.ID, m.Msg, m.Ctime)
	if err != nil {
		log.Error("tx.Exec error(%v)", err)
		return
	}
	return res.LastInsertId()
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
	pageNum := len(dmids) / d.pageSize
	if len(dmids)%d.pageSize > 0 {
		pageNum = pageNum + 1
	}
	for i := 0; i < pageNum; i++ {
		start := i * d.pageSize
		end := (i + 1) * d.pageSize
		if end > len(dmids) {
			end = len(dmids)
		}
		wg.Go(func() (err error) {
			rows, err := d.dmReader.Query(c, fmt.Sprintf(_getContentsSQL, d.hitContent(oid), xstr.JoinInts(dmids[start:end])))
			if err != nil {
				log.Error("db.Query(%s) error(%v)", fmt.Sprintf(_getContentsSQL, d.hitContent(oid), xstr.JoinInts(dmids)), err)
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
	rows, err := d.dmReader.Query(c, fmt.Sprintf(_getContentsSpeSQL, xstr.JoinInts(dmids)))
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

// ContentSpecial get special dm content by dmids.
func (d *Dao) ContentSpecial(c context.Context, dmid int64) (contentSpe *model.ContentSpecial, err error) {
	contentSpe = &model.ContentSpecial{}
	row := d.dmReader.QueryRow(c, _getContentSpeSQL, dmid)
	if err = row.Scan(&contentSpe.ID, &contentSpe.Msg, &contentSpe.Ctime, &contentSpe.Mtime); err != nil {
		log.Error("rows.Scan() error(%v)", err)
	}
	return
}

// DelDMHideState del dm hide state
func (d *Dao) DelDMHideState(c context.Context, tp int32, oid int64, dmid int64) (affect int64, err error) {
	res, err := d.dmWriter.Exec(c, fmt.Sprintf(_delDMHideState, d.hitIndex(oid)), model.StateNormal, oid, dmid, model.StateHide)
	if err != nil {
		log.Error("dmWriter.Exec(%s %d dmid=%d) error(%v)", _delDMHideState, oid, dmid, err)
		return
	}
	return res.RowsAffected()
}

// TxIncrSubMCount update monitor dm count.
func (d *Dao) TxIncrSubMCount(tx *sql.Tx, tp int32, oid int64) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_incrSubMCountSQL, d.hitSubject(oid)), tp, oid)
	if err != nil {
		log.Error("tx.Exec(%s,%d,%d) error(%v)", _incrSubMCountSQL, tp, oid, err)
		return
	}
	return res.RowsAffected()
}

// UpdateSubtitle update subtitle mid
func (d *Dao) UpdateSubtitle(c context.Context, subtitle *model.Subtitle) (err error) {
	if _, err = d.biliDMWriter.Exec(c, fmt.Sprintf(_updateSubtitle, d.hitSubtile(subtitle.Oid)), subtitle.UpMid, subtitle.Status,
		subtitle.PubTime, subtitle.RejectComment, subtitle.ID); err != nil {
		log.Error("biliDMWriter.Exec(query:%v,subtitle:%+v) error(%v)", _updateSubtitle, subtitle, err)
		return
	}
	return
}

// GetSubtitles .
func (d *Dao) GetSubtitles(c context.Context, tp int32, oid int64) (subtitles []*model.Subtitle, err error) {
	rows, err := d.biliDMWriter.Query(c, fmt.Sprintf(_getSubtitles, d.hitSubtile(oid)), oid, tp)
	if err != nil {
		log.Error("biliDMWriter.Query(%s,%d,%d) error(%v)", _getSubtitles, oid, tp, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var subtitle = &model.Subtitle{}
		if err = rows.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Status, &subtitle.Mid, &subtitle.UpMid,
			&subtitle.SubtitleURL, &subtitle.PubTime, &subtitle.RejectComment); err != nil {
			log.Error("biliDMWriter.Scan(%s,%d,%d) error(%v)", _getSubtitles, oid, tp, err)
			return
		}
		subtitles = append(subtitles, subtitle)
	}
	if err = rows.Err(); err != nil {
		log.Error("biliDMWriter.rows.Err()(%s,%d,%d) error(%v)", _getSubtitles, oid, tp, err)
		return
	}
	return
}

// GetSubtitle .
func (d *Dao) GetSubtitle(c context.Context, oid int64, subtitleID int64) (subtitle *model.Subtitle, err error) {
	subtitle = &model.Subtitle{}
	row := d.biliDMWriter.QueryRow(c, fmt.Sprintf(_getSubtitle, d.hitSubtile(oid)), subtitleID)
	if err = row.Scan(&subtitle.ID, &subtitle.Oid, &subtitle.Type, &subtitle.Lan, &subtitle.Status, &subtitle.Mid, &subtitle.UpMid,
		&subtitle.SubtitleURL, &subtitle.PubTime, &subtitle.RejectComment); err != nil {
		if err == sql.ErrNoRows {
			subtitle = nil
			err = nil
		} else {
			log.Error("row.Scan() error(%v)", err)
		}
	}
	return
}

// TxUpdateSubtitle .
func (d *Dao) TxUpdateSubtitle(tx *sql.Tx, subtitle *model.Subtitle) (err error) {
	if _, err = tx.Exec(fmt.Sprintf(_updateSubtitle, d.hitSubtile(subtitle.Oid)), subtitle.UpMid, subtitle.Status,
		subtitle.PubTime, subtitle.RejectComment, subtitle.ID); err != nil {
		log.Error("params(%+v),error(%v)", subtitle, err)
		return
	}
	return
}

// TxAddSubtitlePub .
func (d *Dao) TxAddSubtitlePub(tx *sql.Tx, subtitlePub *model.SubtitlePub) (err error) {
	if _, err = tx.Exec(_addSubtitlePub, subtitlePub.Oid, subtitlePub.Type, subtitlePub.Lan, subtitlePub.SubtitleID, subtitlePub.IsDelete, subtitlePub.SubtitleID, subtitlePub.IsDelete); err != nil {
		log.Error("params(%+v),error(%v)", subtitlePub, err)
		return
	}
	return
}

// MaskMids get mask mids from db.
func (d *Dao) MaskMids(c context.Context) (mids []int64, err error) {
	mids = make([]int64, 0, 100)
	rows, err := d.biliDMWriter.Query(c, _getMaskMids)
	if err != nil {
		log.Error("biliDMWriter.Query(%s) error(%v)", _getMaskMids, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var mid int64
		if err = rows.Scan(&mid); err != nil {
			log.Error("biliDMWriter.Scan(%s) error(%v)", _getMaskMids, err)
			return
		}
		mids = append(mids, mid)
	}
	if err = rows.Err(); err != nil {
		log.Error("biliDMWriter.rows.Err() error(%v)", err)
	}
	return
}
