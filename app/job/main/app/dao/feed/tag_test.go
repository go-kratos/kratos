package feed

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/app/model/feed"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Tags(t *testing.T) {
	type args struct {
		c    context.Context
		aids []int64
		now  time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantTm  map[string][]*feed.Tag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			gotTm, err := d.Tags(tt.args.c, tt.args.aids, tt.args.now)
			So(gotTm, ShouldEqual, tt.wantTm)
			So(err, ShouldEqual, tt.wantErr)
		})
	}
}
