package warden

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go-common/library/log"
)

func Test_logFn(t *testing.T) {
	type args struct {
		code int
		dt   time.Duration
	}
	tests := []struct {
		name string
		args args
		want func(context.Context, ...log.D)
	}{
		{
			name: "ok",
			args: args{code: 0, dt: time.Millisecond},
			want: log.Infov,
		},
		{
			name: "slowlog",
			args: args{code: 0, dt: time.Second},
			want: log.Warnv,
		},
		{
			name: "business error",
			args: args{code: 2233, dt: time.Millisecond},
			want: log.Warnv,
		},
		{
			name: "system error",
			args: args{code: -1, dt: 0},
			want: log.Errorv,
		},
		{
			name: "system error and slowlog",
			args: args{code: -1, dt: time.Second},
			want: log.Errorv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logFn(tt.args.code, tt.args.dt); reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("unexpect log function!")
			}
		})
	}
}
