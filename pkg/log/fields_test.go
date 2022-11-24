/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//nolint:maintidx
func TestStandardFields(t *testing.T) {
	const module = "test_module"

	t.Run("console error", func(t *testing.T) {
		stdErr := newMockWriter()

		logger := New(module, WithStdErr(stdErr))

		logger.Error("Sample error", WithError(errors.New("some error")))

		require.Contains(t, stdErr.Buffer.String(), `Sample error	{"error": "some error"}`)
	})

	t.Run("json error", func(t *testing.T) {
		stdErr := newMockWriter()

		logger := New(module,
			WithStdErr(stdErr), WithEncoding(JSON),
		)

		logger.Error("Sample error", WithError(errors.New("some error")))

		l := unmarshalLogData(t, stdErr.Bytes())

		require.Equal(t, "test_module", l.Logger)
		require.Contains(t, l.Caller, "log/fields_test.go")
		require.Equal(t, "some error", l.Error)
		require.Equal(t, "error", l.Level)
	})

	t.Run("json fields 1", func(t *testing.T) {
		stdOut := newMockWriter()

		logger := New(module, WithStdOut(stdOut), WithEncoding(JSON))

		id := "123"
		name := "Joe"
		topic := "some topic"
		address := "https://localhost:8080"
		token := "someToken"
		path := "some/path"
		url := "some/url"
		duration := time.Second * 20
		txID := "some tx id"
		state := "some state"

		logger.Info("Some message",
			WithDuration(duration),
			WithHTTPStatus(http.StatusNotFound),
			WithID(id),
			WithName(name),
			WithPath(path),
			WithState(state),
			WithToken(token),
			WithTopic(topic),
			WithTxID(txID),
			WithURL(url),
			WithAddress(address),
		)

		t.Logf(stdOut.String())
		l := unmarshalLogData(t, stdOut.Bytes())

		require.Equal(t, 404, l.HTTPStatus)
		require.Equal(t, id, l.ID)
		require.Equal(t, name, l.Name)
		require.Equal(t, topic, l.Topic)
		require.Equal(t, token, l.Token)
		require.Equal(t, path, l.Path)
		require.Equal(t, url, l.URL)
		require.Equal(t, txID, l.TxID)
		require.Equal(t, state, l.State)
		require.Equal(t, address, l.Address)
	})
}

type logData struct {
	Level  string `json:"level"`
	Time   string `json:"time"`
	Logger string `json:"logger"`
	Caller string `json:"caller"`
	Error  string `json:"error"`

	HTTPStatus int    `json:"httpStatus"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	State      string `json:"state"`
	Token      string `json:"token"`
	Topic      string `json:"topic"`
	TxID       string `json:"txID"`
	URL        string `json:"url"`
	Address    string `json:"address"`
}

func unmarshalLogData(t *testing.T, b []byte) *logData {
	t.Helper()

	l := &logData{}

	require.NoError(t, json.Unmarshal(b, l))

	return l
}
