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

// Middleware returns a mux middleware that sets the correlation ID in the header of the HTTP request.
func Middleware() mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return &MuxMiddleware{
			handler: handler,
		}
	}
}

// MuxMiddleware is a mux middleware that sets the correlation ID in the header of the HTTP request.
type MuxMiddleware struct {
	handler http.Handler
}

// ServeHTTP sets the correlation ID in the header of the HTTP request.
func (tw *MuxMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	correlationID := req.Header.Get(api.CorrelationIDHeader)
	if correlationID != "" {
		logger.Debugc(ctx, "Received HTTP request with correlation ID in header", log.WithCorrelationID(correlationID))

		var err error
		ctx, err = correlationid.SetWithValue(ctx, correlationID)
		if err != nil {
			logger.Warnc(ctx, "Failed to set correlation ID in context", log.WithError(err))
		}
	} else {
		var err error
		ctx, correlationID, err = correlationid.Set(ctx)
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

	tw.handler.ServeHTTP(w, req.WithContext(ctx))
}
