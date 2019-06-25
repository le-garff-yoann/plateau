package protocol

import "encoding/json"

// Notification ...
type Notification string

const (
	// NDealChange ...
	NDealChange Notification = "DEAL_CHANGE"
)

// NotificationContainer ...
type NotificationContainer struct {
	Notification `json:"notification"`
	Body         interface{} `json:"body,omitempty"`
}

func (s *NotificationContainer) String() string {
	b, _ := json.Marshal(s)

	return string(b)
}
