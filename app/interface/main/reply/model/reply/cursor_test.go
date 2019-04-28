package reply

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	cases := []struct {
		inputIds []int64
		comp     Comp
		expected []int64
	}{
		{
			inputIds: []int64{6, 5, 3, 2, 1},
			comp:     OrderASC,
			expected: []int64{1, 2, 3, 5, 6},
		},
		{
			inputIds: []int64{5, 3, 6, 1, 2},
			comp:     OrderDESC,
			expected: []int64{6, 5, 3, 2, 1},
		},
		{
			inputIds: []int64{2, 1},
			comp:     nil,
			expected: []int64{2, 1},
		},
		{
			inputIds: []int64{},
			comp:     OrderDESC,
			expected: []int64{},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			got := SortArr(c.inputIds, c.comp)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("err sort, want %v, got %v", c.expected, got)
			}
		})
	}
}
