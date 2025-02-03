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

// Middleware reads the X-Correlation-Id header and, if found, sets the
// dts.correlation_id attribute on the current span.
func Middleware() echo.MiddlewareFunc {
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
				ctx, err = correlationid.SetWithValue(ctx, correlationID)
				if err != nil {
					return err
				}
			} else {
				var err error
				ctx, correlationID, err = correlationid.Set(ctx)
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
