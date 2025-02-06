/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationid

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestTransport_RoundTrip(t *testing.T) {
	t.Run("No span", func(t *testing.T) {
		var rt mockRoundTripperFunc = func(req *http.Request) (*http.Response, error) {
			require.Len(t, req.Header.Get(api.CorrelationIDHeader), 8)

			return &http.Response{}, nil
		}

		transport := NewHTTPTransport(rt)

		ctx, correlationID, err := Set(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, correlationID)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
		require.NoError(t, err)

		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("With span", func(t *testing.T) {
		var correlationID string

		var rt mockRoundTripperFunc = func(req *http.Request) (*http.Response, error) {
			require.Equal(t, correlationID, req.Header.Get(api.CorrelationIDHeader))

			return &http.Response{}, nil
		}

		transport := NewHTTPTransport(rt)

		tp := trace.NewTracerProvider()

		otel.SetTracerProvider(tp)

		ctx, span := tp.Tracer("test").Start(context.Background(), "test")
		require.NotNil(t, span)

		var err error
		ctx, correlationID, err = Set(ctx)
		require.NoError(t, err)

		ctx, correlationID, err = Set(ctx)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
		require.NoError(t, err)

		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

type mockRoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn mockRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
