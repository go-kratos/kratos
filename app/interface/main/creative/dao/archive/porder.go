package archive

import (
	"context"
	"database/sql"
	"net/url"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_flow      = "/videoup/flows"
	_porderSQL = "SELECT id,aid,industry_id,brand_id,brand_name,official,show_type,advertiser,agent,ctime,mtime FROM archive_porder WHERE aid=? AND show_front = 1"
)

// Flows fn
func (d *Dao) Flows(c context.Context) (flows []*archive.Flow, err error) {
	params := url.Values{}
	flows = []*archive.Flow{}
	var res struct {
		Code int             `json:"code"`
		Data []*archive.Flow `json:"data"`
	}
	if err = d.client.Get(c, d.flow, "", params, &res); err != nil {
		log.Error("archive.Flow url(%s) error(%v)", d.flow+"?"+params.Encode(), err)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Code != 0 {
		log.Error("archive.Flow url(%s) res(%v)", d.flow+"?"+params.Encode(), res)
		err = ecode.CreativeArchiveAPIErr
		return
	}
	if res.Data == nil || len(res.Data) == 0 {
		return
	}
	flows = res.Data
	return
}

// Porder NOTE: move to up service
func (d *Dao) Porder(c context.Context, aid int64) (pd *archive.Porder, err error) {
	row := d.db.QueryRow(c, _porderSQL, aid)
	pd = &archive.Porder{}
	if err = row.Scan(&pd.ID, &pd.AID, &pd.IndustryID, &pd.BrandID, &pd.BrandName, &pd.Official, &pd.ShowType, &pd.Advertiser, &pd.Agent, &pd.Ctime, &pd.Mtime); err != nil {
		if err == sql.ErrNoRows {
			pd = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}
