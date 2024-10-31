/*
Copyright Gen Digital Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

const (
	// CorrelationIDHeader is the HTTP header key for the correlation ID.
	CorrelationIDHeader = "X-Correlation-ID"

	// CorrelationIDAttribute is the Open Telemetry span attribute key for the correlation ID.
	CorrelationIDAttribute = "dts.correlation_id"
)
