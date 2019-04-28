package service

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Encrypt(t *testing.T) {
	Convey("Encrypt  param ", t, func() {
		input := []byte("13122119298")
		res, err := Encrypt(input)
		ShouldBeNil(err)
		ShouldNotBeNil(res)
		fmt.Printf("(%+v) \n", res)
	})

}

func TestService_Decrypt(t *testing.T) {
	Convey("Decrypt param ", t, func() {
		input := []byte{70, 251, 31, 100, 140, 98, 156, 101, 27, 21, 201, 100, 159, 107, 90, 147}
		fmt.Printf("input is  (%+v) \n", input)
		res, err := Decrypt(input)
		ShouldBeNil(err)
		ShouldNotBeNil(res)
		fmt.Printf("(%+v) \n", string(res))
	})
}
