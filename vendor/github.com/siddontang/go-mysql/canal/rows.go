package canal

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/siddontang/go-mysql/schema"
)

const (
	UpdateAction = "update"
	InsertAction = "insert"
	DeleteAction = "delete"
)

type RowsEvent struct {
	Table  *schema.Table
	Action string
	// changed row list
	// binlog has three update event version, v0, v1 and v2.
	// for v1 and v2, the rows number must be even.
	// Two rows for one event, format is [before update row, after update row]
	// for update v0, only one row for a event, and we don't support this version.
	Rows [][]interface{}
}

func newRowsEvent(table *schema.Table, action string, rows [][]interface{}) *RowsEvent {
	e := new(RowsEvent)

	e.Table = table
	e.Action = action
	e.Rows = rows

	return e
}

// Get primary keys in one row for a table, a table may use multi fields as the PK
func GetPKValues(table *schema.Table, row []interface{}) ([]interface{}, error) {
	indexes := table.PKColumns
	if len(indexes) == 0 {
		return nil, errors.Errorf("table %s has no PK", table)
	} else if len(table.Columns) != len(row) {
		return nil, errors.Errorf("table %s has %d columns, but row data %v len is %d", table,
			len(table.Columns), row, len(row))
	}

	values := make([]interface{}, 0, len(indexes))

	for _, index := range indexes {
		values = append(values, row[index])
	}

	return values, nil
}

// Get term column's value
func GetColumnValue(table *schema.Table, column string, row []interface{}) (interface{}, error) {
	index := table.FindColumn(column)
	if index == -1 {
		return nil, errors.Errorf("table %s has no column name %s", table, column)
	}

	return row[index], nil
}

// String implements fmt.Stringer interface.
func (r *RowsEvent) String() string {
	return fmt.Sprintf("%s %s %v", r.Action, r.Table, r.Rows)
}
