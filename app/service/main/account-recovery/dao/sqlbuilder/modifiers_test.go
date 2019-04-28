// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestEscape(t *testing.T) {
	cases := map[string]string{
		"foo":  "foo",
		"$foo": "$$foo",
		"$$$":  "$$$$$$",
	}
	var inputs, expects []string

	for s, expected := range cases {
		inputs = append(inputs, s)
		expects = append(expects, expected)

		if actual := Escape(s); actual != expected {
			t.Fatalf("invalid escape result. [expected:%v] [actual:%v]", expected, actual)
		}
	}

	actuals := EscapeAll(inputs...)

	if !reflect.DeepEqual(expects, actuals) {
		t.Fatalf("invalid escape result. [expected:%v] [actual:%v]", expects, actuals)
	}
}

func TestFlatten(t *testing.T) {
	cases := [][2]interface{}{
		{
			"foo",
			[]interface{}{"foo"},
		},
		{
			[]int{1, 2, 3},
			[]interface{}{1, 2, 3},
		},
		{
			[]interface{}{"abc", []int{1, 2, 3}, [3]string{"def", "ghi"}},
			[]interface{}{"abc", 1, 2, 3, "def", "ghi", ""},
		},
	}

	for _, c := range cases {
		input, expected := c[0], c[1]
		actual := Flatten(input)

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("invalid flatten result. [expected:%v] [actual:%v]", expected, actual)
		}
	}
}
