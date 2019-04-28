package grok

import "testing"

func TestReverseList(t *testing.T) {
	var array = []string{"A", "B", "C", "D"}
	var expectedArray = []string{"D", "C", "B", "A"}
	arrayReversed := reverseList(array)
	if !sliceEquals(arrayReversed, expectedArray) {
		t.Fatalf("reversedList is %+v, expected : %+v", arrayReversed, expectedArray)
	}
}

func TestSortGraph(t *testing.T) {
	var g = graph{}
	g["7"] = []string{"11", "8"}
	g["5"] = []string{"11"}
	g["3"] = []string{"8", "10"}
	g["11"] = []string{"2", "9", "10"}
	g["8"] = []string{"9", "10"}
	g["2"] = []string{}
	g["9"] = []string{}
	g["10"] = []string{}

	validOrders := [][]string{
		{"3", "5", "7", "8", "11", "10", "9", "2"},
		{"3", "5", "7", "8", "11", "2", "10", "9"},
		{"3", "5", "7", "8", "11", "9", "2", "10"},
		{"7", "3", "8", "5", "11", "10", "9", "2"},
		{"3", "5", "7", "11", "2", "8", "10", "9"},
		{"5", "7", "11", "2", "3", "8", "10", "9"},
	}

	order, cycle := sortGraph(g)
	if cycle != nil {
		t.Fatal("cycle detected while not expected")
	}
	for _, expectedOrder := range validOrders {
		if sliceEquals(order, expectedOrder) {
			return
		}
	}

	t.Fatalf("sorted graph is %+v, expected a order like: %+v", order, validOrders[0])
}

func TestSortGraphWithCycle(t *testing.T) {
	var g = graph{}
	g["7"] = []string{"11", "8"}
	g["5"] = []string{"11"}
	g["3"] = []string{"8", "10"}
	g["11"] = []string{"2", "9", "10"}
	g["8"] = []string{"9", "10"}
	g["2"] = []string{}
	g["9"] = []string{"3"}
	g["10"] = []string{}

	validCycles := [][]string{
		{"3", "9", "8"},
		{"8", "3", "9"},
		{"9", "8", "3"},
	}

	_, cycle := sortGraph(g)

	if cycle == nil {
		t.Fatal("cycle not detected while sorting graph")
	}

	for _, expectedCycle := range validCycles {
		if sliceEquals(cycle, expectedCycle) {
			return
		}
	}

	t.Fatalf("cycle have %+v, expected %+v", cycle, validCycles[0])
}

func sliceEquals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
