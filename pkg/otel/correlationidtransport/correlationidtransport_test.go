/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidtransport

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/trustbloc/logutil-go/pkg/otel/api"
)

func TestTransport_RoundTrip(t *testing.T) {
	var rt mockRoundTripperFunc = func(req *http.Request) (*http.Response, error) {
		correlationID := req.Header.Get(api.CorrelationIDHeader)

		require.Len(t, correlationID, 8)
		return &http.Response{}, nil
	}

	transport := New(rt, WithCorrelationIDLength(8))

	t.Run("No span", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com", nil)
		require.NoError(t, err)

		resp, err := transport.RoundTrip(req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("With span", func(t *testing.T) {
		tp := trace.NewTracerProvider()

		otel.SetTracerProvider(tp)

		ctx, span := tp.Tracer("test").Start(context.Background(), "test")
		require.NotNil(t, span)

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
