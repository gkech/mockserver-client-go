// Package create ...
package create

// Expectation represents an expectation in mockserver.
type Expectation struct {
	Request  Request   `json:"httpRequest"`
	Response Response  `json:"httpResponse"`
	Times    CallTimes `json:"times"`
}

// Request is the request that an expectation matches against.
type Request struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Body   interface{} `json:"body,omitempty"`
}

// Response is an untyped response from mockserver expectation.
type Response struct {
	Status int         `json:"statusCode"`
	Body   interface{} `json:"body,omitempty"`
}

// CallTimes configures the call times for an expectation.
type CallTimes struct {
	RemainingTimes int  `json:"remainingTimes"`
	Unlimited      bool `json:"unlimited"`
}
