package selector

var globalSelector Builder

// GlobalSelector returns global selector builder.
func GlobalSelector() Builder {
	return globalSelector
}

// SetGlobalSelector set global selector builder.
func SetGlobalSelector(builder Builder) {
	globalSelector = builder
}
