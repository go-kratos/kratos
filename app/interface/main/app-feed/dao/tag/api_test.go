package tag

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go-common/app/interface/main/app-feed/model/tag"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Hots(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		rid int16
		now time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantHs  []*tag.Hot
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHs, err := d.Hots(tt.args.c, tt.args.mid, tt.args.rid, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Hots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHs, tt.wantHs) {
				t.Errorf("Dao.Hots() = %v, want %v", gotHs, tt.wantHs)
			}
		})
	}
}

func TestDao_Add(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		tid int64
		now time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			err := d.Add(tt.args.c, tt.args.mid, tt.args.tid, tt.args.now)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}

func TestDao_Cancel(t *testing.T) {
	type args struct {
		c   context.Context
		mid int64
		tid int64
		now time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.Cancel(tt.args.c, tt.args.mid, tt.args.tid, tt.args.now); (err != nil) != tt.wantErr {
				t.Errorf("Dao.Cancel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDao_Tags(t *testing.T) {
	type args struct {
		c    context.Context
		mid  int64
		aids []int64
		now  time.Time
	}
	tests := []struct {
		name     string
		args     args
		wantTagm map[string][]*tag.Tag
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTagm, err := d.Tags(tt.args.c, tt.args.mid, tt.args.aids, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Tags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTagm, tt.wantTagm) {
				t.Errorf("Dao.Tags() = %v, want %v", gotTagm, tt.wantTagm)
			}
		})
	}
}

func TestDao_Detail(t *testing.T) {
	type args struct {
		c     context.Context
		tagID int
		pn    int
		ps    int
		now   time.Time
	}
	tests := []struct {
		name       string
		args       args
		wantArcids []int64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArcids, err := d.Detail(tt.args.c, tt.args.tagID, tt.args.pn, tt.args.ps, tt.args.now)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.Detail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotArcids, tt.wantArcids) {
				t.Errorf("Dao.Detail() = %v, want %v", gotArcids, tt.wantArcids)
			}
		})
	}
}
