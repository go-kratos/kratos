package dao

import (
	"context"
	"encoding/json"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"go-common/app/interface/main/player/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_Playurl(t *testing.T) {
	convey.Convey("test playurl", t, func(ctx convey.C) {
		// aid 10108138
		mid := int64(908085)
		//cid=47083420&qn=80&type=&otype=json&fnver=0&fnval=8
		arg := &model.PlayurlArg{
			Aid:   10111420,
			Cid:   10151472,
			Qn:    0,
			OType: "xml",
			//Player: 1,
			Fnver: 0,
			Fnval: 8,
		}
		token := ""
		playurl := "http://uat-videodispatch-ugc.bilibili.co/v3/playurl"
		defer gock.OffAll()
		httpMock("GET", playurl).Reply(200).JSON(`{"code":0}`)
		data, err := d.Playurl(context.Background(), mid, arg, playurl, token)
		convey.So(err, convey.ShouldBeNil)
		str, _ := json.Marshal(data)
		convey.Println(string(str))
	})
}
