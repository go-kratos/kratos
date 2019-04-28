package service

import (
	"testing"

	"go-common/app/job/main/up/model/archivemodel"
)

func Test_archiveStateChange(t *testing.T) {
	var (
		testcase = [][]*archivemodel.ArchiveCanal{
			{{State: 0}, {State: -1}},
			{{State: -5}, {State: 0}},
			{{State: 10}, {State: 10}},
			{{State: -5}, {State: -5}},
		}
		testresult = []bool{
			true,
			true,
			false,
			false,
		}
	)

	for i := range testcase {
		var cas = testcase[i]
		if archiveStateChange(cas[0], cas[1]) != testresult[i] {
			t.Errorf("test fail, testcase[%d]=%v, expect=%t", i, cas, testresult[i])
			t.Fail()
		}
	}

}
