package reply

import (
	"context"

	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_selBussinessSQL = "SELECT type, alias, appkey FROM business WHERE state=0"
)

// BusinessDao business dao.
type BusinessDao struct {
	db *sql.DB
}

// NewBusinessDao new BusinessDao and return.
func NewBusinessDao(db *sql.DB) (dao *BusinessDao) {
	dao = &BusinessDao{
		db: db,
	}
	return
}

// ListBusiness gets all business records
func (dao *BusinessDao) ListBusiness(c context.Context) (business []*reply.Business, err error) {
	rows, err := dao.db.Query(c, _selBussinessSQL)
	if err != nil {
		log.Error("sql query error(%v)", err)
		return
	}
	defer rows.Close()
	business = make([]*reply.Business, 0)
	for rows.Next() {
		b := new(reply.Business)
		if err = rows.Scan(&b.Type, &b.Alias, &b.Appkey); err != nil {
			return
		}
		business = append(business, b)
	}
	err = rows.Err()
	return
}
