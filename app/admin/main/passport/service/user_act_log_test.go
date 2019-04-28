package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/passport/conf"
	"go-common/app/admin/main/passport/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_DecryptBindLog(t *testing.T) {
	config := &conf.Config{
		Encode: &conf.Encode{
			AesKey: "0123456789abcdef",
			Salt:   "",
		},
	}
	s := New(config)
	convey.Convey("", t, func() {
		res, err := s.DecryptBindLog(context.Background(), &model.DecryptBindLogParam{EncryptText: []string{"IsjRu7dmHBY8l7bGf6O3rgDegFvh3cVTgWkf2Bn87Oc="}})
		convey.So(err, convey.ShouldBeNil)
		convey.So(res["IsjRu7dmHBY8l7bGf6O3rgDegFvh3cVTgWkf2Bn87Oc="], convey.ShouldEqual, "19921218988")
	})
}
