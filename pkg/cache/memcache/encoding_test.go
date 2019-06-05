package memcache

import (
	"bytes"
	"testing"

	mt "github.com/bilibili/kratos/pkg/cache/memcache/test"
)

func TestEncode(t *testing.T) {
	type TestObj struct {
		Name string
		Age  int32
	}
	testObj := TestObj{"abc", 1}

	ed := newEncodeDecoder()
	tests := []struct {
		name string
		a    *Item
		r    []byte
		e    error
	}{
		{
			"EncodeRawFlagErrItem",
			&Item{
				Object: &TestObj{"abc", 1},
				Flags:  FlagRAW,
			},
			[]byte{},
			ErrItem,
		},
		{
			"EncodeEncodeFlagErrItem",
			&Item{
				Value: []byte("test"),
				Flags: FlagJSON,
			},
			[]byte{},
			ErrItem,
		},
		{
			"EncodeEmpty",
			&Item{
				Value: []byte(""),
				Flags: FlagRAW,
			},
			[]byte(""),
			nil,
		},
		{
			"EncodeMaxSize",
			&Item{
				Value: bytes.Repeat([]byte("A"), 8000000),
				Flags: FlagRAW,
			},
			bytes.Repeat([]byte("A"), 8000000),
			nil,
		},
		{
			"EncodeExceededMaxSize",
			&Item{
				Value: bytes.Repeat([]byte("A"), 8000000+1),
				Flags: FlagRAW,
			},
			nil,
			ErrValueSize,
		},
		{
			"EncodeGOB",
			&Item{
				Object: testObj,
				Flags:  FlagGOB,
			},
			[]byte{38, 255, 131, 3, 1, 1, 7, 84, 101, 115, 116, 79, 98, 106, 1, 255, 132, 0, 1, 2, 1, 4, 78, 97, 109, 101, 1, 12, 0, 1, 3, 65, 103, 101, 1, 4, 0, 0, 0, 10, 255, 132, 1, 3, 97, 98, 99, 1, 2, 0},
			nil,
		},
		{
			"EncodeJSON",
			&Item{
				Object: testObj,
				Flags:  FlagJSON,
			},
			[]byte{123, 34, 78, 97, 109, 101, 34, 58, 34, 97, 98, 99, 34, 44, 34, 65, 103, 101, 34, 58, 49, 125, 10},
			nil,
		},
		{
			"EncodeProtobuf",
			&Item{
				Object: &mt.TestItem{Name: "abc", Age: 1},
				Flags:  FlagProtobuf,
			},
			[]byte{10, 3, 97, 98, 99, 16, 1},
			nil,
		},
		{
			"EncodeGzip",
			&Item{
				Value: bytes.Repeat([]byte("B"), 50),
				Flags: FlagGzip,
			},
			[]byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 114, 34, 25, 0, 2, 0, 0, 255, 255, 252, 253, 67, 209, 50, 0, 0, 0},
			nil,
		},
		{
			"EncodeGOBGzip",
			&Item{
				Object: testObj,
				Flags:  FlagGOB | FlagGzip,
			},
			[]byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 82, 251, 223, 204, 204, 200, 200, 30, 146, 90, 92, 226, 159, 148, 197, 248, 191, 133, 129, 145, 137, 145, 197, 47, 49, 55, 149, 145, 135, 129, 145, 217, 49, 61, 149, 145, 133, 129, 129, 129, 235, 127, 11, 35, 115, 98, 82, 50, 35, 19, 3, 32, 0, 0, 255, 255, 211, 249, 1, 154, 50, 0, 0, 0},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if r, err := ed.encode(test.a); err != test.e {
				t.Fatal(err)
			} else {
				if err == nil {
					if !bytes.Equal(r, test.r) {
						t.Fatalf("not equal, expect %v\n got %v", test.r, r)
					}
				}
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type TestObj struct {
		Name string
		Age  int32
	}
	testObj := &TestObj{"abc", 1}

	ed := newEncodeDecoder()
	tests := []struct {
		name string
		a    *Item
		r    interface{}
		e    error
	}{
		{
			"DecodeGOB",
			&Item{
				Flags: FlagGOB,
				Value: []byte{38, 255, 131, 3, 1, 1, 7, 84, 101, 115, 116, 79, 98, 106, 1, 255, 132, 0, 1, 2, 1, 4, 78, 97, 109, 101, 1, 12, 0, 1, 3, 65, 103, 101, 1, 4, 0, 0, 0, 10, 255, 132, 1, 3, 97, 98, 99, 1, 2, 0},
			},
			testObj,
			nil,
		},
		{
			"DecodeJSON",
			&Item{
				Value: []byte{123, 34, 78, 97, 109, 101, 34, 58, 34, 97, 98, 99, 34, 44, 34, 65, 103, 101, 34, 58, 49, 125, 10},
				Flags: FlagJSON,
			},
			testObj,
			nil,
		},
		{
			"DecodeProtobuf",
			&Item{
				Value: []byte{10, 3, 97, 98, 99, 16, 1},

				Flags: FlagProtobuf,
			},
			&mt.TestItem{Name: "abc", Age: 1},
			nil,
		},
		{
			"DecodeGzip",
			&Item{
				Value: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 114, 34, 25, 0, 2, 0, 0, 255, 255, 252, 253, 67, 209, 50, 0, 0, 0},
				Flags: FlagGzip,
			},
			bytes.Repeat([]byte("B"), 50),
			nil,
		},
		{
			"DecodeGOBGzip",
			&Item{
				Value: []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 82, 251, 223, 204, 204, 200, 200, 30, 146, 90, 92, 226, 159, 148, 197, 248, 191, 133, 129, 145, 137, 145, 197, 47, 49, 55, 149, 145, 135, 129, 145, 217, 49, 61, 149, 145, 133, 129, 129, 129, 235, 127, 11, 35, 115, 98, 82, 50, 35, 19, 3, 32, 0, 0, 255, 255, 211, 249, 1, 154, 50, 0, 0, 0},
				Flags: FlagGOB | FlagGzip,
			},
			testObj,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if (test.a.Flags & FlagProtobuf) > 0 {
				var dd mt.TestItem
				if err := ed.decode(test.a, &dd); err != nil {
					t.Fatal(err)
				}
				if (test.r.(*mt.TestItem).Name != dd.Name) || (test.r.(*mt.TestItem).Age != dd.Age) {
					t.Fatalf("compare failed error, expect %v\n got %v", test.r.(*mt.TestItem), dd)
				}
			} else if test.a.Flags == FlagGzip {
				var dd []byte
				if err := ed.decode(test.a, &dd); err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(dd, test.r.([]byte)) {
					t.Fatalf("compare failed error, expect %v\n got %v", test.r, dd)
				}
			} else {
				var dd TestObj
				if err := ed.decode(test.a, &dd); err != nil {
					t.Fatal(err)
				}
				if (test.r.(*TestObj).Name != dd.Name) || (test.r.(*TestObj).Age != dd.Age) {
					t.Fatalf("compare failed error, expect %v\n got %v", test.r.(*TestObj), dd)
				}
			}
		})
	}
}
