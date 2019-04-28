package feed

import (
	"context"
	"reflect"
	"testing"

	"go-common/app/service/main/archive/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_ArchivesWithPlayer(t *testing.T) {
	type args struct {
		c         context.Context
		aids      []int64
		qn        int
		platform  string
		fnver     int
		fnval     int
		forceHost int
		build     int
	}
	tests := []struct {
		name    string
		args    args
		wantRes map[int64]*archive.ArchiveWithPlayer
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotRes, err := s.ArchivesWithPlayer(tt.args.c, tt.args.aids, tt.args.qn, tt.args.platform, tt.args.fnver, tt.args.fnval, tt.args.forceHost, tt.args.build)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ArchivesWithPlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.ArchivesWithPlayer() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
