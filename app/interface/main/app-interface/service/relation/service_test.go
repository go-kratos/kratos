package relation

import (
	"reflect"
	"testing"

	"go-common/app/interface/main/app-interface/conf"
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
		t.Run(tt.name, func(t *testing.T) {
			if gotS := New(tt.args.c); !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("New() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
