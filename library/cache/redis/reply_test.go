// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

type valueError struct {
	v   interface{}
	err error
}

func ve(v interface{}, err error) valueError {
	return valueError{v, err}
}

var replyTests = []struct {
	name     interface{}
	actual   valueError
	expected valueError
}{
	{
		"ints([v1, v2])",
		ve(Ints([]interface{}{[]byte("4"), []byte("5")}, nil)),
		ve([]int{4, 5}, nil),
	},
	{
		"ints(nil)",
		ve(Ints(nil, nil)),
		ve([]int(nil), ErrNil),
	},
	{
		"strings([v1, v2])",
		ve(Strings([]interface{}{[]byte("v1"), []byte("v2")}, nil)),
		ve([]string{"v1", "v2"}, nil),
	},
	{
		"strings(nil)",
		ve(Strings(nil, nil)),
		ve([]string(nil), ErrNil),
	},
	{
		"byteslices([v1, v2])",
		ve(ByteSlices([]interface{}{[]byte("v1"), []byte("v2")}, nil)),
		ve([][]byte{[]byte("v1"), []byte("v2")}, nil),
	},
	{
		"byteslices(nil)",
		ve(ByteSlices(nil, nil)),
		ve([][]byte(nil), ErrNil),
	},
	{
		"values([v1, v2])",
		ve(Values([]interface{}{[]byte("v1"), []byte("v2")}, nil)),
		ve([]interface{}{[]byte("v1"), []byte("v2")}, nil),
	},
	{
		"values(nil)",
		ve(Values(nil, nil)),
		ve([]interface{}(nil), ErrNil),
	},
	{
		"float64(1.0)",
		ve(Float64([]byte("1.0"), nil)),
		ve(float64(1.0), nil),
	},
	{
		"float64(nil)",
		ve(Float64(nil, nil)),
		ve(float64(0.0), ErrNil),
	},
	{
		"uint64(1)",
		ve(Uint64(int64(1), nil)),
		ve(uint64(1), nil),
	},
	{
		"uint64(-1)",
		ve(Uint64(int64(-1), nil)),
		ve(uint64(0), ErrNegativeInt),
	},
}

func TestReply(t *testing.T) {
	for _, rt := range replyTests {
		if errors.Cause(rt.actual.err) != rt.expected.err {
			t.Errorf("%s returned err %v, want %v", rt.name, rt.actual.err, rt.expected.err)
			continue
		}
		if !reflect.DeepEqual(rt.actual.v, rt.expected.v) {
			t.Errorf("%s=%+v, want %+v", rt.name, rt.actual.v, rt.expected.v)
		}
	}
}

// dial wraps DialDefaultServer() with a more suitable function name for examples.
func dial() (Conn, error) {
	return DialDefaultServer()
}

func ExampleBool() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Do("SET", "foo", 1)
	exists, _ := Bool(c.Do("EXISTS", "foo"))
	fmt.Printf("%#v\n", exists)
	// Output:
	// true
}

func ExampleInt() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Do("SET", "k1", 1)
	n, _ := Int(c.Do("GET", "k1"))
	fmt.Printf("%#v\n", n)
	n, _ = Int(c.Do("INCR", "k1"))
	fmt.Printf("%#v\n", n)
	// Output:
	// 1
	// 2
}

func ExampleInts() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Do("SADD", "set_with_integers", 4, 5, 6)
	ints, _ := Ints(c.Do("SMEMBERS", "set_with_integers"))
	fmt.Printf("%#v\n", ints)
	// Output:
	// []int{4, 5, 6}
}

func ExampleString() {
	c, err := dial()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Do("SET", "hello", "world")
	s, err := String(c.Do("GET", "hello"))
	fmt.Printf("%#v\n", s)
	// Output:
	// "world"
}
