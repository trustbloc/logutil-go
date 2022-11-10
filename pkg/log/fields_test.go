/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
		require.Equal(t, "Sample error", l.Msg)
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
		msg := "Some message"
		hostURL := "https://localhost:8080"
		address := "https://localhost:8080"
		responseBody := []byte("response body")
		token := "someToken"
		totalRequests := 10
		responses := 9
		concurrencyReq := 3
		workers := 4
		path := "some/path"
		url := "some/url"
		json := "{\"some\":\"json object\"}"
		jsonResolution := "json/resolution"
		sleep := time.Second * 10
		duration := time.Second * 20
		event := &mockObject{
			Field1: "event1",
			Field2: 123,
		}
		idToken := "some id token"
		vpToken := "some vp token"
		txID := "some tx id"
		presDefID := "some pd id"
		state := "some state"
		profileID := "some profile id"
		parameter := "param1"
		parameters := &mockObject{Field1: "param1", Field2: 4612}
		brokers := []string{"broker"}
		totalMessages := 3

		dockerComposeCmd := strings.Join([]string{
			"docker-compose",
			"-f",
			"/path/to/composeFile.yaml",
			"up",
			"--force-recreate",
			"-d",
		}, " ")
		certPoolSize := 3

		logger.Info("Some message",
			WithAdditionalMessage(msg),
			WithCertPoolSize(certPoolSize),
			WithCommand(command),
			WithConcurrencyRequests(concurrencyReq),
			WithDockerComposeCmd(dockerComposeCmd),
			WithDuration(duration),
			WithEvent(event),
			WithHTTPStatus(http.StatusNotFound),
			WithHostURL(hostURL),
			WithID(id),
			WithIDToken(idToken),
			WithJSON(json),
			WithJSONResolution(jsonResolution),
			WithName(name),
			WithParameter(parameter),
			WithParameters(parameters),
			WithPath(path),
			WithPresDefID(presDefID),
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
			WithUserLogLevel(DEBUG.String()),
			WithVPToken(vpToken),
			WithWorkers(workers),
			WithAddress(address),
			WithMessageBrokers(brokers),
			WithTotalMessages(totalMessages),
		)

		t.Logf(stdOut.String())
		l := unmarshalLogData(t, stdOut.Bytes())

		require.Equal(t, 404, l.HTTPStatus)
		require.Equal(t, `DEBUG`, l.UserLogLevel)
		require.Equal(t, id, l.ID)
		require.Equal(t, name, l.Name)
		require.Equal(t, command, l.Command)
		require.Equal(t, topic, l.Topic)
		require.Equal(t, msg, l.Msg)
		require.Equal(t, hostURL, l.HostURL)
		require.EqualValues(t, responseBody, l.ResponseBody)
		require.Equal(t, token, l.Token)
		require.Equal(t, totalRequests, l.TotalRequests)
		require.Equal(t, responses, l.Responses)
		require.Equal(t, concurrencyReq, l.ConcurrencyRequests)
		require.Equal(t, workers, l.Workers)
		require.Equal(t, path, l.Path)
		require.Equal(t, url, l.URL)
		require.Equal(t, json, l.JSON)
		require.Equal(t, jsonResolution, l.JSONResolution)
		require.Equal(t, sleep.String(), l.Sleep)
		require.Equal(t, event, l.Event)
		require.Equal(t, dockerComposeCmd, l.DockerComposeCmd)
		require.Equal(t, certPoolSize, l.CertPoolSize)
		require.Equal(t, idToken, l.IDToken)
		require.Equal(t, vpToken, l.VPToken)
		require.Equal(t, txID, l.TxID)
		require.Equal(t, presDefID, l.PresDefID)
		require.Equal(t, state, l.State)
		require.Equal(t, profileID, l.ProfileID)
		require.Equal(t, address, l.Address)
		require.Equal(t, brokers, l.MessageBrokers)
		require.Equal(t, totalMessages, l.TotalMessages)
	})

	t.Run("json fields 2", func(t *testing.T) {
		stdOut := newMockWriter()

		logger := New(module, WithStdOut(stdOut), WithEncoding(JSON))

		logger.Info("Some message")

		l := unmarshalLogData(t, stdOut.Bytes())

		require.Equal(t, `Some message`, l.Msg)
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
	Msg    string `json:"msg"`
	Error  string `json:"error"`

	AdditionalMessage   string      `json:"additionalMessage"`
	CertPoolSize        int         `json:"certPoolSize"`
	Command             string      `json:"command"`
	ConcurrencyRequests int         `json:"concurrencyRequests"`
	DockerComposeCmd    string      `json:"dockerComposeCmd"`
	Duration            string      `json:"duration"`
	Event               *mockObject `json:"event"`
	HTTPStatus          int         `json:"httpStatus"`
	HostURL             string      `json:"hostURL"`
	ID                  string      `json:"id"`
	IDToken             string      `json:"idToken"`
	JSON                string      `json:"json"`
	JSONResolution      string      `json:"jsonResolution"`
	Name                string      `json:"name"`
	Parameter           string      `json:"parameter"`
	Parameters          *mockObject `json:"parameters"`
	Path                string      `json:"path"`
	PresDefID           string      `json:"presDefinitionID"`
	ProfileID           string      `json:"profileID"`
	ResponseBody        string      `json:"responseBody"`
	Responses           int         `json:"responses"`
	Sleep               string      `json:"sleep"`
	State               string      `json:"state"`
	Token               string      `json:"token"`
	Topic               string      `json:"topic"`
	TotalRequests       int         `json:"totalRequests"`
	TxID                string      `json:"transactionID"`
	URL                 string      `json:"url"`
	UserLogLevel        string      `json:"userLogLevel"`
	VPToken             string      `json:"vpToken"`
	Workers             int         `json:"workers"`
	Address             string      `json:"address"`
	MessageBrokers      []string    `json:"message-brokers"`
	TotalMessages       int         `json:"total-messages"`
}

func unmarshalLogData(t *testing.T, b []byte) *logData {
	t.Helper()

	l := &logData{}

	require.NoError(t, json.Unmarshal(b, l))

	return l
}
