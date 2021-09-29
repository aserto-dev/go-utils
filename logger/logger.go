package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/sirupsen/logrus" //nolint
	"google.golang.org/grpc/grpclog"
)

// Config represents logging configuration
type Config struct {
	Prod           bool          `json:"prod"`
	LogLevelParsed zerolog.Level `json:"-"`
	LogLevel       string        `json:"log_level"`
}

// ParseLogLevel parses the log level in the config and
// sets the appropriate value for `LogLevelParsed`.
func (c *Config) ParseLogLevel(defaultLevel zerolog.Level) error {
	var err error
	if len(c.LogLevel) == 0 {
		c.LogLevelParsed = defaultLevel
		return nil
	}

	c.LogLevelParsed, err = zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		return errors.Wrapf(err, "logging.log_level failed to parse")
	}

	return nil
}

// NewLogger returns a new zap logger
func NewLogger(logOutput io.Writer, cfg *Config) (*zerolog.Logger, error) {
	var logger zerolog.Logger
	if cfg.Prod {
		logger = zerolog.New(logOutput).With().Timestamp().Logger()
	} else {
		cw := zerolog.NewConsoleWriter()
		cw.Out = logOutput
		logger = zerolog.New(cw).With().Timestamp().Logger()
	}
	zerolog.SetGlobalLevel(cfg.LogLevelParsed)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Override standard log output with zerolog
	stdLogger := logger.With().Str("log-source", "std").Logger()
	log.SetOutput(NewZerologWriter(&stdLogger))

	// Override logrus with zerolog
	logrusLogger := logger.With().Str("log-source", "logrus").Logger()
	logrus.AddHook(&logrusHook{logger: &logrusLogger, writer: logOutput})
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	// Override GRPC logging with zerolog
	grpclog.SetLoggerV2(NewGRPCZeroLogger(&logger))

	zerolog.ErrorHandler = func(err error) {
		if !strings.Contains(err.Error(), "file already closed") {
			fmt.Fprintf(os.Stderr, "zerolog: could not write event: %v\n", err)
		}
	}

	return &logger, nil
}
