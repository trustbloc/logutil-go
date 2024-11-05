/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationid

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/trustbloc/logutil-go/pkg/log"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
)

const (
	nilTraceID          = "00000000000000000000000000000000"
	correlationIDLength = 8
)

var logger = log.New("correlationid")

type contextKey struct{}

// Set derives the correlation ID from the OpenTelemetry trace ID and sets it on the returned context.
// If no trace ID is available, a random correlation ID is generated.
func Set(ctx context.Context) (context.Context, string, error) {
	var correlationID string

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	if traceID != "" && traceID != nilTraceID {
		correlationID = deriveID(traceID)

		logger.Debugc(ctx, "Derived correlation ID from trace ID", log.WithCorrelationID(correlationID))
	} else {
		var err error
		correlationID, err = generateID()
		if err != nil {
			return nil, "", fmt.Errorf("generate correlation ID: %w", err)
		}

		logger.Debug("Generated correlation ID", log.WithCorrelationID(correlationID))
	}

	return context.WithValue(ctx, contextKey{}, correlationID), correlationID, nil
}

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
	correlationID, ok := req.Context().Value(contextKey{}).(string)
	if ok {
		req = req.Clone(req.Context())
		req.Header.Add(api.CorrelationIDHeader, correlationID)
	}

	return t.defaultTransport.RoundTrip(req)
}

func generateID() (string, error) {
	bytes := make([]byte, correlationIDLength/2) //nolint:gomnd

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}

func deriveID(id string) string {
	hash := sha256.Sum256([]byte(id))

	return strings.ToUpper(hex.EncodeToString(hash[:correlationIDLength/2])) //nolint:gomnd
}
