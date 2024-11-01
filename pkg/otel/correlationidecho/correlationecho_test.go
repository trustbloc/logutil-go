/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package correlationidecho

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	const correlationID1 = "correlationID1"

	m := Middleware()

	handler := m(func(c echo.Context) error {
		return nil
	})
	require.NotNil(t, handler)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Correlation-Id", correlationID1)

	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := handler(ctx)
	require.NoError(t, err)
}
