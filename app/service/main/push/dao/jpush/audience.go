package jpush

const (
	_audienceTag    = "tag"
	_audienceTagAnd = "tag_and"
	_audienceAlias  = "alias"
	_audienceID     = "registration_id"
	_audienceAll    = "all"
)

// Audience .
type Audience struct {
	Object   interface{}
	audience map[string][]string
}

// All .
func (a *Audience) All() {
	a.Object = _audienceAll
}

// SetID .
func (a *Audience) SetID(ids []string) {
	a.set(_audienceID, ids)
}

// SetTag .
func (a *Audience) SetTag(tags []string) {
	a.set(_audienceTag, tags)
}

// SetTagAnd .
func (a *Audience) SetTagAnd(tags []string) {
	a.set(_audienceTagAnd, tags)
}

// SetAlias .
func (a *Audience) SetAlias(alias []string) {
	a.set(_audienceAlias, alias)
}

func (a *Audience) set(key string, v []string) {
	if a.Object == nil {
		a.audience = map[string][]string{key: v}
		a.Object = a.audience
	}
	a.audience[key] = v
}
