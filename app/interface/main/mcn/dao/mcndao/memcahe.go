package mcndao

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	//cache: -nullcache=&mcnmodel.McnSign{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	McnSign(c context.Context, mcnmid int64) (up *mcnmodel.McnSign, err error)

	//cache: -nullcache=&mcnmodel.McnGetDataSummaryReply{IsNull:true} -check_null_code=$!=nil&&$.IsNull
	McnDataSummary(c context.Context, mcnmid int64, generateDate time.Time) (res *mcnmodel.McnGetDataSummaryReply, err error)

	//cache: -nullcache=&mcnmodel.UpPermissionCache{IsNull:true} -check_null_code=$!=nil&&$.IsNull
	UpPermission(c context.Context, signID int64, mid int64) (data *mcnmodel.UpPermissionCache, err error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	//mc: -key=mcnSignCacheKey -expire=d.mcnSignExpire -encode=json
	AddCacheMcnSign(c context.Context, mcnmid int64, up *mcnmodel.McnSign) (err error)
	//mc: -key=mcnSignCacheKey
	CacheMcnSign(c context.Context, mcnmid int64) (up *mcnmodel.McnSign, err error)
	//mc: -key=mcnSignCacheKey
	DelCacheMcnSign(c context.Context, mcnmid int64) (err error)

	//mc: -key=mcnDataCacheKey -expire=d.mcnDataExpire -encode=json
	AddCacheMcnDataSummary(c context.Context, mcnmid int64, data *mcnmodel.McnGetDataSummaryReply, generateDate time.Time) (err error)
	//mc: -key=mcnDataCacheKey
	CacheMcnDataSummary(c context.Context, mcnmid int64, generateDate time.Time) (data *mcnmodel.McnGetDataSummaryReply, err error)
	//mc: -key=mcnDataCacheKey
	DelMcnDataSummary(c context.Context, mcnmid int64, generateDate time.Time) (err error)

	//mc: -key=mcnPublicationPriceKey -expire=0 -encode=json
	AddCachePublicationPrice(c context.Context, signID int64, data *mcnmodel.PublicationPriceCache, mid int64) (err error)
	//mc: -key=mcnPublicationPriceKey
	CachePublicationPrice(c context.Context, signID int64, mid int64) (data *mcnmodel.PublicationPriceCache, err error)
	//mc: -key=mcnPublicationPriceKey
	DelPublicationPrice(c context.Context, signID int64, mid int64) (err error)

	//mc: -key=mcnUpPermissionKey -expire=d.mcnSignExpire  -encode=json
	AddCacheUpPermission(c context.Context, signID int64, data *mcnmodel.UpPermissionCache, mid int64) (err error)
	//mc: -key=mcnUpPermissionKey
	CacheUpPermission(c context.Context, signID int64, mid int64) (data *mcnmodel.UpPermissionCache, err error)
	//mc: -key=mcnUpPermissionKey
	DelUpPermission(c context.Context, signID int64, mid int64) (err error)
}

func mcnSignCacheKey(mcnmid int64) string {
	return fmt.Sprintf("mcn_s_%d", mcnmid)
}

//RawMcnSign raw db .
func (d *Dao) RawMcnSign(c context.Context, mcnmid int64) (up *mcnmodel.McnSign, err error) {
	up, _, err = d.GetMcnSignState("*", mcnmid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
		log.Error("error get state, err=%s", err)
		return
	}
	return
}

//AsyncDelCacheMcnSign delete in async way
func (d *Dao) AsyncDelCacheMcnSign(id int64) (err error) {
	return d.cache.Do(context.Background(), func(c context.Context) {
		d.DelCacheMcnSign(c, id)
	})
}

func mcnDataCacheKey(signID int64, generateDate time.Time) string {
	var key = fmt.Sprintf("mcn_data_%d_%s", signID, generateDate.Format(dateFmt))
	return key
}

//RawMcnDataSummary raw get
func (d *Dao) RawMcnDataSummary(c context.Context, signID int64, generateDate time.Time) (res *mcnmodel.McnGetDataSummaryReply, err error) {
	return d.GetMcnDataSummaryWithDiff(signID, mcnmodel.McnDataTypeDay, generateDate)
}

func mcnPublicationPriceKey(signID int64, mid int64) string {
	return fmt.Sprintf("mcn_pubprice_%d_%d", signID, mid)
}

func mcnUpPermissionKey(signID int64, mid int64) string {
	var s = fmt.Sprintf("mcn_upperm_%d_%d", signID, mid)
	log.Info("key=%s", s)
	return s
}

//RawUpPermission get permissino from db
func (d *Dao) RawUpPermission(c context.Context, signID int64, mid int64) (res *mcnmodel.UpPermissionCache, err error) {
	upList, err := d.GetUpBind("up_mid=? and sign_id=? and state in (?)", mid, signID, UpSignedStates)
	if err != nil {
		log.Error("get up from db fail, err=%v, signid=%d, mid=%d", err, signID, mid)
		return
	}
	if len(upList) == 0 {
		log.Warn("up not found, sign_id=%d, mid=%d", signID, mid)
		return
	}
	res = new(mcnmodel.UpPermissionCache)
	res.Permission = upList[0].Permission
	return
}
