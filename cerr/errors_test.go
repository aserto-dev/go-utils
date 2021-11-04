package cerr_test

import (
	"net/http"
	"testing"

	"github.com/aserto-dev/go-utils/cerr"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestDoubleCerr(t *testing.T) {
	assert := require.New(t)

	err := cerr.ErrAccountNotFound.Err(cerr.ErrRepoAlreadyConnected)

	assert.Contains(err.Error(), "already been connected")
	assert.Contains(err.Error(), "account not found")
}

func TestDoubleCerrWithMsg(t *testing.T) {
	assert := require.New(t)

	err := cerr.ErrAccountNotFound.Err(cerr.ErrRepoAlreadyConnected).Msg("failed to setup")

	assert.Contains(err.Error(), "already been connected")
	assert.Contains(err.Error(), "account not found")
}

func TestWithGrpcStatusCode(t *testing.T) {
	assert := require.New(t)
	err := cerr.ErrAccountNotFound.WithGRPCStatus(codes.Canceled)
	assert.Equal(err.StatusCode, codes.Canceled)
}

func TestWithHttpStatusCode(t *testing.T) {
	assert := require.New(t)
	err := cerr.ErrAccountNotFound.WithHTTPStatus(http.StatusAccepted)
	assert.Equal(err.HttpCode, http.StatusAccepted)
}
