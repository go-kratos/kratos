package service

import (
	"strings"
	"testing"
)

func TestSliceUnique(t *testing.T) {
	var slice1 = []string{"绘画", "哈哈", "22", "2223", "绘画"}
	t.Logf("unique : %v", SliceUnique(Slice2Interface(slice1)))
	tag := "绘画, 哈哈, 22, 2223, 绘画"
	t.Logf("unique : %v", strings.Join(Slice2String(SliceUnique(Slice2Interface(strings.Split(tag, ",")))), ","))

}
