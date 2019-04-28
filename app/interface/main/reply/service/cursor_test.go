package service

import (
	"context"
	"reflect"
	"testing"

	model "go-common/app/interface/main/reply/model/reply"
)

func TestRemove(t *testing.T) {
	cases := []struct {
		inputIds []int64
		id       int64
		expected []int64
	}{
		{
			inputIds: []int64{5, 3, 2, 1},
			id:       6,
			expected: []int64{5, 3, 2, 1},
		},
		{
			inputIds: []int64{6, 5, 3, 2, 1},
			id:       6,
			expected: []int64{5, 3, 2, 1},
		},
		{
			inputIds: []int64{5, 3, 1, 6, 1, 2},
			id:       1,
			expected: []int64{5, 3, 6, 2},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			got := Remove(c.inputIds, c.id)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("err sort, want %v, got %v", c.expected, got)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	cases := []struct {
		inputIds []int64
		expected []int64
	}{
		{
			inputIds: []int64{1, 2, 1, 2, 3, 43},
			expected: []int64{1, 2, 3, 43},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			got := model.SortArr(Unique(c.inputIds), model.OrderASC)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("err sort, want %v, got %v", c.expected, got)
			}
		})
	}

}

func TestNewCursorByReplyID(t *testing.T) {
	s := &Service{}
	s.NewCursorByReplyID(context.Background(), 1123, int8(1), 112, 20, model.OrderDESC)
}

func TestGetRootReplyListByCursor(t *testing.T) {
	s := &Service{}
	s.GetRootReplyListByCursor(context.Background(), &model.CursorParams{})
}

func TestInsertInto(t *testing.T) {
	cases := []struct {
		inputIds []int64
		id       int64
		size     int
		comp     model.Comp
		expected []int64
	}{
		{
			inputIds: []int64{},
			id:       100,
			size:     5,
			expected: []int64{100},
			comp:     model.OrderDESC,
		},
		{
			inputIds: []int64{1, 2, 3, 5},
			size:     5,
			id:       4,
			expected: []int64{1, 2, 3, 4, 5},
			comp:     model.OrderASC,
		},
		{
			inputIds: []int64{2, 5, 1, 3},
			size:     5,
			id:       4,
			expected: []int64{5, 4, 3, 2, 1},
			comp:     model.OrderDESC,
		},
		{
			inputIds: []int64{1, 2, 3, 5},
			size:     5,
			id:       4,
			expected: []int64{5, 4, 3, 2, 1},
			comp:     model.OrderDESC,
		},
		{
			inputIds: []int64{1, 2, 3, 4, 5},
			size:     5,
			id:       4,
			expected: []int64{1, 2, 3, 4, 5},
			comp:     model.OrderDESC,
		},
		{
			inputIds: []int64{1, 2, 3, 4, 5},
			size:     5,
			id:       4,
			expected: []int64{1, 2, 3, 4, 5},
			comp:     model.OrderASC,
		},
		{
			inputIds: []int64{5, 4, 3, 1},
			id:       2,
			size:     5,
			expected: []int64{5, 4, 3, 2, 1},
			comp:     model.OrderDESC,
		},
		{
			inputIds: []int64{5},
			id:       4,
			size:     2,
			expected: []int64{5, 4},
			comp:     model.OrderDESC,
		},
		{
			inputIds: []int64{5},
			id:       4,
			size:     2,
			expected: []int64{4, 5},
			comp:     model.OrderASC,
		},
		{
			inputIds: []int64{5},
			id:       4,
			size:     1,
			expected: []int64{5},
			comp:     model.OrderDESC,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			got := InsertInto(c.inputIds, c.id, c.size, c.comp)
			if !reflect.DeepEqual(c.expected, got) {
				t.Errorf("err sorted insert, want %v, got %v", c.expected, got)
			}
		})
	}
}
