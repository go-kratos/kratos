package dao

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go-common/app/job/main/search/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_getBusinessSQL = "SELECT business, app_ids, asset_db, asset_es, asset_dtb FROM digger_business WHERE business=?"
)

type bns struct {
	d        *Dao
	business string
	bInfo    *model.Bsn
}

func newBusiness(d *Dao, business string) (bs *bns, err error) {
	bs = &bns{
		d:        d,
		business: business,
		bInfo:    new(model.Bsn),
	}
	if err = bs.initBusiness(); err != nil {
		log.Error("d.initBusiness error (%v)", err)
	}
	return
}

func (bs *bns) initBusiness() (err error) {
	var sqlBusiness *model.SQLBusiness
	for {
		if sqlBusiness, err = bs.getBusiness(context.TODO()); err != nil {
			log.Error("initBusiness error (%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
	if sqlBusiness == nil {
		err = errors.New("initBusiness: " + bs.business + " not found in `digger_business`")
		return
	}
	bs.bInfo.Business = sqlBusiness.Business
	bs.bInfo.AppInfo = make([]model.BsnAppInfo, 0)
	// business-appinfo
	if sqlBusiness.AppIds != "" {
		err = json.Unmarshal([]byte(sqlBusiness.AppIds), &bs.bInfo.AppInfo)
	}
	// business-assetdb

	// business-assetes

	// business-assedtb

	return
}

func (bs *bns) getBusiness(c context.Context) (res *model.SQLBusiness, err error) {
	res = new(model.SQLBusiness)
	row := bs.d.SearchDB.QueryRow(c, _getBusinessSQL, bs.business)
	if err = row.Scan(&res.Business, &res.AppIds, &res.AssetDB, &res.AssetES, &res.AssetDtb); err != nil {
		log.Error("business row.Scan error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
			res = nil
			return
		}
	}
	return
}
