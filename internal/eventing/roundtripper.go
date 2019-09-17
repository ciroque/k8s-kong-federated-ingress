package eventing

import (
	"net/http"
	"strings"
)

type RoundTripper struct {
	Headers      []string
	RoundTripper http.RoundTripper
}

func (roundTripper *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := new(http.Request)
	*newRequest = *req
	newRequest.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		newRequest.Header[k] = append([]string(nil), s...)
	}
	for _, s := range roundTripper.Headers {
		split := strings.SplitN(s, ":", 2)
		if len(split) >= 2 {
			newRequest.Header[split[0]] = append([]string(nil), split[1])
		}
	}
	return roundTripper.RoundTripper.RoundTrip(newRequest)
}
