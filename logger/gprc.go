package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

// GRPCZeroLogger is a GRPC logger that uses zerolog
type GRPCZeroLogger struct {
	log *zerolog.Logger
}

// NewGRPCZeroLogger creates a GRPCZeroLogger
func NewGRPCZeroLogger(log *zerolog.Logger) GRPCZeroLogger {
	grpcLogger := log.With().Str("log-source", "grpc").Logger()
	return GRPCZeroLogger{log: &grpcLogger}
}

// Fatal logs a fatal message
func (l GRPCZeroLogger) Fatal(args ...interface{}) {
	l.log.Fatal().Msg(fmt.Sprint(args...))
}

// Fatalf formats and logs a fatal message
func (l GRPCZeroLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatal().Msg(fmt.Sprintf(format, args...))
}

// Fatalln logs a fatal message and a newline
func (l GRPCZeroLogger) Fatalln(args ...interface{}) {
	l.Fatal(args...)
}

// Error logs an error message
func (l GRPCZeroLogger) Error(args ...interface{}) {
	l.log.Error().Msg(fmt.Sprint(args...))
}

// Errorf formats and logs an error message
func (l GRPCZeroLogger) Errorf(format string, args ...interface{}) {
	l.log.Error().Msg(fmt.Sprintf(format, args...))
}

// Errorln logs an error message and a newline
func (l GRPCZeroLogger) Errorln(args ...interface{}) {
	l.Error(args...)
}

// Info logs an info message
func (l GRPCZeroLogger) Info(args ...interface{}) {
	l.log.Info().Msg(fmt.Sprint(args...))
}

// Infof formats and logs an info message
func (l GRPCZeroLogger) Infof(format string, args ...interface{}) {
	l.log.Info().Msg(fmt.Sprintf(format, args...))
}

// Infoln formats and logs an info message and a newline
func (l GRPCZeroLogger) Infoln(args ...interface{}) {
	l.Info(args...)
}

// Warning logs a warning message
func (l GRPCZeroLogger) Warning(args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprint(args...))
}

// Warningf formats and logs a warning message
func (l GRPCZeroLogger) Warningf(format string, args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprintf(format, args...))
}

// Warningln formats and logs a warning message and a newline
func (l GRPCZeroLogger) Warningln(args ...interface{}) {
	l.Warning(args...)
}

// Print prints a message
func (l GRPCZeroLogger) Print(args ...interface{}) {
	l.log.Print(args...)
}

// Printf formats and prints a message
func (l GRPCZeroLogger) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

// Println prints a message and a newline
func (l GRPCZeroLogger) Println(args ...interface{}) {
	l.Print(args...)
}

// V always returns true
func (l GRPCZeroLogger) V(level int) bool {
	return true
}
