package service

import (
	"fmt"
	"testing"

	"go-common/app/job/main/passport-encrypt/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_EncryptAccount(t *testing.T) {
	Convey("Encrypt a right account", t, func() {
		account := &model.OriginAccount{
			Mid:            12047569,
			UserID:         "bili_1710676855",
			Uname:          "Bili_12047569",
			Pwd:            "3686c9d96ae6896fe117319ba6c07087",
			Salt:           "pdMXF856",
			Email:          "",
			Tel:            "13122110209",
			CountryID:      1,
			MobileVerified: 1,
			Isleak:         0,
		}
		enAcc := EncryptAccount(account)
		fmt.Printf("right (+%v)", enAcc)
	})

	Convey("Encrypt a empty tel account", t, func() {
		account := &model.OriginAccount{
			Mid:            12047569,
			UserID:         "bili_1710676855",
			Uname:          "Bili_12047569",
			Pwd:            "3686c9d96ae6896fe117319ba6c07087",
			Salt:           "pdMXF856",
			Email:          "",
			Tel:            "",
			CountryID:      1,
			MobileVerified: 1,
			Isleak:         0,
		}
		enAcc := EncryptAccount(account)
		fmt.Printf("empty tel (+%v)", enAcc)
	})
}

func TestService_DecryptAccount(t *testing.T) {
	Convey("Decrypt right account", t, func() {
		account := &model.EncryptAccount{
			Mid:            12047569,
			UserID:         "bili_1710676855",
			Uname:          "Bili_12047569",
			Pwd:            "3686c9d96ae6896fe117319ba6c07087",
			Salt:           "pdMXF856",
			Email:          "",
			Tel:            []byte{216, 28, 241, 48, 17, 243, 198, 234, 205, 33, 216, 25, 75, 146, 97, 203},
			CountryID:      1,
			MobileVerified: 1,
			Isleak:         0,
		}
		enAcc := DecryptAccount(account)
		fmt.Printf("(+%v)", enAcc)
	})

	Convey("Decrypt empty tel  account", t, func() {
		account := &model.EncryptAccount{
			Mid:            12047569,
			UserID:         "bili_1710676855",
			Uname:          "Bili_12047569",
			Pwd:            "3686c9d96ae6896fe117319ba6c07087",
			Salt:           "pdMXF856",
			Email:          "",
			Tel:            []byte{73, 224, 88, 238, 85, 242, 37, 46, 200, 70, 5, 49, 73, 107, 179, 51},
			CountryID:      1,
			MobileVerified: 1,
			Isleak:         0,
		}
		enAcc := DecryptAccount(account)
		fmt.Printf("(+%v)", enAcc)
	})
}
