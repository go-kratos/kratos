package binding

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestBindQuery(t *testing.T) {
	type TestBind struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	p1 := TestBind{}
	type args struct {
		vars   url.Values
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "test",
			args: args{
				vars:   map[string][]string{"name": {"kratos"}, "url": {"https://go-kratos.dev/"}},
				target: &p1,
			},
			wantErr: false,
			want:    &TestBind{"kratos", "https://go-kratos.dev/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BindQuery(tt.args.vars, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("BindQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("BindQuery() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}

func TestBindForm(t *testing.T) {
	type TestBind struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	p1 := TestBind{}
	type args struct {
		req    *http.Request
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "error not nil",
			args: args{
				req:    &http.Request{Method: "POST"},
				target: &p1,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "error is nil",
			args: args{
				req: &http.Request{
					Method: "POST",
					Header: http.Header{"Content-Type": {"application/x-www-form-urlencoded; param=value"}},
					Body:   io.NopCloser(strings.NewReader("name=kratos&url=https://go-kratos.dev/")),
				},
				target: &p1,
			},
			wantErr: false,
			want:    &TestBind{"kratos", "https://go-kratos.dev/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := BindForm(tt.args.req, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("BindForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("BindForm() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}
