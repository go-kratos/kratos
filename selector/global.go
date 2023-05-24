package selector

var globalSelector = &wrapSelector{}

var _ Builder = (*wrapSelector)(nil)

// wrapSelector wrapped Selector, help override global Selector implementation.
type wrapSelector struct{ Builder }

// GlobalSelector returns global selector builder.
func GlobalSelector() Builder {
	if globalSelector.Builder != nil {
		return globalSelector
	}
	return nil
}

// SetGlobalSelector set global selector builder.
func SetGlobalSelector(builder Builder) {
	globalSelector.Builder = builder
}
