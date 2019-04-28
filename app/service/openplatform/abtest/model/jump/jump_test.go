package jump

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type hashStruct struct {
	value uint64
	exp   int32
}

func Test_Hash(t *testing.T) {
	Convey("Test_Hash: ", t, func() {
		h := hashStruct{
			value: uint64(314978625),
			exp:   int32(18),
		}
		bucket := 100
		v := Hash(h.value, bucket)
		So(v, ShouldEqual, h.exp)
	})
}

type md5Struct struct {
	value string
	exp   uint64
}

func Test_Md5(t *testing.T) {
	Convey("Test_Hash: ", t, func() {
		h := md5Struct{
			value: "987654321",
			exp:   uint64(7979946199622949865),
		}
		v := Md5(h.value)
		So(v, ShouldEqual, h.exp)
	})
}
