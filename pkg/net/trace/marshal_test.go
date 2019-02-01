package trace

import (
	"testing"
)

func TestMarshalSpanV1(t *testing.T) {
	report := &mockReport{}
	t1 := newTracer("service1", report, &Config{DisableSample: true})
	sp1 := t1.New("opt_test").(*span)
	sp1.SetLog(Log("hello", "test123"))
	sp1.SetTag(TagString("tag1", "hell"), TagBool("booltag", true), TagFloat64("float64tag", 3.14159))
	sp1.Finish(nil)
	_, err := marshalSpanV1(sp1)
	if err != nil {
		t.Error(err)
	}
}
