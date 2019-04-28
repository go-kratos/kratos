package external

import (
	"context"
	"encoding/json"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-feed-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func TestService_DynamicNew(t *testing.T) {
	type args struct {
		c      context.Context
		params string
	}
	tests := []struct {
		name    string
		args    args
		wantRes json.RawMessage
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotRes, err := s.DynamicNew(tt.args.c, tt.args.params)
			So(gotRes, ShouldEqual, tt.wantRes)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestService_DynamicCount(t *testing.T) {
	type args struct {
		c      context.Context
		params string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes json.RawMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.DynamicCount(tt.args.c, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DynamicCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.DynamicCount() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestService_DynamicHistory(t *testing.T) {
	type args struct {
		c      context.Context
		params string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantRes json.RawMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.DynamicHistory(tt.args.c, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DynamicHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.DynamicHistory() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
