/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidecho

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/logutil-go/pkg/otel/correlationid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestMiddleware(t *testing.T) {
	const correlationID1 = "correlationID1"

	m := Middleware()

	handler := m(func(e echo.Context) error {
		_, correlationID, err := correlationid.FromContext(e.Request().Context())
		assert.NoError(t, err)
		assert.NotEmpty(t, correlationID)

		return nil
	})
	require.NotNil(t, handler)

	otel.SetTracerProvider(trace.NewTracerProvider())

	ctx, span := otel.GetTracerProvider().Tracer("test").Start(context.Background(), "test")
	defer span.End()

	t.Run("No correlation ID", func(t *testing.T) {
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()

		ectx := echo.New().NewContext(req, rec)

		require.NoError(t, handler(ectx))
	})

	t.Run("With correlation ID", func(t *testing.T) {
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		req.Header.Set("X-Correlation-Id", correlationID1)

		rec := httptest.NewRecorder()

		ectx := echo.New().NewContext(req, rec)

		require.NoError(t, handler(ectx))
	})
}

func TestMiddlewareGenerateNewID(t *testing.T) {
	t.Run("Fixed length correlation ID", func(t *testing.T) {
		m := Middleware(GenerateNewFixedLengthIfNotFound(12))

		handler := m(func(e echo.Context) error {
			_, correlationID, err := correlationid.FromContext(e.Request().Context())
			require.NoError(t, err)
			require.Len(t, correlationID, 12)

			return nil
		})
		require.NotNil(t, handler)

		otel.SetTracerProvider(trace.NewTracerProvider())

		ctx, span := otel.GetTracerProvider().Tracer("test").Start(context.Background(), "test")
		defer span.End()

		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()

		ectx := echo.New().NewContext(req, rec)

		require.NoError(t, handler(ectx))
	})

	t.Run("UUID correlation ID", func(t *testing.T) {
		m := Middleware(GenerateUUIDIfNotFound())

		handler := m(func(e echo.Context) error {
			_, correlationID, err := correlationid.FromContext(e.Request().Context())
			require.NoError(t, err)
			_, err = uuid.Parse(correlationID)
			require.NoError(t, err)

			return nil
		})
		require.NotNil(t, handler)

		otel.SetTracerProvider(trace.NewTracerProvider())

		ctx, span := otel.GetTracerProvider().Tracer("test").Start(context.Background(), "test")
		defer span.End()

		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()

		ectx := echo.New().NewContext(req, rec)

		require.NoError(t, handler(ectx))
	})
}
