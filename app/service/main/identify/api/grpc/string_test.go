package v1

import (
	"fmt"
	"testing"
)

func TestGetCookieInfoReply_String(t *testing.T) {
	m := GetCookieInfoReq{
		Cookie: "==000123000==",
	}
	fmt.Println(fmt.Sprintf("%v", m))
}
