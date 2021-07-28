package endpoint

import (
	"net/url"
	"reflect"
	"testing"
)

func TestEndPoint(t *testing.T) {
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
			args: args{NewEndpoint("grpc", "127.0.0.1", false)},
			want: false,
		},
		{
			name: "grpc://127.0.0.1?isSecure=true",
			args: args{NewEndpoint("http", "127.0.0.1", true)},
			want: true,
		},
		{
			name: "grpc://127.0.0.1",
			args: args{NewEndpoint("grpc", "localhost", false)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSecure(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
