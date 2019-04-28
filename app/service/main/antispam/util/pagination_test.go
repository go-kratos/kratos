package util

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSimplePage(t *testing.T) {
	cases := []struct {
		perPage      int64
		curPage      int64
		expectedFrom int64
		expectedTo   int64
	}{
		{20, 9, 161, 180},
		{0, 1, 1, 20},
		{0, 0, 1, 20},
		{1, 0, 1, 20},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("curPage(%d) perPage(%d)", c.curPage, c.perPage), func(t *testing.T) {
			p := &Pagination{
				CurPage: c.curPage,
				PerPage: c.perPage,
			}
			from, to := p.SimplePage()
			if from != c.expectedFrom || to != c.expectedTo {
				t.Errorf("cond.SimplePage() = from: %d, to: %d, want: %d, %d", from, to, c.expectedFrom, c.expectedTo)
			}
		})
	}
}

func TestPage(t *testing.T) {
	cases := []struct {
		total        int64
		perPage      int64
		curPage      int64
		expectedFrom int64
		expectedTo   int64
	}{
		{66269, 20, 3314, 66261, 66269},
		{66269, 20, 3313, 66241, 66260},
		{81, 20, 9, 0, 0},
		{100, 0, 1, 1, 20},
		{1, 20, 1, 1, 1},
		{0, 20, 1, 0, 0},
		{5, 20, 1, 1, 5},
		{211, 20, 3, 41, 60},
		{100, 100, 1, 1, 100},
		{101, 20, 6, 101, 101},
		{211, 20, 2, 21, 40},
		{211, 20, 1, 1, 20},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("total(%d) curPage(%d) perPage(%d)", c.total, c.curPage, c.perPage), func(t *testing.T) {
			p := &Pagination{
				CurPage: c.curPage,
				PerPage: c.perPage,
			}
			from, to := p.Page(c.total)
			if from != c.expectedFrom || to != c.expectedTo {
				t.Errorf("cond.Page(%d) = from: %d, to: %d, want: %d, %d", c.total, from, to, c.expectedFrom, c.expectedTo)
			}
		})
	}
}

func TestOffsetLimit(t *testing.T) {
	cases := []struct {
		total          int64
		perpage        int64
		curpage        int64
		expectedoffset int64
		expectedlimit  int64
	}{
		{66269, 20, 3314, 66260, 9},
		{66269, 20, 3313, 66240, 20},
		{100, 0, 1, 0, 20},
		{1, 20, 1, 0, 1},
		{0, 20, 1, 0, 0},
		{5, 20, 1, 0, 5},
		{211, 20, 3, 40, 20},
		{100, 100, 1, 0, 100},
		{101, 20, 6, 100, 1},
		{211, 20, 2, 20, 20},
		{211, 20, 1, 0, 20},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("total(%d) curpage(%d) perpage(%d)", c.total, c.curpage, c.perpage), func(t *testing.T) {
			p := &Pagination{
				CurPage: c.curpage,
				PerPage: c.perpage,
			}
			offset, limit := p.OffsetLimit(c.total)
			if offset != c.expectedoffset || limit != c.expectedlimit {
				t.Errorf("cond.offsetlimit(%d) = offset: %d, limit: %d, want %d, %d", c.total, offset, limit, c.expectedoffset, c.expectedlimit)
			}
		})
	}
}

func TestBulkPage(t *testing.T) {
	p := &Pagination{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 9999; i++ {
		p.CurPage = rand.Int63n(10000)
		p.PerPage = rand.Int63n(300)
		total := rand.Int63n(20000)
		from, to := p.Page(total)
		if from < 0 || to < 0 {
			t.Fatalf(`Bulk test page fail, got negative result, 
				total: %d, curPage: %d, perPage: %d, from: %d, to: %d`, total, p.CurPage, p.PerPage, from, to)
		}
	}
}
