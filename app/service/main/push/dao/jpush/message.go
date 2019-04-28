package jpush

// Message .
type Message struct {
	Content     string                 `json:"msg_content"`
	Title       string                 `json:"title,omitempty"`
	ContentType string                 `json:"content_type,omitempty"`
	Extras      map[string]interface{} `json:"extras,omitempty"`
}

// SetContent .
func (m *Message) SetContent(c string) {
	m.Content = c

}

// SetTitle .
func (m *Message) SetTitle(title string) {
	m.Title = title
}

// SetContentType .
func (m *Message) SetContentType(t string) {
	m.ContentType = t
}

// AddExtras .
func (m *Message) AddExtras(key string, value interface{}) {
	if m.Extras == nil {
		m.Extras = make(map[string]interface{})
	}
	m.Extras[key] = value
}
