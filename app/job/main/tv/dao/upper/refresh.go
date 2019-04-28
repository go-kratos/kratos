package upper

import (
	"context"
	"fmt"

	ugcMdl "go-common/app/job/main/tv/model/ugc"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_countUpper   = "SELECT count(1) FROM ugc_uploader WHERE deleted = 0"
	_pickUppers   = "SELECT mid FROM ugc_uploader WHERE deleted = 0 AND mid > ? ORDER BY mid LIMIT 0,%d"
	_setUpperName = "UPDATE ugc_uploader SET ori_name = ?, cms_name = ? WHERE mid = ? AND deleted = 0"
	_setUpperFace = "UPDATE ugc_uploader SET ori_face = ?, cms_face = ? WHERE mid = ? AND deleted = 0"
	_sendUpper    = "UPDATE ugc_uploader SET submit = 1 WHERE mid = ? AND deleted = 0"
	_importUpper  = "REPLACE INTO ugc_uploader (ori_name, cms_name, ori_face, cms_face, mid) VALUES (?,?,?,?,?)"
	_upName       = 1
	_upFace       = 2
)

// CountUP counts the uppers
func (d *Dao) CountUP(c context.Context) (count int64, err error) {
	if err = d.DB.QueryRow(c, _countUpper).Scan(&count); err != nil {
		log.Error("d.CountUP.Query error(%v)", err)
	}
	return
}

// PickUppers picks data by Piece to refresh uppers
func (d *Dao) PickUppers(ctx context.Context, LastID int64, nbData int) (res []int64, myLast int64, err error) {
	var (
		rows  *sql.Rows
		query = fmt.Sprintf(_pickUppers, nbData)
	)
	if rows, err = d.DB.Query(ctx, query, LastID); err != nil {
		log.Error("d.refreshArcMC.Query: %s error(%v)", query, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var li int64
		if err = rows.Scan(&li); err != nil {
			log.Error("refreshArcMC row.Scan() error(%v)", err)
			return
		}
		res = append(res, li)
	}
	if err = rows.Err(); err != nil {
		log.Error("d.PickUppers.Query error(%v)", err)
		return
	}
	// record the max ID in this piece
	if len(res) > 0 {
		myLast = res[len(res)-1]
	}
	return
}

// SendUpper updates one upper's submit to 1
func (d *Dao) SendUpper(ctx context.Context, mid int64) (err error) {
	if _, err = d.DB.Exec(ctx, _sendUpper, mid); err != nil {
		log.Error("SendUpper Error: %v", mid, err)
	}
	return
}

// setUpperV updates the upper's name or face value
func (d *Dao) setUpperV(ctx context.Context, req *ugcMdl.ReqSetUp, oriUp *ugcMdl.Upper) (err error) {
	var sql string
	if req.UpType == _upName {
		sql = _setUpperName
		oriUp.CMSName = req.Value
		oriUp.OriName = req.Value
	} else if req.UpType == _upFace {
		sql = _setUpperFace
		oriUp.CMSFace = req.Value
		oriUp.OriFace = req.Value
	} else {
		return ecode.TvDangbeiWrongType
	}
	if _, err = d.DB.Exec(ctx, sql, req.Value, req.Value, req.MID); err != nil {
		log.Error("SetUpperName Error: %v", req.MID, err)
	}
	return
}

// ImportUp is used to import a new up when we manually add an archive
func (d *Dao) ImportUp(ctx context.Context, up *ugcMdl.EasyUp) (err error) {
	if _, err = d.DB.Exec(ctx, _importUpper, up.Name, up.Name, up.Face, up.Face, up.MID); err != nil {
		log.Error("ImportUp Error: %v", up.MID, err)
	}
	d.addUpMetaCache(ctx, up.ToUpper(nil)) // add this upper into the cache
	return
}

// RefreshUp refreshes the upper's name and face in both DB and cache
func (d *Dao) RefreshUp(ctx context.Context, req *ugcMdl.ReqSetUp) (err error) {
	var oriUp *ugcMdl.Upper
	if oriUp, err = d.LoadUpMeta(ctx, req.MID); err != nil || oriUp == nil {
		return
	}
	if err = d.setUpperV(ctx, req, oriUp); err != nil { // update DB
		return
	}
	d.addUpMetaCache(ctx, oriUp) // update cache
	return
}
