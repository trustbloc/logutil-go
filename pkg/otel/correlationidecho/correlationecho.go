/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidecho

import (
	"github.com/labstack/echo/v4"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Middleware reads the X-Correlation-Id header and, if found, sets the
// dts.correlation_id attribute on the current span.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := c.Request().Header.Get(api.CorrelationIDHeader)
			if correlationID != "" {
				span := trace.SpanFromContext(c.Request().Context())
				span.SetAttributes(attribute.String(api.CorrelationIDAttribute, correlationID))
			}

			return next(c)
		}
	}
}
