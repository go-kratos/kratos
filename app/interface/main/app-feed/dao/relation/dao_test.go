package relation

import (
	. "github.com/smartystreets/goconvey/convey"

	"context"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-feed/conf"
	relation "go-common/app/service/main/relation/model"
)

func TestNew(t *testing.T) {
	type args struct {
		c *conf.Config
	}
	tests := []struct {
		name  string
		args  args
		wantD *Dao
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotD := New(tt.args.c)
			So(gotD, ShouldEqual, tt.wantD)
		})
	}
}

func TestDao_Stats(t *testing.T) {
	type args struct {
		ctx  context.Context
		mids []int64
	}
	tests := []struct {
		name    string
		d       *Dao
		args    args
		wantRes map[int64]*relation.Stat
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.d.Stats(tt.args.ctx, tt.args.mids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Stats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Dao.Stats() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
