/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidmux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/logutil-go/pkg/otel/correlationid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/trustbloc/logutil-go/pkg/otel/api"
)

func TestMuxMiddleware(t *testing.T) {
	const correlationID1 = "correlationID1"

	otel.SetTracerProvider(trace.NewTracerProvider())

	t.Run("with correlation ID in header", func(t *testing.T) {
		m := Middleware(GenerateNewFixedLengthIfNotFound(12))

		handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, correlationID, err := correlationid.FromContext(r.Context())
			assert.NoError(t, err)
			assert.Equal(t, correlationID1, correlationID)

			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(api.CorrelationIDHeader, correlationID1)

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("without correlation ID in header", func(t *testing.T) {
		m := Middleware(GenerateUUIDIfNotFound())

		handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, correlationID, err := correlationid.FromContext(r.Context())
			assert.NoError(t, err)
			assert.NotEmpty(t, correlationID)

			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})
}
