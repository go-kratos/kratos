package region

import (
	"reflect"
	"testing"

	"go-common/app/interface/main/app-feed/conf"
)

func TestNew(t *testing.T) {
	type args struct {
		c *conf.Config
	}
	tests := []struct {
		name  string
		args  args
		wantS *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := New(tt.args.c); !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("New() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestService_md5(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		s    *Service
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.md5(tt.args.v); got != tt.want {
				t.Errorf("Service.md5() = %v, want %v", got, tt.want)
			}
		})
	}
}
