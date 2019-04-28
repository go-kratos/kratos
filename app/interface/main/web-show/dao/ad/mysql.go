package ad

import (
	"context"
	"time"

	"go-common/app/interface/main/web-show/model/resource"
	"go-common/library/log"
)

const (
	_selAds = `SELECT sch.id,sch.resource_id,so.title,so.pic,so.spic,so.url,so.atype FROM schedule as sch INNER JOIN order_applied as o ON sch.applied_id = o.id  
				INNER JOIN material as so ON so.id= o.material_id WHERE sch.stime<? AND sch.etime>? AND o.state =3`
)

// Ads return ads info
func (dao *Dao) Ads(c context.Context) (ads []*resource.Assignment, err error) {
	rows, err := dao.selAdsStmt.Query(c, time.Now(), time.Now())
	if err != nil {
		log.Error("dao.selAdsStmt() err(%v)", err)
		return
	}
	defer rows.Close()
	ads = make([]*resource.Assignment, 0)
	for rows.Next() {
		ad := &resource.Assignment{}
		if err = rows.Scan(&ad.ID, &ad.ResID, &ad.Name, &ad.Pic, &ad.LitPic, &ad.URL, &ad.Atype); err != nil {
			PromError("mysql.Ads", "rows.scan err(%v)", err)
			return
		}
		ads = append(ads, ad)
	}
	return
}
