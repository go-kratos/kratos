package black

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/interface/main/app-feed/conf"
	"go-common/library/cache/redis"
	httpx "go-common/library/net/http/blademaster"

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
		Convey(tt.name, func(t *testing.T) {
			gotD := New(tt.args.c)
			So(gotD, ShouldEqual, tt.wantD)
		})
	}
}

func TestDao_Ping(t *testing.T) {
	type fields struct {
		clientAsyn *httpx.Client
		redis      *redis.Pool
		expireRds  int32
		aCh        chan func()
	}
	type args struct {
		c context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				clientAsyn: tt.fields.clientAsyn,
				redis:      tt.fields.redis,
				expireRds:  tt.fields.expireRds,
				aCh:        tt.fields.aCh,
			}
			if err := d.Ping(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Dao.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDao_AddBlacklist(t *testing.T) {
	type fields struct {
		clientAsyn *httpx.Client
		redis      *redis.Pool
		expireRds  int32
		aCh        chan func()
	}
	type args struct {
		mid int64
		aid int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				clientAsyn: tt.fields.clientAsyn,
				redis:      tt.fields.redis,
				expireRds:  tt.fields.expireRds,
				aCh:        tt.fields.aCh,
			}
			d.AddBlacklist(tt.args.mid, tt.args.aid)
		})
	}
}

func TestDao_DelBlacklist(t *testing.T) {
	type fields struct {
		clientAsyn *httpx.Client
		redis      *redis.Pool
		expireRds  int32
		aCh        chan func()
	}
	type args struct {
		mid int64
		aid int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				clientAsyn: tt.fields.clientAsyn,
				redis:      tt.fields.redis,
				expireRds:  tt.fields.expireRds,
				aCh:        tt.fields.aCh,
			}
			d.DelBlacklist(tt.args.mid, tt.args.aid)
		})
	}
}

func TestDao_BlackList(t *testing.T) {
	type fields struct {
		clientAsyn *httpx.Client
		redis      *redis.Pool
		expireRds  int32
		aCh        chan func()
	}
	type args struct {
		c   context.Context
		mid int64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantAidm map[int64]struct{}
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				clientAsyn: tt.fields.clientAsyn,
				redis:      tt.fields.redis,
				expireRds:  tt.fields.expireRds,
				aCh:        tt.fields.aCh,
			}
			gotAidm, err := d.BlackList(tt.args.c, tt.args.mid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dao.BlackList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAidm, tt.wantAidm) {
				t.Errorf("Dao.BlackList() = %v, want %v", gotAidm, tt.wantAidm)
			}
		})
	}
}

func TestDao_addCache(t *testing.T) {
	type fields struct {
		clientAsyn *httpx.Client
		redis      *redis.Pool
		expireRds  int32
		aCh        chan func()
	}
	type args struct {
		i func()
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				clientAsyn: tt.fields.clientAsyn,
				redis:      tt.fields.redis,
				expireRds:  tt.fields.expireRds,
				aCh:        tt.fields.aCh,
			}
			d.addCache(tt.args.i)
		})
	}
}

func TestDao_cacheproc(t *testing.T) {
	type fields struct {
		clientAsyn *httpx.Client
		redis      *redis.Pool
		expireRds  int32
		aCh        chan func()
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				clientAsyn: tt.fields.clientAsyn,
				redis:      tt.fields.redis,
				expireRds:  tt.fields.expireRds,
				aCh:        tt.fields.aCh,
			}
			d.cacheproc()
		})
	}
}
