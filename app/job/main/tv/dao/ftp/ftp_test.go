package ftp

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {

			return false
		}
	}
	return true
}

func createFile(path string) {
	if !fileExist(path) {
		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}
		_, err = f.Write([]byte("Hello"))
		if err != nil {
			fmt.Println(err)
		}
		f.Close()
	}
}

func TestFtpRetry(t *testing.T) {
	var (
		callback func() error
		retry    = int(0)
		sleep    time.Duration
	)
	convey.Convey("Retry", t, func(ctx convey.C) {
		err := Retry(callback, retry, sleep)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFtpFileMd5(t *testing.T) {
	var (
		path    = "/tmp/testMd5.source"
		md5Path = "/tmp/testMd5.target"
	)
	convey.Convey("FileMd5", t, func(ctx convey.C) {
		createFile(path)
		err := d.FileMd5(path, md5Path)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFtpUploadFile(t *testing.T) {
	var (
		cfg  = d.conf.Search
		path = "/tmp/testMd5.source"
	)
	convey.Convey("UploadFile", t, func(ctx convey.C) {
		createFile(path)
		err := d.UploadFile(path, "testMd5.remote", cfg.FTP.RemotePgcURL)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
