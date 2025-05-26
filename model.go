package vemetric

type TrackEventOpts struct {
	EventName      string                 `json:"name"`
	UserIdentifier string                 `json:"userIdentifier,omitempty"`
	EventData      map[string]any         `json:"customData,omitempty"`
	UserData       UserData               `json:"userData,omitempty"`
}

type UpdateUserOpts struct {
	UserIdentifier string         `json:"userIdentifier"`
	UserData       UserData       `json:"data,omitempty"`
}

type UserData struct {
	Set map[string]any     `json:"set,omitempty"`
	SetOnce map[string]any `json:"setOnce,omitempty"`
	Unset []string         `json:"unset,omitempty"`
}
