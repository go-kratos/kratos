package service

import (
	"context"
	"go-common/app/job/main/app-player/conf"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
		Convey(tt.name, t, func(t *testing.T) {
			gotS := New(tt.args.c)
			So(gotS, ShouldEqual, tt.wantS)
		})
	}
}

func TestService_Close(t *testing.T) {
	tests := []struct {
		name string
		s    *Service
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Close()
		})
	}
}

func TestService_Ping(t *testing.T) {
	type args struct {
		c context.Context
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Ping(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Service.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
