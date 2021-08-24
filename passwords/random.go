package passwords

import (
	"crypto/rand"

	"github.com/pkg/errors"
)

func RandomPassword() (string, error) {
	passBytes := make([]byte, 10)
	_, err := rand.Read(passBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random password")
	}
	return "A5er!0" + string(passBytes), nil
}
