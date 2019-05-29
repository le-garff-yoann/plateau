package body

// Body ...
type Body struct {
	Successes []string `json:"ok,omitempty"`
	Failures  []string `json:"ko,omitempty"`
}

// New ...
func New() *Body {
	return &Body{}
}

// Ok ...
func (s *Body) Ok(ok ...string) *Body {
	for _, e := range ok {
		s.Successes = append(s.Successes, e)
	}

	return s
}

// Ko ...
func (s *Body) Ko(ko ...error) *Body {
	for _, e := range ko {
		s.Failures = append(s.Failures, e.Error())
	}

	return s
}
