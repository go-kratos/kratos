package account

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-feed/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	type args struct {
		c *conf.Config
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, func(t *testing.T) {
			gotD := New(tt.args.c)
			So(gotD, ShouldNotBeNil)
		})
	}
}

func TestDao_Relations2(t *testing.T) {
	type args struct {
		c      context.Context
		owners []int64
		mid    int64
	}
	tests := []struct {
		name        string
		d           *Dao
		args        args
		wantFollows map[int64]bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotFollows := tt.d.Relations3(tt.args.c, tt.args.owners, tt.args.mid); !reflect.DeepEqual(gotFollows, tt.wantFollows) {
				t.Errorf("Dao.Relations2() = %v, want %v", gotFollows, tt.wantFollows)
			}
		})
	}
}
