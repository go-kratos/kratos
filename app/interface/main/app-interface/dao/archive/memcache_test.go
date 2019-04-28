package archive

import (
	"context"
	"go-common/app/service/main/archive/api"
	"reflect"
	"testing"
)

func Test_keyArc(t *testing.T) {
	type args struct {
		aid int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyArc(tt.args.aid); got != tt.want {
				t.Errorf("keyArc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keyStat(t *testing.T) {
	type args struct {
		aid int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyStat(tt.args.aid); got != tt.want {
				t.Errorf("keyStat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDao_arcsCache(t *testing.T) {
	type args struct {
		c    context.Context
		aids []int64
	}
	tests := []struct {
		name       string
		d          *Dao
		args       args
		wantCached map[int64]*api.Arc
		wantMissed []int64
		wantErr    bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCached, gotMissed, err := tt.d.arcsCache(tt.args.c, tt.args.aids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.arcsCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCached, tt.wantCached) {
				t.Errorf("Dao.arcsCache() gotCached = %v, want %v", gotCached, tt.wantCached)
			}
			if !reflect.DeepEqual(gotMissed, tt.wantMissed) {
				t.Errorf("Dao.arcsCache() gotMissed = %v, want %v", gotMissed, tt.wantMissed)
			}
		})
	}
}

func TestDao_statsCache(t *testing.T) {
	type args struct {
		c    context.Context
		aids []int64
	}
	tests := []struct {
		name       string
		d          *Dao
		args       args
		wantCached map[int64]*api.Stat
		wantMissed []int64
		wantErr    bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCached, gotMissed, err := tt.d.statsCache(tt.args.c, tt.args.aids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.statsCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCached, tt.wantCached) {
				t.Errorf("Dao.statsCache() gotCached = %v, want %v", gotCached, tt.wantCached)
			}
			if !reflect.DeepEqual(gotMissed, tt.wantMissed) {
				t.Errorf("Dao.statsCache() gotMissed = %v, want %v", gotMissed, tt.wantMissed)
			}
		})
	}
}

func TestDao_avWithStCaches(t *testing.T) {
	type args struct {
		c    context.Context
		aids []int64
	}
	tests := []struct {
		name         string
		d            *Dao
		args         args
		wantCached   map[int64]*api.Arc
		wantAvMissed []int64
		wantStMissed []int64
		wantErr      bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCached, gotAvMissed, gotStMissed, err := tt.d.avWithStCaches(tt.args.c, tt.args.aids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.avWithStCaches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCached, tt.wantCached) {
				t.Errorf("Dao.avWithStCaches() gotCached = %v, want %v", gotCached, tt.wantCached)
			}
			if !reflect.DeepEqual(gotAvMissed, tt.wantAvMissed) {
				t.Errorf("Dao.avWithStCaches() gotAvMissed = %v, want %v", gotAvMissed, tt.wantAvMissed)
			}
			if !reflect.DeepEqual(gotStMissed, tt.wantStMissed) {
				t.Errorf("Dao.avWithStCaches() gotStMissed = %v, want %v", gotStMissed, tt.wantStMissed)
			}
		})
	}
}
