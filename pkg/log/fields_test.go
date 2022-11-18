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
		command := "some command"
		topic := "some topic"
		hostURL := "https://localhost:8080"
		address := "https://localhost:8080"
		responseBody := []byte("response body")
		token := "someToken"
		totalRequests := 10
		responses := 9
		path := "some/path"
		url := "some/url"
		json := "{\"some\":\"json object\"}"
		sleep := time.Second * 10
		duration := time.Second * 20
		txID := "some tx id"
		state := "some state"
		profileID := "some profile id"
		parameter := "param1"
		parameters := &mockObject{Field1: "param1", Field2: 4612}
		totalMessages := 3

		logger.Info("Some message",
			WithCommand(command),
			WithDuration(duration),
			WithHTTPStatus(http.StatusNotFound),
			WithHostURL(hostURL),
			WithID(id),
			WithJSON(json),
			WithName(name),
			WithParameter(parameter),
			WithParameters(parameters),
			WithPath(path),
			WithProfileID(profileID),
			WithResponseBody(responseBody),
			WithResponses(responses),
			WithSleep(sleep),
			WithState(state),
			WithToken(token),
			WithTopic(topic),
			WithTotalRequests(totalRequests),
			WithTxID(txID),
			WithURL(url),
			WithAddress(address),
			WithTotalMessages(totalMessages),
		)

		t.Logf(stdOut.String())
		l := unmarshalLogData(t, stdOut.Bytes())

		require.Equal(t, 404, l.HTTPStatus)
		require.Equal(t, id, l.ID)
		require.Equal(t, name, l.Name)
		require.Equal(t, command, l.Command)
		require.Equal(t, topic, l.Topic)
		require.Equal(t, hostURL, l.HostURL)
		require.EqualValues(t, responseBody, l.ResponseBody)
		require.Equal(t, token, l.Token)
		require.Equal(t, totalRequests, l.TotalRequests)
		require.Equal(t, responses, l.Responses)
		require.Equal(t, path, l.Path)
		require.Equal(t, url, l.URL)
		require.Equal(t, json, l.JSON)
		require.Equal(t, sleep.String(), l.Sleep)
		require.Equal(t, txID, l.TxID)
		require.Equal(t, state, l.State)
		require.Equal(t, profileID, l.ProfileID)
		require.Equal(t, address, l.Address)
		require.Equal(t, totalMessages, l.TotalMessages)
	})
}

type mockObject struct {
	Field1 string
	Field2 int
}

type logData struct {
	Level  string `json:"level"`
	Time   string `json:"time"`
	Logger string `json:"logger"`
	Caller string `json:"caller"`
	Error  string `json:"error"`

	AdditionalMessage string      `json:"additionalMessage"`
	Command           string      `json:"command"`
	Duration          string      `json:"duration"`
	HTTPStatus        int         `json:"httpStatus"`
	HostURL           string      `json:"hostURL"`
	ID                string      `json:"id"`
	JSON              string      `json:"json"`
	Name              string      `json:"name"`
	Parameter         string      `json:"parameter"`
	Parameters        *mockObject `json:"parameters"`
	Path              string      `json:"path"`
	ProfileID         string      `json:"profileID"`
	ResponseBody      string      `json:"responseBody"`
	Responses         int         `json:"responses"`
	Sleep             string      `json:"sleep"`
	State             string      `json:"state"`
	Token             string      `json:"token"`
	Topic             string      `json:"topic"`
	TotalRequests     int         `json:"totalRequests"`
	TxID              string      `json:"transactionID"`
	URL               string      `json:"url"`
	Workers           int         `json:"workers"`
	Address           string      `json:"address"`
	TotalMessages     int         `json:"total-messages"`
}

func unmarshalLogData(t *testing.T, b []byte) *logData {
	t.Helper()

	l := &logData{}

	require.NoError(t, json.Unmarshal(b, l))

	return l
}
