package replication

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/juju/errors"
	. "github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go/hack"
)

const (
	JSONB_SMALL_OBJECT byte = iota // small JSON object
	JSONB_LARGE_OBJECT             // large JSON object
	JSONB_SMALL_ARRAY              // small JSON array
	JSONB_LARGE_ARRAY              // large JSON array
	JSONB_LITERAL                  // literal (true/false/null)
	JSONB_INT16                    // int16
	JSONB_UINT16                   // uint16
	JSONB_INT32                    // int32
	JSONB_UINT32                   // uint32
	JSONB_INT64                    // int64
	JSONB_UINT64                   // uint64
	JSONB_DOUBLE                   // double
	JSONB_STRING                   // string
	JSONB_OPAQUE       byte = 0x0f // custom data (any MySQL data type)
)

const (
	JSONB_NULL_LITERAL  byte = 0x00
	JSONB_TRUE_LITERAL  byte = 0x01
	JSONB_FALSE_LITERAL byte = 0x02
)

const (
	jsonbSmallOffsetSize = 2
	jsonbLargeOffsetSize = 4

	jsonbKeyEntrySizeSmall = 2 + jsonbSmallOffsetSize
	jsonbKeyEntrySizeLarge = 2 + jsonbLargeOffsetSize

	jsonbValueEntrySizeSmall = 1 + jsonbSmallOffsetSize
	jsonbValueEntrySizeLarge = 1 + jsonbLargeOffsetSize
)

func jsonbGetOffsetSize(isSmall bool) int {
	if isSmall {
		return jsonbSmallOffsetSize
	}

	return jsonbLargeOffsetSize
}

func jsonbGetKeyEntrySize(isSmall bool) int {
	if isSmall {
		return jsonbKeyEntrySizeSmall
	}

	return jsonbKeyEntrySizeLarge
}

func jsonbGetValueEntrySize(isSmall bool) int {
	if isSmall {
		return jsonbValueEntrySizeSmall
	}

	return jsonbValueEntrySizeLarge
}

// decodeJsonBinary decodes the JSON binary encoding data and returns
// the common JSON encoding data.
func (e *RowsEvent) decodeJsonBinary(data []byte) ([]byte, error) {
	d := jsonBinaryDecoder{useDecimal: e.useDecimal}

	if d.isDataShort(data, 1) {
		return nil, d.err
	}

	v := d.decodeValue(data[0], data[1:])
	if d.err != nil {
		return nil, d.err
	}

	return json.Marshal(v)
}

type jsonBinaryDecoder struct {
	useDecimal bool
	err        error
}

func (d *jsonBinaryDecoder) decodeValue(tp byte, data []byte) interface{} {
	if d.err != nil {
		return nil
	}

	switch tp {
	case JSONB_SMALL_OBJECT:
		return d.decodeObjectOrArray(data, true, true)
	case JSONB_LARGE_OBJECT:
		return d.decodeObjectOrArray(data, false, true)
	case JSONB_SMALL_ARRAY:
		return d.decodeObjectOrArray(data, true, false)
	case JSONB_LARGE_ARRAY:
		return d.decodeObjectOrArray(data, false, false)
	case JSONB_LITERAL:
		return d.decodeLiteral(data)
	case JSONB_INT16:
		return d.decodeInt16(data)
	case JSONB_UINT16:
		return d.decodeUint16(data)
	case JSONB_INT32:
		return d.decodeInt32(data)
	case JSONB_UINT32:
		return d.decodeUint32(data)
	case JSONB_INT64:
		return d.decodeInt64(data)
	case JSONB_UINT64:
		return d.decodeUint64(data)
	case JSONB_DOUBLE:
		return d.decodeDouble(data)
	case JSONB_STRING:
		return d.decodeString(data)
	case JSONB_OPAQUE:
		return d.decodeOpaque(data)
	default:
		d.err = errors.Errorf("invalid json type %d", tp)
	}

	return nil
}

func (d *jsonBinaryDecoder) decodeObjectOrArray(data []byte, isSmall bool, isObject bool) interface{} {
	offsetSize := jsonbGetOffsetSize(isSmall)
	if d.isDataShort(data, 2*offsetSize) {
		return nil
	}

	count := d.decodeCount(data, isSmall)
	size := d.decodeCount(data[offsetSize:], isSmall)

	if d.isDataShort(data, int(size)) {
		return nil
	}

	keyEntrySize := jsonbGetKeyEntrySize(isSmall)
	valueEntrySize := jsonbGetValueEntrySize(isSmall)

	headerSize := 2*offsetSize + count*valueEntrySize

	if isObject {
		headerSize += count * keyEntrySize
	}

	if headerSize > size {
		d.err = errors.Errorf("header size %d > size %d", headerSize, size)
		return nil
	}

	var keys []string
	if isObject {
		keys = make([]string, count)
		for i := 0; i < count; i++ {
			// decode key
			entryOffset := 2*offsetSize + keyEntrySize*i
			keyOffset := d.decodeCount(data[entryOffset:], isSmall)
			keyLength := int(d.decodeUint16(data[entryOffset+offsetSize:]))

			// Key must start after value entry
			if keyOffset < headerSize {
				d.err = errors.Errorf("invalid key offset %d, must > %d", keyOffset, headerSize)
				return nil
			}

			if d.isDataShort(data, keyOffset+keyLength) {
				return nil
			}

			keys[i] = hack.String(data[keyOffset : keyOffset+keyLength])
		}
	}

	if d.err != nil {
		return nil
	}

	values := make([]interface{}, count)
	for i := 0; i < count; i++ {
		// decode value
		entryOffset := 2*offsetSize + valueEntrySize*i
		if isObject {
			entryOffset += keyEntrySize * count
		}

		tp := data[entryOffset]

		if isInlineValue(tp, isSmall) {
			values[i] = d.decodeValue(tp, data[entryOffset+1:entryOffset+valueEntrySize])
			continue
		}

		valueOffset := d.decodeCount(data[entryOffset+1:], isSmall)

		if d.isDataShort(data, valueOffset) {
			return nil
		}

		values[i] = d.decodeValue(tp, data[valueOffset:])
	}

	if d.err != nil {
		return nil
	}

	if !isObject {
		return values
	}

	m := make(map[string]interface{}, count)
	for i := 0; i < count; i++ {
		m[keys[i]] = values[i]
	}

	return m
}

func isInlineValue(tp byte, isSmall bool) bool {
	switch tp {
	case JSONB_INT16, JSONB_UINT16, JSONB_LITERAL:
		return true
	case JSONB_INT32, JSONB_UINT32:
		return !isSmall
	}

	return false
}

func (d *jsonBinaryDecoder) decodeLiteral(data []byte) interface{} {
	if d.isDataShort(data, 1) {
		return nil
	}

	tp := data[0]

	switch tp {
	case JSONB_NULL_LITERAL:
		return nil
	case JSONB_TRUE_LITERAL:
		return true
	case JSONB_FALSE_LITERAL:
		return false
	}

	d.err = errors.Errorf("invalid literal %c", tp)

	return nil
}

func (d *jsonBinaryDecoder) isDataShort(data []byte, expected int) bool {
	if d.err != nil {
		return true
	}

	if len(data) < expected {
		d.err = errors.Errorf("data len %d < expected %d", len(data), expected)
	}

	return d.err != nil
}

func (d *jsonBinaryDecoder) decodeInt16(data []byte) int16 {
	if d.isDataShort(data, 2) {
		return 0
	}

	v := ParseBinaryInt16(data[0:2])
	return v
}

func (d *jsonBinaryDecoder) decodeUint16(data []byte) uint16 {
	if d.isDataShort(data, 2) {
		return 0
	}

	v := ParseBinaryUint16(data[0:2])
	return v
}

func (d *jsonBinaryDecoder) decodeInt32(data []byte) int32 {
	if d.isDataShort(data, 4) {
		return 0
	}

	v := ParseBinaryInt32(data[0:4])
	return v
}

func (d *jsonBinaryDecoder) decodeUint32(data []byte) uint32 {
	if d.isDataShort(data, 4) {
		return 0
	}

	v := ParseBinaryUint32(data[0:4])
	return v
}

func (d *jsonBinaryDecoder) decodeInt64(data []byte) int64 {
	if d.isDataShort(data, 8) {
		return 0
	}

	v := ParseBinaryInt64(data[0:8])
	return v
}

func (d *jsonBinaryDecoder) decodeUint64(data []byte) uint64 {
	if d.isDataShort(data, 8) {
		return 0
	}

	v := ParseBinaryUint64(data[0:8])
	return v
}

func (d *jsonBinaryDecoder) decodeDouble(data []byte) float64 {
	if d.isDataShort(data, 8) {
		return 0
	}

	v := ParseBinaryFloat64(data[0:8])
	return v
}

func (d *jsonBinaryDecoder) decodeString(data []byte) string {
	if d.err != nil {
		return ""
	}

	l, n := d.decodeVariableLength(data)

	if d.isDataShort(data, l+n) {
		return ""
	}

	data = data[n:]

	v := hack.String(data[0:l])
	return v
}

func (d *jsonBinaryDecoder) decodeOpaque(data []byte) interface{} {
	if d.isDataShort(data, 1) {
		return nil
	}

	tp := data[0]
	data = data[1:]

	l, n := d.decodeVariableLength(data)

	if d.isDataShort(data, l+n) {
		return nil
	}

	data = data[n : l+n]

	switch tp {
	case MYSQL_TYPE_NEWDECIMAL:
		return d.decodeDecimal(data)
	case MYSQL_TYPE_TIME:
		return d.decodeTime(data)
	case MYSQL_TYPE_DATE, MYSQL_TYPE_DATETIME, MYSQL_TYPE_TIMESTAMP:
		return d.decodeDateTime(data)
	default:
		return hack.String(data)
	}

	return nil
}

func (d *jsonBinaryDecoder) decodeDecimal(data []byte) interface{} {
	precision := int(data[0])
	scale := int(data[1])

	v, _, err := decodeDecimal(data[2:], precision, scale, d.useDecimal)
	d.err = err

	return v
}

func (d *jsonBinaryDecoder) decodeTime(data []byte) interface{} {
	v := d.decodeInt64(data)

	if v == 0 {
		return "00:00:00"
	}

	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}

	intPart := v >> 24
	hour := (intPart >> 12) % (1 << 10)
	min := (intPart >> 6) % (1 << 6)
	sec := intPart % (1 << 6)
	frac := v % (1 << 24)

	return fmt.Sprintf("%s%02d:%02d:%02d.%06d", sign, hour, min, sec, frac)
}

func (d *jsonBinaryDecoder) decodeDateTime(data []byte) interface{} {
	v := d.decodeInt64(data)
	if v == 0 {
		return "0000-00-00 00:00:00"
	}

	// handle negative?
	if v < 0 {
		v = -v
	}

	intPart := v >> 24
	ymd := intPart >> 17
	ym := ymd >> 5
	hms := intPart % (1 << 17)

	year := ym / 13
	month := ym % 13
	day := ymd % (1 << 5)
	hour := (hms >> 12)
	minute := (hms >> 6) % (1 << 6)
	second := hms % (1 << 6)
	frac := v % (1 << 24)

	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d", year, month, day, hour, minute, second, frac)

}

func (d *jsonBinaryDecoder) decodeCount(data []byte, isSmall bool) int {
	if isSmall {
		v := d.decodeUint16(data)
		return int(v)
	}

	return int(d.decodeUint32(data))
}

func (d *jsonBinaryDecoder) decodeVariableLength(data []byte) (int, int) {
	// The max size for variable length is math.MaxUint32, so
	// here we can use 5 bytes to save it.
	maxCount := 5
	if len(data) < maxCount {
		maxCount = len(data)
	}

	pos := 0
	length := uint64(0)
	for ; pos < maxCount; pos++ {
		v := data[pos]
		length |= uint64(v&0x7F) << uint(7*pos)

		if v&0x80 == 0 {
			if length > math.MaxUint32 {
				d.err = errors.Errorf("variable length %d must <= %d", length, math.MaxUint32)
				return 0, 0
			}

			pos += 1
			// TODO: should consider length overflow int here.
			return int(length), pos
		}
	}

	d.err = errors.New("decode variable length failed")

	return 0, 0
}
