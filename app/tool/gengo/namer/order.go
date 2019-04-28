package namer

import (
	"sort"

	"go-common/app/tool/gengo/types"
)

// Orderer produces an ordering of types given a Namer.
type Orderer struct {
	Namer
}

// OrderUniverse assigns a name to every type in the Universe, including Types,
// Functions and Variables, and returns a list sorted by those names.
func (o *Orderer) OrderUniverse(u types.Universe) []*types.Type {
	list := tList{
		namer: o.Namer,
	}
	for _, p := range u {
		for _, t := range p.Types {
			list.types = append(list.types, t)
		}
		for _, f := range p.Functions {
			list.types = append(list.types, f)
		}
		for _, v := range p.Variables {
			list.types = append(list.types, v)
		}
	}
	sort.Sort(list)
	return list.types
}

// OrderTypes assigns a name to every type, and returns a list sorted by those
// names.
func (o *Orderer) OrderTypes(typeList []*types.Type) []*types.Type {
	list := tList{
		namer: o.Namer,
		types: typeList,
	}
	sort.Sort(list)
	return list.types
}

type tList struct {
	namer Namer
	types []*types.Type
}

func (t tList) Len() int           { return len(t.types) }
func (t tList) Less(i, j int) bool { return t.namer.Name(t.types[i]) < t.namer.Name(t.types[j]) }
func (t tList) Swap(i, j int)      { t.types[i], t.types[j] = t.types[j], t.types[i] }
