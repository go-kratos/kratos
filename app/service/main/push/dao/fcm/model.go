package fcm

import "time"

// Message represents fcm request message
type (
	Message struct {
		// Data parameter specifies the custom key-value pairs of the message's payload.
		//
		// For example, with data:{"score":"3x1"}:
		//
		// On iOS, if the message is sent via APNS, it represents the custom data fields.
		// If it is sent via FCM connection server, it would be represented as key value dictionary
		// in AppDelegate application:didReceiveRemoteNotification:.
		// On Android, this would result in an intent extra named score with the string value 3x1.
		// The key should not be a reserved word ("from" or any word starting with "google" or "gcm").
		// Do not use any of the words defined in this table (such as collapse_key).
		// Values in string types are recommended. You have to convert values in objects
		// or other non-string data types (e.g., integers or booleans) to string.
		//
		Data interface{} `json:"data,omitempty"`

		// To this parameter specifies the recipient of a message.
		//
		// The value must be a registration token, notification key, or topic.
		// Do not set this field when sending to multiple topics. See Condition.
		To string `json:"to,omitempty"`

		// RegistrationIDs for all registration ids
		// This parameter specifies a list of devices
		// (registration tokens, or IDs) receiving a multicast message.
		// It must contain at least 1 and at most 1000 registration tokens.
		// Use this parameter only for multicast messaging, not for single recipients.
		// Multicast messages (sending to more than 1 registration tokens)
		// are allowed using HTTP JSON format only.
		RegistrationIDs []string `json:"registration_ids,omitempty"`

		// CollapseKey This parameter identifies a group of messages
		// (e.g., with collapse_key: "Updates Available") that can be collapsed,
		// so that only the last message gets sent when delivery can be resumed.
		// This is intended to avoid sending too many of the same messages when the
		// device comes back online or becomes active (see delay_while_idle).
		CollapseKey string `json:"collapse_key,omitempty"`

		// Priority Sets the priority of the message. Valid values are "normal" and "high."
		// On iOS, these correspond to APNs priorities 5 and 10.
		// By default, notification messages are sent with high priority, and data messages
		// are sent with normal priority. Normal priority optimizes the client app's battery
		// consumption and should be used unless immediate delivery is required. For messages
		// with normal priority, the app may receive the message with unspecified delay.
		// When a message is sent with high priority, it is sent immediately, and the app
		// can wake a sleeping device and open a network connection to your server.
		// For more information, see Setting the priority of a message.
		Priority string `json:"priority,omitempty"`

		// Notification parameter specifies the predefined, user-visible key-value pairs of
		// the notification payload. See Notification payload support for detail.
		// For more information about notification message and data message options, see
		// Notification
		Notification Notification `json:"notification,omitempty"`

		// ContentAvailable On iOS, use this field to represent content-available
		// in the APNS payload. When a notification or message is sent and this is set
		// to true, an inactive client app is awoken. On Android, data messages wake
		// the app by default. On Chrome, currently not supported.
		ContentAvailable bool `json:"content_available,omitempty"`

		// DelayWhenIdle When this parameter is set to true, it indicates that
		// the message should not be sent until the device becomes active.
		// The default value is false.
		DelayWhileIdle bool `json:"delay_while_idle,omitempty"`

		// TimeToLive This parameter specifies how long (in seconds) the message
		// should be kept in FCM storage if the device is offline. The maximum time
		// to live supported is 4 weeks, and the default value is 4 weeks.
		// For more information, see
		// https://firebase.google.com/docs/cloud-messaging/concept-options#ttl
		TimeToLive int `json:"time_to_live,omitempty"`

		// RestrictedPackageName This parameter specifies the package name of the
		// application where the registration tokens must match in order to
		// receive the message.
		RestrictedPackageName string `json:"restricted_package_name,omitempty"`

		// DryRun This parameter, when set to true, allows developers to test
		// a request without actually sending a message.
		// The default value is false
		DryRun bool `json:"dry_run,omitempty"`

		// Condition to set a logical expression of conditions that determine the message target
		// This parameter specifies a logical expression of conditions that determine the message target.
		// Supported condition: Topic, formatted as "'yourTopic' in topics". This value is case-insensitive.
		// Supported operators: &&, ||. Maximum two operators per topic message supported.
		Condition string `json:"condition,omitempty"`

		// Currently for iOS 10+ devices only. On iOS, use this field to represent mutable-content in the APNS payload.
		// When a notification is sent and this is set to true, the content of the notification can be modified before
		// it is displayed, using a Notification Service app extension. This parameter will be ignored for Android and web.
		MutableContent bool `json:"mutable_content,omitempty"`

		Android Android `json:"android,omitempty"`
	}

	Android struct {
		Priority string `json:"priority,omitempty"`
	}

	// Result Downstream result from FCM, sent in the "results" field of the Response packet
	Result struct {
		// String specifying a unique ID for each successfully processed message.
		MessageID string `json:"message_id"`

		// Optional string specifying the canonical registration token for the
		// client app that the message was processed and sent to. Sender should
		// use this value as the registration token for future requests.
		// Otherwise, the messages might be rejected.
		RegistrationID string `json:"registration_id"`

		// String specifying the error that occurred when processing the message
		// for the recipient. The possible values can be found in table 9 here:
		// https://firebase.google.com/docs/cloud-messaging/http-server-ref#table9
		Error string `json:"error"`
	}

	// Response represents fcm response message - (tokens and topics)
	Response struct {
		Ok         bool
		StatusCode int

		// MulticastID a unique ID (number) identifying the multicast message.
		MulticastID int `json:"multicast_id"`

		// Success number of messages that were processed without an error.
		Success int `json:"success"`

		// Fail number of messages that could not be processed.
		Fail int `json:"failure"`

		// CanonicalIDs number of results that contain a canonical registration token.
		// A canonical registration ID is the registration token of the last registration
		// requested by the client app. This is the ID that the server should use
		// when sending messages to the device.
		CanonicalIDs int `json:"canonical_ids"`

		// Results Array of objects representing the status of the messages processed. The objects are listed in the same order as the request (i.e., for each registration ID in the request, its result is listed in the same index in the response).
		// message_id: String specifying a unique ID for each successfully processed message.
		// registration_id: Optional string specifying the canonical registration token for the client app that the message was processed and sent to. Sender should use this value as the registration token for future requests. Otherwise, the messages might be rejected.
		// error: String specifying the error that occurred when processing the message for the recipient. The possible values can be found in table 9.
		Results []Result `json:"results,omitempty"`

		// The topic message ID when FCM has successfully received the request and will attempt to deliver to all subscribed devices.
		MsgID int `json:"message_id,omitempty"`

		// Error that occurred when processing the message. The possible values can be found in table 9.
		Err string `json:"error,omitempty"`

		// RetryAfter
		RetryAfter string
	}

	// Notification notification message payload
	Notification struct {
		// Title indicates notification title. This field is not visible on iOS phones and tablets.
		Title string `json:"title,omitempty"`

		// Body indicates notification body text.
		Body string `json:"body,omitempty"`

		// Sound indicates a sound to play when the device receives a notification.
		// Sound files can be in the main bundle of the client app or in the
		// Library/Sounds folder of the app's data container.
		// See the iOS Developer Library for more information.
		// http://apple.co/2jaGqiE
		Sound string `json:"sound,omitempty"`

		// Badge indicates the badge on the client app home icon.
		Badge string `json:"badge,omitempty"`

		// Icon indicates notification icon. Sets value to myicon for drawable resource myicon.
		// If you don't send this key in the request, FCM displays the launcher icon specified
		// in your app manifest.
		Icon string `json:"icon,omitempty"`

		// Tag indicates whether each notification results in a new entry in the notification
		// drawer on Android. If not set, each request creates a new notification.
		// If set, and a notification with the same tag is already being shown,
		// the new notification replaces the existing one in the notification drawer.
		Tag string `json:"tag,omitempty"`

		// Color indicates color of the icon, expressed in #rrggbb format
		Color string `json:"color,omitempty"`

		// ClickAction indicates the action associated with a user click on the notification.
		// When this is set, an activity with a matching intent filter is launched when user
		// clicks the notification.
		ClickAction string `json:"click_action,omitempty"`

		// BodyLockKey indicates the key to the body string for localization. Use the key in
		// the app's string resources when populating this value.
		BodyLocKey string `json:"body_loc_key,omitempty"`

		// BodyLocArgs indicates the string value to replace format specifiers in the body
		// string for localization. For more information, see Formatting and Styling.
		BodyLocArgs string `json:"body_loc_args,omitempty"`

		// TitleLocKey indicates the key to the title string for localization.
		// Use the key in the app's string resources when populating this value.
		TitleLocKey string `json:"title_loc_key,omitempty"`

		// TitleLocArgs indicates the string value to replace format specifiers in the title string for
		// localization. For more information, see
		// https://developer.android.com/guide/topics/resources/string-resource.html#FormattingAndStyling
		TitleLocArgs string `json:"title_loc_args,omitempty"`
	}
)

// GetRetryAfterTime converts the retry after response header to a time.Duration
func (r *Response) GetRetryAfterTime() (time.Duration, error) {
	return time.ParseDuration(r.RetryAfter)
}
