/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/trustbloc/logutil-go/pkg/log"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"go.opentelemetry.io/otel/baggage"
)

var logger = log.New("correlationid")

type options struct {
	generateFixedLengthID bool
	generateUUID          bool
	value                 string
	correlationIDLength   int
}

// Opt is an option for the FromContext function.
type Opt func(*options)

// GenerateUUIDIfNotFound configures the FromContext function to generate a UUID as the correlation ID.
func GenerateUUIDIfNotFound() Opt {
	return func(o *options) {
		o.generateUUID = true
	}
}

// GenerateNewFixedLengthIfNotFound configures the FromContext function to generate a new
// correlation ID if none is found in the context.
func GenerateNewFixedLengthIfNotFound(length int) Opt {
	return func(o *options) {
		o.generateFixedLengthID = true
		o.correlationIDLength = length
	}
}

// WithValue configures the FromContext function to use the provided correlation ID.
func WithValue(correlationID string) Opt {
	return func(o *options) {
		o.value = correlationID
	}
}

// FromContext returns the correlation ID from the given context. If a correlation ID is not found
// in the context then:
//   - If GenerateUUIDIfNotFound option is set, a new UUID is generated and set on the returned context.
//   - If GenerateNewFixedLengthIfNotFound option is set, a new fixed-length correlation ID
//     is generated and set on the returned context.
//   - If WithValue is set then the given correlation ID is set on the returned context.
//   - If none of the above options is specified then the existing context and empty string are returned.
func FromContext(ctx context.Context, opts ...Opt) (context.Context, string, error) {
	options := &options{}

	for _, opt := range opts {
		opt(options)
	}

	b := baggage.FromContext(ctx)

	m := b.Member(api.CorrelationIDHeader)
	if m.Value() != "" {
		if options.value == "" || m.Value() == options.value {
			logger.Debugc(ctx, "Found correlation ID in baggage")

			return ctx, m.Value(), nil
		}
	}

	if !options.generateFixedLengthID && !options.generateUUID && options.value == "" {
		return ctx, "", nil
	}

	correlationID := options.value

	if correlationID == "" {
		var err error
		correlationID, err = generateID(options)
		if err != nil {
			return nil, "", fmt.Errorf("generate correlation ID: %w", err)
		}

		logger.Debug("Generated correlation ID", log.WithCorrelationID(correlationID))
	} else {
		logger.Debug("Using correlation ID from options", log.WithCorrelationID(correlationID))
	}

	m, err := baggage.NewMember(api.CorrelationIDHeader, correlationID)
	if err != nil {
		return nil, "", fmt.Errorf("create baggage member: %w", err)
	}

	b, err = baggage.New(m)
	if err != nil {
		return nil, "", fmt.Errorf("create baggage: %w", err)
	}

	return baggage.ContextWithBaggage(ctx, b), correlationID, nil
}

func generateID(options *options) (string, error) {
	if options.generateUUID {
		return uuid.NewString(), nil
	}

	bytes := make([]byte, options.correlationIDLength/2) //nolint:gomnd

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}
