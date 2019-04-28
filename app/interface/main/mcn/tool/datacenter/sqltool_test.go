package datacenter

import "testing"

func TestSelect(t *testing.T) {
	var q = &Query{}
	q.Select("a,b   b2,c as  c2")
	t.Logf("query=%s", q)
	q = &Query{}
	q.Select("   * ")
	t.Logf("query=%s", q)
}

func TestWhere(t *testing.T) {
	var q = &Query{}
	q.Select("a,b,c as yeah").Where(
		ConditionMapType{
			"field1": ConditionIn(1, 2, 3, 4),
			"field3": ConditionIn("OK"),
		},
		ConditionMapType{
			"field2": ConditionGt(100),
		},
		ConditionMapType{
			"field1": ConditionGte(100),
		})
	t.Logf("query=%s", q)
}

func TestSort(t *testing.T) {
	var q = &Query{}
	q.Order("field1 desc, field2")
	t.Logf("query=%s", q)
}
