/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log Fields.
const (
	FieldAddress       = "address"
	FieldCommand       = "command"
	FieldDuration      = "duration"
	FieldHTTPStatus    = "httpStatus"
	FieldHostURL       = "hostURL"
	FieldID            = "id"
	FieldIDToken       = "idToken"
	FieldJSON          = "json"
	FieldName          = "name"
	FieldParameter     = "parameter"
	FieldParameters    = "parameters"
	FieldPath          = "path"
	FieldProfileID     = "profileID"
	FieldResponse      = "response"
	FieldResponseBody  = "responseBody"
	FieldResponses     = "responses"
	FieldService       = "service"
	FieldSleep         = "sleep"
	FieldState         = "state"
	FieldToken         = "token"
	FieldTopic         = "topic"
	FieldTotalMessages = "total-messages"
	FieldTotalRequests = "totalRequests"
	FieldTxID          = "transactionID"
	FieldURL           = "url"
)

// ObjectMarshaller uses reflection to marshal an object's fields.
type ObjectMarshaller struct {
	key string
	obj interface{}
}

// NewObjectMarshaller returns a new ObjectMarshaller.
func NewObjectMarshaller(key string, obj interface{}) *ObjectMarshaller {
	return &ObjectMarshaller{key: key, obj: obj}
}

// MarshalLogObject marshals the object's fields.
func (m *ObjectMarshaller) MarshalLogObject(e zapcore.ObjectEncoder) error {
	return e.AddReflected(m.key, m.obj)
}

// WithError sets the error field.
func WithError(err error) zap.Field {
	return zap.Error(err)
}

// WithID sets the id field.
func WithID(id string) zap.Field {
	return zap.String(FieldID, id)
}

// WithName sets the name field.
func WithName(name string) zap.Field {
	return zap.String(FieldName, name)
}

// WithCommand sets the command field.
func WithCommand(command string) zap.Field {
	return zap.String(FieldCommand, command)
}

// WithParameter sets the parameter field.
func WithParameter(value string) zap.Field {
	return zap.String(FieldParameter, value)
}

// WithParameters sets the parameters field.
func WithParameters(value interface{}) zap.Field {
	return zap.Inline(NewObjectMarshaller(FieldParameters, value))
}

// WithHTTPStatus sets the http-status field.
func WithHTTPStatus(value int) zap.Field {
	return zap.Int(FieldHTTPStatus, value)
}

// WithResponseBody sets the response body field.
func WithResponseBody(value []byte) zap.Field {
	return zap.String(FieldResponseBody, string(value))
}

// WithTopic sets the topic field.
func WithTopic(value string) zap.Field {
	return zap.String(FieldTopic, value)
}

// WithHostURL sets the hostURL field.
func WithHostURL(hostURL string) zap.Field {
	return zap.String(FieldHostURL, hostURL)
}

// WithToken sets the token field.
func WithToken(token string) zap.Field {
	return zap.String(FieldToken, token)
}

// WithTotalRequests sets the total requests field.
func WithTotalRequests(totalRequests int) zap.Field {
	return zap.Int(FieldTotalRequests, totalRequests)
}

// WithResponse sets the response field.
func WithResponse(value []byte) zap.Field {
	return zap.String(FieldResponse, string(value))
}

// WithResponses sets the responses field.
func WithResponses(responses int) zap.Field {
	return zap.Int(FieldResponses, responses)
}

// WithPath sets the path field.
func WithPath(path string) zap.Field {
	return zap.String(FieldPath, path)
}

// WithURL sets the url field.
func WithURL(url string) zap.Field {
	return zap.String(FieldURL, url)
}

// WithJSON sets the json field.
func WithJSON(json string) zap.Field {
	return zap.String(FieldJSON, json)
}

// WithSleep sets the sleep field.
func WithSleep(sleep time.Duration) zap.Field {
	return zap.Duration(FieldSleep, sleep)
}

// WithDuration sets the duration field.
func WithDuration(value time.Duration) zap.Field {
	return zap.Duration(FieldDuration, value)
}

// WithIDToken sets the id token field.
func WithIDToken(idToken string) zap.Field {
	return zap.String(FieldIDToken, idToken)
}

// WithTxID sets the transaction id field.
func WithTxID(txID string) zap.Field {
	return zap.String(FieldTxID, txID)
}

// WithService sets the service field.
func WithService(value string) zap.Field {
	return zap.String(FieldService, value)
}

// WithState sets the state field.
func WithState(state string) zap.Field {
	return zap.String(FieldState, state)
}

// WithProfileID sets the presentation definition id field.
func WithProfileID(id string) zap.Field {
	return zap.String(FieldProfileID, id)
}

// WithAddress sets the address field.
func WithAddress(address string) zap.Field {
	return zap.String(FieldAddress, address)
}

// WithTotalMessages sets the total messages field.
func WithTotalMessages(totalMessages int) zap.Field {
	return zap.Int(FieldTotalMessages, totalMessages)
}
