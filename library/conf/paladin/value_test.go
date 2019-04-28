package paladin

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testUnmarshler struct {
	Text string
	Int  int
}

func TestValueUnmarshal(t *testing.T) {
	s := `
		int = 100
		text = "hello"	
	`
	v := Value{val: s, raw: s}
	obj := new(testUnmarshler)
	assert.Nil(t, v.UnmarshalTOML(obj))
	// error
	v = Value{val: nil, raw: ""}
	assert.NotNil(t, v.UnmarshalTOML(obj))
}

func TestValue(t *testing.T) {
	var tests = []struct {
		in  interface{}
		out interface{}
	}{
		{
			"text",
			"text",
		},
		{
			time.Duration(time.Second * 10),
			"10s",
		},
		{
			int64(100),
			int64(100),
		},
		{
			float64(100.1),
			float64(100.1),
		},
		{
			true,
			true,
		},
		{
			nil,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.in), func(t *testing.T) {
			v := Value{val: test.in, raw: fmt.Sprint(test.in)}
			switch test.in.(type) {
			case nil:
				s, err := v.String()
				assert.NotNil(t, err)
				assert.Equal(t, s, "", test.in)
				i, err := v.Int64()
				assert.NotNil(t, err)
				assert.Equal(t, i, int64(0), test.in)
				f, err := v.Float64()
				assert.NotNil(t, err)
				assert.Equal(t, f, float64(0.0), test.in)
				b, err := v.Bool()
				assert.NotNil(t, err)
				assert.Equal(t, b, false, test.in)
			case string:
				val, err := v.String()
				assert.Nil(t, err)
				assert.Equal(t, val, test.out.(string), test.in)
			case int64:
				val, err := v.Int()
				assert.Nil(t, err)
				assert.Equal(t, val, int(test.out.(int64)), test.in)
				val32, err := v.Int32()
				assert.Nil(t, err)
				assert.Equal(t, val32, int32(test.out.(int64)), test.in)
				val64, err := v.Int64()
				assert.Nil(t, err)
				assert.Equal(t, val64, test.out.(int64), test.in)
			case float64:
				val32, err := v.Float32()
				assert.Nil(t, err)
				assert.Equal(t, val32, float32(test.out.(float64)), test.in)
				val64, err := v.Float64()
				assert.Nil(t, err)
				assert.Equal(t, val64, test.out.(float64), test.in)
			case bool:
				val, err := v.Bool()
				assert.Nil(t, err)
				assert.Equal(t, val, test.out.(bool), test.in)
			case time.Duration:
				v.val = test.out
				val, err := v.Duration()
				assert.Nil(t, err)
				assert.Equal(t, val, test.in.(time.Duration), test.out)
			}
		})
	}
}

func TestValueSlice(t *testing.T) {
	var tests = []struct {
		in  interface{}
		out interface{}
	}{
		{
			nil,
			nil,
		},
		{
			[]interface{}{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			[]interface{}{1, 2, 3},
			[]int64{1, 2, 3},
		},
		{
			[]interface{}{1.1, 1.2, 1.3},
			[]float64{1.1, 1.2, 1.3},
		},
		{
			[]interface{}{true, false, true},
			[]bool{true, false, true},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprint(test.in), func(t *testing.T) {
			v := Value{val: test.in, raw: fmt.Sprint(test.in)}
			switch test.in.(type) {
			case nil:
				var s []string
				assert.NotNil(t, v.Slice(&s))
			case []string:
				var s []string
				assert.Nil(t, v.Slice(&s))
				assert.Equal(t, s, test.out)
			case []int64:
				var s []int64
				assert.Nil(t, v.Slice(&s))
				assert.Equal(t, s, test.out)
			case []float64:
				var s []float64
				assert.Nil(t, v.Slice(&s))
				assert.Equal(t, s, test.out)
			case []bool:
				var s []bool
				assert.Nil(t, v.Slice(&s))
				assert.Equal(t, s, test.out)
			}
		})
	}
}

func BenchmarkValueInt(b *testing.B) {
	v := &Value{val: int64(100), raw: "100"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.Int64()
		}
	})
}
func BenchmarkValueFloat(b *testing.B) {
	v := &Value{val: float64(100.1), raw: "100.1"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.Float64()
		}
	})
}
func BenchmarkValueBool(b *testing.B) {
	v := &Value{val: true, raw: "true"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.Bool()
		}
	})
}
func BenchmarkValueString(b *testing.B) {
	v := &Value{val: "text", raw: "text"}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v.String()
		}
	})
}

func BenchmarkValueSlice(b *testing.B) {
	v := &Value{val: []interface{}{1, 2, 3}, raw: "100"}
	b.RunParallel(func(pb *testing.PB) {
		var slice []int64
		for pb.Next() {
			v.Slice(&slice)
		}
	})
}
