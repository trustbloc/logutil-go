/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationid

import (
	"net/http"

	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"go.opentelemetry.io/otel/baggage"
)

// Transport is an HTTP RoundTripper that adds a correlation ID to the request header.
type Transport struct {
	defaultTransport http.RoundTripper
}

// NewHTTPTransport creates a new HTTP Transport.
func NewHTTPTransport(defaultTransport http.RoundTripper) *Transport {
	return &Transport{
		defaultTransport: defaultTransport,
	}
}

// RoundTrip executes a single HTTP transaction.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	b := baggage.FromContext(ctx)

	m := b.Member(api.CorrelationIDHeader)
	if m.Value() != "" {
		logger.Infoc(ctx, "Found correlation ID in baggage")

		req = req.Clone(ctx)
		req.Header.Add(api.CorrelationIDHeader, m.Value())
	}

	return t.defaultTransport.RoundTrip(req)
}
