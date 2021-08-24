package logger

import (
	"log"

	"github.com/rs/zerolog"
)

// ZerologWriter implements io.Writer for a zerolog logger
type ZerologWriter struct {
	logger *zerolog.Logger
}

// NewZerologWriter creates a new ZerologWriter
func NewZerologWriter(logger *zerolog.Logger) *ZerologWriter {
	return &ZerologWriter{logger: logger}
}

func (z *ZerologWriter) Write(p []byte) (n int, err error) {
	z.logger.Info().Msg(string(p))
	return len(p), nil
}

// NewSTDLogger creates a standard logger that writes to
// a zerolog logger
func NewSTDLogger(logger *zerolog.Logger) *log.Logger {
	return log.New(NewZerologWriter(logger), "", 0)
}
