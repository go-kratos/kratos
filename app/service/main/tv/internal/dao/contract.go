package dao

import (
	"context"
	"database/sql"

	"go-common/app/service/main/tv/internal/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_getUserContractByMid        = "SELECT `id`, `mid`, `contract_id`, `order_no`, `is_deleted`, `ctime`, `mtime` FROM `tv_user_contract` WHERE `mid`=? AND `is_deleted`=0"
	_getUserContractByContractId = "SELECT `id`, `mid`, `contract_id`, `order_no`, `is_deleted`, `ctime`, `mtime` FROM `tv_user_contract` WHERE `contract_id`=? AND `is_deleted`=0"

	_deleteUserContract = "UPDATE `tv_user_contract` SET `is_deleted`=1 WHERE `id`=?"

	_insertUserContract = "INSERT INTO tv_user_contract (`mid`, `contract_id`, `order_no`) VALUES (?,?,?)"
)

// UserContractByMid quires one row from tv_user_contract.
func (d *Dao) UserContractByMid(c context.Context, mid int64) (uc *model.UserContract, err error) {
	row := d.db.QueryRow(c, _getUserContractByMid, mid)
	uc = &model.UserContract{}
	err = row.Scan(&uc.ID, &uc.Mid, &uc.ContractId, &uc.OrderNo, &uc.IsDeleted, &uc.Ctime, &uc.Mtime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error("rows.Scan(%s) error(%v)", _getUserContractByMid, err)
		err = errors.WithStack(err)
		return nil, err
	}
	return uc, nil
}

// UserContractByContractId quires one row from tv_user_contract.
func (d *Dao) UserContractByContractId(c context.Context, contractId string) (uc *model.UserContract, err error) {
	row := d.db.QueryRow(c, _getUserContractByContractId, contractId)
	uc = &model.UserContract{}
	err = row.Scan(&uc.ID, &uc.Mid, &uc.ContractId, &uc.OrderNo, &uc.IsDeleted, &uc.Ctime, &uc.Mtime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error("rows.Scan(%s) error(%v)", _getUserContractByContractId, err)
		err = errors.WithStack(err)
		return nil, err
	}
	return uc, nil
}

// TxDeleteUserContract deletes one user contract record.
func (d *Dao) TxDeleteUserContract(ctx context.Context, tx *xsql.Tx, id int32) (err error) {
	if _, err = tx.Exec(_deleteUserContract, id); err != nil {
		log.Error("rows.Scan(%s) error(%v)", _deleteUserContract, err)
		err = errors.WithStack(err)
		return
	}
	return
}

// TxInsertUserContract insert one row into tv_user_contract.
func (d *Dao) TxInsertUserContract(ctx context.Context, tx *xsql.Tx, uc *model.UserContract) (id int64, err error) {
	var (
		res sql.Result
	)
	if res, err = tx.Exec(_insertUserContract, uc.Mid, uc.ContractId, uc.OrderNo); err != nil {
		log.Error("d.TxInsertUserContract(%+v) err(%+v)", uc, err)
		err = errors.WithStack(err)
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Error("d.TxInsertUserContract(%+v) err(%+v)", uc, err)
		err = errors.WithStack(err)
		return
	}
	return
}
