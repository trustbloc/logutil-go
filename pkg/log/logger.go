/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package log

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	timestampKey  = "ts"
	levelKey      = "level"
	moduleKey     = "logger"
	callerKey     = "caller"
	messageKey    = "msg"
	stacktraceKey = "stacktrace"
)

// DefaultEncoding sets the default logger encoding.
// It may be overridden at build time using the -ldflags option.
var DefaultEncoding = JSON //nolint gochecknoglobals

// Level defines a log level for logging messages.
type Level int

// String returns string representation of given log level.
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARN"
	case ERROR:
		return "ERROR"
	case PANIC:
		return "PANIC"
	case FATAL:
		return "FATAL"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

// ParseLevel returns the level from the given string.
func ParseLevel(level string) (Level, error) {
	switch level {
	case "DEBUG", "debug":
		return DEBUG, nil
	case "INFO", "info":
		return INFO, nil
	case "WARN", "warn", "WARNING", "warning":
		return WARNING, nil
	case "ERROR", "error":
		return ERROR, nil
	case "PANIC", "panic":
		return PANIC, nil
	case "FATAL", "fatal":
		return FATAL, nil
	default:
		return ERROR, errors.New("logger: invalid log level")
	}
}

// Log levels.
const (
	DEBUG   = Level(zapcore.DebugLevel)
	INFO    = Level(zapcore.InfoLevel)
	WARNING = Level(zapcore.WarnLevel)
	ERROR   = Level(zapcore.ErrorLevel)
	PANIC   = Level(zapcore.PanicLevel)
	FATAL   = Level(zapcore.FatalLevel)

	minLogLevel  = DEBUG
	defaultLevel = INFO
)

var levels = newModuleLevels() //nolint: gochecknoglobals

type options struct {
	encoding   Encoding
	stdOut     zapcore.WriteSyncer
	stdErr     zapcore.WriteSyncer
	fields     []zap.Field
	callerSkip int
}

// Encoding defines the log encoding.
type Encoding = string

// Log encodings.
const (
	Console Encoding = "console"
	JSON    Encoding = "json"
)

const defaultModuleName = ""

// Option is a logger option.
type Option func(o *options)

// WithStdOut sets the output for logs of type DEBUG, INFO, and WARN.
func WithStdOut(stdOut zapcore.WriteSyncer) Option {
	return func(o *options) {
		o.stdOut = stdOut
	}
}

// WithStdErr sets the output for logs of type ERROR, PANIC, and FATAL.
func WithStdErr(stdErr zapcore.WriteSyncer) Option {
	return func(o *options) {
		o.stdErr = stdErr
	}
}

// WithFields sets the fields that will be output with every log.
func WithFields(fields ...zap.Field) Option {
	return func(o *options) {
		o.fields = fields
	}
}

// WithEncoding sets the output encoding (console or json).
func WithEncoding(encoding Encoding) Option {
	return func(o *options) {
		o.encoding = encoding
	}
}

// WithCallerSkip sets the caller skip for ctx logger.
func WithCallerSkip(callerSkip int) Option {
	return func(o *options) {
		o.callerSkip = callerSkip
	}
}

// Log uses the Zap Logger to log messages in a structured way. Functions are also included to
// log context-specific fields, such as OpenTelemetry trace and span IDs.
type Log struct {
	*zap.Logger
	ctxLogger *zap.Logger
	module    string
}

// New creates a Zap Logger to log messages in a structured way.
func New(module string, opts ...Option) *Log {
	options := getOptions(opts)

	return &Log{
		Logger: newZap(module, options.encoding, options.stdOut, options.stdErr).
			With(options.fields...),
		ctxLogger: newZap(module, options.encoding, options.stdOut, options.stdErr).
			WithOptions(zap.AddCallerSkip(options.callerSkip)).
			With(options.fields...),
		module: module,
	}
}

// IsEnabled returns true if given log level is enabled.
func (l *Log) IsEnabled(level Level) bool {
	return levels.isEnabled(l.module, level)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Log) With(fields ...zap.Field) *Log {
	if len(fields) == 0 {
		return l
	}

	return &Log{
		Logger:    l.Logger.With(fields...),
		ctxLogger: l.ctxLogger.With(fields...),
		module:    l.module,
	}
}

// Debugc logs a message at Debug level, including the provided fields and any implicit context
// fields (such as OpenTelemetry trace ID and span ID).
func (l *Log) Debugc(ctx context.Context, msg string, fields ...zap.Field) {
	l.ctxLogger.Debug(msg, append(fields, WithTracing(ctx))...)
}

// Infoc logs a message at Info level, including the provided fields and any implicit context
// fields (such as OpenTelemetry trace ID and span ID).
func (l *Log) Infoc(ctx context.Context, msg string, fields ...zap.Field) {
	l.ctxLogger.Info(msg, append(fields, WithTracing(ctx))...)
}

// Warnc logs a message at Warning level, including the provided fields and any implicit context
// fields (such as OpenTelemetry trace ID and span ID).
func (l *Log) Warnc(ctx context.Context, msg string, fields ...zap.Field) {
	l.ctxLogger.Warn(msg, append(fields, WithTracing(ctx))...)
}

// Errorc logs a message at Error level, including the provided fields and any implicit context
// fields (such as OpenTelemetry trace ID and span ID).
func (l *Log) Errorc(ctx context.Context, msg string, fields ...zap.Field) {
	l.ctxLogger.Error(msg, append(fields, WithTracing(ctx))...)
}

// Panicc logs a message at Panic level, including the provided fields and any implicit context
// fields (such as OpenTelemetry trace ID and span ID).
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Log) Panicc(ctx context.Context, msg string, fields ...zap.Field) {
	l.ctxLogger.Panic(msg, append(fields, WithTracing(ctx))...)
}

// Fatalc logs a message at Fatal level, including the provided fields and any implicit context
// fields (such as OpenTelemetry trace ID and span ID).
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Log) Fatalc(ctx context.Context, msg string, fields ...zap.Field) {
	l.ctxLogger.Fatal(msg, append(fields, WithTracing(ctx))...)
}

// SetLevel sets the log level for given module and level.
func SetLevel(module string, level Level) {
	levels.Set(module, level)
}

// SetDefaultLevel sets the default log level.
func SetDefaultLevel(level Level) {
	levels.SetDefault(level)
}

// GetLevel returns the log level for the given module.
func GetLevel(module string) Level {
	return levels.Get(module)
}

// SetSpec sets the log levels for individual modules as well as the default log level.
// The format of the spec is as follows:
//
//	module1=level1:module2=level2:module3=level3:defaultLevel
//
// Valid log levels are: critical, error, warning, info, debug
//
// Example:
//
//	module1=error:module2=debug:module3=warning:info
func SetSpec(spec string) error {
	logLevelByModule := strings.Split(spec, ":")

	defaultLogLevel := minLogLevel - 1

	var moduleLevelPairs []moduleLevelPair

	for _, logLevelByModulePart := range logLevelByModule {
		if strings.Contains(logLevelByModulePart, "=") {
			moduleAndLevelPair := strings.Split(logLevelByModulePart, "=")

			logLevel, err := ParseLevel(moduleAndLevelPair[1])
			if err != nil {
				return err
			}

			moduleLevelPairs = append(moduleLevelPairs,
				moduleLevelPair{moduleAndLevelPair[0], logLevel})
		} else {
			if defaultLogLevel >= minLogLevel {
				return errors.New("multiple default values found")
			}

			level, err := ParseLevel(logLevelByModulePart)
			if err != nil {
				return err
			}

			defaultLogLevel = level
		}
	}

	if defaultLogLevel >= minLogLevel {
		levels.Set("", defaultLogLevel)
	} else {
		levels.Set("", INFO)
	}

	for _, moduleLevelPair := range moduleLevelPairs {
		levels.Set(moduleLevelPair.module, moduleLevelPair.logLevel)
	}

	return nil
}

// GetSpec returns the log spec which specifies the log level of each individual module. The spec is
// in the following format:
//
//	module1=level1:module2=level2:module3=level3:defaultLevel
//
// Example:
//
//	module1=error:module2=debug:module3=warning:info
func GetSpec() string {
	var spec string

	var defaultDebugLevel string

	for module, level := range getAllLevels() {
		if module == "" {
			defaultDebugLevel = level.String()
		} else {
			spec += fmt.Sprintf("%s=%s:", module, level.String())
		}
	}

	return spec + defaultDebugLevel
}

func getAllLevels() map[string]Level {
	metadataLevels := levels.All()

	// Convert to the Level type in this package
	levels := make(map[string]Level)
	for module, logLevel := range metadataLevels {
		levels[module] = logLevel
	}

	return levels
}

type moduleLevelPair struct {
	module   string
	logLevel Level
}

func newModuleLevels() *moduleLevels {
	return &moduleLevels{levels: make(map[string]Level)}
}

// moduleLevels maintains log levels based on modules.
type moduleLevels struct {
	levels  map[string]Level
	rwmutex sync.RWMutex
}

// Get returns the log level for given module and level.
func (l *moduleLevels) Get(module string) Level {
	l.rwmutex.RLock()
	defer l.rwmutex.RUnlock()

	level, exists := l.levels[module]
	if !exists {
		level, exists = l.levels[defaultModuleName]
		// no configuration exists, default to info
		if !exists {
			return defaultLevel
		}
	}

	return level
}

// All returns all set log levels.
func (l *moduleLevels) All() map[string]Level {
	l.rwmutex.RLock()
	levels := l.levels
	l.rwmutex.RUnlock()

	levelsCopy := make(map[string]Level)

	for module, logLevel := range levels {
		levelsCopy[module] = logLevel
	}

	return levelsCopy
}

func (l *moduleLevels) Set(module string, level Level) {
	l.rwmutex.Lock()
	l.levels[module] = level
	l.rwmutex.Unlock()
}

func (l *moduleLevels) SetDefault(level Level) {
	l.Set(defaultModuleName, level)
}

// isEnabled will return true if logging is enabled for given module and level.
func (l *moduleLevels) isEnabled(module string, level Level) bool {
	return level >= l.Get(module)
}

func newZap(module string, encoding Encoding, stdOut, stdErr zapcore.WriteSyncer) *zap.Logger {
	encoder := newZapEncoder(encoding)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.Lock(stdErr),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel && levels.isEnabled(module, Level(lvl))
			}),
		),
		zapcore.NewCore(encoder, zapcore.Lock(stdOut),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl < zapcore.ErrorLevel && levels.isEnabled(module, Level(lvl))
			}),
		),
	)

	return zap.New(core, zap.AddCaller()).Named(module)
}

func newZapEncoder(encoding Encoding) zapcore.Encoder {
	defaultCfg := zapcore.EncoderConfig{
		TimeKey:        timestampKey,
		LevelKey:       levelKey,
		NameKey:        moduleKey,
		CallerKey:      callerKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     messageKey,
		StacktraceKey:  stacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	switch strings.ToLower(encoding) {
	case JSON:
		cfg := defaultCfg
		cfg.EncodeLevel = zapcore.LowercaseLevelEncoder

		return zapcore.NewJSONEncoder(cfg)
	case Console:
		cfg := defaultCfg
		cfg.EncodeName = func(moduleName string, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(fmt.Sprintf("[%s]", moduleName))
		}

		return zapcore.NewConsoleEncoder(cfg)
	default:
		panic("unsupported encoding " + encoding)
	}
}

func getOptions(opts []Option) *options {
	options := &options{
		encoding:   DefaultEncoding,
		stdOut:     os.Stdout,
		stdErr:     os.Stderr,
		callerSkip: 1,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
