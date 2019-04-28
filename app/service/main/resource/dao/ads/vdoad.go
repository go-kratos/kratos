package ads

import (
	"context"
	"database/sql"
	"time"

	"go-common/app/service/main/resource/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_videoAdsSQL = `SELECT name,contract_id,aid,season_id,typeid,ad_cid,ad_strategy,ad_url,ad_order,skipable,note,agency_name,agency_country,
	agency_area,price,verified,state,front_aid,target,platform,type,user_set,play_count,mtime FROM video_ads WHERE state=0 AND verified=1 AND starttime<? AND endtime>? ORDER BY mtime,ctime ASC`
)

// VideoAds get video_ads
func (dao *Dao) VideoAds(c context.Context) (ads []*model.VideoAD, err error) {
	var (
		rows *xsql.Rows
		now  = time.Now()
	)
	if rows, err = dao.db.Query(c, _videoAdsSQL, now, now); err != nil {
		log.Error("dao.Exec(%v, %v), err (%v)", now, now, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ad := &model.VideoAD{}
		var (
			agencyCountry sql.NullInt64
			agencyArea    sql.NullInt64
			price         sql.NullFloat64
		)
		if err = rows.Scan(&ad.Name, &ad.ContractID, &ad.Aids, &ad.SeasonID, &ad.TypeID, &ad.AdCid, &ad.AdStrategy,
			&ad.AdURL, &ad.AdOrder, &ad.Skipable, &ad.Note, &ad.AgencyName, &agencyCountry, &agencyArea,
			&price, &ad.Verified, &ad.State, &ad.FrontAid, &ad.Target, &ad.Platform, &ad.Type, &ad.UserSet, &ad.PlayCount, &ad.MTime); err != nil {
			log.Error("rows.Scan(), err (%v)", err)
			return
		}
		ad.AgencyCountry = int(agencyCountry.Int64)
		ad.AgencyArea = int(agencyArea.Int64)
		ad.Price = float32(price.Float64)
		ads = append(ads, ad)
	}
	err = rows.Err()
	return
}
