package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidName(t *testing.T) {
	Convey("ValidName", t, func() {
		So(ValidName("FasdA0asd"+string(1)+"asdas"), ShouldBeFalse)
		So(ValidName("FasdA0asd,asdas"), ShouldBeFalse)
		So(ValidName("Fasd啊三0a😀asdas"), ShouldBeFalse)
		So(ValidName("Fasd啊asdas"), ShouldBeFalse)
		So(ValidName("Fasd啊\asdas"), ShouldBeFalse)
		So(ValidName("Fasd啊-asdas_"), ShouldBeTrue)
		So(ValidName("Fasd啊\xF0\x9Fasdas_"), ShouldBeFalse)
		So(ValidName("Fasd啊xC2\xA0Fasdas_"), ShouldBeFalse)
	})
}
