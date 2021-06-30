package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_atomicValue_Bool(t *testing.T) {
	var vlist []interface{}
	vlist = []interface{}{"1", "t", "T", "true", "TRUE", "True", true, 1, int32(1)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Bool()
		assert.NoError(t, err, b)
		assert.True(t, b, b)
	}

	vlist = []interface{}{"0", "f", "F", "false", "FALSE", "False", false, 0, int32(0)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Bool()
		assert.NoError(t, err, b)
		assert.False(t, b, b)
	}

	vlist = []interface{}{uint16(1), "bbb", "-1"}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Bool()
		assert.Error(t, err, b)
	}
}

func Test_atomicValue_Int(t *testing.T) {
	var vlist []interface{}
	vlist = []interface{}{"123123", float64(123123), int64(123123), int32(123123), 123123}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Int()
		assert.NoError(t, err, b)
		assert.Equal(t, int64(123123), b, b)
	}

	vlist = []interface{}{uint16(1), "bbb", "-x1", true}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Int()
		assert.Error(t, err, b)
	}
}

func Test_atomicValue_Float(t *testing.T) {
	var vlist []interface{}
	vlist = []interface{}{"123123.1", float64(123123.1)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Float()
		assert.NoError(t, err, b)
		assert.Equal(t, float64(123123.1), b, b)
	}

	vlist = []interface{}{float32(1123123), uint16(1), "bbb", "-x1"}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Float()
		assert.Error(t, err, b)
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
	var vlist []interface{}
	vlist = []interface{}{"1", float64(1), int64(1), 1, int64(1)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.String()
		assert.NoError(t, err, b)
		assert.Equal(t, "1", b, b)
	}

	v := atomicValue{}
	v.Store(true)
	b, err := v.String()
	assert.NoError(t, err, b)
	assert.Equal(t, "true", b, b)

	v = atomicValue{}
	v.Store(ts{
		Name: "test",
		Age:  10,
	})
	b, err = v.String()
	assert.NoError(t, err, b)
	assert.Equal(t, "test10", b, "test Stringer should be equal")
}

func Test_atomicValue_Duration(t *testing.T) {
	var vlist []interface{}
	vlist = []interface{}{int64(5)}
	for _, x := range vlist {
		v := atomicValue{}
		v.Store(x)
		b, err := v.Duration()
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(5), b)
	}
}

func Test_atomicValue_Scan(t *testing.T) {
	var err error
	v := atomicValue{}
	err = v.Scan(&struct {
		A string `json:"a"`
	}{"a"})
	assert.NoError(t, err)

	err = v.Scan(&struct {
		A string `json:"a"`
	}{"a"})
	assert.NoError(t, err)
}
