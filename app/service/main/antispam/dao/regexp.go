package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/service/main/antispam/util"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	columnsRegexp = `id, admin_id, area, name, operation, content, state, ctime, mtime`

	selectRegexpCountsSQL           = `SELECT COUNT(1) FROM regexps %s`
	selectRegexpsByCondSQL          = `SELECT ` + columnsRegexp + ` FROM regexps %s`
	selectRegexpByIDsSQL            = `SELECT ` + columnsRegexp + ` FROM regexps WHERE id IN(%s)`
	selectRegexpByContentsSQL       = `SELECT ` + columnsRegexp + ` FROM regexps WHERE content IN(%s)`
	selectRegexpByAreaAndContentSQL = `SELECT ` + columnsRegexp + ` FROM regexps WHERE area = %s AND content IN(%s)`

	insertRegexpSQL = `INSERT INTO regexps(id, admin_id, area, name, operation, content, state) VALUES(?, ?, ?, ?, ?, ?, ?)`
	updateRegexpSQL = `UPDATE regexps SET admin_id = ?, name = ?, content = ?, operation = ?, state = ?, mtime = ? WHERE id = ?`
)

const (
	// OperationLimit .
	OperationLimit int = iota
	// OperationPutToWhiteList .
	OperationPutToWhiteList
	// OperationRestrictLimit .
	OperationRestrictLimit
	// OperationIgnore .
	OperationIgnore
)

// RegexpDaoImpl .
type RegexpDaoImpl struct{}

// Regexp .
type Regexp struct {
	ID        int64     `db:"id"`
	Area      int       `db:"area"`
	Name      string    `db:"name"`
	AdminID   int64     `db:"admin_id"`
	Operation int       `db:"operation"`
	Content   string    `db:"content"`
	State     int       `db:"state"`
	CTime     time.Time `db:"ctime"`
	MTime     time.Time `db:"mtime"`
}

// NewRegexpDao .
func NewRegexpDao() *RegexpDaoImpl {
	return &RegexpDaoImpl{}
}

// GetByCond .
func (*RegexpDaoImpl) GetByCond(ctx context.Context,
	cond *Condition) (regexps []*Regexp, totalCounts int64, err error) {
	sqlConds := make([]string, 0)
	if cond.Area != "" {
		sqlConds = append(sqlConds, fmt.Sprintf("area = %s", cond.Area))
	}
	if cond.State != "" {
		sqlConds = append(sqlConds, fmt.Sprintf("state = %s", cond.State))
	}
	var optionSQL string
	if len(sqlConds) > 0 {
		optionSQL = fmt.Sprintf("WHERE %s", strings.Join(sqlConds, " AND "))
	}

	var limitSQL string
	if cond.Pagination != nil {
		queryCountsSQL := fmt.Sprintf(selectRegexpCountsSQL, optionSQL)
		totalCounts, err = GetTotalCounts(ctx, db, queryCountsSQL)
		if err != nil {
			return nil, 0, err
		}
		offset, limit := cond.OffsetLimit(totalCounts)
		if limit == 0 {
			return nil, 0, ErrResourceNotExist
		}
		limitSQL = fmt.Sprintf("LIMIT %d, %d", offset, limit)
	}
	if cond.OrderBy != "" {
		optionSQL = fmt.Sprintf("%s ORDER BY %s %s", optionSQL, cond.OrderBy, cond.Order)
	}
	if limitSQL != "" {
		optionSQL = fmt.Sprintf("%s %s", optionSQL, limitSQL)
	}
	querySQL := fmt.Sprintf(selectRegexpsByCondSQL, optionSQL)
	log.Info("OptionSQL(%s), GetByCondSQL(%s)", optionSQL, querySQL)
	regexps, err = queryRegexps(ctx, db, querySQL)
	if err != nil {
		return nil, totalCounts, err
	}
	return regexps, totalCounts, nil
}

// Update .
func (rdi *RegexpDaoImpl) Update(ctx context.Context, r *Regexp) (*Regexp, error) {
	err := updateRegexp(ctx, db, r)
	if err != nil {
		return nil, err
	}
	return rdi.GetByID(ctx, r.ID)
}

// Insert .
func (rdi *RegexpDaoImpl) Insert(ctx context.Context, r *Regexp) (*Regexp, error) {
	err := insertRegexp(ctx, db, r)
	if err != nil {
		return nil, err
	}
	return rdi.GetByID(ctx, r.ID)
}

// GetByID .
func (rdi *RegexpDaoImpl) GetByID(ctx context.Context, id int64) (*Regexp, error) {
	rs, err := rdi.GetByIDs(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	if rs[0] == nil {
		return nil, ErrResourceNotExist
	}
	return rs[0], nil
}

// GetByIDs .
func (*RegexpDaoImpl) GetByIDs(ctx context.Context, ids []int64) ([]*Regexp, error) {
	rs, err := queryRegexps(ctx, db, fmt.Sprintf(selectRegexpByIDsSQL, util.IntSliToSQLVarchars(ids)))
	if err != nil {
		return nil, err
	}
	res := make([]*Regexp, len(ids))
	for i, id := range ids {
		for _, r := range rs {
			if r.ID == id {
				res[i] = r
			}
		}
	}
	return res, nil
}

// GetByContents .
func (*RegexpDaoImpl) GetByContents(ctx context.Context, contents []string) ([]*Regexp, error) {
	if len(contents) == 0 {
		log.Error("%v", ErrParams)
		return nil, ErrParams
	}
	rs, err := queryRegexps(ctx, db, fmt.Sprintf(selectRegexpByContentsSQL, util.StrSliToSQLVarchars(contents)))
	if err != nil {
		return nil, err
	}
	res := make([]*Regexp, len(contents))
	for i, c := range contents {
		for _, r := range rs {
			if strings.EqualFold(r.Content, c) {
				res[i] = r
			}
		}
	}
	return res, nil
}

// GetByAreaAndContent .
func (*RegexpDaoImpl) GetByAreaAndContent(ctx context.Context, cond *Condition) (*Regexp, error) {
	rs, err := queryRegexps(ctx, db, fmt.Sprintf(selectRegexpByAreaAndContentSQL,
		cond.Area, util.StrSliToSQLVarchars(cond.Contents)))
	if err != nil {
		return nil, err
	}
	return rs[0], nil
}

func insertRegexp(ctx context.Context, executer Executer, r *Regexp) error {
	res, err := executer.Exec(ctx, insertRegexpSQL,
		r.ID,
		r.AdminID,
		r.Area,
		r.Name,
		r.Operation,
		r.Content,
		r.State,
	)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Error("%v", err)
		return err
	}
	r.ID = lastID
	return nil
}

func updateRegexp(ctx context.Context, executer Executer, r *Regexp) error {
	_, err := executer.Exec(ctx, updateRegexpSQL,
		r.AdminID,
		r.Name,
		r.Content,
		r.Operation,
		r.State,
		time.Now(),

		r.ID,
	)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

func queryRegexps(ctx context.Context, q Querier, rawSQL string) ([]*Regexp, error) {
	rows, err := q.Query(ctx, rawSQL)
	if err == sql.ErrNoRows {
		err = ErrResourceNotExist
	}
	if err != nil {
		log.Error("Error: %v, sql: %s", err, rawSQL)
		return nil, err
	}
	defer rows.Close()

	rs, err := mapRowToRegexps(rows)
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 {
		log.Error("Error: %v, sql: %s", ErrResourceNotExist, rawSQL)
		return nil, ErrResourceNotExist
	}
	return rs, nil
}

func mapRowToRegexps(rows *sql.Rows) (rs []*Regexp, err error) {
	rs = make([]*Regexp, 0)
	for rows.Next() {
		r := Regexp{}
		err = rows.Scan(
			&r.ID,
			&r.AdminID,
			&r.Area,
			&r.Name,
			&r.Operation,
			&r.Content,
			&r.State,
			&r.CTime,
			&r.MTime,
		)
		if err != nil {
			log.Error("%v", err)
			return nil, err
		}
		rs = append(rs, &r)
	}
	if err = rows.Err(); err != nil {
		log.Error("%v", err)
		return nil, err
	}
	return rs, nil
}
