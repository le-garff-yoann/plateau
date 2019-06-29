package store

// DuplicateError ...
type DuplicateError string

func (s DuplicateError) Error() string {
	return string(s)
}

// DontExistError ...
type DontExistError string

func (s DontExistError) Error() string {
	return string(s)
}

// PlayerParticipationError ...
type PlayerParticipationError string

func (s PlayerParticipationError) Error() string {
	return string(s)
}
