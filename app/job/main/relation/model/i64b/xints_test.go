package i64b

import (
	"testing"
)

func TestMarshalAndUnmarshal(t *testing.T) {
	var a = Int64Bytes{1, 2, 3}
	data := make([]byte, a.Size())
	n, err := a.MarshalTo(data)
	if n != 24 {
		t.Logf("marshal size must be 24")
		t.FailNow()
	}
	if err != nil {
		t.Fatalf("err:%v", err)
	}

	var b Int64Bytes
	err = b.Unmarshal(data)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	if b[0] != 1 || b[1] != 2 || b[2] != 3 {
		t.Logf("unmarshal failed!b:%v", b)
		t.FailNow()
	}
}

func TestUncompleteMarshal(t *testing.T) {
	var a = Int64Bytes{1, 2, 3}
	data := make([]byte, a.Size()-7)
	n, err := a.MarshalTo(data)
	if n != 16 {
		t.Logf("marshal size must be 16")
		t.FailNow()
	}
	if err != nil {
		t.Fatalf("err:%v", err)
	}

	var b Int64Bytes
	err = b.Unmarshal(data)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	if b[0] != 1 || b[1] != 2 {
		t.Logf("unmarshal failed!b:%v", b)
		t.FailNow()
	}
}

func TestNilMarshal(t *testing.T) {
	var a = Int64Bytes{1, 2, 3}
	var data []byte
	n, err := a.MarshalTo(data)
	if n != 0 {
		t.Logf("marshal size must be 0")
		t.FailNow()
	}
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	var b Int64Bytes
	err = b.Unmarshal(data)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	if b != nil {
		t.Logf("unmarshal failed!b:%v", b)
		t.FailNow()
	}

}
