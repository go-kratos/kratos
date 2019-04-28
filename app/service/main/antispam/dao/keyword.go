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
	columnKeywords = "id, area, content, regexp_name, tag, hit_counts, state, origin_content, ctime, mtime"

	selectKeywordCountsSQL            = `SELECT COUNT(1) FROM keywords %s`
	selectKeywordsByCondSQL           = `SELECT ` + columnKeywords + ` FROM keywords %s`
	selectKeywordByIDsSQL             = `SELECT ` + columnKeywords + ` FROM keywords WHERE id IN (%s)`
	selectKeywordNeedRecycledSQL      = `SELECT ` + columnKeywords + ` FROM keywords FORCE INDEX(ix_ctime) WHERE state = %s AND hit_counts < %s AND tag IN(%s) AND ctime BETWEEN '%s' AND '%s' LIMIT %d`
	selectKeywordByOffsetLimitSQL     = `SELECT ` + columnKeywords + ` FROM keywords WHERE area = %s AND id > %s AND tag IN(%s) AND state = 0 LIMIT %s`
	selectKeywordByAreaAndContentsSQL = `SELECT ` + columnKeywords + ` FROM keywords WHERE area = %s AND content IN(%s)`

	insertKeywordSQL      = `INSERT INTO keywords(area, content, regexp_name, tag, hit_counts, origin_content) VALUES(?, ?, ?, ?, ?, ?)`
	updateKeywordSQL      = `UPDATE keywords SET content = ?, regexp_name = ?, tag = ?, hit_counts = ?, state = ?, origin_content = ?, ctime = ?, mtime = ? WHERE id = ?`
	deleteKeywordByIDsSQL = `UPDATE keywords SET state = 1, hit_counts = 0, mtime = ? WHERE id IN (%s)`
)

const (
	// KeywordTagDefaultLimit .
	KeywordTagDefaultLimit int = iota
	// KeywordTagRestrictLimit .
	KeywordTagRestrictLimit
	// KeywordTagWhite .
	KeywordTagWhite
	// KeywordTagBlack .
	KeywordTagBlack
)

// KeywordDaoImpl .
type KeywordDaoImpl struct{}

// Keyword .
type Keyword struct {
	ID            int64  `db:"id"`
	Area          int    `db:"area"`
	Tag           int    `db:"tag"`
	State         int    `db:"state"`
	HitCounts     int64  `db:"hit_counts"`
	RegexpName    string `db:"regexp_name"`
	Content       string `db:"content"`
	OriginContent string `db:"origin_content"`

	CTime time.Time `db:"ctime"`
	MTime time.Time `db:"mtime"`
}

// NewKeywordDao .
func NewKeywordDao() *KeywordDaoImpl {
	return &KeywordDaoImpl{}
}

// GetRubbish .
func (*KeywordDaoImpl) GetRubbish(ctx context.Context, cond *Condition) (keywords []*Keyword, err error) {
	querySQL := fmt.Sprintf(selectKeywordNeedRecycledSQL,
		cond.State,
		cond.HitCounts,
		util.StrSliToSQLVarchars(cond.Tags),
		cond.StartTime,
		cond.EndTime,
		cond.PerPage,
	)
	log.Info("get rubbish keywords rawSQL: %s", querySQL)
	ks, err := queryKeywords(ctx, db, querySQL)
	if err != nil {
		return nil, err
	}
	return ks, nil
}

// GetByOffsetLimit .
func (*KeywordDaoImpl) GetByOffsetLimit(ctx context.Context, cond *Condition) (keywords []*Keyword, err error) {
	return queryKeywords(ctx, db, fmt.Sprintf(selectKeywordByOffsetLimitSQL, cond.Area,
		cond.Offset, util.StrSliToSQLVarchars(cond.Tags), cond.Limit))
}

// GetByCond .
func (*KeywordDaoImpl) GetByCond(ctx context.Context, cond *Condition) (keywords []*Keyword, totalCounts int64, err error) {
	sqlConds := make([]string, 0)
	if cond.Search != "" {
		sqlConds = append(sqlConds, fmt.Sprintf("content LIKE '%%%s%%'", cond.Search))
	}
	if len(cond.Contents) > 0 {
		sqlConds = append(sqlConds, fmt.Sprintf("content IN (%s)", util.StrSliToSQLVarchars(cond.Tags)))
	}
	if cond.LastModifiedTime != "" {
		sqlConds = append(sqlConds, fmt.Sprintf("mtime >= '%s'", cond.LastModifiedTime))
		cond.OrderBy = ""
	}
	if cond.StartTime != "" || cond.EndTime != "" {
		if cond.StartTime != "" && cond.EndTime != "" {
			sqlConds = append(sqlConds, fmt.Sprintf("ctime BETWEEN '%s' AND '%s'", cond.StartTime, cond.EndTime))
		} else if cond.StartTime != "" {
			sqlConds = append(sqlConds, fmt.Sprintf("ctime >= '%s'", cond.StartTime))
		} else {
			sqlConds = append(sqlConds, fmt.Sprintf("ctime <= '%s'", cond.EndTime))
		}
	}
	if cond.State != "" {
		sqlConds = append(sqlConds, fmt.Sprintf("state = %s", cond.State))
	}
	if cond.Area != "" {
		sqlConds = append(sqlConds, fmt.Sprintf("area = %s", cond.Area))
	}
	if len(cond.Tags) > 0 {
		sqlConds = append(sqlConds, fmt.Sprintf("tag IN(%s)", util.StrSliToSQLVarchars(cond.Tags)))
	}

	var optionSQL string
	if len(sqlConds) > 0 {
		optionSQL = fmt.Sprintf("WHERE %s", strings.Join(sqlConds, " AND "))
	}

	var limitSQL string
	if cond.Pagination != nil {
		queryCountsSQL := fmt.Sprintf(selectKeywordCountsSQL, optionSQL)
		log.Info("queryCounts sql: %s", queryCountsSQL)
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
	querySQL := fmt.Sprintf(selectKeywordsByCondSQL, optionSQL)
	log.Info("OptionSQL(%s), GetByCondSQL(%s)", optionSQL, querySQL)
	keywords, err = queryKeywords(ctx, db, querySQL)
	if err != nil {
		return nil, 0, err
	}
	if totalCounts == 0 {
		totalCounts = int64(len(keywords))
	}
	return keywords, totalCounts, nil
}

// GetByAreaAndContents .
func (*KeywordDaoImpl) GetByAreaAndContents(ctx context.Context,
	cond *Condition) ([]*Keyword, error) {
	querySQL := fmt.Sprintf(selectKeywordByAreaAndContentsSQL,
		cond.Area, util.StrSliToSQLVarchars(cond.Contents))
	ks, err := queryKeywords(ctx, db, querySQL)
	if err != nil {
		return nil, err
	}
	res := make([]*Keyword, len(cond.Contents))
	for i, c := range cond.Contents {
		for _, k := range ks {
			if strings.EqualFold(k.Content, c) {
				res[i] = k
			}
		}
	}
	return res, nil
}

// GetByAreaAndContent .
func (kdi *KeywordDaoImpl) GetByAreaAndContent(ctx context.Context,
	cond *Condition) (*Keyword, error) {
	ks, err := kdi.GetByAreaAndContents(ctx, cond)
	if err != nil {
		return nil, err
	}
	if ks[0] == nil {
		return nil, ErrResourceNotExist
	}
	return ks[0], nil
}

// Update .
func (kdi *KeywordDaoImpl) Update(ctx context.Context,
	k *Keyword) (*Keyword, error) {
	if err := updateKeyword(ctx, db, k); err != nil {
		return nil, err
	}
	return kdi.GetByID(ctx, k.ID)
}

// Insert .
func (kdi *KeywordDaoImpl) Insert(ctx context.Context, k *Keyword) (*Keyword, error) {
	if err := insertKeyword(ctx, db, k); err != nil {
		return nil, err
	}
	return kdi.GetByID(ctx, k.ID)
}

// DeleteByIDs .
func (kdi *KeywordDaoImpl) DeleteByIDs(ctx context.Context, ids []int64) ([]*Keyword, error) {
	if err := deleteKeywordByIDs(ctx, db, ids); err != nil {
		return nil, err
	}
	return kdi.GetByIDs(ctx, ids)
}

// GetByID .
func (kdi *KeywordDaoImpl) GetByID(ctx context.Context, id int64) (*Keyword, error) {
	ks, err := kdi.GetByIDs(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	if ks[0] == nil {
		return nil, ErrResourceNotExist
	}
	return ks[0], nil
}

// GetByIDs .
func (*KeywordDaoImpl) GetByIDs(ctx context.Context, ids []int64) ([]*Keyword, error) {
	ks, err := queryKeywords(ctx, db,
		fmt.Sprintf(selectKeywordByIDsSQL, util.IntSliToSQLVarchars(ids)))
	if err != nil {
		return nil, err
	}
	res := make([]*Keyword, len(ids))
	for i, id := range ids {
		for _, k := range ks {
			if k.ID == id {
				res[i] = k
			}
		}
	}
	return res, nil
}

func insertKeyword(ctx context.Context, executer Executer, k *Keyword) error {
	defaultHitCount := 1
	res, err := executer.Exec(ctx,
		insertKeywordSQL,

		k.Area,
		k.Content,
		k.RegexpName,
		k.Tag,
		defaultHitCount,
		k.OriginContent,
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
	k.ID = lastID
	return nil
}

func updateKeyword(ctx context.Context, executer Executer, k *Keyword) error {
	_, err := executer.Exec(ctx,
		updateKeywordSQL,

		k.Content,
		k.RegexpName,
		k.Tag,
		k.HitCounts,
		k.State,
		k.OriginContent,
		k.CTime,
		time.Now(),

		k.ID,
	)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	return nil
}

func deleteKeywordByIDs(ctx context.Context, executer Executer, ids []int64) error {
	rawSQL := fmt.Sprintf(deleteKeywordByIDsSQL, util.IntSliToSQLVarchars(ids))
	if _, err := executer.Exec(ctx, rawSQL, time.Now()); err != nil {
		log.Error("Error: %v, RawSQL: %s", err, rawSQL)
		return err
	}
	return nil
}

func queryKeywords(ctx context.Context, q Querier, rawSQL string) ([]*Keyword, error) {
	// NOTICE: this MotherFucker Query() will never return `ErrNoRows` when there is no rows found !
	rows, err := q.Query(ctx, rawSQL)
	if err == sql.ErrNoRows {
		return nil, ErrResourceNotExist
	} else if err != nil {
		log.Error("ctx: %+v, Error: %v, RawSQL: %s", ctx, err, rawSQL)
		return nil, err
	}
	defer rows.Close()

	log.Info("Query sql: %q", rawSQL)
	ks, err := mapRowToKeywords(rows)
	if err != nil {
		return nil, err
	}
	if len(ks) == 0 {
		return nil, ErrResourceNotExist
	}
	return ks, nil
}

func mapRowToKeywords(rows *sql.Rows) (ks []*Keyword, err error) {
	for rows.Next() {
		k := Keyword{}
		err = rows.Scan(
			&k.ID,
			&k.Area,
			&k.Content,
			&k.RegexpName,
			&k.Tag,
			&k.HitCounts,
			&k.State,
			&k.OriginContent,
			&k.CTime,
			&k.MTime,
		)
		if err != nil {
			log.Error("%v", err)
			return nil, err
		}
		ks = append(ks, &k)
	}
	if err = rows.Err(); err != nil {
		log.Error("%v", err)
		return nil, err
	}
	return ks, nil
}
