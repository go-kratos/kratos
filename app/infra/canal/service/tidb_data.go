package service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"go-common/app/infra/canal/model"

	pb "github.com/pingcap/tidb-tools/tidb_binlog/slave_binlog_proto/go-binlog"
)

// lower case column field type in mysql
// https://dev.mysql.com/doc/refman/8.0/en/data-types.html
// for numeric type: int bigint smallint tinyint float double decimal bit
// for string type: text longtext mediumtext char tinytext varchar
// blob longblog mediumblog binary tinyblob varbinary
// enum set
// for json type: json

// for text and char type, string_value is set
// for blob and binary type, bytes_value is set
// for enum, set, uint64_value is set
// for json, bytes_value is set

func tidbMakeData(m *msg) (data *model.Data, err error) {
	action := m.mu.GetType()
	if (action != pb.MutationType_Insert) && (action != pb.MutationType_Delete) && (action != pb.MutationType_Update) {
		err = errInvalidAction
		return
	}
	data = &model.Data{
		Action: strings.ToLower(action.String()),
		Table:  m.table,
	}
	var keys []string
	switch action {
	case pb.MutationType_Insert, pb.MutationType_Delete:
		var values = m.mu.GetRow().GetColumns()
		for i, c := range m.columns {
			for _, key := range m.keys {
				if c.Name == key {
					keys = append(keys, columnToString(values[i]))
					break
				}
			}
			if m.ignore[c.Name] {
				continue
			}
			if data.New == nil {
				data.New = make(map[string]interface{}, len(m.columns))
			}
			if strings.Contains(c.GetMysqlType(), "binary") {
				data.New[c.Name] = base64.StdEncoding.EncodeToString(values[i].GetBytesValue())
				continue
			}
			data.New[c.Name] = columnToValue(values[i])
		}
	case pb.MutationType_Update:
		if m.mu.Row == nil || m.mu.ChangeRow == nil {
			err = errInvalidUpdate
			return
		}
		var oldValues = m.mu.GetChangeRow().GetColumns()
		var newValues = m.mu.GetRow().GetColumns()
		for i, c := range m.columns {
			for _, key := range m.keys {
				if c.Name == key {
					keys = append(keys, columnToString(newValues[i]))
					break
				}
			}
			if m.ignore[c.Name] {
				continue
			}
			if data.New == nil {
				data.New = make(map[string]interface{}, len(m.columns))
			}
			if data.Old == nil {
				data.Old = make(map[string]interface{}, len(m.columns))
			}
			if strings.Contains(c.GetMysqlType(), "binary") {
				data.Old[c.Name] = base64.StdEncoding.EncodeToString(oldValues[i].GetBytesValue())
				data.New[c.Name] = base64.StdEncoding.EncodeToString(newValues[i].GetBytesValue())
				continue
			}
			data.Old[c.Name] = columnToValue(oldValues[i])
			data.New[c.Name] = columnToValue(newValues[i])
		}
	}
	if len(keys) == 0 {
		data.Key = columnToString(m.mu.GetRow().GetColumns()[0])
	} else {
		data.Key = strings.Join(keys, ",")
	}
	if data.New == nil && data.Old == nil {
		data = nil
	}
	return
}

func columnToValue(c *pb.Column) interface{} {
	if c.GetIsNull() {
		return nil
	}
	if c.Int64Value != nil {
		return c.GetInt64Value()
	}
	if c.Uint64Value != nil {
		return c.GetUint64Value()
	}
	if c.DoubleValue != nil {
		return c.GetDoubleValue()
	}
	if c.StringValue != nil {
		return c.GetStringValue()
	}
	return c.GetBytesValue()
}

func columnToString(c *pb.Column) string {
	return fmt.Sprint(columnToValue(c))
}
