package verify

const (
	// TypeJSON defines a type as described in mockserver documentation.
	TypeJSON = "JSON"
)

// Request is the request that an expectation matches against.
type Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   Body   `json:"body,omitempty"`
}

// Times configures the call times for verifying a request.
type Times struct {
	AtLeast int `json:"atLeast,omitempty"`
	AtMost  int `json:"atMost,omitempty"`
}

// Body is the body of the request that is verified.
type Body struct {
	Type      string `json:"type"`
	MatchType string `json:"matchType"`
}

// Expectation configures the call times for verifying a request.
type Expectation struct {
	Request Request `json:"httpRequest"`
	Times   Times   `json:"times"`
}
