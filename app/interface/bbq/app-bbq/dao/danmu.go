package dao

import (
	"context"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/log"
)

const (
	_queryDanmu = "select oid,mid,offset,content from bullet_content where id = ?"
)

//RawBullet ...
func (d *Dao) RawBullet(ctx context.Context, danmu int64) (v *model.Danmu, err error) {
	v = new(model.Danmu)
	if err = d.dbDM.QueryRow(ctx, _queryDanmu, danmu).Scan(&v.OID, &v.MID, &v.Offset, &v.Content); err != nil {
		log.Errorw(ctx, "content", "query danmu", "err", err, "danmu id", danmu)
		return
	}
	return
}
