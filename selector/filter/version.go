package filter

import (
	"context"

	"github.com/go-kratos/kratos/v2/selector"
)

// Version is node verion filter.
func Version(version string) selector.Filter {
	kf := func(ctx context.Context) Keep {
		return func(node selector.Node) bool {
			return node.Version() == version
		}
	}
	return BaseFilter(kf)
}
