package enums

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Email string

func (email Email) IsValid() error {
	if !emailRegex.MatchString(string(email)) {
		return errors.New("Email addres not valid!")
	}

	return nil
}
