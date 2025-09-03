package vemetric

type TrackEventOpts struct {
	EventName       string                 `json:"name"`
	UserIdentifier  string                 `json:"userIdentifier,omitempty"`
	UserDisplayName string                 `json:"displayName,omitempty"`
	// A map of key-value pairs to set on the event.
	EventData       map[string]any         `json:"customData,omitempty"`
	// Lets you set data on the user while tracking an event.
	UserData        UserData               `json:"userData,omitempty"`
}

type UpdateUserOpts struct {
	UserIdentifier  string         `json:"userIdentifier"`
	UserDisplayName string         `json:"displayName,omitempty"`
	UserData        UserData       `json:"data,omitempty"`
}

type UserData struct {
	// A map of key-value pairs to set on the user, overwriting any existing values.
	Set map[string]any     `json:"set,omitempty"`
	// A map of key-value pairs to set on the user, but only if the key does not already exist.
	SetOnce map[string]any `json:"setOnce,omitempty"`
	// An array of keys to remove from the user.
	Unset []string         `json:"unset,omitempty"`
}
