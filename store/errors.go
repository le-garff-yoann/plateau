package store

// DuplicateError is the error that must be returned
// by the `Transaction` methods when a new entry is already existing.
type DuplicateError string

func (s DuplicateError) Error() string {
	return string(s)
}

// DontExistError is the error that must be returned
// by the `Transaction` methods when the search entry is non-existent.
type DontExistError string

func (s DontExistError) Error() string {
	return string(s)
}

// PlayerParticipationError is the error that must be returned
// by the `Transaction` methods when entry or exit of a `protocol.Player`
// from a `protocol.Match` is impossible.
type PlayerParticipationError string

func (s PlayerParticipationError) Error() string {
	return string(s)
}
