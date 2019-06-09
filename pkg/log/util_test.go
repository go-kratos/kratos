package log

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestFuncName(t *testing.T) {
	name := funcName(1)
	if !strings.Contains(name, "util_test.go:11") {
		t.Errorf("expect contains util_test.go:11 got %s", name)
	}
}

func Test_toMap(t *testing.T) {
	type args struct {
		args []D
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			args: args{[]D{KVString("test", "hello")}},
			want: map[string]interface{}{"test": "hello"},
		},
		{
			args: args{[]D{KVInt64("test", 123)}},
			want: map[string]interface{}{"test": int64(123)},
		},
		{
			args: args{[]D{KVFloat32("test", float32(1.01))}},
			want: map[string]interface{}{"test": float32(1.01)},
		},
		{
			args: args{[]D{KVFloat32("test", float32(1.01))}},
			want: map[string]interface{}{"test": float32(1.01)},
		},
		{
			args: args{[]D{KVDuration("test", time.Second)}},
			want: map[string]interface{}{"test": time.Second},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toMap(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
