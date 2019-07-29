package store

// MatchNotification notify that a `protocol.Match` has changed.
//
// It currently carries no information and simply
// reports that the match should probably be reloaded on the client side.
type MatchNotification struct{}

// MatchNotificationsIterator represents the iterator in its most classical form.
// He is specialized to return only the `protocol.Match`.
//	- `Next()` fetches the next state of the match.
//	- `Close()` closes the iterator and stops calls to `Next()`.
type MatchNotificationsIterator interface {
	Next(*MatchNotification) bool
	Close() error
}
