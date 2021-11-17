package email

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidEmail(t *testing.T) {
	assert := require.New(t)

	assert.Nil(ValidateFormat("foo@bat.com"))
}

func TestNoAt(t *testing.T) {
	assert := require.New(t)

	assert.NotNil(ValidateFormat("foo.bat.com"))
}

func TestNoDotDomain(t *testing.T) {
	assert := require.New(t)

	assert.NotNil(ValidateFormat("foo@batcom"))
}
