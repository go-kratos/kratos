package config

import (
	"fmt"
	"testing"
	"time"
)

func Test_atomicValue_Bool(t *testing.T) {
	vlist := []interface{}{"1", "t", "T", "true", "TRUE", "True", true, 1, int32(1)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Bool()
		if err != nil {
			t.Fatal(`err is not nil`)
		}
		if !b {
			t.Fatal(`b is not equal to true`)
		}
	}

	vlist = []interface{}{"0", "f", "F", "false", "FALSE", "False", false, 0, int32(0)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Bool()
		if err != nil {
			t.Fatal(`err is not nil`)
		}
		if b {
			t.Fatal(`b is not equal to false`)
		}
	}

	vlist = []interface{}{"bbb", "-1"}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		_, err := v.Bool()
		if err == nil {
			t.Fatal(`err is nil`)
		}
	}
}

func Test_atomicValue_Int(t *testing.T) {
	vlist := []interface{}{"123123", float64(123123), int64(123123), int32(123123), 123123}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Int()
		if err != nil {
			t.Fatal(`err is not nil`)
		}
		if b != 123123 {
			t.Fatal(`b is not equal to 123123`)
		}
	}

	vlist = []interface{}{"bbb", "-x1", true}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		_, err := v.Int()
		if err == nil {
			t.Fatal(`err is nil`)
		}
	}
}

func Test_atomicValue_Float(t *testing.T) {
	vlist := []interface{}{"123123.1", float64(123123.1)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Float()
		if err != nil {
			t.Fatal(`err is not nil`)
		}
		if b != float64(123123.1) {
			t.Fatal(`b is not equal to 123123.1`)
		}
	}

	vlist = []interface{}{"bbb", "-x1"}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		_, err := v.Float()
		if err == nil {
			t.Fatal(`err is nil`)
		}
	}
}

type ts struct {
	Name string
	Age  int
}

func (t ts) String() string {
	return fmt.Sprintf("%s%d", t.Name, t.Age)
}

func Test_atomicValue_String(t *testing.T) {
	vlist := []interface{}{"1", float64(1), int64(1), 1, int64(1)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.String()
		if err != nil {
			t.Fatal(`err is not nil`)
		}
		if b != "1" {
			t.Fatal(`b is not equal to 1`)
		}
	}

	v := atomicValue{}
	v.Store(true)
	b, err := v.String()
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	if b != "true" {
		t.Fatal(`b is not equal to "true"`)
	}

	v = atomicValue{}
	v.Store(ts{
		Name: "test",
		Age:  10,
	})
	b, err = v.String()
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	if b != "test10" {
		t.Fatal(`b is not equal to "test10"`)
	}
}

func Test_atomicValue_Duration(t *testing.T) {
	vlist := []interface{}{int64(5)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Duration()
		if err != nil {
			t.Fatal(`err is not nil`)
		}
		if b != time.Duration(5) {
			t.Fatal(`b is not equal to time.Duration(5)`)
		}
	}
}

func Test_atomicValue_Slice(t *testing.T) {
	vlist := []interface{}{int64(5)}
	v := atomicValue{}
	v.Store(vlist)
	slices, err := v.Slice()
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	for _, v := range slices {
		b, err := v.Duration()
		if err != nil {
			t.Fatal(`err is not nil`)
		}
		if b != time.Duration(5) {
			t.Fatal(`b is not equal to time.Duration(5)`)
		}
	}
}

func Test_atomicValue_Map(t *testing.T) {
	vlist := make(map[string]interface{})
	vlist["5"] = int64(5)
	vlist["text"] = "text"
	v := atomicValue{}
	v.Store(vlist)
	m, err := v.Map()
	if err != nil {
		t.Fatal(`err is not nil`)
	}
	for k, v := range m {
		if k == "5" {
			b, err := v.Duration()
			if err != nil {
				t.Fatal(`err is not nil`)
			}
			if b != time.Duration(5) {
				t.Fatal(`b is not equal to time.Duration(5)`)
			}
		} else {
			b, err := v.String()
			if err != nil {
				t.Fatal(`err is not nil`)
			}
			if b != "text" {
				t.Fatal(`b is not equal to "text"`)
			}
		}
	}
}

func Test_atomicValue_Scan(t *testing.T) {
	var err error
	v := atomicValue{}
	err = v.Scan(&struct {
		A string `json:"a"`
	}{"a"})
	if err != nil {
		t.Fatal(`err is not nil`)
	}

	err = v.Scan(&struct {
		A string `json:"a"`
	}{"a"})
	if err != nil {
		t.Fatal(`err is not nil`)
	}
}
