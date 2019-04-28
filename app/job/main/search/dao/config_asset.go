package dao

// import (
// 	"context"

// 	"go-common/app/job/main/search/model"
// 	"go-common/database/sql"
// 	"go-common/xstr"
// )

// const (
// 	_getAssetSQL = "SELECT id, name, type, src FROM digger_asset where id in (?)"
// )

// func (d *Dao) getAsset(c context.Context, ids []int64) (res *model.SQLAsset, err error) {
// 	res = new(model.SQLAsset)
// 	row := d.SearchDB.QueryRow(c, _getAssetSQL, xstr.JoinInts(ids))
// 	if err = row.Scan(&res.ID, &res.Name, &res.Type, &res.Src); err != nil {
// 		if err == sql.ErrNoRows {
// 			err = nil
// 			res = nil
// 		}
// 	}
// 	return
// }
