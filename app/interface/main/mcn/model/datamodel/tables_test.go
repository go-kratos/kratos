package datamodel

import (
	"testing"
	"time"
)

type timeTestCase struct {
	In  []byte
	Out time.Time
}

func timeHelper(tm time.Time, err error) time.Time {
	return tm
}
func TestLogTime_UnmarshalJSON(t *testing.T) {
	var (
		testcase = []timeTestCase{
			{[]byte(`"2018-11-11"`), timeHelper(time.ParseInLocation("2006-01-02", "2018-11-11", time.Local))},
			{[]byte(`1542795906`), time.Unix(1542795906, 0)},
		}
	)

	for _, testcase := range testcase {
		var ltm LogTime
		var err = ltm.UnmarshalJSON(testcase.In)
		if err != nil {
			t.Errorf("err=%v", err)
			t.Fail()
			continue
		}

		if int64(ltm.Time()) != testcase.Out.Unix() {
			t.Errorf("expect=%d, get=%d", testcase.Out, ltm.Time())
			t.Fail()
			continue
		}
	}
}
