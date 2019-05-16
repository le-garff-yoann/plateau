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
