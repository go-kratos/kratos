package service

import (
	"encoding/hex"
	"fmt"
	"testing"

	"go-common/app/interface/main/passport-login/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Encrypt(t *testing.T) {
	once.Do(startService)
	Convey("Encrypt  param ", t, func() {
		input := []byte("nianshoud@126.com")
		res, err := Encrypt(input, []byte(aesKey))
		ShouldBeNil(err)
		ShouldNotBeNil(res)
		fmt.Printf("hex is (%+v) \n", hex.EncodeToString(res))
	})
}

func TestService_Decrypt(t *testing.T) {
	once.Do(startService)
	Convey("Decrypt param ", t, func() {
		var (
			in  []byte
			err error
			res []byte
		)
		a := new(model.User)
		fmt.Printf("(%+v)", a == nil)
		if in, err = hex.DecodeString("A77CF15FF3B95CCD3F34B879532D44D5"); err != nil {
			fmt.Printf("err , (%+v)\n", err)
		} else {
			fmt.Printf("in, (%+v)\n", in)
		}
		res, err = Decrypt(in, []byte(aesKey))
		ShouldBeNil(err)
		ShouldNotBeNil(res)
		fmt.Printf("(%+v) ,len is %d\n", string(res), len(in))
	})
}
