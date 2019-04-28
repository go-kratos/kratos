// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sqlbuilder

import (
	"fmt"
	"testing"
)

func TestFlavor(t *testing.T) {
	cases := map[Flavor]string{
		0:          "<invalid>",
		MySQL:      "MySQL",
		PostgreSQL: "PostgreSQL",
	}

	for f, expected := range cases {
		if actual := f.String(); actual != expected {
			t.Fatalf("invalid flavor name. [expected:%v] [actual:%v]", expected, actual)
		}
	}
}

func ExampleFlavor() {
	// Create a flavored builder.
	sb := PostgreSQL.NewSelectBuilder()
	sb.Select("name").From("user").Where(
		sb.E("id", 1234),
		sb.G("rank", 3),
	)
	sql, args := sb.Build()

	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// SELECT name FROM user WHERE id = $1 AND rank > $2
	// [1234 3]
}
