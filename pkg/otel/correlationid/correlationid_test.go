/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationid

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/logutil-go/pkg/otel/api"
	"go.opentelemetry.io/otel/baggage"
)

func TestSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		origCtx := context.Background()

		ctx, correlationID, err := FromContext(origCtx)
		require.NoError(t, err)
		require.Empty(t, correlationID)
		require.Equal(t, origCtx, ctx)

		ctx, correlationID, err = FromContext(context.Background(), GenerateNewFixedLengthIfNotFound(12))
		require.NoError(t, err)
		require.NotEmpty(t, correlationID)
		require.Len(t, correlationID, 12)

		b := baggage.FromContext(ctx)
		m := b.Member(api.CorrelationIDHeader)
		require.Equal(t, correlationID, m.Value())

		t.Run("With existing correlation ID", func(t *testing.T) {
			ctx2, correlationID2, err := FromContext(ctx)
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

	t.Run("Generate UUID", func(t *testing.T) {
		ctx := context.Background()

		ctx2, correlationID, err := FromContext(ctx, GenerateUUIDIfNotFound())
		require.NoError(t, err)
		_, err = uuid.Parse(correlationID)
		require.NoError(t, err)

		b := baggage.FromContext(ctx2)
		m := b.Member(api.CorrelationIDHeader)
		require.Equal(t, correlationID, m.Value())

		ctx3, correlationID, err := FromContext(ctx2, WithValue("id1"))
		require.NoError(t, err)
		require.NotEqual(t, ctx2, ctx3)
		require.Equal(t, "id1", correlationID)
	})

	t.Run("ID in options", func(t *testing.T) {
		ctx := context.Background()

		ctx2, correlationID, err := FromContext(ctx, WithValue("id1"))
		require.NoError(t, err)
		require.Equal(t, "id1", correlationID)

		b := baggage.FromContext(ctx2)
		m := b.Member(api.CorrelationIDHeader)
		require.Equal(t, "id1", m.Value())

		ctx3, correlationID, err := FromContext(ctx2, WithValue("id1"))
		require.NoError(t, err)
		require.Equal(t, ctx2, ctx3)
		require.Equal(t, "id1", correlationID)
	})
}
