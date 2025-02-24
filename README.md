[![Release](https://img.shields.io/github/release/trustbloc/logutil-go.svg?style=flat-square)](https://github.com/trustbloc/logutil-go/releases/latest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/trustbloc/logutil-go/main/LICENSE)
[![Godocs](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/trustbloc/logutil-go)

[![Build Status](https://github.com/trustbloc/logutil-go/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/trustbloc/logutil-go/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/trustbloc/logutil-go/branch/main/graph/badge.svg)](https://codecov.io/gh/trustbloc/logutil-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/trustbloc/logutil-go)](https://goreportcard.com/report/github.com/trustbloc/logutil-go)

# logutil-go

General purpose field enabled logging module. This allows logs to be attribute labelled instead of having long string logs.

Following are the functions available in the logger:

- **Debug** - Logs a message at the debug level.
- **Debugc** - Logs a message at the debug level with the given context.
- **Info** - Logs a message at the info level.
- **Infoc** - Logs a message at the info level with the given context.
- **Warn** - Logs a message at the warn level.
- **Warnc** - Logs a message at the warn level with the given context.
- **Error** - Logs a message at the error level.
- **Errorc** - Logs a message at the error level with the given context.

The methods that accept a context parameter are used to log trace information.  The trace information is extracted from the context and logged as part of the log message.

Following are the additional fields that are logged when trace information is found in the provided context:

- **trace_id**: <traceID> (from [OTel](https://www.w3.org/TR/trace-context/))
- **span_id**: <spanID> (from [OTel](https://www.w3.org/TR/trace-context/))
- **parent_span_id**: <parentSpanID> (from [OTel](https://www.w3.org/TR/trace-context/))
- **correlation_id**: <correlationID> (from [Bagage](https://www.w3.org/TR/baggage/))

For example:

```
{"level":"debug","ts":"2024-11-04T19:37:32.844Z","logger":"controller","caller":"controller.go:63","msg":"Received request","trace_id":"b20283308c97befd8606ab8932e1d476","span_id":"f32eec4232b5d3e4","parent_span_id":"221262d93c002aab","correlation_id":"2A1E11A0"}
```

## Correlation ID

The correlation ID is used to correlate logs across services. The correlation ID is passed in the request header and is propagated to all the services that are called as part of the request. The correlation ID is logged as part of the log message. The following functions are available to work with the correlation ID:

- **FromContext** returns the correlation ID from the given context. A correlation ID is searched in the [Bagage](https://www.w3.org/TR/baggage/) member, 'X-Correlation-Id',
 If a correlation ID is not found in the baggage then:
    - If _GenerateUUIDIfNotFound_ option is set, a new UUID is generated and set on the returned context.
    - If _GenerateNewFixedLengthIfNotFound_ option is set, a new fixed-length correlation ID
      is generated and set on the returned context.
    - If _WithValue_ is set then the given correlation ID is set on the returned context.
    - If none of the above options is specified then the existing context and empty string are returned.
- **correlationid.HTTPTransport** is a RoundTripper that sets the X-Correlation-Id request header for outgoing requests.
- **correlationidecho.Middleware** is middleware for the Echo HTTP server that extracts the X-Correlation-Id request header and sets it in the request context Baggage.
- **correlationidmux.Middleware** is middleware for the Gorilla Mux HTTP server that extracts the X-Correlation-Id request header and sets it in the request context Baggage.

## Enabling OTel tracing, including Baggage

In order to ensure tracing information is propogated across services, the following steps are required:

``` go
// Propagate trace context via traceparent and tracestate headers (https://www.w3.org/TR/trace-context/)
// and baggage items (https://www.w3.org/TR/baggage/).
otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
```
