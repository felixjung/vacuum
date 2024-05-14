package remote

import (
	"fmt"
	"net/http"
	"time"
)

func NewAuthenticatedClient(authorizationHeader string) *http.Client {
	return &http.Client{
		Timeout: time.Second * 120,
		Transport: &addHeaderTransport{
			T:          http.DefaultTransport,
			authHeader: authorizationHeader,
		},
	}
}

type addHeaderTransport struct {
	T          http.RoundTripper
	authHeader string
}

func (adt *addHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if adt.authHeader != "" {
		req.Header.Add("Authorization", adt.authHeader)
	}
	return adt.T.RoundTrip(req)
}

type RemoteHandlerFunc func(url string) (*http.Response, error)

func NewAuthenticatedRemoteHandlerFunc(authoorizationHeader string) RemoteHandlerFunc {
	c := NewAuthenticatedClient(authoorizationHeader)

	return func(url string) (*http.Response, error) {
		resp, err := c.Get(url)
		if err != nil {
			return nil, fmt.Errorf("fetch remote ref: %v", err)
		}

		return resp, nil
	}
}
