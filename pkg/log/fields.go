/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/otel/trace"
)

// Log Fields.
const (
	FieldAddress      = "address"
	FieldDuration     = "duration"
	FieldHTTPStatus   = "httpStatus"
	FieldID           = "id"
	FieldName         = "name"
	FieldPath         = "path"
	FieldResponse     = "response"
	FieldState        = "state"
	FieldToken        = "token"
	FieldTopic        = "topic"
	FieldTxID         = "txID"
	FieldURL          = "url"
	FieldTraceID      = "trace_id"
	FieldSpanID       = "span_id"
	FieldParentSpanID = "parent_span_id"
)

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

// WithTracing adds OpenTelemetry fields, i.e. traceID, spanID, and (optionally) parentSpanID fields.
// If the provided context doesn't contain OpenTelemetry data then the fields are not logged.
func WithTracing(ctx context.Context) zap.Field {
	return zap.Inline(&otelMarshaller{ctx: ctx})
}

// otelMarshaller is an OpenTelemetry marshaller which adds Open-Telemetry
// trace and span IDs (as well as parent span ID if exists) to the log message.
type otelMarshaller struct {
	ctx context.Context
}

type childSpan interface {
	Parent() trace.SpanContext
}

const (
	nilTraceID = "00000000000000000000000000000000"
	nilSpanID  = "0000000000000000"
)

func (m *otelMarshaller) MarshalLogObject(e zapcore.ObjectEncoder) error {
	s := trace.SpanFromContext(m.ctx)

	traceID := s.SpanContext().TraceID().String()
	if traceID == "" || traceID == nilTraceID {
		return nil
	}

	e.AddString(FieldTraceID, traceID)

	spanID := s.SpanContext().SpanID().String()
	if spanID != "" && spanID != nilSpanID {
		e.AddString(FieldSpanID, spanID)
	}

	cspan, ok := s.(childSpan)
	if ok {
		parentSpanID := cspan.Parent().SpanID().String()
		if parentSpanID != "" && parentSpanID != nilSpanID {
			e.AddString(FieldParentSpanID, parentSpanID)
		}
	}

	return nil
}
