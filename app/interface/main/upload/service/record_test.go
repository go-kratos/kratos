package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"go-common/app/interface/main/upload/conf"
	"go-common/app/interface/main/upload/model"

	. "github.com/smartystreets/goconvey/convey"
)

func loadbs() []byte {
	client := &http.Client{}
	resp, err := client.Get("https://i0.hdslb.com/bfs/album/b11defd6410e9fa5b6962c3c5f0402be2608db8c.jpg")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return bs
}

func TestGenImageUpload(t *testing.T) {
	Convey("create image and upload it", t, func() {
		res, err := svr.GenImageUpload(context.TODO(), "b4cfeeadca80f6f5", "c605dd5324f91ea1", "hello world", 2, true)
		So(err, ShouldBeNil)
		t.Logf("result:%+v", res)
	})
}

func TestUpload(t *testing.T) {
	bs := loadbs()
	Convey("test Upload", t, func() {
		now := time.Now().Unix()
		sha1 := sha1.New()
		sha1.Write([]byte(fmt.Sprintf("i love bilibili %s:%d", conf.Conf.Auths[0].AppSercet, now)))
		token := fmt.Sprintf("%s:%d", hex.EncodeToString(sha1.Sum([]byte(""))), now)
		result, err := svr.Upload(context.Background(), conf.Conf.Auths[0].AppKey, token, "", bs)
		if err != nil {
			t.Fatal(err.Error())
		}
		So("54aeb138b7fea2fe812aa8548f96cf1c0e4596ff", ShouldResemble, result.Etag)
	})
}

func TestUploadRecord(t *testing.T) {
	bs := loadbs()
	Convey("test UploadRecord", t, func() {
		ap := &model.UploadParam{
			ContentType: "image/jpeg",
			Bucket:      "static",
			FileName:    "",
			Dir:         "",
		}
		result, err := svr.UploadRecord(context.Background(), model.UploadInternal, 11, ap, bs)
		if err != nil {
			t.Fatal(err.Error())
		}
		So("54aeb138b7fea2fe812aa8548f96cf1c0e4596ff", ShouldResemble, result.Etag)
	})
}
