package gojson_test

import (
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/contrib/encoding/gojson"
)

type Person struct {
	Name string
}

func TestGoJson_Encode_Decode(t *testing.T) {
	c := new(gojson.GoJson)
	data, err := c.Marshal(Person{Name: "a"})
	if err != nil {
		t.Errorf("Marshal() should be nil, but got %s", err)
	}
	p := &Person{}
	err = c.Unmarshal(data, p)
	if err != nil {
		t.Errorf("Unmarshal() should be nil, but got %s", err)
	}
	if !reflect.DeepEqual(p.Name, "a") {
		t.Errorf("ID should be 'a', but got '%s'", p.Name)
	}
}

func TestGoJson_Name(t *testing.T) {
	c := new(gojson.GoJson)
	if !reflect.DeepEqual("gojson", c.Name()) {
		t.Errorf("Name() should be gojson, but got %s", c.Name())
	}
}
