package dsn

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	xtime "go-common/library/time"
)

type cfg1 struct {
	Name     string  `dsn:"query.name"`
	Def      string  `dsn:"query.def,hello"`
	DefSlice []int   `dsn:"query.defslice,1,2,3,4"`
	Ignore   string  `dsn:"-"`
	FloatNum float64 `dsn:"query.floatNum"`
}

type cfg2 struct {
	Timeout xtime.Duration `dsn:"query.timeout"`
}

type cfg3 struct {
	Username string         `dsn:"username"`
	Timeout  xtime.Duration `dsn:"query.timeout"`
}

type cfg4 struct {
	Timeout xtime.Duration `dsn:"query.timeout,1s"`
}

func TestDecodeQuery(t *testing.T) {
	type args struct {
		query       url.Values
		v           interface{}
		assignFuncs map[string]assignFunc
	}
	tests := []struct {
		name    string
		args    args
		want    url.Values
		cfg     interface{}
		wantErr bool
	}{
		{
			name: "test generic",
			args: args{
				query: url.Values{
					"name":     {"hello"},
					"Ignore":   {"test"},
					"floatNum": {"22.33"},
					"adb":      {"123"},
				},
				v: &cfg1{},
			},
			want: url.Values{
				"Ignore": {"test"},
				"adb":    {"123"},
			},
			cfg: &cfg1{
				Name:     "hello",
				Def:      "hello",
				DefSlice: []int{1, 2, 3, 4},
				FloatNum: 22.33,
			},
		},
		{
			name: "test go-common/library/time",
			args: args{
				query: url.Values{
					"timeout": {"1s"},
				},
				v: &cfg2{},
			},
			want: url.Values{},
			cfg:  &cfg2{xtime.Duration(time.Second)},
		},
		{
			name: "test empty go-common/library/time",
			args: args{
				query: url.Values{},
				v:     &cfg2{},
			},
			want: url.Values{},
			cfg:  &cfg2{},
		},
		{
			name: "test go-common/library/time",
			args: args{
				query: url.Values{},
				v:     &cfg4{},
			},
			want: url.Values{},
			cfg:  &cfg4{xtime.Duration(time.Second)},
		},
		{
			name: "test build-in value",
			args: args{
				query: url.Values{
					"timeout": {"1s"},
				},
				v:           &cfg3{},
				assignFuncs: map[string]assignFunc{"username": stringsAssignFunc("hello")},
			},
			want: url.Values{},
			cfg: &cfg3{
				Timeout:  xtime.Duration(time.Second),
				Username: "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bindQuery(tt.args.query, tt.args.v, tt.args.assignFuncs)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeQuery() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.args.v, tt.cfg) {
				t.Errorf("DecodeQuery() = %v, want %v", tt.args.v, tt.cfg)
			}
		})
	}
}
