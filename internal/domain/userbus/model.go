package userbus

import (
	"github.com/google/uuid"
	"net/mail"
	"time"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	Roles        []Role
	PasswordHash []byte
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

type NewUser struct {
	Name     string
	Email    mail.Address
	Roles    []Role
	Password string
}

type UpdateUser struct {
	Name     *string
	Email    *mail.Address
	Roles    []Role
	Password *string
	Enabled  *bool
}
