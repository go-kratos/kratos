// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/mysql"
)

var ErrTableNotExist = errors.New("table is not exist")
var ErrMissingTableMeta = errors.New("missing table meta")
var HAHealthCheckSchema = "mysql.ha_health_check"

const (
	TYPE_NUMBER    = iota + 1 // tinyint, smallint, mediumint, int, bigint, year
	TYPE_FLOAT                // float, double
	TYPE_ENUM                 // enum
	TYPE_SET                  // set
	TYPE_STRING               // other
	TYPE_DATETIME             // datetime
	TYPE_TIMESTAMP            // timestamp
	TYPE_DATE                 // date
	TYPE_TIME                 // time
	TYPE_BIT                  // bit
	TYPE_JSON                 // json
)

type TableColumn struct {
	Name       string
	Type       int
	Collation  string
	RawType    string
	IsAuto     bool
	IsUnsigned bool
	EnumValues []string
	SetValues  []string
}

type Index struct {
	Name        string
	Columns     []string
	Cardinality []uint64
}

type Table struct {
	Schema string
	Name   string

	Columns   []TableColumn
	Indexes   []*Index
	PKColumns []int
}

func (ta *Table) String() string {
	return fmt.Sprintf("%s.%s", ta.Schema, ta.Name)
}

func (ta *Table) AddColumn(name string, columnType string, collation string, extra string) {
	index := len(ta.Columns)
	ta.Columns = append(ta.Columns, TableColumn{Name: name, Collation: collation})
	ta.Columns[index].RawType = columnType

	if strings.HasPrefix(columnType, "float") ||
		strings.HasPrefix(columnType, "double") ||
		strings.HasPrefix(columnType, "decimal") {
		ta.Columns[index].Type = TYPE_FLOAT
	} else if strings.HasPrefix(columnType, "enum") {
		ta.Columns[index].Type = TYPE_ENUM
		ta.Columns[index].EnumValues = strings.Split(strings.Replace(
			strings.TrimSuffix(
				strings.TrimPrefix(
					columnType, "enum("),
				")"),
			"'", "", -1),
			",")
	} else if strings.HasPrefix(columnType, "set") {
		ta.Columns[index].Type = TYPE_SET
		ta.Columns[index].SetValues = strings.Split(strings.Replace(
			strings.TrimSuffix(
				strings.TrimPrefix(
					columnType, "set("),
				")"),
			"'", "", -1),
			",")
	} else if strings.HasPrefix(columnType, "datetime") {
		ta.Columns[index].Type = TYPE_DATETIME
	} else if strings.HasPrefix(columnType, "timestamp") {
		ta.Columns[index].Type = TYPE_TIMESTAMP
	} else if strings.HasPrefix(columnType, "time") {
		ta.Columns[index].Type = TYPE_TIME
	} else if "date" == columnType {
		ta.Columns[index].Type = TYPE_DATE
	} else if strings.HasPrefix(columnType, "bit") {
		ta.Columns[index].Type = TYPE_BIT
	} else if strings.HasPrefix(columnType, "json") {
		ta.Columns[index].Type = TYPE_JSON
	} else if strings.Contains(columnType, "int") || strings.HasPrefix(columnType, "year") {
		ta.Columns[index].Type = TYPE_NUMBER
	} else {
		ta.Columns[index].Type = TYPE_STRING
	}

	if strings.Contains(columnType, "unsigned") || strings.Contains(columnType, "zerofill") {
		ta.Columns[index].IsUnsigned = true
	}

	if extra == "auto_increment" {
		ta.Columns[index].IsAuto = true
	}
}

func (ta *Table) FindColumn(name string) int {
	for i, col := range ta.Columns {
		if col.Name == name {
			return i
		}
	}
	return -1
}

func (ta *Table) GetPKColumn(index int) *TableColumn {
	return &ta.Columns[ta.PKColumns[index]]
}

func (ta *Table) AddIndex(name string) (index *Index) {
	index = NewIndex(name)
	ta.Indexes = append(ta.Indexes, index)
	return index
}

func NewIndex(name string) *Index {
	return &Index{name, make([]string, 0, 8), make([]uint64, 0, 8)}
}

func (idx *Index) AddColumn(name string, cardinality uint64) {
	idx.Columns = append(idx.Columns, name)
	if cardinality == 0 {
		cardinality = uint64(len(idx.Cardinality) + 1)
	}
	idx.Cardinality = append(idx.Cardinality, cardinality)
}

func (idx *Index) FindColumn(name string) int {
	for i, colName := range idx.Columns {
		if name == colName {
			return i
		}
	}
	return -1
}

func IsTableExist(conn mysql.Executer, schema string, name string) (bool, error) {
	query := fmt.Sprintf("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '%s' and TABLE_NAME = '%s' LIMIT 1", schema, name)
	r, err := conn.Execute(query)
	if err != nil {
		return false, errors.Trace(err)
	}

	return r.RowNumber() == 1, nil
}

func NewTableFromSqlDB(conn *sql.DB, schema string, name string) (*Table, error) {
	ta := &Table{
		Schema:  schema,
		Name:    name,
		Columns: make([]TableColumn, 0, 16),
		Indexes: make([]*Index, 0, 8),
	}

	if err := ta.fetchColumnsViaSqlDB(conn); err != nil {
		return nil, errors.Trace(err)
	}

	if err := ta.fetchIndexesViaSqlDB(conn); err != nil {
		return nil, errors.Trace(err)
	}

	return ta, nil
}

func NewTable(conn mysql.Executer, schema string, name string) (*Table, error) {
	ta := &Table{
		Schema:  schema,
		Name:    name,
		Columns: make([]TableColumn, 0, 16),
		Indexes: make([]*Index, 0, 8),
	}

	if err := ta.fetchColumns(conn); err != nil {
		return nil, errors.Trace(err)
	}

	if err := ta.fetchIndexes(conn); err != nil {
		return nil, errors.Trace(err)
	}

	return ta, nil
}

func (ta *Table) fetchColumns(conn mysql.Executer) error {
	r, err := conn.Execute(fmt.Sprintf("show full columns from `%s`.`%s`", ta.Schema, ta.Name))
	if err != nil {
		return errors.Trace(err)
	}

	for i := 0; i < r.RowNumber(); i++ {
		name, _ := r.GetString(i, 0)
		colType, _ := r.GetString(i, 1)
		collation, _ := r.GetString(i, 2)
		extra, _ := r.GetString(i, 6)

		ta.AddColumn(name, colType, collation, extra)
	}

	return nil
}

func (ta *Table) fetchColumnsViaSqlDB(conn *sql.DB) error {
	r, err := conn.Query(fmt.Sprintf("show full columns from `%s`.`%s`", ta.Schema, ta.Name))
	if err != nil {
		return errors.Trace(err)
	}

	defer r.Close()

	var unusedVal interface{}
	unused := &unusedVal

	for r.Next() {
		var name, colType, extra string
		var collation sql.NullString
		err := r.Scan(&name, &colType, &collation, &unused, &unused, &unused, &extra, &unused, &unused)
		if err != nil {
			return errors.Trace(err)
		}
		ta.AddColumn(name, colType, collation.String, extra)
	}

	return r.Err()
}

func (ta *Table) fetchIndexes(conn mysql.Executer) error {
	r, err := conn.Execute(fmt.Sprintf("show index from `%s`.`%s`", ta.Schema, ta.Name))
	if err != nil {
		return errors.Trace(err)
	}
	var currentIndex *Index
	currentName := ""

	for i := 0; i < r.RowNumber(); i++ {
		indexName, _ := r.GetString(i, 2)
		if currentName != indexName {
			currentIndex = ta.AddIndex(indexName)
			currentName = indexName
		}
		cardinality, _ := r.GetUint(i, 6)
		colName, _ := r.GetString(i, 4)
		currentIndex.AddColumn(colName, cardinality)
	}

	return ta.fetchPrimaryKeyColumns()

}

func (ta *Table) fetchIndexesViaSqlDB(conn *sql.DB) error {
	r, err := conn.Query(fmt.Sprintf("show index from `%s`.`%s`", ta.Schema, ta.Name))
	if err != nil {
		return errors.Trace(err)
	}

	defer r.Close()

	var currentIndex *Index
	currentName := ""

	var unusedVal interface{}
	unused := &unusedVal

	for r.Next() {
		var indexName, colName string
		var cardinality interface{}

		err := r.Scan(
			&unused,
			&unused,
			&indexName,
			&unused,
			&colName,
			&unused,
			&cardinality,
			&unused,
			&unused,
			&unused,
			&unused,
			&unused,
			&unused,
		)
		if err != nil {
			return errors.Trace(err)
		}

		if currentName != indexName {
			currentIndex = ta.AddIndex(indexName)
			currentName = indexName
		}

		c := toUint64(cardinality)
		currentIndex.AddColumn(colName, c)
	}

	return ta.fetchPrimaryKeyColumns()
}

func toUint64(i interface{}) uint64 {
	switch i := i.(type) {
	case int:
		return uint64(i)
	case int8:
		return uint64(i)
	case int16:
		return uint64(i)
	case int32:
		return uint64(i)
	case int64:
		return uint64(i)
	case uint:
		return uint64(i)
	case uint8:
		return uint64(i)
	case uint16:
		return uint64(i)
	case uint32:
		return uint64(i)
	case uint64:
		return uint64(i)
	}

	return 0
}

func (ta *Table) fetchPrimaryKeyColumns() error {
	if len(ta.Indexes) == 0 {
		return nil
	}

	pkIndex := ta.Indexes[0]
	if pkIndex.Name != "PRIMARY" {
		return nil
	}

	ta.PKColumns = make([]int, len(pkIndex.Columns))
	for i, pkCol := range pkIndex.Columns {
		ta.PKColumns[i] = ta.FindColumn(pkCol)
	}

	return nil
}
