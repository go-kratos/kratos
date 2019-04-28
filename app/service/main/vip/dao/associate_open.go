package dao

import (
	"context"

	"go-common/app/service/main/vip/model"
	"go-common/library/database/sql"

	"github.com/pkg/errors"
)

const (
	_openInfoOpenIDSQL            = "SELECT mid,open_id,app_id FROM oauth_open_id WHERE open_id = ? AND app_id = ?;"
	_openInfoByMidSQL             = "SELECT mid,open_id,app_id FROM oauth_open_id WHERE mid = ? AND app_id = ?;"
	_bindByOutOpenIDSQL           = "SELECT id,mid,app_id,out_open_id,state,ver FROM oauth_user_bind WHERE out_open_id = ? AND  app_id =?;"
	_bindByMidSQL                 = "SELECT id,mid,app_id,out_open_id,state,ver FROM oauth_user_bind WHERE mid = ? AND app_id =?;"
	_addBindSQL                   = "INSERT IGNORE INTO oauth_user_bind (mid,app_id,out_open_id)VALUES(?,?,?);"
	_updateBindByOutOpenID        = "UPDATE oauth_user_bind SET mid = ?,app_id = ?,out_open_id = ?,ver = ver + 1 WHERE out_open_id = ? AND app_id = ? AND ver = ?;"
	_updateBindStateSQL           = "UPDATE oauth_user_bind SET state = ?,ver = ver + 1 WHERE mid = ? AND app_id = ?  AND ver = ?;"
	_addOpenInfoSQL               = "INSERT IGNORE INTO oauth_open_id (mid,open_id,app_id)VALUES(?,?,?);"
	_bindInfoByOutOpenIDAndMidSQL = "SELECT mid,app_id,out_open_id,state,ver FROM oauth_user_bind WHERE mid = ? AND out_open_id = ? AND app_id =?;"
	_deleteBindInfoSQL            = "DELETE FROM oauth_user_bind WHERE id = ?;"
)

//RawOpenInfoByOpenID get info by open id.
func (d *Dao) RawOpenInfoByOpenID(c context.Context, openID string, appID int64) (res *model.OpenInfo, err error) {
	res = new(model.OpenInfo)
	if err = d.db.QueryRow(c, _openInfoOpenIDSQL, openID, appID).
		Scan(&res.Mid, &res.OpenID, &res.AppID); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao associate openinfo openid(%s,%d)", openID, appID)
	}
	return
}

// OpenInfoByMid get open info by mid.
func (d *Dao) OpenInfoByMid(c context.Context, mid int64, appID int64) (res *model.OpenInfo, err error) {
	res = new(model.OpenInfo)
	if err = d.db.QueryRow(c, _openInfoByMidSQL, mid, appID).
		Scan(&res.Mid, &res.OpenID, &res.AppID); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao associate openinfo bymid(%d,%d)", mid, appID)
	}
	return
}

//ByOutOpenID get open bind info by out_open_id.
func (d *Dao) ByOutOpenID(c context.Context, outOpenID string, appID int64) (res *model.OpenBindInfo, err error) {
	res = new(model.OpenBindInfo)
	if err = d.db.QueryRow(c, _bindByOutOpenIDSQL, outOpenID, appID).
		Scan(&res.ID, &res.Mid, &res.AppID, &res.OutOpenID, &res.State, &res.Ver); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao associate by out_open_id (%s,%d)", outOpenID, appID)
	}
	return
}

//RawBindInfoByMid get open bind info by mid.
func (d *Dao) RawBindInfoByMid(c context.Context, mid int64, appID int64) (res *model.OpenBindInfo, err error) {
	res = new(model.OpenBindInfo)
	if err = d.db.QueryRow(c, _bindByMidSQL, mid, appID).
		Scan(&res.ID, &res.Mid, &res.AppID, &res.OutOpenID, &res.State, &res.Ver); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao associate by mid (%d,%d)", mid, appID)
	}
	return
}

// TxAddBind add bind.
func (d *Dao) TxAddBind(tx *sql.Tx, arg *model.OpenBindInfo) (err error) {
	if _, err = tx.Exec(_addBindSQL, arg.Mid, arg.AppID, arg.OutOpenID); err != nil {
		err = errors.Wrapf(err, "dao add bind(%+v)", arg)
	}
	return
}

// TxUpdateBindByOutOpenID update bind by out_open_id.
func (d *Dao) TxUpdateBindByOutOpenID(tx *sql.Tx, arg *model.OpenBindInfo) (err error) {
	if _, err = tx.Exec(_updateBindByOutOpenID, arg.Mid, arg.AppID, arg.OutOpenID, arg.OutOpenID, arg.AppID, arg.Ver); err != nil {
		err = errors.Wrapf(err, "dao update bind by out_open_id(%+v)", arg)
	}
	return
}

// UpdateBindState update bind state.
func (d *Dao) UpdateBindState(c context.Context, arg *model.OpenBindInfo) (err error) {
	if _, err = d.db.Exec(c, _updateBindStateSQL, arg.State, arg.Mid, arg.AppID, arg.Ver); err != nil {
		err = errors.Wrapf(err, "dao update bind state(%+v)", arg)
	}
	return
}

// AddOpenInfo add open info.
func (d *Dao) AddOpenInfo(c context.Context, a *model.OpenInfo) (err error) {
	if _, err = d.db.Exec(c, _addOpenInfoSQL, a.Mid, a.OpenID, a.AppID); err != nil {
		err = errors.Wrapf(err, "dao insert open infp state(%+v)", a)
	}
	return
}

//BindInfoByOutOpenIDAndMid get bind info by out_open_id AND mid.
func (d *Dao) BindInfoByOutOpenIDAndMid(c context.Context, mid int64, outOpenID string, appID int64) (res *model.OpenBindInfo, err error) {
	res = new(model.OpenBindInfo)
	if err = d.db.QueryRow(c, _bindInfoByOutOpenIDAndMidSQL, mid, outOpenID, appID).
		Scan(&res.Mid, &res.AppID, &res.OutOpenID, &res.State, &res.Ver); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			err = nil
			return
		}
		err = errors.Wrapf(err, "dao associate by out_open_id (%s,%d)", outOpenID, appID)
	}
	return
}

// TxDeleteBindInfo bind info.
func (d *Dao) TxDeleteBindInfo(tx *sql.Tx, id int64) (err error) {
	if _, err = tx.Exec(_deleteBindInfoSQL, id); err != nil {
		err = errors.Wrapf(err, "dao delete bind(%d)", id)
	}
	return
}
