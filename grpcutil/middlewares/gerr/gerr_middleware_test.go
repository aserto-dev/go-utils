package gerr

import (
	"bytes"
	"testing"

	"github.com/aserto-dev/go-utils/cerr"
	"github.com/aserto-dev/go-utils/grpcutil/middlewares/test"
	"github.com/aserto-dev/go-utils/logger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func TestUnaryServerWithWrappedError(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler("output", errors.Wrap(cerr.ErrAlreadyMember, "unimportant error"))

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)
	_, err := NewErrorMiddleware().Unary()(ctx, "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)
	assert.Contains(err.Error(), "already a tenant member")
}

func TestUnaryServerWithFields(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler(
		"output",
		errors.Wrap(cerr.ErrAlreadyMember.Str("my-field", "deadbeef"), "another error"),
	)

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	logOutput := buf.String()

	assert.Contains(logOutput, "deadbeef")
}

func TestUnaryServerWithDoubleCerr(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler(
		"output",
		cerr.ErrSCC.Err(cerr.ErrAlreadyMember.Str("my-field", "deadbeef").Msg("old message")).Msg("new message"),
	)

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	logOutput := buf.String()

	assert.Contains(logOutput, "new message")
	assert.Contains(logOutput, "deadbeef")
}

func TestSimpleInnerError(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler("output", cerr.ErrSCC.Err(errors.New("deadbeef")).Msg("failed to setup initial tag"))

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	logOutput := buf.String()

	assert.Contains(logOutput, "deadbeef")
}

func TestDirectResult(t *testing.T) {
	assert := require.New(t)
	handler := test.NewHandler(
		"output",
		cerr.ErrSCC.Err(cerr.ErrRepoAlreadyConnected).Msg("failed to setup initial tag"),
	)

	buf := bytes.NewBufferString("")
	testLogger := logger.TestLogger(buf)

	ctx := grpc.NewContextWithServerTransportStream(
		test.RequestIDContext(t),
		test.ServerTransportStream(""),
	)

	_, err := NewErrorMiddleware().Unary()(testLogger.WithContext(ctx), "xyz", test.UnaryInfo, handler.Unary)
	assert.Error(err)

	s := status.Convert(err)

	errDetailsFound := false
	for _, detail := range s.Details() {
		switch t := detail.(type) {
		case *errdetails.ErrorInfo:
			errDetailsFound = true
			assert.Contains(t.Metadata, "msg")
			assert.Contains(t.Metadata["msg"], "failed to setup")
		}
	}

	assert.True(errDetailsFound)
	assert.Contains(s.Message(), "there was an error interacting")
	assert.Contains(err.Error(), "already been")
}
