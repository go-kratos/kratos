package host

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetQuery(t *testing.T) {
	type args struct {
		url *url.URL
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "grpc://127.0.0.1?isSecure=false",
			args: args{&url.URL{Scheme: "grpc", Host: "127.0.0.1", RawQuery: "isSecure=false"}},
			want: false,
		},
		{
			name: "grpc://127.0.0.1?isSecure=true",
			args: args{&url.URL{Scheme: "grpc", Host: "127.0.0.1", RawQuery: "isSecure=true"}},
			want: true,
		},
		{
			name: "grpc://127.0.0.1",
			args: args{&url.URL{Scheme: "grpc", Host: "127.0.0.1"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetQuery(tt.args.url); !reflect.DeepEqual(got.IsSecure, tt.want) {
				t.Errorf("GetQuery() = %v, want %v", got.IsSecure, tt.want)
			}
		})
	}
}
