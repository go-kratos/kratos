package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestCreateTag
func TestDao_CreateTag(t *testing.T) {
	Convey("CreateTag", t, func() {
		once.Do(startService)
		tx := d.db.Begin()
		err := d.CreateTag(context.TODO(), tx, 1, "10003")
		tx.Commit()
		So(err, ShouldBeNil)
	})
}
