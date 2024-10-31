/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidtransport

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/trustbloc/logutil-go/pkg/otel/api"
)

const (
	nilTraceID                 = "00000000000000000000000000000000"
	defaultCorrelationIDLength = 8
)

// Transport is an http.RoundTripper that adds a correlation ID to the request.
type Transport struct {
	defaultTransport    http.RoundTripper
	correlationIDLength int
}

type Opt func(*Transport)

// WithCorrelationIDLength sets the length of the correlation ID.
func WithCorrelationIDLength(length int) Opt {
	return func(t *Transport) {
		t.correlationIDLength = length
	}
}

// New creates a new Transport.
func New(defaultTransport http.RoundTripper, opts ...Opt) *Transport {
	t := &Transport{
		defaultTransport:    defaultTransport,
		correlationIDLength: defaultCorrelationIDLength,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// RoundTrip executes a single HTTP transaction.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var correlationID string

	span := trace.SpanFromContext(req.Context())

	traceID := span.SpanContext().TraceID().String()
	if traceID == "" || traceID == nilTraceID {
		var err error
		correlationID, err = t.generateID()
		if err != nil {
			return nil, fmt.Errorf("generate correlation ID: %w", err)
		}
	} else {
		correlationID = t.shortenID(traceID)
	}

	clonedReq := req.Clone(req.Context())
	clonedReq.Header.Add(api.CorrelationIDHeader, correlationID)

	return t.defaultTransport.RoundTrip(clonedReq)
}

func (t *Transport) generateID() (string, error) {
	bytes := make([]byte, t.correlationIDLength/2) //nolint:gomnd

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func (t *Transport) shortenID(id string) string {
	hash := sha256.Sum256([]byte(id))
	return strings.ToUpper(hex.EncodeToString(hash[:t.correlationIDLength/2])) //nolint:gomnd
}
