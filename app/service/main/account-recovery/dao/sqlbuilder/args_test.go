// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sqlbuilder

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestArgs(t *testing.T) {
	cases := map[string][]interface{}{
		"abc ? def\n[123]":                   {"abc $? def", 123},
		"abc ? def\n[456]":                   {"abc $0 def", 456},
		"abc  def\n[]":                       {"abc $1 def", 123},
		"abc ? def\n[789]":                   {"abc ${s} def", Named("s", 789)},
		"abc  def \n[]":                      {"abc ${unknown} def ", 123},
		"abc $ def\n[]":                      {"abc $$ def", 123},
		"abcdef\n[]":                         {"abcdef$", 123},
		"abc ? ? ? ? def\n[123 456 123 456]": {"abc $? $? $0 $? def", 123, 456, 789},
		"abc ? raw ? raw def\n[123 123]":     {"abc $? $? $0 $? def", 123, Raw("raw"), 789},
	}

	for expected, c := range cases {
		args := new(Args)

		for i := 1; i < len(c); i++ {
			args.Add(c[i])
		}

		sql, values := args.Compile(c[0].(string))
		actual := fmt.Sprintf("%v\n%v", sql, values)

		if actual != expected {
			t.Fatalf("invalid compile result. [expected:%v] [actual:%v]", expected, actual)
		}
	}

	old := DefaultFlavor
	DefaultFlavor = PostgreSQL
	defer func() {
		DefaultFlavor = old
	}()

	// PostgreSQL flavor compiled sql.
	for expected, c := range cases {
		args := new(Args)

		for i := 1; i < len(c); i++ {
			args.Add(c[i])
		}

		sql, values := args.Compile(c[0].(string))
		actual := fmt.Sprintf("%v\n%v", sql, values)
		expected = toPostgreSQL(expected)

		if actual != expected {
			t.Fatalf("invalid compile result. [expected:%v] [actual:%v]", expected, actual)
		}
	}
}

func toPostgreSQL(sql string) string {
	parts := strings.Split(sql, "?")
	buf := &bytes.Buffer{}
	buf.WriteString(parts[0])

	for i, p := range parts[1:] {
		fmt.Fprintf(buf, "$%v", i+1)
		buf.WriteString(p)
	}

	return buf.String()
}
