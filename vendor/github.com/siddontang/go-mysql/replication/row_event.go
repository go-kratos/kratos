package replication

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/shopspring/decimal"
	. "github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go/hack"
	log "github.com/sirupsen/logrus"
)

type errMissingTableMapEvent error

type TableMapEvent struct {
	tableIDSize int

	TableID uint64

	Flags uint16

	Schema []byte
	Table  []byte

	ColumnCount uint64
	ColumnType  []byte
	ColumnMeta  []uint16

	//len = (ColumnCount + 7) / 8
	NullBitmap []byte
}

func (e *TableMapEvent) Decode(data []byte) error {
	pos := 0
	e.TableID = FixedLengthInt(data[0:e.tableIDSize])
	pos += e.tableIDSize

	e.Flags = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	schemaLength := data[pos]
	pos++

	e.Schema = data[pos : pos+int(schemaLength)]
	pos += int(schemaLength)

	//skip 0x00
	pos++

	tableLength := data[pos]
	pos++

	e.Table = data[pos : pos+int(tableLength)]
	pos += int(tableLength)

	//skip 0x00
	pos++

	var n int
	e.ColumnCount, _, n = LengthEncodedInt(data[pos:])
	pos += n

	e.ColumnType = data[pos : pos+int(e.ColumnCount)]
	pos += int(e.ColumnCount)

	var err error
	var metaData []byte
	if metaData, _, n, err = LengthEnodedString(data[pos:]); err != nil {
		return errors.Trace(err)
	}

	if err = e.decodeMeta(metaData); err != nil {
		return errors.Trace(err)
	}

	pos += n

	if len(data[pos:]) != bitmapByteSize(int(e.ColumnCount)) {
		return io.EOF
	}

	e.NullBitmap = data[pos:]

	return nil
}

func bitmapByteSize(columnCount int) int {
	return int(columnCount+7) / 8
}

// see mysql sql/log_event.h
/*
	0 byte
	MYSQL_TYPE_DECIMAL
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR

	1 byte
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_BLOB
	MYSQL_TYPE_GEOMETRY

	//maybe
	MYSQL_TYPE_TIME2
	MYSQL_TYPE_DATETIME2
	MYSQL_TYPE_TIMESTAMP2

	2 byte
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT
	MYSQL_TYPE_NEWDECIMAL
	MYSQL_TYPE_VAR_STRING
	MYSQL_TYPE_STRING

	This enumeration value is only used internally and cannot exist in a binlog.
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_ENUM
	MYSQL_TYPE_SET
	MYSQL_TYPE_TINY_BLOB
	MYSQL_TYPE_MEDIUM_BLOB
	MYSQL_TYPE_LONG_BLOB
*/
func (e *TableMapEvent) decodeMeta(data []byte) error {
	pos := 0
	e.ColumnMeta = make([]uint16, e.ColumnCount)
	for i, t := range e.ColumnType {
		switch t {
		case MYSQL_TYPE_STRING:
			var x uint16 = uint16(data[pos]) << 8 //real type
			x += uint16(data[pos+1])              //pack or field length
			e.ColumnMeta[i] = x
			pos += 2
		case MYSQL_TYPE_NEWDECIMAL:
			var x uint16 = uint16(data[pos]) << 8 //precision
			x += uint16(data[pos+1])              //decimals
			e.ColumnMeta[i] = x
			pos += 2
		case MYSQL_TYPE_VAR_STRING,
			MYSQL_TYPE_VARCHAR,
			MYSQL_TYPE_BIT:
			e.ColumnMeta[i] = binary.LittleEndian.Uint16(data[pos:])
			pos += 2
		case MYSQL_TYPE_BLOB,
			MYSQL_TYPE_DOUBLE,
			MYSQL_TYPE_FLOAT,
			MYSQL_TYPE_GEOMETRY,
			MYSQL_TYPE_JSON:
			e.ColumnMeta[i] = uint16(data[pos])
			pos++
		case MYSQL_TYPE_TIME2,
			MYSQL_TYPE_DATETIME2,
			MYSQL_TYPE_TIMESTAMP2:
			e.ColumnMeta[i] = uint16(data[pos])
			pos++
		case MYSQL_TYPE_NEWDATE,
			MYSQL_TYPE_ENUM,
			MYSQL_TYPE_SET,
			MYSQL_TYPE_TINY_BLOB,
			MYSQL_TYPE_MEDIUM_BLOB,
			MYSQL_TYPE_LONG_BLOB:
			return errors.Errorf("unsupport type in binlog %d", t)
		default:
			e.ColumnMeta[i] = 0
		}
	}

	return nil
}

func (e *TableMapEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "TableID: %d\n", e.TableID)
	fmt.Fprintf(w, "TableID size: %d\n", e.tableIDSize)
	fmt.Fprintf(w, "Flags: %d\n", e.Flags)
	fmt.Fprintf(w, "Schema: %s\n", e.Schema)
	fmt.Fprintf(w, "Table: %s\n", e.Table)
	fmt.Fprintf(w, "Column count: %d\n", e.ColumnCount)
	fmt.Fprintf(w, "Column type: \n%s", hex.Dump(e.ColumnType))
	fmt.Fprintf(w, "NULL bitmap: \n%s", hex.Dump(e.NullBitmap))
	fmt.Fprintln(w)
}

// RowsEventStmtEndFlag is set in the end of the statement.
const RowsEventStmtEndFlag = 0x01

type RowsEvent struct {
	//0, 1, 2
	Version int

	tableIDSize int
	tables      map[uint64]*TableMapEvent
	needBitmap2 bool

	Table *TableMapEvent

	TableID uint64

	Flags uint16

	//if version == 2
	ExtraData []byte

	//lenenc_int
	ColumnCount uint64
	//len = (ColumnCount + 7) / 8
	ColumnBitmap1 []byte

	//if UPDATE_ROWS_EVENTv1 or v2
	//len = (ColumnCount + 7) / 8
	ColumnBitmap2 []byte

	//rows: invalid: int64, float64, bool, []byte, string
	Rows [][]interface{}

	parseTime  bool
	useDecimal bool
}

func (e *RowsEvent) Decode(data []byte) error {
	pos := 0
	e.TableID = FixedLengthInt(data[0:e.tableIDSize])
	pos += e.tableIDSize

	e.Flags = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if e.Version == 2 {
		dataLen := binary.LittleEndian.Uint16(data[pos:])
		pos += 2

		e.ExtraData = data[pos : pos+int(dataLen-2)]
		pos += int(dataLen - 2)
	}

	var n int
	e.ColumnCount, _, n = LengthEncodedInt(data[pos:])
	pos += n

	bitCount := bitmapByteSize(int(e.ColumnCount))
	e.ColumnBitmap1 = data[pos : pos+bitCount]
	pos += bitCount

	if e.needBitmap2 {
		e.ColumnBitmap2 = data[pos : pos+bitCount]
		pos += bitCount
	}

	var ok bool
	e.Table, ok = e.tables[e.TableID]
	if !ok {
		if len(e.tables) > 0 {
			return errors.Errorf("invalid table id %d, no corresponding table map event", e.TableID)
		} else {
			return errMissingTableMapEvent(errors.Errorf("invalid table id %d, no corresponding table map event", e.TableID))
		}
	}

	var err error

	// ... repeat rows until event-end
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("parse rows event panic %v, data %q, parsed rows %#v, table map %#v\n%s", r, data, e, e.Table, Pstack())
		}
	}()

	for pos < len(data) {
		if n, err = e.decodeRows(data[pos:], e.Table, e.ColumnBitmap1); err != nil {
			return errors.Trace(err)
		}
		pos += n

		if e.needBitmap2 {
			if n, err = e.decodeRows(data[pos:], e.Table, e.ColumnBitmap2); err != nil {
				return errors.Trace(err)
			}
			pos += n
		}
	}

	return nil
}

func isBitSet(bitmap []byte, i int) bool {
	return bitmap[i>>3]&(1<<(uint(i)&7)) > 0
}

func (e *RowsEvent) decodeRows(data []byte, table *TableMapEvent, bitmap []byte) (int, error) {
	row := make([]interface{}, e.ColumnCount)

	pos := 0

	// refer: https://github.com/alibaba/canal/blob/c3e38e50e269adafdd38a48c63a1740cde304c67/dbsync/src/main/java/com/taobao/tddl/dbsync/binlog/event/RowsLogBuffer.java#L63
	count := 0
	for i := 0; i < int(e.ColumnCount); i++ {
		if isBitSet(bitmap, i) {
			count++
		}
	}
	count = (count + 7) / 8

	nullBitmap := data[pos : pos+count]
	pos += count

	nullbitIndex := 0

	var n int
	var err error
	for i := 0; i < int(e.ColumnCount); i++ {
		if !isBitSet(bitmap, i) {
			continue
		}

		isNull := (uint32(nullBitmap[nullbitIndex/8]) >> uint32(nullbitIndex%8)) & 0x01
		nullbitIndex++

		if isNull > 0 {
			row[i] = nil
			continue
		}

		row[i], n, err = e.decodeValue(data[pos:], table.ColumnType[i], table.ColumnMeta[i])

		if err != nil {
			return 0, err
		}
		pos += n
	}

	e.Rows = append(e.Rows, row)
	return pos, nil
}

func (e *RowsEvent) parseFracTime(t interface{}) interface{} {
	v, ok := t.(fracTime)
	if !ok {
		return t
	}

	if !e.parseTime {
		// Don't parse time, return string directly
		return v.String()
	}

	// return Golang time directly
	return v.Time
}

// see mysql sql/log_event.cc log_event_print_value
func (e *RowsEvent) decodeValue(data []byte, tp byte, meta uint16) (v interface{}, n int, err error) {
	var length int = 0

	if tp == MYSQL_TYPE_STRING {
		if meta >= 256 {
			b0 := uint8(meta >> 8)
			b1 := uint8(meta & 0xFF)

			if b0&0x30 != 0x30 {
				length = int(uint16(b1) | (uint16((b0&0x30)^0x30) << 4))
				tp = byte(b0 | 0x30)
			} else {
				length = int(meta & 0xFF)
				tp = b0
			}
		} else {
			length = int(meta)
		}
	}

	switch tp {
	case MYSQL_TYPE_NULL:
		return nil, 0, nil
	case MYSQL_TYPE_LONG:
		n = 4
		v = ParseBinaryInt32(data)
	case MYSQL_TYPE_TINY:
		n = 1
		v = ParseBinaryInt8(data)
	case MYSQL_TYPE_SHORT:
		n = 2
		v = ParseBinaryInt16(data)
	case MYSQL_TYPE_INT24:
		n = 3
		v = ParseBinaryInt24(data)
	case MYSQL_TYPE_LONGLONG:
		n = 8
		v = ParseBinaryInt64(data)
	case MYSQL_TYPE_NEWDECIMAL:
		prec := uint8(meta >> 8)
		scale := uint8(meta & 0xFF)
		v, n, err = decodeDecimal(data, int(prec), int(scale), e.useDecimal)
	case MYSQL_TYPE_FLOAT:
		n = 4
		v = ParseBinaryFloat32(data)
	case MYSQL_TYPE_DOUBLE:
		n = 8
		v = ParseBinaryFloat64(data)
	case MYSQL_TYPE_BIT:
		nbits := ((meta >> 8) * 8) + (meta & 0xFF)
		n = int(nbits+7) / 8

		//use int64 for bit
		v, err = decodeBit(data, int(nbits), int(n))
	case MYSQL_TYPE_TIMESTAMP:
		n = 4
		t := binary.LittleEndian.Uint32(data)
		v = e.parseFracTime(fracTime{time.Unix(int64(t), 0), 0})
	case MYSQL_TYPE_TIMESTAMP2:
		v, n, err = decodeTimestamp2(data, meta)
		v = e.parseFracTime(v)
	case MYSQL_TYPE_DATETIME:
		n = 8
		i64 := binary.LittleEndian.Uint64(data)
		d := i64 / 1000000
		t := i64 % 1000000
		v = e.parseFracTime(fracTime{time.Date(int(d/10000),
			time.Month((d%10000)/100),
			int(d%100),
			int(t/10000),
			int((t%10000)/100),
			int(t%100),
			0,
			time.UTC), 0})
	case MYSQL_TYPE_DATETIME2:
		v, n, err = decodeDatetime2(data, meta)
		v = e.parseFracTime(v)
	case MYSQL_TYPE_TIME:
		n = 3
		i32 := uint32(FixedLengthInt(data[0:3]))
		if i32 == 0 {
			v = "00:00:00"
		} else {
			sign := ""
			if i32 < 0 {
				sign = "-"
			}
			v = fmt.Sprintf("%s%02d:%02d:%02d", sign, i32/10000, (i32%10000)/100, i32%100)
		}
	case MYSQL_TYPE_TIME2:
		v, n, err = decodeTime2(data, meta)
	case MYSQL_TYPE_DATE:
		n = 3
		i32 := uint32(FixedLengthInt(data[0:3]))
		if i32 == 0 {
			v = "0000-00-00"
		} else {
			v = fmt.Sprintf("%04d-%02d-%02d", i32/(16*32), i32/32%16, i32%32)
		}

	case MYSQL_TYPE_YEAR:
		n = 1
		v = int(data[0]) + 1900
	case MYSQL_TYPE_ENUM:
		l := meta & 0xFF
		switch l {
		case 1:
			v = int64(data[0])
			n = 1
		case 2:
			v = int64(binary.BigEndian.Uint16(data))
			n = 2
		default:
			err = fmt.Errorf("Unknown ENUM packlen=%d", l)
		}
	case MYSQL_TYPE_SET:
		n = int(meta & 0xFF)
		nbits := n * 8

		v, err = decodeBit(data, nbits, n)
	case MYSQL_TYPE_BLOB:
		v, n, err = decodeBlob(data, meta)
	case MYSQL_TYPE_VARCHAR,
		MYSQL_TYPE_VAR_STRING:
		length = int(meta)
		v, n = decodeString(data, length)
	case MYSQL_TYPE_STRING:
		v, n = decodeString(data, length)
	case MYSQL_TYPE_JSON:
		// Refer: https://github.com/shyiko/mysql-binlog-connector-java/blob/master/src/main/java/com/github/shyiko/mysql/binlog/event/deserialization/AbstractRowsEventDataDeserializer.java#L404
		length = int(FixedLengthInt(data[0:meta]))
		n = length + int(meta)
		v, err = e.decodeJsonBinary(data[meta:n])
	case MYSQL_TYPE_GEOMETRY:
		// MySQL saves Geometry as Blob in binlog
		// Seem that the binary format is SRID (4 bytes) + WKB, outer can use
		// MySQL GeoFromWKB or others to create the geometry data.
		// Refer https://dev.mysql.com/doc/refman/5.7/en/gis-wkb-functions.html
		// I also find some go libs to handle WKB if possible
		// see https://github.com/twpayne/go-geom or https://github.com/paulmach/go.geo
		v, n, err = decodeBlob(data, meta)
	default:
		err = fmt.Errorf("unsupport type %d in binlog and don't know how to handle", tp)
	}
	return
}

func decodeString(data []byte, length int) (v string, n int) {
	if length < 256 {
		length = int(data[0])

		n = int(length) + 1
		v = hack.String(data[1:n])
	} else {
		length = int(binary.LittleEndian.Uint16(data[0:]))
		n = length + 2
		v = hack.String(data[2:n])
	}

	return
}

const digitsPerInteger int = 9

var compressedBytes = []int{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}

func decodeDecimalDecompressValue(compIndx int, data []byte, mask uint8) (size int, value uint32) {
	size = compressedBytes[compIndx]
	databuff := make([]byte, size)
	for i := 0; i < size; i++ {
		databuff[i] = data[i] ^ mask
	}
	value = uint32(BFixedLengthInt(databuff))
	return
}

func decodeDecimal(data []byte, precision int, decimals int, useDecimal bool) (interface{}, int, error) {
	//see python mysql replication and https://github.com/jeremycole/mysql_binlog
	integral := (precision - decimals)
	uncompIntegral := int(integral / digitsPerInteger)
	uncompFractional := int(decimals / digitsPerInteger)
	compIntegral := integral - (uncompIntegral * digitsPerInteger)
	compFractional := decimals - (uncompFractional * digitsPerInteger)

	binSize := uncompIntegral*4 + compressedBytes[compIntegral] +
		uncompFractional*4 + compressedBytes[compFractional]

	buf := make([]byte, binSize)
	copy(buf, data[:binSize])

	//must copy the data for later change
	data = buf

	// Support negative
	// The sign is encoded in the high bit of the the byte
	// But this bit can also be used in the value
	value := uint32(data[0])
	var res bytes.Buffer
	var mask uint32 = 0
	if value&0x80 == 0 {
		mask = uint32((1 << 32) - 1)
		res.WriteString("-")
	}

	//clear sign
	data[0] ^= 0x80

	pos, value := decodeDecimalDecompressValue(compIntegral, data, uint8(mask))
	res.WriteString(fmt.Sprintf("%d", value))

	for i := 0; i < uncompIntegral; i++ {
		value = binary.BigEndian.Uint32(data[pos:]) ^ mask
		pos += 4
		res.WriteString(fmt.Sprintf("%09d", value))
	}

	res.WriteString(".")

	for i := 0; i < uncompFractional; i++ {
		value = binary.BigEndian.Uint32(data[pos:]) ^ mask
		pos += 4
		res.WriteString(fmt.Sprintf("%09d", value))
	}

	if size, value := decodeDecimalDecompressValue(compFractional, data[pos:], uint8(mask)); size > 0 {
		res.WriteString(fmt.Sprintf("%0*d", compFractional, value))
		pos += size
	}

	if useDecimal {
		f, err := decimal.NewFromString(hack.String(res.Bytes()))
		return f, pos, err
	}

	f, err := strconv.ParseFloat(hack.String(res.Bytes()), 64)
	return f, pos, err
}

func decodeBit(data []byte, nbits int, length int) (value int64, err error) {
	if nbits > 1 {
		switch length {
		case 1:
			value = int64(data[0])
		case 2:
			value = int64(binary.BigEndian.Uint16(data))
		case 3:
			value = int64(BFixedLengthInt(data[0:3]))
		case 4:
			value = int64(binary.BigEndian.Uint32(data))
		case 5:
			value = int64(BFixedLengthInt(data[0:5]))
		case 6:
			value = int64(BFixedLengthInt(data[0:6]))
		case 7:
			value = int64(BFixedLengthInt(data[0:7]))
		case 8:
			value = int64(binary.BigEndian.Uint64(data))
		default:
			err = fmt.Errorf("invalid bit length %d", length)
		}
	} else {
		if length != 1 {
			err = fmt.Errorf("invalid bit length %d", length)
		} else {
			value = int64(data[0])
		}
	}
	return
}

func decodeTimestamp2(data []byte, dec uint16) (interface{}, int, error) {
	//get timestamp binary length
	n := int(4 + (dec+1)/2)
	sec := int64(binary.BigEndian.Uint32(data[0:4]))
	usec := int64(0)
	switch dec {
	case 1, 2:
		usec = int64(data[4]) * 10000
	case 3, 4:
		usec = int64(binary.BigEndian.Uint16(data[4:])) * 100
	case 5, 6:
		usec = int64(BFixedLengthInt(data[4:7]))
	}

	if sec == 0 {
		return formatZeroTime(int(usec), int(dec)), n, nil
	}

	return fracTime{time.Unix(sec, usec*1000), int(dec)}, n, nil
}

const DATETIMEF_INT_OFS int64 = 0x8000000000

func decodeDatetime2(data []byte, dec uint16) (interface{}, int, error) {
	//get datetime binary length
	n := int(5 + (dec+1)/2)

	intPart := int64(BFixedLengthInt(data[0:5])) - DATETIMEF_INT_OFS
	var frac int64 = 0

	switch dec {
	case 1, 2:
		frac = int64(data[5]) * 10000
	case 3, 4:
		frac = int64(binary.BigEndian.Uint16(data[5:7])) * 100
	case 5, 6:
		frac = int64(BFixedLengthInt(data[5:8]))
	}

	if intPart == 0 {
		return formatZeroTime(int(frac), int(dec)), n, nil
	}

	tmp := intPart<<24 + frac
	//handle sign???
	if tmp < 0 {
		tmp = -tmp
	}

	// var secPart int64 = tmp % (1 << 24)
	ymdhms := tmp >> 24

	ymd := ymdhms >> 17
	ym := ymd >> 5
	hms := ymdhms % (1 << 17)

	day := int(ymd % (1 << 5))
	month := int(ym % 13)
	year := int(ym / 13)

	second := int(hms % (1 << 6))
	minute := int((hms >> 6) % (1 << 6))
	hour := int((hms >> 12))

	return fracTime{time.Date(year, time.Month(month), day, hour, minute, second, int(frac*1000), time.UTC), int(dec)}, n, nil
}

const TIMEF_OFS int64 = 0x800000000000
const TIMEF_INT_OFS int64 = 0x800000

func decodeTime2(data []byte, dec uint16) (string, int, error) {
	//time  binary length
	n := int(3 + (dec+1)/2)

	tmp := int64(0)
	intPart := int64(0)
	frac := int64(0)
	switch dec {
	case 1:
	case 2:
		intPart = int64(BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		frac = int64(data[3])
		if intPart < 0 && frac > 0 {
			/*
			   Negative values are stored with reverse fractional part order,
			   for binary sort compatibility.

			     Disk value  intpart frac   Time value   Memory value
			     800000.00    0      0      00:00:00.00  0000000000.000000
			     7FFFFF.FF   -1      255   -00:00:00.01  FFFFFFFFFF.FFD8F0
			     7FFFFF.9D   -1      99    -00:00:00.99  FFFFFFFFFF.F0E4D0
			     7FFFFF.00   -1      0     -00:00:01.00  FFFFFFFFFF.000000
			     7FFFFE.FF   -1      255   -00:00:01.01  FFFFFFFFFE.FFD8F0
			     7FFFFE.F6   -2      246   -00:00:01.10  FFFFFFFFFE.FE7960

			     Formula to convert fractional part from disk format
			     (now stored in "frac" variable) to absolute value: "0x100 - frac".
			     To reconstruct in-memory value, we shift
			     to the next integer value and then substruct fractional part.
			*/
			intPart++     /* Shift to the next integer value */
			frac -= 0x100 /* -(0x100 - frac) */
		}
		tmp = intPart<<24 + frac*10000
	case 3:
	case 4:
		intPart = int64(BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		frac = int64(binary.BigEndian.Uint16(data[3:5]))
		if intPart < 0 && frac > 0 {
			/*
			   Fix reverse fractional part order: "0x10000 - frac".
			   See comments for FSP=1 and FSP=2 above.
			*/
			intPart++       /* Shift to the next integer value */
			frac -= 0x10000 /* -(0x10000-frac) */
		}
		tmp = intPart<<24 + frac*100

	case 5:
	case 6:
		tmp = int64(BFixedLengthInt(data[0:6])) - TIMEF_OFS
	default:
		intPart = int64(BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		tmp = intPart << 24
	}

	if intPart == 0 {
		return "00:00:00", n, nil
	}

	hms := int64(0)
	sign := ""
	if tmp < 0 {
		tmp = -tmp
		sign = "-"
	}

	hms = tmp >> 24

	hour := (hms >> 12) % (1 << 10) /* 10 bits starting at 12th */
	minute := (hms >> 6) % (1 << 6) /* 6 bits starting at 6th   */
	second := hms % (1 << 6)        /* 6 bits starting at 0th   */
	secPart := tmp % (1 << 24)

	if secPart != 0 {
		return fmt.Sprintf("%s%02d:%02d:%02d.%06d", sign, hour, minute, second, secPart), n, nil
	}

	return fmt.Sprintf("%s%02d:%02d:%02d", sign, hour, minute, second), n, nil
}

func decodeBlob(data []byte, meta uint16) (v []byte, n int, err error) {
	var length int
	switch meta {
	case 1:
		length = int(data[0])
		v = data[1 : 1+length]
		n = length + 1
	case 2:
		length = int(binary.LittleEndian.Uint16(data))
		v = data[2 : 2+length]
		n = length + 2
	case 3:
		length = int(FixedLengthInt(data[0:3]))
		v = data[3 : 3+length]
		n = length + 3
	case 4:
		length = int(binary.LittleEndian.Uint32(data))
		v = data[4 : 4+length]
		n = length + 4
	default:
		err = fmt.Errorf("invalid blob packlen = %d", meta)
	}

	return
}

func (e *RowsEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "TableID: %d\n", e.TableID)
	fmt.Fprintf(w, "Flags: %d\n", e.Flags)
	fmt.Fprintf(w, "Column count: %d\n", e.ColumnCount)

	fmt.Fprintf(w, "Values:\n")
	for _, rows := range e.Rows {
		fmt.Fprintf(w, "--\n")
		for j, d := range rows {
			if _, ok := d.([]byte); ok {
				fmt.Fprintf(w, "%d:%q\n", j, d)
			} else {
				fmt.Fprintf(w, "%d:%#v\n", j, d)
			}
		}
	}
	fmt.Fprintln(w)
}

type RowsQueryEvent struct {
	Query []byte
}

func (e *RowsQueryEvent) Decode(data []byte) error {
	//ignore length byte 1
	e.Query = data[1:]
	return nil
}

func (e *RowsQueryEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Query: %s\n", e.Query)
	fmt.Fprintln(w)
}
