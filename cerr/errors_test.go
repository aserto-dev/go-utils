package cerr_test

import (
	"testing"

	"github.com/aserto-dev/go-utils/cerr"
	"github.com/stretchr/testify/require"
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
