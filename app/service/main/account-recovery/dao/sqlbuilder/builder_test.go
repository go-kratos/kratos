// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package sqlbuilder

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func ExampleBuildf() {
	sb := NewSelectBuilder()
	sb.Select("id").From("user")

	explain := Buildf("EXPLAIN %v LEFT JOIN SELECT * FROM banned WHERE state IN (%v, %v)", sb, 1, 2)
	sql, args := explain.Build()
	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// EXPLAIN SELECT id FROM user LEFT JOIN SELECT * FROM banned WHERE state IN (?, ?)
	// [1 2]
}

func ExampleBuild() {
	sb := NewSelectBuilder()
	sb.Select("id").From("user").Where(sb.In("status", 1, 2))

	b := Build("EXPLAIN $? LEFT JOIN SELECT * FROM $? WHERE created_at > $? AND state IN (${states}) AND modified_at BETWEEN $2 AND $?",
		sb, Raw("banned"), 1514458225, 1514544625, Named("states", List([]int{3, 4, 5})))
	sql, args := b.Build()

	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// EXPLAIN SELECT id FROM user WHERE status IN (?, ?) LEFT JOIN SELECT * FROM banned WHERE created_at > ? AND state IN (?, ?, ?) AND modified_at BETWEEN ? AND ?
	// [1 2 1514458225 3 4 5 1514458225 1514544625]
}

func ExampleBuildNamed() {
	b := BuildNamed("SELECT * FROM ${table} WHERE status IN (${status}) AND name LIKE ${name} AND created_at > ${time} AND modified_at < ${time} + 86400",
		map[string]interface{}{
			"time":   sql.Named("start", 1234567890),
			"status": List([]int{1, 2, 5}),
			"name":   "Huan%",
			"table":  Raw("user"),
		})
	sql, args := b.Build()

	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// SELECT * FROM user WHERE status IN (?, ?, ?) AND name LIKE ? AND created_at > @start AND modified_at < @start + 86400
	// [1 2 5 Huan% {{} start 1234567890}]
}

func ExampleWithFlavor() {
	sql, args := WithFlavor(Buildf("SELECT * FROM foo WHERE id = %v", 1234), PostgreSQL).Build()

	fmt.Println(sql)
	fmt.Println(args)

	// Explicitly use MySQL as the flavor.
	sql, args = WithFlavor(Buildf("SELECT * FROM foo WHERE id = %v", 1234), PostgreSQL).BuildWithFlavor(MySQL)

	fmt.Println(sql)
	fmt.Println(args)

	// Output:
	// SELECT * FROM foo WHERE id = $1
	// [1234]
	// SELECT * FROM foo WHERE id = ?
	// [1234]
}

func TestBuildWithPostgreSQL(t *testing.T) {
	sb1 := PostgreSQL.NewSelectBuilder()
	sb1.Select("col1", "col2").From("t1").Where(sb1.E("id", 1234), sb1.G("level", 2))

	sb2 := PostgreSQL.NewSelectBuilder()
	sb2.Select("col3", "col4").From("t2").Where(sb2.E("id", 4567), sb2.LE("level", 5))

	// Use DefaultFlavor (MySQL) instead of PostgreSQL.
	sql, args := Build("SELECT $1 AS col5 LEFT JOIN $0 LEFT JOIN $2", sb1, 7890, sb2).Build()

	if expected := "SELECT ? AS col5 LEFT JOIN SELECT col1, col2 FROM t1 WHERE id = ? AND level > ? LEFT JOIN SELECT col3, col4 FROM t2 WHERE id = ? AND level <= ?"; sql != expected {
		t.Fatalf("invalid sql. [expected:%v] [actual:%v]", expected, sql)
	}

	if expected := []interface{}{7890, 1234, 2, 4567, 5}; !reflect.DeepEqual(args, expected) {
		t.Fatalf("invalid args. [expected:%v] [actual:%v]", expected, args)
	}

	old := DefaultFlavor
	DefaultFlavor = PostgreSQL
	defer func() {
		DefaultFlavor = old
	}()

	sql, args = Build("SELECT $1 AS col5 LEFT JOIN $0 LEFT JOIN $2", sb1, 7890, sb2).Build()

	if expected := "SELECT $1 AS col5 LEFT JOIN SELECT col1, col2 FROM t1 WHERE id = $2 AND level > $3 LEFT JOIN SELECT col3, col4 FROM t2 WHERE id = $4 AND level <= $5"; sql != expected {
		t.Fatalf("invalid sql. [expected:%v] [actual:%v]", expected, sql)
	}

	if expected := []interface{}{7890, 1234, 2, 4567, 5}; !reflect.DeepEqual(args, expected) {
		t.Fatalf("invalid args. [expected:%v] [actual:%v]", expected, args)
	}
}
