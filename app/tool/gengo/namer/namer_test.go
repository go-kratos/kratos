package namer

import (
	"reflect"
	"testing"

	"go-common/app/tool/gengo/types"
)

func TestNameStrategy(t *testing.T) {
	u := types.Universe{}

	// Add some types.
	base := u.Type(types.Name{Package: "foo/bar", Name: "Baz"})
	base.Kind = types.Struct

	tmp := u.Type(types.Name{Package: "", Name: "[]bar.Baz"})
	tmp.Kind = types.Slice
	tmp.Elem = base

	tmp = u.Type(types.Name{Package: "", Name: "map[string]bar.Baz"})
	tmp.Kind = types.Map
	tmp.Key = types.String
	tmp.Elem = base

	tmp = u.Type(types.Name{Package: "foo/other", Name: "Baz"})
	tmp.Kind = types.Struct
	tmp.Members = []types.Member{{
		Embedded: true,
		Type:     base,
	}}

	tmp = u.Type(types.Name{Package: "", Name: "chan Baz"})
	tmp.Kind = types.Chan
	tmp.Elem = base

	u.Type(types.Name{Package: "", Name: "string"})

	o := Orderer{NewPublicNamer(0)}
	order := o.OrderUniverse(u)
	orderedNames := make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}
	expect := []string{"Baz", "Baz", "ChanBaz", "MapStringToBaz", "SliceBaz", "String"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}

	o = Orderer{NewRawNamer("my/package", nil)}
	order = o.OrderUniverse(u)
	orderedNames = make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}

	expect = []string{"[]bar.Baz", "bar.Baz", "chan bar.Baz", "map[string]bar.Baz", "other.Baz", "string"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}

	o = Orderer{NewRawNamer("foo/bar", nil)}
	order = o.OrderUniverse(u)
	orderedNames = make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}

	expect = []string{"Baz", "[]Baz", "chan Baz", "map[string]Baz", "other.Baz", "string"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}

	o = Orderer{NewPublicNamer(1)}
	order = o.OrderUniverse(u)
	orderedNames = make([]string, len(order))
	for i, t := range order {
		orderedNames[i] = o.Name(t)
	}
	expect = []string{"BarBaz", "ChanBarBaz", "MapStringToBarBaz", "OtherBaz", "SliceBarBaz", "String"}
	if e, a := expect, orderedNames; !reflect.DeepEqual(e, a) {
		t.Errorf("Wanted %#v, got %#v", e, a)
	}
}
