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
	FieldAddress    = "address"
	FieldDuration   = "duration"
	FieldHTTPStatus = "httpStatus"
	FieldID         = "id"
	FieldName       = "name"
	FieldPath       = "path"
	FieldResponse   = "response"
	FieldState      = "state"
	FieldToken      = "token"
	FieldTopic      = "topic"
	FieldTxID       = "txID"
	FieldURL        = "url"
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

// WithHTTPStatus sets the http-status field.
func WithHTTPStatus(value int) zap.Field {
	return zap.Int(FieldHTTPStatus, value)
}

// WithTopic sets the topic field.
func WithTopic(value string) zap.Field {
	return zap.String(FieldTopic, value)
}

// WithToken sets the token field.
func WithToken(token string) zap.Field {
	return zap.String(FieldToken, token)
}

// WithResponse sets the response field.
func WithResponse(value []byte) zap.Field {
	return zap.String(FieldResponse, string(value))
}

// WithPath sets the path field.
func WithPath(path string) zap.Field {
	return zap.String(FieldPath, path)
}

// WithURL sets the url field.
func WithURL(url string) zap.Field {
	return zap.String(FieldURL, url)
}

// WithDuration sets the duration field.
func WithDuration(value time.Duration) zap.Field {
	return zap.Duration(FieldDuration, value)
}

// WithTxID sets the transaction id field.
func WithTxID(txID string) zap.Field {
	return zap.String(FieldTxID, txID)
}

// WithState sets the state field.
func WithState(state string) zap.Field {
	return zap.String(FieldState, state)
}

// WithAddress sets the address field.
func WithAddress(address string) zap.Field {
	return zap.String(FieldAddress, address)
}
