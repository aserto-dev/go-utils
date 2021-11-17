package email

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// recommended by Auth0: https://community.auth0.com/t/email-regex-verification/28334
	emailRegexp = regexp.MustCompile(`(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$)`)
)

func ValidateFormat(email string) error {

	if !emailRegexp.MatchString(strings.ToLower(email)) {
		return errors.New("invalid email format")
	}
	return nil
}
