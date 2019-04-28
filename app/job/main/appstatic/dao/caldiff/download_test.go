package caldiff

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_DownloadFile(t *testing.T) {
	var (
		newPath = fmt.Sprintf("%s/%s", d.c.Cfg.Diff.Folder, "123.zip")
		url     = "http://i0.hdslb.com/bfs/face/ca14680bcfe4a956d1b6c06fbc1f6a6529257746.jpg"
	)
	Convey("TestDao_DownloadFile", t, WithDao(func(d *Dao) {
		httpMock("GET", url).Reply(200).Body(bytes.NewReader([]byte("test")))
		data, err := d.DownloadFile(ctx, url, newPath)
		So(err, ShouldBeNil)
		So(data, ShouldBeGreaterThan, 0)
		fmt.Println(data)
	}))
	Convey("CreateFile Error", t, WithDao(func(d *Dao) {
		httpMock("GET", url).Reply(200).Body(bytes.NewReader([]byte("test")))
		_, err := d.DownloadFile(ctx, url, "/test/test.txt")
		So(err, ShouldNotBeNil)
		fmt.Println(err)
	}))
}
