package toml

import (
	"reflect"
	"testing"
)

type TestCfg struct {
	Etcd struct {
		Endpoints []string `toml:"Endpoints" `
	} `toml:"Etcd" `
}

var config = TestCfg{
	Etcd: struct {
		Endpoints []string `toml:"Endpoints" `
	}{[]string{"127.0.0.1:2379"}},
}

func Test_codec_Marshal(t *testing.T) {
	type args struct {
		v TestCfg
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Marshal",
			args: args{v: config},
			want: []byte(`[Etcd]
Endpoints = ['127.0.0.1:2379']
`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := codec{}
			got, err := co.Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), string(tt.want)) {
				t.Errorf("Marshal() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_codec_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: Name,
			want: Name,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := codec{}
			if got := co.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_codec_Unmarshal(t *testing.T) {
	var cfg TestCfg

	type args struct {
		data []byte
		v    TestCfg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    TestCfg
	}{
		{
			name: "",
			args: args{data: []byte(`
[Etcd]
Endpoints = ["127.0.0.1:2379"]
`), v: cfg},
			want: config,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			co := codec{}
			if err := co.Unmarshal(tt.args.data, &tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", tt.args.v, tt.want)
			}
		})
	}
}
