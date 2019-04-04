package paladin

import (
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	type args struct {
		v   *Value
		def bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{v: &Value{val: true}},
			want: true,
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: false,
		},
		{
			name: "default",
			args: args{v: &Value{}, def: true},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bool(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	type args struct {
		v   *Value
		def int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "ok",
			args: args{v: &Value{val: int64(2233)}},
			want: 2233,
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: 0,
		},
		{
			name: "default",
			args: args{v: &Value{}, def: 2233},
			want: 2233,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt32(t *testing.T) {
	type args struct {
		v   *Value
		def int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "ok",
			args: args{v: &Value{val: int64(2233)}},
			want: 2233,
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: 0,
		},
		{
			name: "default",
			args: args{v: &Value{}, def: 2233},
			want: 2233,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int32(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("Int32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	type args struct {
		v   *Value
		def int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "ok",
			args: args{v: &Value{val: int64(2233)}},
			want: 2233,
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: 0,
		},
		{
			name: "default",
			args: args{v: &Value{}, def: 2233},
			want: 2233,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32(t *testing.T) {
	type args struct {
		v   *Value
		def float32
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			name: "ok",
			args: args{v: &Value{val: float64(2233)}},
			want: float32(2233),
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: 0,
		},
		{
			name: "default",
			args: args{v: &Value{}, def: float32(2233)},
			want: float32(2233),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float32(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("Float32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	type args struct {
		v   *Value
		def float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "ok",
			args: args{v: &Value{val: float64(2233)}},
			want: float64(2233),
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: 0,
		},
		{
			name: "default",
			args: args{v: &Value{}, def: float64(2233)},
			want: float64(2233),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		v   *Value
		def string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{v: &Value{val: "test"}},
			want: "test",
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: "",
		},
		{
			name: "default",
			args: args{v: &Value{}, def: "test"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	type args struct {
		v   *Value
		def time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "ok",
			args: args{v: &Value{val: "1s"}},
			want: time.Second,
		},
		{
			name: "fail",
			args: args{v: &Value{}},
			want: 0,
		},
		{
			name: "default",
			args: args{v: &Value{}, def: time.Second},
			want: time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Duration(tt.args.v, tt.args.def); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
