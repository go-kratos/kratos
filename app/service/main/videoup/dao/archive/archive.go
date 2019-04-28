package archive

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// insert
	_inArcSQL        = "INSERT INTO archive (mid,typeid,title,author,cover,content,tag,attribute,copyright,state,round,pubtime,ctime,mtime,reject_reason) VALUES (?,?,?,?,?,?,?,?,?,-30,0,?,?,?,'')"
	_inAddSQL        = "INSERT INTO archive_addit (aid,mission_id,up_from,ipv6,source,order_id,flow_id,advertiser,flow_remark,description,desc_format_id,dynamic) VALUES (?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE mission_id=?,source=?,order_id=?,flow_id=?,advertiser=?,flow_remark=?,description=?,desc_format_id=?,dynamic=?"
	_inAddRdrSQL     = "INSERT INTO archive_addit (aid,redirect_url) VALUES (?,?) ON DUPLICATE KEY UPDATE redirect_url=?"
	_inAddReaSQL     = "INSERT INTO archive_addit (aid,recheck_reason) VALUES (?,?) ON DUPLICATE KEY UPDATE recheck_reason=?"
	_inAddMissionSQL = "INSERT INTO archive_addit (aid,mission_id) VALUES (?,?) ON DUPLICATE KEY UPDATE mission_id=?"
	// update
	_upArcSQL      = "UPDATE archive SET typeid=?,title=?,cover=?,content=?,tag=?,copyright=? WHERE id=?"
	_upArcMidSQL   = "UPDATE archive SET mid=?,state=? WHERE id=?"
	_upArcStateSQL = "UPDATE archive SET state=? WHERE id=?"
	_upArcAttrSQL  = "UPDATE archive SET attribute=attribute&(~(1<<?))|(?<<?) WHERE id=?"
	_upTagSQL      = "UPDATE archive SET tag=? WHERE id=?"
	// select
	_arcSQL              = "SELECT id,mid,typeid,copyright,author,title,cover,reject_reason,content,tag,duration,round,attribute,access,state,pubtime,ctime,mtime FROM archive WHERE id=?"
	_arcAddSQL           = "SELECT aid,mission_id,from_ip,ipv6,up_from,recheck_reason,redirect_url,source,order_id,flow_id,advertiser,flow_remark,description,desc_format_id,dynamic FROM archive_addit WHERE aid=?"
	_arcMidsSQL          = "SELECT id,mid FROM archive WHERE id IN (%s)"
	_arcUpAllSQL         = `SELECT id FROM archive WHERE mid = ? AND state != -100 ORDER BY id DESC LIMIT ?,?`
	_arcUpOpenSQL        = `SELECT id FROM archive WHERE mid = ? AND (state >= 0 OR state = -6) ORDER BY id DESC LIMIT ?,?`
	_arcUpUnOpenSQL      = `SELECT id FROM archive WHERE mid = ? AND state < 0 AND state != -100 AND state != -6 ORDER BY id DESC LIMIT ?,?`
	_arcUpAllCountSQL    = "SELECT count(*) FROM archive WHERE mid = ? AND state != -100"
	_arcUpOpenCountSQL   = "SELECT count(*) FROM archive WHERE mid = ? AND (state >= 0 OR state = -6)"
	_arcUpUnOpenCountSQL = "SELECT count(*) FROM archive WHERE mid = ? AND state < 0 AND state != -100 AND state != -6"
	_simpleArcSQL        = "SELECT id,title,mid FROM archive WHERE id = ?"
	_getRecoSQL          = "SELECT reco_aid FROM archive_recommend WHERE state= 0 and aid=? ORDER BY ctime asc"
	_rejectArcsSQL       = "SELECT id,mid,title,reject_reason,mtime FROM archive WHERE mid = ? AND state = ? AND mtime > ? ORDER BY mtime DESC LIMIT ?,?"
	_rejectArcsCountSQL  = "SELECT count(*) FROM archive WHERE mid = ? AND state = ? AND mtime > ?"
	_delRecoSQL          = "UPDATE archive_recommend SET state=1 WHERE aid=?"
	_batchAddRecoSQL     = "INSERT IGNORE INTO archive_recommend (aid,reco_aid) VALUES %s on duplicate key update state=0"
	//POI 元数据
	_arcPOISQL   = "SELECT data from archive_biz   WHERE aid=? AND type= ?"
	_arcVoteSQL  = "SELECT data from archive_biz   WHERE aid=? AND type= 2"
	_inADDBizSQL = "INSERT INTO archive_biz (aid,type,data) VALUES (?,?,?) ON DUPLICATE KEY UPDATE data=?"
)

// TxAddArchive insert archive.
func (d *Dao) TxAddArchive(tx *xsql.Tx, a *archive.Archive) (aid int64, err error) {
	var now = time.Now()
	res, err := tx.Exec(_inArcSQL, a.Mid, a.TypeID, a.Title, a.Author, a.Cover, a.Desc, a.Tag, a.Attribute, a.Copyright, now, now, now)
	if err != nil {
		log.Error("d.inArc.Exec() error(%v)", err)
		return
	}
	if aid, err = res.LastInsertId(); err != nil {
		log.Error("res.LastInsertId() error(%v)", err)
		return
	}
	if strings.Contains(xstr.JoinInts(d.c.KeepArc.Aids)+",", strconv.FormatInt(aid, 10)+",") {
		keepAid := aid
		aid, err = d.TxAddArchive(tx, a)
		a.State = archive.StateForbidUpDelete
		a.Mid = d.c.KeepArc.Mid // 内部归属mid
		a.Aid = keepAid         // 内部保留aid
		d.TxUpArchiveMid(tx, a)
		return
	}
	return
}

// TxUpArchive update archive.
func (d *Dao) TxUpArchive(tx *xsql.Tx, a *archive.Archive) (rows int64, err error) {
	res, err := tx.Exec(_upArcSQL, a.TypeID, a.Title, a.Cover, a.Desc, a.Tag, a.Copyright, a.Aid)
	if err != nil {
		log.Error("d.upArc.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArchiveMid update archive mid.
func (d *Dao) TxUpArchiveMid(tx *xsql.Tx, a *archive.Archive) (rows int64, err error) {
	res, err := tx.Exec(_upArcMidSQL, a.Mid, a.State, a.Aid)
	if err != nil {
		log.Error("d.upArcMid.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArchiveState update Archive state.
func (d *Dao) TxUpArchiveState(tx *xsql.Tx, aid int64, state int8) (rows int64, err error) {
	res, err := tx.Exec(_upArcStateSQL, state, aid)
	if err != nil {
		log.Error("d.upVideoState.Exec error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAddit update archive addit.
func (d *Dao) TxUpAddit(tx *xsql.Tx, aid, missionID, orderID, flowID, descFormatID int64, ipv6 []byte, source, advertiser, flowRemark, desc, dynamic string, upFrom int8) (rows int64, err error) {
	if ipv6 == nil {
		ipv6 = []byte{}
	}
	res, err := tx.Exec(_inAddSQL, aid, missionID, upFrom, ipv6, source, orderID, flowID, advertiser, flowRemark, desc, descFormatID, dynamic, missionID, source, orderID, flowID, advertiser, flowRemark, desc, descFormatID, dynamic)
	if err != nil {
		log.Error("d.inArcAddit.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArchiveBiz update archive biz.
func (d *Dao) TxUpArchiveBiz(tx *xsql.Tx, aid, bizType int64, data string) (rows int64, err error) {
	res, err := tx.Exec(_inADDBizSQL, aid, bizType, data, data)
	if err != nil {
		log.Error("d.TxUpArchiveBiz.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAdditReason update archive recheck_reason
func (d *Dao) TxUpAdditReason(tx *xsql.Tx, aid int64, reason string) (rows int64, err error) {
	res, err := tx.Exec(_inAddReaSQL, aid, reason, reason)
	if err != nil {
		log.Error("d.inAdditReason.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpAdditRedirect update archive redirect url.
func (d *Dao) TxUpAdditRedirect(tx *xsql.Tx, aid int64, redirectURL string) (rows int64, err error) {
	res, err := tx.Exec(_inAddRdrSQL, aid, redirectURL, redirectURL)
	if err != nil {
		log.Error("d._inAdditRedirect.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpArcAttr update attribute by aid.
func (d *Dao) TxUpArcAttr(tx *xsql.Tx, aid int64, bit uint, val int32) (rows int64, err error) {
	res, err := tx.Exec(_upArcAttrSQL, bit, val, bit, aid)
	attSql := fmt.Sprintf("UPDATE archive SET attribute=attribute&(~(1<<%d))|(%d<<%d) WHERE id=%d", bit, val, bit, aid)
	log.Info("aid(%d) attribute update log sql (%s)", aid, attSql)
	if err != nil {
		log.Error("d.upArcAttr.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// TxUpTag update tag by aid.
func (d *Dao) TxUpTag(tx *xsql.Tx, aid int64, tag string) (rows int64, err error) {
	res, err := tx.Exec(_upTagSQL, tag, aid)
	if err != nil {
		log.Error("d.upTag.Exec() error(%v)", err)
		return
	}
	rows, err = res.RowsAffected()
	return
}

// Archive get a archive by avid.
func (d *Dao) Archive(c context.Context, aid int64) (a *archive.Archive, err error) {
	var (
		row         = d.rddb.QueryRow(c, _arcSQL, aid)
		reason, tag sql.NullString
	)
	a = &archive.Archive{}
	if err = row.Scan(&a.Aid, &a.Mid, &a.TypeID, &a.Copyright, &a.Author, &a.Title, &a.Cover, &reason, &a.Desc, &tag, &a.Duration,
		&a.Round, &a.Attribute, &a.Access, &a.State, &a.PTime, &a.CTime, &a.MTime); err != nil {
		if err == xsql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	a.RejectReason = reason.String
	a.Tag = tag.String
	return
}

// POI get a archive POI by avid.
func (d *Dao) POI(c context.Context, aid int64) (data []byte, err error) {
	var (
		row = d.rddb.QueryRow(c, _arcPOISQL, aid, archive.BIZPOI)
	)
	if err = row.Scan(&data); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// Vote get a archive Vote by avid.
func (d *Dao) Vote(c context.Context, aid int64) (data []byte, err error) {
	var (
		row = d.rddb.QueryRow(c, _arcVoteSQL, aid)
	)
	if err = row.Scan(&data); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
		return
	}
	return
}

// Addit get a archive addit by avid.
func (d *Dao) Addit(c context.Context, aid int64) (ad *archive.Addit, err error) {
	row := d.rddb.QueryRow(c, _arcAddSQL, aid)
	ad = &archive.Addit{}
	if err = row.Scan(&ad.Aid, &ad.MissionID, &ad.FromIP, &ad.IPv6, &ad.UpFrom, &ad.RecheckReason, &ad.RedirectURL, &ad.Source, &ad.OrderID, &ad.FlowID, &ad.Advertiser, &ad.FlowRemark, &ad.Desc, &ad.DescFormatID, &ad.Dynamic); err != nil {
		if err == xsql.ErrNoRows {
			err = nil
			ad = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Mids multi get archive mid by aids.
func (d *Dao) Mids(c context.Context, aids []int64) (mm map[int64]int64, err error) {
	rows, err := d.rddb.Query(c, fmt.Sprintf(_arcMidsSQL, xstr.JoinInts(aids)))
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	mm = make(map[int64]int64, len(aids))
	for rows.Next() {
		var aid, mid int64
		if err = rows.Scan(&aid, &mid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		mm[aid] = mid
	}
	return
}

// ArchivesUpAll  get archive all aids by mid.
func (d *Dao) ArchivesUpAll(c context.Context, mid int64, offset int, ps int) (aids []int64, err error) {
	rows, err := d.rddb.Query(c, _arcUpAllSQL, mid, offset, ps)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// ArchivesUpOpen  get archive open aids by mid.
func (d *Dao) ArchivesUpOpen(c context.Context, mid int64, offset int, ps int) (aids []int64, err error) {
	rows, err := d.rddb.Query(c, _arcUpOpenSQL, mid, offset, ps)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// ArchivesUpUnOpen  get archive unopen aids by mid.
func (d *Dao) ArchivesUpUnOpen(c context.Context, mid int64, offset int, ps int) (aids []int64, err error) {
	rows, err := d.rddb.Query(c, _arcUpUnOpenSQL, mid, offset, ps)
	if err != nil {
		log.Error("db.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// ArchiveAllUpCount  get all archive count by mid.
func (d *Dao) ArchiveAllUpCount(c context.Context, mid int64) (count int64, err error) {
	row := d.rddb.QueryRow(c, _arcUpAllCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ArchiveOpenUpCount  get open archive count by mid.
func (d *Dao) ArchiveOpenUpCount(c context.Context, mid int64) (count int64, err error) {
	row := d.rddb.QueryRow(c, _arcUpOpenCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// ArchiveUnOpenUpCount  get un open archive count by mid.
func (d *Dao) ArchiveUnOpenUpCount(c context.Context, mid int64) (count int64, err error) {
	row := d.rddb.QueryRow(c, _arcUpUnOpenCountSQL, mid)
	if err = row.Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// SimpleArchive get a archive by avid.
func (d *Dao) SimpleArchive(c context.Context, aid int64) (a *archive.SimpleArchive, err error) {
	row := d.rddb.QueryRow(c, _simpleArcSQL, aid)
	a = &archive.SimpleArchive{}
	if err = row.Scan(&a.Aid, &a.Title, &a.Mid); err != nil {
		if err == xsql.ErrNoRows {
			a = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// Recos fn
func (d *Dao) Recos(c context.Context, aid int64) (aids []int64, err error) {
	rows, err := d.db.Query(c, _getRecoSQL, aid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", aid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aid int64
		if err = rows.Scan(&aid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// RecoUpdate fn
func (d *Dao) RecoUpdate(c context.Context, aid int64, recoIDs []int64) (effCnt int64, err error) {
	if len(recoIDs) == 0 {
		_, err = d.db.Query(c, _delRecoSQL, aid)
		if err != nil {
			log.Error("d.db.Query error(%v)| aid (%d)", err, aid)
			return
		}
		return
	}
	_, err = d.db.Query(c, _delRecoSQL, aid)
	if err != nil {
		log.Error("d.db.Query(%d) error(%v)", aid, err)
		return
	}
	var (
		batchVals = make([]string, 0, len(recoIDs))
	)
	for _, recoID := range recoIDs {
		batchVals = append(batchVals, fmt.Sprintf("(%d,%d)", aid, recoID))
	}
	res, err := d.db.Exec(c, fmt.Sprintf(_batchAddRecoSQL, strings.Join(batchVals, ",")))
	if err != nil {
		log.Error("d.db.Exe _batchAddRecoSQL batchVals(%+v) error(%+v)", batchVals, err)
		return
	}
	return res.RowsAffected()
}

// UpMissionID update mission_id for  archive.
func (d *Dao) UpMissionID(c context.Context, aa *archive.ArcMissionParam) (rows int64, err error) {
	res, err := d.db.Exec(c, _inAddMissionSQL, aa.AID, aa.MissionID, aa.MissionID)
	if err != nil {
		log.Error("UpMissionID.Exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// RejectedArchives list rejected archives
func (d *Dao) RejectedArchives(c context.Context, mid int64, state, offset, limit int32, start *time.Time) (arcs []*archive.Archive, count int32, err error) {
	row := d.slaveDB.QueryRow(c, _rejectArcsCountSQL, mid, state, start)
	if err = row.Scan(&count); err != nil {
		log.Error("rows.Scan error(%v)", err)
		return
	}
	if count == 0 {
		return
	}
	rows, err := d.slaveDB.Query(c, _rejectArcsSQL, mid, state, start, offset, limit)
	if err != nil {
		log.Error("db.Query error(%v)", err)
		return
	}
	for rows.Next() {
		a := archive.Archive{}
		if err = rows.Scan(&a.Aid, &a.Mid, &a.Title, &a.RejectReason, &a.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		arcs = append(arcs, &a)
	}
	return
}
