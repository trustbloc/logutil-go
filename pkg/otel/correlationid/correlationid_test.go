/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/trustbloc/logutil-go/pkg/otel/api"
)

func TestSet(t *testing.T) {
	t.Run("No trace ID", func(t *testing.T) {
		ctx, correlationID, err := Set(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, correlationID)

		b := baggage.FromContext(ctx)
		m := b.Member(api.CorrelationIDHeader)
		require.Equal(t, correlationID, m.Value())

		t.Run("With existing correlation ID", func(t *testing.T) {
			ctx2, correlationID2, err := Set(ctx)
			require.NoError(t, err)
			require.Equal(t, ctx, ctx2)
			require.Equal(t, correlationID, correlationID2)
		})

		t.Run("Nested contexts", func(t *testing.T) {
			type key struct{}

			ctx2, cancel := context.WithCancel(context.WithValue(ctx, key{}, "test"))
			defer cancel()

			b := baggage.FromContext(ctx2)
			m := b.Member(api.CorrelationIDHeader)
			require.Equal(t, correlationID, m.Value())
		})
	})

	t.Run("With trace ID", func(t *testing.T) {
		tp := trace.NewTracerProvider()

		otel.SetTracerProvider(tp)

		ctx, span := tp.Tracer("test").Start(context.Background(), "test")
		require.NotNil(t, span)

		ctx, correlationID, err := Set(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, correlationID)

		b := baggage.FromContext(ctx)
		m := b.Member(api.CorrelationIDHeader)
		require.Equal(t, correlationID, m.Value())
	})
}

func TestSetWithValue(t *testing.T) {
	ctx := context.Background()

	ctx2, err := SetWithValue(ctx, "id1")
	require.NoError(t, err)

	b := baggage.FromContext(ctx2)
	m := b.Member(api.CorrelationIDHeader)
	require.Equal(t, "id1", m.Value())

	ctx3, err := SetWithValue(ctx2, "id1")
	require.NoError(t, err)
	require.Equal(t, ctx2, ctx3)
}
