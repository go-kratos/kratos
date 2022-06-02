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
			if got := IsSecure(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEndpoint(t *testing.T) {
	type args struct {
		scheme string
		host   string
	}
	tests := []struct {
		name string
		args args
		want *url.URL
	}{
		{
			name: "https://github.com/go-kratos/kratos/",
			args: args{"https", "github.com/go-kratos/kratos/"},
			want: &url.URL{Scheme: "https", Host: "github.com/go-kratos/kratos/"},
		},
		{
			name: "https://go-kratos.dev/",
			args: args{"https", "go-kratos.dev/"},
			want: &url.URL{Scheme: "https", Host: "go-kratos.dev/"},
		},
		{
			name: "https://www.google.com/",
			args: args{"https", "www.google.com/"},
			want: &url.URL{Scheme: "https", Host: "www.google.com/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEndpoint(tt.args.scheme, tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseEndpoint(t *testing.T) {
	type args struct {
		endpoints []string
		scheme    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "kratos",
			args:    args{endpoints: []string{"https://github.com/go-kratos/kratos"}, scheme: "https"},
			want:    "github.com",
			wantErr: false,
		},
		{
			name:    "test",
			args:    args{endpoints: []string{"http://go-kratos.dev/"}, scheme: "https"},
			want:    "",
			wantErr: false,
		},
		{
			name:    "localhost:8080",
			args:    args{endpoints: []string{"grpcs://localhost:8080/"}, scheme: "grpcs"},
			want:    "localhost:8080",
			wantErr: false,
		},
		{
			name:    "localhost:8081",
			args:    args{endpoints: []string{"grpcs://localhost:8080/"}, scheme: "grpc"},
			want:    "",
			wantErr: false,
		},

		// Legacy
		{
			name:    "google",
			args:    args{endpoints: []string{"grpc://www.google.com/?isSecure=true"}, scheme: "grpcs"},
			want:    "www.google.com",
			wantErr: false,
		},
		{
			name:    "baidu",
			args:    args{endpoints: []string{"http://www.baidu.com/?isSecure=true"}, scheme: "https"},
			want:    "www.baidu.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEndpoint(tt.args.endpoints, tt.args.scheme)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseEndpoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}
