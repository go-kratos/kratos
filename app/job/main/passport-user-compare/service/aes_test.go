package service

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"encoding/hex"

	"go-common/app/job/main/passport-user-compare/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Encrypt(t *testing.T) {
	Convey("Encrypt  param ", t, func() {
		input := []byte("15253340367")
		res, err := Encrypt(input, []byte("bili_@F1C2^Y_enc"))
		ShouldBeNil(err)
		ShouldNotBeNil(res)
		fmt.Printf("res is (%+v)\n", res)
		fmt.Printf("hex is (%+v) \n", strings.ToUpper(hex.EncodeToString(res)))
	})

}

func TestService_Decrypt(t *testing.T) {
	Convey("Decrypt param ", t, func() {
		var (
			in  []byte
			err error
			res []byte
		)
		sns := &model.OriginAccountSns{
			SinaUID: 5208921734,
		}
		th := &model.UserThirdBind{
			OpenID: "5208921734",
		}
		if s, err1 := strconv.ParseInt(th.OpenID, 10, 64); err1 != nil {
			fmt.Printf("error happen %+v", err1)
		} else {
			fmt.Println(s == sns.SinaUID)
		}
		a := new(model.UserBase)
		fmt.Printf("(%+v)", a == nil)
		if in, err = hex.DecodeString("3DEF2C03D0C822C57C5E3A931C087F27"); err != nil {
			fmt.Printf("err , (%+v)\n", err)
		} else {
			fmt.Printf("in, (%+v)\n", in)
		}
		res, err = Decrypt(in, []byte("bili_@F1C2^Y_enc"))
		ShouldBeNil(err)
		ShouldNotBeNil(res)
		fmt.Printf("(%+v) ,len is %d\n", string(res), len(in))

	})

}
