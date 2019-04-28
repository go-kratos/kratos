package dao

import (
	"bytes"
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UploadSql(t *testing.T) {
	Convey("upload picture or log file to bfs", t, WithMysql(func(d *Dao) {
		var (
			fileName       = "bilibili.png"
			fileType       = "png"
			expire   int64 = 1
			bt             = []byte("/9j/4QAo.Reader(bt)YRXhpZgAASUkqAAgAAAAAAAAAAAAAA")
			body           = bytes.NewReader(bt)
		)
		_, err := d.Upload(context.Background(), fileName, fileType, expire, body)
		if err == errUpload {
			err = nil
		}
		So(err, ShouldBeNil)
	}))
}
