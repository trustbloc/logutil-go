/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidecho

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/trustbloc/logutil-go/pkg/log"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
)

var logger = log.New("correlationid-echo")

// Middleware reads the X-Correlation-Id header and, if found, sets the
// dts.correlation_id attribute on the current span.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if correlationID := c.Request().Header.Get(api.CorrelationIDHeader); correlationID != "" {
				ctx := c.Request().Context()

				span := trace.SpanFromContext(ctx)
				span.SetAttributes(attribute.String(api.CorrelationIDAttribute, correlationID))

				logger.Infoc(ctx, "Received HTTP request", log.WithCorrelationID(correlationID))
			}

			return next(c)
		}
	}
}
