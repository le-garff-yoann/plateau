package protocol

// Response ...
type Response string

const (
	// ResOK ...
	ResOK Response = "OK"
	// ResAccepted ...
	ResAccepted Response = "ACCEPTED"
	// ResBadRequest ...
	ResBadRequest Response = "BAD_REQUEST"
	// ResForbidden ...
	ResForbidden Response = "FORBIDDEN"
	// ResInternalError ...
	ResInternalError Response = "INTERNAL_ERROR"
	// ResNotImplemented ...
	ResNotImplemented Response = "NOT_IMPLEMENTED"
)

func (s Response) String() string {
	return string(s)
}

// ResponseContainer ...
type ResponseContainer struct {
	Response `json:"response"`
	Body     interface{} `json:"body,omitempty"`
}

func (s ResponseContainer) String() string {
	return string(s.Response)
}
