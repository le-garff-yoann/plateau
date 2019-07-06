package protocol

// Response represents a response sent to a
// `Player` who sent a `Request`.
type Response string

const (
	// ResOK means that the processing of the associated
	// `Request` has successfully completed.
	ResOK Response = "OK"
	// ResBadRequest means that the associated
	// `Request` is malformed.
	ResBadRequest Response = "BAD_REQUEST"
	// ResForbidden means that the associated
	// `Request` is fordidden.
	ResForbidden Response = "FORBIDDEN"
	// ResInternalError means that the associated
	// `Request` generated an unexpected server-side error.
	ResInternalError Response = "INTERNAL_ERROR"
	// ResNotImplemented that the associated
	// `Request` is not recognized.
	ResNotImplemented Response = "NOT_IMPLEMENTED"
)

func (s Response) String() string {
	return string(s)
}

// ResponseContainer is `Response`'s container.
type ResponseContainer struct {
	Response `json:"response"`
	Body     interface{} `json:"body,omitempty"`
}

func (s ResponseContainer) String() string {
	return string(s.Response)
}
