package apns2

import (
	"encoding/json"
)

// NOTE these structs and "Table" refer to https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/ApplePushService.html

// Payload payload.
type Payload struct {
	Aps    Aps    `json:"aps"`
	URL    string `json:"url"` // bilibili schedule
	TaskID string `json:"task_id"`
	Token  string `json:"tid"`
	Image  string `json:"image_url,omitempty"`
}

// Marshal marshals payload.
func (p *Payload) Marshal() []byte {
	payload, _ := json.Marshal(p)
	return payload
}

// Aps Apple Push Service request meta.
type Aps struct {
	// If this property is included, the system displays a standard alert or a banner, based on the user’s setting.
	// You can specify a string or a dictionary as the value of alert.
	// If you specify a string, it becomes the message text of an alert with two buttons: Close and View.
	// If the user taps View, the app launches.
	// If you specify a dictionary, refer to Table 5-2 for descriptions of the keys of this dictionary.
	// The JSON \U notation is not supported. Put the actual UTF-8 character in the alert text instead.
	Alert Alert `json:"alert,omitempty"`

	// The number to display as the badge of the app icon.
	// If this property is absent, the badge is not changed. To remove the badge, set the value of this property to 0.
	Badge int `json:"badge,omitempty"`

	// The name of a sound file in the app bundle or in the Library/Sounds folder of the app’s data container.
	// The sound in this file is played as an alert. If the sound file doesn’t exist or default is specified
	// as the value, the default alert sound is played. The audio must be in one of the audio data formats
	// that are compatible with system sounds; see Preparing Custom Alert Sounds for details.
	Sound string `json:"sound,omitempty"`

	// Provide this key with a value of 1 to indicate that new content is available.
	// Including this key and value means that when your app is launched in the background or resumed,
	// application:didReceiveRemoteNotification:fetchCompletionHandler: is called.
	ContentAvailable int `json:"content-available,omitempty"`

	// Provide this key with a string value that represents the identifier property of the
	// UIMutableUserNotificationCategory object you created to define custom actions.
	// To learn more about using custom actions, see Registering Your Actionable Notification Types.
	Category string `json:"category,omitempty"`
	// MutableContent .
	MutableContent int `json:"mutable-content,omitempty"`
}

// Alert alert message.
type Alert struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`

	// could support any more other field
}

// Response reponse message.
type Response struct {
	ApnsID string

	// Http status. (refer to Table 6-4)
	StatusCode int

	// The APNs error string indicating the reason for the notification failure (if
	// any). The error code is specified as a string. For a list of possible
	// values, see the Reason constants above.
	// If the notification was accepted, this value will be "".
	Reason string

	// If the value of StatusCode is 410, this is the last time at which APNs
	// confirmed that the device token was no longer valid for the topic.
	// Timestamp time.Time
}
