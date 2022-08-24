package cerr_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aserto-dev/go-utils/cerr"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestWithEmptyMsg(t *testing.T) {
	assert := require.New(t)

	err := cerr.ErrAccountNotFound.Msg("")

	fields := err.Fields()
	assert.Nil(fields["msg"])

	err = cerr.ErrAccountNotFound.Msg("bla")

	fields = err.Fields()
	assert.NotNil(fields["msg"])
}

func TestError(t *testing.T) {
	assert := require.New(t)

	err := cerr.ErrAccountNotFound.Msg("bla").Err(errors.New("boom"))
	err2 := cerr.ErrAccountNotFound.Msg("bla").Msg("ala")
	err3 := cerr.ErrAccountNotFound.Err(errors.New("boom")).Msg("bla").Msg("ala")
	err4 := cerr.ErrAccountNotFound.Err(errors.New("boom")).Err(errors.New("pow")).Msg("bla").Msg("ala")
	err5 := cerr.ErrAccountNotFound.Err(errors.New("boom"))
	err6 := cerr.ErrAccountNotFound.Err(errors.New("boom")).Err(errors.New("pow"))
	err7 := cerr.ErrAccountNotFound.Msg("bla")

	assert.Equal(err.Error(), "E10012 account not found: boom: bla")
	assert.Equal(err2.Error(), "E10012 account not found: bla: ala")
	assert.Equal(err3.Error(), "E10012 account not found: boom: bla: ala")
	assert.Equal(err4.Error(), "E10012 account not found: boom: pow: bla: ala")
	assert.Equal(err5.Error(), "E10012 account not found: boom")
	assert.Equal(err6.Error(), "E10012 account not found: boom: pow")
	assert.Equal(err7.Error(), "E10012 account not found: bla")
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

func TestFromGRPCStatus(t *testing.T) {
	assert := require.New(t)

	initialErr := cerr.ErrUserNotFound
	initialErr = initialErr.Str("email", "testuser@aserto.com").Msg("foo")

	grpcStatus := status.New(initialErr.StatusCode, initialErr.Error())
	grpcStatus, err := grpcStatus.WithDetails(&errdetails.ErrorInfo{
		Reason:   "1234",
		Metadata: initialErr.Data(),
		Domain:   initialErr.Code,
	})

	if err != nil {
		assert.Fail(err.Error())
	}

	transformedErr := cerr.FromGRPCStatus(*grpcStatus)

	assert.True(initialErr.SameAs(transformedErr))

	assert.Equal(initialErr.Error(), transformedErr.Error())
	assert.Equal(initialErr.Message, transformedErr.Message)
}

func TestUnwrapNilErr(t *testing.T) {
	assert := require.New(t)

	err := cerr.UnwrapAsertoError(nil)

	assert.Nil(err)
}

func TestEquals(t *testing.T) {
	assert := require.New(t)

	err1 := cerr.ErrConnection.Msgf("error 1").Str("key1", "val1").Err(errors.New("boom"))
	err2 := cerr.ErrConnection.Msgf("error 2").Str("key2", "val2").Err(errors.New("zoom"))

	assert.True(cerr.Equals(err1, err2))

}

func TestEqualsNil(t *testing.T) {
	assert := require.New(t)

	assert.True(cerr.Equals(nil, nil))
}

func TestEqualsOneNil(t *testing.T) {
	assert := require.New(t)

	assert.False(cerr.Equals(cerr.ErrAccountNotFound, nil))
}

func TestEqualsNormalErrorOneNil(t *testing.T) {
	assert := require.New(t)

	assert.False(cerr.Equals(errors.New("boom"), nil))
}

func TestEqualsErrCerr(t *testing.T) {
	assert := require.New(t)

	assert.False(cerr.Equals(errors.New("boom"), cerr.ErrAccountNotFound))
}

func TestEqualsFalse(t *testing.T) {
	assert := require.New(t)

	assert.False(cerr.Equals(cerr.ErrAlreadyMember, cerr.ErrAccountNotFound))
}

func TestEqualsNormalErrors(t *testing.T) {
	assert := require.New(t)

	assert.False(cerr.Equals(errors.New("boom1"), errors.New("boom2")))
}

func TestCodeToAsertoError(t *testing.T) {
	assert := require.New(t)

	asertoErr := cerr.CodeToAsertoError("E10009")

	assert.NotNil(asertoErr)
	assert.True(cerr.Equals(asertoErr, cerr.ErrGithubAccessToken))
}

func TestCodeToAsertoErrorInvalidCode(t *testing.T) {
	assert := require.New(t)

	asertoErr := cerr.CodeToAsertoError("E20009")

	assert.Nil(asertoErr)
}
