package service

import (
	"testing"

	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/passport-game-data/model"
)

func TestParseDiffLog(t *testing.T) {
	Convey("parse log text", t, func() {
		str := `local({\"mid\":80793085,\"userid\":\"adeqdiffer\",\"uname\":\"adeqdiffer\",\"pwd\":\"33a5fd6290550b88cc229275e9f790f7\",\"salt\":\"SFhrkmK3\",\"email\":\"\",\"tel\":\"\",\"country_id\":1,\"mobile_verified\":0,\"isleak\":0,\"modify_time\":\"2018-01-21T21:36:50+08:00\"}) local_encrypted({\"mid\":80793085,\"userid\":\"adeqdiffer\",\"uname\":\"adeqdiffer\",\"pwd\":\"7f0aa1b3dadda0c483aa78c3f3b048cf\",\"salt\":\"SFhrkmK3\",\"email\":\"\",\"tel\":\"\",\"country_id\":1,\"mobile_verified\":0,\"isleak\":0,\"ctime\":\"0001-01-01T00:00:00Z\",\"mtime\":\"2018-01-21T21:36:50+08:00\"}) cloud({\"mid\":80793085,\"userid\":\"adeqdiffer\",\"uname\":\"adeqdiffer\",\"pwd\":\"7f0aa1b3dadda0c483aa78c3f3b048cf\",\"salt\":\"SFhrkmK3\",\"email\":\"\",\"tel\":\"\",\"country_id\":1,\"mobile_verified\":0,\"isleak\":0,\"ctime\":\"2017-11-16T19:47:16+08:00\",\"mtime\":\"2018-01-21T21:36:50+08:00\"}`

		str = `local({\"mid\":83768597,\"userid\":\"difficenemy\",\"uname\":\"difficenemy\",\"pwd\":\"7e9f9a98269eb6fcc717f2d6e3a25fc2\",\"salt\":\"8pscksH6\",\"email\":\"\",\"tel\":\"\",\"country_id\":1,\"mobile_verified\":0,\"isleak\":0,\"modify_time\":\"2018-01-10T18:18:05+08:00\"}) local_encrypted({\"mid\":83768597,\"userid\":\"difficenemy\",\"uname\":\"difficenemy\",\"pwd\":\"08fa599f4497e876f7b4c7861f748361\",\"salt\":\"8pscksH6\",\"email\":\"\",\"tel\":\"\",\"country_id\":1,\"mobile_verified\":0,\"isleak\":0,\"ctime\":\"0001-01-01T00:00:00Z\",\"mtime\":\"2018-01-10T18:18:05+08:00\"}) cloud({\"mid\":83768597,\"userid\":\"difficenemy\",\"uname\":\"difficenemy\",\"pwd\":\"08fa599f4497e876f7b4c7861f748361\",\"salt\":\"8pscksH6\",\"email\":\"\",\"tel\":\"\",\"country_id\":1,\"mobile_verified\":0,\"isleak\":0,\"ctime\":\"2017-11-16T23:59:21+08:00\",\"mtime\":\"2018-01-10T18:18:05+08:00\"})`

		res := replace(str)

		t.Logf("res: %s", res)
		cRes := new(model.CompareRes)

		err := json.Unmarshal([]byte(res), &cRes)
		So(err, ShouldBeNil)

		rStr, _ := json.Marshal(cRes)
		t.Logf("res: %s, cRes: %s", res, rStr)
	})
}
