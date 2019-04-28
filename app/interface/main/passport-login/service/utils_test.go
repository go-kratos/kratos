package service

import (
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_doEncrypt(t *testing.T) {
	once.Do(startService)
	Convey("doEncrypt", t, func() {
		res := s.doEncrypt("13122119298")
		ShouldNotBeNil(res)
	})
}

func TestService_doDecrypt(t *testing.T) {
	once.Do(startService)
	Convey("doEncrypt", t, func() {
		if in, err := hex.DecodeString("A77CF15FF3B95CCD3F34B879532D44D5"); err != nil {
			fmt.Printf("err , (%+v)\n", err)
		} else {
			if res, err := s.doDecrypt(in); err != nil {
				ShouldBeNil(err)
			} else {
				ShouldNotBeNil(res)
			}
		}
	})
}

func TestService_doHash(t *testing.T) {
	once.Do(startService)
	Convey("doHash", t, func() {
		res := s.doHash("13122119298")
		ShouldNotBeNil(res)
	})
}

func TestService_addBase64Padding(t *testing.T) {
	once.Do(startService)
	Convey("addBase64Padding", t, func() {
		s := addBase64Padding("a")
		fmt.Println(s)
	})
}
