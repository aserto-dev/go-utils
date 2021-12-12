package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type LevelWriter struct {
	io.Writer
	ErrorWriter io.Writer
}

func (lw *LevelWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	w := lw.Writer
	if l > zerolog.InfoLevel {
		w = lw.ErrorWriter
	}
	return w.Write(p)
}
