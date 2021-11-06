package filter

import (
	"context"

	"github.com/go-kratos/kratos/v2/selector"
)

//
type Keep func(node selector.Node) bool

type KeepFactory func(ctx context.Context) Keep

func BaseFilter(keepFactory KeepFactory) selector.Filter {
	return func(ctx context.Context, nodes *[]selector.Node) {
		if nodes == nil {
			return
		}

		keep := keepFactory(ctx)
		length := len(*nodes)
		for i := 0; i < length; i++ {
			if !keep((*nodes)[i]) {
				if i == length-1 {
					length--
					break
				}
				for ; length > i; length-- {
					if keep((*nodes)[length-1]) {
						(*nodes)[i] = (*nodes)[length-1]
						length--
						break
					}
				}
			}
		}
		*nodes = (*nodes)[:length]
	}
}
