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
	"strings"

	"github.com/trustbloc/logutil-go/pkg/log"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

const (
	nilTraceID          = "00000000000000000000000000000000"
	correlationIDLength = 8
)

var logger = log.New("correlationid")

// Set derives the correlation ID from the OpenTelemetry trace ID and sets it on the returned context.
// If no trace ID is available, a random correlation ID is generated.
func Set(ctx context.Context) (context.Context, string, error) {
	var correlationID string

	b := baggage.FromContext(ctx)

	m := b.Member(api.CorrelationIDHeader)
	if m.Value() != "" {
		correlationID = m.Value()

		return ctx, correlationID, nil
	}

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

	ctx, err := SetWithValue(ctx, correlationID)
	return ctx, correlationID, err
}

// SetWithValue sets the correlation ID on the returned context.
func SetWithValue(ctx context.Context, correlationID string) (context.Context, error) {
	b := baggage.FromContext(ctx)

	m := b.Member(api.CorrelationIDHeader)
	if m.Value() == correlationID {
		logger.Infoc(ctx, "Found correlation ID in baggage")

		return ctx, nil
	}

	logger.Infoc(ctx, "Setting correlation ID in baggage", log.WithCorrelationID(correlationID))

	m, err := baggage.NewMember(api.CorrelationIDHeader, correlationID)
	if err != nil {
		return nil, fmt.Errorf("create baggage member: %w", err)
	}

	b, err = baggage.New(m)
	if err != nil {
		return nil, fmt.Errorf("create baggage: %w", err)
	}

	return baggage.ContextWithBaggage(ctx, b), nil
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
