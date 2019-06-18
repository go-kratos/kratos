package memcache

import (
	"bytes"
	"encoding/json"
	"testing"

	test "github.com/bilibili/kratos/pkg/cache/memcache/test"
	"github.com/gogo/protobuf/proto"
)

func TestConnRaw(t *testing.T) {
	item := &Item{
		Key:        "test",
		Value:      []byte("test"),
		Flags:      FlagRAW,
		Expiration: 60,
		cas:        0,
	}
	if err := testConnASCII.Set(item); err != nil {
		t.Errorf("conn.Store() error(%v)", err)
	}
}

func TestConnSerialization(t *testing.T) {
	type TestObj struct {
		Name string
		Age  int32
	}

	tests := []struct {
		name string
		a    *Item
		e    error
	}{

		{
			"JSON",
			&Item{
				Key:        "test_json",
				Object:     &TestObj{"json", 1},
				Expiration: 60,
				Flags:      FlagJSON,
			},
			nil,
		},
		{
			"JSONGzip",
			&Item{
				Key:        "test_json_gzip",
				Object:     &TestObj{"jsongzip", 2},
				Expiration: 60,
				Flags:      FlagJSON | FlagGzip,
			},
			nil,
		},
		{
			"GOB",
			&Item{
				Key:        "test_gob",
				Object:     &TestObj{"gob", 3},
				Expiration: 60,
				Flags:      FlagGOB,
			},
			nil,
		},
		{
			"GOBGzip",
			&Item{
				Key:        "test_gob_gzip",
				Object:     &TestObj{"gobgzip", 4},
				Expiration: 60,
				Flags:      FlagGOB | FlagGzip,
			},
			nil,
		},
		{
			"Protobuf",
			&Item{
				Key:        "test_protobuf",
				Object:     &test.TestItem{Name: "protobuf", Age: 6},
				Expiration: 60,
				Flags:      FlagProtobuf,
			},
			nil,
		},
		{
			"ProtobufGzip",
			&Item{
				Key:        "test_protobuf_gzip",
				Object:     &test.TestItem{Name: "protobufgzip", Age: 7},
				Expiration: 60,
				Flags:      FlagProtobuf | FlagGzip,
			},
			nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := testConnASCII.Set(tc.a); err != nil {
				t.Fatal(err)
			}
			if r, err := testConnASCII.Get(tc.a.Key); err != tc.e {
				t.Fatal(err)
			} else {
				if (tc.a.Flags & FlagProtobuf) > 0 {
					var no test.TestItem
					if err := testConnASCII.Scan(r, &no); err != nil {
						t.Fatal(err)
					}
					if (tc.a.Object.(*test.TestItem).Name != no.Name) || (tc.a.Object.(*test.TestItem).Age != no.Age) {
						t.Fatalf("compare failed error, %v %v", tc.a.Object.(*test.TestItem), no)
					}
				} else {
					var no TestObj
					if err := testConnASCII.Scan(r, &no); err != nil {
						t.Fatal(err)
					}
					if (tc.a.Object.(*TestObj).Name != no.Name) || (tc.a.Object.(*TestObj).Age != no.Age) {
						t.Fatalf("compare failed error, %v %v", tc.a.Object.(*TestObj), no)
					}
				}

			}
		})
	}
}

func BenchmarkConnJSON(b *testing.B) {
	st := &struct {
		Name string
		Age  int
	}{"json", 10}
	itemx := &Item{Key: "json", Object: st, Flags: FlagJSON}
	var (
		eb  bytes.Buffer
		je  *json.Encoder
		ir  bytes.Reader
		jd  *json.Decoder
		jr  reader
		nst test.TestItem
	)
	jd = json.NewDecoder(&jr)
	je = json.NewEncoder(&eb)
	eb.Grow(_encodeBuf)
	// NOTE reuse bytes.Buffer internal buf
	// DON'T concurrency call Scan
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eb.Reset()
		if err := je.Encode(itemx.Object); err != nil {
			return
		}
		data := eb.Bytes()
		ir.Reset(data)
		jr.Reset(&ir)
		jd.Decode(&nst)
	}
}

func BenchmarkConnProtobuf(b *testing.B) {
	st := &test.TestItem{Name: "protobuf", Age: 10}
	itemx := &Item{Key: "protobuf", Object: st, Flags: FlagJSON}
	var (
		eb  bytes.Buffer
		nst test.TestItem
		ped *proto.Buffer
	)
	ped = proto.NewBuffer(eb.Bytes())
	eb.Grow(_encodeBuf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ped.Reset()
		pb, ok := itemx.Object.(proto.Message)
		if !ok {
			return
		}
		if err := ped.Marshal(pb); err != nil {
			return
		}
		data := ped.Bytes()
		ped.SetBuf(data)
		ped.Unmarshal(&nst)
	}
}
