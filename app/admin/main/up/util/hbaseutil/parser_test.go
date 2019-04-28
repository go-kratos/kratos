package hbaseutil

import (
	"encoding/binary"
	"encoding/json"
	"github.com/tsuna/gohbase/hrpc"
	"reflect"
	"testing"
	"time"
)

func uint64ToByte(value uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, value)
	return buf
}

func uint32ToByte(value uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, value)
	return buf
}

func uint16ToByte(value uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, value)
	return buf
}

type testStruct struct {
	A         int64   `family:"f" qualifier:"q64"`
	B         *int32  `family:"f" qualifier:"q32"`
	C         int16   `family:"f" qualifier:"q16"`
	D         int     `qualifier:"q"`
	S         *string `qualifier:"s"`
	FailField bool    `family:"f" qualifier:"fail"`

	MapInt    map[string]int    `family:"m1"`
	MapString map[string]string `family:"m2"`
}

func (t *testStruct) equal(o testStruct) bool {
	return t.A == o.A &&
		*t.B == *o.B &&
		t.C == o.C &&
		t.D == o.D &&
		*t.S == *o.S &&
		t.FailField == o.FailField &&
		reflect.DeepEqual(t.MapInt, o.MapInt) &&
		reflect.DeepEqual(t.MapString, o.MapString)

}

var testcase = [][]*hrpc.Cell{
	{
		{Family: []byte("f"), Qualifier: []byte("q64"), Value: uint64ToByte(10000000000)},
		{Family: []byte("f"), Qualifier: []byte("q32"), Value: uint32ToByte(1000000)},
		{Family: []byte("f"), Qualifier: []byte("q16"), Value: uint16ToByte(100)},
		{Family: []byte("f"), Qualifier: []byte("q"), Value: uint64ToByte(1000000)},
		{Family: []byte("f"), Qualifier: []byte("s"), Value: []byte("just test")},
		{Family: []byte("f"), Qualifier: []byte("fail"), Value: []byte("1")},
		{Family: []byte("m1"), Qualifier: []byte("k1"), Value: uint32ToByte(1)},
		{Family: []byte("m1"), Qualifier: []byte("k2"), Value: uint16ToByte(2)},
		{Family: []byte("m2"), Qualifier: []byte("k1"), Value: []byte("1")},
		{Family: []byte("m2"), Qualifier: []byte("k2"), Value: []byte("2")},
	},
	{
		{Family: []byte("f"), Qualifier: []byte("q64"), Value: uint64ToByte(10000000000)},
		{Family: []byte("f"), Qualifier: []byte("q32"), Value: uint64ToByte(10000000000)},
		{Family: []byte("f"), Qualifier: []byte("q16"), Value: uint64ToByte(10000000000)},
		{Family: []byte("f"), Qualifier: []byte("q"), Value: uint64ToByte(10000000000)},
		{Family: []byte("f"), Qualifier: []byte("s"), Value: []byte("just test2")},
		{Family: []byte("f"), Qualifier: []byte("fail"), Value: []byte("1")},
	},
}
var resultb = []int32{
	1000000,
	10000000000 & 0xffffffff,
}
var results = []string{
	"just test",
	"just test2",
}
var resultcase = []testStruct{
	{A: 10000000000, B: &resultb[0], C: 100, D: 1000000, S: &results[0],
		MapInt:    map[string]int{"k1": 1, "k2": 2},
		MapString: map[string]string{"k1": "1", "k2": "2"}},
	{A: 10000000000, B: &resultb[1], C: -7168, D: int(10000000000), S: &results[1]},
}

func TestParser_Parse(t *testing.T) {
	var parser = Parser{}

	for i, cells := range testcase {
		var result testStruct
		var err = parser.Parse(cells, &result)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%v, expect=%v", result, resultcase[i])
		if !resultcase[i].equal(result) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}
}

var testcase2 = [][]*hrpc.Cell{
	{
		{Family: []byte("f"), Qualifier: []byte("q64"), Value: []byte("10000000000")},
		{Family: []byte("f"), Qualifier: []byte("q32"), Value: []byte("1000000")},
		{Family: []byte("f"), Qualifier: []byte("q16"), Value: []byte("100")},
		{Family: []byte("f"), Qualifier: []byte("q"), Value: []byte("1000000")},
		{Family: []byte("f"), Qualifier: []byte("s"), Value: []byte("just test")},
		{Family: []byte("f"), Qualifier: []byte("fail"), Value: []byte("1")},
		{Family: []byte("m1"), Qualifier: []byte("k1"), Value: []byte("1")},
		{Family: []byte("m1"), Qualifier: []byte("k2"), Value: []byte("2")},
		{Family: []byte("m2"), Qualifier: []byte("k1"), Value: []byte("1")},
		{Family: []byte("m2"), Qualifier: []byte("k2"), Value: []byte("2")},
	},
}
var resultcase2 = []testStruct{
	{A: 10000000000, B: &resultb[0], C: 100, D: 1000000, S: &results[0],
		MapInt:    map[string]int{"k1": 1, "k2": 2},
		MapString: map[string]string{"k1": "1", "k2": "2"}},
}

func TestParser_ParseCustomParseInt(t *testing.T) {
	var parser = Parser{
		ParseIntFunc: StringToUint,
	}

	for i, cells := range testcase2 {
		var result testStruct
		var err = parser.Parse(cells, &result)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%v, expect=%v", result, resultcase[i])
		if !resultcase2[i].equal(result) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}

}

func TestParser_StringInt(t *testing.T) {
	var testcase2 = [][]*hrpc.Cell{
		{
			{Family: []byte("f"), Qualifier: []byte("q64"), Value: []byte("9223372036854775807")},
			{Family: []byte("f"), Qualifier: []byte("q32"), Value: []byte("1000000")},
			{Family: []byte("f"), Qualifier: []byte("q16"), Value: []byte("10000")},
			{Family: []byte("f"), Qualifier: []byte("q"), Value: []byte("100")},
		},
		{
			{Family: []byte("f"), Qualifier: []byte("q64"), Value: []byte("-9223372036854775808")},
			{Family: []byte("f"), Qualifier: []byte("q32"), Value: []byte("-2")},
			{Family: []byte("f"), Qualifier: []byte("q16"), Value: []byte("-3")},
			{Family: []byte("f"), Qualifier: []byte("q"), Value: []byte("-4")},
			{Family: []byte("f"), Qualifier: []byte("u64"), Value: []byte("18446744073709551615")},
			{Family: []byte("f"), Qualifier: []byte("u32"), Value: []byte("2147483648")},
			{Family: []byte("f"), Qualifier: []byte("u16"), Value: []byte("32768")},
			{Family: []byte("f"), Qualifier: []byte("u8"), Value: []byte("128")},
		},
	}
	type testStruct struct {
		A   int64  `family:"f" qualifier:"q64"`
		B   int32  `family:"f" qualifier:"q32"`
		C   int16  `family:"f" qualifier:"q16"`
		D   int8   `qualifier:"q"`
		U64 uint64 `family:"f" qualifier:"u64"`
		U32 uint32 `family:"f" qualifier:"u32"`
		U16 uint16 `family:"f" qualifier:"u16"`
		U8  uint8  `family:"f" qualifier:"u8"`
	}
	var expect = []testStruct{
		{
			A: 9223372036854775807, B: 1000000, C: 10000, D: 100,
		},
		{
			A: -9223372036854775808, B: -2, C: -3, D: -4,
			U64: 18446744073709551615, U32: 2147483648, U16: 32768, U8: 128,
		},
	}
	var parser = Parser{
		ParseIntFunc: StringToUint,
	}

	for i, cells := range testcase2 {
		var result testStruct
		var err = parser.Parse(cells, &result)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%+v, expect=%+v", result, expect[i])
		if !reflect.DeepEqual(expect[i], result) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}

}

func TestParser_InterfaceInterface(t *testing.T) {
	var parser = Parser{}

	for i, cells := range testcase {
		var st testStruct
		var result interface{} = &st
		var err = parser.Parse(cells, result)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%v, expect=%v", result, resultcase[i])
		if !resultcase[i].equal(*result.(*testStruct)) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}
}

type testOnlyQualifier struct {
	A         int64             `qualifier:"q64"`
	B         *int32            `qualifier:"q32"`
	C         int16             `qualifier:"q16"`
	D         int               `qualifier:"q"`
	S         *string           `qualifier:"s"`
	FailField bool              `qualifier:"fail"`
	MapInt    map[string]int    `family:"m1"`
	MapString map[string]string `family:"m2"`
}

var testcase3 = [][]*hrpc.Cell{
	{
		{Family: []byte("f"), Qualifier: []byte("q64"), Value: uint64ToByte(10000000000)},
		{Family: []byte("f"), Qualifier: []byte("q32"), Value: uint32ToByte(1000000)},
		{Family: []byte("f"), Qualifier: []byte("q16"), Value: uint16ToByte(100)},
		{Family: []byte("f"), Qualifier: []byte("q"), Value: uint64ToByte(1000000)},
		{Family: []byte("f"), Qualifier: []byte("s"), Value: []byte("just test")},
		{Family: []byte("f"), Qualifier: []byte("fail"), Value: []byte("1")},
		{Family: []byte("m1"), Qualifier: []byte("k1"), Value: uint32ToByte(1)},
		{Family: []byte("m1"), Qualifier: []byte("k2"), Value: uint16ToByte(2)},
		{Family: []byte("m2"), Qualifier: []byte("k1"), Value: []byte("1")},
		{Family: []byte("m2"), Qualifier: []byte("k2"), Value: []byte("2")},
	},
}

var resultcase3 = []testOnlyQualifier{
	{A: 10000000000, B: &resultb[0], C: 100, D: 1000000, S: &results[0],
		MapInt:    map[string]int{"k1": 1, "k2": 2},
		MapString: map[string]string{"k1": "1", "k2": "2"}},
}

func (t *testOnlyQualifier) equal(o testOnlyQualifier) bool {
	return t.A == o.A &&
		*t.B == *o.B &&
		t.C == o.C &&
		t.D == o.D &&
		*t.S == *o.S &&
		t.FailField == o.FailField &&
		reflect.DeepEqual(t.MapInt, o.MapInt) &&
		reflect.DeepEqual(t.MapString, o.MapString)

}
func TestParser_OnlyQualifier(t *testing.T) {
	var parser = Parser{}

	for i, cells := range testcase3 {
		var result testOnlyQualifier
		var err = parser.Parse(cells, &result)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%v, expect=%v", result, resultcase3[i])
		if !resultcase3[i].equal(result) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}

}

func TestParser_PartialFamilyPartPartialQualifier(t *testing.T) {
	var testcase = [][]*hrpc.Cell{
		{
			{Family: []byte("f"), Qualifier: []byte("q64"), Value: uint64ToByte(10000000000)},
			{Family: []byte("f"), Qualifier: []byte("q32"), Value: uint32ToByte(1000000)},
			{Family: []byte("f"), Qualifier: []byte("q16"), Value: uint16ToByte(100)},
			{Family: []byte("f"), Qualifier: []byte("q"), Value: uint64ToByte(1000000)},
			{Family: []byte("f"), Qualifier: []byte("s"), Value: []byte("just test")},
			{Family: []byte("f"), Qualifier: []byte("fail"), Value: []byte("1")},
		},
	}

	type testStruct struct {
		A   int            `family:"f" qualifier:"q64"`
		Map map[string]int `family:"f"`
	}
	var resultcase = []testStruct{
		{A: 10000000000,
			Map: map[string]int{
				"q32":  1000000,
				"q16":  100,
				"q":    1000000,
				"fail": int('1'),
			}},
	}

	var parser = Parser{}

	for i, cells := range testcase {

		var result testStruct
		var resInterface interface{} = &result
		var start = time.Now()
		var err = parser.Parse(cells, &resInterface)
		var elapse = time.Since(start)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%v, expect=%v, parse time=%v", resInterface, resultcase[i], elapse)
		if !reflect.DeepEqual(resultcase[i], result) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}

}

func TestParser_SetBasicValue(t *testing.T) {
	var (
		i32s int32 = -2
		i64s int64 = -3
		i16s int16 = -4

		i32d int32
		i64d int64
		i16d int16

		i64Bytes = uint64ToByte(uint64(i64s))
		i32Bytes = uint32ToByte(uint32(i32s))
		i16Bytes = uint16ToByte(uint16(i16s))
	)

	var rvi32 = reflect.ValueOf(&i32d)
	var e = setBasicValue(i32Bytes, rvi32, "rvi32", ByteBigEndianToUint64)
	if e != nil {
		t.Errorf("fail, err=%+v", e)
		t.FailNow()
	}
	if i32s != i32d {
		t.Errorf("fail,expect=%d, got=%d", i32d, i32s)
		t.FailNow()
	}

	var rvi64 = reflect.ValueOf(&i64d)
	e = setBasicValue(i64Bytes, rvi64, "rvi64", ByteBigEndianToUint64)
	if e != nil {
		t.Errorf("fail, err=%+v", e)
		t.FailNow()
	}
	if i64s != i64d {
		t.Errorf("fail,expect=%d, got=%d", i64d, i64s)
		t.FailNow()
	}
	var rvi16 = reflect.ValueOf(&i16d)
	e = setBasicValue(i16Bytes, rvi16, "rvi16", ByteBigEndianToUint64)
	if e != nil {
		t.Errorf("fail, err=%+v", e)
		t.FailNow()
	}
	if i16s != i16d {
		t.Errorf("fail,expect=%d, got=%d", i16d, i16s)
		t.FailNow()
	}

}

func TestParser_Minus(t *testing.T) {
	var i64 int64 = -2
	var i32 int32 = -3
	var i16 int16 = -4
	var i8 int8 = -128
	var testcase = [][]*hrpc.Cell{
		{
			{Family: []byte("f"), Qualifier: []byte("i64"), Value: uint64ToByte(uint64(i64))},
			{Family: []byte("f"), Qualifier: []byte("i32"), Value: uint32ToByte(uint32(i32))},
			{Family: []byte("f"), Qualifier: []byte("i16"), Value: uint16ToByte(uint16(i16))},
			{Family: []byte("f"), Qualifier: []byte("i8"), Value: []byte{uint8(i8)}},
		},
	}

	type testStruct struct {
		I64 int64 `family:"f" qualifier:"i64"`
		I32 int32 `family:"f" qualifier:"i32"`
		I16 int16 `family:"f" qualifier:"i16"`
		I8  int8  `family:"f" qualifier:"i8"`
	}

	var resultcase = []testStruct{
		{
			I64: i64,
			I32: i32,
			I16: i16,
			I8:  i8,
		},
	}

	var parser = Parser{}

	for i, cells := range testcase {

		var result testStruct
		var resInterface interface{} = &result
		var start = time.Now()
		var err = parser.Parse(cells, &resInterface)
		var elapse = time.Since(start)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%v, expect=%v, parse time=%v", resInterface, resultcase[i], elapse)
		if !reflect.DeepEqual(resultcase[i], result) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}

}

func TestParser_Overflow(t *testing.T) {
	var testcase = [][]*hrpc.Cell{
		{
			{Family: []byte("f"), Qualifier: []byte("i64"), Value: uint64ToByte(0xff01020304050607)},
			{Family: []byte("f"), Qualifier: []byte("i32"), Value: uint64ToByte(0xff01020304050607)},
			{Family: []byte("f"), Qualifier: []byte("i16"), Value: uint64ToByte(0xff01020304050607)},
			{Family: []byte("f"), Qualifier: []byte("i8"), Value: uint64ToByte(0xff01020304050607)},
		},
	}

	type testStruct struct {
		I64 int64 `family:"f" qualifier:"i64"`
		I32 int32 `family:"f" qualifier:"i32"`
		I16 int16 `family:"f" qualifier:"i16"`
		I8  int8  `family:"f" qualifier:"i8"`
	}

	var resultcase = []testStruct{
		{
			I64: -(0xffffffffffffffff - 0xff01020304050607 + 1),
			I32: 0xff01020304050607 & 0xffffffff,
			I16: 0xff01020304050607 & 0xffff,
			I8:  0xff01020304050607 & 0xff,
		},
	}

	var parser = Parser{}

	for i, cells := range testcase {

		var result testStruct
		var resInterface interface{} = &result
		var start = time.Now()
		var err = parser.Parse(cells, &resInterface)
		var elapse = time.Since(start)
		if err != nil {
			t.Logf("err=%v", err)
			t.Fail()
		}
		t.Logf("result=%v, expect=%v, parse time=%v", resInterface, resultcase[i], elapse)
		if !reflect.DeepEqual(resultcase[i], result) {
			t.Logf("fail case: index=%d", i)
			t.Fail()
		}

	}
}

func BenchmarkParser(b *testing.B) {
	var testcase = [][]*hrpc.Cell{
		{
			{Family: []byte("f"), Qualifier: []byte("q64"), Value: uint64ToByte(10000000000)},
			{Family: []byte("f"), Qualifier: []byte("q32"), Value: uint32ToByte(1000000)},
			{Family: []byte("f"), Qualifier: []byte("q16"), Value: uint16ToByte(100)},
			{Family: []byte("f"), Qualifier: []byte("q"), Value: uint64ToByte(1000000)},
			{Family: []byte("f"), Qualifier: []byte("fail"), Value: uint64ToByte(100)},
		},
	}

	type testStruct struct {
		A   int            `family:"f" qualifier:"q64"`
		B   int32          `family:"f" qualifier:"q32"`
		C   int16          `family:"f" qualifier:"q16"`
		D   int            `qualifier:"q"`
		S   string         `qualifier:"s"`
		Map map[string]int `family:"f"`
	}
	var parser Parser
	b.Logf("bench parser")
	for i := 0; i < b.N; i++ {
		var result testStruct
		parser.Parse(testcase[0], &result)
	}
}

func BenchmarkJson(b *testing.B) {
	var text = []byte("{\"a\": 123, \"b\": \"1234\", \"c\": \"1234\", \"d\": \"1234\", \"e\" :{\"a\": 123, \"b\": \"1234\"}}")
	var result struct {
		A int    `json:"a"`
		B string `json:"b"`
		C string `json:"c"`
		D string `json:"d"`
		E struct {
			A int    `json:"a"`
			B string `json:"b"`
		} `json:"e"`
	}
	b.Logf("bench json")
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(text, &result); err != nil {
			b.FailNow()
		}
	}
}
