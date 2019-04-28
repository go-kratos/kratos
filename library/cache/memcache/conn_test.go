package memcache

import (
	"bytes"
	"encoding/json"
	"errors"
	test "go-common/library/cache/memcache/test"
	"testing"
	"time"

	"github.com/bouk/monkey"
	"github.com/gogo/protobuf/proto"
)

var s = []string{"test", "test1"}
var c Conn

var item = &Item{
	Key:        "test",
	Value:      []byte("test"),
	Flags:      FlagRAW,
	Expiration: 60,
	cas:        0,
}

var item2 = &Item{
	Key:        "test1",
	Value:      []byte("test"),
	Flags:      0,
	Expiration: 1000,
	cas:        0,
}

var item3 = &Item{
	Key:        "test2",
	Value:      []byte("0"),
	Flags:      0,
	Expiration: 60,
	cas:        0,
}

type TestItem struct {
	Name string
	Age  int
}

func (t *TestItem) Compare(nt *TestItem) bool {
	return t.Name == nt.Name && t.Age == nt.Age
}

func prepareEnv(t *testing.T) {
	if c != nil {
		return
	}
	var err error
	cnop := DialConnectTimeout(time.Duration(2 * time.Second))
	rdop := DialReadTimeout(time.Duration(2 * time.Second))
	wrop := DialWriteTimeout(time.Duration(2 * time.Second))
	c, err = Dial("tcp", testMemcacheAddr, cnop, rdop, wrop)
	if err != nil {
		t.Errorf("Dial() error(%v)", err)
		t.FailNow()
	}
	c.Delete("test")
	c.Delete("test1")
	c.Delete("test2")
}

func TestRaw(t *testing.T) {
	prepareEnv(t)
	if err := c.Set(item); err != nil {
		t.Errorf("conn.Store() error(%v)", err)
	}
}

func TestAdd(t *testing.T) {
	var (
		key  = "test_add"
		item = &Item{
			Key:        key,
			Value:      []byte("0"),
			Flags:      0,
			Expiration: 60,
			cas:        0,
		}
	)
	prepareEnv(t)
	c.Delete(key)
	if err := c.Add(item); err != nil {
		t.Errorf("c.Add() error(%v)", err)
	}
	if err := c.Add(item); err != ErrNotStored {
		t.Errorf("c.Add() error(%v)", err)
	}
}

func TestSetErr(t *testing.T) {
	prepareEnv(t)
	//set
	st := &TestItem{Name: "jsongzip", Age: 10}
	itemx := &Item{Key: "jsongzip", Object: st}
	if err := c.Set(itemx); err != ErrItem {
		t.Errorf("conn.Set() error(%v)", err)
	}
}

func TestSetErr2(t *testing.T) {
	prepareEnv(t)
	//set
	itemx := &Item{Key: "jsongzip", Flags: FlagJSON | FlagGzip}
	if err := c.Set(itemx); err != ErrItem {
		t.Errorf("conn.Set() error(%v)", err)
	}
}

func TestSetErr3(t *testing.T) {
	prepareEnv(t)
	//set
	itemx := &Item{Key: "jsongzip", Value: []byte("test"), Flags: FlagJSON}
	if err := c.Set(itemx); err != ErrItem {
		t.Errorf("conn.Set() error(%v)", err)
	}
}

func TestJSONGzip(t *testing.T) {
	prepareEnv(t)
	//set
	st := &TestItem{Name: "jsongzip", Age: 10}
	itemx := &Item{Key: "jsongzip", Object: st, Flags: FlagJSON | FlagGzip}
	if err := c.Set(itemx); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err := c.Get("jsongzip"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		var nst TestItem
		scanAndCompare(t, r, st, &nst)
	}
}

func TestJSON(t *testing.T) {
	prepareEnv(t)
	st := &TestItem{Name: "json", Age: 10}
	itemx := &Item{Key: "json", Object: st, Flags: FlagJSON}
	if err := c.Set(itemx); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err := c.Get("json"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		var nst TestItem
		scanAndCompare(t, r, st, &nst)
	}
}

func BenchmarkJSON(b *testing.B) {
	st := &TestItem{Name: "json", Age: 10}
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

func BenchmarkProtobuf(b *testing.B) {
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

func TestGob(t *testing.T) {
	prepareEnv(t)
	st := &TestItem{Name: "gob", Age: 10}
	itemx := &Item{Key: "gob", Object: st, Flags: FlagGOB}
	if err := c.Set(itemx); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err := c.Get("gob"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		var nst TestItem
		scanAndCompare(t, r, st, &nst)
	}
}

func TestGobGzip(t *testing.T) {
	prepareEnv(t)
	st := &TestItem{Name: "gobgzip", Age: 10}
	itemx := &Item{Key: "gobgzip", Object: st, Flags: FlagGOB | FlagGzip}
	if err := c.Set(itemx); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err := c.Get("gobgzip"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		var nst TestItem
		scanAndCompare(t, r, st, &nst)
	}
}

func TestGzip(t *testing.T) {
	prepareEnv(t)
	st := &TestItem{Name: "gzip", Age: 123}
	itemx := &Item{Key: "gzip", Object: st, Flags: FlagGOB | FlagGzip}
	if err := c.Set(itemx); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err := c.Get("gzip"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		var nst TestItem
		scanAndCompare(t, r, st, &nst)
	}
}

func TestProtobuf(t *testing.T) {
	prepareEnv(t)
	var (
		err error
		// value []byte
		r   *Item
		nst test.TestItem
	)
	st := &test.TestItem{Name: "proto", Age: 3021}
	itemx := &Item{Key: "proto", Object: st, Flags: FlagProtobuf}
	if err = c.Set(itemx); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err = c.Get("proto"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	}
	if err = c.Scan(r, &nst); err != nil {
		t.Errorf("decode() error(%v)", err)
		t.FailNow()
	} else {
		scanAndCompare2(t, r, st, &nst)
	}
}

func TestProtobufGzip(t *testing.T) {
	prepareEnv(t)
	var (
		err error
		// value []byte
		r   *Item
		nst test.TestItem
	)
	st := &test.TestItem{Name: "protogzip", Age: 3021}
	itemx := &Item{Key: "protogzip", Object: st, Flags: FlagProtobuf | FlagGzip}
	if err = c.Set(itemx); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err = c.Get("protogzip"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	}
	if err = c.Scan(r, &nst); err != nil {
		t.Errorf("decode() error(%v)", err)
		t.FailNow()
	} else {
		scanAndCompare2(t, r, st, &nst)
	}
}

func TestGet(t *testing.T) {
	prepareEnv(t)
	// get
	if r, err := c.Get("test"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else if r.Key != "test" || !bytes.Equal(r.Value, []byte("test")) || r.Flags != 0 {
		t.Error("conn.Get() error, value")
	}
}

func TestGetHasErr(t *testing.T) {
	prepareEnv(t)

	st := &TestItem{Name: "json", Age: 10}
	itemx := &Item{Key: "test", Object: st, Flags: FlagJSON}
	c.Set(itemx)

	expected := errors.New("some error")
	monkey.Patch(scanGetReply, func(line []byte, item *Item) (size int, err error) {
		return 0, expected
	})

	if _, err := c.Get("test"); err.Error() != expected.Error() {
		t.Errorf("conn.Get() unexpected error(%v)", err)
	}
	if err := c.(*conn).err; err.Error() != expected.Error() {
		t.Errorf("unexpected error(%v)", err)
	}
}

func TestGet2(t *testing.T) {
	prepareEnv(t)
	// get not exist
	if _, err := c.Get("not_exist"); err != ErrNotFound {
		t.Errorf("conn.Get() error(%v)", err)
	}
}

func TestGetMulti(t *testing.T) {
	prepareEnv(t)
	// getMulti
	if _, err := c.GetMulti(s); err != nil {
		t.Errorf("conn.GetMulti() error(%v)", err)
	}
}

func TestGetMulti2(t *testing.T) {
	prepareEnv(t)
	//set
	if err := c.Set(item); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if err := c.Set(item2); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if res, err := c.GetMulti(s); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		if len(res) != 2 {
			t.Error("conn.Get() error, length", len(res))
		}
		reply := res["test"]
		compareItem2(t, reply, item)
		reply = res["test1"]
		compareItem2(t, reply, item2)
	}
}

func TestIncrement(t *testing.T) {
	// set
	if err := c.Set(item3); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	// incr
	if d, err := c.Increment("test2", 4); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	} else {
		if d != 4 {
			t.Error("conn.IncrDecr value error")
		}
	}
}

func TestDecrement(t *testing.T) {
	// decr
	if d, err := c.Decrement("test2", 3); err != nil {
		t.Errorf("conn.Store() error(%v)", err)
	} else {
		if d != 1 {
			t.Error("conn.IncrDecr value error", d)
		}
	}
}

func TestTouch(t *testing.T) {
	// touch
	if err := c.Touch("test2", 1); err != nil {
		t.Errorf("conn.Touch error(%v)", err)
	}
}

func TestCompareAndSwap(t *testing.T) {
	prepareEnv(t)
	if err := c.Set(item3); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	//cas
	if r, err := c.Get("test2"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		r.Value = []byte("fuck")
		if err := c.CompareAndSwap(r); err != nil {
			t.Errorf("conn.CompareAndSwap() error(%v)", err)
		}
		if r, err := c.Get("test2"); err != nil {
			t.Errorf("conn.Get() error(%v)", err)
		} else {
			itemx := &Item{Key: "test2", Value: []byte("fuck"), Flags: 0}
			compareItem2(t, r, itemx)
		}
	}
}

func TestReplace(t *testing.T) {
	prepareEnv(t)
	if err := c.Set(item); err != nil {
		t.Errorf("conn.Set() error(%v)", err)
	}
	if r, err := c.Get("test"); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		r.Value = []byte("go")
		if err := c.Replace(r); err != nil {
			t.Errorf("conn.CompareAndSwap() error(%v)", err)
		}
		if r, err := c.Get("test"); err != nil {
			t.Errorf("conn.Get() error(%v)", err)
		} else {
			itemx := &Item{Key: "test", Value: []byte("go"), Flags: 0}
			compareItem2(t, r, itemx)
		}

	}
}

func scanAndCompare(t *testing.T, item *Item, st *TestItem, nst *TestItem) {
	if err := c.Scan(item, nst); err != nil {
		t.Errorf("decode() error(%v)", err)
		t.FailNow()
	}
	if !st.Compare(nst) {
		t.Errorf("st: %v, use of closed network connection nst: %v", st, &nst)
		t.FailNow()
	}
}

func scanAndCompare2(t *testing.T, item *Item, st *test.TestItem, nst *test.TestItem) {
	if err := c.Scan(item, nst); err != nil {
		t.Errorf("decode() error(%v)", err)
		t.FailNow()
	}
	if st.Age != nst.Age || st.Name != nst.Name {
		t.Errorf("st: %v, use of closed network connection nst: %v", st, &nst)
		t.FailNow()
	}
}

func compareItem2(t *testing.T, r, item *Item) {
	if r.Key != item.Key || !bytes.Equal(r.Value, item.Value) || r.Flags != item.Flags {
		t.Error("conn.Get() error, value")
	}
}

func Test_legalKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test empty key",
			want: false,
		},
		{
			name: "test too large key",
			args: args{func() string {
				var data []byte
				for i := 0; i < 255; i++ {
					data = append(data, 'k')
				}
				return string(data)
			}()},
			want: false,
		},
		{
			name: "test invalid char",
			args: args{"hello world"},
			want: false,
		},
		{
			name: "test invalid char",
			args: args{string([]byte{0x7f})},
			want: false,
		},
		{
			name: "test normal key",
			args: args{"hello"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := legalKey(tt.args.key); got != tt.want {
				t.Errorf("legalKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
