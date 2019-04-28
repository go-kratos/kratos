package types

import (
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	mapType := &Type{
		Name: Name{Package: "", Name: "map[string]string"},
		Kind: Map,
		Key:  String,
		Elem: String,
	}
	m := []Member{
		{
			Name:     "Baz",
			Embedded: true,
			Type: &Type{
				Name: Name{Package: "pkg", Name: "Baz"},
				Kind: Struct,
				Members: []Member{
					{Name: "Foo", Type: String},
					{
						Name:     "Qux",
						Embedded: true,
						Type: &Type{
							Name:    Name{Package: "pkg", Name: "Qux"},
							Kind:    Struct,
							Members: []Member{{Name: "Zot", Type: String}},
						},
					},
				},
			},
		},
		{Name: "Bar", Type: String},
		{
			Name:     "NotSureIfLegal",
			Embedded: true,
			Type:     mapType,
		},
	}
	e := []Member{
		{Name: "Bar", Type: String},
		{Name: "NotSureIfLegal", Type: mapType, Embedded: true},
		{Name: "Foo", Type: String},
		{Name: "Zot", Type: String},
	}
	if a := FlattenMembers(m); !reflect.DeepEqual(e, a) {
		t.Errorf("Expected \n%#v\n, got \n%#v\n", e, a)
	}
}
