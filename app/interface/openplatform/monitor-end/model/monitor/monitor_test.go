package monitor

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestCase .
type TestCase struct {
	tag      string
	testData *Log
	expected int
}

// TestLogData .
func TestLogData(t *testing.T) {
	var (
		tcs = []TestCase{
			TestCase{
				tag:      "empty log",
				testData: &Log{},
				expected: 1,
			},
			TestCase{
				tag:      "normal data",
				testData: &Log{Product: "test", LogType: "1", Event: "test", SubEvent: "test"},
				expected: 0,
			},
		}
	)
	for _, tc := range tcs {
		Convey(tc.tag, t, func() {
			_, _, _, err := tc.testData.LogData()
			if tc.expected == 0 {
				So(err, ShouldBeNil)
			} else {
				So(err, ShouldNotBeNil)
			}
		})
	}
}
