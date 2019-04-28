package mysql

import (
	"math"
	"strconv"

	"github.com/juju/errors"
	"github.com/siddontang/go/hack"
)

func formatTextValue(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case int8:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case uint8:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case []byte:
		return v, nil
	case string:
		return hack.Slice(v), nil
	default:
		return nil, errors.Errorf("invalid type %T", value)
	}
}

func formatBinaryValue(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case int8:
		return Uint64ToBytes(uint64(v)), nil
	case int16:
		return Uint64ToBytes(uint64(v)), nil
	case int32:
		return Uint64ToBytes(uint64(v)), nil
	case int64:
		return Uint64ToBytes(uint64(v)), nil
	case int:
		return Uint64ToBytes(uint64(v)), nil
	case uint8:
		return Uint64ToBytes(uint64(v)), nil
	case uint16:
		return Uint64ToBytes(uint64(v)), nil
	case uint32:
		return Uint64ToBytes(uint64(v)), nil
	case uint64:
		return Uint64ToBytes(uint64(v)), nil
	case uint:
		return Uint64ToBytes(uint64(v)), nil
	case float32:
		return Uint64ToBytes(math.Float64bits(float64(v))), nil
	case float64:
		return Uint64ToBytes(math.Float64bits(v)), nil
	case []byte:
		return v, nil
	case string:
		return hack.Slice(v), nil
	default:
		return nil, errors.Errorf("invalid type %T", value)
	}
}
func formatField(field *Field, value interface{}) error {
	switch value.(type) {
	case int8, int16, int32, int64, int:
		field.Charset = 63
		field.Type = MYSQL_TYPE_LONGLONG
		field.Flag = BINARY_FLAG | NOT_NULL_FLAG
	case uint8, uint16, uint32, uint64, uint:
		field.Charset = 63
		field.Type = MYSQL_TYPE_LONGLONG
		field.Flag = BINARY_FLAG | NOT_NULL_FLAG | UNSIGNED_FLAG
	case float32, float64:
		field.Charset = 63
		field.Type = MYSQL_TYPE_DOUBLE
		field.Flag = BINARY_FLAG | NOT_NULL_FLAG
	case string, []byte:
		field.Charset = 33
		field.Type = MYSQL_TYPE_VAR_STRING
	default:
		return errors.Errorf("unsupport type %T for resultset", value)
	}
	return nil
}

func BuildSimpleTextResultset(names []string, values [][]interface{}) (*Resultset, error) {
	r := new(Resultset)

	r.Fields = make([]*Field, len(names))

	var b []byte
	var err error

	for i, vs := range values {
		if len(vs) != len(r.Fields) {
			return nil, errors.Errorf("row %d has %d column not equal %d", i, len(vs), len(r.Fields))
		}

		var row []byte
		for j, value := range vs {
			if i == 0 {
				field := &Field{}
				r.Fields[j] = field
				field.Name = hack.Slice(names[j])

				if err = formatField(field, value); err != nil {
					return nil, errors.Trace(err)
				}
			}
			b, err = formatTextValue(value)

			if err != nil {
				return nil, errors.Trace(err)
			}

			row = append(row, PutLengthEncodedString(b)...)
		}

		r.RowDatas = append(r.RowDatas, row)
	}

	return r, nil
}

func BuildSimpleBinaryResultset(names []string, values [][]interface{}) (*Resultset, error) {
	r := new(Resultset)

	r.Fields = make([]*Field, len(names))

	var b []byte
	var err error

	bitmapLen := ((len(names) + 7 + 2) >> 3)

	for i, vs := range values {
		if len(vs) != len(r.Fields) {
			return nil, errors.Errorf("row %d has %d column not equal %d", i, len(vs), len(r.Fields))
		}

		var row []byte
		nullBitmap := make([]byte, bitmapLen)

		row = append(row, 0)
		row = append(row, nullBitmap...)

		for j, value := range vs {
			if i == 0 {
				field := &Field{}
				r.Fields[j] = field
				field.Name = hack.Slice(names[j])

				if err = formatField(field, value); err != nil {
					return nil, errors.Trace(err)
				}
			}
			if value == nil {
				nullBitmap[(i+2)/8] |= (1 << (uint(i+2) % 8))
				continue
			}

			b, err = formatBinaryValue(value)

			if err != nil {
				return nil, errors.Trace(err)
			}

			if r.Fields[j].Type == MYSQL_TYPE_VAR_STRING {
				row = append(row, PutLengthEncodedString(b)...)
			} else {
				row = append(row, b...)
			}
		}

		copy(row[1:], nullBitmap)

		r.RowDatas = append(r.RowDatas, row)
	}

	return r, nil
}

func BuildSimpleResultset(names []string, values [][]interface{}, binary bool) (*Resultset, error) {
	if binary {
		return BuildSimpleBinaryResultset(names, values)
	} else {
		return BuildSimpleTextResultset(names, values)
	}
}
