package protocol

// Notification ...
type Notification string

const (
	// NTransactionChange ...
	NTransactionChange Notification = "TRANSACTION_CHANGE"
)

// NotificationContainer ...
type NotificationContainer struct {
	Notification `json:"notification"`
	Body         interface{} `json:"body,omitempty"`
}
