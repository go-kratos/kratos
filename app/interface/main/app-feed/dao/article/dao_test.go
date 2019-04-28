package article

import (
	"context"
	"go-common/app/interface/main/app-feed/conf"
	article "go-common/app/interface/openplatform/article/model"
	"reflect"
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
		wantD *Dao
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotD := New(tt.args.c); !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("New() = %v, want %v", gotD, tt.wantD)
			}
		})
	}
	Convey("new", t, func() {
		So(t, ShouldNotBeNil)
	})
}

func TestDao_Articles(t *testing.T) {
	type args struct {
		c    context.Context
		aids []int64
	}
	tests := []struct {
		name    string
		d       *Dao
		args    args
		wantMs  map[int64]*article.Meta
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMs, err := tt.d.Articles(tt.args.c, tt.args.aids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Articles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMs, tt.wantMs) {
				t.Errorf("Dao.Articles() = %v, want %v", gotMs, tt.wantMs)
			}
		})
	}
}
