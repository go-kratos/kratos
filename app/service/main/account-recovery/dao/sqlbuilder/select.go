// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sqlbuilder

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// NewSelectBuilder creates a new SELECT builder.
func NewSelectBuilder() *SelectBuilder {
	return DefaultFlavor.NewSelectBuilder()
}

func newSelectBuilder() *SelectBuilder {
	args := &Args{}
	return &SelectBuilder{
		Cond: Cond{
			Args: args,
		},
		limit:  -1,
		offset: -1,
		args:   args,
	}
}

// SelectBuilder is a builder to build SELECT.
type SelectBuilder struct {
	Cond

	distinct    bool
	tables      []string
	selectCols  []string
	whereExprs  []string
	havingExprs []string
	groupByCols []string
	orderByCols []string
	order       string
	limit       int
	offset      int

	args *Args
}

// Distinct marks this SELECT as DISTINCT.
func (sb *SelectBuilder) Distinct() *SelectBuilder {
	sb.distinct = true
	return sb
}

// Select sets columns in SELECT.
func (sb *SelectBuilder) Select(col ...string) *SelectBuilder {
	sb.selectCols = EscapeAll(col...)
	return sb
}

// From sets table names in SELECT.
func (sb *SelectBuilder) From(table ...string) *SelectBuilder {
	sb.tables = table
	return sb
}

// Where sets expressions of WHERE in SELECT.
func (sb *SelectBuilder) Where(andExpr ...string) *SelectBuilder {
	sb.whereExprs = append(sb.whereExprs, andExpr...)
	return sb
}

// Having sets expressions of HAVING in SELECT.
func (sb *SelectBuilder) Having(andExpr ...string) *SelectBuilder {
	sb.havingExprs = append(sb.havingExprs, andExpr...)
	return sb
}

// GroupBy sets columns of GROUP BY in SELECT.
func (sb *SelectBuilder) GroupBy(col ...string) *SelectBuilder {
	sb.groupByCols = EscapeAll(col...)
	return sb
}

// OrderBy sets columns of ORDER BY in SELECT.
func (sb *SelectBuilder) OrderBy(col ...string) *SelectBuilder {
	sb.orderByCols = EscapeAll(col...)
	return sb
}

// Asc sets order of ORDER BY to ASC.
func (sb *SelectBuilder) Asc() *SelectBuilder {
	sb.order = "ASC"
	return sb
}

// Desc sets order of ORDER BY to DESC.
func (sb *SelectBuilder) Desc() *SelectBuilder {
	sb.order = "DESC"
	return sb
}

// Limit sets the LIMIT in SELECT.
func (sb *SelectBuilder) Limit(limit int) *SelectBuilder {
	sb.limit = limit
	return sb
}

// Offset sets the LIMIT offset in SELECT.
func (sb *SelectBuilder) Offset(offset int) *SelectBuilder {
	sb.offset = offset
	return sb
}

// As returns an AS expression.
func (sb *SelectBuilder) As(col, alias string) string {
	return fmt.Sprintf("%v AS %v", col, Escape(alias))
}

// BuilderAs returns an AS expression wrapping a complex SQL.
// According to SQL syntax, SQL built by builder is surrounded by parens.
func (sb *SelectBuilder) BuilderAs(builder Builder, alias string) string {
	return fmt.Sprintf("(%v) AS %v", sb.Var(builder), Escape(alias))
}

// String returns the compiled SELECT string.
func (sb *SelectBuilder) String() string {
	s, _ := sb.Build()
	return s
}

// Build returns compiled SELECT string and args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (sb *SelectBuilder) Build() (sql string, args []interface{}) {
	return sb.BuildWithFlavor(sb.args.Flavor)
}

// BuildWithFlavor returns compiled SELECT string and args with flavor and initial args.
// They can be used in `DB#Query` of package `database/sql` directly.
func (sb *SelectBuilder) BuildWithFlavor(flavor Flavor, initialArg ...interface{}) (sql string, args []interface{}) {
	buf := &bytes.Buffer{}
	buf.WriteString("SELECT ")

	if sb.distinct {
		buf.WriteString("DISTINCT ")
	}

	buf.WriteString(strings.Join(sb.selectCols, ", "))
	buf.WriteString(" FROM ")
	buf.WriteString(strings.Join(sb.tables, ", "))

	if len(sb.whereExprs) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(sb.whereExprs, " AND "))
	}

	if len(sb.groupByCols) > 0 {
		buf.WriteString(" GROUP BY ")
		buf.WriteString(strings.Join(sb.groupByCols, ", "))

		if len(sb.havingExprs) > 0 {
			buf.WriteString(" HAVING ")
			buf.WriteString(strings.Join(sb.havingExprs, " AND "))
		}
	}

	if len(sb.orderByCols) > 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(sb.orderByCols, ", "))

		if sb.order != "" {
			buf.WriteRune(' ')
			buf.WriteString(sb.order)
		}
	}

	if sb.limit >= 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(sb.limit))

		if sb.offset >= 0 {
			buf.WriteString(" OFFSET ")
			buf.WriteString(strconv.Itoa(sb.offset))
		}
	}

	return sb.Args.CompileWithFlavor(buf.String(), flavor, initialArg...)
}

// SetFlavor sets the flavor of compiled sql.
func (sb *SelectBuilder) SetFlavor(flavor Flavor) (old Flavor) {
	old = sb.args.Flavor
	sb.args.Flavor = flavor
	return
}

// Copy the builder
func (sb *SelectBuilder) Copy() *SelectBuilder {
	return &SelectBuilder{
		Cond:        sb.Cond,
		distinct:    sb.distinct,
		tables:      sb.tables,
		selectCols:  sb.selectCols,
		whereExprs:  sb.whereExprs,
		havingExprs: sb.havingExprs,
		groupByCols: sb.groupByCols,
		orderByCols: sb.orderByCols,
		order:       sb.order,
		limit:       sb.limit,
		offset:      sb.offset,
		args:        sb.args.Copy(),
	}
}
