package protocol

import "encoding/json"

// Notification is a message intended to be sent
// to clients as notification of an event.
type Notification string

const (
	// NDealChange signals a change on a `Deal`.
	NDealChange Notification = "DEAL_CHANGE"
)

// NotificationContainer is `Notification`'s container.
type NotificationContainer struct {
	Notification `json:"notification"`
	Body         interface{} `json:"body,omitempty"`
}

func (s *NotificationContainer) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
}
