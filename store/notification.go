package store

// MatchNotification notify that a `protocol.Match` has changed.
//
// It currently carries no information and simply
// reports that the match should probably be reloaded on the client side.
type MatchNotification struct {
	ID string `json:"id"`
}
