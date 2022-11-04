package selector

var globalSelector Builder

var _ Builder = (*wrapSelector)(nil)

// wrapSelector wrapped Selector, help override global Selector implementation
type wrapSelector struct {
	Builder
}

// GlobalSelector returns global selector builder.
func GlobalSelector() Builder {
	return globalSelector
}

// SetGlobalSelector set global selector builder.
func SetGlobalSelector(builder Builder) {
	if globalSelector == nil {
		globalSelector = &wrapSelector{builder}
	}
	globalSelector.(*wrapSelector).Builder = builder
}
