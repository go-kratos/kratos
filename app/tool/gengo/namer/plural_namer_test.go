package namer

import (
	"testing"

	"go-common/app/tool/gengo/types"
)

func TestPluralNamer(t *testing.T) {
	exceptions := map[string]string{
		// The type name is already in the plural form
		"Endpoints": "endpoints",
	}
	public := NewPublicPluralNamer(exceptions)
	private := NewPrivatePluralNamer(exceptions)

	cases := []struct {
		typeName        string
		expectedPrivate string
		expectedPublic  string
	}{
		{
			"I",
			"i",
			"I",
		},
		{
			"Pod",
			"pods",
			"Pods",
		},
		{
			"Entry",
			"entries",
			"Entries",
		},
		{
			"Endpoints",
			"endpoints",
			"Endpoints",
		},
		{
			"Bus",
			"buses",
			"Buses",
		},
		{
			"Fizz",
			"fizzes",
			"Fizzes",
		},
		{
			"Search",
			"searches",
			"Searches",
		},
		{
			"Autograph",
			"autographs",
			"Autographs",
		},
		{
			"Dispatch",
			"dispatches",
			"Dispatches",
		},
		{
			"Earth",
			"earths",
			"Earths",
		},
		{
			"City",
			"cities",
			"Cities",
		},
		{
			"Ray",
			"rays",
			"Rays",
		},
		{
			"Fountain",
			"fountains",
			"Fountains",
		},
		{
			"Life",
			"lives",
			"Lives",
		},
		{
			"Leaf",
			"leaves",
			"Leaves",
		},
	}
	for _, c := range cases {
		testType := &types.Type{Name: types.Name{Name: c.typeName}}
		if e, a := c.expectedPrivate, private.Name(testType); e != a {
			t.Errorf("Unexpected result from private plural namer. Expected: %s, Got: %s", e, a)
		}
		if e, a := c.expectedPublic, public.Name(testType); e != a {
			t.Errorf("Unexpected result from public plural namer. Expected: %s, Got: %s", e, a)
		}
	}
}
