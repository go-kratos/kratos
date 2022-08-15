package transport

import (
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
)

var globalSelector selector.Builder = wrr.NewBuilder()

func GlobalSelector() selector.Builder {
	return globalSelector
}

func SetGlobalSelector() selector.Builder {
	return globalSelector
}
