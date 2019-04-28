package model

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSON(t *testing.T) {
	Convey("json", t, func() {
		level := &AreaLevel{
			Level: 20,
			Area: map[string]int8{
				"reply": 30,
				"live":  10,
			},
		}
		bytes, err := json.Marshal(level)
		So(err, ShouldBeNil)
		So(string(bytes), ShouldBeEmpty)
	})
}
