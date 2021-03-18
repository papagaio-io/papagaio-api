package model

import (
	"errors"
	"regexp"
)

type User struct {
	Email string `json:"email"`
	//Role  string `json:"role"`
}

const constMailRegex = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var emailRegex = regexp.MustCompile(constMailRegex)

func (user User) IsValid() error {
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
