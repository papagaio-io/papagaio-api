package dto

import "errors"

type ChangeUserRoleRequestDto struct {
	UserID   *uint64  `json:"id"`
	UserRole UserRole `json:"role"`
}

type UserRole string

const (
	Administrator UserRole = "ADMINISTRATOR"
	Developer     UserRole = "DEVELOPER"
)

func (request *ChangeUserRoleRequestDto) IsValid() error {
	if request.UserID == nil {
		return errors.New("UserId is nil")
	}

	if request.UserRole != Administrator && request.UserRole != Developer {
		return errors.New("UserRole is not valid")
	}

	return nil
}
