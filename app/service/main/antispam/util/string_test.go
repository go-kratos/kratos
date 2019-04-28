package util

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleIntSliToSQLVarchars() {
	fmt.Println(IntSliToSQLVarchars([]int64{1, 2, 3}))
	// Output: 1,2,3
}

func ExampleStrToIntSli() {
	fmt.Println(StrToIntSli("1,2,3", ","))
	// Output: [1 2 3] <nil>
}

func ExampleStrSliToSQLVarchars() {
	fmt.Println(StrSliToSQLVarchars([]string{"default", "deleted", "modified"}))
	// Output: 'default','deleted','modified'
}

func TestStrSliToSQLVarchars(t *testing.T) {
	cases := []struct {
		s        []string
		expected string
	}{
		{[]string{"foo", "bar"}, "'foo','bar'"},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("inputStr(%v)", c.s), func(t *testing.T) {
			got := StrSliToSQLVarchars(c.s)
			if got != c.expected {
				t.Errorf("StrSliToSQLVarchars(%v) = %s, want: %s", c.s, got, c.expected)
			}
		})
	}
}

func TestStrToIntSli(t *testing.T) {
	cases := []struct {
		s           string
		delimiter   string
		expectedSli []int64
		expectedErr error
	}{
		{"1,2,3", ",", []int64{1, 2, 3}, nil},
		{"1 2 3", " ", []int64{1, 2, 3}, nil},
		{"1|2|3", "|", []int64{1, 2, 3}, nil},
	}
	for _, c := range cases {
		assert := assert.New(t)
		t.Run(fmt.Sprintf("inputString(%v) delimiter(%s)", c.s, c.delimiter), func(t *testing.T) {
			got, err := StrToIntSli(c.s, c.delimiter)
			assert.Equal(c.expectedSli, got, "")
			assert.Equal(c.expectedErr, err, "")
		})
	}
}

func TestIntSliToStr(t *testing.T) {
	cases := []struct {
		s         []int64
		delimiter string
		expected  string
	}{
		{[]int64{1, 2, 3}, ",", "1,2,3"},
		{[]int64{1, 2, 3}, " ", "1 2 3"},
		{[]int64{1, 2, 3}, "|", "1|2|3"},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("inputSli(%v) delimiter(%s)", c.s, c.delimiter), func(t *testing.T) {
			got := intSliToStr(c.s, c.delimiter)
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("IntSliToStr(%v, %s) = %s, want %s", c.s, c.delimiter, got, c.expected)
			}
		})
	}
}

func TestStripPrefix(t *testing.T) {
	cases := []struct {
		name           string
		content        string
		expectedOutput string
	}{
		{
			"need strip prefix",
			"回复 @画鸾凰 :我知道 但我的不是正版的 上不了工坊 只能要模型软件 格式是LPK的模型软件 你找一下看看 模型列表下面应该有在哪个文件夹里面 找到可以发我QQ1918882322 如果找不到就算吧",
			"我知道 但我的不是正版的 上不了工坊 只能要模型软件 格式是LPK的模型软件 你找一下看看 模型列表下面应该有在哪个文件夹里面 找到可以发我QQ1918882322 如果找不到就算吧",
		},
		{
			"empty reply body",
			"回复 @画鸾凰 :",
			"",
		},
		{
			"not need strip",
			"我知道 但我的不是正版的 上不了工坊 只能要模型软件 格式是LPK的模型软件 你找一下看看 模型列表下面应该有在哪个文件夹里面 找到可以发我QQ1918882322 如果找不到就算吧",
			"我知道 但我的不是正版的 上不了工坊 只能要模型软件 格式是LPK的模型软件 你找一下看看 模型列表下面应该有在哪个文件夹里面 找到可以发我QQ1918882322 如果找不到就算吧",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := StripPrefix(c.content, "回复 @", ":")
			if actual != c.expectedOutput {
				t.Fatalf("Strip Prefix failed, expected %q \t\n got %q", c.expectedOutput, actual)
			}
		})
	}
}

func TestSameChar(t *testing.T) {
	cases := []struct {
		content        string
		expectedResult bool
	}{
		{"~~~~~~~", true},
		{"666666666", true},
		{"666666666~~~", false},
		{"WWWWWWW", true},
		{"XXXxxx", true},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("content(%s)", c.content), func(t *testing.T) {
			if rs := SameChar(c.content); rs != c.expectedResult {
				t.Errorf("SameChar(%s) = %v, want %v", c.content, rs, c.expectedResult)
			}
		})
	}
}
