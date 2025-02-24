/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidmux

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/trustbloc/logutil-go/pkg/log"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"github.com/trustbloc/logutil-go/pkg/otel/correlationid"
)

var logger = log.New("correlationid-mux")

type options struct {
	generateFixedLengthID bool
	generateUUID          bool
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

// GenerateNewFixedLengthIfNotFound configures the FromContext function to generate
// a new correlation ID if none is found in the context.
func GenerateNewFixedLengthIfNotFound(length int) Opt {
	return func(o *options) {
		o.generateFixedLengthID = true
		o.correlationIDLength = length
	}
}

// Middleware returns a mux middleware that sets the correlation ID in the header of the HTTP request.
func Middleware(opts ...Opt) mux.MiddlewareFunc {
	options := &options{
		generateUUID: true,
	}

	for _, opt := range opts {
		opt(options)
	}

	var copts []correlationid.Opt

	if options.generateFixedLengthID {
		copts = append(copts, correlationid.GenerateNewFixedLengthIfNotFound(options.correlationIDLength))
	}

	if options.generateUUID {
		copts = append(copts, correlationid.GenerateUUIDIfNotFound())
	}

	return func(handler http.Handler) http.Handler {
		return &MuxMiddleware{
			options: copts,
			handler: handler,
		}
	}
}

// MuxMiddleware is a mux middleware that sets the correlation ID in the header of the HTTP request.
type MuxMiddleware struct {
	options []correlationid.Opt
	handler http.Handler
}

// ServeHTTP sets the correlation ID in the header of the HTTP request.
func (m *MuxMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	correlationID := req.Header.Get(api.CorrelationIDHeader)
	if correlationID != "" {
		logger.Debugc(ctx, "Received HTTP request with correlation ID in header", log.WithCorrelationID(correlationID))

		var err error
		ctx, _, err = correlationid.FromContext(ctx, correlationid.WithValue(correlationID))
		if err != nil {
			logger.Warnc(ctx, "Failed to set correlation ID in context", log.WithError(err))
		}
	} else {
		var err error
		ctx, correlationID, err = correlationid.FromContext(ctx, m.options...)
		if err != nil {
			logger.Warnc(ctx, "Failed to set correlation ID in context", log.WithError(err))
		} else {
			logger.Debugc(ctx, "Generated new correlation ID since none was found in the HTTP header")
		}
	}

	if correlationID != "" {
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String(api.CorrelationIDAttribute, correlationID))
	}

	m.handler.ServeHTTP(w, req.WithContext(ctx))
}
