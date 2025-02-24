/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidecho

import (
	"github.com/labstack/echo/v4"
	"github.com/trustbloc/logutil-go/pkg/log"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"github.com/trustbloc/logutil-go/pkg/otel/correlationid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var logger = log.New("correlationid-echo")

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
		o.generateFixedLengthID = false
	}
}

// GenerateNewFixedLengthIfNotFound configures the FromContext function to generate
// a new correlation ID if none is found in the context.
func GenerateNewFixedLengthIfNotFound(length int) Opt {
	return func(o *options) {
		o.generateFixedLengthID = true
		o.correlationIDLength = length
		o.generateUUID = false
	}
}

// Middleware reads the X-Correlation-Id header and, if found, sets the
// dts.correlation_id attribute on the current span.
func Middleware(opts ...Opt) echo.MiddlewareFunc {
	options := &options{
		generateUUID: true,
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			defer func() {
				c.SetRequest(req)
			}()

			ctx := req.Context()

			correlationID := c.Request().Header.Get(api.CorrelationIDHeader)
			if correlationID != "" {
				logger.Debugc(ctx, "Received HTTP request with correlation ID in header", log.WithCorrelationID(correlationID))

				var err error
				ctx, _, err = correlationid.FromContext(ctx, correlationid.WithValue(correlationID))
				if err != nil {
					return err
				}
			} else {
				var err error
				ctx, correlationID, err = correlationid.FromContext(ctx, getOptions(options)...)
				if err != nil {
					return err
				}

				logger.Debugc(ctx, "Generated new correlation ID since none was found in the HTTP header")
			}

			span := trace.SpanFromContext(ctx)
			span.SetAttributes(attribute.String(api.CorrelationIDAttribute, correlationID))

			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	}
}

func getOptions(opts *options) []correlationid.Opt {
	var copts []correlationid.Opt

	if opts.generateFixedLengthID {
		copts = append(copts, correlationid.GenerateNewFixedLengthIfNotFound(opts.correlationIDLength))
	}

	if opts.generateUUID {
		copts = append(copts, correlationid.GenerateUUIDIfNotFound())
	}

	return copts
}
