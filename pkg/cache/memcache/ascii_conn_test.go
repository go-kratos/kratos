package memcache

import (
	"bytes"
	"strconv"
	"strings"

	"testing"
)

func TestASCIIConnAdd(t *testing.T) {
	tests := []struct {
		name string
		a    *Item
		e    error
	}{
		{
			"Add",
			&Item{
				Key:        "test_add",
				Value:      []byte("0"),
				Flags:      0,
				Expiration: 60,
			},
			nil,
		},
		{
			"Add_Large",
			&Item{
				Key:        "test_add_large",
				Value:      bytes.Repeat(space, _largeValue+1),
				Flags:      0,
				Expiration: 60,
			},
			nil,
		},
		{
			"Add_Exist",
			&Item{
				Key:        "test_add",
				Value:      []byte("0"),
				Flags:      0,
				Expiration: 60,
			},
			ErrNotStored,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := testConnASCII.Add(test.a); err != test.e {
				t.Fatal(err)
			}
			if b, err := testConnASCII.Get(test.a.Key); err != nil {
				t.Fatal(err)
			} else {
				compareItem(t, test.a, b)
			}
		})
	}
}

func TestASCIIConnGet(t *testing.T) {
	tests := []struct {
		name string
		a    *Item
		k    string
		e    error
	}{
		{
			"Get",
			&Item{
				Key:        "test_get",
				Value:      []byte("0"),
				Flags:      0,
				Expiration: 60,
			},
			"test_get",
			nil,
		},
		{
			"Get_NotExist",
			&Item{
				Key:        "test_get_not_exist",
				Value:      []byte("0"),
				Flags:      0,
				Expiration: 60,
			},
			"test_get_not_exist!",
			ErrNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := testConnASCII.Add(test.a); err != nil {
				t.Fatal(err)
			}
			if b, err := testConnASCII.Get(test.a.Key); err != nil {
				t.Fatal(err)
			} else {
				compareItem(t, test.a, b)
			}
		})
	}
}

//func TestGetHasErr(t *testing.T) {
//	prepareEnv(t)
//
//	st := &TestItem{Name: "json", Age: 10}
//	itemx := &Item{Key: "test", Object: st, Flags: FlagJSON}
//	c.Set(itemx)
//
//	expected := errors.New("some error")
//	monkey.Patch(scanGetReply, func(line []byte, item *Item) (size int, err error) {
//		return 0, expected
//	})
//
//	if _, err := c.Get("test"); err.Error() != expected.Error() {
//		t.Errorf("conn.Get() unexpected error(%v)", err)
//	}
//	if err := c.(*asciiConn).err; err.Error() != expected.Error() {
//		t.Errorf("unexpected error(%v)", err)
//	}
//}

func TestASCIIConnGetMulti(t *testing.T) {
	tests := []struct {
		name string
		a    []*Item
		k    []string
		e    error
	}{
		{"getMulti_Add",
			[]*Item{
				{
					Key:        "get_multi_1",
					Value:      []byte("test"),
					Flags:      FlagRAW,
					Expiration: 60,
					cas:        0,
				},
				{
					Key:        "get_multi_2",
					Value:      []byte("test2"),
					Flags:      FlagRAW,
					Expiration: 60,
					cas:        0,
				},
			},
			[]string{"get_multi_1", "get_multi_2"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, i := range test.a {
				if err := testConnASCII.Set(i); err != nil {
					t.Fatal(err)
				}
			}
			if r, err := testConnASCII.GetMulti(test.k); err != nil {
				t.Fatal(err)
			} else {
				reply := r["get_multi_1"]
				compareItem(t, reply, test.a[0])
				reply = r["get_multi_2"]
				compareItem(t, reply, test.a[1])
			}

		})
	}

}

func TestASCIIConnSet(t *testing.T) {
	tests := []struct {
		name string
		a    *Item
		e    error
	}{
		{
			"SetLowerBound",
			&Item{
				Key:        strings.Repeat("a", 1),
				Value:      []byte("4"),
				Flags:      0,
				Expiration: 60,
			},
			nil,
		},
		{
			"SetUpperBound",
			&Item{
				Key:        strings.Repeat("a", 250),
				Value:      []byte("3"),
				Flags:      0,
				Expiration: 60,
			},
			nil,
		},
		{
			"SetIllegalKeyZeroLength",
			&Item{
				Key:        "",
				Value:      []byte("2"),
				Flags:      0,
				Expiration: 60,
			},
			ErrMalformedKey,
		},
		{
			"SetIllegalKeyLengthExceededLimit",
			&Item{
				Key:        " ",
				Value:      []byte("1"),
				Flags:      0,
				Expiration: 60,
			},
			ErrMalformedKey,
		},
		{
			"SeJsonItem",
			&Item{
				Key: "set_obj",
				Object: &struct {
					Name string
					Age  int
				}{"json", 10},
				Expiration: 60,
				Flags:      FlagJSON,
			},
			nil,
		},
		{
			"SeErrItemJSONGzip",
			&Item{
				Key:        "set_err_item",
				Expiration: 60,
				Flags:      FlagJSON | FlagGzip,
			},
			ErrItem,
		},
		{
			"SeErrItemBytesValueWrongFlag",
			&Item{
				Key:        "set_err_item",
				Value:      []byte("2"),
				Expiration: 60,
				Flags:      FlagJSON,
			},
			ErrItem,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := testConnASCII.Set(test.a); err != test.e {
				t.Fatal(err)
			}
		})
	}
}

func TestASCIIConnCompareAndSwap(t *testing.T) {
	tests := []struct {
		name string
		a    *Item
		b    *Item
		c    *Item
		k    string
		e    error
	}{
		{
			"CompareAndSwap",
			&Item{
				Key:        "test_cas",
				Value:      []byte("2"),
				Flags:      0,
				Expiration: 60,
			},
			nil,
			&Item{
				Key:        "test_cas",
				Value:      []byte("3"),
				Flags:      0,
				Expiration: 60,
			},
			"test_cas",
			nil,
		},
		{
			"CompareAndSwapErrCASConflict",
			&Item{
				Key:        "test_cas_conflict",
				Value:      []byte("2"),
				Flags:      0,
				Expiration: 60,
			},
			&Item{
				Key:        "test_cas_conflict",
				Value:      []byte("1"),
				Flags:      0,
				Expiration: 60,
			},
			&Item{
				Key:        "test_cas_conflict",
				Value:      []byte("3"),
				Flags:      0,
				Expiration: 60,
			},
			"test_cas_conflict",
			ErrCASConflict,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := testConnASCII.Set(test.a); err != nil {
				t.Fatal(err)
			}
			r, err := testConnASCII.Get(test.k)
			if err != nil {
				t.Fatal(err)
			}

			if test.b != nil {
				if err := testConnASCII.Set(test.b); err != nil {
					t.Fatal(err)
				}
			}

			r.Value = test.c.Value
			if err := testConnASCII.CompareAndSwap(r); err != nil {
				if err != test.e {
					t.Fatal(err)
				}
			} else {
				if fr, err := testConnASCII.Get(test.k); err != nil {
					t.Fatal(err)
				} else {
					compareItem(t, fr, test.c)
				}
			}
		})
	}

	t.Run("TestCompareAndSwapErrNotFound", func(t *testing.T) {
		ti := &Item{
			Key:        "test_cas_notfound",
			Value:      []byte("2"),
			Flags:      0,
			Expiration: 60,
		}
		if err := testConnASCII.Set(ti); err != nil {
			t.Fatal(err)
		}
		r, err := testConnASCII.Get(ti.Key)
		if err != nil {
			t.Fatal(err)
		}

		r.Key = "test_cas_notfound_boom"
		r.Value = []byte("3")
		if err := testConnASCII.CompareAndSwap(r); err != nil {
			if err != ErrNotFound {
				t.Fatal(err)
			}
		}
	})
}

func TestASCIIConnReplace(t *testing.T) {
	tests := []struct {
		name string
		a    *Item
		b    *Item
		e    error
	}{
		{
			"TestReplace",
			&Item{
				Key:        "test_replace",
				Value:      []byte("2"),
				Flags:      0,
				Expiration: 60,
			},
			&Item{
				Key:        "test_replace",
				Value:      []byte("3"),
				Flags:      0,
				Expiration: 60,
			},
			nil,
		},
		{
			"TestReplaceErrNotStored",
			&Item{
				Key:        "test_replace_not_stored",
				Value:      []byte("2"),
				Flags:      0,
				Expiration: 60,
			},
			&Item{
				Key:        "test_replace_not_stored_boom",
				Value:      []byte("3"),
				Flags:      0,
				Expiration: 60,
			},
			ErrNotStored,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := testConnASCII.Set(test.a); err != nil {
				t.Fatal(err)
			}
			if err := testConnASCII.Replace(test.b); err != nil {
				if err == test.e {
					return
				}
				t.Fatal(err)
			}
			if r, err := testConnASCII.Get(test.b.Key); err != nil {
				t.Fatal(err)
			} else {
				compareItem(t, r, test.b)
			}
		})
	}
}

func TestASCIIConnIncrDecr(t *testing.T) {
	tests := []struct {
		fn   func(key string, delta uint64) (uint64, error)
		name string
		k    string
		v    uint64
		w    uint64
	}{
		{
			testConnASCII.Increment,
			"Incr_10",
			"test_incr",
			10,
			10,
		},
		{
			testConnASCII.Increment,
			"Incr_10(2)",
			"test_incr",
			10,
			20,
		},
		{
			testConnASCII.Decrement,
			"Decr_10",
			"test_incr",
			10,
			10,
		},
	}
	if err := testConnASCII.Add(&Item{
		Key:   "test_incr",
		Value: []byte("0"),
	}); err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if a, err := test.fn(test.k, test.v); err != nil {
				t.Fatal(err)
			} else {
				if a != test.w {
					t.Fatalf("want %d, got %d", test.w, a)
				}
			}
			if b, err := testConnASCII.Get(test.k); err != nil {
				t.Fatal(err)
			} else {
				if string(b.Value) != strconv.FormatUint(test.w, 10) {
					t.Fatalf("want %s, got %d", b.Value, test.w)
				}
			}
		})
	}
}

func TestASCIIConnTouch(t *testing.T) {
	tests := []struct {
		name string
		k    string
		a    *Item
		e    error
	}{
		{
			"Touch",
			"test_touch",
			&Item{
				Key:        "test_touch",
				Value:      []byte("0"),
				Expiration: 60,
			},
			nil,
		},
		{
			"Touch_NotExist",
			"test_touch_not_exist",
			nil,
			ErrNotFound,
		},
	}
	for _, test := range tests {
		if test.a != nil {
			if err := testConnASCII.Add(test.a); err != nil {
				t.Fatal(err)
			}
			if err := testConnASCII.Touch(test.k, 1); err != test.e {
				t.Fatal(err)
			}
		}
	}
}

func TestASCIIConnDelete(t *testing.T) {
	tests := []struct {
		name string
		k    string
		a    *Item
		e    error
	}{
		{
			"Delete",
			"test_delete",
			&Item{
				Key:        "test_delete",
				Value:      []byte("0"),
				Expiration: 60,
			},
			nil,
		},
		{
			"Delete_NotExist",
			"test_delete_not_exist",
			nil,
			ErrNotFound,
		},
	}
	for _, test := range tests {
		if test.a != nil {
			if err := testConnASCII.Add(test.a); err != nil {
				t.Fatal(err)
			}
			if err := testConnASCII.Delete(test.k); err != test.e {
				t.Fatal(err)
			}
			if _, err := testConnASCII.Get(test.k); err != ErrNotFound {
				t.Fatal(err)
			}
		}
	}
}

func compareItem(t *testing.T, a, b *Item) {
	if a.Key != b.Key || !bytes.Equal(a.Value, b.Value) || a.Flags != b.Flags {
		t.Fatalf("compareItem: a(%s, %d, %d) : b(%s, %d, %d)", a.Key, len(a.Value), a.Flags, b.Key, len(b.Value), b.Flags)
	}
}
