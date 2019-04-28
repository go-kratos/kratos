package jpush

// Notification .
type Notification struct {
	Android *AndroidNotification `json:"android,omitempty"`
	IOS     *IOSNotification     `json:"ios,omitempty"`
}

// AndroidNotification .
type AndroidNotification struct {
	Alert      string      `json:"alert"`
	Title      string      `json:"title,omitempty"`
	AlertType  int         `json:"alert_type,omitempty"`
	BuilderID  int         `json:"builder_id,omitempty"`
	Style      int         `json:"style,omitempty"`
	BigPicPath string      `json:"big_pic_path,omitempty"`
	Extras     interface{} `json:"extras,omitempty"`
}

// IOSAlert .
type IOSAlert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// IOSNotification .
type IOSNotification struct {
	Alert            interface{} `json:"alert"`
	Sound            string      `json:"sound,omitempty"`
	Badge            int32       `json:"badge,omitempty"`
	ContentAvailable bool        `json:"content-available,omitempty"`
	MutableContent   bool        `json:"mutable-content,omitempty"`
	Category         string      `json:"category,omitempty"`
	Extras           interface{} `json:"extras,omitempty"`
}
