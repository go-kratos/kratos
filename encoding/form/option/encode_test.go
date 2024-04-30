package option

import "testing"

func TestEncodeOption_UseProtoTextAsKey(t *testing.T) {
	e := Encode().UseProtoTextAsKey(true)
	if e.ForceProtoTextAsKey != true {
		t.Error("expect true")
	}
	e = Encode().UseProtoTextAsKey(false)
	if e.ForceProtoTextAsKey != false {
		t.Error("expect false")
	}
}

func TestMergeEncodeOptions(t *testing.T) {
	opt := MergeEncodeOptions(nil)
	if opt == nil {
		t.Fatal("expect not nil")
	}
	if opt.ForceProtoTextAsKey {
		t.Error("expect false")
	}

	opt = MergeEncodeOptions(Encode().UseProtoTextAsKey(true))
	if opt == nil {
		t.Fatal("expect not nil")
	}
	if !opt.ForceProtoTextAsKey {
		t.Error("expect true")
	}

	opt = MergeEncodeOptions(&EncodeOption{}, Encode().UseProtoTextAsKey(true))
	if opt == nil {
		t.Fatal("expect not nil")
	}
	if !opt.ForceProtoTextAsKey {
		t.Error("expect true")
	}
}
