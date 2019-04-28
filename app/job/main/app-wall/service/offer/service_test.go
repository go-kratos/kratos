package offer

import (
	"context"
	"flag"
	"path/filepath"
	"reflect"
	"testing"

	"go-common/app/job/main/app-wall/conf"

	cluster "github.com/bsm/sarama-cluster"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/app-wall-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

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

func TestService_NewConsumer(t *testing.T) {
	tests := []struct {
		name    string
		s       *Service
		want    *cluster.Consumer
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			got, err := tt.s.NewConsumer()
			So(err, ShouldEqual, tt.wantErr)
			So(got, ShouldResemble, tt.want)
		})
	}
}
