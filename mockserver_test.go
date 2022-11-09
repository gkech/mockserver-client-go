package mockserver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/gkech/mockserver-client-go/create"
	"github.com/gkech/mockserver-client-go/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateExpectation(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mhc := mock.NewMockHttpClient(ctrl)

	c := Client{
		host:   "http://mockserver",
		client: mhc,
	}

	tests := map[string]struct {
		httpRequestBody *bytes.Reader
		httpStatusCode  int
		httpResponse    []byte
		hasError        bool
		expErr          error
	}{
		"success: expectation created": {
			httpRequestBody: bytes.NewReader([]byte(`{"httpRequest":{"method":"POST","path":"/some/resource"},"httpResponse":{"statusCode":201,"body":{"field":"value"}},"times":{"remainingTimes":1,"unlimited":false}}`)),
			httpStatusCode:  201,
			httpResponse:    []byte(`{}`),
		},
		"failure: expectation not created according to the status code": {
			httpRequestBody: bytes.NewReader([]byte(`{"httpRequest":{"method":"POST","path":"/some/resource"},"httpResponse":{"statusCode":201,"body":{"field":"value"}},"times":{"remainingTimes":1,"unlimited":false}}`)),
			httpStatusCode:  404,
			httpResponse:    []byte(`{}`),
			hasError:        true,
			expErr:          errors.New("error creating mockserver expectation. received status code: 404 with body: {}"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodPut, "http://mockserver/mockserver/expectation", tt.httpRequestBody)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			res := createHTTPResponse(tt.httpStatusCode, tt.httpResponse)
			err = res.Body.Close()
			require.NoError(t, err)

			mhc.EXPECT().Do(eqHTTPRequest(req)).Return(res, nil)

			err = c.CreateExpectation(defaultExpectation())
			if tt.hasError {
				assert.EqualError(t, err, tt.expErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func defaultExpectation() create.Expectation {
	return create.Expectation{Request: create.Request{
		Method: "POST",
		Path:   "/some/resource",
	},
		Response: create.Response{
			Status: 201,
			Body: struct {
				Field string `json:"field"`
			}{
				Field: "value",
			},
		},
		Times: create.CallTimes{
			RemainingTimes: 1,
			Unlimited:      false,
		}}
}

func createHTTPResponse(statusCode int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBuffer(body)),
	}
}

type httpRequestMatcher struct {
	r *http.Request
}

func (m httpRequestMatcher) Matches(arg interface{}) bool {
	x := arg.(*http.Request)

	m.r.GetBody = nil
	x.GetBody = nil

	return reflect.DeepEqual(m.r, x)
}

func (m httpRequestMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.r)
}

func eqHTTPRequest(r *http.Request) gomock.Matcher {
	return httpRequestMatcher{r: r}
}
