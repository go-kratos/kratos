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
	columnRules = "id, area, limit_type, limit_scope, dur_sec, allowed_counts, ctime, mtime"

	selectRuleCountsSQL                 = `SELECT COUNT(1) FROM rate_limit_rules %s`
	selectRulesByCondSQL                = `SELECT ` + columnRules + ` FROM rate_limit_rules %s`
	selectRuleByIDsSQL                  = `SELECT ` + columnRules + ` FROM rate_limit_rules WHERE id IN(%s)`
	selectRulesByAreaSQL                = `SELECT ` + columnRules + ` FROM rate_limit_rules WHERE area = %s`
	selectRulesByAreaAndTypeSQL         = `SELECT ` + columnRules + ` FROM rate_limit_rules WHERE area = %s AND limit_type = %s`
	selectRulesByAreaAndTypeAndScopeSQL = `SELECT ` + columnRules + ` FROM rate_limit_rules WHERE area = %s AND limit_type = %s AND limit_scope = %s`

	insertRuleSQL = `INSERT INTO rate_limit_rules(area, limit_type, limit_scope, dur_sec, allowed_counts) VALUES(?, ?, ?, ?, ?)`
	updateRuleSQL = `UPDATE rate_limit_rules SET dur_sec = ?, allowed_counts = ?, mtime = ? WHERE area = ? AND limit_type = ? AND limit_scope = ?`
)

// Rule .
type Rule struct {
	ID            int64 `db:"id"`
	Area          int   `db:"area"`
	LimitType     int   `db:"limit_type"`
	LimitScope    int   `db:"limit_scope"`
	DurationSec   int64 `db:"dur_sec"`
	AllowedCounts int64 `db:"allowed_counts"`

	CTime time.Time `db:"ctime"`
	MTime time.Time `db:"mtime"`
}

// RuleDaoImpl .
type RuleDaoImpl struct{}

const (
	// LimitTypeDefaultLimit .
	LimitTypeDefaultLimit int = iota
	// LimitTypeRestrictLimit .
	LimitTypeRestrictLimit
	// LimitTypeWhite .
	LimitTypeWhite
	// LimitTypeBlack .
	LimitTypeBlack
)

const (
	// LimitScopeGlobal .
	LimitScopeGlobal int = iota
	// LimitScopeLocal .
	LimitScopeLocal
)

// NewRuleDao .
func NewRuleDao() *RuleDaoImpl {
	return &RuleDaoImpl{}
}

func updateRule(ctx context.Context, executer Executer, r *Rule) error {
	_, err := executer.Exec(ctx,
		updateRuleSQL,

		r.DurationSec,
		r.AllowedCounts,
		time.Now(),

		r.Area,
		r.LimitType,
		r.LimitScope,
	)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

func insertRule(ctx context.Context, executer Executer, r *Rule) error {
	res, err := executer.Exec(ctx,
		insertRuleSQL,

		r.Area,
		r.LimitType,
		r.LimitScope,
		r.DurationSec,
		r.AllowedCounts,
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

// GetByCond .
func (*RuleDaoImpl) GetByCond(ctx context.Context, cond *Condition) (rules []*Rule, totalCounts int64, err error) {
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
		queryCountsSQL := fmt.Sprintf(selectRuleCountsSQL, optionSQL)
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
	querySQL := fmt.Sprintf(selectRulesByCondSQL, optionSQL)
	log.Info("OptionSQL(%s), GetByCondSQL(%s)", optionSQL, querySQL)
	rules, err = queryRules(ctx, db, querySQL)
	if err != nil {
		return nil, totalCounts, err
	}
	return rules, totalCounts, nil
}

// Update .
func (rdi *RuleDaoImpl) Update(ctx context.Context, r *Rule) (*Rule, error) {
	if err := updateRule(ctx, db, r); err != nil {
		return nil, err
	}
	return rdi.GetByAreaAndTypeAndScope(ctx, &Condition{
		Area:       fmt.Sprintf("%d", r.Area),
		LimitType:  fmt.Sprintf("%d", r.LimitType),
		LimitScope: fmt.Sprintf("%d", r.LimitScope),
	})
}

// Insert .
func (rdi *RuleDaoImpl) Insert(ctx context.Context, r *Rule) (*Rule, error) {
	if err := insertRule(ctx, db, r); err != nil {
		return nil, err
	}
	return rdi.GetByID(ctx, r.ID)
}

// GetByID .
func (rdi *RuleDaoImpl) GetByID(ctx context.Context, id int64) (*Rule, error) {
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
func (*RuleDaoImpl) GetByIDs(ctx context.Context, ids []int64) ([]*Rule, error) {
	rs, err := queryRules(ctx, db, fmt.Sprintf(selectRuleByIDsSQL, util.IntSliToSQLVarchars(ids)))
	if err != nil {
		return nil, err
	}
	res := make([]*Rule, len(ids))
	for i, id := range ids {
		for _, r := range rs {
			if r.ID == id {
				res[i] = r
			}
		}
	}
	return res, nil
}

// GetByAreaAndLimitType .
func (*RuleDaoImpl) GetByAreaAndLimitType(ctx context.Context, cond *Condition) ([]*Rule, error) {
	return queryRules(ctx, db, fmt.Sprintf(selectRulesByAreaAndTypeSQL, cond.Area, cond.LimitType))
}

// GetByAreaAndTypeAndScope .
func (*RuleDaoImpl) GetByAreaAndTypeAndScope(ctx context.Context, cond *Condition) (*Rule, error) {
	rs, err := queryRules(ctx, db, fmt.Sprintf(selectRulesByAreaAndTypeAndScopeSQL,
		cond.Area,
		cond.LimitType,
		cond.LimitScope,
	))
	if err != nil {
		return nil, err
	}
	return rs[0], nil
}

// GetByArea .
func (*RuleDaoImpl) GetByArea(ctx context.Context, cond *Condition) ([]*Rule, error) {
	return queryRules(ctx, db, fmt.Sprintf(selectRulesByAreaSQL, cond.Area))
}

func queryRules(ctx context.Context, q Querier, rawSQL string) ([]*Rule, error) {
	log.Info("Query sql: %q", rawSQL)
	rows, err := q.Query(ctx, rawSQL)
	if err == sql.ErrNoRows {
		err = ErrResourceNotExist
	}
	if err != nil {
		log.Error("Error: %v, RawSQL: %s", err, rawSQL)
		return nil, err
	}
	defer rows.Close()

	rs, err := mapRowToRules(rows)
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 {
		return nil, ErrResourceNotExist
	}
	return rs, nil
}

func mapRowToRules(rows *sql.Rows) (rs []*Rule, err error) {
	for rows.Next() {
		r := Rule{}
		err = rows.Scan(
			&r.ID,
			&r.Area,
			&r.LimitType,
			&r.LimitScope,
			&r.DurationSec,
			&r.AllowedCounts,
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
