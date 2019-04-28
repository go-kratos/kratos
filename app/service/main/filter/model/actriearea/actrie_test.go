package actriearea

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	matcher := NewMatcher()
	matcher.Insert("你是毛泽东吗", 30, []int64{0}, 265)
	matcher.Insert("毛泽东啊", 30, []int64{0}, 265)
	//matcher.Insert("邀请码", 30, []int64{0}, 265)
	matcher.Build()
	fmt.Println(matcher.Filter("你是毛泽东啊", 0, 15))
}
