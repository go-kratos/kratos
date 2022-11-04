package selector

var globalSelector = &wrapSelector{}

// wrapSelector wrapped Selector.
// help override global Selector implementation
type wrapSelector struct {
	Builder
}

// GlobalSelector returns global selector builder.
func GlobalSelector() Builder {
	return globalSelector
}

// SetGlobalSelector set global selector builder.
func SetGlobalSelector(builder Builder) {
	globalSelector.Builder = builder
}
