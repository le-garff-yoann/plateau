package body

// Body represents the standard return of an HTTP or *plateau* request.
type Body struct {
	Successes []string `json:"ok,omitempty"`
	Failures  []string `json:"ko,omitempty"`
}

// New returns a new `Body`.
func New() *Body {
	return &Body{}
}

// Ok adds one or more `string` as successful messages.
func (s *Body) Ok(ok ...string) *Body {
	s.Successes = append(s.Successes, ok...)

	return s
}

// Ko adds one or more `error` as error messages.
func (s *Body) Ko(ko ...error) *Body {
	for _, e := range ko {
		if e != nil {
			s.Failures = append(s.Failures, e.Error())
		}
	}

	return s
}
