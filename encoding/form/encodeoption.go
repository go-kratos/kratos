package form

// EncodeOption is the encoding options.
type EncodeOption struct {
	ForceProtoTextAsKey bool // forces to use proto field text name.
}

// Encode creates a new EncodeOption instance.
func Encode() *EncodeOption {
	return &EncodeOption{}
}

// UseProtoTextAsKey forces to use proto field text name as key.
func (opt *EncodeOption) UseProtoTextAsKey(isUse bool) *EncodeOption {
	opt.ForceProtoTextAsKey = isUse
	return opt
}

// MergeEncodeOptions merges the options.
func MergeEncodeOptions(opts ...*EncodeOption) *EncodeOption {
	opt := new(EncodeOption)
	for _, o := range opts {
		if o == nil {
			continue
		}

		if o.ForceProtoTextAsKey {
			opt.ForceProtoTextAsKey = true
		}
	}
	return opt
}
